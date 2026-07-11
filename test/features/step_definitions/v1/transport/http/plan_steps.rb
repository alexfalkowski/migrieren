# frozen_string_literal: true

When('I request a migration plan with HTTP:') do |table|
  @response = request_plan_with_http(table)
end

Then('I should receive a migration plan from HTTP:') do |table|
  expect(@response.code).to eq(200)

  resp = JSON.parse(@response.body)
  plan = resp['plan']
  status = plan['status']
  rows = table.rows_hash

  expect(resp['meta'].length).to be > 0
  expect(status['database']).to eq(rows['database'])
  expect(status.fetch('version', 0)).to eq(rows['version'].to_i)
  expect(migration_state(status['state'])).to eq(rows['state'])
  expect(plan.fetch('latest_version', 0)).to eq(rows['latest_version'].to_i)
  expect(plan.fetch('target_version', 0)).to eq(rows['target_version'].to_i)
  expect(migration_direction(plan['direction'])).to eq(rows['direction'])
  expect(plan['pending_versions'] || []).to eq(migration_versions(rows['pending_versions']))
end

def request_plan_with_http(table)
  rows = table.rows_hash

  opts = Migrieren.http_options(
    headers: {
      user_agent: 'Migrieren-ruby-client/1.0 HTTP/1.0',
      content_type: :json, accept: :json
    }
  )
  target_version = rows['target_version'].to_i if rows.key?('target_version')

  Migrieren::V1.server_http.plan_migrations(rows['database'], opts, target_version:)
end
