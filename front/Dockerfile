FROM node:22.14
WORKDIR /usr/app
COPY ./app ./
RUN yarn install
RUN yarn build
CMD ["yarn", "start"]

EXPOSE 3000