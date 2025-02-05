# frozen_string_literal: true

Nonnative.configure do |config|
  config.load_file('nonnative.yml')
end

BeforeAll do
  Migrieren.pg.destroy
end

After('@clean') do
  Migrieren.pg.destroy
end
