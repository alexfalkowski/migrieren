# frozen_string_literal: true

require 'securerandom'
require 'yaml'
require 'base64'

require 'pg'
require 'grpc/health/v1/health_services_pb'

require 'migrieren/pg'
require 'migrieren/v1/http'
require 'migrieren/v1/service_services_pb'

module Migrieren
  class << self
    def observability
      @observability ||= Nonnative::Observability.new('http://localhost:11000')
    end

    def server_config
      @server_config ||= Nonnative.configurations('.config/server.yml')
    end

    def health_grpc
      @health_grpc ||= Grpc::Health::V1::Health::Stub.new('localhost:12000', :this_channel_is_insecure, channel_args: Migrieren.user_agent)
    end

    def pg
      @pg ||= Migrieren::PG.new
    end

    def user_agent
      @user_agent ||= Nonnative::Header.grpc_user_agent('Migrieren-ruby-client/1.0 gRPC/1.0')
    end

    def token
      Nonnative::Header.auth_bearer(Base64.decode64(File.read('secrets/token')))
    end
  end

  module V1
    class << self
      def server_http
        @server_http ||= Migrieren::V1::HTTP.new('http://localhost:11000')
      end

      def server_grpc
        @server_grpc ||= Migrieren::V1::Service::Stub.new('localhost:12000', :this_channel_is_insecure, channel_args: Migrieren.user_agent)
      end
    end
  end
end
