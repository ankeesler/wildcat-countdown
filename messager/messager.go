package messager

// Messager is an object from which a caller can get a message.
type Messager struct{}

// New instantiates a new Messager.
func New() *Messager { return &Messager{} }

// Message returns a message created by a Getter.
func (g *Messager) Message() (string, error) {
	return "Go 'Cats!", nil
}
