package vhosts

import (
	"encoding/gob"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Vhost struct {
	Hostname   string          // hostname is the hostname of the vhost
	Path       string          // path is the path of the vhost
	WebsiteID  string          // websiteID is the websiteID of the vhost
	Middleware FiberMiddleware // middleware is the middleware for the vhost
}

// vhosts contains all the vhosts protected by mutex lock for concurrent access safety
type Vhosts struct {
	// vhosts is the list of vhosts
	Vhosts []Vhost
	// mutex is the mutex lock for concurrent access safety
	mutex sync.RWMutex
}

type FiberMiddleware func(*fiber.Ctx) error

// NewVhost returns a new vhost with the given hostname, path, websiteID and middleware
func NewVhost(hostname, path, websiteID string, middleware FiberMiddleware) Vhost {
	return Vhost{
		Hostname:   hostname,
		Path:       path,
		WebsiteID:  websiteID,
		Middleware: middleware,
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

// get middleware returns the middleware for the vhost with the given hostname
func (v *Vhosts) getMiddleware(hostname string) (func(*fiber.Ctx) error, bool) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	for _, vhost := range v.Vhosts {
		if vhost.Hostname == hostname {
			return vhost.Middleware, true
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

	// Create the file at the given path
	saveFile, err := createFile(path)
	if err != nil {
		return err
	}
	defer saveFile.Close()

	// gob encode the vhosts list
	encoder := gob.NewEncoder(saveFile)
	err = encoder.Encode(vhosts)
	if err != nil {
		return err
	}

	return nil
}

// load loads the vhosts list from a file at the given path
func load(path string) ([]Vhost, error) {

	// load using gob decoding

	// Open the file at the given path
	loadFile, err := openFile(path)
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

// utility functions

// createFile creates a file at the given path
func createFile(path string) (*os.File, error) {

	// Create the file at the given path
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// openFile opens a file at the given path
func openFile(path string) (*os.File, error) {

	// Open the file at the given path
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// vhosts is the vhosts list
var vhosts *Vhosts

// init initializes the vhosts list
func init() {
	vhosts = &Vhosts{}
}
