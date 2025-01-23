# frozen_string_literal: true

module Migrieren
  class PG
    def initialize
      uri = URI.parse('postgres://test:test@localhost:5432/test?sslmode=disable')
      @conn = ::PG.connect(uri.hostname, uri.port, nil, nil, uri.path[1..], uri.user, uri.password)

      @conn.set_notice_processor { |_| }
    end

    def destroy
      @conn.exec('DROP TABLE IF EXISTS accounts')
      @conn.exec('DROP TABLE IF EXISTS schema_migrations')
    end
  end
end
