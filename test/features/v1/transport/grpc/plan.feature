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
