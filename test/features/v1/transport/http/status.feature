@startup
Feature: HTTP status API
  These endpoints allow users to inspect current migration status.

  @clean
  Scenario: Report migration status
    When I request to migrate with HTTP:
      | database | postgres |
      | version  |        1 |
    And I request migration status with HTTP:
      | database | postgres |
    Then I should receive a migration status from HTTP:
      | database | postgres |
      | version  |        1 |
      | state    | clean    |

  @clean
  Scenario: Report unapplied migration status
    When I request migration status with HTTP:
      | database | postgres |
    Then I should receive a migration status from HTTP:
      | database | postgres  |
      | version  |         0 |
      | state    | unapplied |

  Scenario: Report missing migration status
    When I request migration status with HTTP:
      | database | missing |
    Then I should receive a not found migration from HTTP

  Scenario Outline: Report misconfigured migration status
    When I request migration status with HTTP:
      | database | <database> |
    Then I should receive an invalid migration from HTTP
    And I should receive failure diagnostics from HTTP:
      | error | <error> |
      | logs  | <logs>  |
      | stage | <stage> |

    Examples:
      | database    | error          | logs  | stage |
      | missing_url | invalid_config | empty | url   |
      | invalid_url | invalid_config | empty | url   |
