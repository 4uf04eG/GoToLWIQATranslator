package lexer

type Pair struct {
	Key   interface{}
	Value interface{}
}

const (
	Keyword    = "Keyword"
	Comment    = "Comment"
	Operator   = "Operator"
	Delimiter  = "Delimiter"
	Identifier = "Identifier"
	Literal    = "Literal"
    EndOfLine  = "EOL"
	Invalid    = "Invalid"
)

var patterns = []Pair{
	{Keyword,   "switch|case|default|var|for|break|continue|return|if|else"},
	{Comment,   "(//|/\\*).*"},
	{Operator,  ":=|==|<=|>=|=|\\+\\+|--|\\+|-|\\*|/|%|>|<|!"},
	{Delimiter, "[{}():;,]"},
    {Literal,   "\\d+[\\.exobEXOB]?\\d*|true|false|\"[^\"]*\"|'\\\\?.'"},
	{Identifier,"[a-zA-Z_]\\w*"},
    {EndOfLine, "[\r\n]"},
	{Invalid,   "\\S+?"},
}
