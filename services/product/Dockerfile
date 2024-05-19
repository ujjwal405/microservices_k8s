FROM golang:alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin ./server


FROM scratch AS production
COPY --from=builder /app/bin /bin/app
EXPOSE 8080
ENTRYPOINT ["/bin/app"]