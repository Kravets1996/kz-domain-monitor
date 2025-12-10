FROM golang:1.25.5 AS build

WORKDIR /code
COPY ./ /code

RUN CGO_ENABLED=0 GOOS=linux go build -o kz-domain-monitor

FROM alpine:latest

RUN apk update && apk --no-cache add tzdata

WORKDIR /app
COPY --from=build /code/kz-domain-monitor /app/kz-domain-monitor

CMD ["/app/kz-domain-monitor"]