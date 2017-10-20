FROM node:8

ARG GOOGLE_MAPS_API_KEY
ARG TRIMET_API_KEY

COPY package.json /opt/trimet/package.json
COPY yarn.lock /opt/trimet/yarn.lock
WORKDIR /opt/trimet/
RUN yarn install
COPY . /opt/trimet/

ENV GOOGLE_MAPS_API_KEY=$GOOGLE_MAPS_API_KEY
ENV TRIMET_API_KEY=$TRIMET_API_KEY
RUN yarn dist

EXPOSE 8080 9876
VOLUME ["/opt/trimet"]
CMD ["yarn", "start"]