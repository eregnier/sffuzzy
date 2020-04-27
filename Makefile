test:
	@export DEBUG=0 && go test

test-trace:
	@clear && export DEBUG=1 && go test