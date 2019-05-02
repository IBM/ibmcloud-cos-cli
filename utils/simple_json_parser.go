package utils

import (
	"errors"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

// Excerpted from goparsec tool that provides the
// parser combinator library in Golang

// A library to construct top-down recursive backtracking parsers using parser-combinators. Before
// proceeding you might want to take at peep at theory of parser combinators. As for this package,
// it provides:

// A standard set of combinators.
// Regular expression based simple-scanner.
// Standard set of tokenizers based on the simple-scanner.
// To construct syntax-trees based on detailed grammar try with AST struct

// Standard set of combinators are exported as methods to AST.
// Generate dot-graph EG: dotfile for html.
// Pretty print on the console.
// Make debugging easier.

// More info at https://github.com/prataprc/goparsec
var (
	commaSeparator = parsec.Atom(",", "COMMA")
	eqSeparator    = parsec.Atom("=", "EQ")
	listStart      = parsec.Atom("[", "LISTSTART")
	listEnd        = parsec.Atom("]", "LISTEND")
	objectStart    = parsec.Atom("{", "OBJECTSTART")
	objectEnd      = parsec.Atom("}", "OBJECTEND")
	trueToken      = parsec.Atom("true", "BOOL")
	falseToken     = parsec.Atom("false", "BOOL")
	simpleString   = parsec.Token(`[-.a-zA-Z0-9]+`, "SIMPLESTRING")
	key            = parsec.OrdChoice(valueNodify, parsec.String(), simpleString)
	value          parsec.Parser // recursive
	listCore       = parsec.Kleene(coreNodify, &value, commaSeparator)
	list           = parsec.And(listNodify, listStart, listCore, listEnd)
	keyValuePair   = parsec.And(keyValuePairNodify, key, eqSeparator, &value)
	objectCore     = parsec.Kleene(coreNodify, keyValuePair, commaSeparator)
	object         = parsec.And(objectNodify, objectStart, objectCore, objectEnd)

	root = parsec.Maybe(rootNodify, objectCore)
)

// Initialize parsec to compose basic set of terminal parsers, a.k.a tokenizers
// and compose them together as a tree of parsers, using combinators like: And,
// OrdChoice, Kleene, Many, Maybe.
func init() {
	value = parsec.OrdChoice(valueNodify,
		trueToken,
		falseToken,
		parsec.String(),
		parsec.Float(),
		parsec.Int(),
		list,
		object,
		simpleString,
	)
}

// GetAsJson parses the string as valid Json or not
func GetAsJson(input string) ([]byte, error) {
	scanner := parsec.NewScanner([]byte(input))
	node, scanner := root(scanner)
	scanner.SkipWS() // consume any white space in the end
	if !scanner.Endof() {
		return nil, errors.New("Error Parsing input at position " + strconv.Itoa(scanner.GetCursor()))
	}

	if str, ok := node.(string); ok {
		return []byte(str), nil
	}
	return nil, errors.New("unexpected node type")
}

// Nodify, callback function is supplied while combining parser
// functions. If the underlying parsing logic matches with i/p text,
// then callback will be dispatched with list of matching ParsecNode.
// Value returned by callback function will further be used as
// ParsecNode item in higher-level list of ParsecNodes.
func coreNodify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	if len(nodes) < 1 {
		return ""
	}
	str := nodeHandle(nodes[0])
	for i := 1; i < len(nodes); i++ {
		str += "," + nodeHandle(nodes[i])
	}
	return str
}

// Parse list and nodify
func listNodify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return "[" + nodeHandle(nodes[1]) + "]"
}

// Parse Key Value Pair and nodify
func keyValuePairNodify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	k := nodeHandle(nodes[0])
	v := nodeHandle(nodes[2])
	return k + `:` + v
}

// Parse Object and nodify
func objectNodify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return "{" + nodeHandle(nodes[1]) + "}"
}

// Parse value and nodify
func valueNodify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return nodeHandle(nodes[0])
}

// Parse root and nodify
func rootNodify(nodes []parsec.ParsecNode) parsec.ParsecNode {
	return "{" + nodeHandle(nodes[0]) + "}"
}

// in our parser only options are parsec.Terminal for terminals
// after process results, strings
func nodeHandle(node interface{}) string {
	if t, ok := node.(*parsec.Terminal); ok {
		return terminalHandle(t)
	}

	if s, ok := node.(string); ok {
		return s
	}
	panic("UNEXPECTED TYPE")
}

// terminalHandle checks if this is a terminal or not. Parsers for
// standard set of tokens are supplied along with this package.
// Most of these parsers return Terminal type as ParseNode.
func terminalHandle(t *parsec.Terminal) string {
	if !t.IsTerminal() {
		panic("NOT A TERMINAL")
	}
	if t.GetName() == "SIMPLESTRING" {
		return `"` + t.GetValue() + `"`
	} else {
		return t.Value
	}
}
