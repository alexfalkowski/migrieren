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
      | postgres |       2 |
      | postgres |       1 |

  Scenario Outline: Migrate missing databases
    When I request to migrate with HTTP:
      | database | <database> |
      | version  | <version>  |
    Then I should receive a not found migration from HTTP

    Examples:
      | database | version |
      | missing  |       1 |

  Scenario Outline: Migrate misconfigured databases
    When I request to migrate with HTTP:
      | database | <database> |
      | version  | <version>  |
    Then I should receive an invalid migration from HTTP

    Examples:
      | database       | version |
      | invalid_source |       1 |
      | missing_url    |       1 |
      | invalid_url    |       1 |

  @failure
  Scenario: Migrate erroneous databases
    Given I should see "postgres" as unhealthy
    When I request to migrate with HTTP:
      | database | postgres |
      | version  |        1 |
    Then I should receive an invalid migration from HTTP
