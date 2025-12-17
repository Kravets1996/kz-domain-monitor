FROM golang:1.25.5 AS build

ARG APP_VERSION=dev

WORKDIR /code
COPY ./ /code

RUN CGO_ENABLED=0 GOOS=linux go build -o kz-domain-monitor  -ldflags "-X 'main.Version=${APP_VERSION}'"

FROM alpine:latest

RUN apk update && apk --no-cache add tzdata

WORKDIR /app
COPY --from=build /code/kz-domain-monitor /app/kz-domain-monitor
ENV PATH="/app:${PATH}"

CMD ["/app/kz-domain-monitor"]