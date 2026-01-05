## Hamlet's Syntactic Grammer

> The following grammer is represented in a flavour of [BNF](https://en.wikipedia.org/wiki/Backus-Naur_form), having some special _specifier_ keywords.

**Literal** &rarr; _one-of { string | character | NumericalLiteral | BooleanLiteral | nil }_

**BooleanLiteral** &rarr; _one-of { true | false }_

**NumericalLiteral** &rarr; _one-of { real | integer }_

**Expression** &rarr; _one-of { GroupedExpression | UnaryExpression | BinaryExpression | Literal }_

**Operator** &rarr; _one-of { '==' | '!=' | '<' | '<=' | '>' | '>=' | '+'  | '-'  | '*' | '/' }_

**GroupedExpression** &rarr; '(' Expression ')'

**UnaryExpression** &rarr; _one-of { '~' | '!' } Expression_

**BinaryExpression** &rarr; _Expression Operator Expression_
