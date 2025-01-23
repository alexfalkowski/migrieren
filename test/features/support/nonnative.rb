# frozen_string_literal: true

Nonnative.configure do |config|
  config.load_file('nonnative.yml')
end

Before('@failure') do
  service = Nonnative.pool.service_by_name('postgres')
  service.proxy.close_all

  sleep 1
end

After('@failure') do
  service = Nonnative.pool.service_by_name('postgres')
  service.proxy.reset

  sleep 1
end

After('@clean') do
  Migrieren.pg.destroy
end
