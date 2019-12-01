
start-locakstack:
	@docker-compose up -d

test: start-locakstack
	@cd src/backend/persons && \
	go test -v ./...
