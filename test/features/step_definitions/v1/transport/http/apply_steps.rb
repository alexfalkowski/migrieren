# frozen_string_literal: true

When('I request to apply migrations with HTTP:') do |table|
  @response = request_apply_with_http(table)
end

Then('I should receive truncated migration logs from HTTP:') do |table|
  rows = table.rows_hash
  resp = JSON.parse(@response.body)
  logs = resp['migration']['logs'] || []

  expect(logs.length).to eq(rows['max'].to_i)
  expect(logs.first).to match(/\Amigration logs truncated \(showing last \d+ of \d+\)\z/)
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
