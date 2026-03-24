# frozen_string_literal: true

When('I request to migrate with gRPC which performs in {int} ms') do |time|
  metadata = { 'request-id' => SecureRandom.uuid }
  request = Migrieren::V1::MigrateRequest.new(database: 'postgres', version: rand(1..2))

  expect { Migrieren::V1.server_grpc.migrate(request, { metadata: }) }.to perform_under(time).ms
end
