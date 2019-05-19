ANSI_RED="\033[0;31m"
ANSI_GREEN="\033[0;32m"
ANSI_BLUE="\033[0;34m"
ANSI_RESET="\033[0m"

ifneq ("$(wildcard /usr/local/opt/coreutils/libexec/gnubin/echo)","")
	ECHO="/usr/local/opt/coreutils/libexec/gnubin/echo"
else
	ECHO="/bin/echo"
endif

.PHONY: all test test-html docs

all: libreLIS.bin
docs:
	xdg-open "http://localhost:6060/pkg/github.com/gjvnq/LibreLIS/libLIS"
docs-server:
	godoc -http=:6060
test: coverage.out
test-html: coverage.out
	@$(ECHO) -e $(ANSI_GREEN)"Generating coverage report..."$(ANSI_RESET)
	go tool cover -html=coverage.out
	@$(ECHO) -e $(ANSI_BLUE)"Finished target"$(ANSI_RESET)
	
libLIS/libLIS.a: libLIS/*.go libLIS/Makefile
	cd libLIS && make libLIS.a

libreLIS.bin: *.go libLIS/libLIS.a
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Fixing imports..."$(ANSI_RESET)
	goimports -w .
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Formatting code..."$(ANSI_RESET)
	go fmt
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Compiling code..."$(ANSI_RESET)
	go build -o $@
	@$(ECHO) -e $(ANSI_BLUE)"["$@"] Finished target $@"$(ANSI_RESET)

coverage.out: *.go libLIS/libLIS.a
	@$(ECHO) -e $(ANSI_GREEN)"["$@"] Testing code..."$(ANSI_RESET)
	go test -cover -coverprofile=coverage.out
	@$(ECHO) -e $(ANSI_BLUE)"["$@"] Finished target"$(ANSI_RESET)

run: libreLIS.bin
	./libreLIS.bin

run-dev: libreLIS.bin
	#go get github.com/codegangsta/gin
	gin --all --port 8080 --appPort 8081 run "./libreLIS.bin --dev"