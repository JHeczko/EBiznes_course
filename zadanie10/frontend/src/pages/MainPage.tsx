import "./mainPage.css";
import {Link} from "react-router-dom";

function MainPage() {
    return (
        <main className="main-page">
            {/* 1. HERO SECTION */}
            <section className="hero">
                <div className="hero-content">
                    <h1>Your space, your rules</h1>
                    <p>Crafting furniture for those who value minimalism and timeless quality.</p>
                    <Link className="cta-button" to="/products">Browse Catalog</Link>
                </div>
            </section>

            {/* 2. CATEGORIES */}
            <section className="categories">
                <div className="category-item">
                    <h2>Living Room</h2>
                    <p>Comfort at the highest level.</p>
                </div>
                <div className="category-item">
                    <h2>Bedroom</h2>
                    <p>The perfect place for your rest.</p>
                </div>
                <div className="category-item">
                    <h2>Home Office</h2>
                    <p>Design that inspires productivity.</p>
                </div>
            </section>

            {/* 3. BRAND VALUES */}
            <section className="brand-values">
                <div className="value-box">
                    <h3>Eco Materials</h3>
                    <p>We care for the planet as much as for your home.</p>
                </div>
                <div className="value-box">
                    <h3>Fast Delivery</h3>
                    <p>Your order will arrive safely and on time.</p>
                </div>
                <div className="value-box">
                    <h3>Handcrafted</h3>
                    <p>Every piece is carefully finished by skilled artisans.</p>
                </div>
            </section>
        </main>
    );
}

export default MainPage;