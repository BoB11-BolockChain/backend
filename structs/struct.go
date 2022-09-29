package structs

type Ability struct {
	Tactic         string
	Name           string
	Technique_id   string
	Description    string
	Technique_name string
	Ability_id     string
	Repeatable     bool
	Singleton      bool
	Executors      []Executor
}

type Executor struct {
	Platform string
}
