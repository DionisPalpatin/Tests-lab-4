package e2e_tests

import (
	"net/http"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/gavv/httpexpect/v2"
)

var (
	expectReset *httpexpect.Expect
)

func TestResetPassword(t *testing.T) {
	client := &http.Client{}
	expectReset = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(nil),
	})

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeResetPasswordScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/reset.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run reset password feature tests")
	}
}

func resetPasswordWith2FA(ctx *godog.ScenarioContext) {
	var response *httpexpect.Response

	ctx.Step(`^User send "([^"]*)" request to "([^"]*)"$`, func(method, endpoint string) error {
		response = expectReset.Request(method, endpoint).
			WithJSON(map[string]string{
				"email":    os.Getenv("USER_EMAIL"),
				"password": os.Getenv("USER_PASSWORD"),
			}).
			Expect()
		return nil
	})

	ctx.Step(`^the response on /password/reset code should be (\d+)$`, func(statusCode int) error {
		response.Status(statusCode)
		return nil
	})

	ctx.Step(`^the response on /password/reset should match json:$`, func(expectedJSON *godog.DocString) error {
		response.JSON().Object().IsEqual(map[string]interface{}{
			"message": "Reset password code send successfully",
		})
		return nil
	})

	ctx.Step(`^user send "([^"]*)" request to "([^"]*)"$`, func(method, endpoint string) error {
		response = expectReset.Request(method, endpoint).
			WithJSON(map[string]string{
				"email":       os.Getenv("USER_EMAIL"),
				"code":        os.Getenv("TEST_CODE"),
				"newPassword": os.Getenv("NEW_USER_PASSWORD"),
			}).Expect()
		return nil
	})

	ctx.Step(`^the response on /password/reset/verify code should be (\d+)$`, func(statusCode int) error {
		response.Status(statusCode)
		return nil
	})

	ctx.Step(`^the response on /password/reset/verify should match json:$`, func(expectedJSON *godog.DocString) error {
		response.JSON().Object().IsEqual(map[string]interface{}{
			"message": "Password is changed successfully",
		})
		return nil
	})
}

func InitializeResetPasswordScenario(ctx *godog.ScenarioContext) {
	resetPasswordWith2FA(ctx)
}
