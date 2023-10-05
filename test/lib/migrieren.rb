# frozen_string_literal: true

require 'securerandom'
require 'yaml'
require 'base64'

require 'grpc/health/v1/health_services_pb'

require 'migrieren/v1/http'
require 'migrieren/v1/service_services_pb'

module Migrieren
  class << self
    def observability
      @observability ||= Nonnative::Observability.new('http://localhost:8080')
    end

    def server_config
      @server_config ||= YAML.load_file('.config/server.yml')
    end

    def health_grpc
      @health_grpc ||= Grpc::Health::V1::Health::Stub.new('localhost:9090', :this_channel_is_insecure)
    end
  end

  module V1
    class << self
      def server_http
        @server_http ||= Migrieren::V1::HTTP.new('http://localhost:8080')
      end

      def server_grpc
        @server_grpc ||= Migrieren::V1::Service::Stub.new('localhost:9090', :this_channel_is_insecure)
      end
    end
  end
end
