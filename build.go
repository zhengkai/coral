package coral

// BuildSimple ...
func BuildSimple(loadFn LoadFunc) (c Cache) {

	c = &simple{
		store:  make(map[interface{}]*entry),
		load:   make(map[interface{}]*entry),
		loadFn: loadFn,
	}

	return
}
