FROM golang:1 AS builder

RUN mkdir /app
ADD . /app/
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/web-resume /app/cmd/web-resume/main.go

FROM alpine:latest

RUN mkdir -p /usr/local/web-resume/config
RUN mkdir -p /usr/local/web-resume/static
RUN mkdir -p /usr/local/web-resume/templates

COPY --from=builder /app/bin/web-resume /usr/local/bin/web-resume
COPY ./static /usr/local/web-resume/static
COPY ./templates /usr/local/web-resume/templates

VOLUME /usr/local/webresume/config

EXPOSE 8080

ENV CONFIG_PATH=/usr/local/web-resume/config
ENV GIN_MODE=release
ENV LOG_FORMAT=json
ENV LOG_LEVEL=info
ENV PORT=8080
ENV STATIC_PATH=/usr/local/web-resume/static
ENV TEMPLATES_PATH=/usr/local/web-resume/templates

CMD ["/usr/local/bin/web-resume"]