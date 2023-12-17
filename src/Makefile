SRC=cmd/cutter/main.go
TESTS=tests/cutter_test.go

export CONFIG_PATH=config/local.yaml

all: clean
	go run $(SRC)

test: clean
	go run $(SRC) &
	cd tests && go run gotest.tools/gotestsum@latest --format testname
	cd cmd && bash cleanports.sh

clean:
	cd cmd && bash cleanports.sh
