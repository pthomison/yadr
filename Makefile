

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

export 
export OCI_NAMESPACE="myorg/myrepo"
export OCI_USERNAME="myuser"
export OCI_PASSWORD="mypass"