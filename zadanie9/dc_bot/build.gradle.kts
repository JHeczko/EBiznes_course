plugins {
    kotlin("jvm") version "2.3.10"
    id("io.ktor.plugin") version "3.4.1"
}

group = "org.zadanie"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
    maven("https://snapshots.kord.dev")
}

application {
    mainClass.set("org.zadanie.MainKt")
}

dependencies {
    testImplementation(kotlin("test"))
    implementation("io.ktor:ktor-server-core:3.4.1")
    implementation("io.ktor:ktor-server-netty:3.4.1")
    implementation("dev.kord:kord-core:0.18.1")
    implementation("io.ktor:ktor-server-cors:3.4.1")
}

kotlin {
    jvmToolchain(23)
}

tasks.test {
    useJUnitPlatform()
}