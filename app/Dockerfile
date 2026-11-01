FROM golang:1.24-alpine AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/bank-transfer-api \
    ./cmd/server


FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /out/bank-transfer-api /app/bank-trnsfer-api

EXPOSE 8080
 
USER nonroot:nonroot

ENTRYPOINT ["/app/bank-transfer-api"]
