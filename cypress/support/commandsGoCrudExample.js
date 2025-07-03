const DEFAULT_CREDS = { login: 'admin', password: 'admin' };

// Custom command to get URL path
Cypress.Commands.add('getPath', () => cy.location('pathname'));

Cypress.Commands.add('resetDatabase', () => {
  cy.request({
    url: `http://localhost:8080/tools/db-clear`,
  }).then((response) => {
    expect(response.status).to.eq(200)
    cy.request({
      url: `http://localhost:8080/tools/seed`,
    }).then((response) => {
      expect(response.status).to.eq(200)
    });
  });
});

Cypress.Commands.add('login', ({ login, password } = DEFAULT_CREDS) => {
  cy.session([login, password], () => {
    cy.visit(`http://localhost:8080/login`);
    cy.get('#login').type(login);
    cy.get('#password').type(password);
    cy.get('button[type="submit"]').click();
    cy.get('h1').should('contain', 'Users');
    cy.getPath().should('eq', '/admin/users');
  });
});
