package vhosts

import (
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

// vhosts is the vhosts list
var Vhs *Vhosts

// Vhost is a virtual host
type Vhost struct {
	Hostname     string            // hostname is the hostname of the vhost
	Path         string            // path is the path of the vhost
	WebsiteID    string            // websiteID is the websiteID of the vhost
	ErrorHandler FiberErrorHandler // errorHandler is the error handler for the vhost
	Handler      FiberHandler      // middleware is the middleware for the vhost
	LastModified int64             // lastModified is the last modified time of the vhost
}

// vhosts contains all the vhosts protected by mutex lock for concurrent access safety
type Vhosts struct {
	// vhosts is the list of vhosts
	Vhosts []Vhost
	// LastModified is the last modified time of the vhosts file
	LastModified int64
	// Version is the version of the vhosts file ( quick way to check if the vhosts file has changed )
	Version int64
	// Checksum is the checksum of the vhosts file ( quick way to check if the vhosts file has changed )
	Checksum string
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
		LastModified: time.Now().Unix(),
	}
}

// add adds a vhost to the vhosts list
func (v *Vhosts) Add(vhost Vhost) error {
	// lookup the vhost by hostname and return error if it already exists
	_, ok := v.Get(vhost.Hostname)
	if ok {
		return errors.New("vhost already exists")
	}
	v.mutex.Lock()
	defer v.mutex.Unlock()
	v.Vhosts = append(v.Vhosts, vhost)
	// update the vhosts list version and last modified time
	v.Version = +1
	v.LastModified = time.Now().Unix()

	return nil
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
func (v *Vhosts) Remove(hostname string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	for i, vhost := range v.Vhosts {
		if vhost.Hostname == hostname {
			v.Vhosts = append(v.Vhosts[:i], v.Vhosts[i+1:]...)
			// update the vhosts list version and last modified time
			v.Version = +1
			v.LastModified = time.Now().Unix()
			return nil
		}
	}
	return errors.New("vhost not found")
}

// NumberOfVhosts returns the length of the vhosts list
func (v *Vhosts) NumberOfVhosts() int {
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

// GetVhostnames returns the hostnames list ( []string )
func GetVhostnames(v ...*Vhosts) []string {
	var vh []Vhost
	for _, vhost := range v {
		vh = append(vh, vhost.getVhosts()...)
	}
	var hostnames []string
	for _, vhost := range vh {
		hostnames = append(hostnames, vhost.Hostname)
	}
	return hostnames
}

// GetVhostnames returns the hostnames list ( []string )
func (v *Vhosts) GetVhostnames() []string {
	var hostnames []string
	for _, vhost := range v.Vhosts {
		hostnames = append(hostnames, vhost.Hostname)
	}
	return hostnames
}

// getHandler returns the handler for the given hostname
func (v *Vhosts) getHandler(hostname string) (FiberHandler, bool) {
	v.mutex.RLock()
	defer v.mutex.RUnlock()
	for _, vhost := range v.Vhosts {
		if vhost.Hostname == hostname {
			return vhost.Handler, true
		}
	}
	return nil, false
}

// Save saves the vhosts to the given file
func (v *Vhosts) Save(file string) error {
	v.mutex.RLock()
	defer v.mutex.RUnlock()

	// hash the vhosts list
	hash, err := Hash(v.Vhosts)
	if err != nil {
		return err
	}

	// update the vhosts list checksum
	v.Checksum = hash

	return save(file, v)
}

// save saves the vhosts to the given file
func save(file string, v *Vhosts) error {
	return EncodeAsGob(file, v)
}

// EncodeAsGob encodes the given vhosts as gob and saves it to the given file
func EncodeAsGob(file string, v *Vhosts) error {

	// Open the file at the given path
	saveFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer saveFile.Close()

	encoder := gob.NewEncoder(saveFile)
	err = encoder.Encode(v)
	if err != nil {
		return err
	}
	return nil
}

// Load loads the vhosts from the given file
func (v *Vhosts) Load(file string) error {

	// does the file we're trying to load exist?
	if !doesFileExist(file) {
		return errors.New("file doesn't exist")
	}

	// load the vhosts from the given file
	err := load(file, v)
	if err != nil {
		return err
	}

	// // set the vhosts
	// v.mutex.Lock()
	// defer v.mutex.Unlock()
	// v.Vhosts = vhosts

	return nil
}

// load loads the vhosts from the given file into the pointer to vhosts
func load(file string, vhPtr *Vhosts) error {

	// Open the file at the given path
	loadFile, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}

	// gob decode the vhosts list
	decoder := gob.NewDecoder(loadFile)
	err = decoder.Decode(vhPtr)
	if err != nil {
		return err
	}

	// verify the vhosts list checksum
	hash, err := Hash(vhPtr.Vhosts)
	if err != nil {
		return err
	}
	if hash != vhPtr.Checksum {
		return errors.New("vhosts list checksum doesn't match")
	}

	return nil

}

// Hash returns the hash of the given vhosts list
func Hash(vhosts []Vhost) (string, error) {

	var hashes []string

	// hash the vhosts list
	for _, vhost := range vhosts {

		var vhostHash string
		// create a new sha256 hash
		h := sha256.New()
		// hash the vhost hostname
		_, err := h.Write([]byte(vhost.Hostname))
		if err != nil {
			return "", err
		}
		vhostHash = string(h.Sum(nil))

		// hash the websiteID
		_, err = h.Write([]byte(vhost.WebsiteID))
		if err != nil {
			return "", err
		}
		websiteIdHash := string(h.Sum(nil))

		// combine the vhost hash and the websiteID hash
		vhostHash = vhostHash + websiteIdHash

		// add the vhost hash to the hashes list
		hashes = append(hashes, vhostHash)

	}

	// sort the hashes list alphabetically ( so that the order of the vhosts doesn't matter )
	sort.Strings(hashes)

	// create a new sha256 hash
	h := sha256.New()

	// combine all the hashes into one string and hash it
	var combinedString string
	for _, hash := range hashes {
		combinedString = combinedString + hash
	}
	_, err := h.Write([]byte(combinedString))
	if err != nil {
		return "", err
	}

	// return the hash as a string
	return string(h.Sum(nil)), nil

}

// init initializes the vhosts list
func init() {
	Vhs = &Vhosts{}
}

// Init initializes the vhosts list from a file at the given path
func InitVHostDataFile(path string) error {
	return Vhs.Load(path)
}

// Initialize initializes the vhosts list with some vhosts defaults map of hostname to middleware ( map[string]func(*fiber.Ctx) error )
func Initialize(listOfHostnames map[string]map[string]interface{}) {

	// Add the vhosts to the vhosts list
	for hostname, middleware := range listOfHostnames {
		Vhs.Add(NewVhost(hostname, "", "", middleware["handler"].(func(*fiber.Ctx) error), middleware["errorHandler"].(func(*fiber.Ctx, error) error)))
	}
}

// utility functions
// vhostReset resets the vhosts list
func vhostReset() {
	Vhs = &Vhosts{}
}

// doesFileExist checks if a file exists at the given path
func doesFileExist(path string) bool {
	// return true if the file already exists, if not return false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetVhosts returns the vhosts
func GetVhosts() *Vhosts {

	return Vhs
}

// SetHandler sets the handler for the given hostname
func SetHandler(hostname string, handler FiberHandler) error {
	vhost, ok := Vhs.Get(hostname)
	if !ok {
		return errors.New("vhost not found")
	}
	vhost.Handler = handler
	return nil
}

// SetErrorHandler sets the error handler for the given hostname
func SetErrorHandler(hostname string, errorHandler FiberErrorHandler) error {
	vhost, ok := Vhs.Get(hostname)
	if !ok {
		return errors.New("vhost not found")
	}
	vhost.ErrorHandler = errorHandler
	return nil
}
