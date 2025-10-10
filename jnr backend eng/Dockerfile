FROM golang:1.20-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /delivery cmd/server/main.go

FROM alpine:3.18
COPY --from=build /delivery /delivery
EXPOSE 8080
ENTRYPOINT ["/delivery"]
