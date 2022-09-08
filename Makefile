APP = party-day-bot
SERVICE_PATH = github.com/Dsmit05
BR = `git rev-parse --symbolic-full-name --abbrev-ref HEAD`
VER = `git describe --tags --abbrev=0`
TIMESTM = `date -u '+%Y-%m-%d_%H:%M:%S%p'`
FORMAT = $(VER)-$(TIMESTM)
DOCTAG = $(VER)-$(BR)

.PHONY: build
build:
	CGO_ENABLED=0 go build -o $(APP) -ldflags "-X $(SERVICE_PATH)/$(APP)/internal/config.buildVersion=$(FORMAT)" cmd/bot/main.go

.PHONY: build-image
build-image:
	sudo docker build -t $(APP):$(DOCTAG) .

.PHONY: run-app
run-app:
	docker run -d --name=$(APP)-$(VER) -p 8082:8082 -p 8082:8082 $(APP):$(DOCTAG)

.PHONY: del-app
del-app:
	docker rm $(APP)-$(VER)

.PHONY: recreate-bot
recreate-api:
	docker-compose up -d --force-recreate --no-deps --build bot

.PHONY: cover
cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: bench
bench:
	go test -bench=. -benchmem -benchtime=5x ./...

.PHONY: mock
mock:
	mockgen -source=internal/api/interface.go \
	-destination=internal/api/mocks/mock_for_api.go