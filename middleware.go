package vhosts

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// FiberHandler is the handler for the vhost middleware
type FiberHandler func(*fiber.Ctx) error

// FiberErrorHandler is the error handler for the vhost middleware
type FiberErrorHandler func(*fiber.Ctx, error) error

// VhostsHandler is the handler for the vhosts middleware

func XVhost(vh *Vhosts) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// debug("vhosts middleware")
		log.Debugf("vhosts middleware %s", c.Hostname())

		// log the number of vhosts
		log.Debugf("vhosts count %d", len(vh.Vhosts))

		// Get the hostname from the request
		hostname := c.Hostname()

		// Get the vhost with the given hostname
		fVhost, ok := vh.Get(hostname)
		if !ok {
			log.Debugf("vhost not found for hostname %s", hostname)
			// Return a 404 if the vhost doesn't exist
			return c.SendStatus(404)
		}

		log.Debugf("vhost found for hostname %s", hostname)
		log.Debugf("vhost websiteID %s", fVhost.WebsiteID)
		log.Debugf("vhost path %s", fVhost.Path)
		log.Debugf("vhost lastModified %d", fVhost.LastModified)
		log.Debugf("vhost errorHandler %v", fVhost.ErrorHandler)
		log.Debugf("vhost handler %v", fVhost.Handler)

		// Set some values on the context
		c.Locals("vhost.hostname", hostname)                // Hostname
		c.Locals("vhost.websiteID", fVhost.WebsiteID)       // Website ID
		c.Locals("vhost.errorHandler", fVhost.ErrorHandler) // Error Handler

		// Call the vhost's middleware
		return fVhost.Handler(c)
	}
}
