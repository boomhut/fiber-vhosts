package vhosts

import "github.com/gofiber/fiber/v2"

// VhostsHandler is the handler for the vhosts middleware
func VhostsHandler(c *fiber.Ctx) error {
	// Get the vhosts from the context

	// Get the hostname from the request
	hostname := c.Hostname()

	// Get the vhost with the given hostname
	vhost, ok := vhosts.Get(hostname)
	if !ok {
		// Return a 404 if the vhost doesn't exist
		return c.SendStatus(404)
	}

	// Call the vhost's middleware
	return vhost.Middleware(c)
}
