
run:
	docker run -it --rm -v "$(PWD):/hacking" -w "/hacking" -p "5000:5000" -e "GOCACHE=/tmp/" -u "1000:1000" golang:latest go run main.go

test:
	docker run -it --rm -v "$(PWD):/hacking" -w "/hacking" golang:latest go test -v ./...


bash:
	docker run -it --rm -v "$(PWD):/hacking" -w "/hacking" golang:latest /bin/bash

loadtestlayers:
	docker run -it --rm -d --name registry -p 5000:5000 registry:2 
	docker push 127.0.0.1:5000/fedora:latest
	docker cp registry:/var/lib/registry/docker/registry/v2/ ./tmp-data

cleartest:
	docker rm -f registry
	rm -rf ./tmp-data || true

clearregistrydata:
	sudo rm -rf ./registry-data

runtestregistry:
	docker pull fedora:32 
	docker tag fedora:32 127.0.0.1:5000/fedora:latest
	docker run -it --rm --name registry -p 5000:5000 registry:2 
# 	docker push 127.0.0.1:5000/fedora:latest

pushtestregistry:
	docker push 127.0.0.1:5000/fedora:latest