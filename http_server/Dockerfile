FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY go.mod ./
# fetch module metadata first (creates go.sum if present)
RUN go mod download
# now bring in source and let tidy add any missing deps to go.sum
COPY main.go .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /out/enclave-http main.go

FROM alpine:3.20
RUN adduser -D enclave && apk add --no-cache ca-certificates
USER enclave
COPY --from=build /out/enclave-http /enclave-http
EXPOSE 8000
CMD ["/enclave-http"]