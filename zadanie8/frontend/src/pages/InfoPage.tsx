import "./infoPage.css";
import workshopImg from "../../public/woodworker.jpg";

function InfoPage() {
    return (
        <main className="info-page">
            <section className="info-hero">
                <div className="hero-content">
                    <h1>Our Story</h1>
                    <div className="underline"></div>
                </div>
            </section>

            <section className="info-grid">
                <div className="info-text-card">
                    <h2>Crafting Comfort Since 2010</h2>
                    <p>
                        It all started in a small workshop. Our passion for wood
                        and modern design led us to create furniture that doesn't
                        just fill a room – it defines it.
                    </p>
                    <p>
                        Every piece is a result of hours of hand-polishing and
                        a deep respect for sustainable forestry.
                    </p>
                </div>

                <div className="info-image-container">
                    <img
                        src={workshopImg}
                        alt="Our workshop"
                        className="info-image"
                    />
                </div>
            </section>

            {/* Kontakt - stylistyka Brand Values z Twojego kodu */}
            <section className="info-contact">
                <div className="contact-box">
                    <h3>Visit Our Studio</h3>
                    <p>Design Street 12, Warsaw</p>
                    <p>Mon - Fri: 9:00 - 18:00</p>
                </div>
                <div className="contact-box">
                    <h3>Get in Touch</h3>
                    <p>hello@furniturestore.com</p>
                    <p>+48 123 456 789</p>
                </div>
            </section>
        </main>
    );
}

export default InfoPage;