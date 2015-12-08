package cvar

import "fmt"

type Registry map[string]interface{}

func New() Registry { return make(Registry) }

type options struct {
	Saved, ServerSide bool
}

func (o *options) Apply(os ...Option) {
	for _, opt := range os {
		opt(o)
	}
}

type String struct {
	val  string
	opts options
}

func (s *String) Get() string { return s.val }

func (r Registry) NewString(name, def string, os ...Option) (*String, error) {
	if _, ok := r[name]; ok {
		return nil, ErrAlreadyRegistered(name)
	}
	cvar := &String{val: def}
	cvar.opts.Apply(os...)
	r[name] = cvar
	return cvar, nil
}

type Float struct {
	val  float64
	opts options
}

func (f *Float) Get() float64 { return f.val }

func (r Registry) NewFloat(name string, def float64, os ...Option) (*Float, error) {
	if _, ok := r[name]; ok {
		return nil, ErrAlreadyRegistered(name)
	}
	cvar := &Float{val: def}
	cvar.opts.Apply(os...)
	r[name] = cvar
	return cvar, nil
}

type Option func(*options)

var Saved Option = func(o *options) { o.Saved = true }
var ServerSide Option = func(o *options) { o.ServerSide = true }

type ErrAlreadyRegistered string

func (e ErrAlreadyRegistered) Error() string {
	return fmt.Sprintf("command: %v is already registered", e)
}
