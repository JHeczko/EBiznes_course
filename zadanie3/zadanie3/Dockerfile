FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

SHELL ["/bin/bash", "-c"]

WORKDIR /turbo_bot
COPY ./ /turbo_bot

RUN apt-get update && \
    apt-get install -y software-properties-common curl zip unzip gnupg && \
    apt-get install -y openjdk-21-jdk

RUN curl -s https://get.sdkman.io | bash && \
    source "$HOME/.sdkman/bin/sdkman-init.sh" && \
    sdk install kotlin && \
    sdk install gradle 8.11

ENV SDKMAN_DIR=/root/.sdkman
ENV PATH="${SDKMAN_DIR}/candidates/gradle/current/bin:${SDKMAN_DIR}/candidates/kotlin/current/bin:${PATH}"

EXPOSE 8000

CMD ["gradle", "run"]