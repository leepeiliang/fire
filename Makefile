.PHONY: all
all: fire

.PHONY: fire
fire:
	GOOS=linux go build -o bin/fire cmd/app/server.go
.PHONY: fire_image
fire_image:fire
	sudo docker build -t harbor.meta42.indc.vnet.com/metaedge/fire:v1.3.10 .
	docker push harbor.meta42.indc.vnet.com/metaedge/fire:v1.3.10
clean:
	rm -f  ./bin/fire
