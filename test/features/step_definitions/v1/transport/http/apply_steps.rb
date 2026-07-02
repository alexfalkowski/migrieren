# frozen_string_literal: true

When('I request to apply migrations with HTTP:') do |table|
  @response = request_apply_with_http(table)
end

def request_apply_with_http(table)
  rows = table.rows_hash
  opts = Migrieren.http_options(
    headers: {
      user_agent: 'Migrieren-ruby-client/1.0 HTTP/1.0',
      content_type: :json, accept: :json
    }
  )

  Migrieren::V1.server_http.apply_migrations(rows['database'], opts)
end
