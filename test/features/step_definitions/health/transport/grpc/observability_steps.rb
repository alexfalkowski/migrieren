# frozen_string_literal: true

When('the system requests the health status with gRPC') do
  @response = Migrieren.health_grpc.check(deadline: Time.now + 6)
end

Then('the system should respond with a healthy status with gRPC') do
  expect(@response.status).to eq(:SERVING)
end
