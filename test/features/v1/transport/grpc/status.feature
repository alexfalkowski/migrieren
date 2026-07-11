@startup
Feature: gRPC status API
  These endpoints allow users to inspect current migration status.

  @clean
  Scenario: Report migration status
    When I request to migrate with gRPC:
      | database | postgres |
      | version  |        1 |
    And I request migration status with gRPC:
      | database | postgres |
    Then I should receive a migration status from gRPC:
      | database | postgres |
      | version  |        1 |
      | state    | clean    |

  @clean
  Scenario: Report unapplied migration status
    When I request migration status with gRPC:
      | database | postgres |
    Then I should receive a migration status from gRPC:
      | database | postgres  |
      | version  |         0 |
      | state    | unapplied |

  Scenario: Report missing migration status
    When I request migration status with gRPC:
      | database | missing |
    Then I should receive a not found migration from gRPC

  Scenario Outline: Report misconfigured migration status
    When I request migration status with gRPC:
      | database | <database> |
    Then I should receive an invalid migration from gRPC
    And I should receive failure diagnostics from gRPC:
      | error | <error> |
      | logs  | <logs>  |
      | stage | <stage> |

    Examples:
      | database    | error          | logs  | stage |
      | missing_url | invalid_config | empty | url   |
      | invalid_url | invalid_config | empty | url   |

  @clean
  Scenario: Report dirty migration status
    When I request to migrate with gRPC:
      | database | postgres |
      | version  |        3 |
    Then I should receive an invalid migration from gRPC
    When I request migration status with gRPC:
      | database | postgres |
    Then I should receive a migration status from gRPC:
      | database | postgres |
      | version  |        3 |
      | state    | dirty    |

  @clean
  Scenario: Report a dirty migration status for an unapplied database
    Given the postgres database has a dirty unapplied migration
    When I request migration status with gRPC:
      | database | postgres |
    Then I should receive a migration status from gRPC:
      | database | postgres |
      | version  |        0 |
      | state    | dirty    |
