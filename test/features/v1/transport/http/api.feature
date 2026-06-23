@startup
Feature: HTTP API
  These endpoints allows users to migrate different databases.

  Scenario Outline: Migrate valid databases
    When I request to migrate with HTTP:
      | database | <database> |
      | version  | <version>  |
    Then I should receive a successful migration from HTTP:
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
    When I request to migrate with HTTP:
      | database | <database> |
      | version  | <version>  |
    Then I should receive a not found migration from HTTP

    Examples:
      | database | version |
      | missing  |       1 |

  Scenario: List configured databases
    When I request configured databases with HTTP
    Then I should receive configured databases from HTTP:
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

  Scenario: Report misconfigured migration status
    When I request migration status with HTTP:
      | database | missing_url |
    Then I should receive an invalid migration from HTTP

  Scenario Outline: Reject <case> migration version
    When I request to migrate with HTTP:
      | database | postgres  |
      | version  | <version> |
    Then I should receive an invalid argument migration from HTTP

    Examples:
      | case      | version             |
      | zero      |                   0 |
      | oversized | 9223372036854775808 |

  @clean
  Scenario Outline: Migrate misconfigured databases
    When I request to migrate with HTTP:
      | database | <database> |
      | version  | <version>  |
    Then I should receive an invalid migration from HTTP

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
    When I request to migrate with HTTP:
      | database | postgres |
      | version  |        1 |
    Then I should receive an invalid migration from HTTP
