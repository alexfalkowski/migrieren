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
      # @return [Object] whatever `Nonnative::HTTPClient#post` returns (typically a response wrapper)
      def migrate(database, version, opts = {})
        post('/migrieren.v1.Service/Migrate', { database:, version: }.to_json, opts)
      end
    end
  end
end
