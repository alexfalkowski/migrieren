# frozen_string_literal: true

When('I request migration status with gRPC:') do |table|
  @response = request_status_with_grpc(table)
end

Then('I should receive a migration status from gRPC:') do |table|
  rows = table.rows_hash

  expect(@response.meta.length).to be > 0
  expect(@response.status.database).to eq(rows['database'])
  expect(@response.status.version).to eq(rows['version'].to_i)
  expect(migration_state(@response.status.state)).to eq(rows['state'])
end

def request_status_with_grpc(table)
  rows = table.rows_hash
  request = Migrieren::V1::StatusRequest.new(database: rows['database'])

  Migrieren::V1.server_grpc.status(request, Migrieren.grpc_options)
rescue StandardError => e
  e
end
