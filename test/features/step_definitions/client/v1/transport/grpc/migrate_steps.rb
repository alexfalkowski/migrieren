# frozen_string_literal: true

Given('the client is configured with {string} config') do |app|
  @app = app
end

When('the client tries to migrate the database') do
  @status = migrate(@app)
end

Then('the client should have succesfully migrated the database') do
  expect(@status.exitstatus).to eq(0)
end

Then('the client should have unsuccesfully migrated the database') do
  expect(@status.exitstatus).to eq(1)
end

def migrate(app)
  env = {
    'MIGRIEREN_CONFIG_FILE' => ".config/#{app}.client.yml"
  }
  cmd = Nonnative.go_executable(%w[cover], 'reports', '../migrieren', 'migrate')
  pid = spawn(env, cmd, %i[out err] => ['reports/client.log', 'a'])

  _, status = Process.waitpid2(pid)
  status
end
