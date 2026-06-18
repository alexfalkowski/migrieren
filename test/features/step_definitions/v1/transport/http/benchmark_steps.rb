# frozen_string_literal: true

When('I request to migrate with HTTP which performs in {int} ms') do |time|
  opts = Migrieren.http_options(
    headers: {
      user_agent: 'Bezeichner-ruby-client/1.0 HTTP/1.0',
      content_type: :json, accept: :json
    }
  )

  expect do
    response = Migrieren::V1.server_http.migrate('postgres', rand(1..2), opts)

    expect(response.code).to eq(200)
  end.to perform_under(time).ms
end
