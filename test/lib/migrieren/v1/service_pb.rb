# frozen_string_literal: true
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: migrieren/v1/service.proto

require 'google/protobuf'

descriptor_data = "\n\x1amigrieren/v1/service.proto\x12\x0cmigrieren.v1\"U\n\tMigration\x12\x1a\n\x08\x64\x61tabase\x18\x01 \x01(\tR\x08\x64\x61tabase\x12\x18\n\x07version\x18\x02 \x01(\x04R\x07version\x12\x12\n\x04logs\x18\x03 \x03(\tR\x04logs\"F\n\x0eMigrateRequest\x12\x1a\n\x08\x64\x61tabase\x18\x01 \x01(\tR\x08\x64\x61tabase\x12\x18\n\x07version\x18\x02 \x01(\x04R\x07version\"\xbe\x01\n\x0fMigrateResponse\x12;\n\x04meta\x18\x01 \x03(\x0b\x32\'.migrieren.v1.MigrateResponse.MetaEntryR\x04meta\x12\x35\n\tmigration\x18\x02 \x01(\x0b\x32\x17.migrieren.v1.MigrationR\tmigration\x1a\x37\n\tMetaEntry\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value:\x02\x38\x01\x32S\n\x07Service\x12H\n\x07Migrate\x12\x1c.migrieren.v1.MigrateRequest\x1a\x1d.migrieren.v1.MigrateResponse\"\x00\x42\x45Z3github.com/alexfalkowski/migrieren/api/migrieren/v1\xea\x02\rMigrieren::V1b\x06proto3"

pool = Google::Protobuf::DescriptorPool.generated_pool
pool.add_serialized_file(descriptor_data)

module Migrieren
  module V1
    Migration = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("migrieren.v1.Migration").msgclass
    MigrateRequest = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("migrieren.v1.MigrateRequest").msgclass
    MigrateResponse = ::Google::Protobuf::DescriptorPool.generated_pool.lookup("migrieren.v1.MigrateResponse").msgclass
  end
end
