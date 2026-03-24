# frozen_string_literal: true

When('I request to migrate with HTTP which performs in {int} ms') do |time|
  opts = {
    headers: {
      request_id: SecureRandom.uuid, user_agent: 'Bezeichner-ruby-client/1.0 HTTP/1.0',
      content_type: :json, accept: :json
    },
    read_timeout: 10, open_timeout: 10
  }

  expect { Migrieren::V1.server_http.migrate('postgres', rand(1..2), opts) }.to perform_under(time).ms
end
