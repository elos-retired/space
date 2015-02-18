package space

import (
	"sync"

	"github.com/elos/data"
	"github.com/elos/models"
	"github.com/elos/models/user"
	"github.com/robertkrimen/otto"
)

type ObjectKind string

type Credentials struct {
	ID  string
	Key string
}

type Space struct {
	c *Credentials
	*data.Access

	m       sync.Mutex
	objects map[Object]bool
}

func NewSpace(c *Credentials, store data.Store) (space *Space, err error) {
	space = new(Space)
	space.c = c
	space.objects = make(map[Object]bool)

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
	o.Set("Routine", func() interface{} { return NewRoutine(s) })
	o.Set("Task", func() interface{} { return NewTask(s) })

	o.Set("FindUser", s.FindUser)
	o.Set("FindAction", s.FindAction)
	o.Set("FindRoutine", s.FindRoutine)

	u := s.FindUser(s.c.ID)
	o.Set("me", u)
}

func (s *Space) FindUser(id string) *User {
	m, _ := s.Access.Unmarshal(models.UserKind, data.AttrMap{
		"id": id,
	})
	s.Access.PopulateByID(m)
	return UserModel(s, m.(models.User))
}

func (s *Space) FindAction(id string) *Action {
	m, _ := s.Access.Unmarshal(models.ActionKind, data.AttrMap{
		"id": id,
	})
	s.Access.PopulateByID(m)
	return ActionModel(s, m.(models.Action))
}

func (s *Space) FindRoutine(id string) *Routine {
	m, _ := s.Access.Unmarshal(models.RoutineKind, data.AttrMap{
		"id": id,
	})
	s.Access.PopulateByID(m)
	return RoutineModel(s, m.(models.Routine))
}

func (s *Space) FindTask(id string) *Task {
	m, _ := s.Access.Unmarshal(models.TaskKind, data.AttrMap{
		"id": id,
	})
	s.Access.PopulateByID(m)
	return TaskModel(s, m.(models.Task))
}

func (s *Space) Register(o Object) {
	s.m.Lock()
	defer s.m.Unlock()

	s.objects[o] = true
}

func (s *Space) Reload() {
	s.m.Lock()
	defer s.m.Unlock()

	for object := range s.objects {
		object.Reload()
	}
}
