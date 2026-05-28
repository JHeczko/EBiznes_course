package org.zadanie

import java.net.URI
import java.net.http.HttpClient
import java.net.http.HttpRequest
import java.net.http.HttpResponse
import io.ktor.server.request.receiveText
import dev.kord.common.entity.Snowflake
import dev.kord.core.Kord
import dev.kord.core.on
import io.ktor.http.ContentType
import io.ktor.server.engine.embeddedServer
import io.ktor.server.netty.Netty
import io.ktor.server.response.respondText
import io.ktor.server.routing.get
import io.ktor.server.routing.post
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
        Category("Monitory"),    // ID: 4
        Category("Jedzenie")    // ID: 5
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
        val pythonApiUrl: String = System.getenv("PYTHON_API_URL") ?: "http://backend-service:8001/ask"

        launch {
            val kord = Kord(token)

            kord.on<dev.kord.core.event.message.MessageCreateEvent> {
                println("EVENT TRIGGERED. Message: ${message.content}")

                if (message.author?.isBot == true) return@on
                if (message.content.isBlank()) return@on

                val text = message.content.trim()

                if (text.startsWith("!ask ")) {
                    val question = text.substringAfter("!ask ").trim()

                    if (question.isNotEmpty()) {
                        try {
                            // KLUCZOWA ZMIANA: Wymuszamy HTTP/1.1
                            val client = HttpClient.newBuilder()
                                .version(HttpClient.Version.HTTP_1_1)
                                .build()

                            // NAPRAWA: Zwykły escape zamiast potrójnych cudzysłowów
                            val escapedQuestion = question.replace("\"", "\\\"").replace("\n", "\\n")
                            val jsonBody = "{\"question\": \"$escapedQuestion\"}"

                            // NAPRAWA: Dodany nagłówek Accept
                            val request = HttpRequest.newBuilder()
                                .uri(URI.create(pythonApiUrl))
                                .header("Content-Type", "application/json")
                                .header("Accept", "application/json")
                                .POST(HttpRequest.BodyPublishers.ofString(jsonBody))
                                .build()

                            val response = client.send(request, HttpResponse.BodyHandlers.ofString())

                            val responseText = response.body()
                                .substringAfter("\"response\":\"")
                                .substringBeforeLast("\"}")
                                .replace("\\n", "\n")
                                .replace("\\\"", "\"")

                            message.channel.createMessage("🤖 AI: $responseText")

                        } catch (e: Exception) {
                            message.channel.createMessage("Błąd połączenia z backendem AI: ${e.message}")
                            println(e)
                        }
                    } else {
                        message.channel.createMessage("Użycie: !ask <twoje pytanie>")
                    }
                }

                if (text.startsWith("!ping:")) {
                    val text_message = text.substringAfter("!ping")
                    if (text_message.isEmpty()){
                        message.channel.createMessage("pong")
                    } else {
                        message.channel.createMessage("pong: ${text_message}")
                    }
                }

                if (text.startsWith("!cat")) {
                    var messege_back = "-----All avaible categories-----\n"
                    for (cat in categories) {
                        messege_back += " - ${cat.cat_name}\n\t - Category ID: ${cat.cat_id}\n"
                    }
                    message.channel.createMessage(messege_back)
                }

                if (text.startsWith("!prod")) {
                    val cat_id = text.substringAfter("!prod ").trim().toIntOrNull()
                    if (cat_id == null) {
                        message.channel.createMessage("Please enter a category number!")
                    } else {
                        val category_name = categories.find { x -> x.cat_id == cat_id }?.cat_name
                        val filtered_prods = products.filter { x -> x.cat_id == cat_id }

                        if(filtered_prods.isEmpty()) {
                            message.channel.createMessage("No avaible products for category ${category_name ?: cat_id}!")
                        } else {
                            var out_string = "Available products from category ${category_name}:\n"
                            for (prod in products) {
                                if (prod.cat_id == cat_id) {
                                    out_string += " - ${prod.prod_name}: ${prod.price} PLN\n"
                                }
                            }
                            message.channel.createMessage(out_string)
                        }
                    }
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

                post("/api/ask-bot") {
                    val userText = call.receiveText()

                    try {
                        // KLUCZOWA ZMIANA: Wymuszamy HTTP/1.1 również tutaj
                        val client = HttpClient.newBuilder()
                            .version(HttpClient.Version.HTTP_1_1)
                            .build()

                        // NAPRAWA: Zwykły escape zamiast potrójnych cudzysłowów
                        val escapedQuestion = userText.replace("\"", "\\\"").replace("\n", "\\n")
                        val jsonBody = "{\"question\": \"$escapedQuestion\"}"

                        // NAPRAWA: Dodany nagłówek Accept
                        val request = HttpRequest.newBuilder()
                            .uri(URI.create(pythonApiUrl))
                            .header("Content-Type", "application/json")
                            .header("Accept", "application/json")
                            .POST(HttpRequest.BodyPublishers.ofString(jsonBody))
                            .build()

                        val response = client.send(request, HttpResponse.BodyHandlers.ofString())

                        call.respondText(response.body(), contentType = ContentType.Application.Json)

                    } catch (e: Exception) {
                        call.respondText("""{"error": "${e.message}"}""", contentType = ContentType.Application.Json)
                    }
                }
            }

            try {
                kord.login {
                    @OptIn(PrivilegedIntent::class)
                    intents += Intent.MessageContent
                }
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }
        println("Ktor server working on port http://localhost:8000")
    }.start(wait = true)
}