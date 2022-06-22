PKGS=$(shell go list ./...)

installdeps:
	go mod vendor && go mod tidy

test: installdeps
	@set -e; go test -v ${PKGS}

