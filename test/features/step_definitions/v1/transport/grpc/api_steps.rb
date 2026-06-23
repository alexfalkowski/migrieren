# frozen_string_literal: true

When('I request to migrate with gRPC:') do |table|
  @response = request_with_grpc(table)
end

When('I request migration status with gRPC:') do |table|
  @response = request_status_with_grpc(table)
end

When('I cancel a migration with gRPC:') do |table|
  @response = cancel_with_grpc(table)
end

Then('I should receive a successful migration from gRPC:') do |table|
  rows = table.rows_hash

  expect(@response.meta.length).to be > 0
  expect(@response.migration.database).to eq(rows['database'])
  expect(@response.migration.version).to eq(rows['version'].to_i)
  expect(@response.migration.logs.length).to be >= 0

  expect_postgres_migration(rows['version'].to_i) if rows['database'] == 'postgres'
end

Then('I should receive a not found migration from gRPC') do
  expect(@response).to be_a(GRPC::NotFound)
end

Then('I should receive a migration status from gRPC:') do |table|
  rows = table.rows_hash

  expect(@response.meta.length).to be > 0
  expect(@response.status.database).to eq(rows['database'])
  expect(@response.status.version).to eq(rows['version'].to_i)
  expect(migration_state(@response.status.state)).to eq(rows['state'])
end

Then('I should receive an invalid argument migration from gRPC') do
  expect(@response).to be_a(GRPC::InvalidArgument)
end

Then('I should receive a stopped deadline migration from gRPC') do
  expect(@response).to be_a(GRPC::DeadlineExceeded).or be_a(GRPC::Cancelled)
end

Then('I should receive a canceled migration from gRPC') do
  expect(@response).to be_a(GRPC::Cancelled)
end

Then('I should receive bounded migration logs from gRPC') do
  logs = @response.migration.logs

  expect(logs.length).to be <= 100
  expect(logs.first).to eq('migration logs truncated')
end

Then('I should receive an invalid migration from gRPC') do
  expect(@response).to be_a(GRPC::Internal)
end

Then('I should receive failure diagnostics from gRPC:') do |table|
  rows = table.rows_hash
  metadata = @response.metadata
  log_count = grpc_metadata_value(metadata, 'migration-log-count').to_i

  expect(grpc_metadata_value(metadata, 'migration-error')).to eq(rows['error'])
  expect(grpc_metadata_value(metadata, 'migration-stage')).to eq(empty_to_nil(rows['stage']))

  if rows['logs'] == 'present'
    expect(log_count).to be > 0
    expect(grpc_metadata_value(metadata, 'migration-log-last')).not_to be_empty
  else
    expect(log_count).to eq(0)
    expect(grpc_metadata_value(metadata, 'migration-log-last')).to be_nil
  end
end

def request_with_grpc(table)
  rows = table.rows_hash
  deadline = Time.now + 1 if rows['database'] == 'timeout'
  request = Migrieren::V1::MigrateRequest.new(database: rows['database'], version: rows['version'].to_i)

  Migrieren::V1.server_grpc.migrate(request, Migrieren.grpc_options(deadline:))
rescue StandardError => e
  e
end

def request_status_with_grpc(table)
  rows = table.rows_hash
  request = Migrieren::V1::StatusRequest.new(database: rows['database'])

  Migrieren::V1.server_grpc.status(request, Migrieren.grpc_options)
rescue StandardError => e
  e
end

def cancel_with_grpc(table)
  rows = table.rows_hash
  request = Migrieren::V1::MigrateRequest.new(database: rows['database'], version: rows['version'].to_i)
  operation = grpc_migrate_operation(request)
  thread = execute_grpc_operation(operation)

  sleep 0.2
  operation.cancel
  thread.value
end

def grpc_migrate_operation(request)
  Migrieren::V1.server_grpc.request_response(
    '/migrieren.v1.Service/Migrate',
    request,
    Migrieren::V1::MigrateRequest.method(:encode),
    Migrieren::V1::MigrateResponse.method(:decode),
    **Migrieren.grpc_options,
    return_op: true
  )
end

def execute_grpc_operation(operation)
  Thread.new do
    operation.execute
  rescue StandardError => e
    e
  end
end

def grpc_metadata_value(metadata, key)
  value = metadata[key]

  value.is_a?(Array) ? value.first : value
end

def empty_to_nil(value)
  value.empty? ? nil : value
end
