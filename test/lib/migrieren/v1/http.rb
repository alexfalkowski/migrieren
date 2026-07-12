# frozen_string_literal: true

module Migrieren
  module V1
    ##
    # HTTP façade client for the `migrieren.v1.Service` API.
    #
    # This class is part of the Ruby feature-test harness and is used by
    # Cucumber steps to call the service via its HTTP RPC façade.
    #
    # Under the hood it delegates request execution to `Nonnative::HTTPClient`,
    # and exposes small convenience methods for the endpoints used in tests.
    #
    # Endpoint helpers normally return `RestClient::Response` objects for
    # successful and non-2xx HTTP responses, so callers inspect `code`, `body`,
    # and `headers` directly. Passing `raw_response: true` returns a
    # `RestClient::RawResponse` instead. Timeouts, broken connections, and other
    # transport failures propagate as RestClient or system exceptions rather
    # than being converted into responses.
    #
    # The service is expected to be reachable at the base URL provided when the
    # client is constructed (see `Migrieren::V1.server_http`).
    class HTTP < Nonnative::HTTPClient
      ##
      # Calls the `Migrate` RPC via the HTTP façade.
      #
      # This maps to the HTTP RPC route:
      # `POST /migrieren.v1.Service/Migrate`
      #
      # The request body is JSON with the following shape:
      # `{ "database": String, "version": Integer }`
      #
      # @param database [String] logical database name as configured in the service
      # @param version [Integer] target migration version (encoded as JSON number)
      # @param opts [Hash] optional request options passed through to `post`
      # @return [RestClient::Response, RestClient::RawResponse] the HTTP response;
      #   RawResponse is returned when opts enables raw_response
      def migrate(database, version, opts = {})
        post('/migrieren.v1.Service/Migrate', { database:, version: }.to_json, opts)
      end

      ##
      # Calls the `ApplyMigrations` RPC via the HTTP façade.
      #
      # This maps to the HTTP RPC route:
      # `POST /migrieren.v1.Service/ApplyMigrations`
      #
      # The request body is JSON with the following shape:
      # `{ "database": String }`
      #
      # @param database [String] logical database name as configured in the service
      # @param opts [Hash] optional request options passed through to `post`
      # @return [RestClient::Response, RestClient::RawResponse] the HTTP response;
      #   RawResponse is returned when opts enables raw_response
      def apply_migrations(database, opts = {})
        post('/migrieren.v1.Service/ApplyMigrations', { database: }.to_json, opts)
      end

      ##
      # Calls the `PlanMigrations` RPC via the HTTP façade.
      #
      # This maps to the HTTP RPC route:
      # `POST /migrieren.v1.Service/PlanMigrations`
      #
      # The request body is JSON with the following shape:
      # `{ "database": String, "target_version"?: Integer }`
      #
      # @param database [String] logical database name as configured in the service
      # @param opts [Hash] optional request options passed through to `post`
      # @param target_version [Integer, nil] optional explicit migration version
      #   to preview; when nil, the request preserves latest-up planning
      # @return [RestClient::Response, RestClient::RawResponse] the HTTP response;
      #   RawResponse is returned when opts enables raw_response
      def plan_migrations(database, opts = {}, target_version: nil)
        payload = { database: }
        payload[:target_version] = target_version unless target_version.nil?

        post('/migrieren.v1.Service/PlanMigrations', payload.to_json, opts)
      end

      ##
      # Calls the `Status` RPC via the HTTP façade.
      #
      # This maps to the HTTP RPC route:
      # `POST /migrieren.v1.Service/Status`
      #
      # The request body is JSON with the following shape:
      # `{ "database": String }`
      #
      # @param database [String] logical database name as configured in the service
      # @param opts [Hash] optional request options passed through to `post`
      # @return [RestClient::Response, RestClient::RawResponse] the HTTP response;
      #   RawResponse is returned when opts enables raw_response
      def status(database, opts = {})
        post('/migrieren.v1.Service/Status', { database: }.to_json, opts)
      end

      ##
      # Calls the `ListDatabases` RPC via the HTTP façade.
      #
      # This maps to the HTTP RPC route:
      # `POST /migrieren.v1.Service/ListDatabases`
      #
      # @param opts [Hash] optional request options passed through to `post`
      # @return [RestClient::Response, RestClient::RawResponse] the HTTP response;
      #   RawResponse is returned when opts enables raw_response
      def list_databases(opts = {})
        post('/migrieren.v1.Service/ListDatabases', {}.to_json, opts)
      end

      protected

      ##
      # Authenticates every HTTP RPC call before delegating to the base client.
      #
      # Adds a route-scoped `migrieren` Bearer token (see
      # {Migrieren.authorize_http}) unless the caller already set an
      # `:authorization` header, so authentication is transparent to the endpoint
      # helpers above while rejection scenarios opt out by passing their own
      # header.
      #
      # @param pathname [String] the HTTP RPC path
      # @param payload [String] the JSON request body
      # @param opts [Hash] request options passed to `Nonnative::HTTPClient#post`
      # @return [RestClient::Response, RestClient::RawResponse] the base client's
      #   response; RawResponse is returned when opts enables raw_response
      def post(pathname, payload, opts = {})
        super(pathname, payload, Migrieren.authorize_http(pathname, opts))
      end
    end
  end
end
