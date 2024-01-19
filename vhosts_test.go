package vhosts

import (
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestAdd(t *testing.T) {
	v := &Vhosts{}
	Vhost := Vhost{Hostname: "test.com"}
	v.Add(Vhost)

	if len(v.Vhosts) != 1 {
		t.Errorf("Expected length of Vhosts to be 1, got %d", len(v.Vhosts))
	}
}

func TestGet(t *testing.T) {
	v := &Vhosts{}
	Vhost := Vhost{Hostname: "test.com"}
	v.Add(Vhost)

	result, ok := v.Get("test.com")
	if !ok || result.Hostname != "test.com" {
		t.Errorf("Expected to get Vhost with hostname 'test.com', got %v", result)
	}

	// test getting a vhost that doesn't exist
	_, ok = v.Get("test2.com")
	if ok {
		t.Errorf("Expected to not get Vhost with hostname 'test2.com'")
	}

}

func TestRemove(t *testing.T) {
	v := &Vhosts{}
	Vhost := Vhost{Hostname: "test.com"}
	v.Add(Vhost)

	v.Remove("test.com")
	if len(v.Vhosts) != 0 {
		t.Errorf("Expected length of Vhosts to be 0, got %d", len(v.Vhosts))
	}

	// test removing a vhost that doesn't exist
	v.Remove("test2.com")
	if len(v.Vhosts) != 0 {
		t.Errorf("Expected length of Vhosts to be 0, got %d", len(v.Vhosts))
	}
}

func TestLength(t *testing.T) {
	v := &Vhosts{}
	Vhost := Vhost{Hostname: "test.com"}
	v.Add(Vhost)

	if v.length() != 1 {
		t.Errorf("Expected length of Vhosts to be 1, got %d", v.length())
	}
}

func TestGetVhosts(t *testing.T) {
	v := &Vhosts{}
	Vhost := Vhost{Hostname: "test.com"}
	v.Add(Vhost)

	if len(v.getVhosts()) != 1 {
		t.Errorf("Expected length of Vhosts to be 1, got %d", len(v.getVhosts()))
	}

}

func TestGetMiddleware(t *testing.T) {
	v := &Vhosts{}
	Vhost := Vhost{Hostname: "test.com", Middleware: func(c *fiber.Ctx) error { return nil }}
	v.Add(Vhost)

	_, ok := v.getMiddleware("test.com")
	if !ok {
		t.Errorf("Expected to get middleware for hostname 'test.com'")
	}

	// test getting middleware for a vhost that doesn't exist
	_, ok = v.getMiddleware("test2.com")
	if ok {
		t.Errorf("Expected to not get middleware for hostname 'test2.com'")
	}

}

func TestSaveAndLoad(t *testing.T) {
	v := &Vhosts{}
	Vhost := Vhost{Hostname: "test.com"}
	v.Add(Vhost)

	err := v.Save("test.gob")
	if err != nil {
		t.Errorf("Failed to save Vhosts: %v", err)
	}

	err = v.Load("test.gob")
	if err != nil {
		t.Errorf("Failed to load Vhosts: %v", err)
	}

	if v.length() != 1 {
		t.Errorf("Expected length of Vhosts to be 1, got %d", v.length())
	}

	// add another Vhost to the list and save it again to test overwriting the file on disk
	Vhost = NewVhost("test2.com", "", "", nil)
	v.Add(Vhost)

	err = v.Save("test.gob")
	if err != nil {
		t.Errorf("Failed to save Vhosts: %v", err)
	}

	// test getting the vhost from the list
	_, ok := v.Get("test2.com")
	if !ok {
		t.Errorf("Expected to get Vhost with hostname 'test2.com'")
	}

	// Cleanup
	os.Remove("test.gob")
}

// load from a file that doesn't exist
func TestLoadFileDoesNotExist(t *testing.T) {
	v := &Vhosts{}
	err := v.Load("testx.gob")
	if err == nil {
		t.Errorf("Expected to get error when loading from file that doesn't exist")
	}

}

// Mock middleware function for testing
func mockMiddleware(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func TestInitVHostDataFile(t *testing.T) {

	// create a file to test
	file, err := os.Create("testfile.bin")
	if err != nil {
		t.Errorf("Failed to create test file: %s", err)
	}
	file.Close()

	// put some data in the file to test
	file, err = os.OpenFile("testfile.bin", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Errorf("Failed to open test file: %s", err)
	}

	vhtest := []Vhost{
		{
			Hostname: "test.com",
			Middleware: func(c *fiber.Ctx) error {
				return c.Status(200).SendString("Hello, World!")
			},
		},
		{
			Hostname: "test2.com",
			Middleware: func(c *fiber.Ctx) error {
				return c.Status(200).SendString("Hello, World!")
			},
		},
	}

	err = gobEncode(file, vhtest)
	if err != nil {
		t.Errorf("Failed to encode test data: %s", err)
	}

	// close the file
	file.Close()

	// use the gob encoder to encode some data to the file

	err = InitVHostDataFile("testfile.bin")
	if err != nil {
		t.Errorf("InitVHostDataFile failed: %s", err)
	}
}

func TestInitialize(t *testing.T) {
	listOfHostnames := map[string]func(*fiber.Ctx) error{
		"localhost": mockMiddleware,
	}

	Initialize(listOfHostnames)

	// Add your assertions here to verify the vhosts have been initialized correctly
	// This will depend on how you can access and check the vhosts data
}

// Final cleanup after all tests have run
func TestFinalCleanup(t *testing.T) {
	os.Remove("test.gob")
	os.Remove("test.txt")
	os.Remove("testfile.bin")
}
