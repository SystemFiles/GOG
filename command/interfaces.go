package command

type Runnable interface {
	Init([]string) error
	Run() error
	Name() string
	Help()
}