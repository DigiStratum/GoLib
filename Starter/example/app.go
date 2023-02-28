package main

type app struct {
	svc		ServiceIfc
}

func NewApp() *app {
	return  &app{
		svc:		NewService(),
	}
}

func (r *app) DoSomething() {
	r.svc.DoSomething()
	r.svc.Start()
	r.svc.DoSomething()
}

