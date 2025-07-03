describe('User Management', () => {
  before(() => {
    cy.resetDatabase();
  });

  beforeEach(() => {
    cy.login();
  });

  after(() => {
    cy.resetDatabase();
  });

  it('Displays user management page correctly', () => {
    cy.visit(`http://localhost:8080/admin/users`);
    cy.get('h1').should('contain', 'Users');
    cy.get('th').eq(0).should('contain', 'ID');
    cy.get('th').eq(1).should('contain', 'Login');
    cy.get('th').eq(2).should('contain', 'Actions');
    cy.get('table tr').should('have.length.gt', 1);
    cy.get('table tr').eq(1).find('a').should('contain', 'Edit');
    cy.get('table tr').eq(1).find('form').should('exist');
  });

  it('Validates new user form requirements', () => {
    cy.visit(`http://localhost:8080/admin/users/new`);
    cy.get('#login').type('1');
    cy.get('#password').type('1');
    cy.get('button[type="submit"]').click();
    cy.contains('Validation error').should('be.visible');
    cy.contains('Field is too short').should('be.visible');
  });

  it('Enforces unique usernames', () => {
    const uniqueName = `testuser_${Date.now()}`;

    // Create first user
    cy.visit(`http://localhost:8080/admin/users/new`);
    cy.get('#login').type(uniqueName);
    cy.get('#password').type('password');
    cy.get('button[type="submit"]').click();
    cy.contains('User was added.').should('be.visible');
    cy.contains('tbody', uniqueName).should('exist');

    // Attempt duplicate
    cy.visit(`http://localhost:8080/admin/users/new`);
    cy.get('#login').type(uniqueName);
    cy.get('#password').type('password');
    cy.get('button[type="submit"]').click();
    cy.contains('ERROR: duplicate key value violates unique constraint "user_login_key" (SQLSTATE 23505)').should('be.visible');
  });

  it('Creates and deletes users', () => {
    const uniqueName = `testuser_${Date.now()}`;

    // Create user
    cy.visit(`http://localhost:8080/admin/users/new`);
    cy.get('#login').type(uniqueName);
    cy.get('#password').type('password');
    cy.get('button[type="submit"]').click();
    cy.contains('User was added.').should('be.visible');
    cy.contains('tbody', uniqueName).should('exist');

    // Delete user
    cy.get(`[data-selenium="delete-${uniqueName}"]`).click();
    cy.contains('User was deleted').should('be.visible');
    cy.contains('tbody', uniqueName).should('not.exist');
  });

  it('Creates and edits users', () => {
    const uniqueName = `testuser_${Date.now()}`;
    const updatedName = `${uniqueName}_edited`;

    // Create user
    cy.visit(`http://localhost:8080/admin/users/new`);
    cy.get('#login').type(uniqueName);
    cy.get('#password').type('password');
    cy.get('button[type="submit"]').click();
    cy.contains('User was added.').should('be.visible');

    // Edit user
    cy.get(`[data-selenium="edit-${uniqueName}"]`).click();
    cy.get('h1').should('contain', 'Edit User');
    cy.get('#login').clear().type(updatedName);
    cy.get('button[type="submit"]').click();
    cy.getPath().should('eq', '/admin/users');
    cy.contains('tbody', updatedName).should('exist');
  });

  it('Shows error if user does not exist', () => {
    cy.visit(`http://localhost:8080/admin/users/0`, { 'failOnStatusCode': false });
    cy.get('body').should('contain', 'Error');
    cy.get('body').should('contain', 'User not found');
});
});
