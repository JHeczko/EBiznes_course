FROM ubuntu:24.04

SHELL ["/bin/bash", "-c"]

# gradle project import
WORKDIR /app
COPY ./projekt .

RUN apt-get update && \
    apt-get install -y software-properties-common curl zip unzip gnupg

# Python 3.10
RUN add-apt-repository ppa:deadsnakes/ppa && \
    apt-get update && \
    apt-get install -y python3.10

RUN update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.10 1

# Java 8 (Temurin / Adoptium)
RUN apt-get install -y openjdk-8-jdk

# Kotlin and gradle via SDKMAN
RUN curl -s https://get.sdkman.io | bash && \
    source "$HOME/.sdkman/bin/sdkman-init.sh" && \
    sdk install kotlin && \
    sdk install gradle 6.9.4 && \
    gradle build --no-daemon

# add gradle to ENV in order to CMD it
ENV PATH="/root/.sdkman/candidates/gradle/current/bin:$PATH"

CMD ["gradle", "run", "--no-daemon"]