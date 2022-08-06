package core

// define repository interface for storage and parsing
// and core implementation
// every repository will have to implement this interface
type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
