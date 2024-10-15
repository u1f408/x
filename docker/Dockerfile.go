FROM docker.io/alpine:edge AS builder
WORKDIR /build
RUN apk add --no-cache build-base go

ARG BUILD_VERSION
COPY . .
RUN go mod download
RUN make VERSION=${BUILD_VERSION}

FROM scratch
COPY --from=builder /build/bin/* /
