REGISTRY ?= bsdavidson
HOST ?=

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

dev:
	docker-compose up

lint:
	gometalinter --cyclo-over=15 --fast --vendor
	yarn run lint

node_modules: package.json
	yarn install

node_modules/react-map-gl/dist/index.js: node_modules
	cd node_modules/react-map-gl && \
		yarn install && \
		yarn run build

test: test-go test-web

test-go:
	go test ./...

test-web: node_modules node_modules/react-map-gl/dist/index.js
	yarn test


.PHONY: dev build deploy test-web test test-go lint
