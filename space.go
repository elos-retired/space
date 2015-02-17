package space

import (
	"github.com/elos/data"
	"github.com/elos/models/user"
	"github.com/robertkrimen/otto"
)

type ObjectKind string

type Credentials struct {
	ID  string
	Key string
}

type Space struct {
	*data.Access
}

func NewSpace(c *Credentials, store data.Store) (space *Space, err error) {
	space = new(Space)

	client, authed, err := user.Authenticate(store, c.ID, c.Key)
	if !authed {
		return
	}

	space.Access = data.NewAccess(client, store)

	return
}

func (s *Space) Expose(o *otto.Otto) {
	o.Set("User", func() interface{} { return NewUser(s) })
	o.Set("Action", func() interface{} { return NewAction(s) })
}
