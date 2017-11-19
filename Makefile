REGISTRY ?= bsdavidson
HOST ?=

dev:
	docker-compose up

build:
	docker-compose build
	docker build --tag $(REGISTRY)/trimetric .
	docker push $(REGISTRY)/trimetric

deploy:
ifndef HOST
	$(error HOST variable is required)
endif
	scp docker-compose.yml $(HOST):trimetric/docker-compose.yml
	ssh $(HOST) 'cd trimetric && docker-compose pull api && docker-compose up -d --no-deps api'