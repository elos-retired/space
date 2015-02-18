package space

import "github.com/elos/models"

type Routine struct {
	space *Space         `json:"-"`
	model models.Routine `json:"-"`

	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
	StartTime        string   `json:"start_time"`
	EndTime          string   `json:"end_time"`
	UserID           string   `json:"user_id"`
	TaskIDs          []string `json:"task_ids"`
	CompletedTaskIDs []string `json:"completed_task_ids"`
	ActionIDs        []string `json:"action_ids"`
	CurrentActionID  string   `json:"current_action_id"`
}

func (this *Routine) Save() {
	transferAttrs(this, this.model)
	this.space.Save(this.model)
	this.space.Reload()
}

func (this *Routine) Delete() error {
	transferAttrs(this, this.model)
	return this.space.Delete(this.model)
}

func (this *Routine) Reload() error {
	this.space.Access.PopulateByID(this.model)
	transferAttrs(this.model, this)
	return nil
}

func (this *Routine) Model() models.Routine {
	return this.model
}

func NewRoutine(s *Space) *Routine {
	r, _ := s.Access.ModelFor(models.RoutineKind)
	r.SetID(s.NewID())
	return RoutineModel(s, r.(models.Routine))
}

func RoutineModel(s *Space, m models.Routine) *Routine {
	r := &Routine{
		space: s,
		model: m,
	}

	transferAttrs(r.model, r)
	s.Register(r)
	return r
}

func (this *Routine) AddTask(t *Task) {
	other := t.Model()
	transferAttrs(t, other)
	this.model.IncludeTask(other)
	this.space.Access.Save(other)
	this.space.Access.Save(this.model)
	this.space.Reload()
}
