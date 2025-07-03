describe('Page Management', () => {
    before(() => {
        cy.resetDatabase();
    });

    beforeEach(() => {
        cy.login();
    });

    after(() => {
        cy.resetDatabase();
    });

    it('Displays Page management page correctly', () => {
        cy.visit(`http://localhost:8080/admin/pages`);
        cy.get('h1').should('contain', 'Pages');
        cy.get('th').eq(0).should('contain', 'ID');
        cy.get('th').eq(1).should('contain', 'Slug');
        cy.get('th').eq(2).should('contain', 'Actions');
        cy.get('table tr').should('have.length.gt', 1);
        cy.get('table tr').eq(1).find('a').should('contain', 'Edit');
        cy.get('table tr').eq(1).find('form').should('exist');
    });

    it('Enforces unique slugs', () => {
        const uniqueName = `testslug_${Date.now()}`;

        // Create first slug
        cy.visit(`http://localhost:8080/admin/pages/new`);
        cy.get('#slug').type(uniqueName);
        cy.get('#content').type('Sample content');
        cy.get('button[type="submit"]').click();
        cy.contains('Page was added.').should('be.visible');
        cy.contains('tbody', uniqueName).should('exist');

        // Attempt duplicate
        cy.visit(`http://localhost:8080/admin/pages/new`);
        cy.get('#slug').type(uniqueName);
        cy.get('#content').type('Sample content');
        cy.get('button[type="submit"]').click();
        cy.contains('ERROR: duplicate key value violates unique constraint').should('be.visible');
    });

    it('Creates and deletes pages', () => {
        const uniqueName = `testpage_${Date.now()}`;

        // Create page
        cy.visit(`http://localhost:8080/admin/pages/new`);
        cy.get('#slug').type(uniqueName);
        cy.get('#content').type('Sample content');
        cy.get('button[type="submit"]').click();
        cy.contains('Page was added.').should('be.visible');
        cy.contains('tbody', uniqueName).should('exist');

        // Delete page
        cy.get(`[data-selenium="delete-${uniqueName}"]`).click();
        cy.contains('Page was deleted').should('be.visible');
        cy.contains('tbody', uniqueName).should('not.exist');
    });

    it('Creates and edits pages', () => {
        const uniqueName = `testpage_${Date.now()}`;
        const updatedName = `${uniqueName}_edited`;

        // Create page
        cy.visit(`http://localhost:8080/admin/pages/new`);
        cy.get('#slug').type(uniqueName);
        cy.get('#content').type('Sample content');
        cy.get('button[type="submit"]').click();
        cy.contains('Page was added.').should('be.visible');

        // Edit page
        cy.get(`[data-selenium="edit-${uniqueName}"]`).click();
        cy.get('h1').should('contain', 'Edit Page');
        cy.get('#slug').clear().type(updatedName);
        cy.get('button[type="submit"]').click();
        cy.getPath().should('eq', '/admin/pages');
        cy.contains('tbody', updatedName).should('exist');
    });

    it('Creates and shows page', () => {
        const uniqueName = `testpage_${Date.now()}`;

        // Create page
        cy.visit(`http://localhost:8080/admin/pages/new`);
        cy.get('#slug').type(uniqueName);
        cy.get('#content').type('Sample content');
        cy.get('button[type="submit"]').click();
        cy.contains('Page was added.').should('be.visible');

        // View page
        cy.get(`[data-selenium="show-${uniqueName}"]`).click();
        cy.get('h1').should('contain', 'Showing Page ');

        cy.visit(`http://localhost:8080`);

        cy.get(`[data-selenium="page-list"]`).should('exist');

        cy.get('[data-selenium="page-list"] li').should('have.length', 4)

        cy.get('[data-selenium="page-list"]').within(() => {
            // this scope only looks at your <ul> children
            cy.get('li').eq(0).get('a').should('contain.text', 'about')
            cy.get('li').eq(1).get('a').should('contain.text', uniqueName)
        })

        cy.get('a').contains(uniqueName).click();
        cy.get('h1').should('contain', uniqueName);
        cy.get('body').should('contain', 'Sample content');
    });

    it('Shows error if page does not exist', () => {
        cy.visit(`http://localhost:8080/admin/pages/0`, { 'failOnStatusCode': false });
        cy.get('body').should('contain', 'Error');
        cy.get('body').should('contain', 'Page not found');
    });
});
