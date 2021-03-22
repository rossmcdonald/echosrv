echo : main.go
	go build -o echo main.go

image :
	docker build --squash -t rossmcd/echo-json:latest .
	docker tag rossmcd/echo-json:latest rossmcd/echo-json:$(git rev-parse --short HEAD)

push :
	docker push rossmcd/echo-json:latest
	docker push rossmcd/echo-json:$(shell git rev-parse --short HEAD)

clean :
	rm echo
