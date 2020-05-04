package prepare

type Task interface {
	ID() string
	Run() (interface{}, error)
}
