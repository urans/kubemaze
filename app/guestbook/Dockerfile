FROM golang:1.24 AS builder
ARG TARGETOS
ARG TARGETARCH
COPY . /workspace
WORKDIR /workspace
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -a -o manager cmd/main.go


FROM alpine:3.21
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532
ENTRYPOINT ["/manager"]
