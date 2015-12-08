package command

import "fmt"

type Registry map[string]Func

func New() Registry { return make(Registry) }

type ErrAlreadyRegistered string

func (e ErrAlreadyRegistered) Error() string {
	return fmt.Sprintf("command: %v is already registered", e)
}

type Func func() error

func (r Registry) Add(name string, fn Func) error {
	if _, ok := r[name]; ok {
		return ErrAlreadyRegistered(name)
	}
	r[name] = fn
	return nil
}
