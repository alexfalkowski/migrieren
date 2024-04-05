@manual
Feature: Client

  Client allows the system to migrate the databse to specific version

  Scenario: Migrate valid database
    Given the client is configured with "valid" config
    And I start the system
    When the client tries to migrate the database
    Then the client should have succesfully migrated the database
    And I should see a log entry of "finished call with code OK" in the file "reports/client.log"

  Scenario: Migrate missing databases
    Given the client is configured with "missing" config
    And I start the system
    When the client tries to migrate the database
    Then the client should have unsuccesfully migrated the database
    And I should see a log entry of "not found" in the file "reports/client.log"

  Scenario: Migrate misconfigured databases
    Given the client is configured with "misconfigured" config
    And I start the system
    When the client tries to migrate the database
    Then the client should have unsuccesfully migrated the database
    And I should see a log entry of "invalid config" in the file "reports/client.log"

   Scenario: Migrate with invalid host
    Given the client is configured with "invalid" config
    And I start the system
    When the client tries to migrate the database
    Then the client should have unsuccesfully migrated the database
    And I should see a log entry of "2020: connect: connection refused" in the file "reports/client.log"
