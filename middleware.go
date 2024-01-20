package vhosts

import "github.com/gofiber/fiber/v2"

// FiberHandler is the handler for the vhost middleware
type FiberHandler func(*fiber.Ctx) error

// FiberErrorHandler is the error handler for the vhost middleware
type FiberErrorHandler func(*fiber.Ctx, error) error

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

	// Set some values on the context
	c.Locals("vhost.hostname", hostname)               // Hostname
	c.Locals("vhost.websiteID", vhost.WebsiteID)       // Website ID
	c.Locals("vhost.errorHandler", vhost.ErrorHandler) // Error Handler

	// Call the vhost's middleware
	return vhost.Handler(c)
}
