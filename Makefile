.PHONY: dev
dev:
	cd cmd/api/ && \
		PG_DB_CONN_STRING="postgres://weather:weather@localhost:6432/weather?sslmode=disable" \
		PORT=8080 \
		go run *.go

.PHONY: tests
tests:
	go test -v ./...

.PHONY: integration_tests
integration_tests:
	go test -tags=integration -run=Integration ./...
