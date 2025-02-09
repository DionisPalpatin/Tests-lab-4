package e2e_tests

import (
	"github.com/joho/godotenv"
	"github.com/pquerna/otp/totp"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/gavv/httpexpect/v2"
)

var (
	expectReset      *httpexpect.Expect
	testPasswordCode string
)

func TestResetPassword(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	testPasswordCode, err = totp.GenerateCode(os.Getenv("TOTP_SECRET"), time.Now())
	if err != nil {
		t.Fatalf("Failed to generate TOTP code: %v", err)
	}

	client := &http.Client{}
	expectReset = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeResetPasswordScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/test_password_reset.feature"},
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
			"message": "Enter the TOTP code from your app",
		})
		return nil
	})

	ctx.Step(`^user send "([^"]*)" request to "([^"]*)"$`, func(method, endpoint string) error {
		response = expectReset.Request(method, endpoint).
			WithJSON(map[string]string{
				"email":        os.Getenv("USER_EMAIL"),
				"code":         testPasswordCode,
				"new_password": os.Getenv("NEW_USER_PASSWORD"),
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
