FROM alpine

RUN apk add go make git

RUN mkdir /app /go
ENV GOPATH=/go

ADD main.go go.mod go.sum /app/
ADD Makefile /app/

WORKDIR /app

RUN make deps
RUN make
RUN apk del go make git && rm -rf /go && rm -rf /root/.cache

EXPOSE 8889

ENTRYPOINT ["./echosrv"]
