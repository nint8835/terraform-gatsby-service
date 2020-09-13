FROM golang:1.14-alpine AS builder

WORKDIR /build
COPY . /build
RUN go build

FROM alpine
WORKDIR /app
# ARG gin_mode=release
# ENV GIN_MODE=${gin_mode}
COPY --from=builder /build/terraform-gatsby-service .
RUN addgroup -S terraform && \
    adduser -S terraform -G terraform && \
    chown -R terraform /app
ENTRYPOINT [ "/app/terraform-gatsby-service" ]
