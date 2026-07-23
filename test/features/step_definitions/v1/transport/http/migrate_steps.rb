# frozen_string_literal: true

When('I request to migrate with HTTP:') do |table|
  @response = request_with_http(table)
end

Then('I should receive a successful migration from HTTP:') do |table|
  expect(@response.code).to eq(200)

  resp = JSON.parse(@response.body)
  migration = resp['migration']
  rows = table.rows_hash
  logs = migration['logs'] || []

  expect(resp['meta'].length).to be > 0
  expect(migration['database']).to eq(rows['database'])
  expect(migration['version']).to eq(rows['version'].to_i)
  expect(logs.length).to be >= 0

  expect_postgres_migration(rows['version'].to_i) if rows['database'] == 'postgres'
  expect_log_migration(rows['version'].to_i) if rows['database'] == 'logs'
end

Then('I should receive an invalid argument migration from HTTP') do
  expect(@response.code).to eq(400)
end

Then('I should receive a timed out migration from HTTP') do
  expect(@response).to(
    be_a(RestClient::Exceptions::ReadTimeout).or(be_a(Timeout::Error))
  )
end

def request_with_http(table)
  rows = table.rows_hash
  opts = Migrieren.http_options(
    headers: {
      user_agent: 'Migrieren-ruby-client/1.0 HTTP/1.0',
      content_type: :json, accept: :json
    }
  )

  Migrieren::V1.server_http.migrate(rows['database'], rows['version'].to_i, opts)
rescue StandardError => e
  e
end
