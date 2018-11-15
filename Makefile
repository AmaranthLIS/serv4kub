APP?=serv4kub
PORT?=8000

RELEASE?=0.0.1
# COMMIT?=$(shell git rev-parse --short HEAD)
COMMIT?=33301
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROJECT?=play4j/serv4kub
USERNAME?=lis7
CONTAINER_IMAGE_NAME?=${USERNAME}/${APP}
CONTAINER_IMAGE?=docker.io/${CONTAINER_IMAGE_NAME}


GOOS?=linux
GOARCH?=amd64



# ----- ----- ----- ----- -----
clean:
	rm -f ${APP}

build: clean
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} \
	go build \
		-ldflags "-s -w -X ${PROJECT}/version.Release=${RELEASE} \
		-X ${PROJECT}/version.Commit=${COMMIT} -X ${PROJECT}/version.BuildTime=${BUILD_TIME}" \
		-o ${APP}

#run: build
#	PORT=${PORT} ./${APP}

container: build
	docker build -t $(CONTAINER_IMAGE_NAME) .

run: container
	docker stop $(CONTAINER_IMAGE_NAME) || true && docker rm $(CONTAINER_IMAGE_NAME) || true
	docker run --name ${APP} -p ${PORT}:${PORT} --rm \
		-e "PORT=${PORT}" \
		$(CONTAINER_IMAGE_NAME)


test:
	go test -v -race ./...


push: container
	docker push $(CONTAINER_IMAGE_NAME)

minikube: push
	for t in $(shell find ./kubernetes/serv4kub -type f -name "*.yaml"); do \
        cat $$t | \
        	gsed -E "s/\{\{(\s*)\.Release(\s*)\}\}/$(RELEASE)/g" | \
        	gsed -E "s/\{\{(\s*)\.ServiceName(\s*)\}\}/$(APP)/g"; \
        echo ---; \
    done > tmp.yaml
	kubectl apply -f tmp.yaml