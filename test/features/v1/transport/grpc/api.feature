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
      | postgres |       2 |

  Scenario Outline: Migrate missing databases
    When I request to migrate with gRPC:
      | database | <database> |
      | version  | <version>  |
    Then I should receive a not found migration from gRPC

    Examples:
      | database | version |
      | missing  |       1 |

  @clean
  Scenario Outline: Migrate misconfigured databases
    When I request to migrate with gRPC:
      | database | <database> |
      | version  | <version>  |
    Then I should receive an invalid migration from gRPC

    Examples:
      | database       | version |
      | missing_source |       1 |
      | invalid_source |       1 |
      | missing_url    |       1 |
      | invalid_url    |       1 |
      | postgres       |       3 |

  Scenario: Migrate erroneous databases
    Given I set the proxy for service 'postgres' to 'close_all'
    And I should see "postgres" as unhealthy
    When I request to migrate with gRPC:
      | database | postgres |
      | version  |        1 |
    Then I should receive an invalid migration from gRPC
    And I should reset the proxy for service 'postgres'
