package semantic

type Scope int

type Variables map[Scope][]*Variable

type Variable struct {
	Name  string
	Type  string
}

const (
	Undefined = "Undefined"
	Int       = "Integer"
	Float     = "Float"
	String    = "String"
	Bool      = "Boolean"
)

var types = map[int]string{
	0: Int,
	1: Float,
	2: String,
	3: Bool,
}

func intToType(keyCode int) string {
	if keyCode >= 0 && keyCode < len(types) {
		return types[keyCode]
	}

	return Undefined
}
