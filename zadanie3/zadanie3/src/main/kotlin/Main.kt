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

import org.zadanie.Product
import org.zadanie.Category


fun main() {
    val categories: List<Category> = listOf(
        Category("Laptopy"),    // ID: 1
        Category("Podzespoly"), // ID: 2
        Category("Akcesoria"),  // ID: 3
        Category("Monitory")    // ID: 4
    )

    val products: List<Product> = listOf(
        // Kategoria: Laptopy (cat_id = 1)
        Product(prod_name = "Asus Zenbook 12", price = 2700.0, cat_id = 1),
        Product(prod_name = "MacBook Air M2", price = 5200.0, cat_id = 1),
        Product(prod_name = "Lenovo Legion 5", price = 4500.0, cat_id = 1),

        // Kategoria: Podzespoly (cat_id = 2)
        Product(prod_name = "RTX 4070 Super", price = 2999.0, cat_id = 2),
        Product(prod_name = "Ryzen 7 7800X3D", price = 1600.0, cat_id = 2),

        // Kategoria: Akcesoria (cat_id = 3)
        Product(prod_name = "Logitech G Pro X", price = 450.0, cat_id = 3),
        Product(prod_name = "Klawiatura Keychron K2", price = 380.0, cat_id = 3),

        // Kategoria: Monitory (cat_id = 4)
        Product(prod_name = "Gigabyte M27Q", price = 1300.0, cat_id = 4),
        Product(prod_name = "Samsung Odyssey G5", price = 1100.0, cat_id = 4)
    )

    embeddedServer(Netty, port = 8000) {
        val token: String = System.getenv("DISCORD_BOT_TOKEN") ?: ""

        launch {
            val kord = Kord(token)

            // reading messeges event handler
            kord.on<dev.kord.core.event.message.MessageCreateEvent> {
                println("EVENT TRIGGERED. Message: ${message.content}")

                if (message.author?.isBot == true) return@on
                if (message.content.isBlank()) return@on

                val text = message.content.trim()

                // ping handler
                if (text.startsWith("!ping:")) {
                    message.channel.createMessage("pong: ${text.subSequence(6, text.length)}")
                    println("PingPong: $text")
                }

                // zad 4.0
                if (text.startsWith("!cat")) {
                    var messege_back = "-----All avaible categories-----\n"

                    for (cat in categories) {
                        messege_back += " - ${cat.cat_name}\n\t - Category ID: ${cat.cat_id}\n"
                    }

                    message.channel.createMessage(messege_back)
                    println("Listed categories")
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