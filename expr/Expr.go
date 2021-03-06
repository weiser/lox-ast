package expr
// DO NOT MODIFY. GENERATED VIA `go run cmd/tool/generateAst.go expr`
// TODO:  MAKE `cmd/tool/generateAst.go` format code
import . "github.com/weiser/lox/token"
type Expr struct {
}
type ExprInterface interface {
		Accept(evi ExprVisitorInterface) interface{}
	}
type ExprVisitorInterface interface {
VisitExpr(e *Expr) interface{}
VisitAssign(e *Assign) interface{}
VisitBinary(e *Binary) interface{}
VisitCall(e *Call) interface{}
VisitGet(e *Get) interface{}
VisitGrouping(e *Grouping) interface{}
VisitLiteral(e *Literal) interface{}
VisitLogical(e *Logical) interface{}
VisitSet(e *Set) interface{}
VisitUnary(e *Unary) interface{}
VisitVariable(e *Variable) interface{}
}
func (o *Expr) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitExpr(o)
}
type Assign struct {
*Expr
Name Token
Value ExprInterface
}
func (o *Assign) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitAssign(o)
}
type Binary struct {
*Expr
Left ExprInterface
Operator Token
Right ExprInterface
}
func (o *Binary) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitBinary(o)
}
type Call struct {
*Expr
Callee ExprInterface
Paren Token
Arguments []ExprInterface
}
func (o *Call) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitCall(o)
}
type Get struct {
*Expr
Object ExprInterface
Name Token
}
func (o *Get) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitGet(o)
}
type Grouping struct {
*Expr
Expression ExprInterface
}
func (o *Grouping) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitGrouping(o)
}
type Literal struct {
*Expr
Value interface{}
}
func (o *Literal) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitLiteral(o)
}
type Logical struct {
*Expr
Left ExprInterface
Operator Token
Right ExprInterface
}
func (o *Logical) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitLogical(o)
}
type Set struct {
*Expr
Object ExprInterface
Name Token
Value ExprInterface
}
func (o *Set) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitSet(o)
}
type Unary struct {
*Expr
Operator Token
Right ExprInterface
}
func (o *Unary) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitUnary(o)
}
type Variable struct {
*Expr
Name Token
}
func (o *Variable) Accept(evi ExprVisitorInterface) interface{} {
return evi.VisitVariable(o)
}
