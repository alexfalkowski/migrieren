# frozen_string_literal: true

When('I try to start the server with config {string}') do |config|
  command = ['../migrieren', 'server', '-config', "file:.config/#{config}"]
  stdout = +''
  stderr = +''

  @server_status = nil
  @server_timeout = false

  Open3.popen3(*command) do |stdin, out, err, wait_thr|
    stdin.close

    begin
      Timeout.timeout(3) do
        @server_status = wait_thr.value
      end
    rescue Timeout::Error
      @server_timeout = true
      Process.kill('TERM', wait_thr.pid)
      @server_status = wait_thr.value
    end

    stdout = out.read
    stderr = err.read
  end

  @server_output = "#{stdout}\n#{stderr}"
end

Then('the server should fail to start') do
  expect(@server_timeout).to be(false), @server_output
  expect(@server_status).not_to be_success, @server_output
end
