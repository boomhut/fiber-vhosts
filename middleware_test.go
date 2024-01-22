package vhosts

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Test the VhostsHandler middleware
func TestVhostsHandler(t *testing.T) {
	app := fiber.New()

	// Create a new vhost
	vhost := Vhost{
		Hostname:     "test.com",
		Handler:      mockMiddleware,
		ErrorHandler: mockErrorHandler,
	}

	// Create a new vhosts
	vhosts = &Vhosts{}
	vhosts.Add(vhost)

	// Register the vhosts middleware
	app.Use(VhostsHandler)

	// Create a new request
	req := httptest.NewRequest("GET", "http://test.com/", nil)

	// Perform the request
	resp, err := app.Test(req)
	assert.Equal(t, nil, err, "app.Test(req)")

	assert.Equal(t, 200, resp.StatusCode, "Status code")
	// check the response body
	body, err := io.ReadAll(resp.Body)

	// print the response body
	t.Logf("Response Body: %s", string(body))

	assert.Equal(t, nil, err, "ioutil.ReadAll(resp.Body)")
	assert.Equal(t, "Hello, World!", string(body), "Response body")

	// Create a new request for a vhost that doesn't exist
	req = httptest.NewRequest("GET", "http://test2.com/", nil)

	// Perform the request
	resp, err = app.Test(req)
	assert.Equal(t, nil, err, "app.Test(req)")
	assert.Equal(t, 404, resp.StatusCode, "Status code")

}

func mockMiddleware(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func mockErrorHandler(c *fiber.Ctx, err error) error {
	return c.SendStatus(404)
}
