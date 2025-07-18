describe('Login Functionality', () => {
    before(() => {
        cy.resetDatabase();
    });

    after(() => {
        cy.resetDatabase();
    });

    it('Rejects invalid credentials', () => {
        cy.visit(`http://localhost:8080/login`);
        cy.get('#login').type('admin');
        cy.get('#password').type('wrong_password');
        cy.get('button[type="submit"]').click();
        cy.getPath().should('eq', '/login');
        cy.contains('body', 'Invalid username or password').should('be.visible');
    });

    it('Logs in with valid credentials', () => {
        cy.login();
    });
});
