# frozen_string_literal: true

When('I request to migrate with HTTP:') do |table|
  @response = request_with_http(table)
end

Then('I should receive a successful migration from HTTP:') do |table|
  expect(@response.code).to eq(200)

  resp = JSON.parse(@response.body)
  migration = resp['migration']
  rows = table.rows_hash

  expect(resp['meta'].length).to be > 0
  expect(migration['database']).to eq(rows['database'])
  expect(migration['version']).to eq(rows['version'])
  expect(migration['logs'].length).to be > 0
end

Then('I should receive a not found migration from HTTP') do
  expect(@response.code).to eq(404)
end

Then('I should receive an invalid migration from HTTP') do
  expect(@response.code).to eq(500)
end

def request_with_http(table)
  rows = table.rows_hash
  opts = {
    headers: {
      request_id: SecureRandom.uuid, user_agent: Migrieren.server_config['transport']['http']['user_agent'],
      content_type: :json, accept: :json
    },
    read_timeout: 10, open_timeout: 10
  }

  Migrieren::V1.server_http.migrate(rows['database'], rows['version'], opts)
end
