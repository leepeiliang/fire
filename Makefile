.PHONY: all
all: fire

.PHONY: fire
fire:
	GOOS=linux go build -o bin/fire cmd/app/server.go
.PHONY: fire_image
fire_image:fire
	sudo docker build -t harbor.dev.21vianet.com/metaedge/fire:v1.3.10 .
	docker push harbor.dev.21vianet.com/metaedge/fire:v1.3.10
clean:
	rm -f  ./bin/fire
