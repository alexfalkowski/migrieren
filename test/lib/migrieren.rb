# frozen_string_literal: true

require 'securerandom'
require 'yaml'
require 'base64'

require 'pg'
require 'grpc/health/v1/health_services_pb'

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
#   Migrieren::V1.server_http.migrate('test', 1)
#
# @example Using the gRPC API stub
#   req = Migrieren::V1::MigrateRequest.new(database: 'test', version: 1)
#   Migrieren::V1.server_grpc.migrate(req)
#
module Migrieren
  class << self
    ##
    # Returns the observability client used by tests.
    #
    # This client is provided by the `nonnative` test utilities and is used by
    # observability-related steps to query or assert on telemetry emitted by the
    # service.
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
    # The harness loads `.config/server.yml` via `Nonnative.configurations` and
    # memoizes the resulting configuration object.
    #
    # @return [Object] the configuration structure returned by `nonnative`
    def server_config
      @server_config ||= Nonnative.configurations('.config/server.yml')
    end

    ##
    # Returns a gRPC Health Check stub connected to the test server.
    #
    # This is used by health-related steps to query standard gRPC health
    # endpoints.
    #
    # @return [Grpc::Health::V1::Health::Stub] a memoized gRPC health stub
    def health_grpc
      @health_grpc ||= Grpc::Health::V1::Health::Stub.new('localhost:12000', :this_channel_is_insecure, channel_args: Migrieren.user_agent)
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
        @server_grpc ||= Migrieren::V1::Service::Stub.new('localhost:12000', :this_channel_is_insecure, channel_args: Migrieren.user_agent)
      end
    end
  end
end
