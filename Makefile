include .env

server:
	go build -o cmd/server/server cmd/server/main.go

agent:
	go build -o cmd/agent/agent cmd/agent/main.go

build:
	make server
	make agent

test_1: build
	./metricstest -test.v -test.run=^TestIteration1$$ -binary-path=cmd/server/server

test_2: test_1
	./metricstest -test.v -test.run=^TestIteration2[AB]*$$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent

test_3: test_2
	./metricstest -test.v -test.run=^TestIteration3[AB]*$$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server

test_4: test_3
	./metricstest -test.v -test.run=^TestIteration4$$ \
	-agent-binary-path=cmd/agent/agent \
	-binary-path=cmd/server/server \
	-server-port=${SERVER_PORT} \
	-source-path=.

test_5: test_4
	./metricstest -test.v -test.run=^TestIteration5$$\
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=${SERVER_PORT} \
            -source-path=.

test_6: test_5
	./metricstest -test.v -test.run=^TestIteration6$$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=${SERVER_PORT} \
            -source-path=.

test_7: test_6
	./metricstest -test.v -test.run=^TestIteration7$$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=${SERVER_PORT} \
            -source-path=.

test_8: test_7
	./metricstest -test.v -test.run=^TestIteration8$$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=${SERVER_PORT} \
            -source-path=.

test_9: test_8
	./metricstest -test.v -test.run=^TestIteration9$$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -file-storage-path=${TEMP_FILE} \
            -server-port=${SERVER_PORT} \
            -source-path=.

test_10: test_9
	./metricstest -test.v -test.run=^TestIteration10[AB]$$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://monogo:monogo@localhost:5432/monogodb?sslmode=disable' \
            -server-port=${SERVER_PORT} \
            -source-path=.

test_11: test_10
	./metricstest -test.v -test.run=^TestIteration11$$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://monogo:monogo@localhost:5432/monogodb?sslmode=disable' \
            -server-port=${SERVER_PORT} \
            -source-path=.
test_12: test_11
	./metricstest -test.v -test.run=^TestIteration12$$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://monogo:monogo@localhost:5432/monogodb?sslmode=disable' \
            -server-port=${SERVER_PORT} \
            -source-path=.

test_13: test_12
	./metricstest -test.v -test.run=^TestIteration13$$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://monogo:monogo@localhost:5432/monogodb?sslmode=disable' \
            -server-port=${SERVER_PORT} \
            -source-path=.

test_14: test_13
	./metricstest -test.v -test.run=^TestIteration14$$ \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://monogo:monogo@localhost:5432/monogodb?sslmode=disable' \
            -key="${TEMP_FILE}" \
            -server-port=${SERVER_PORT} \
            -source-path=.

db_init:
	docker run -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_DB=${POSTGRES_DB} -e POSTGRES_PORT=${POSTGRES_PORT} -d -p 5432:5432 postgres

run_serv:
	ADDRESS="${ADDRESS}" LOG_LEVEL="${LOG_LEVEL}" LOG_PATH="${LOG_PATH}" DATABASE_DSN="${DATABASE_DSN}" go run cmd/server/main.go

run_agent:
	ADDRESS="${ADDRESS}" REPORT_INTERVAL="${REPORT_INTERVAL}" POLL_INTERVAL="${POLL_INTERVAL}" go run cmd/agent/main.go

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func coverage.out
	rm -rf coverage.out

profile:
	curl http://localhost:8080/debug/pprof/heap\?seconds\=30 -o profiles/rename_me.pprof

generate_certificate:
	./make_cert.sh

.PHONY: build server agent