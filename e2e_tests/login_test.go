package e2e_tests

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/gavv/httpexpect/v2"
)

var (
	expectLogin *httpexpect.Expect
)

func TestLogin(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	client := &http.Client{}
	expectLogin = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://auth-app:8080",
		Client:   client,
		Reporter: httpexpect.NewRequireReporter(t),
	})

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeLoginScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"./features/login.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run login feature tests")
	}
}

func loginWith2FA(ctx *godog.ScenarioContext) {
	var response *httpexpect.Response

	ctx.Step(`^User send "([^"]*)" request to "([^"]*)"$`, func(method, endpoint string) error {
		response = expectLogin.Request(method, endpoint).
			WithJSON(map[string]string{
				"email":    os.Getenv("USER_EMAIL"),
				"password": os.Getenv("USER_PASSWORD"),
			}).
			Expect()
		return nil
	})

	ctx.Step(`^the response on /login code should be (\d+)$`, func(statusCode int) error {
		response.Status(statusCode)
		return nil
	})

	ctx.Step(`^the response on /login should match json:$`, func(expectedJSON *godog.DocString) error {
		response.JSON().Object().IsEqual(map[string]interface{}{
			"message": "Auth email send successfully",
		})
		return nil
	})

	ctx.Step(`^user send "([^"]*)" request to "([^"]*)"$`, func(method, endpoint string) error {
		response = expectLogin.Request(method, endpoint).
			WithJSON(map[string]string{
				"email": os.Getenv("USER_EMAIL"),
				"code":  os.Getenv("TEST_CODE"),
			}).Expect()
		return nil
	})

	ctx.Step(`^the response on /login/verify code should be (\d+)$`, func(statusCode int) error {
		response.Status(statusCode)
		return nil
	})

	ctx.Step(`^the response on /login/verify should match json:$`, func(expectedJSON *godog.DocString) error {
		response.JSON().Object().IsEqual(map[string]interface{}{
			"message": "Auth is made successfully",
		})
		return nil
	})
}

func InitializeLoginScenario(ctx *godog.ScenarioContext) {
	loginWith2FA(ctx)
}
