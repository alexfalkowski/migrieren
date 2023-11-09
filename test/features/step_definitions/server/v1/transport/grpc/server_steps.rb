# frozen_string_literal: true

When('I request to migrate with gRPC:') do |table|
  @response = request_with_grpc(table)
end

Then('I should receive a successful migration from gRPC:') do |table|
  rows = table.rows_hash

  expect(@response.meta.length).to be > 0
  expect(@response.migration.database).to eq(rows['database'])
  expect(@response.migration.version).to eq(rows['version'].to_i)
  expect(@response.migration.logs.length).to be > 0
end

Then('I should receive a not found migration from gRPC') do
  expect(@response).to be_a(GRPC::NotFound)
end

Then('I should receive an invalid migration from gRPC') do
  expect(@response).to be_a(GRPC::Internal)
end

def request_with_grpc(table)
  rows = table.rows_hash
  @request_id = SecureRandom.uuid
  metadata = { 'request-id' => @request_id }

  request = Migrieren::V1::MigrateRequest.new(database: rows['database'], version: rows['version'].to_i)
  Migrieren::V1.server_grpc.migrate(request, { metadata: })
rescue StandardError => e
  e
end
