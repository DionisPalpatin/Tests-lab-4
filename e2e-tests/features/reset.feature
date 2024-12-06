Feature: Reset password with 2FA

Scenario: User reset password with 2FA
  When User send "POST" request to "/password/reset"
  Then the response on /password/reset code should be 200
  And the response on /password/reset should match json:
      """
      {
        "message": "Reset password code send successfully"
      }
      """
  And user send "POST" request to "/password/reset/verify"
  Then the response on /password/reset/verify code should be 200
  And the response on /password/reset/verify should match json:
      """
      {
        "message": "Password is changed successfully"
      }
      """