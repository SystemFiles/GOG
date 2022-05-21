package cmd

type Runnable interface {
	Init([]string) error
	Run() error
	Name() string
	Alias() string
	Help()
}