package-ui:
	cd ui && \
	npm run build && \
	cp -r dist ../server/public