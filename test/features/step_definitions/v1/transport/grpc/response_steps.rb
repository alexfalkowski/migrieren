# frozen_string_literal: true

Then('I should receive a not found migration from gRPC') do
  expect(@response).to be_a(GRPC::NotFound)
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

def grpc_metadata_value(metadata, key)
  value = metadata[key]

  value.is_a?(Array) ? value.first : value
end
