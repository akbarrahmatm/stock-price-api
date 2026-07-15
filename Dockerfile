FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod ./
RUN go mod download || true
COPY *.go ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /stock-api .

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /stock-api /stock-api
EXPOSE 8000
ENTRYPOINT ["/stock-api"]
