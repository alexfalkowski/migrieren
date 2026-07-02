# frozen_string_literal: true

When('I request configured databases with HTTP') do
  @response = request_databases_with_http
end

Then('I should receive configured databases from HTTP:') do |table|
  expect(@response.code).to eq(200)

  resp = JSON.parse(@response.body)

  expect(resp['meta'].length).to be > 0
  expect(resp['databases'].map { |database| database['name'] }).to eq(table.hashes.map { |row| row['database'] })
end

def request_databases_with_http
  opts = Migrieren.http_options(
    headers: {
      user_agent: 'Migrieren-ruby-client/1.0 HTTP/1.0',
      content_type: :json, accept: :json
    }
  )

  Migrieren::V1.server_http.list_databases(opts)
end
