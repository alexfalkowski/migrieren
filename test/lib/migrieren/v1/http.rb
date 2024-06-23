# frozen_string_literal: true

module Migrieren
  module V1
    class HTTP < Nonnative::HTTPClient
      def migrate(database, version, opts = {})
        post('/v1/migrate', { database:, version: }.to_json, opts)
      end
    end
  end
end
