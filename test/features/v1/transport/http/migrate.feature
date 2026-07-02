@startup
Feature: HTTP migrate API
  These endpoints allow users to migrate configured databases.

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

    @github
    Examples:
      | database | version |
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

  @clean
  Scenario Outline: Return failure diagnostics
    When I request to migrate with HTTP:
      | database | <database> |
      | version  | <version>  |
    Then I should receive an invalid migration from HTTP
    And I should receive failure diagnostics from HTTP:
      | error | <error> |
      | logs  | <logs>  |
      | stage | <stage> |

    Examples:
      | database        | version | error             | logs    | stage  |
      | missing_source  |       1 | invalid_config    | empty   | source |
      | missing_url     |       1 | invalid_config    | empty   | url    |
      | postgres        |       3 | invalid_migration | present |        |

  @reset
  Scenario: Migrate erroneous databases
    Given I set the proxy for service 'postgres' to 'close_all'
    And I should see "postgres" as unhealthy
    When I request to migrate with HTTP:
      | database | postgres |
      | version  |        1 |
    Then I should receive an invalid migration from HTTP
