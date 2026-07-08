@startup
Feature: HTTP apply API
  These endpoints allow users to apply pending migrations.

  @clean
  Scenario: Apply all pending migrations
    When I request to apply migrations with HTTP:
      | database | logs |
    Then I should receive a successful migration from HTTP:
      | database | logs |
      | version  |   40 |

  @clean
  Scenario: Apply all pending migrations when already current
    When I request to apply migrations with HTTP:
      | database | logs |
    And I request to apply migrations with HTTP:
      | database | logs |
    Then I should receive a successful migration from HTTP:
      | database | logs |
      | version  |   40 |

  @clean
  Scenario: Apply migrations truncates logs to the configured maximum
    When I request to apply migrations with HTTP:
      | database | logs |
    Then I should receive truncated migration logs from HTTP:
      | max | 20 |

  Scenario: Apply missing databases
    When I request to apply migrations with HTTP:
      | database | missing |
    Then I should receive a not found migration from HTTP

  Scenario Outline: Apply misconfigured databases
    When I request to apply migrations with HTTP:
      | database | <database> |
    Then I should receive an invalid migration from HTTP

    Examples:
      | database       |
      | missing_source |
      | missing_url    |

  @clean
  Scenario Outline: Return apply failure diagnostics
    When I request to apply migrations with HTTP:
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
