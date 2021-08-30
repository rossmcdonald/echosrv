echosrv : main.go
	go build -o echosrv main.go

deps :
	go mod download

image :
	docker build --squash -t rossmcd/echo-json:latest .
	docker tag rossmcd/echo-json:latest rossmcd/echo-json:$(shell git rev-parse --short HEAD)

push :
	docker push rossmcd/echo-json:latest
	docker push rossmcd/echo-json:$(shell git rev-parse --short HEAD)

clean :
	rm echosrv
