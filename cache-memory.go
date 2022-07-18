package cachememory

type CacheMemory interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Delete(key string)
}

type Cache struct {
	m map[string]interface{}
}

func New() CacheMemory {
	var c Cache
	c.m = make(map[string]interface{})
	return &c
}

func (c Cache) Set(key string, value interface{}) {
	c.m[key] = value
}

func (c Cache) Get(key string) interface{} {
	return c.m[key]
}

func (c Cache) Delete(key string) {
	delete(c.m, key)
}
