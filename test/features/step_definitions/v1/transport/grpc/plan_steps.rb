# frozen_string_literal: true

When('I request a migration plan with gRPC:') do |table|
  @response = request_plan_with_grpc(table)
end

Then('I should receive a migration plan from gRPC:') do |table|
  rows = table.rows_hash
  plan = @response.plan
  status = plan.status

  expect(@response.meta.length).to be > 0
  expect(status.database).to eq(rows['database'])
  expect(status.version).to eq(rows['version'].to_i)
  expect(migration_state(status.state)).to eq(rows['state'])
  expect(plan.latest_version).to eq(rows['latest_version'].to_i)
  expect(plan.target_version).to eq(rows['target_version'].to_i)
  expect(migration_direction(plan.direction)).to eq(rows['direction'])
  expect(plan.pending_versions.to_a).to eq(migration_versions(rows['pending_versions']))
end

def request_plan_with_grpc(table)
  rows = table.rows_hash
  attrs = { database: rows['database'] }
  attrs[:target_version] = rows['target_version'].to_i if rows.key?('target_version')
  request = Migrieren::V1::PlanMigrationsRequest.new(**attrs)

  Migrieren::V1.server_grpc.plan_migrations(request, Migrieren.grpc_options)
rescue StandardError => e
  e
end
