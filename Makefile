appb:
	@echo "Building app....."
	@cd zendesk-app/ && zat package && open ./tmp
appr:
	@echo "Clean TMP folder"
	@cd zendesk-app/ && zat clean
run:
	@go get -v && GOOS=linux && go build && go run server.go
run-docker: 
	docker rmi -f zendesk-integration:1.0
	docker-compose up -d