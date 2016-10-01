.PHONY: build
build:
	docker build -t ssh-hack .

.PHONY: run
run:
	docker run -v /var/run/docker.sock:/var/run/docker.sock --rm -it --name ssh ssh-hack
