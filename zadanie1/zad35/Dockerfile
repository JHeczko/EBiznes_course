FROM ubuntu:24.04

SHELL ["/bin/bash", "-c"]

RUN apt-get update && \
    apt-get install -y software-properties-common curl zip unzip gnupg

# Python 3.10
RUN add-apt-repository ppa:deadsnakes/ppa && \
    apt-get update && \
    apt-get install -y python3.10

RUN update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.10 1

# Java 8 (Temurin / Adoptium)
RUN apt-get install -y openjdk-8-jdk

# Kotlin via SDKMAN
RUN curl -s https://get.sdkman.io | bash && \
    source "$HOME/.sdkman/bin/sdkman-init.sh" && \
    sdk install kotlin

CMD ["/bin/bash"]