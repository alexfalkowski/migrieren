@startup
Feature: gRPC API authentication
  The gRPC API requires a verified token and an authorized subject.

  Scenario: Accept a request with an authorized token
    When I request configured databases with gRPC and an authorized token
    Then I should receive an authorized response from gRPC

  Scenario: Reject a request without a token
    When I request configured databases with gRPC and no token
    Then I should receive an unauthenticated response from gRPC

  Scenario: Reject a request with an invalid token
    When I request configured databases with gRPC and an invalid token
    Then I should receive an unauthenticated response from gRPC

  Scenario: Reject a request from an unauthorized subject
    When I request configured databases with gRPC and an unauthorized token
    Then I should receive a forbidden response from gRPC
