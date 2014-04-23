package php

import (
	"stephensearles.com/php/ast"
	"stephensearles.com/php/token"
)

func (p *Parser) parseInstantiation() ast.Expression {
	p.expectCurrent(token.NewOperator)
	p.next()

	expr := &ast.NewExpression{}
	expr.Class = p.parseOperand()

	if p.peek().typ == token.OpenParen {
		p.expect(token.OpenParen)
		if p.peek().typ != token.CloseParen {
			expr.Arguments = append(expr.Arguments, p.parseNextExpression())
			for p.peek().typ == token.Comma {
				p.expect(token.Comma)
				expr.Arguments = append(expr.Arguments, p.parseNextExpression())
			}
		}
		p.expect(token.CloseParen)
	}
	return expr
}

func (p *Parser) parseClass() ast.Class {
	if p.current.typ == token.Abstract {
		p.expect(token.Class)
	}
	if p.current.typ == token.Final {
		p.expect(token.Class)
	}
	p.expect(token.Identifier)
	name := p.current.val
	if p.peek().typ == token.Extends {
		p.expect(token.Extends)
		p.expect(token.Identifier)
	}
	if p.peek().typ == token.Implements {
		p.expect(token.Implements)
		p.expect(token.Identifier)
		for p.peek().typ == token.Comma {
			p.expect(token.Comma)
			p.expect(token.Identifier)
		}
	}
	p.expect(token.BlockBegin)
	return p.parseClassFields(ast.Class{Name: name})
}

func (p *Parser) parseObjectLookup(r ast.Expression) (expr ast.Expression) {
	p.expectCurrent(token.ObjectOperator)
	prop := &ast.PropertyExpression{
		Receiver: r,
	}
	switch p.next(); p.current.typ {
	case token.BlockBegin:
		prop.Name = p.parseNextExpression()
		p.expect(token.BlockEnd)
	case token.VariableOperator:
		prop.Name = p.parseExpression()
	case token.Identifier:
		prop.Name = ast.Identifier{Value: p.current.val}
	}
	expr = prop
	switch pk := p.peek(); pk.typ {
	case token.OpenParen:
		expr = &ast.MethodCallExpression{
			Receiver:               r,
			FunctionCallExpression: p.parseFunctionCall(prop.Name),
		}
	}
	expr = p.parseOperation(p.parenLevel, expr)
	return
}

func (p *Parser) parseVisibility() (vis ast.Visibility, found bool) {
	switch p.peek().typ {
	case token.Private:
		vis = ast.Private
	case token.Public:
		vis = ast.Public
	case token.Protected:
		vis = ast.Protected
	default:
		return ast.Public, false
	}
	p.next()
	return vis, true
}

func (p *Parser) parseAbstract() bool {
	if p.peek().typ == token.Abstract {
		p.next()
		return true
	}
	return false
}

func (p *Parser) parseClassFields(c ast.Class) ast.Class {
	// Starting on BlockBegin
	c.Methods = make([]ast.Method, 0)
	c.Properties = make([]ast.Property, 0)
	for p.peek().typ != token.BlockEnd {
		vis, _, _, abstract := p.parseClassMemberSettings()
		p.next()
		switch p.current.typ {
		case token.Function:
			if abstract {
				f := p.parseFunctionDefinition()
				m := ast.Method{
					Visibility:   vis,
					FunctionStmt: &ast.FunctionStmt{FunctionDefinition: f},
				}
				c.Methods = append(c.Methods, m)
				p.expect(token.StatementEnd)
			} else {
				c.Methods = append(c.Methods, ast.Method{
					Visibility:   vis,
					FunctionStmt: p.parseFunctionStmt(),
				})
			}
		case token.Var:
			p.expect(token.VariableOperator)
			fallthrough
		case token.VariableOperator:
			p.expect(token.Identifier)
			prop := ast.Property{
				Visibility: vis,
				Name:       "$" + p.current.val,
			}
			if p.peek().typ == token.AssignmentOperator {
				p.expect(token.AssignmentOperator)
				prop.Initialization = p.parseNextExpression()
			}
			c.Properties = append(c.Properties, prop)
			p.expect(token.StatementEnd)
		case token.Const:
			constant := ast.Constant{}
			p.expect(token.Identifier)
			constant.Variable = ast.NewVariable(p.current.val)
			if p.peek().typ == token.AssignmentOperator {
				p.expect(token.AssignmentOperator)
				constant.Value = p.parseNextExpression()
			}
			c.Constants = append(c.Constants, constant)
			p.expect(token.StatementEnd)
		default:
			p.errorf("unexpected class member %v", p.current)
		}
	}
	p.expect(token.BlockEnd)
	return c
}

func (p *Parser) parseInterface() *ast.Interface {
	i := &ast.Interface{
		Inherits: make([]string, 0),
	}
	p.expect(token.Identifier)
	i.Name = p.current.val
	if p.peek().typ == token.Extends {
		p.expect(token.Extends)
		for {
			p.expect(token.Identifier)
			i.Inherits = append(i.Inherits, p.current.val)
			if p.peek().typ != token.Comma {
				break
			}
			p.expect(token.Comma)
		}
	}
	p.expect(token.BlockBegin)
	for p.peek().typ != token.BlockEnd {
		vis, _ := p.parseVisibility()
		if p.peek().typ == token.Static {
			p.next()
		}
		p.next()
		switch p.current.typ {
		case token.Function:
			f := p.parseFunctionDefinition()
			m := ast.Method{
				Visibility:   vis,
				FunctionStmt: &ast.FunctionStmt{FunctionDefinition: f},
			}
			i.Methods = append(i.Methods, m)
			p.expect(token.StatementEnd)
		default:
			p.errorf("unexpected interface member %v", p.current)
		}
	}
	p.expect(token.BlockEnd)
	return i
}

func (p *Parser) parseClassMemberSettings() (vis ast.Visibility, static, final, abstract bool) {
	var foundVis bool
	vis = ast.Public
	for {
		switch p.peek().typ {
		case token.Abstract:
			if abstract {
				p.errorf("found multiple abstract declarations")
			}
			abstract = true
			p.next()
		case token.Private, token.Public, token.Protected:
			if foundVis {
				p.errorf("found multiple visibility declarations")
			}
			vis, foundVis = p.parseVisibility()
		case token.Final:
			if final {
				p.errorf("found multiple final declarations")
			}
			final = true
			p.next()
		case token.Static:
			if static {
				p.errorf("found multiple static declarations")
			}
			static = true
			p.next()
		default:
			return
		}
	}
	return
}
