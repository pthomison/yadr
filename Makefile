define golang_docker_exec 
docker run \
-it --rm \
-v "$(PWD):/hacking" \
-w "/hacking" \
-p "5000:5000" \
-e "GOCACHE=/tmp/" \
-e "CGO_ENABLED=0" \
-u "1000:1000" \
golang:latest
endef

run:
	clear
	$(golang_docker_exec) go run main.go --data-directory /hacking/data

build:
	clear
	$(golang_docker_exec) go build -o yadr main.go

image:
	clear
	docker build . -t yadr:latest

runimage: image
	docker run -it --rm -p "5000:5000" yadr

cleardata:
	rm -rf ./data

runclean: cleardata run

test: pushtest pulltest

pushtest:
	docker pull fedora:32
	docker tag fedora:32 127.0.0.1:5000/fedora:latest
	docker push 127.0.0.1:5000/fedora:latest

pulltest:
	docker rmi 127.0.0.1:5000/fedora:latest
	docker pull 127.0.0.1:5000/fedora:latest
	docker run 127.0.0.1:5000/fedora:latest bash -c "echo This Image Runs"

compliance:
	docker run \
	-it \
	--network host \
	-e "OCI_ROOT_URL=http://127.0.0.1:5000" \
	-e "OCI_NAMESPACE=fedora" \
	-e "OCI_TEST_PULL=1" \
	-e "OCI_TEST_PUSH=1" \
	--name conform \
	pthomison/oci-conformance-registry-tester:latest || true
	docker cp conform:/report.html ./conformance-report.html
	docker rm -f conform