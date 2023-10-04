.PHONY: dev
dev:
	cd cmd/api/ && \
		PG_DB_CONN_STRING="postgres://weather:weather@localhost:6432/weather?sslmode=disable" \
		PORT=8080 \
		go run *.go

.PHONY: tests
tests:
	go test -v ./...

.PHONY: build_images
build_images:
	docker build -f Dockerfile.api -t go-sample-rest/weather-service .
	docker build -f Dockerfile.flyway -t go-sample-rest/flyway .

.PHONY: integration_tests
integration_tests:
	go test -tags=integration -run=Integration ./...

.PHONY: start_minikube
start_minikube: build_images
	minikube start
	minikube image load go-sample-rest/weather-service
	minikube image load go-sample-rest/flyway
	kubectl apply --context minikube -f minikube.yaml

.PHONY: stop_minikube
stop_minikube:
	minikube stop
	minikube delete
