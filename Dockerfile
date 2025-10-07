FROM golang:1.25-alpine AS builder
RUN apk add --no-cache ca-certificates
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

# ---- Distroless Runtime Stage ----
FROM gcr.io/distroless/static:nonroot
# Copy CA certs and the binary
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/main /app/main
WORKDIR /app
USER nonroot:nonroot
ENTRYPOINT ["/app/main"]