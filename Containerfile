FROM golang:1.21-alpine AS builder

# Set us up a non-root user
ENV USER=scheduler
ENV UID=1000
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

# Set us up the build environment
WORKDIR $GOPATH/src/github.com/supertylerc/scheduler/
COPY pkg ./pkg
COPY internal ./internal
COPY cmd ./cmd
COPY go.mod ./
COPY go.sum ./
COPY vendor ./vendor
RUN GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-w -s" -o /go/bin/scheduler cmd/scheduler/main.go

FROM gcr.io/distroless/static
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /go/bin/scheduler /go/bin/scheduler
USER scheduler:scheduler
ENTRYPOINT ["/go/bin/scheduler"]
