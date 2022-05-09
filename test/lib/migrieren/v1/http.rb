# frozen_string_literal: true

module Migrieren
  module V1
    class HTTP < Nonnative::HTTPClient
      def migrate(database, version, headers = {})
        headers.merge!(content_type: :json, accept: :json)

        get("v1/migrate/#{database}/#{version}", headers, 10)
      end
    end
  end
end
