FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && adduser -D -H appuser

WORKDIR /app
COPY --from=build /out/server /app/server

USER appuser
EXPOSE 8080

ENTRYPOINT ["/app/server"]

