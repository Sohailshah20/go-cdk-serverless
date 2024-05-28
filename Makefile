build:
	@cd lambda && go build -o bootstrap
	@cd lambda && zip function.zip bootstrap

.PHONY: build