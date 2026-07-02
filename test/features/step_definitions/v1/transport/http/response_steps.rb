# frozen_string_literal: true

Then('I should receive a not found migration from HTTP') do
  expect(@response.code).to eq(404)
end

Then('I should receive an invalid migration from HTTP') do
  expect(@response.code).to eq(500)
end

Then('I should receive failure diagnostics from HTTP:') do |table|
  rows = table.rows_hash
  log_count = http_header_value(@response, 'migration-log-count').to_i

  expect(http_header_value(@response, 'migration-error')).to eq(rows['error'])
  expect(http_header_value(@response, 'migration-stage')).to eq(empty_to_nil(rows['stage']))

  if rows['logs'] == 'present'
    expect(log_count).to be > 0
    expect(http_header_value(@response, 'migration-log-last')).not_to be_empty
  else
    expect(log_count).to eq(0)
    expect(http_header_value(@response, 'migration-log-last')).to be_nil
  end
end

def http_header_value(response, key)
  response.headers[key.tr('-', '_').to_sym]
end
