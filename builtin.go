package tgtl

func p(env *Environment, args ...Value) (Value, Effect) {
	for _, arg := range args {
		print(arg, " ")
	}
	print("\n")
	return nil, nil
}

func print_(env *Environment, args ...Value) (Value, Effect) {
	var msg string
	erra := Args(args, &msg)
	if erra != nil {
		return env.FailString("printf: ${1}", erra)
	}
	extra := []Value{}
	if len(args) > 1 {
		extra = args[1:len(args)]
	}
	n, err := env.Printi(msg, extra...)
	if err == nil {
		return Int(n), nil
	}
	return Int(n), ErrorFromError(err)
}

func write(env *Environment, args ...Value) (Value, Effect) {
	var msg string
	erra := Args(args, &msg)
	if erra != nil {
		return env.FailString("write: ${1}", erra)
	}
	n, err := env.Write(msg)
	if err == nil {
		return Int(n), nil
	}
	return Int(n), ErrorFromError(err)
}

func iadd(env *Environment, args ...Value) (Value, Effect) {
	var i, j int
	err := Args(args, &i, &j)
	if err != nil {
		return env.Fail(err)
	}
	return Int(i + j), nil
}

func isub(env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 int
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	return Int(v1 - v2), nil
}

func imul(env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 int
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	return Int(v1 * v2), nil
}

func idiv(env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 int
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	if v2 == 0 {
		return nil, ErrorFromString("division by 0")
	}
	return Int(v1 / v2), nil
}

func igt(env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 int
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	return Bool(v1 > v2), nil
}

func ilt(env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 int
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	return Bool(v1 < v2), nil
}

func ige(env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 int
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	return Bool(v1 >= v2), nil
}

func ile(env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 int
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	return Bool(v1 <= v2), nil
}

func ieq(env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 int
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	return Bool(v1 == v2), nil
}

func seq(env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 string
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	return Bool(v1 == v2), nil
}

func teq(env *Environment, args ...Value) (Value, Effect) {
	var t1, t2 Type
	err := Args(args, &t1, &t2)
	if err != nil {
		return env.Fail(err)
	}
	return Bool(t1 == t2), nil
}

func updateIntByName(update func(in Int) Int, env *Environment, args ...Value) (Int, Effect) {
	var name Word
	err := Args(args, &name)
	if err != nil {
		return Int(0), err
	}
	val := env.Lookup(name.String())
	vi, ok := val.(Int)
	if !ok {
		return Int(0), ErrorFromString("Not an integer.")
	}
	newi := update(vi)
	env.Set(name.String(), newi)
	return newi, nil
}

func inc(env *Environment, args ...Value) (Value, Effect) {
	return updateIntByName(func(in Int) Int {
		return in + 1
	}, env, args...)
}

func dec(env *Environment, args ...Value) (Value, Effect) {
	return updateIntByName(func(in Int) Int {
		return in - 1
	}, env, args...)
}

func str(env *Environment, args ...Value) (Value, Effect) {
	var v1 Value
	err := Args(args, &v1)
	if err != nil {
		return env.Fail(err)
	}
	return String(v1.String()), nil
}

func int_(env *Environment, args ...Value) (Value, Effect) {
	var v1 Value
	err := Args(args, &v1)
	if err != nil {
		return env.Fail(err)
	}
	rs := []rune(v1.String() + " ")
	index := 0
	return ParseInteger(rs, &index)
}

func boolBinop(op func(b1, b2 bool) bool, env *Environment, args ...Value) (Value, Effect) {
	var v1, v2 bool
	err := Args(args, &v1, &v2)
	if err != nil {
		return env.Fail(err)
	}
	return Bool(op(v1, v2)), nil
}

func isnil(env *Environment, args ...Value) (Value, Effect) {
	if len(args) < 1 {
		return env.FailString("isnil requires 1 argument")
	}
	return Bool(args[0] == nil), nil
}

func band(env *Environment, args ...Value) (Value, Effect) {
	return boolBinop(func(b1, b2 bool) bool {
		return b1 && b2
	}, env, args...)
}

func bor(env *Environment, args ...Value) (Value, Effect) {
	return boolBinop(func(b1, b2 bool) bool {
		return b1 || b2
	}, env, args...)
}

func bxor(env *Environment, args ...Value) (Value, Effect) {
	return boolBinop(func(b1, b2 bool) bool {
		return b1 != b2
	}, env, args...)
}

func bnot(env *Environment, args ...Value) (Value, Effect) {
	var v1 bool
	err := Args(args, &v1)
	if err != nil {
		return env.Fail(err)
	}
	return Bool(!v1), nil
}

func val(env *Environment, args ...Value) (Value, Effect) {
	if len(args) < 1 {
		return env.FailString("val requres at least one argument.")
	}
	return List(args), nil
}

func ret(env *Environment, args ...Value) (Value, Effect) {
	if len(args) < 1 {
		return env.Return(nil)
	} else if len(args) == 1 {
		return env.Return(args[0])
	} else {
		return env.Return(List(args))
	}
}

func fail(env *Environment, args ...Value) (Value, Effect) {
	if len(args) < 1 {
		return env.Fail(ErrorFromString("fail"))
	} else {
		return env.FailString(args[0].String(), args[1:len(args)]...)
	}
}

func break_(env *Environment, args ...Value) (Value, Effect) {
	if len(args) < 1 {
		return env.Break(nil)
	} else if len(args) == 1 {
		return env.Break(args[0])
	} else {
		return env.Break(List(args))
	}
}

func nop(env *Environment, args ...Value) (Value, Effect) {
	return nil, nil
}

func typeof_(env *Environment, args ...Value) (Value, Effect) {
	var val Value
	err := Args(args, &val)
	if err != nil {
		return nil, err
	}
	return TypeOf(val), nil
}

func type_(env *Environment, args ...Value) (Value, Effect) {
	var val Value
	err := Args(args, &val)
	if err != nil {
		return nil, err
	}
	name := val.String()
	return Type(name), nil
}

func to(env *Environment, args ...Value) (Value, Effect) {
	var name string

	if len(args) < 2 {
		return env.FailString("to needs at least 2 arguments")
	}
	err := Convert(args[0], &name)
	if err != nil {
		return env.Fail(err)
	}
	block, ok := (args[len(args)-1]).(Block)
	if !ok {
		return env.FailString("to: last argument must be a block")
	}

	last := args[len(args)-1]
	block, isBlock := last.(Block)
	if !isBlock {
		return env.FailString("Not a block")
	}
	params := args[1 : len(args)-1]
	defined := Defined{name, params, block}
	env.Define(name, defined, 1)
	return defined, nil
}

func do(env *Environment, args ...Value) (Value, Effect) {
	var name string
	var doArgs List
	err := Args(args, &name, &doArgs)
	if err != nil {
		return env.Fail(err)
	}
	fun := env.Lookup(name)
	if fun == nil {
		return env.FailString("Cannot evaluate unknown order: " + name)
	}
	eva, ok := fun.(Evaler)
	if !ok {
		return env.FailString("Cannot evaluate: " + name)
	}
	return eva.Eval(env, doArgs...)
}

func if_(env *Environment, args ...Value) (Value, Effect) {
	var cond, ok, haveElse bool
	var ifBlock, elseBlock Block

	if len(args) < 2 {
		return env.FailString("if needs at least 2 arguments")
	}
	if len(args) > 4 {
		return env.FailString("if needs at most 4 arguments")
	}
	err := Convert(args[0], &cond)
	if err != nil {
		return env.Fail(err)
	}
	ifBlock, ok = (args[1]).(Block)
	if !ok {
		return env.FailString("if: second argument must be a block")
	}
	elseIndex := 2
	if 2 < len(args) {
		// look for an else keyword but don't mind if it really is else
		_, ok = (args[2]).(Word)
		if ok {
			// block after else keyword
			elseIndex = 3
		}
	}
	if elseIndex < len(args) {
		// There should be an else block...
		elseBlock, ok = (args[elseIndex]).(Block)
		if !ok {
			return env.FailString("if: missing else block")
		}
		haveElse = true
	}
	if cond {
		return ifBlock.Eval(env, args...)
	} else {
		if haveElse {
			return elseBlock.Eval(env, args...)
		} else {
			return nil, nil
		}
	}
}

func switch_(env *Environment, args ...Value) (Value, Effect) {
	var defaultBlock Block
	var haveDefault bool = false
	if len(args) < 3 {
		return env.FailString("switch needs at least 3 arguments")
	}
	compareTo := args[0]
	for i := 2; i < len(args); i += 2 {
		case_ := args[i-1]
		block, blockOk := args[i].(Block)
		if !blockOk {
			return env.FailString("switch: argument ${1} is not a block",
				Int(i))
		}
		if kw, kwOk := case_.(Word); kwOk && kw.String() == "default" {
			if haveDefault {
				return env.FailString("switch: duplicate default block ${1}",
					Int(i))
			}
			haveDefault = true
			defaultBlock = block
		} else {
			if compareTo.String() == case_.String() {
				return block.Eval(env, args...)
			}
		}
	}
	if haveDefault {
		return defaultBlock.Eval(env, args...)
	}
	return nil, nil
}

func while(env *Environment, args ...Value) (Value, Effect) {
	var blockRes Value
	var blockEff Effect
	if len(args) != 2 {
		return env.FailString("while needs exactly 3 arguments")
	}
	cond, condOk := args[0].(Block)
	block, blockOk := args[1].(Block)
	if !condOk {
		return env.FailString("while condition must be a block")
	}
	if !blockOk {
		return env.FailString("while body must be a block")
	}

	for res, eff := cond.Eval(env, args...); ValToBool(res); res, eff = cond.Eval(env, args...) {
		if eff != nil && eff.Flow() > NormalFlow {
			return res, eff
		}
		blockRes, blockEff = block.Eval(env, args...)
		if blockEff != nil && blockEff.Flow() > NormalFlow {
			return blockRes, blockEff
		}
	}
	return blockRes, blockEff
}

func rescue(env *Environment, args ...Value) (Value, Effect) {
	var block Block
	err := Args(args, &block)
	if err != nil {
		return env.Fail(err)
	}
	return env.Prevent(block)
}

func set(env *Environment, args ...Value) (Value, Effect) {
	if len(args) < 2 {
		return env.FailString("set needs at 2 arguments")
	}
	if args[0] == nil {
		return env.FailString("set $1 is nil")
	}
	return env.Set(args[0].String(), args[1])
}

func let(env *Environment, args ...Value) (Value, Effect) {
	if len(args) < 2 {
		return env.FailString("def needs at 2 arguments")
	}
	if args[0] == nil {
		return env.FailString("def $1 is nil")
	}
	return env.Define(args[0].String(), args[1], 1)
}

func get(env *Environment, val ...Value) (Value, Effect) {
	if len(val) < 1 {
		return env.FailString("get needs at least 1 argument")
	}
	target := val[0].String()
	return env.Lookup(target), nil
}

func list(env *Environment, args ...Value) (Value, Effect) {
	return List(args), nil
}

func sadd(env *Environment, args ...Value) (Value, Effect) {
	var value Value
	var str String
	err := Args(args, &str, &value)
	if err != nil {
		return env.Fail(err)
	}
	str = str + String(value.String())
	return str, nil
}

func sget(env *Environment, args ...Value) (Value, Effect) {
	var index int
	var str String
	err := Args(args, &str, &index)
	if err != nil {
		return env.Fail(err)
	}
	runes := []rune(str)
	if (index < 0) || (index >= len(runes)) {
		return env.FailString("index out of range")
	}
	return Int(runes[index]), nil
}

func runes(env *Environment, args ...Value) (Value, Effect) {
	var str String
	err := Args(args, &str)
	if err != nil {
		return env.Fail(err)
	}
	res := List{}
	runes := []rune(str)
	for i := 0; i < len(runes); i++ {
		res = append(res, Int(runes[i]))
	}
	return res, nil
}

func wire(env *Environment, args ...Value) (Value, Effect) {
	var str String
	for i := 0; i < len(args); i++ {
		var ch Int
		err := Convert(args[i], &ch)
		if err != nil {
			return str, err
		}
		str = str + String([]rune{rune(ch)})
	}
	return str, nil
}

func slen(env *Environment, args ...Value) (Value, Effect) {
	var str String
	err := Args(args, &str)
	if err != nil {
		return env.Fail(err)
	}
	runes := []rune(str)
	return Int(len(runes)), nil
}

func ladd(env *Environment, args ...Value) (Value, Effect) {
	var value Value
	var list List
	err := Args(args, &list, &value)
	if err != nil {
		return env.Fail(err)
	}
	list = append(list, value)
	return list, nil
}

func lget(env *Environment, args ...Value) (Value, Effect) {
	var index int
	var list List
	err := Args(args, &list, &index)
	if err != nil {
		return env.Fail(err)
	}
	if (index < 0) || (index >= len(list)) {
		return env.FailString("index out of range")
	}
	return list[index], nil
}

func lset(env *Environment, args ...Value) (Value, Effect) {
	var index int
	var list List
	var val Value
	err := Args(args, &list, &index, &val)
	if err != nil {
		return env.Fail(err)
	}
	if (index < 0) || (index >= len(list)) {
		return env.FailString("index out of range")
	}
	list[index] = val
	return list[index], nil
}

func llen(env *Environment, args ...Value) (Value, Effect) {
	var list List
	err := Args(args, &list)
	if err != nil {
		return env.Fail(err)
	}
	return Int(len(list)), nil
}

func lsort(env *Environment, args ...Value) (Value, Effect) {
	var list List
	err := Args(args, &list)
	if err != nil {
		return env.Fail(err)
	}
	return list.SortStrings(), nil
}

func leach(env *Environment, args ...Value) (Value, Effect) {
	var list List
	var key Word
	var name Word
	var block Block
	err := Args(args, &list, &key, &name, &block)
	if err != nil {
		return env.Fail(err)
	}
	for i, v := range list {
		env.Define(key.String(), Int(i), 0)
		env.Define(name.String(), v, 0)
		bval, berr := block.Eval(env, args...)
		if berr != nil {
			return bval, berr
		}
	}
	return list, nil
}

func lslice(env *Environment, args ...Value) (Value, Effect) {
	var list List
	var from Int
	var to Int
	err := Args(args, &list, &from, &to)
	if err != nil {
		return env.Fail(err)
	}
	length := Int(len(list))
	if length == 0 {
		return list, nil
	}
	if from < 0 {
		from = length - from
	}
	if to < 0 {
		from = length - from
	}
	if from >= length {
		from = length - 1
	}
	if to >= length {
		to = length - 1
	}
	if from > to {
		from, to = to, from
	}
	return list[from:to], nil
}

func map_(env *Environment, args ...Value) (Value, Effect) {
	res := make(Map)
	for i := 1; i < len(args); i += 2 {
		key := args[i-1]
		val := args[i]
		res[key.String()] = val
	}
	return res, nil
}

func mget(env *Environment, args ...Value) (Value, Effect) {
	var index string
	var hmap Map
	err := Args(args, &hmap, &index)
	if err != nil {
		return env.Fail(err)
	}
	return hmap[index], nil
}

func mset(env *Environment, args ...Value) (Value, Effect) {
	var index string
	var hmap Map
	var val Value
	err := Args(args, &hmap, &index, &val)
	if err != nil {
		return env.Fail(err)
	}
	hmap[index] = val
	return hmap[index], nil
}

func mkeys(env *Environment, args ...Value) (Value, Effect) {
	var hmap Map
	err := Args(args, &hmap)
	if err != nil {
		return env.Fail(err)
	}
	res := List{}
	for k, _ := range hmap {
		res = append(res, String(k))
	}
	return res, nil
}

func meach(env *Environment, args ...Value) (Value, Effect) {
	var map_ Map
	var key Word
	var name Word
	var block Block
	err := Args(args, &map_, &key, &name, &block)
	if err != nil {
		return env.Fail(err)
	}
	miter := Map{}

	for k, v := range map_ {
		miter[k] = v
	}
	for k, v := range miter {
		env.Define(key.String(), String(k), 0)
		env.Define(name.String(), v, 0)
		bval, beff := block.Eval(env, args...)
		if beff != nil {
			return bval, beff
		}
	}
	return map_, nil
}

func expand(env *Environment, args ...Value) (Value, Effect) {
	var msg string
	err := Args(args, &msg)
	if err != nil {
		return env.Fail(err)
	}
	res := env.Interpolate(msg)
	return String(res), nil
}

func help(env *Environment, args ...Value) (Value, Effect) {
	var name string
	err := Args(args, &name)
	if err != nil {
		return env.Fail(err)
	}
	helpMap := env.Lookup("HELP")
	if helpMap == nil {
		env.Printi("help: $1:No help available 1.\n", String(name))
		return nil, err
	}
	if name == "all" {
		keys := helpMap.(Map).SortedKeys()
		for _, k := range keys {
			v := helpMap.(Map)[k.String()]
			env.Printi("$1: $2\n", k, v)
		}
		return nil, nil
	}

	msg, ok := helpMap.(Map)[name]
	if ok {
		env.Printi("help: $1:\n$2\n", String(name), msg)

	} else {
		env.Printi("help: $1:No help available 2.\n", String(name))
	}
	return msg, nil
}

func explain(env *Environment, args ...Value) (Value, Effect) {
	var name string
	var help String
	err := Args(args, &name, &help)
	if err != nil {
		return env.Fail(err)
	}
	helpMap := env.Lookup("HELP")
	if helpMap == nil {
		helpMap = make(Map)
	}
	helpMap.(Map)[name] = help
	env.Define("HELP", helpMap, -1)
	return help, nil
}

func overload(env *Environment, args ...Value) (Value, Effect) {
	var name string
	var target Value
	err := Args(args, &name, &target)
	if err != nil {
		return env.Fail(err)
	}
	if len(args) < 3 {
		return env.FailString("overload needs at least 3 arguments")
	}
	return env.Overload(name, target, args[2:len(args)])
}

func (env *Environment) Register(name string,
	f func(e *Environment, args ...Value) (Value, Effect), help string) {
	env.Define(name, Proc(f), -1)
	explain(env, String(name), String(help))
}

func (env *Environment) RegisterBuiltins() {
	env.Define("true", Bool(true), -1)
	env.Define("false", Bool(false), -1)
	env.Register("sadd", sadd, "returns a string  with $2 appended to string $1")
	env.Register("sget", sget, "gets a rune from a string by index")
	env.Register("slen", slen, "returns the length of a string")
	env.Register("iadd", iadd, "adds two integers together")
	env.Register("band", band, `returns true if $1 and $2 arguments are true`)
	env.Register("bor", bor, `returns true if $1 or $2 arguments are true`)
	env.Register("bxor", bxor, `returns true if $1 and $2 are different booleans`)
	env.Register("bnot", bnot, `returns true if $1 is false and false otherwise`)
	env.Register("ladd", ladd, "returns a list with $2 appended to List $1")
	env.Register("list", list, "creates a new array list")
	env.Register("lget", lget, "gets a value from a list by index")
	env.Register("lset", lset, "sets a value to a list by index and value")
	env.Register("llen", llen, "returns the length of a list")
	env.Register("lsort", lsort, "returns the List $1 sorted by string value")
	env.Register("leach", leach, "calls the block $4 for each entry in the list")
	env.Register("lslice", lslice, "slices the list $1 from $2 to $3")
	env.Register("iadd", iadd, "adds and Ints to and Int")
	env.Register("isub", isub, "subtracts an Int from an Int")
	env.Register("imul", imul, "multiplies an Ints by an Int")
	env.Register("idiv", idiv, "divides an Int by an Int")
	env.Register("ilt", ilt, "checks if $1 < $2, where $1 and $2 must be Int")
	env.Register("ile", ile, "checks if $1 <= $2, where $1 and $2 must be Int")
	env.Register("igt", igt, "checks if $1 > $2, where $1 and $2 must be Int")
	env.Register("ige", ige, "checks if $1 >= $2, where $1 and $2 must be Int")
	env.Register("ieq", ieq, "checks if $1 == $2, where $1 and $2 must be Int")
	env.Register("seq", seq, "checks if [str $1] == [str $2]")
	env.Register("str", str, "converts $1 to String")
	env.Register("wire", wire, "converts unicode character indexes or runes to String")
	env.Register("runes", runes, "converts String to alist of character indexes or runes")
	env.Register("int", int_, "converts $1 to Int")
	env.Register("inc", inc, "increments the named integer $1")
	env.Register("dec", dec, "decrements the named integer $1")
	env.Register("map", map_, "creates a new hash map")
	env.Register("mget", mget, "gets a value from a map by key")
	env.Register("mset", mset, "sets a value to a map by key and value")
	env.Register("mkeys", mkeys, "returns all keys of a map as an unsorted list")
	env.Register("meach", meach, "calls the block $4 for each entry in the map")

	env.Register("p", p, "print debug output")
	env.Register("print", print_, "print to the environnment's current writer with interpolation")
	env.Register("write", write, "write to the environnment's current writer")
	env.Register("to", to, "define a procedure")
	env.Register("do", do, "execute a command $1 with arguments in $2 as array")
	env.Register("ret", ret, "return from a procedure")
	env.Register("return", ret, "return from a procedure")
	env.Register("break", break_, "return from a block")
	env.Register("val", val, "gets the value of a value")
	env.Register("let", let, "creates a new vavariable with given value")
	env.Register("set", set, "sets an existing variable")
	env.Register("get", get, "get the contents of a variable")
	env.Register("help", help, "get help for a procedure")
	env.Register("explain", explain, "set the help for a procedure")
	env.Register("expand", expand, "interpolate strings from environment")
	env.Register("fail", fail, "fail execution of a procedure")
	env.Register("rescue", rescue, "call $1 as the error handler on failure")
	env.Register("if", if_, "if runs $1 if $0 is true, otherwise runs $2")
	env.Register("isnil", isnil, "returns true if $1 is nil, false if not")
	env.Register("switch", switch_, "selects one of many cases")
	env.Register("type", type_, "returns $1 converted to a type")
	env.Register("teq", teq, "checks if $1 and $2 are exactly the same type")
	env.Register("typeof", typeof_, "returns the type of $1 or Unknown if not known")
	env.Register("nop", nop, "does nothing and returns nil")
	env.Register("overload", overload, "creates a command overload named $1 targeting $2 for the types following $2")
}

// This function registers builtins that make Tgtl turing complete
// Not to be used in situations where this is undesirable.
func (env *Environment) RegisterTuringCompleteBuiltins() {
	env.Register("while", while, "executes $2 while $1 returns true")
}
