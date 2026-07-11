@startup
Feature: HTTP plan API
  These endpoints allow users to inspect pending migrations without applying them.

  @clean
  Scenario: Plan all pending migrations
    When I request a migration plan with HTTP:
      | database | logs |
    Then I should receive a migration plan from HTTP:
      | database         | logs      |
      | version          | 0         |
      | state            | unapplied |
      | latest_version   | 40        |
      | target_version   | 40        |
      | direction        | up        |
      | pending_versions | 1..40     |

  @clean
  Scenario: Plan migrations when already current
    When I request to apply migrations with HTTP:
      | database | logs |
    And I request a migration plan with HTTP:
      | database | logs |
    Then I should receive a migration plan from HTTP:
      | database         | logs  |
      | version          | 40    |
      | state            | clean |
      | latest_version   | 40    |
      | target_version   | 40    |
      | direction        | none  |
      | pending_versions |       |

  @clean
  Scenario Outline: Plan an explicit migration target without applying it
    When I request to migrate with HTTP:
      | database | postgres          |
      | version  | <current_version> |
    And I request a migration plan with HTTP:
      | database       | postgres          |
      | target_version | <target_version> |
    Then I should receive a migration plan from HTTP:
      | database         | postgres           |
      | version          | <current_version>  |
      | state            | clean              |
      | latest_version   | 3                  |
      | target_version   | <target_version>   |
      | direction        | <direction>        |
      | pending_versions | <pending_versions> |
    And I request migration status with HTTP:
      | database | postgres          |
    Then I should receive a migration status from HTTP:
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
    When I request a migration plan with HTTP:
      | database       | postgres |
      | target_version |        2 |
    Then I should receive a migration plan from HTTP:
      | database         | postgres  |
      | version          | 0         |
      | state            | unapplied |
      | latest_version   | 3         |
      | target_version   | 2         |
      | direction        | up        |
      | pending_versions | 1,2       |
    And I request migration status with HTTP:
      | database | postgres  |
    Then I should receive a migration status from HTTP:
      | database | postgres  |
      | version  | 0         |
      | state    | unapplied |

  @clean
  Scenario Outline: Reject <case> plan target version
    When I request a migration plan with HTTP:
      | database       | postgres          |
      | target_version | <target_version> |
    Then I should receive an invalid argument migration from HTTP

    Examples:
      | case      | target_version      |
      | zero      |                   0 |
      | oversized | 9223372036854775808 |

  @clean
  Scenario: Reject a plan target that is not in the migration source
    When I request a migration plan with HTTP:
      | database       | postgres |
      | target_version |        4 |
    Then I should receive an invalid migration from HTTP
    And I should receive failure diagnostics from HTTP:
      | error | invalid_migration |
      | logs  | empty             |
      | stage |                   |

  @clean
  Scenario: Reject an explicit target from a dirty migration state
    When I request to migrate with HTTP:
      | database | postgres |
      | version  |        3 |
    Then I should receive an invalid migration from HTTP
    When I request a migration plan with HTTP:
      | database       | postgres |
      | target_version |        1 |
    Then I should receive an invalid migration from HTTP
    When I request migration status with HTTP:
      | database | postgres |
    Then I should receive a migration status from HTTP:
      | database | postgres |
      | version  |        3 |
      | state    | dirty    |

  Scenario: Plan missing databases
    When I request a migration plan with HTTP:
      | database | missing |
    Then I should receive a not found migration from HTTP

  Scenario Outline: Plan misconfigured databases
    When I request a migration plan with HTTP:
      | database | <database> |
    Then I should receive an invalid migration from HTTP

    Examples:
      | database       |
      | invalid_source |
      | missing_source |
      | missing_url    |
      | invalid_url    |

  @clean
  Scenario Outline: Return plan failure diagnostics
    When I request a migration plan with HTTP:
      | database | <database> |
    Then I should receive an invalid migration from HTTP
    And I should receive failure diagnostics from HTTP:
      | error | <error> |
      | logs  | <logs>  |
      | stage | <stage> |

    Examples:
      | database       | error          | logs  | stage  |
      | missing_source | invalid_config | empty | source |
      | missing_url    | invalid_config | empty | url    |
