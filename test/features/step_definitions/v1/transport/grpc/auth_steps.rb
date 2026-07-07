# frozen_string_literal: true

AUTH_GRPC_METHOD = '/migrieren.v1.Service/ListDatabases'

When('I request configured databases with gRPC and an authorized token') do
  @response = request_databases_with_grpc_authorization(Migrieren.grpc_authorization(AUTH_GRPC_METHOD))
end

When('I request configured databases with gRPC and no token') do
  @response = request_databases_with_grpc_authorization('')
end

When('I request configured databases with gRPC and an invalid token') do
  @response = request_databases_with_grpc_authorization('Bearer not-a-real-token')
end

When('I request configured databases with gRPC and an unauthorized token') do
  @response = request_databases_with_grpc_authorization(Migrieren.grpc_authorization(AUTH_GRPC_METHOD, 'guest'))
end

Then('I should receive an authorized response from gRPC') do
  expect(@response.databases.length).to be > 0
end

Then('I should receive an unauthenticated response from gRPC') do
  expect(@response).to be_a(GRPC::Unauthenticated)
end

Then('I should receive a forbidden response from gRPC') do
  expect(@response).to be_a(GRPC::PermissionDenied)
end

def request_databases_with_grpc_authorization(authorization)
  request = Migrieren::V1::ListDatabasesRequest.new

  Migrieren::V1.server_grpc.list_databases(request, Migrieren.grpc_options(metadata: { 'authorization' => authorization }))
rescue StandardError => e
  e
end
