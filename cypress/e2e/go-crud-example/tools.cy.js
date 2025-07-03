describe('Tools', () => {
  it('/tools/db-clear', () => {
    cy.request({
      url: `http://localhost:8080/tools/db-clear`,
      followRedirect: false,       // do not follow so we can inspect the 3xx
      failOnStatusCode: false      // prevent Cypress from failing on 3xx
    }).then((response) => {
      expect(response.status).to.eq(303)

      // Location header (absolute or relative)
      expect(response.headers).to.have.property('location')
      expect(response.headers.location).to.equal('/tools')
    });
  })

  /**
   * Verifies that GET /tools/sql?q=select+1
   * returns the expected JSON structure.
   */
  it('select 1', () => {
    cy.request({
      method: 'GET',
      url: 'http://localhost:8080/tools/sql',
      qs: { q: 'select 1' },          // Cypress will URL-encode this to ?q=select+1
    }).then((response) => {
      expect(response.status).to.eq(200)

      expect(response.body).to.deep.equal({
        out: [{ '?column?': 1 }],
        q: 'select 1'
      })
    })
  })

  it('delete from page', () => {
    cy.request({
      method: 'GET',
      url: 'http://localhost:8080/tools/sql',
      qs: { q: 'delete from page' },
    }).then((response) => {
      expect(response.status).to.eq(200)

      expect(response.body).to.deep.equal({
        out: [],
        q: 'delete from page'
      })
    })
  })

  it('INSERT INTO public.page...', () => {
    cy.request({
      method: 'GET',
      url: 'http://localhost:8080/tools/sql',
      qs: {
        q:
          "INSERT INTO public.page (slug, content, created_at, updated_at) " +
          "VALUES ( 'contact', 'Get in touch with us at contact@example.com', " +
          "NOW(), NOW() );"
      },
    }).then((response) => {
      expect(response.status).to.eq(200)

      expect(response.body).to.deep.equal({
        out: [],
        q: "INSERT INTO public.page (slug, content, created_at, updated_at) " +
          "VALUES ( 'contact', 'Get in touch with us at contact@example.com', " +
          "NOW(), NOW() );"
      })
    })
  })

  it('select * from page', () => {
    cy.request({
      method: 'GET',
      url: 'http://localhost:8080/tools/sql',
      qs: {
        q:
          "select * from page"
      },
    }).then((response) => {
      expect(response.status).to.eq(200)

      const jsonData = response.body;

      // Check for required fields
      expect(jsonData).to.have.all.keys('out', 'q');

      // Verify query value
      expect(jsonData.q).to.eq('select * from page');

      // Verify 'out' is an array with one element
      expect(jsonData.out).to.be.an('array').with.lengthOf(1);

      // Verify fields in the first element
      const firstElement = jsonData.out[0];
      expect(firstElement).to.include.keys(
        'content',
        'created_at',
        'id',
        'slug',
        'updated_at'
      );
    })
  })

  it('handles errors', () => {
    cy.request({
      method: 'GET',
      url: 'http://localhost:8080/tools/sql',
      qs: { q: 'select 1 from 1' },
      failOnStatusCode: false
    }).then((response) => {
      expect(response.status).to.eq(500)

      expect(response.body).to.deep.equal({
        "error": 'Error executing SQL query: ERROR: syntax error at or near "1" (SQLSTATE 42601)'
      })
    })
  })
})
