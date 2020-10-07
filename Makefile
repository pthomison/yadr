
run:
	docker run \
	-it --rm \
	-v "$(PWD):/hacking" \
	-w "/hacking" \
	-p "5000:5000" \
	-e "GOCACHE=/tmp/" \
	-u "1000:1000" \
	golang:latest \
	go run main.go

runclean:
	clear
	rm -rf ./data
	docker run \
	-it --rm \
	-v "$(PWD):/hacking" \
	-w "/hacking" \
	-p "5000:5000" \
	-e "GOCACHE=/tmp/" \
	-u "1000:1000" \
	golang:latest \
	go run main.go

delve:
	docker build . -t delve-image
	docker run \
	-it --rm \
	--entrypoint='' \
	-v "$(PWD):/go/src/github.com/pthomison/yadr" \
	-w "/go/src/github.com/pthomison/yadr" \
	-p "5000:5000" \
	-e "GOCACHE=/tmp/" \
	-e "HOME=/tmp/" \
	-u "1000:1000" \
	delve-image:latest \
	bash -c "go mod download && /go/bin/dlv debug main.go"

test: pushtest pulltest

pushtest:
	docker pull fedora:32
	docker tag fedora:32 127.0.0.1:5000/fedora:latest
	docker push 127.0.0.1:5000/fedora:latest

pulltest:
	docker rmi 127.0.0.1:5000/fedora:latest
	docker pull 127.0.0.1:5000/fedora:latest
	docker run 127.0.0.1:5000/fedora:latest bash -c "echo This Image Runs"