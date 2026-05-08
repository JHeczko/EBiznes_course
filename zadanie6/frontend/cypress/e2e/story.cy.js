// ============================================================
// Task 6 – E2E Tests (Cypress)
// Furniture shop – 20 test cases, 50+ assertions
// Run: npx cypress run (or npx cypress open)
// ============================================================

const BASE = "http://localhost:5173";
//const BASE = "https://resemblant-evalyn-nonphonemically.ngrok-free.dev"

// ─────────────────────────────────────────────
// TC-01 Main page load good
// ─────────────────────────────────────────────
describe("TC-01: Main page", () => {
    it("show header, sections and CTA button", () => {
        cy.visit(BASE);

        // A1
        cy.get("header").should("be.visible");
        // A2
        cy.get("header").contains("Store").should("be.visible");
        // A3
        cy.get(".hero h1").should("contain", "Your space");
        // A4
        cy.get(".cta-button").should("have.attr", "href", "/products");
        // A5
        cy.get(".category-item").should("have.length", 3);
        // A6
        cy.get(".brand-values .value-box").should("have.length", 3);
    });
});

// ─────────────────────────────────────────────
// TC-02 Navigation – links in head
// ─────────────────────────────────────────────
describe("TC-02: Navigation via header", () => {
    it("go to /products and come back to /", () => {
        cy.visit(BASE);

        // A7
        cy.get("a.menu-element").contains("Products").click();
        // A8
        cy.url().should("include", "/products");

        cy.get("#store-name").click();
        // A9
        cy.url().should("eq", `${BASE}/`);
    });
});

// ─────────────────────────────────────────────
// TC-03 Navigation – Cart link
// ─────────────────────────────────────────────
describe("TC-03: Navigation – cart", () => {
    it("go to /cart from top menu", () => {
        cy.visit(BASE);

        cy.get("a.menu-element").contains("Cart").click();
        // A10
        cy.url().should("include", "/cart");
        // A11
        cy.get("main.cart-page h1").should("contain", "Cart");
    });
});

// ─────────────────────────────────────────────
// TC-04 Navigation – Info link
// ─────────────────────────────────────────────
describe("TC-04: Navigation – info", () => {
    it("go to /info page", () => {
        cy.visit(BASE);

        cy.get("a.menu-element").contains("Info").click();
        // A12
        cy.url().should("include", "/info");
        // A13
        cy.get("main.info-page").should("be.visible");
    });
});

// ─────────────────────────────────────────────
// TC-05 Products page – list loading fine
// ─────────────────────────────────────────────
describe("TC-05: Product list", () => {
    it("show minimum one product with price and category", () => {
        cy.visit(`${BASE}/products`);

        // A14
        cy.get(".product-card", { timeout: 8000 }).should("have.length.gte", 1);
        // A15
        cy.get(".product-card").first().find("h3").should("not.be.empty");
        // A16
        cy.get(".product-card").first().find(".price-tag").should("be.visible");
        // A17
        cy.get(".product-card").first().find(".category").should("contain", "Category:");
    });
});

// ─────────────────────────────────────────────
// TC-06 Add product to cart (toast success check)
// ─────────────────────────────────────────────
describe("TC-06: Add product to cart", () => {
    it("click 'Add' and toast pop up", () => {
        cy.visit(`${BASE}/products`);

        cy.get(".product-card", { timeout: 8000 }).first().find(".add-to-cart-btn").click();
        // A18
        cy.get(".go3958317564", { timeout: 4000 }).should("exist");
        // A19
        cy.contains("Added", { timeout: 4000 }).should("be.visible");
    });
});

// ─────────────────────────────────────────────
// TC-07 Cart – show products after add
// ─────────────────────────────────────────────
describe("TC-07: Cart – view after adding", () => {
    it("cart must have minimum one item", () => {
        cy.visit(`${BASE}/products`);
        cy.get(".product-card", { timeout: 8000 }).first().find(".add-to-cart-btn").click();
        cy.wait(500);

        cy.visit(`${BASE}/cart`);
        // A20
        cy.get(".cart-item", { timeout: 6000 }).should("have.length.gte", 1);
        // A21
        cy.get(".cart-item").first().find("h3").should("not.be.empty");
    });
});

// ─────────────────────────────────────────────
// TC-08 Cart – change quantity (+)
// ─────────────────────────────────────────────
describe("TC-08: Change quantity up", () => {
    it("click + button and number go big", () => {
        cy.visit(`${BASE}/cart`);

        cy.get(".cart-item", { timeout: 6000 }).first().then(($item) => {
            const before = parseInt($item.find("h3").last().text(), 10);

            cy.wrap($item).find(".quantity-control-button").last().click();

            cy.get(".cart-item").first().find("h3").last().invoke("text").then((text) => {
                const after = parseInt(text, 10);
                // A22
                expect(after).to.eq(before + 1);
            });
        });
    });
});

// ─────────────────────────────────────────────
// TC-09 Cart – change quantity (-)
// ─────────────────────────────────────────────
describe("TC-09: Change quantity down", () => {
    it("click - and number go small (but not less than 1)", () => {
        cy.visit(`${BASE}/cart`);

        cy.get(".cart-item", { timeout: 6000 }).first().then(($item) => {
            const cur = parseInt($item.find("h3").last().text(), 10);
            if (cur === 1) {
                cy.wrap($item).find(".quantity-control-button").last().click();
                cy.wait(300);
            }

            cy.get(".cart-item").first().then(($item2) => {
                const before2 = parseInt($item2.find("h3").last().text(), 10);
                cy.wrap($item2).find(".quantity-control-button").first().click();

                cy.get(".cart-item").first().find("h3").last().invoke("text").then((text) => {
                    const after = parseInt(text, 10);
                    // A23
                    expect(after).to.eq(before2 - 1);
                    // A24
                    expect(after).to.be.gte(1);
                });
            });
        });
    });
});

// ─────────────────────────────────────────────
// TC-10 Cart – delete item
// ─────────────────────────────────────────────
describe("TC-10: Remove item from cart", () => {
    it("click trash icon and item gone with toast", () => {
        cy.visit(`${BASE}/cart`);

        cy.get(".cart-item", { timeout: 6000 }).then(($items) => {
            const countBefore = $items.length;

            cy.get(".cart-item").first().find(".delete-button").click();

            // A25
            cy.get(".cart-item").should("have.length", countBefore - 1);
            // A26
            cy.contains("Deleted", { timeout: 3000 }).should("be.visible");
        });
    });
});

// ─────────────────────────────────────────────
// TC-11 Cart empty – "Empty cart" message
// ─────────────────────────────────────────────
describe("TC-11: Empty cart state", () => {
    it("show message when nothing in cart", () => {
        cy.request("DELETE", "http://localhost:13000/cart/1").then(() => {
            cy.visit(`${BASE}/cart`);
            // A27
            cy.get(".empty-cart-message").should("be.visible");
            // A28
            cy.get(".checkout-button").should("not.exist");
        });
    });
});

// ─────────────────────────────────────────────
// TC-12 Cart – Checkout button visible when > 0 items
// ─────────────────────────────────────────────
describe("TC-12: Checkout button visibility", () => {
    it("checkout button show when cart not empty", () => {
        cy.visit(`${BASE}/products`);
        cy.get(".product-card", { timeout: 8000 }).first().find(".add-to-cart-btn").click();
        cy.wait(600);

        cy.visit(`${BASE}/cart`);
        // A29
        cy.get(".checkout-button", { timeout: 5000 }).should("be.visible");
        // A30
        cy.get(".checkout-button").should("contain.text", "Checkout");
    });
});

// ─────────────────────────────────────────────
// TC-13 Info page – content sections
// ─────────────────────────────────────────────
describe("TC-13: Info page – content", () => {
    it("have hero, grid and contact info", () => {
        cy.visit(`${BASE}/info`);

        // A31
        cy.get(".info-hero h1").should("contain", "Our Story");
        // A32
        cy.get(".info-grid").should("be.visible");
        // A33
        cy.get(".info-image").should("be.visible");
        // A34
        cy.get(".info-contact .contact-box").should("have.length", 2);
        // A35
        cy.get(".info-contact").should("contain", "Warsaw");
    });
});

// ─────────────────────────────────────────────
// TC-14 Unknown route – fallback (404)
// ─────────────────────────────────────────────
describe("TC-14: Wrong URL", () => {
    it("app no crash when visit /unknown", () => {
        cy.visit(`${BASE}/unknown`, { failOnStatusCode: false });
        // A36
        cy.get("header").should("exist");
        // A37
        cy.get("body").should("be.visible");
    });
});

// ─────────────────────────────────────────────
// TC-15 Header – cart counter
// ─────────────────────────────────────────────
describe("TC-15: Cart link in head", () => {
    it("cart link in header is there", () => {
        cy.visit(BASE);

        // A38
        cy.get("header a.menu-element").contains("Cart").should("be.visible");
        // A39
        cy.get("header a.menu-element").contains("Cart").should("have.attr", "href", "/cart");
    });
});

// ─────────────────────────────────────────────
// TC-16 Products – every card has "Add" button
// ─────────────────────────────────────────────
describe("TC-16: 'Add' buttons on list", () => {
    it("every product has button for add to cart", () => {
        cy.visit(`${BASE}/products`);

        cy.get(".product-card", { timeout: 8000 }).each(($card) => {
            // A40
            cy.wrap($card).find(".add-to-cart-btn").should("be.visible");
        });
    });
});

// ─────────────────────────────────────────────
// TC-17 Cart – item price calculate good
// ─────────────────────────────────────────────
describe("TC-17: Price math in cart", () => {
    it("show item price = quantity x unit price", () => {
        cy.visit(`${BASE}/cart`);

        cy.get(".cart-item", { timeout: 6000 }).first().then(($item) => {
            const qty    = parseInt($item.find("h3").last().text(), 10);
            const priceText = $item.find("h5").text().replace(" zł", "").replace(",", ".");
            const price  = parseFloat(priceText);

            cy.request("GET", "http://localhost:13000/cart/1").then((res) => {
                const firstItem = res.body[0];
                const unitPrice = firstItem.product.price;
                const expected  = parseFloat((qty * unitPrice).toFixed(2));
                // A41
                expect(price).to.eq(expected);
            });
        });
    });
});

// ─────────────────────────────────────────────
// TC-18 API – GET /products return 200
// ─────────────────────────────────────────────
describe("TC-18: Verify product API", () => {
    it("GET /products return 200 and array", () => {
        cy.request("GET", "http://localhost:13000/products").then((res) => {
            // A42
            expect(res.status).to.eq(200);
            // A43
            expect(res.body).to.be.an("array");
            // A44
            expect(res.body.length).to.be.gte(1);
            // A45
            expect(res.body[0]).to.have.property("name");
            // A46
            expect(res.body[0]).to.have.property("price");
        });
    });
});

// ─────────────────────────────────────────────
// TC-19 API – GET /category return 200
// ─────────────────────────────────────────────
describe("TC-19: Verify category API", () => {
    it("GET /category return 200 and not empty array", () => {
        cy.request("GET", "http://localhost:13000/category").then((res) => {
            // A47
            expect(res.status).to.eq(200);
            // A48
            expect(res.body).to.be.an("array");
            // A49
            expect(res.body[0]).to.have.property("id");
            // A50
            expect(res.body[0]).to.have.property("name");
        });
    });
});

// ─────────────────────────────────────────────
// TC-20 Responsiveness – product page on mobile
// ─────────────────────────────────────────────
describe("TC-20: Responsiveness – mobile view", () => {
    it("product page look okay on iPhone 12", () => {
        cy.viewport("iphone-x");
        cy.visit(`${BASE}/products`);

        // A51
        cy.get(".products-page").should("be.visible");
        // A52
        cy.get(".product-card", { timeout: 8000 }).should("have.length.gte", 1);
        // A53
        cy.get("header").should("be.visible");
    });
});