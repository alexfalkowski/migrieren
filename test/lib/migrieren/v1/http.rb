# frozen_string_literal: true

module Migrieren
  module V1
    class HTTP < Nonnative::HTTPClient
      def migrate(database, version, opts = {})
        get("v1/migrate/#{database}/#{version}", opts)
      end
    end
  end
end
