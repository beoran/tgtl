package tgtl

type Tgtl struct {
	index int
	input string
}

type Proc func(*Environment, ...Value) Value

func (pv Proc) String() string {
	return "proc"
}

// Evaler is an interface to a Value that can evaluate itself
// based on and possibly modifying a given environment.
type Evaler interface {
	Eval(*Environment, ...Value) Value
}

// Lazyer is an interface to a Value that does not automatically
// evaluate itself in List and similar contexts.
type Lazyer interface {
	// Lazy is a marker method only.
	Lazy()
}

// Value is an interface for the basic unit of data that
// the TGTL interpreter works with.
// It is used both for run time data and parsing data.
type Value interface {
	// Since TGTL is TCL-like, all values must be
	// convertable to string
	String() string
	// Furthermore any TGTL value must be able to evaluate itself.
	Evaler
}

// Typer is an interface that Values can optionally implement to
// allow them to report their type.
type Typer interface {
	Type() Type
}

type Int int

type Bool bool

type String string

type Word string

type Type string

type Comment string

type Error struct {
	Message  string
	Index    int
	Children List
}

type List []Value

type Map map[string]Value

func (m Map) Keys() List {
	res := List{}
	for k, _ := range m {
		res = append(res, String(k))
	}
	return res
}

func (m Map) SortedKeys() List {
	return m.Keys().SortStrings()
}

type Getter struct {
	Key Value
}

type Evaluation struct {
	Command
}

type Command struct {
	Order      Value
	Parameters List
}

type Block struct {
	Statements List
}

type Defined struct {
	Name   string
	Params List
	Block
}

type Wrapper struct {
	Kind    Type
	Handle  interface{}
	Methods Map
}

type Object struct {
	Wrapper
	Fields   Map
	Embedded Map
}

type Overload Map

func (bv Bool) String() string {
	if bv {
		return "true"
	}
	return "false"
}

func (sv String) String() string {
	return string(sv)
}

func (cv Comment) String() string {
	return string(cv)
}

func (sv Word) String() string {
	return string(sv)
}

func (tv Type) String() string {
	return string(tv)
}

func (gv Getter) String() string {
	return gv.Key.String()
}

func (cv Command) String() string {
	return cv.Order.String() + " " + cv.Parameters.String()
}

func (gv Evaluation) String() string {
	return gv.Command.String()
}

func (bv Block) String() string {
	return "{" + bv.Statements.String() + "}"
}

func (dv Defined) String() string {
	return "to " + dv.Name + " (" +
		dv.Params.String() + ") " + dv.Block.String()
}

func (cv Overload) String() string {
	return Map(cv).String()
}

// Lazy marks that Blocks should be lazily evaluated.
func (Block) Lazy() {}

func (iv Int) String() string {
	return Itoa(int(iv))
}

func (lv Map) String() string {
	aid := "[map "
	for k, v := range lv {
		aid += " "
		aid += k
		aid += " "
		aid += v.String()
	}
	aid += "]"
	return aid
}

func (lv List) String() string {
	aid := "[list"
	for _, v := range lv {
		aid += " "
		if v == nil {
			aid += "nil"
		} else {
			aid += v.String()
		}
	}
	aid += "]"
	return aid
}

func (ev Error) String() string {
	return string(ev.Message)
}

func (ev Error) Error() string {
	return ev.Message
}

// Implement the effect interface
func (ev *Error) Flow() Flow {
	if ev == nil {
		return NormalFlow
	}
	return FailFlow
}

func (ev *Error) Unwrap() Value {
	if ev == nil {
		return nil
	}
	return ev
}

func NewError(message string, index int, children ...Value) *Error {
	return &Error{message, index, children}
}

func ErrorFromString(message string) *Error {
	return NewError(message, -1)
}

func ErrorFromError(err error, children ...Value) *Error {
	if err == nil {
		return nil
	}
	return NewError(err.Error(), -1, children...)
}

// Break is used for break flows
type Break struct {
	Value // value returned by break
}

func (bv Break) Flow() Flow {
	return BreakFlow
}

func (bv Break) Unwrap() Value {
	return bv.Value
}

// Return is used for return flows
type Return struct {
	Value // value returned
}

func (rv Return) Flow() Flow {
	return ReturnFlow
}

func (bv Return) Unwrap() Value {
	return bv.Value
}

// Rescue is used to evaluate rescue commands
type Rescue struct {
	Block // A rescue is a special block
}

// Rescue essentially protects the block from
// rescue recursion by pushing the stack once.
func (r Rescue) Eval(env *Environment, args ...Value) Value {
	// ignore the stack depth
	// protection here to be sure the
	// rescue block is executed
	_ = env.Push()
	env.Rescuing = true
	env.Printi("Rescuing.\n")
	defer env.Pop()
	defer func() {
		env.Rescuing = false
	}()
	return r.Block.Eval(env, args...)
}

func (iv Wrapper) String() string {
	aid := "[interface"
	aid += " " + iv.Kind.String()
	aid += " " + iv.Methods.String()
	aid += "]"
	return aid
}

func (sv Object) String() string {
	aid := "[struct"
	aid += " " + sv.Kind.String()
	aid += " " + sv.Methods.String()
	aid += " " + sv.Fields.String()
	aid += " " + sv.Embedded.String()
	aid += "]"
	return aid
}

func NewTgtl(input string) Tgtl {
	return Tgtl{0, input}
}

func (wv Word) Eval(env *Environment, args ...Value) Value {
	return wv
}

func (sv String) Eval(env *Environment, args ...Value) Value {
	return sv
}

func (tv Type) Eval(env *Environment, args ...Value) Value {
	return tv
}

func (iv Int) Eval(env *Environment, args ...Value) Value {
	return iv
}

func (bv Bool) Eval(env *Environment, args ...Value) Value {
	return bv
}

func (ev Error) Eval(env *Environment, args ...Value) Value {
	return ev
}

// Eval of a List expands arguments, except Lazyer elements.
func (lv List) Eval(env *Environment, args ...Value) Value {
	res := List{}
	for _, s := range lv {
		_, isLazy := s.(Lazyer)
		if isLazy {
			res = append(res, s)
		} else {
			val := s.Eval(env, args...)
			flow := ValueFlow(val)
			if flow > NormalFlow {
				return val
			}
			res = append(res, val)
		}
	}
	return res
}

func (mv Map) Eval(env *Environment, args ...Value) Value {
	return mv
}

func (cv Comment) Eval(env *Environment, args ...Value) Value {
	return nil
}

func (bv Block) Eval(env *Environment, args ...Value) Value {
	var res Value
	// set parameters to $1 ... $(len(args))
	for i, a := range args {
		name := Itoa(i + 1)
		env.Define(name, a, 0)
	}
	// Set $argc to amount of arguments
	// and $argv to arguments as well
	env.Define("argc", Int(len(args)), 0)
	env.Define("argv", List(args), 0)
	for _, s := range bv.Statements {
		// Call the statement.
		res = s.Eval(env, args...)
		// if the flow is not normal anymore,
		// end the block execution at this point
		flow := ValueFlow(res)
		if flow > NormalFlow {
			// If it was a break or less, unwrap and done,
			if flow <= BreakFlow {
				return res.(Effect).Unwrap()
			} else if flow == FailFlow {
				// If it is a fail try to rescue it
				return env.Rescue(res)
			}
			return res
		}
		env.Define("RESULT", res, 0)
	}
	return res
}

func (pv Proc) Eval(env *Environment, args ...Value) Value {
	return pv(env, args...)
}

func (cv Command) Eval(env *Environment, args ...Value) Value {
	val := cv.Order.Eval(env)
	flow := ValueFlow(val)
	if flow > NormalFlow {
		return val
	}
	name := val.String()
	fun := env.Lookup(name)
	if fun == nil {
		return ErrorFromString("Cannot evaluate nil order: " + name)
	}
	eva, ok := fun.(Evaler)
	if !ok {
		return ErrorFromString("Cannot evaluate: " + name)
	}
	err := env.Push()
	// stack depth protection
	if err != nil {
		return env.Rescue(env.Fail(err))
	}
	defer env.Pop()
	fargs := cv.Parameters
	// Expand Evaluation arguments, but not block elements.
	eargs := fargs.Eval(env, args...)
	flow = ValueFlow(eargs)
	if flow > NormalFlow {
		return eargs
	}
	return eva.Eval(env, eargs.(List)...)
}

func (gv Getter) Eval(env *Environment, args ...Value) Value {
	val := gv.Key.Eval(env)
	flow := ValueFlow(val)
	if flow > NormalFlow {
		return val
	}
	return env.Lookup(val.String())
}

func (ev Evaluation) Eval(env *Environment, args ...Value) Value {
	err := env.Push()
	// stack depth protection
	if err != nil {
		return env.Fail(err)
	}
	defer env.Pop()
	val := ev.Command.Eval(env, args...)
	return val
}

func (dv Defined) Eval(env *Environment, args ...Value) Value {
	err := env.Push()
	// stack depth protection
	if err != nil {
		return env.Rescue(env.Fail(err))
	}
	if len(dv.Params) > len(args) {
		return env.FailString("Not enough arguments")
	}
	for i := 0; i < len(dv.Params); i++ {
		env.Define(dv.Params[i].String(), args[i], 0)
	}
	// $0 contains the name of the defined procedure
	env.Define("0", String(dv.Name), 0)
	defer env.Pop()
	val := dv.Block.Eval(env, args...)
	flow := ValueFlow(val)
	if flow < ReturnFlow {
		return val
	} else if flow == ReturnFlow {
		return val.(Effect).Unwrap()
	} else { // failures pass through
		return val
	}
}

func (iv Wrapper) Eval(env *Environment, args ...Value) Value {
	// Object like values such as interfaces or structs
	// have methods that are called with the method
	// name as the first word argument, which is used
	// to dispatch the function.
	// The dispathed function receives the object
	// as it's first argument
	var name Word
	err := Args(args, &name)
	if err != nil {
		return env.Fail(err)
	}
	method, ok := iv.Methods[name.String()]
	if !ok {
		return env.FailString("No such method ${1}", name)
	}
	args[0] = iv
	return method.Eval(env, args...)
}

func (sv Object) Eval(env *Environment, args ...Value) Value {
	// See interface for this dispatch
	var name Word
	err := Args(args, &name)
	if err != nil {
		return env.Fail(err)
	}
	method, ok := sv.Methods[name.String()]
	if !ok {
		return env.FailString("No such method ${1}", name)
	}
	args[0] = sv
	return method.Eval(env, args...)
}

func TypeOf(val Value) Type {
	if typer, ok := val.(Typer); ok {
		return typer.Type()
	} else {
		return Type("Unknown")
	}
}

func (cv Overload) Eval(env *Environment, args ...Value) Value {
	signature := ""
	for _, arg := range args {
		signature += "_" + TypeOf(arg).String()
	}
	target, ok := cv[signature]
	if ok {
		return target.Eval(env, args...)
	}
	return env.FailString("No overload defined for signature: " + signature)
}

// Implement Typer interface for commonly used Values
func (String) Type() Type     { return Type("String") }
func (Bool) Type() Type       { return Type("Bool") }
func (Int) Type() Type        { return Type("Int") }
func (Error) Type() Type      { return Type("Error") }
func (List) Type() Type       { return Type("List") }
func (Map) Type() Type        { return Type("Map") }
func (Proc) Type() Type       { return Type("Proc") }
func (Word) Type() Type       { return Type("Word") }
func (Defined) Type() Type    { return Type("Defined") }
func (Block) Type() Type      { return Type("Block") }
func (Command) Type() Type    { return Type("Command") }
func (Getter) Type() Type     { return Type("Getter") }
func (Evaluation) Type() Type { return Type("Evaluation") }
func (t Type) Type() Type     { return Type("Type") }
func (s Object) Type() Type   { return s.Kind }
func (i Wrapper) Type() Type  { return i.Kind }
func (t Type) Overload() Type { return Type("Overload") }

type EffectValue interface {
	Effect
	Value
}

func ValueFlow(v Value) Flow {
	effect, ok := v.(Effect)
	if ok {
		return effect.Flow()
	}
	return NormalFlow
}
