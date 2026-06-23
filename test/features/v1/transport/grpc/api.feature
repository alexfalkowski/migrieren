@startup
Feature: gRPC API
  These endpoints allows users to migrate different databases.

  Scenario Outline: Migrate valid databases
    When I request to migrate with gRPC:
      | database | <database> |
      | version  | <version>  |
    Then I should receive a successful migration from gRPC:
      | database | <database> |
      | version  | <version>  |

    Examples:
      | database | version |
      | postgres |       1 |
      | postgres |       1 |
      | postgres |       2 |
      | postgres |       1 |
      | github   |       1 |
      | github   |       2 |

  Scenario Outline: Migrate missing databases
    When I request to migrate with gRPC:
      | database | <database> |
      | version  | <version>  |
    Then I should receive a not found migration from gRPC

    Examples:
      | database | version |
      | missing  |       1 |

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

  Scenario: Report misconfigured migration status
    When I request migration status with gRPC:
      | database | missing_url |
    Then I should receive an invalid migration from gRPC

  Scenario Outline: Reject <case> migration version
    When I request to migrate with gRPC:
      | database | postgres  |
      | version  | <version> |
    Then I should receive an invalid argument migration from gRPC

    Examples:
      | case      | version             |
      | zero      |                   0 |
      | oversized | 9223372036854775808 |

  @clean
  Scenario: Stop migration when request deadline expires
    When I request to migrate with gRPC:
      | database | timeout |
      | version  |       1 |
    Then I should receive a stopped deadline migration from gRPC
    And I should not see a completed timeout migration

  @clean
  Scenario: Stop migration when request is canceled
    When I cancel a migration with gRPC:
      | database | timeout |
      | version  |       1 |
    Then I should receive a canceled migration from gRPC
    And I should not see a completed timeout migration

  @clean
  Scenario: Return bounded migration logs
    When I request to migrate with gRPC:
      | database | logs |
      | version  |   40 |
    Then I should receive bounded migration logs from gRPC

  @clean
  Scenario Outline: Return failure diagnostics
    When I request to migrate with gRPC:
      | database | <database> |
      | version  | <version>  |
    Then I should receive an invalid migration from gRPC
    And I should receive failure diagnostics from gRPC:
      | error | <error> |
      | logs  | <logs>  |
      | stage | <stage> |

    Examples:
      | database                         | version | error             | logs    | stage  |
      | missing_source                   |       1 | invalid_config    | empty   | source |
      | missing_url                      |       1 | invalid_config    | empty   | url    |
      | invalid_incomplete_quoted_table  |       1 | invalid_config    | empty   |        |
      | postgres                         |       3 | invalid_migration | present |        |

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
  Scenario Outline: Migrate misconfigured databases
    When I request to migrate with gRPC:
      | database | <database> |
      | version  | <version>  |
    Then I should receive an invalid migration from gRPC

    Examples:
      | database                         | version |
      | missing_source                   |       1 |
      | invalid_source                   |       1 |
      | missing_url                      |       1 |
      | invalid_url                      |       1 |
      | invalid_db                       |       1 |
      | invalid_quoted_table             |       1 |
      | invalid_incomplete_quoted_table  |       1 |
      | invalid_port                     |       1 |
      | postgres                         |       3 |

  @reset
  Scenario: Migrate erroneous databases
    Given I set the proxy for service 'postgres' to 'close_all'
    And I should see "postgres" as unhealthy
    When I request to migrate with gRPC:
      | database | postgres |
      | version  |        1 |
    Then I should receive an invalid migration from gRPC
