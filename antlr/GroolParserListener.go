package antlr

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/golang-collections/collections/stack"
	"github.com/newm4n/grool/antlr/parser"
	"github.com/newm4n/grool/model"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

var (
	log = logrus.WithFields(logrus.Fields{
		"Package": "antlr",
		"Struct":  "GroolParserListener",
	})
)

func NewGroolParserListener() *GroolParserListener {
	return &GroolParserListener{
		Stack:       stack.New(),
		RuleEntries: make(map[string]*model.RuleEntry),
		ParseErrors: make([]error, 0),
	}
}

type GroolParserListener struct {
	parser.BasegroolListener
	ParseErrors []error

	RuleEntries map[string]*model.RuleEntry
	Stack       *stack.Stack
}

func (s *GroolParserListener) AddError(e error) {
	log.Errorf("Got error : %v", e)
	s.ParseErrors = append(s.ParseErrors, e)
}

// VisitTerminal is called when a terminal node is visited.
func (s *GroolParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *GroolParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *GroolParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *GroolParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterRoot is called when production root is entered.
func (s *GroolParserListener) EnterRoot(ctx *parser.RootContext) {}

// ExitRoot is called when production root is exited.
func (s *GroolParserListener) ExitRoot(ctx *parser.RootContext) {}

// EnterRuleEntry is called when production ruleEntry is entered.
func (s *GroolParserListener) EnterRuleEntry(ctx *parser.RuleEntryContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	entry := &model.RuleEntry{}
	s.Stack.Push(entry)
}

// ExitRuleEntry is called when production ruleEntry is exited.
func (s *GroolParserListener) ExitRuleEntry(ctx *parser.RuleEntryContext) {
	entry := s.Stack.Pop().(*model.RuleEntry)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	// check for duplicate rule.
	if _, ok := s.RuleEntries[entry.RuleName]; ok {
		s.AddError(fmt.Errorf("duplicate rule name '%s'", entry.RuleName))
		return
	}
	// if everything ok, add the rule entry.
	s.RuleEntries[entry.RuleName] = entry
}

// EnterRuleName is called when production ruleName is entered.
func (s *GroolParserListener) EnterRuleName(ctx *parser.RuleNameContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	ruleName := ctx.GetText()
	entry := s.Stack.Peek().(*model.RuleEntry)
	entry.RuleName = ruleName
}

// ExitRuleName is called when production ruleName is exited.
func (s *GroolParserListener) ExitRuleName(ctx *parser.RuleNameContext) {}

// EnterRuleDescription is called when production ruleDescription is entered.
func (s *GroolParserListener) EnterRuleDescription(ctx *parser.RuleDescriptionContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	ruleDescription := ctx.GetText()
	entry := s.Stack.Peek().(*model.RuleEntry)
	entry.RuleDescription = ruleDescription
}

// ExitRuleDescription is called when production ruleDescription is exited.
func (s *GroolParserListener) ExitRuleDescription(ctx *parser.RuleDescriptionContext) {}

// EnterWhenScope is called when production whenScope is entered.
func (s *GroolParserListener) EnterWhenScope(ctx *parser.WhenScopeContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	whenScope := &model.WhenScope{}
	s.Stack.Push(whenScope)
}

// ExitWhenScope is called when production whenScope is exited.
func (s *GroolParserListener) ExitWhenScope(ctx *parser.WhenScopeContext) {
	whenScope := s.Stack.Pop().(*model.WhenScope)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	ruleEntry := s.Stack.Peek().(*model.RuleEntry)
	ruleEntry.WhenScope = whenScope
}

// EnterThenScope is called when production thenScope is entered.
func (s *GroolParserListener) EnterThenScope(ctx *parser.ThenScopeContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	thenScope := &model.ThenScope{}
	s.Stack.Push(thenScope)
}

// ExitThenScope is called when production thenScope is exited.
func (s *GroolParserListener) ExitThenScope(ctx *parser.ThenScopeContext) {
	thenScope := s.Stack.Pop().(*model.ThenScope)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	ruleEntry := s.Stack.Peek().(*model.RuleEntry)
	ruleEntry.ThenScope = thenScope
}

// EnterAssignExpressions is called when production assignExpressions is entered.
func (s *GroolParserListener) EnterAssignExpressions(ctx *parser.AssignExpressionsContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	assigns := &model.AssignExpressions{
		ExpressionList: make([]*model.AssignExpression, 0),
	}
	s.Stack.Push(assigns)
}

// ExitAssignExpressions is called when production assignExpressions is exited.
func (s *GroolParserListener) ExitAssignExpressions(ctx *parser.AssignExpressionsContext) {
	assigns := s.Stack.Pop().(*model.AssignExpressions)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	thenScope := s.Stack.Peek().(*model.ThenScope)
	thenScope.AssignExpressions = assigns
}

// EnterAssignExpression is called when production assignExpression is entered.
func (s *GroolParserListener) EnterAssignExpression(ctx *parser.AssignExpressionContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	assign := &model.AssignExpression{}
	s.Stack.Push(assign)
}

// ExitAssignExpression is called when production assignExpression is exited.
func (s *GroolParserListener) ExitAssignExpression(ctx *parser.AssignExpressionContext) {
	assign := s.Stack.Pop().(*model.AssignExpression)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	assigns := s.Stack.Peek().(*model.AssignExpressions)
	assigns.ExpressionList = append(assigns.ExpressionList, assign)
}

// EnterAssignment is called when production assignment is entered.
func (s *GroolParserListener) EnterAssignment(ctx *parser.AssignmentContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	assignment := &model.Assignment{}
	s.Stack.Push(assignment)
}

// ExitAssignment is called when production assignment is exited.
func (s *GroolParserListener) ExitAssignment(ctx *parser.AssignmentContext) {
	assignment := s.Stack.Pop().(*model.Assignment)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	assign := s.Stack.Peek().(*model.AssignExpression)
	assign.Assignment = assignment
}

// EnterExpression is called when production expression is entered.
func (s *GroolParserListener) EnterExpression(ctx *parser.ExpressionContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	expression := &model.Expression{}
	s.Stack.Push(expression)
}

// ExitExpression is called when production expression is exited.
func (s *GroolParserListener) ExitExpression(ctx *parser.ExpressionContext) {
	expr := s.Stack.Pop().(*model.Expression)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	holder := s.Stack.Peek().(model.ExpressionHolder)
	err := holder.AcceptExpression(expr)
	if err != nil {
		s.AddError(err)
	}
}

// EnterPredicate is called when production predicate is entered.
func (s *GroolParserListener) EnterPredicate(ctx *parser.PredicateContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	predicate := &model.Predicate{}
	s.Stack.Push(predicate)
}

// ExitPredicate is called when production predicate is exited.
func (s *GroolParserListener) ExitPredicate(ctx *parser.PredicateContext) {
	predicate := s.Stack.Pop().(*model.Predicate)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	expr := s.Stack.Peek().(*model.Expression)
	expr.Predicate = predicate
}

// EnterExpressionAtom is called when production expressionAtom is entered.
func (s *GroolParserListener) EnterExpressionAtom(ctx *parser.ExpressionAtomContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	exprAtom := &model.ExpressionAtom{}
	s.Stack.Push(exprAtom)
}

// ExitExpressionAtom is called when production expressionAtom is exited.
func (s *GroolParserListener) ExitExpressionAtom(ctx *parser.ExpressionAtomContext) {
	exprAtom := s.Stack.Pop().(*model.ExpressionAtom)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	holder := s.Stack.Peek().(model.ExpressionAtomHolder)
	err := holder.AcceptExpressionAtom(exprAtom)
	if err != nil {
		s.AddError(err)
	}
}

// EnterFunctionCall is called when production functionCall is entered.
func (s *GroolParserListener) EnterFunctionCall(ctx *parser.FunctionCallContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	funcCall := &model.FunctionCall{}
	s.Stack.Push(funcCall)
}

// ExitFunctionCall is called when production functionCall is exited.
func (s *GroolParserListener) ExitFunctionCall(ctx *parser.FunctionCallContext) {
	funcCall := s.Stack.Pop().(*model.FunctionCall)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	holder := s.Stack.Peek().(model.FunctionCallHolder)
	err := holder.AcceptFunctionCall(funcCall)
	if err != nil {
		s.AddError(err)
	}
}

// EnterFunctionArgs is called when production functionArgs is entered.
func (s *GroolParserListener) EnterFunctionArgs(ctx *parser.FunctionArgsContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	funcArg := &model.FunctionArgument{
		Arguments: make([]*model.ArgumentHolder, 0),
	}
	s.Stack.Push(funcArg)
}

// ExitFunctionArgs is called when production functionArgs is exited.
func (s *GroolParserListener) ExitFunctionArgs(ctx *parser.FunctionArgsContext) {
	funcArgs := s.Stack.Pop().(*model.FunctionArgument)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	funcCall := s.Stack.Peek().(*model.FunctionCall)
	funcCall.FunctionArguments = funcArgs
}

// EnterLogicalOperator is called when production logicalOperator is entered.
func (s *GroolParserListener) EnterLogicalOperator(ctx *parser.LogicalOperatorContext) {
}

// ExitLogicalOperator is called when production logicalOperator is exited.
func (s *GroolParserListener) ExitLogicalOperator(ctx *parser.LogicalOperatorContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	expr := s.Stack.Peek().(*model.Expression)
	if ctx.GetText() == "&&" {
		expr.LogicalOperator = model.LogicalOperatorAnd
	} else if ctx.GetText() == "||" {
		expr.LogicalOperator = model.LogicalOperatorOr
	} else {
		s.AddError(fmt.Errorf("unknown logical operator %s", ctx.GetText()))
	}
}

// EnterVariable is called when production variable is entered.
func (s *GroolParserListener) EnterVariable(ctx *parser.VariableContext) {}

// ExitVariable is called when production variable is exited.
func (s *GroolParserListener) ExitVariable(ctx *parser.VariableContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	varName := ctx.GetText()
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	holder := s.Stack.Peek().(model.VariableHolder)
	err := holder.AcceptVariable(varName)
	if err != nil {
		s.AddError(err)
	}
}

// EnterMathOperator is called when production mathOperator is entered.
func (s *GroolParserListener) EnterMathOperator(ctx *parser.MathOperatorContext) {
}

// ExitMathOperator is called when production mathOperator is exited.
func (s *GroolParserListener) ExitMathOperator(ctx *parser.MathOperatorContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	expr := s.Stack.Peek().(*model.ExpressionAtom)
	if ctx.GetText() == "+" {
		expr.MathOperator = model.MathOperatorPlus
	} else if ctx.GetText() == "-" {
		expr.MathOperator = model.MathOperatorMinus
	} else if ctx.GetText() == "/" {
		expr.MathOperator = model.MathOperatorDiv
	} else if ctx.GetText() == "*" {
		expr.MathOperator = model.MathOperatorMul
	} else {
		s.AddError(fmt.Errorf("unknown mathematic operator %s", ctx.GetText()))
	}
}

// EnterComparisonOperator is called when production comparisonOperator is entered.
func (s *GroolParserListener) EnterComparisonOperator(ctx *parser.ComparisonOperatorContext) {}

// ExitComparisonOperator is called when production comparisonOperator is exited.
func (s *GroolParserListener) ExitComparisonOperator(ctx *parser.ComparisonOperatorContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	predicate := s.Stack.Peek().(*model.Predicate)
	if ctx.GetText() == "==" {
		predicate.ComparisonOperator = model.ComparisonOperatorEQ
	} else if ctx.GetText() == "!=" {
		predicate.ComparisonOperator = model.ComparisonOperatorNEQ
	} else if ctx.GetText() == "<" {
		predicate.ComparisonOperator = model.ComparisonOperatorLT
	} else if ctx.GetText() == "<=" {
		predicate.ComparisonOperator = model.ComparisonOperatorLTE
	} else if ctx.GetText() == ">" {
		predicate.ComparisonOperator = model.ComparisonOperatorGT
	} else if ctx.GetText() == ">=" {
		predicate.ComparisonOperator = model.ComparisonOperatorGTE
	} else {
		s.AddError(fmt.Errorf("unknown comparison operator %s", ctx.GetText()))
	}
}

// EnterConstant is called when production constant is entered.
func (s *GroolParserListener) EnterConstant(ctx *parser.ConstantContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	cons := &model.Constant{}
	s.Stack.Push(cons)
}

// ExitConstant is called when production constant is exited.
func (s *GroolParserListener) ExitConstant(ctx *parser.ConstantContext) {
	cons := s.Stack.Pop().(*model.Constant)
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	if cons.DataType == model.DecimalDataType && ctx.GetText()[:1] == "-" {
		cons.DecimalValue = cons.DecimalValue * -1
	}
	if ctx.REAL_LITERAL() != nil {
		cons.DataType = model.FloatDataType
		flo, err := strconv.ParseFloat(ctx.GetText(), 64)
		if err != nil {
			s.AddError(fmt.Errorf("string to float conversion error. String is not real type '%s'", ctx.GetText()))
			return
		} else {
			cons.FloatValue = flo
		}
	} else if ctx.NULL_LITERAL() != nil {
		if ctx.NOT() != nil {
			cons.IsNull = false
		} else {
			cons.IsNull = true
		}
	}
	holder := s.Stack.Peek().(model.ConstantHolder)
	err := holder.AcceptConstant(cons)
	if err != nil {
		s.AddError(err)
	}
}

// EnterDecimalLiteral is called when production decimalLiteral is entered.
func (s *GroolParserListener) EnterDecimalLiteral(ctx *parser.DecimalLiteralContext) {}

// ExitDecimalLiteral is called when production decimalLiteral is exited.
func (s *GroolParserListener) ExitDecimalLiteral(ctx *parser.DecimalLiteralContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	cons := s.Stack.Peek().(*model.Constant)
	i64, err := strconv.ParseInt(ctx.GetText(), 10, 64)
	if err != nil {
		s.AddError(fmt.Errorf("string to integer conversion error. literal is not a decimal '%s'", ctx.GetText()))
	} else {
		cons.DecimalValue = i64
		cons.DataType = model.DecimalDataType
	}
}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *GroolParserListener) EnterStringLiteral(ctx *parser.StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *GroolParserListener) ExitStringLiteral(ctx *parser.StringLiteralContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	cons := s.Stack.Peek().(*model.Constant)
	cons.StringValue = ctx.GetText()
	cons.DataType = model.StringDataType
}

// EnterBooleanLiteral is called when production booleanLiteral is entered.
func (s *GroolParserListener) EnterBooleanLiteral(ctx *parser.BooleanLiteralContext) {
}

// ExitBooleanLiteral is called when production booleanLiteral is exited.
func (s *GroolParserListener) ExitBooleanLiteral(ctx *parser.BooleanLiteralContext) {
	// return immediately when there's an error
	if len(s.ParseErrors) > 0 {
		return
	}
	cons := s.Stack.Peek().(*model.Constant)
	val := strings.ToLower(ctx.GetText())
	if val == "true" {
		cons.BooleanValue = true
		cons.DataType = model.BooleanDataType
	} else if val == "false" {
		cons.BooleanValue = false
		cons.DataType = model.BooleanDataType
	} else {
		s.AddError(fmt.Errorf("unknown boolear literal '%s'", ctx.GetText()))
	}
}
