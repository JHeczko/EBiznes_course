FROM sbtscala/scala-sbt:eclipse-temurin-alpine-21.0.8_9_1.12.7_3.8.2

WORKDIR "/biznes_app"

COPY ./src /biznes_app

RUN apk add --no-cache bash

RUN adduser -S scala-server && \
    addgroup -S scala-server && \
    chown -R scala-server /biznes_app && \
    chgrp -R scala-server /biznes_app

USER scala-server

EXPOSE 9000 

CMD ["sbt", "run"]