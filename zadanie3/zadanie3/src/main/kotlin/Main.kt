package org.zadanie

import dev.kord.common.entity.Snowflake
import dev.kord.core.Kord
import io.ktor.http.ContentType
import io.ktor.server.engine.embeddedServer
import io.ktor.server.netty.Netty
import io.ktor.server.response.respondText
import io.ktor.server.routing.get
import io.ktor.server.routing.routing
import kotlinx.coroutines.launch

//TIP To <b>Run</b> code, press <shortcut actionId="Run"/> or
// click the <icon src="AllIcons.Actions.Execute"/> icon in the gutter.
fun main() {
    embeddedServer(Netty, port = 8000) {
        val token: String = System.getenv("DISCORD_BOT_TOKEN") ?: ""

        launch {
            val kord = Kord(token)

            routing {
                get("/send") {
                    val messegeText: String = call.request.queryParameters["mess"] ?: "Nic nie dostałem ale i tak wysyłam"
                    val channelId: String = System.getenv("CHANNEL_ID") ?: ""

                    kord.getChannelOf<dev.kord.core.entity.channel.TextChannel>(Snowflake(channelId))
                        ?.createMessage(messegeText)

                    call.respondText("Send to discord $messegeText", contentType = ContentType.Text.Plain)
                }
            }
            kord.login()
        }
        println("Ktos server working on port http://localhost:8000")
    }.start(wait = true)
}