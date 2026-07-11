@startup
Feature: gRPC plan API
  These endpoints allow users to inspect pending migrations without applying them.

  @clean
  Scenario: Plan all pending migrations
    When I request a migration plan with gRPC:
      | database | logs |
    Then I should receive a migration plan from gRPC:
      | database         | logs      |
      | version          | 0         |
      | state            | unapplied |
      | latest_version   | 40        |
      | target_version   | 40        |
      | direction        | up        |
      | pending_versions | 1..40     |

  @clean
  Scenario: Plan migrations when already current
    When I request to apply migrations with gRPC:
      | database | logs |
    And I request a migration plan with gRPC:
      | database | logs |
    Then I should receive a migration plan from gRPC:
      | database         | logs  |
      | version          | 40    |
      | state            | clean |
      | latest_version   | 40    |
      | target_version   | 40    |
      | direction        | none  |
      | pending_versions |       |

  @clean
  Scenario Outline: Plan an explicit migration target without applying it
    When I request to migrate with gRPC:
      | database | postgres          |
      | version  | <current_version> |
    And I request a migration plan with gRPC:
      | database       | postgres         |
      | target_version | <target_version> |
    Then I should receive a migration plan from gRPC:
      | database         | postgres         |
      | version          | <current_version> |
      | state            | clean            |
      | latest_version   | 3                |
      | target_version   | <target_version> |
      | direction        | <direction>      |
      | pending_versions | <pending_versions> |
    And I request migration status with gRPC:
      | database | postgres          |
    Then I should receive a migration status from gRPC:
      | database | postgres          |
      | version  | <current_version> |
      | state    | clean             |

    Examples:
      | current_version | target_version | direction | pending_versions |
      |               1 |              1 | none      |                  |
      |               1 |              2 | up        | 2                |
      |               2 |              1 | down      | 2                |

  @clean
  Scenario: Plan an explicit target from an unapplied database
    When I request a migration plan with gRPC:
      | database       | postgres |
      | target_version |        2 |
    Then I should receive a migration plan from gRPC:
      | database         | postgres  |
      | version          | 0         |
      | state            | unapplied |
      | latest_version   | 3         |
      | target_version   | 2         |
      | direction        | up        |
      | pending_versions | 1,2       |
    And I request migration status with gRPC:
      | database | postgres  |
    Then I should receive a migration status from gRPC:
      | database | postgres  |
      | version  | 0         |
      | state    | unapplied |

  @clean
  Scenario Outline: Reject <case> plan target version
    When I request a migration plan with gRPC:
      | database       | postgres          |
      | target_version | <target_version> |
    Then I should receive an invalid argument migration from gRPC

    Examples:
      | case      | target_version      |
      | zero      |                   0 |
      | oversized | 9223372036854775808 |

  @clean
  Scenario: Reject a plan target that is not in the migration source
    When I request a migration plan with gRPC:
      | database       | postgres |
      | target_version |        4 |
    Then I should receive an invalid migration from gRPC
    And I should receive failure diagnostics from gRPC:
      | error | invalid_migration |
      | logs  | empty             |
      | stage |                   |

  @clean
  Scenario: Reject an explicit target from a dirty migration state
    When I request to migrate with gRPC:
      | database | postgres |
      | version  |        3 |
    Then I should receive an invalid migration from gRPC
    When I request a migration plan with gRPC:
      | database       | postgres |
      | target_version |        1 |
    Then I should receive an invalid migration from gRPC
    When I request migration status with gRPC:
      | database | postgres |
    Then I should receive a migration status from gRPC:
      | database | postgres |
      | version  |        3 |
      | state    | dirty    |

  Scenario: Plan missing databases
    When I request a migration plan with gRPC:
      | database | missing |
    Then I should receive a not found migration from gRPC

  Scenario Outline: Plan misconfigured databases
    When I request a migration plan with gRPC:
      | database | <database> |
    Then I should receive an invalid migration from gRPC

    Examples:
      | database       |
      | invalid_source |
      | missing_source |
      | missing_url    |
      | invalid_url    |

  @clean
  Scenario Outline: Return plan failure diagnostics
    When I request a migration plan with gRPC:
      | database | <database> |
    Then I should receive an invalid migration from gRPC
    And I should receive failure diagnostics from gRPC:
      | error | <error> |
      | logs  | <logs>  |
      | stage | <stage> |

    Examples:
      | database       | error          | logs  | stage  |
      | missing_source | invalid_config | empty | source |
      | missing_url    | invalid_config | empty | url    |
      | invalid_source | invalid_config | empty | source |
      | invalid_url    | invalid_config | empty | url    |
