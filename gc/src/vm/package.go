package vm

import "gc/src/objects"

type Package struct {
	Globals map[string]*objects.Object
}
