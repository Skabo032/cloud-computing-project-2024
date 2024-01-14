FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o library-city

FROM scratch

COPY --from=builder /app/library-city /library-city

EXPOSE 8080

ENV PORT=:8080

ENTRYPOINT ["/library-city"]