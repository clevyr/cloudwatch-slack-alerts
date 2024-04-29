#syntax=docker/dockerfile:1.7

FROM --platform=$BUILDPLATFORM golang:1.22.2-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY *.go ./
# Set Golang build envs based on Docker platform string
ARG TARGETPLATFORM
RUN --mount=type=cache,target=/root/.cache <<EOT
  set -eux

  case "$TARGETPLATFORM" in
    'linux/amd64') export GOARCH=amd64 ;;
    'linux/arm/v6') export GOARCH=arm GOARM=6 ;;
    'linux/arm/v7') export GOARCH=arm GOARM=7 ;;
    'linux/arm64') export GOARCH=arm64 ;;
    *) echo "Unsupported target: $TARGETPLATFORM" && exit 1 ;;
  esac

  export CGO_ENABLED=0
  go build -ldflags='-w -s' -tags lambda.norpc -trimpath -o cloudwatch-slack-alerts .
EOT

FROM alpine:3.19 AS rie
WORKDIR /app
ARG TARGETPLATFORM
ARG RIE_VERSION=v1.18
RUN <<EOT
  set -eux

  case "$TARGETPLATFORM" in
    'linux/amd64') export SUFFIX=x86_64 ;;
    'linux/arm64') export SUFFIX=arm64 ;;
    *) echo "Unsupported target: $TARGETPLATFORM" && exit 1 ;;
  esac

  wget \
    -O aws-lambda-rie \
    "https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/download/$RIE_VERSION/aws-lambda-rie-$SUFFIX"
  chmod +x aws-lambda-rie
EOT

FROM gcr.io/distroless/static:nonroot AS base
WORKDIR /
COPY --from=builder /app/cloudwatch-slack-alerts .
ENTRYPOINT ["./cloudwatch-slack-alerts"]

FROM base AS local
COPY --from=rie /app/aws-lambda-rie .
ENTRYPOINT ["./aws-lambda-rie"]
CMD ["./cloudwatch-slack-alerts"]

FROM base
