SRC=cmd/cutter/main.go
TESTS=tests/cutter_test.go

all:
	go run $(SRC)

test:
	cd cmd && bash cleanports.sh
	go run $(SRC) &
	cd tests && go run gotest.tools/gotestsum@latest --format testname
	cd cmd && bash cleanports.sh
