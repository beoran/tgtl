# tgtl

A Tiny Go Tool Language.

Tgtl is a TCL-like embedded scripting language implemented in Go.

## Introduction

Tgtl is A Tiny Go Tool Language, an embeddable, interpreted scripting language,
implemented in Go, that somewhat resembles TCL and shell scripts, but with
typed values, and LISP like semantics.

One notable feature is that the language itself has no key words, nor
flow control by itself, but that this is provided by the commands that
TGTL executes. This makes it possible to disable certain commands,
or example, in case where it is desirable for the script to be not Turing
complete.

The syntax is extremelty simple, and based on a rescursive descent LL1 parser,
where the parser only considers the next character and the current state
of parsing to deterine the meaning of the code. Code that is easy to parse
by the computer is easy to understand by humans also, which is why the
limitations of LL1 parsing are acceptable.

## Grammar
The formal grammar of TGTL is as follows:

	SCRIPT        -> STATEMENTS .
	STATEMENTS    -> STATEMENT OPTSTATEMENTS .
	OPTSTATEMENTS -> rs STATEMENT OPTSTATEMENTS | .
	STATEMENT     -> OPTWS EXPRESSION .
	EXPRESSION    -> COMMAND | BLOCK | comment | .
	OPTWS         -> ws | .
	COMMAND       -> ORDER PARAMETERS .
	ORDER         -> LITERAL | EVALUATION .
	BLOCK         -> ob STATEMENTS cb .
	PARAMETERS    -> ws PARAMETER OPTPARAMETERS | .
	PARAMETER     -> LITERAL | BLOCK | GETTER | EVALUATION | .
	EVALUATION    -> oe COMMAND ce .
	GETTER        -> get TARGET .
	TARGET 		  -> GETTER | LITERAL .
	LITERAL       -> word | string | integer .
	rs			-> /[\n\r]+/ .
	ws			-> /[\t ]+/  .
	word 		-> /[^ \t\n\$\(\)\{\}\]\[]+/
	string 		-> /"[^"]+"/ | /`[^`]+`/
	integer     -> [+-]?[0-9]+
	comment 	-> /#[^\n]+\n/ .
	get			-> '$' .
	oe 			-> '[' .
	ce          -> ']' .
	ob 			-> '{' .
	cb 			-> '}' .

