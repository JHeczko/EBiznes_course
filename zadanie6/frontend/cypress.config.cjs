const { defineConfig } = require("cypress");

module.exports = defineConfig({
    e2e: {
        // frontend URL
        // baseUrl: "https://resemblant-evalyn-nonphonemically.ngrok-free.dev",
        baseUrl: "http://localhost:5173",
        specPattern: "cypress/e2e/**/*.cy.{js,ts}",
        supportFile: false,
        video: false,
        screenshotOnRunFailure: true,
        defaultCommandTimeout: 8000,
    },
});