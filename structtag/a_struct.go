package structtag

type A struct {
	Age  int `yaml:"age"`
	Name string
}

type B struct {
	A
	addr string
}
