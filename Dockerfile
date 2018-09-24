FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
ADD ./vendor /go/src/
RUN go build -o main .
FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
EXPOSE 8000
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
