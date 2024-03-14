
## iter 1

metricstest -test.v -test.run="^TestIteration1$" \
            -binary-path=cmd/server/server

### iter 2 

metricstest -test.v -test.run="^TestIteration2[AB]*$" \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent

### iter 3

metricstest -test.v -test.run="^TestIteration3[AB]*$" \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server

### iter 4

SERVER_PORT=8085
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration4$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.


### iter 5

SERVER_PORT=8081
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration5$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.


### iter 6

SERVER_PORT=8081
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration6$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.

### iter 7

SERVER_PORT=8081
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration7$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.

### iter 8

SERVER_PORT=8080
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./temp
          metricstest -test.v -test.run="^TestIteration8$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -server-port=$SERVER_PORT \
            -source-path=.

### iter 9

 SERVER_PORT=8080
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./metrics.json
          metricstest -test.v -test.run="^TestIteration9$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -file-storage-path=$TEMP_FILE \
            -server-port=$SERVER_PORT \
            -source-path=.

### iter 10

 SERVER_PORT=8080
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./metrics.json
          metricstest -test.v -test.run="^TestIteration10[AB]$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://postgres:1@localhost:5432/video?sslmode=disable' \
            -server-port=$SERVER_PORT \
            -source-path=.

### iter 11

SERVER_PORT=8080
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./metrics.json
          metricstest -test.v -test.run="^TestIteration11$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://postgres:1@localhost:5432/video?sslmode=disable' \
            -server-port=$SERVER_PORT \
            -source-path=.

### iter 12

 SERVER_PORT=8080
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./metrics.json
          metricstest -test.v -test.run="^TestIteration12$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://postgres:1@localhost:5432/video?sslmode=disable' \
            -server-port=$SERVER_PORT \
            -source-path=.

### iter 13

SERVER_PORT=8080
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./metrics.json
          metricstest -test.v -test.run="^TestIteration13$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://postgres:1@localhost:5432/video?sslmode=disable' \
            -server-port=$SERVER_PORT \
            -source-path=.


### iter 14
SERVER_PORT=8080
          ADDRESS="localhost:${SERVER_PORT}"
          TEMP_FILE=./metrics.json
          metricstest -test.v -test.run="^TestIteration14$" \
            -agent-binary-path=cmd/agent/agent \
            -binary-path=cmd/server/server \
            -database-dsn='postgres://postgres:1@localhost:5432/video?sslmode=disable' \
            -key="${TEMP_FILE}" \
            -server-port=$SERVER_PORT \
            -source-path=.


### coverage

go test ./... -coverprofile cover.out && go tool cover -func cover.out && go tool cover -html cover.out

### build all

go build -o cmd/server/server cmd/server/*.go && go build -o cmd/agent/agent cmd/agent/*.go

### generate protoc

protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  proto/metrics.proto


### mockgen

mockgen -destination=internal/server/repository/mockrepo/mockrepo.go -package=mockrepo github.com/benderr/metrics/internal/server/repository MetricRepository