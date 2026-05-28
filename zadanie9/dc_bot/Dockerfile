#FROM ubuntu:24.04
#
#ENV DEBIAN_FRONTEND=noninteractive
#
#SHELL ["/bin/bash", "-c"]
#
#WORKDIR /turbo_bot
#COPY ./ /turbo_bot
#
#RUN apt-get update && \
#    apt-get install -y software-properties-common curl zip unzip gnupg && \
#    apt-get install -y openjdk-21-jdk
#
#RUN curl -s https://get.sdkman.io | bash && \
#    source "$HOME/.sdkman/bin/sdkman-init.sh" && \
#    sdk install kotlin && \
#    sdk install gradle 8.11
#
#ENV SDKMAN_DIR=/root/.sdkman
#ENV PATH="${SDKMAN_DIR}/candidates/gradle/current/bin:${SDKMAN_DIR}/candidates/kotlin/current/bin:${PATH}"
#
#EXPOSE 8000
#
#CMD ["gradle", "run"]


# =========================================================
# ETAP 1: Budowanie aplikacji (Gruby obraz z Gradle)
# =========================================================
FROM gradle:8.11-jdk23 AS builder

WORKDIR /home/gradle/project

# 1. Kopiujemy TYLKO pliki konfiguracyjne Gradle
# Dzięki temu Docker scachuje pobrane biblioteki i nie będzie ich pobierał przy zmianie kodu bota
COPY --chown=gradle:gradle build.gradle* settings.gradle* gradle.properties* ./

# Pobieramy zależności (sztuczka na wymuszenie cache'owania warstwy)
RUN gradle dependencies --no-daemon || true

# 2. Kopiujemy resztę kodu źródłowego (dopiero teraz!)
COPY --chown=gradle:gradle . .

# 3. Budujemy gotowy plik JAR (pomijamy testy dla szybkości)
# Jeśli używasz Ktora, najpewniej masz task 'buildFatJar' lub 'shadowJar'
RUN gradle buildFatJar --no-daemon || gradle shadowJar --no-daemon || gradle build -x test --no-daemon

# =========================================================
# ETAP 2: Uruchamianie aplikacji (Lekki obraz produkcyjny)
# =========================================================
FROM eclipse-temurin:23-jre-noble

WORKDIR /turbo_bot

# Kopiujemy z ETAPU 1 tylko zbudowany plik .jar (zmień nazwę/ścieżkę jeśli Twój projekt generuje inaczej)
# Zazwyczaj dla bota/Ktora będzie to w build/libs/*all.jar lub po prostu *.jar
COPY --from=builder /home/gradle/project/build/libs/*.jar ./bot.jar

ENV DEBIAN_FRONTEND=noninteractive
EXPOSE 8000

# Odpalamy czystym procesem Javy - startuje w milisekundy!
CMD ["java", "-jar", "bot.jar"]