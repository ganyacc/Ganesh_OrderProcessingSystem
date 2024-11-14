# Build stage
FROM golang:1.21.0 AS build-stage

WORKDIR /app




# Copy .env file and all Go source code
COPY config/ ./config/
COPY database/ ./database/
COPY entities/ ./entities/
COPY handler/ ./handler/
COPY logger/ ./logger/
COPY pkg/ ./pkg/
COPY repository/ ./repository/
COPY server/ ./server/
COPY testCases/ ./testCases/
COPY config.yaml ./
COPY *.go ./
COPY go.mod go.sum ./

#COPY *.sh ./


RUN go mod download




# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Test stage
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Production stage
FROM alpine:latest AS build-release-stage

WORKDIR /app





COPY --from=build-stage /app/config/ ./config/
COPY --from=build-stage /app/database/ ./database/
COPY --from=build-stage /app/entities/ ./entities/
COPY --from=build-stage /app/handler/ ./handler/
COPY --from=build-stage /app/logger/ ./logger/
COPY --from=build-stage /app/repository/ ./repository/
COPY --from=build-stage /app/server/ ./server/
COPY --from=build-stage /app/testCases/ ./testCases/

COPY --from=build-stage /app/config.yaml ./
COPY --from=build-stage /app/main /main

EXPOSE 8080

#HEALTHCHECK --interval=1m --timeout=10s CMD curl --fail http://localhost:8080/healthcheck || exit 1


ENTRYPOINT ["/main"]
