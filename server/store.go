package server

type store struct {
	store map[string]string
}

func CreateStore() store {
	return store{store: map[string]interface{}{}}
}
