.PHONY: local-run docker-run docker-rm run-e2e


local-run:
	go run cmd/main.go


docker-run:
	docker compose up -d


docker-rm:
	docker compose down
	docker image rm testing-auth-app alpine:latest


run-e2e: docker-run
	local-run
	go test e2e_tests/login_test.go && go test e2e_tests/reset_password_test.go
