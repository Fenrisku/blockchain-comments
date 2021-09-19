# .PHONY: all dev clean build env-up env-down run
#export BYFN_CA1_PRIVATE_KEY=$(cd fixtures/network-base/crypto-config/peerOrganizations/org1.example.com/ca && ls *_sk) 
#export BYFN_CA2_PRIVATE_KEY=$(cd fixtures/network-base/crypto-config/peerOrganizations/org2.example.com/ca && ls *_sk)



all: build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd fixtures/network-base && docker-compose -f docker-compose-cli.yaml -f docker-compose-ca.yaml -f docker-compose-couch.yaml -f docker-compose-etcdraft2.yaml up --force-recreate -d

	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd fixtures/network-base && docker-compose -f docker-compose-cli.yaml -f docker-compose-ca.yaml -f docker-compose-couch.yaml -f docker-compose-etcdraft2.yaml down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./comments

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@docker volume prune
	@rm -rf /tmp/example-* example
	@docker rm -f -v `docker ps -a --no-trunc | grep "example" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "example" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"
	

