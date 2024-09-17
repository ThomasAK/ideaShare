package-ui:
	cd ui && \
	npm run build && \
	cp -r dist ../server/public

docker:
	docker build . -t ideashare:dev

local:
	docker compose up -d

reset:
	docker compose down -v

lint:
	(cd server && go fmt ./...) && \
	(cd ui && npm run lint:fix)