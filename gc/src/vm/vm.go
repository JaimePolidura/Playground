package vm

var Current VM

type VM struct {
	Threads  []*Thread
	Packages map[string]*Package
	GC       interface{}
}
