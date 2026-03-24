@startup @benchmark
Feature: Benchmark HTTP API
  Make sure these endpoints perform at their best.

  Scenario: Migrate database in a good time frame and memory.
    When I request to migrate with HTTP which performs in 5000 ms
    And the process 'server' should consume less than '70mb' of memory
