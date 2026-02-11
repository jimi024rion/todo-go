FROM node:20.19.2-alpine3.20

RUN apk update \
	&& apk --no-cache add openjdk21 \
	&& rm -rf /var/cache/apk/*

RUN npm install -g firebase-tools
