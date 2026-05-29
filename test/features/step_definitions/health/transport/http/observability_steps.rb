# frozen_string_literal: true

When('the system requests the {string} with HTTP') do |name|
  @observability_name = name
  @observability_opts = {
    headers: { request_id: SecureRandom.uuid },
    read_timeout: 10, open_timeout: 10
  }
end

Then('the system should respond with a healthy status with HTTP') do
  wait_for do
    response = Nonnative.observability.send(@observability_name, @observability_opts)
    [response.code, response.body.strip]
  end.to eq([200, 'SERVING'])
end

Then('the system should respond with an unhealthy status with HTTP') do
  wait_for do
    response = Nonnative.observability.send(@observability_name, @observability_opts)
    response.code
  end.to eq(503)
end

Then('the system should respond with metrics') do
  response = Nonnative.observability.send(@observability_name, @observability_opts)

  expect(response.code).to eq(200)
  expect(response.body).to include('go_info')
end
