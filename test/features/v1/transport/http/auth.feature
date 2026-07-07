@startup
Feature: HTTP API authentication
  The HTTP RPC facade requires a verified token and an authorized subject.

  Scenario: Accept a request with an authorized token
    When I request configured databases with HTTP and an authorized token
    Then I should receive an authorized response from HTTP

  Scenario: Reject a request without a token
    When I request configured databases with HTTP and no token
    Then I should receive an unauthenticated response from HTTP

  Scenario: Reject a request with an invalid token
    When I request configured databases with HTTP and an invalid token
    Then I should receive an unauthenticated response from HTTP

  Scenario: Reject a request from an unauthorized subject
    When I request configured databases with HTTP and an unauthorized token
    Then I should receive a forbidden response from HTTP
