package tododb

type TodoDB interface {
	GetAllTodos() ([]string, error)
	SaveTodo(string) error
	DeleteTodo(string) error
	GetHealthStatus() map[string]string
	RegisterMetrics()
}
