echo : main.go
	go build -o echo main.go

image :
	docker build --squash -t rossmcd/echo-json:latest .

push :
	docker push rossmcd/echo-json:latest

clean :
	rm echo
