@startup
Feature: gRPC apply API
  These endpoints allow users to apply pending migrations.

  @clean
  Scenario: Apply all pending migrations
    When I request to apply migrations with gRPC:
      | database | logs |
    Then I should receive a successful migration from gRPC:
      | database | logs |
      | version  |   40 |

  @clean
  Scenario: Apply all pending migrations when already current
    When I request to apply migrations with gRPC:
      | database | logs |
    And I request to apply migrations with gRPC:
      | database | logs |
    Then I should receive a successful migration from gRPC:
      | database | logs |
      | version  |   40 |

  @clean
  Scenario: Apply migrations truncates logs to the configured maximum
    When I request to apply migrations with gRPC:
      | database | logs |
    Then I should receive truncated migration logs from gRPC:
      | max | 20 |

  Scenario: Apply missing databases
    When I request to apply migrations with gRPC:
      | database | missing |
    Then I should receive a not found migration from gRPC

  Scenario Outline: Apply misconfigured databases
    When I request to apply migrations with gRPC:
      | database | <database> |
    Then I should receive an invalid migration from gRPC

    Examples:
      | database       |
      | missing_source |
      | missing_url    |
