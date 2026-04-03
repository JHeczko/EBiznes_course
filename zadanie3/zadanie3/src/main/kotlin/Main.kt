package org.zadanie

import dev.kord.common.entity.Snowflake
import dev.kord.core.Kord
import dev.kord.core.on
import io.ktor.http.ContentType
import io.ktor.server.engine.embeddedServer
import io.ktor.server.netty.Netty
import io.ktor.server.response.respondText
import io.ktor.server.routing.get
import io.ktor.server.routing.routing
import kotlinx.coroutines.launch
import dev.kord.gateway.Intent
import dev.kord.gateway.PrivilegedIntent


fun main() {
    embeddedServer(Netty, port = 8000) {
        val token: String = System.getenv("DISCORD_BOT_TOKEN") ?: ""

        launch {
            val kord = Kord(token)

            kord.on<dev.kord.core.event.message.MessageCreateEvent> {
                println("EVENT TRIGGERED: ${message.content}")

                if (message.author?.isBot == true) return@on

                val text = message.content
                if (text.isNotBlank()) {
                    message.channel.createMessage("Odbijam: $text")
                    println("Dostałem i odesłałem: $text")
                }else{
                    println("Literalnie pusty string $text")
                }

            }

            routing {
                get("/send") {
                    val messegeText: String = call.request.queryParameters["mess"] ?: "Nic nie dostałem ale i tak wysyłam"
                    val channelId: String = System.getenv("CHANNEL_ID") ?: ""

                    kord.getChannelOf<dev.kord.core.entity.channel.TextChannel>(Snowflake(channelId))
                        ?.createMessage(messegeText)

                    call.respondText("Send to discord $messegeText", contentType = ContentType.Text.Plain)
                }
            }

            try {
                println("LOGGING IN...")
                kord.login {
                    @OptIn(PrivilegedIntent::class)
                    intents += Intent.MessageContent
                }
                println("LOGGED!")
            } catch (e: Exception) {
                println("ERROR: ${e.message}")
                e.printStackTrace()
            }
        }
        println("Ktor server working on port http://localhost:8000")
    }.start(wait = true)
}