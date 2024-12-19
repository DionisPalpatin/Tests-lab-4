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
	expectRequest *httpexpect.Expect
)

func TestLogin(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	expectRequest = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Client:   &http.Client{},
		Reporter: httpexpect.NewRequireReporter(t),
	})

	suite := godog.TestSuite{
		ScenarioInitializer: testLoginScenarioInitializer,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"./features/test_login.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("Some errors occurred while running test (login)")
	}
}

func testLoginScenarioInitializer(ctx *godog.ScenarioContext) {
	login(ctx)
}

func login(ctx *godog.ScenarioContext) {
	var response *httpexpect.Response

	ctx.Step(`^User send "([^"]*)" request to "([^"]*)"$`, func(method, endpoint string) error {
		response = expectRequest.Request(method, endpoint).
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
		response = expectRequest.Request(method, endpoint).
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
