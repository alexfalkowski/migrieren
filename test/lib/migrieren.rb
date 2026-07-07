# frozen_string_literal: true

require 'securerandom'
require 'yaml'
require 'base64'
require 'open3'
require 'timeout'

require 'pg'

require 'migrieren/pg'
require 'migrieren/v1/http'
require 'migrieren/v1/service_services_pb'

##
# Public entrypoints for the Ruby feature-test harness.
#
# This module is used by Cucumber step definitions under `test/features/**` to
# talk to a locally-running Migrieren server over both HTTP and gRPC.
#
# The methods exposed here are intentionally convenience wrappers that:
# - provide stable default endpoints (`http://localhost:11000` and
#   `localhost:12000`),
# - memoize clients/stubs for reuse across steps, and
# - centralize shared configuration such as gRPC user-agent headers.
#
# @example Using the HTTP façade client
#   Migrieren::V1.server_http.migrate('postgres', 1)
#
# @example Using the gRPC API stub
#   req = Migrieren::V1::MigrateRequest.new(database: 'postgres', version: 1)
#   Migrieren::V1.server_grpc.migrate(req)
#
module Migrieren
  class << self
    ##
    # Returns an observability client for feature-harness code.
    #
    # This client is provided by the `nonnative` test utilities and is available
    # to feature-harness code that needs to query or assert on telemetry emitted
    # by the service.
    #
    # The client points at the service's HTTP endpoint.
    #
    # @return [Nonnative::Observability] a memoized observability client
    def observability
      @observability ||= Nonnative::Observability.new('http://localhost:11000')
    end

    ##
    # Returns the parsed server configuration for the test environment.
    #
    # The harness loads `.config/server.yml` via `Nonnative::ConfigurationFile.load` and
    # memoizes the resulting configuration object.
    #
    # @return [Object] the configuration structure returned by `nonnative`
    def server_config
      @server_config ||= Nonnative::ConfigurationFile.load('.config/server.yml')
    end

    ##
    # Returns a gRPC Health Check stub connected to the test server.
    #
    # This is used by health-related steps to query standard gRPC health
    # endpoints.
    #
    # @return [Nonnative::GRPCHealth] a memoized gRPC health client
    def health_grpc
      @health_grpc ||= Nonnative.grpc_health(host: 'localhost', port: 12_000, service: 'migrieren.v1.Service')
    end

    ##
    # Returns a helper for manipulating the Postgres test database directly.
    #
    # The helper uses the `pg` gem to connect to the local test database and is
    # typically used for setup/teardown chores (for example dropping tables).
    #
    # @return [Migrieren::PG] a memoized Postgres helper
    def pg
      @pg ||= Migrieren::PG.new
    end

    ##
    # Returns gRPC channel arguments that set an explicit user-agent.
    #
    # This allows the test harness to assert on and/or preserve a consistent
    # gRPC user-agent header across calls.
    #
    # @return [Hash] gRPC channel arguments compatible with `grpc` Ruby stubs
    def user_agent
      @user_agent ||= Nonnative::Header.grpc_user_agent('Migrieren-ruby-client/1.0 gRPC/1.0')
    end

    ##
    # Returns bounded per-call options for gRPC feature-harness requests.
    #
    # The default deadline is slightly longer than the service transport
    # timeout in `.config/server.yml`, so ordinary requests can finish while a
    # stalled endpoint still fails before an outer Cucumber or CI timeout.
    #
    # Each call includes a generated `request-id` metadata value. Caller-provided
    # metadata is merged afterward, so scenarios can override that value when
    # they need a specific request identifier.
    #
    # @param metadata [Hash] request metadata merged after the generated default
    # @param deadline [Time, nil] optional deadline override for scenarios that
    #   intentionally exercise request cancellation
    # @return [Hash] options compatible with Ruby gRPC unary calls
    def grpc_options(metadata: {}, deadline: nil)
      {
        metadata: { 'request-id' => SecureRandom.uuid }.merge(metadata),
        deadline: deadline || (Time.now + 6)
      }
    end

    ##
    # Returns bounded per-call options for HTTP feature-harness requests.
    #
    # Each call includes a generated `request_id` header. Caller-provided
    # headers are merged afterward, so scenarios can override that value or add
    # transport-specific headers such as content type and user agent.
    #
    # @param headers [Hash] HTTP headers merged after the generated request id
    # @param read_timeout [Integer] read timeout in seconds
    # @param open_timeout [Integer] connection open timeout in seconds
    # @return [Hash] options compatible with `Nonnative::HTTPClient` calls
    def http_options(headers: {}, read_timeout: 10, open_timeout: 10)
      {
        headers: { request_id: SecureRandom.uuid }.merge(headers),
        read_timeout:,
        open_timeout:
      }
    end

    ##
    # Lifetime, in seconds, of tokens minted by the feature harness.
    #
    # It is comfortably below the server's configured token expiration
    # (`transport.*.token.ssh.exp` in `.config/server.yml`) so generated tokens
    # never exceed the verifier's signed-lifetime cap.
    TOKEN_EXPIRATION = 300

    ##
    # Returns a memoized `nonnative` SSH token generator for a signing key.
    #
    # The service verifies go-service SSH tokens (`transport.*.token.ssh`), and
    # `nonnative` mints matching tokens. SSH tokens fix `sub == kid == key`, so
    # the key id becomes the verified user id the access policy is evaluated
    # against. Two keys exist under `secrets/`: `migrieren` (granted by the
    # Casbin policy) and `guest` (verifiable but not granted).
    #
    # @param key [String] the signing key id, matching a `secrets/ssh_<key>`
    #   OpenSSH private key and a server-side public key
    # @return [Nonnative::Token] a memoized token generator for that key
    def auth_token(key = 'migrieren')
      (@auth_tokens ||= {})[key] ||=
        Nonnative.token(kind: 'ssh', issuer: 'migrieren', key:, private_key: "secrets/ssh_#{key}", expiration: TOKEN_EXPIRATION)
    end

    ##
    # Builds a `Bearer` Authorization header value for an HTTP RPC route.
    #
    # The audience is bound to the route (`"POST <path>"`) so the token cannot be
    # replayed against a different endpoint.
    #
    # @param path [String] the HTTP RPC path, for example `/migrieren.v1.Service/Status`
    # @param key [String] the signing key id (see {auth_token})
    # @return [String] an Authorization header value such as `"Bearer <token>"`
    def http_authorization(path, key = 'migrieren')
      Nonnative::Header.auth_bearer(auth_token(key).generate(aud: Nonnative::Token.http_audience('POST', path), sub: key))[:authorization]
    end

    ##
    # Builds a `Bearer` Authorization metadata value for a gRPC method.
    #
    # The audience is bound to the gRPC full method so the token is scoped to a
    # single RPC.
    #
    # @param full_method [String] the gRPC full method, for example `/migrieren.v1.Service/Status`
    # @param key [String] the signing key id (see {auth_token})
    # @return [String] an Authorization metadata value such as `"Bearer <token>"`
    def grpc_authorization(full_method, key = 'migrieren')
      "Bearer #{auth_token(key).generate(aud: Nonnative::Token.grpc_audience(full_method.to_s), sub: key)}"
    end

    ##
    # Returns HTTP request options with a route-scoped Authorization header.
    #
    # This is the default authentication path used by {Migrieren::V1::HTTP}: it
    # mints a `migrieren` token bound to `"POST <path>"` unless the caller already
    # set an `:authorization` header. Scenarios exercising rejection pass an
    # explicit header (empty, malformed, or `guest`-signed) to opt out.
    #
    # @param path [String] the HTTP RPC path being called
    # @param opts [Hash] request options passed to `Nonnative::HTTPClient#post`
    # @return [Hash] the options with an Authorization header when one was absent
    def authorize_http(path, opts)
      headers = opts[:headers] || {}
      return opts if headers.key?(:authorization)

      opts.merge(headers: headers.merge(authorization: http_authorization(path)))
    end
  end

  ##
  # gRPC client interceptor that attaches a route-scoped Bearer token.
  #
  # It authenticates every unary call made through {Migrieren::V1.server_grpc} by
  # setting an `authorization` metadata entry scoped to the call's full method,
  # unless the call already supplies one. Scenarios exercising rejection pass an
  # explicit `authorization` metadata value (empty, malformed, or `guest`-signed)
  # to opt out.
  class GRPCAuthorization < GRPC::ClientInterceptor
    ##
    # Injects the Authorization metadata for a unary request/response call.
    #
    # @param method [String] the gRPC full method being invoked
    # @param metadata [Hash] the mutable per-call metadata
    # @return [Object] the result of the intercepted call
    def request_response(method:, metadata:, **)
      metadata['authorization'] ||= Migrieren.grpc_authorization(method)
      yield
    end
  end

  ##
  # Versioned API clients for the Migrieren service.
  #
  # `V1` mirrors the `migrieren.v1` API surface used by the feature-test harness.
  module V1
    class << self
      ##
      # Returns the HTTP façade client for the v1 API.
      #
      # The client targets the default local server endpoint and is memoized for
      # reuse across steps.
      #
      # @return [Migrieren::V1::HTTP] a memoized HTTP client
      def server_http
        @server_http ||= Migrieren::V1::HTTP.new('http://localhost:11000')
      end

      ##
      # Returns the gRPC stub for the v1 API.
      #
      # The stub targets the default local server endpoint and uses
      # {Migrieren.user_agent} channel args.
      #
      # @return [Migrieren::V1::Service::Stub] a memoized gRPC stub
      def server_grpc
        @server_grpc ||= Migrieren::V1::Service::Stub.new(
          'localhost:12000', :this_channel_is_insecure,
          channel_args: Migrieren.user_agent, interceptors: [Migrieren::GRPCAuthorization.new]
        )
      end
    end
  end
end
