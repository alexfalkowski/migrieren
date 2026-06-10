Feature: Configuration
  Invalid configuration should stop the server from starting.

  Scenario Outline: Reject invalid migration target lists
    When I try to start the server with config "<config>"
    Then the server should fail to start

    Examples:
      | config                         |
      | empty_databases.yml            |
      | duplicate_database_names.yml   |
