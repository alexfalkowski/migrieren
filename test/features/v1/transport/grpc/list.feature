@startup
Feature: gRPC list API
  These endpoints allow users to list configured databases.

  Scenario: List configured databases
    When I request configured databases with gRPC
    Then I should receive configured databases from gRPC:
      | database                        |
      | postgres                        |
      | invalid_source                  |
      | missing_source                  |
      | missing_url                     |
      | invalid_url                     |
      | invalid_db                      |
      | invalid_quoted_table            |
      | invalid_incomplete_quoted_table |
      | invalid_port                    |
      | github                          |
      | timeout                         |
      | logs                            |
