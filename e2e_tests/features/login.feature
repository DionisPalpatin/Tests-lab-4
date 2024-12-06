Feature: Login with 2FA

Scenario: User login with 2FA
  When User send "POST" request to "/login"
  Then the response on /login code should be 200
  And the response on /login should match json:
      """
      {
        "message": "Auth email send successfully"
      }
      """
  And user send "POST" request to "/login/verify"
  Then the response on /login/verify code should be 200
  And the response on /login/verify should match json:
      """
      {
        "message": "Auth is made successfully"
      }
      """