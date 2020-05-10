package server

type store struct {
	store map[string]interface{}
}

func CreateStore() store {
	return store{store: map[string]interface{}{}}
}
