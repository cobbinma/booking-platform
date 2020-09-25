# syntax = docker/dockerfile:experimental
FROM node:14.4.0-alpine3.11 as build

ENV CI=true

RUN apk add --update yarn;

WORKDIR /root

COPY . .

RUN yarn config set no-progress && yarn

ENV NODE_PRODUCTION=true

FROM node:14.4.0-alpine3.11

ENV CI=true
ENV NODE_PRODUCTION=true

RUN apk add --update \
    bash \
    curl \
    tini \
    yarn \
    ;

WORKDIR /root

COPY package.json yarn.lock ./

COPY lib/booking-ui lib/booking-ui
COPY --from=build /root/lib/booking-ui/dist ./lib/booking-ui/dist
COPY --from=build /root/lib/booking-ui/build ./lib/booking-ui/build

RUN yarn config set no-progress && yarn install --production --frozen-lockfile

WORKDIR lib/booking-ui

ENTRYPOINT ["/sbin/tini", "--"]

CMD ["yarn", "start"]

EXPOSE 3000
