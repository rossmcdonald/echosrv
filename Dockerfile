FROM alpine

RUN apk add go make

RUN mkdir /app

ADD main.go /app/
ADD Makefile /app/

WORKDIR /app

RUN make echo 
RUN apk del go make

EXPOSE 8888

ENTRYPOINT ["./echo"]
