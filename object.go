package space

import (
	"encoding/json"

	"github.com/elos/data"
	"github.com/elos/models"
)

type Object struct {
	*Space     `json:"-"`
	data.Model `json:"-"`
}

func NewObject(s *Space, k data.Kind) *Object {
	m, _ := s.ModelFor(k)
	return &Object{
		Space: s,
		Model: m,
	}
}

func (this *Object) transferAttrs() {
	bytes, _ := json.Marshal(this)
	json.Unmarshal(bytes, this.Model)
}

func (this *Object) Save() {
	this.transferAttrs()
	this.Space.Save(this.Model)
}

func (this *Object) Delete() {
	this.transferAttrs() // id has changed?
	this.Space.Delete(this.Model)
}

type User struct {
	*Object

	ID              string   `json:"id"`
	Name            string   `json:"name"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
	Key             string   `json:"key"`
	EventIDs        []string `json:"event_ids"`
	TaskIDs         []string `json:"task_ids"`
	CurrentActionID string   `json:"current_action_ids`
	ActionableKind  string   `json:"actionable_kind"`
	ActionableID    string   `json:"actionable_id"`
}

func NewUser(s *Space) *User {
	return &User{
		Object: NewObject(s, models.UserKind),
	}
}

type Action struct {
	*Object

	ID        string
	Name      string
	CreatedAt string
	UpdatedAt string
}

func NewAction(s *Space) *Action {
	return &Action{
		Object: NewObject(s, models.ActionKind),
	}
}
