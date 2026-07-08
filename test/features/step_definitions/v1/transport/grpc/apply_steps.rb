# frozen_string_literal: true

When('I request to apply migrations with gRPC:') do |table|
  @response = request_apply_with_grpc(table)
end

Then('I should receive truncated migration logs from gRPC:') do |table|
  rows = table.rows_hash
  logs = @response.migration.logs

  expect(logs.length).to eq(rows['max'].to_i)
  expect(logs.first).to match(/\Amigration logs truncated \(showing last \d+ of \d+\)\z/)
end

def request_apply_with_grpc(table)
  rows = table.rows_hash
  request = Migrieren::V1::ApplyMigrationsRequest.new(database: rows['database'])

  Migrieren::V1.server_grpc.apply_migrations(request, Migrieren.grpc_options)
rescue StandardError => e
  e
end
