# frozen_string_literal: true

When('I request configured databases with gRPC') do
  @response = request_databases_with_grpc
end

Then('I should receive configured databases from gRPC:') do |table|
  expect(@response.meta.length).to be > 0
  expect(@response.databases.map(&:name)).to eq(table.hashes.map { |row| row['database'] })
end

def request_databases_with_grpc
  request = Migrieren::V1::ListDatabasesRequest.new

  Migrieren::V1.server_grpc.list_databases(request, Migrieren.grpc_options)
rescue StandardError => e
  e
end
