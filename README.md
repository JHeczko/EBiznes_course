# E-Biznes 2026

## Zadanie 1 Docker

- 3.0 obraz ubuntu z Pythonem w wersji 3.10
- 3.5 obraz ubuntu:24.04 z Javą w wersji 8 oraz Kotlinem
- 4.0 do powyższego należy dodać najnowszego Gradle’a oraz paczkę JDBC SQLite w ramach projektu na Gradle (build.gradle)
- 4.5 stworzyć przykład typu HelloWorld oraz uruchomienie aplikacji przez CMD oraz gradle
- 5.0 dodać konfigurację docker-compose

Git commits:
- 3.0: https://github.com/JHeczko/EBiznes_course/commit/ae874f40e1ea87b8f0cdbaf893c0f09dfba14762
- 3.5: https://github.com/JHeczko/EBiznes_course/commit/c8fe2332282fe32e3ab2f23eb449cbca6e9284c6
- 4.0: https://github.com/JHeczko/EBiznes_course/commit/391b788ce8acf3f14550799f1bb6a983bf594dc4
- 4.5: https://github.com/JHeczko/EBiznes_course/commit/849fbd2c168d2a923e67a4c02e4465f1963184c8
- 5.0: https://github.com/JHeczko/EBiznes_course/commit/265422a6cf531639dea18a61c690efb45a750596

Dockerhub:
- 3.0: https://hub.docker.com/repository/docker/jheczko/zad30/general
- 3.5: https://hub.docker.com/repository/docker/jheczko/zad35/general
- 4.0: https://hub.docker.com/repository/docker/jheczko/zad40/general
- 4.5: https://hub.docker.com/repository/docker/jheczko/zad45/general
- 5.0: bez dockerhuba

Kod: zadanie1

## Zadanie 2 Scala
Należy stworzyć aplikację na frameworku Play lub Scalatra.

- 3.0 Należy stworzyć kontroler do Produktów
- 3.5 Do kontrolera należy stworzyć endpointy zgodnie z CRUD - dane
pobierane z listy
- 4.0 Należy stworzyć kontrolery do Kategorii oraz Koszyka + endpointy
zgodnie z CRUD
- 4.5 Należy aplikację uruchomić na dockerze (stworzyć obraz) oraz dodać
skrypt uruchamiający aplikację via ngrok
- 5.0 Należy dodać konfigurację CORS dla dwóch hostów dla metod CRUD

Git commits:
- 3.0: https://github.com/JHeczko/EBiznes_course/commit/db44f3f59a6474b92adec3615148b2cf2bb1eaf2
- 3.5: https://github.com/JHeczko/EBiznes_course/commit/c8db801625b89cd3e93e73660158f46ce9a95167
- 4.0: https://github.com/JHeczko/EBiznes_course/commit/ed8509629e11e30ee3453181ca3ace63a8566581
- 4.5: https://github.com/JHeczko/EBiznes_course/commit/26ae1868774f4402ae62462791d6ff9923fe2ae2
- 5.0: https://github.com/JHeczko/EBiznes_course/commit/c90328bd50b36449980113a310ec7c528bd2aae2

Video:
- https://drive.google.com/file/d/1Skt6fDH0UW0RJKDeX7j2TSBFlBSrdWTM/view?usp=sharing

Kontrolery mogą bazować na listach zamiast baz danych. CRUD: show all,
show by id (get), update (put), delete (delete), add (post).


## Zadanie 3 - Kotlin

- 3.0 Należy stworzyć aplikację kliencką w Kotlinie we frameworku Ktor, która pozwala na przesyłanie wiadomości na platformę Discord
- 3.5 Aplikacja jest w stanie odbierać wiadomości użytkowników z platformy Discord skierowane do aplikacji (bota)
- 4.0 Zwróci listę kategorii na określone żądanie użytkownika
- 4.5 Zwróci listę produktów wg żądanej kategorii
- 5.0 Aplikacja obsłuży dodatkowo jedną z platform: Slack lub Messenger

**Aplikację należy uruchomić na dockerze.**

Github:
- 3.0: https://github.com/JHeczko/EBiznes_course/commit/ff5f69c251c5a9972bcebfc62a4dc86b67d2c3d1
- 3.5: https://github.com/JHeczko/EBiznes_course/commit/6b2b0e5d98fe706934a1416b97b3eda694022a6a
- 4.0: https://github.com/JHeczko/EBiznes_course/commit/40902e5dffdd4d6bf1f0c8883c688f639572018a
- 4.5: https://github.com/JHeczko/EBiznes_course/commit/626af46b5ed6c988f8a4da493d4655bf00bde885

Video:
- https://drive.google.com/file/d/1LPeTigU6hT8rBU51Q9hNF4SiIMxHN4VO/view?usp=sharing

## Zadanie 4 - GoLang
Należy stworzyć projekt w **echo w Go**. Należy wykorzystać `gorm` do stworzenia kilka modeli, gdzie pomiędzy dwoma musi być relacja. Należy zaimplementować proste endpointy do dodawania oraz wyświetlania danych za pomocą modeli. Jako bazę danych można wybrać dowolną, sugerowałbym jednak pozostać przy sqlite.

- 3.0 Należy stworzyć aplikację we frameworki echo w j. Go, która będzie
miała kontroler Produktów zgodny z CRUD
- 3.5 Należy stworzyć model Produktów wykorzystując gorm oraz
wykorzystać model do obsługi produktów (CRUD) w kontrolerze (zamiast
listy)
- 4.0 Należy dodać model Koszyka oraz dodać odpowiedni endpoint
- 4.5 Należy stworzyć model kategorii i dodać relację między kategorią,
a produktem
- 5.0 pogrupować zapytania w gorm’owe scope'y

Github:
- 3.0 & 3.5: https://github.com/JHeczko/EBiznes_course/commit/8a1ed419c192fe8c761927dcb1eaf6b669d65cb0
- 4.0 & 4.5: https://github.com/JHeczko/EBiznes_course/commit/02de5c4694462827ebb67b3aadd38721d30c8d6f
- 5.0: https://github.com/JHeczko/EBiznes_course/commit/0abe3d79202b435cc84e275a17c75a8a4aa02de3

Demo: https://drive.google.com/file/d/1in7W6-iRclpbwEH2frtpNNUIAoRpqPYa/view?usp=sharing