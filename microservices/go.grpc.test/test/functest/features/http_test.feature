Feature: Testing the REST api

  Scenario Outline: run a simple http test
    Given the backend is running
    and the frontend is running
    and the nginx is running
    When I request <name>
    Then the <message> is returned

    Examples: Cities
      | name   | message               |
      | Venice | Greeting: Test Venice |
      | London | Greeting: Test London |
    
    Examples: Animals
      | name | message            |
      | Dog  | Greeting: Test Dog |
      | Cat  | Greeting: Test Cat |
