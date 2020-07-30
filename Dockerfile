FROM golang:1.14-alpine AS builder

WORKDIR /build
COPY . /build
RUN go build

FROM alpine
WORKDIR /app
COPY --from=builder /build/terraform-gatsby-service .
ENTRYPOINT [ "/app/terraform-gatsby-service" ]
