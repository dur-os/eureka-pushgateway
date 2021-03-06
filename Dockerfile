FROM golang:1.11.0 as builder

WORKDIR /go/pad
ADD ./ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags  "-s -w" -o "eureka-pushgateway"

FROM alpine:latest
#RUN apt-get update && apt-get install -y wget

#FROM prom/pushgateway:v0.4.0
WORKDIR /pushgateway

COPY --from=builder /go/pad/eureka-pushgateway ./
ADD start.sh ./

RUN chmod +x start.sh

EXPOSE 9092

CMD [ "./start.sh"]