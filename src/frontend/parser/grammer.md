## Hamlet's Syntactic Grammer

> The following grammer is represented in a flavour of [BNF](https://en.wikipedia.org/wiki/Backus-Naur_form), having some special _specifier_ keywords.

**Literal** &rarr; _one-of { string | character | NumericalLiteral | BooleanLiteral | empty }_

**BooleanLiteral** &rarr; _one-of { true | false }_

**NumericalLiteral** &rarr; _one-of { real | integer }_

**BinOp** &rarr; _one-of { EqualityOp | ComparisonOp | TermOp | FactorOp }_

**EqualityOp** &rarr; _one-of { '==' | '!=' }_

**ComparisonOp** &rarr; _one-of { '<' | '<=' | '>' | '>=' }_

**TermOp** &rarr; _one-of { '+' | '-' }_

**FactorOp** &rarr; _one-of { '*' | '/' }_

**UniOp** &rarr; _one-of { '~' | '!' }_

**Expression** &rarr; _EqualityExpression_

**EqualityExpression** &rarr; _ComparisonExpression zero-or-many-of { one-of EqualityOp ComparisonExpression }_

**ComparisonExpression** &rarr; _TermExpression zero-or-many-of { one-of ComparisonOp TermExpression }_

**TermExpression** &rarr; _FactorizedExpression zero-or-many-of { one-of TermOp FactorizedExpression }_

**FactorizedExpression** &rarr; _UnaryExpression zero-or-many-of { one-of FactorOp UnaryExpression }_

**UnaryExpression** &rarr; _one-of { one-of UniOp UnaryExpression | PrimaryExpression }_

**PrimaryExpression** &rarr; _one-of { Literal | '(' Expression ')' }_
