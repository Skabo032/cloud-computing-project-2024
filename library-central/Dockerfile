FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o library-central

FROM scratch

COPY --from=builder /app/library-central /library-central

EXPOSE 8080

ENV PORT=:8080

ENTRYPOINT ["/library-central"]