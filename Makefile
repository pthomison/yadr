
run:
	sudo docker run -it --rm -v "$(PWD):/hacking" -w "/hacking" -p "5000:5000" golang:latest go run main.go


bash:
	sudo docker run -it --rm -v "$(PWD):/hacking" -w "/hacking" golang:latest /bin/bash