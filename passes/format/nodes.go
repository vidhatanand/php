package format

import (
	"io"
	"strings"

	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
)

func (f *formatWalker) Walk(node ast.Node) error {
	switch n := node.(type) {
	case *ast.IfStmt:
		f.printTab()
		f.printToken(token.If)
		f.print(" ")
		f.printToken(token.OpenParen)
		f.print("<expression>")
		f.printToken(token.CloseParen)
		f.print(" ")
		f.printToken(token.BlockBegin)
		f.print("\n")
		f.tabLevel += 1
		f.printTab()
		f.Walk(n.TrueBranch)
		f.print("\n")
		f.tabLevel -= 1
		f.printToken(token.BlockEnd)
		if n.FalseBranch != nil {
			f.print(" ")
			f.printToken(token.Else)
			f.print(" ")
			f.printToken(token.BlockBegin)
			f.print("\n")
			f.tabLevel += 1
			f.Walk(n.FalseBranch)
			f.tabLevel -= 1
			f.print("\n")
			f.printToken(token.BlockEnd)
			f.print("\n")
		}
	}
	return nil
}

func (f *formatWalker) print(s string) {
	io.WriteString(f.w, s)
}

func (f *formatWalker) printToken(t token.Token) {
	if s, ok := tokenMap[t]; ok {
		io.WriteString(f.w, s)
		return
	}
	io.WriteString(f.w, t.String())
}

func (f *formatWalker) printTab() {
	io.WriteString(f.w, strings.Repeat(f.Indent, f.tabLevel))
}

var tokenMap = map[token.Token]string{
	token.Class:               "class",
	token.UnaryOperator:       "clone",
	token.Const:               "const",
	token.Abstract:            "abstract",
	token.Interface:           "interface",
	token.Implements:          "implements",
	token.Extends:             "extends",
	token.NewOperator:         "new",
	token.If:                  "if",
	token.Else:                "else",
	token.ElseIf:              "elseif",
	token.While:               "while",
	token.Do:                  "do",
	token.For:                 "for",
	token.Foreach:             "foreach",
	token.Switch:              "switch",
	token.EndIf:               "endif;",
	token.EndFor:              "endfor;",
	token.EndForeach:          "endforeach;",
	token.EndWhile:            "endwhile;",
	token.EndSwitch:           "endswitch;",
	token.Case:                "case",
	token.Break:               "break",
	token.Continue:            "continue",
	token.Default:             "default",
	token.Function:            "function",
	token.Static:              "static",
	token.Final:               "final",
	token.Self:                "self",
	token.Parent:              "parent",
	token.Return:              "return",
	token.BlockBegin:          "{",
	token.BlockEnd:            "}",
	token.StatementEnd:        ";",
	token.OpenParen:           "(",
	token.CloseParen:          ")",
	token.Comma:               ",",
	token.Echo:                "echo",
	token.Throw:               "throw",
	token.Try:                 "try",
	token.Catch:               "catch",
	token.Finally:             "finally",
	token.Private:             "private",
	token.Public:              "public",
	token.Protected:           "protected",
	token.InstanceofOperator:  "instanceof",
	token.Global:              "global",
	token.List:                "list",
	token.Array:               "array",
	token.Exit:                "exit",
	token.IgnoreErrorOperator: "@",
	token.Null:                "null",
	token.Var:                 "var",

	token.Use:       "use",
	token.Namespace: "namespace",

	token.ObjectOperator:          "->",
	token.ScopeResolutionOperator: "::",

	token.ArrayKeyOperator: "=>",

	token.AssignmentOperator:    "=",
	token.NegationOperator:      "!",
	token.AdditionOperator:      "+",
	token.SubtractionOperator:   "-",
	token.ConcatenationOperator: ".",

	token.AndOperator:        "&&",
	token.OrOperator:         "||",
	token.AmpersandOperator:  "&",
	token.BitwiseXorOperator: "^",
	token.BitwiseNotOperator: "~",
	token.BitwiseOrOperator:  "|",
	token.TernaryOperator1:   "?",
	token.TernaryOperator2:   ":",
	token.WrittenAndOperator: "and",
	token.WrittenXorOperator: "xor",
	token.WrittenOrOperator:  "or",
	token.AsOperator:         "as",

	token.ArrayLookupOperatorLeft:  "[",
	token.ArrayLookupOperatorRight: "]",

	token.VariableOperator: "$",
}