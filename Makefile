.PHONY: build

build:
	sam build

api: build
	sam local start-api