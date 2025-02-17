package ast

type Modifier func(Node) (Node, error)

// Walker walks the input AST and applies modifier to each node.
// It returns a AST copy and does not mutate the input AST.
func Walker(node Node, modifier Modifier) (Node, error) {
	if node == nil {
		return nil, nil
	}

	switch n := node.(type) {
	case *Program:
		var newStatements []Statement
		for _, stmt := range n.Statements {
			modifiedStmt, err := Walker(stmt, modifier)
			if err != nil {
				return nil, err
			}
			if modifiedStmt != nil {
				newStatements = append(newStatements, modifiedStmt.(Statement))
			}
		}
		return modifier(&Program{Statements: newStatements})

	case *LetStatement:
		mRight, err := Walker(n.Right, modifier)
		if err != nil {
			return nil, err
		}
		return modifier(&LetStatement{
			Token: n.Token,
			Name:  n.Name, Right: mRight.(Expression)})

	case *ExpressionStatement:
		mExpr, err := Walker(n.Expr, modifier)
		if err != nil {
			return nil, err
		}
		return modifier(&ExpressionStatement{
			Token: n.Token,
			Expr:  mExpr.(Expression)})

	case *PrefixExpression:
		mRight, err := Walker(n.Right, modifier)
		if err != nil {
			return nil, err
		}
		return modifier(&PrefixExpression{
			Token:    n.Token,
			Operator: n.Operator, Right: mRight.(Expression)})

	case *InfixExpression:
		mLeft, err := Walker(n.Left, modifier)
		if err != nil {
			return nil, err
		}
		mRight, err := Walker(n.Right, modifier)
		if err != nil {
			return nil, err
		}
		return modifier(&InfixExpression{
			Token: n.Token,
			Left:  mLeft.(Expression), Operator: n.Operator, Right: mRight.(Expression)})

	case *CallExpression:
		mFunc, err := Walker(n.Function, modifier)
		if err != nil {
			return nil, err
		}
		mArgs := make([]Expression, len(n.Arguments))
		for i, arg := range n.Arguments {
			mArg, err := Walker(arg, modifier)
			if err != nil {
				return nil, err
			}
			mArgs[i] = mArg.(Expression)
		}
		return modifier(&CallExpression{
			Token:    n.Token,
			Function: mFunc.(Expression), Arguments: mArgs})

	case *FunctionLiteral:
		mParams := make([]*Identifier, len(n.Parameters))
		for i, param := range n.Parameters {
			mParam, err := Walker(param, modifier)
			if err != nil {
				return nil, err
			}
			mParams[i] = mParam.(*Identifier)
		}
		mBody, err := Walker(n.Body, modifier)
		if err != nil {
			return nil, err
		}
		return modifier(&FunctionLiteral{
			Token:      n.Token,
			Parameters: mParams, Body: mBody.(*BlockStatement)})

	case *IfElseConditional:
		mCondition, err := Walker(n.Condition, modifier)
		if err != nil {
			return nil, err
		}
		mConsequence, err := Walker(n.Consequence, modifier)
		if err != nil {
			return nil, err
		}
		var mAlternative *BlockStatement
		if n.Alternative != nil {
			mAlternativeNode, err := Walker(n.Alternative, modifier)
			if err != nil {
				return nil, err
			}
			mAlternative = mAlternativeNode.(*BlockStatement)
		}
		return modifier(&IfElseConditional{
			Token:     n.Token,
			Condition: mCondition.(Expression), Consequence: mConsequence.(*BlockStatement), Alternative: mAlternative})

	case *BlockStatement:
		var newStatements []Statement
		for _, stmt := range n.Statements {
			modifiedStmt, err := Walker(stmt, modifier)
			if err != nil {
				return nil, err
			}
			if modifiedStmt != nil {
				newStatements = append(newStatements, modifiedStmt.(Statement))
			}
		}
		return modifier(&BlockStatement{Statements: newStatements})

	default:
		return modifier(n)
	}
}
