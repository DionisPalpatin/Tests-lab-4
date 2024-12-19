Feature: Login with 2FA (email, messenger, etc)

Scenario: User login with 2FA
    When user send "POST" request on "/login"
    Then response code should be 200
    And JSON response on /login should match:
        """
        {
            "message": "Code send to email successfully"
        }
        """

    And user send "POST" request on "/login/verify"
    Then response code should be 200
    And JSON response on /login/verify should match:
        """
        {
            "message": "Successfully sign in"
        }
        """