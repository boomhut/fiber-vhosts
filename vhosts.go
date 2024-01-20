package vhosts

import (
	"encoding/gob"
	"errors"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Vhost struct {
	Hostname     string            // hostname is the hostname of the vhost
	Path         string            // path is the path of the vhost
	WebsiteID    string            // websiteID is the websiteID of the vhost
	ErrorHandler FiberErrorHandler // errorHandler is the error handler for the vhost
	Handler      FiberHandler      // middleware is the middleware for the vhost
}

// vhosts contains all the vhosts protected by mutex lock for concurrent access safety
type Vhosts struct {
	// vhosts is the list of vhosts
	Vhosts []Vhost
	// mutex is the mutex lock for concurrent access safety
	mutex sync.RWMutex
}

// NewVhost returns a new vhost with the given hostname, path, websiteID and middleware
func NewVhost(hostname, path, websiteID string, handler FiberHandler, errorHandler FiberErrorHandler) Vhost {
	return Vhost{
		Hostname:     hostname,
		Path:         path,
		WebsiteID:    websiteID,
		Handler:      handler,
		ErrorHandler: errorHandler,
	}
}

// add adds a vhost to the vhosts list
func (v *Vhosts) Add(vhost Vhost) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.Vhosts = append(v.Vhosts, vhost)
}

// get returns the vhost with the given hostname
func (v *Vhosts) Get(hostname string) (Vhost, bool) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	for _, vhost := range v.Vhosts {
		if vhost.Hostname == hostname {
			return vhost, true
		}
	}
	return Vhost{}, false
}

// remove removes the vhost with the given hostname
func (v *Vhosts) Remove(hostname string) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	for i, vhost := range v.Vhosts {
		if vhost.Hostname == hostname {
			v.Vhosts = append(v.Vhosts[:i], v.Vhosts[i+1:]...)
			break
		}
	}
}

// length returns the length of the vhosts list
func (v *Vhosts) length() int {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return len(v.Vhosts)
}

// getVhosts returns the vhosts list
func (v *Vhosts) getVhosts() []Vhost {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return v.Vhosts
}

// getVhostnames returns the vhostnames list ( []string )
func (v *Vhosts) GetVhostnames() []string {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	var vhostnames []string
	for _, vhost := range v.Vhosts {
		vhostnames = append(vhostnames, vhost.Hostname)
	}
	return vhostnames
}

// get middleware returns the middleware for the vhost with the given hostname
func (v *Vhosts) getHandler(hostname string) (func(*fiber.Ctx) error, bool) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	for _, vhost := range v.Vhosts {
		if vhost.Hostname == hostname {
			return vhost.Handler, true
		}
	}
	return nil, false
}

// Save saves the vhosts list to a file at the given path
func (v *Vhosts) Save(path string) error {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	return save(path, v.Vhosts)
}

// Load loads the vhosts list from a file at the given path
func (v *Vhosts) Load(path string) error {

	// Check if the file already exists, if it doesn't return an error
	if !doesFileExist(path) {
		return errors.New("file doesn't exist")
	}

	// Load the vhosts list from the file at the given path
	vhosts, err := load(path)
	if err != nil {
		return err
	}

	// Set the vhosts list
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.Vhosts = vhosts

	return nil
}

// save saves the vhosts list to a file at the given path
func save(path string, vhosts []Vhost) error {

	// save using gob encoding
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	err = gobEncode(file, vhosts)
	if err != nil {
		return err
	}

	return nil

}

// gobEncode(f *os.File, vhosts []Vhost) error {
func gobEncode(f *os.File, vhosts []Vhost) error {
	encoder := gob.NewEncoder(f)
	err := encoder.Encode(vhosts)
	if err != nil {
		return err
	}
	return nil
}

// load loads the vhosts list from a file at the given path
func load(path string) ([]Vhost, error) {

	// load using gob decoding

	// Open the file at the given path
	loadFile, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer loadFile.Close()

	// gob decode the vhosts list
	var vhosts []Vhost
	decoder := gob.NewDecoder(loadFile)
	err = decoder.Decode(&vhosts)
	if err != nil {
		return nil, err
	}

	return vhosts, nil
}

// vhosts is the vhosts list
var vhosts *Vhosts

// init initializes the vhosts list
func init() {
	vhosts = &Vhosts{}
}

// Init initializes the vhosts list from a file at the given path
func InitVHostDataFile(path string) error {
	return vhosts.Load(path)
}

// Initialize initializes the vhosts list with some vhosts defaults map of hostname to middleware ( map[string]func(*fiber.Ctx) error )
func Initialize(listOfHostnames map[string]map[string]interface{}) {

	// Add the vhosts to the vhosts list
	for hostname, middleware := range listOfHostnames {
		vhosts.Add(NewVhost(hostname, "", "", middleware["handler"].(func(*fiber.Ctx) error), middleware["errorHandler"].(func(*fiber.Ctx, *error) error)))
	}
}
