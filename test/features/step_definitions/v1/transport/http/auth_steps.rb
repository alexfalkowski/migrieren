# frozen_string_literal: true

AUTH_HTTP_PATH = '/migrieren.v1.Service/ListDatabases'

When('I request configured databases with HTTP and an authorized token') do
  @response = request_databases_with_http_authorization(Migrieren.http_authorization(AUTH_HTTP_PATH))
end

When('I request configured databases with HTTP and no token') do
  @response = request_databases_with_http_authorization('')
end

When('I request configured databases with HTTP and an invalid token') do
  @response = request_databases_with_http_authorization('Bearer not-a-real-token')
end

When('I request configured databases with HTTP and an unauthorized token') do
  @response = request_databases_with_http_authorization(Migrieren.http_authorization(AUTH_HTTP_PATH, 'guest'))
end

Then('I should receive an authorized response from HTTP') do
  expect(@response.code).to eq(200)
  expect(JSON.parse(@response.body)['databases'].length).to be > 0
end

Then('I should receive an unauthenticated response from HTTP') do
  expect(@response.code).to eq(401)
end

Then('I should receive a forbidden response from HTTP') do
  expect(@response.code).to eq(403)
end

def request_databases_with_http_authorization(authorization)
  opts = Migrieren.http_options(
    headers: {
      user_agent: 'Migrieren-ruby-client/1.0 HTTP/1.0',
      content_type: :json, accept: :json,
      authorization:
    }
  )

  Migrieren::V1.server_http.list_databases(opts)
end
