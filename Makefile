GO_FILES := $(shell ls *.go | grep -v '_test.go')
TEST_FILES := $(shell ls *_test.go)
MAIN_FILE := main.go

.PHONY: all run test clean restart-redis

all: run

stop-redis:
	sudo systemctl stop redis

restart-redis:
	sudo systemctl start redis

run: stop-redis
	go run $(GO_FILES)

test:
	go test -v

clean: restart-redis
	go clean
	rm -f *.aof
