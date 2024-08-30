package-ui:
	cd ui && \
	npm run build && \
	cp -r dist ../server/public

docker:
	docker build . -t ideashare:dev