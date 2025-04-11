# frozen_string_literal: true

module Migrieren
  module V1
    class HTTP < Nonnative::HTTPClient
      def migrate(database, version, opts = {})
        post('/migrieren.v1.Service/Migrate', { database:, version: }.to_json, opts)
      end
    end
  end
end
