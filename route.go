package web

type uriFragment struct {
	isVariable  bool
	hasWildcard bool
	// fragment will be the key name, if it's a variable/named URI parameter
	fragment string
}
