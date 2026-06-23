# frozen_string_literal: true

When('I request to migrate with HTTP:') do |table|
  @response = request_with_http(table)
end

When('I request migration status with HTTP:') do |table|
  @response = request_status_with_http(table)
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
end

Then('I should receive a not found migration from HTTP') do
  expect(@response.code).to eq(404)
end

Then('I should receive a migration status from HTTP:') do |table|
  expect(@response.code).to eq(200)

  resp = JSON.parse(@response.body)
  status = resp['status']
  rows = table.rows_hash

  expect(resp['meta'].length).to be > 0
  expect(status['database']).to eq(rows['database'])
  expect(status.fetch('version', 0)).to eq(rows['version'].to_i)
  expect(migration_state(status['state'])).to eq(rows['state'])
end

Then('I should receive an invalid argument migration from HTTP') do
  expect(@response.code).to eq(400)
end

Then('I should receive an invalid migration from HTTP') do
  expect(@response.code).to eq(500)
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
end

def request_status_with_http(table)
  rows = table.rows_hash
  opts = Migrieren.http_options(
    headers: {
      user_agent: 'Migrieren-ruby-client/1.0 HTTP/1.0',
      content_type: :json, accept: :json
    }
  )

  Migrieren::V1.server_http.status(rows['database'], opts)
end
