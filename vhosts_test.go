package vhosts

import (
	"os"
	"testing"
)

func TestNewVhost(t *testing.T) {
	vhost := NewVhost("localhost", "/", "1", mockMiddleware, mockErrorHandler)
	if vhost.Hostname != "localhost" {
		t.Errorf("Expected hostname to be 'localhost', got '%s'", vhost.Hostname)
	}
}

func TestVhosts_Add(t *testing.T) {
	vhosts := &Vhosts{}
	vhost := NewVhost("localhost", "/", "1", mockMiddleware, mockErrorHandler)
	err := vhosts.Add(vhost)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(vhosts.Vhosts) != 1 {
		t.Errorf("Expected 1 vhost, got %d", len(vhosts.Vhosts))
	}

	// add the same vhost again
	err = vhosts.Add(vhost)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if len(vhosts.Vhosts) != 1 {
		t.Errorf("Expected 1 vhost, got %d", len(vhosts.Vhosts))
	}

}

// Test GetVhosts() *Vhosts
func TestGetVhosts(t *testing.T) {
	vhosts := Vhs
	vhost := NewVhost("localhost", "/", "1", mockMiddleware, mockErrorHandler)
	vhosts.Add(vhost)
	gotVhosts := GetVhosts()
	if len(gotVhosts.Vhosts) != 2 {
		t.Errorf("Expected 1 vhost, got %d", len(gotVhosts.Vhosts))
	}
}

func TestVhosts_Get(t *testing.T) {
	vhosts := &Vhosts{}
	vhost := NewVhost("localhost", "/", "1", mockMiddleware, mockErrorHandler)
	vhosts.Add(vhost)
	gotVhost, ok := vhosts.Get("localhost")
	if !ok || gotVhost.Hostname != "localhost" {
		t.Errorf("Expected to get vhost 'localhost', got '%s'", gotVhost.Hostname)
	}

	// TestGetVhostnames
	vhostnames := GetVhostnames(vhosts)
	if len(vhostnames) != 1 {
		t.Errorf("Expected 1 vhost, got %d", len(vhostnames))
	}

	// add a second vhost and test again
	vhost2 := NewVhost("localhost2", "/", "1", mockMiddleware, mockErrorHandler)
	vhosts.Add(vhost2)
	vhostnames = GetVhostnames(vhosts)
	if len(vhostnames) != 2 {
		t.Errorf("Expected 2 vhosts, got %d", len(vhostnames))
	}

}

func TestVhosts_Remove(t *testing.T) {
	vhosts := &Vhosts{}
	vhost := NewVhost("localhost", "/", "1", mockMiddleware, mockErrorHandler)
	vhosts.Add(vhost)
	err := vhosts.Remove("localhost")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(vhosts.Vhosts) != 0 {
		t.Errorf("Expected 0 vhosts, got %d", len(vhosts.Vhosts))
	}

	// remove a vhost that doesn't exist
	err = vhosts.Remove("localhost")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	if len(vhosts.Vhosts) != 0 {
		t.Errorf("Expected 0 vhosts, got %d", len(vhosts.Vhosts))
	}

}

func TestVhosts_NumberOfVhosts(t *testing.T) {
	vhosts := &Vhosts{}
	vhost := NewVhost("localhost", "/", "1", mockMiddleware, mockErrorHandler)
	vhosts.Add(vhost)
	if vhosts.NumberOfVhosts() != 1 {
		t.Errorf("Expected 1 vhost, got %d", vhosts.NumberOfVhosts())
	}
}

func TestVhosts_getVhosts(t *testing.T) {
	vhosts := &Vhosts{}
	vhost := NewVhost("localhost", "/", "1", mockMiddleware, mockErrorHandler)
	vhosts.Add(vhost)
	gotVhosts := vhosts.getVhosts()
	if len(gotVhosts) != 1 {
		t.Errorf("Expected 1 vhost, got %d", len(gotVhosts))
	}
}

func TestVhosts_getHandler(t *testing.T) {
	vhosts := &Vhosts{}
	vhost := NewVhost("localhost", "/", "1", mockMiddleware, mockErrorHandler)
	vhosts.Add(vhost)
	handler, ok := vhosts.getHandler("localhost")
	if !ok || handler == nil {
		t.Errorf("Expected to get handler for 'localhost', got nil")
	}
}

func TestVhosts_Save(t *testing.T) {
	vhosts := &Vhosts{}
	vhost := NewVhost("localhost", "/", "1", mockMiddleware, mockErrorHandler)
	vhost2 := NewVhost("secondhost", "/", "1", mockMiddleware, mockErrorHandler)
	vhosts.Add(vhost)
	vhosts.Add(vhost2)
	err := vhosts.Save("test.bin")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestVhosts_Load(t *testing.T) {
	vhosts := &Vhosts{}
	err := vhosts.Load("test.bin")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test that the vhosts list was loaded correctly
	if len(vhosts.Vhosts) != 2 {
		t.Errorf("Expected 2 vhosts, got %d", len(vhosts.Vhosts))
	}

	// Get the vhostnames
	vhostnames := GetVhostnames(vhosts)
	if len(vhostnames) != 2 {
		t.Errorf("Expected 2 vhosts, got %d", len(vhostnames))
	}

	// Test that the vhosts list was loaded correctly
	if vhostnames[0] != "localhost" {
		t.Errorf("Expected hostname 'localhost', got '%s'", vhostnames[0])
	}
	if vhostnames[1] != "secondhost" {
		t.Errorf("Expected hostname 'secondhost', got '%s'", vhostnames[1])
	}
}

// Test load with file that doesn't exist
func TestVhosts_Load_FileDoesntExist(t *testing.T) {
	vhosts := &Vhosts{}
	err := vhosts.Load("test2.bin")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

// test load with altered checksum (should fail)
func TestVhosts_Load_AlteredChecksum(t *testing.T) {

	// modify the checksum in the file to make it invalid
	// this should cause an error when loading the file

	// load the file
	ftm, err := os.OpenFile("test.bin", os.O_RDWR, 0644)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	defer ftm.Close()

	// change the checksum 5th byte before the end of the file
	// (checksum is 32 bytes long)
	ftm.Seek(-5, 2)
	ftm.Write([]byte("66"))

	// close the file
	ftm.Close()

	vhosts := &Vhosts{}
	err = vhosts.Load("test.bin")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestInitVHostDataFile(t *testing.T) {

	// call TestVhosts_Save to create the test.bin file
	TestVhosts_Save(t)

	err := InitVHostDataFile("test.bin")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

// Test Initialize (func Initialize(listOfHostnames map[string]map[string]interface{}))
func TestInitialize(t *testing.T) {

	// reset
	vhostReset()

	// create a map of hostname to middleware
	listOfHostnames := make(map[string]map[string]interface{})
	listOfHostnames["localhost"] = make(map[string]interface{})
	listOfHostnames["localhost"]["handler"] = mockMiddleware
	listOfHostnames["localhost"]["errorHandler"] = mockErrorHandler

	// initialize the vhosts list
	Initialize(listOfHostnames)

	// test that the vhosts list was initialized correctly
	if len(Vhs.Vhosts) != 1 {
		t.Errorf("Expected 1 vhost, got %d", len(Vhs.Vhosts))
	}

	// test that the vhost was initialized correctly
	if Vhs.Vhosts[0].Hostname != "localhost" {
		t.Errorf("Expected hostname 'localhost', got '%s'", Vhs.Vhosts[0].Hostname)
	}

}

// Test SetHandler (func SetHandler(hostname string, handler FiberHandler))
func TestSetHandler(t *testing.T) {

	// reset
	vhostReset()

	// create a map of hostname to middleware
	listOfHostnames := make(map[string]map[string]interface{})
	listOfHostnames["localhost"] = make(map[string]interface{})
	listOfHostnames["localhost"]["handler"] = mockMiddleware
	listOfHostnames["localhost"]["errorHandler"] = mockErrorHandler

	// initialize the vhosts list
	Initialize(listOfHostnames)

	// test that the vhosts list was initialized correctly
	if len(Vhs.Vhosts) != 1 {
		t.Errorf("Expected 1 vhost, got %d", len(Vhs.Vhosts))
	}

	// test that the vhost was initialized correctly
	if Vhs.Vhosts[0].Hostname != "localhost" {
		t.Errorf("Expected hostname 'localhost', got '%s'", Vhs.Vhosts[0].Hostname)
	}

	// set the handler for localhost
	Vhs.SetHandler("localhost", mockMiddleware)

	// test that the handler was set correctly
	if Vhs.Vhosts[0].Handler == nil {
		t.Errorf("Expected handler to be set, got nil")
	}

}

// Test SetErrorHandler (func SetErrorHandler(hostname string, errorHandler FiberErrorHandler))
func TestSetErrorHandler(t *testing.T) {

	// reset
	vhostReset()

	// create a map of hostname to middleware
	listOfHostnames := make(map[string]map[string]interface{})
	listOfHostnames["localhost"] = make(map[string]interface{})
	listOfHostnames["localhost"]["handler"] = mockMiddleware
	listOfHostnames["localhost"]["errorHandler"] = mockErrorHandler

	// initialize the vhosts list
	Initialize(listOfHostnames)

	// test that the vhosts list was initialized correctly
	if len(Vhs.Vhosts) != 1 {
		t.Errorf("Expected 1 vhost, got %d", len(Vhs.Vhosts))
	}

	// test that the vhost was initialized correctly
	if Vhs.Vhosts[0].Hostname != "localhost" {
		t.Errorf("Expected hostname 'localhost', got '%s'", Vhs.Vhosts[0].Hostname)
	}

	// set the error handler for localhost
	Vhs.SetErrorHandler("localhost", mockErrorHandler)

	// test that the error handler was set correctly
	if Vhs.Vhosts[0].ErrorHandler == nil {
		t.Errorf("Expected error handler to be set, got nil")
	}

}

// Set handler for vhost that doesn't exist
func TestSetHandler_VhostDoesntExist(t *testing.T) {

	// reset
	vhostReset()

	// create a map of hostname to middleware
	listOfHostnames := make(map[string]map[string]interface{})
	listOfHostnames["localhost"] = make(map[string]interface{})
	listOfHostnames["localhost"]["handler"] = mockMiddleware
	listOfHostnames["localhost"]["errorHandler"] = mockErrorHandler

	// initialize the vhosts list
	Initialize(listOfHostnames)

	// test that the vhosts list was initialized correctly
	if len(Vhs.Vhosts) != 1 {
		t.Errorf("Expected 1 vhost, got %d", len(Vhs.Vhosts))
	}

	// test that the vhost was initialized correctly
	if Vhs.Vhosts[0].Hostname != "localhost" {
		t.Errorf("Expected hostname 'localhost', got '%s'", Vhs.Vhosts[0].Hostname)
	}

	// set the handler for localhost
	err := Vhs.SetHandler("localhost2", mockMiddleware)

	// test that the error was returned
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

}

// Set error handler for vhost that doesn't exist
func TestSetErrorHandler_VhostDoesntExist(t *testing.T) {

	// reset
	vhostReset()

	// create a map of hostname to middleware
	listOfHostnames := make(map[string]map[string]interface{})
	listOfHostnames["localhost"] = make(map[string]interface{})
	listOfHostnames["localhost"]["handler"] = mockMiddleware
	listOfHostnames["localhost"]["errorHandler"] = mockErrorHandler

	// initialize the vhosts list
	Initialize(listOfHostnames)

	// test that the vhosts list was initialized correctly
	if len(Vhs.Vhosts) != 1 {
		t.Errorf("Expected 1 vhost, got %d", len(Vhs.Vhosts))
	}

	// test that the vhost was initialized correctly
	if Vhs.Vhosts[0].Hostname != "localhost" {
		t.Errorf("Expected hostname 'localhost', got '%s'", Vhs.Vhosts[0].Hostname)
	}

	// set the error handler for localhost
	err := Vhs.SetErrorHandler("localhost2", mockErrorHandler)

	// test that the error was returned
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

}

// Clean up
func TestVhosts_CleanUp(t *testing.T) {
	os.Remove("test.bin")
}
