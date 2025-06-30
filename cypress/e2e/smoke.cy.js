describe("Server Smoke Test", () => {
    it("Loads homepage", () => {
        cy.visit("/");
        cy.get('h1').contains("Home").should("be.visible"); // Update with your content
    });
});
