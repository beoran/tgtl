package tgtl

// Maximum amount of frames,
// to prevent unlimited recursion.
const FRAMES_MAX = 80

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type Frame struct {
	Variables Map
	Out       Writer
	In        Reader
	Rescuer   Value
}

type Environment struct {
	Frames   []*Frame
	Out      Writer
	In       Reader
	Rescuing bool
}

// Looks up the value of a variable and the frame it is in
func (env Environment) LookupFrame(name string) (Value, *Frame) {
	for i := len(env.Frames) - 1; i >= 0; i-- {
		frame := env.Frames[i]
		val, ok := frame.Variables[name]
		if ok {
			return val, frame
		}
	}
	return nil, nil
}

func (env Environment) Lookup(name string) Value {
	val, _ := env.LookupFrame(name)
	return val
}

func (env *Environment) Push() *Error {
	env.Frames = append(env.Frames, &Frame{make(Map), env.Out, env.In, nil})
	if len(env.Frames) >= FRAMES_MAX && !env.Rescuing {
		return ErrorFromString("PROGRAM HAS DISAPPEARED INTO THE BLACK LAGOON - too much recursion or function calls")
	}
	return nil
}

func (env *Environment) Pop() {
	l := len(env.Frames)
	if l > 0 {
		env.Frames = env.Frames[0 : l-1]
	}
}

// Depth returns the amount of frames on the frame stack
func (env *Environment) Depth() int {
	return len(env.Frames)
}

// Frame returns a frame pointer based on the level.
// 0 is the top-level index. A negative level wil refer to the
// outermost frame.
// Returns nil if the level is somehow out of range or
// if no frames have been pushed yet
func (env *Environment) Frame(level int) *Frame {
	if len(env.Frames) < 1 {
		return nil
	}
	if level > 0 {
		l := len(env.Frames)
		index := l - level - 1
		if index < 0 || index >= l {
			return nil
		}
		return env.Frames[index]
	} else {
		return env.Frames[0]
	}
}

// Top returns a pointer to the top or
// inner most Frame of the environment's frame stack
func (env *Environment) Top() *Frame {
	if len(env.Frames) > 0 {
		index := len(env.Frames) - 1
		return env.Frames[index]
	}
	return nil
}

// Botttom returns a pointer to the bottom or// outer most Frame of the environment
func (env *Environment) Bottom() *Frame {
	if len(env.Frames) > 0 {
		return env.Frames[0]
	}
	return nil
}

// Defines the variable in the given scope level
func (env *Environment) Define(name string, val Value, level int) Value {
	frame := env.Frame(level)
	if frame == nil {
		return ErrorFromString("no such frame available.")
	}
	frame.Variables[name] = val
	return val
}

// Looks up the variable and sets it in the scope where it is found.
// Returns an error if no such variable could be found.
func (env *Environment) Set(name string, val Value) Value {
	_, frame := env.LookupFrame(name)
	if frame == nil {
		return ErrorFromString("no such variable")
	}
	frame.Variables[name] = val
	return val
}

func (env *Environment) Rescuer() Value {
	frame := env.Top()
	if frame == nil {
		return nil
	}
	return frame.Rescuer
}

// Prevent sets the rescue block to use for
// the top frame of the environment.
// It returns the previous rescuer.
func (env *Environment) Prevent(block Block) Value {
	frame := env.Frame(1)
	if frame == nil {
		return env.FailString("Could not set rescuer")
	}
	old := frame.Rescuer
	frame.Rescuer = Rescue{block}
	return old
}

//
func (env *Environment) Rescue(res Value) Value {
	flow := ValueFlow(res)
	if flow < FailFlow {
		return res
	}
	// if there is no rescue installed,
	// just return as is.
	if env.Rescuer() == nil {
		return res
	}
	eff := res.(Effect)
	// failures become normal returns
	// if the rescue didn't fail.
	val := eff.Unwrap()
	rres := env.Rescuer().Eval(env, val, res)
	reff, rok := rres.(Effect)
	if !rok {
		return env.Return(rres)
	} else {
		// Here, unpack the effect and replace it with a return
		// to avoid recursion loops of fail in rescue
		return env.Return(reff.Unwrap())
	}
}

func (env *Environment) Return(val Value) EffectValue {
	return Return{val}
}

func (env *Environment) Fail(err *Error) EffectValue {
	return err
}

func (env *Environment) Break(val Value) EffectValue {
	return Break{val}
}

func (env *Environment) FailString(msg string, args ...Value) Value {
	return env.Fail(env.ErrorFromString(msg, args...))
}

func (env Environment) Interpolate(s string, args ...Value) string {
	runes := []rune(s)
	res := []rune{}
	name := []rune{}
	inName := 0
	for i, a := range args {
		env.Define(Itoa(i+1), a, 0)
	}
	apply := func() {
		inName = 0
		val := env.Lookup(string(name))
		if val == nil {
			res = append(res, '!', 'n', 'i', 'l')
		} else {
			add := []rune(val.String())
			res = append(res, add...)
		}
		name = []rune{}
	}
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch r {
		case '$':
			if inName == 0 {
				inName = 1
			} else if inName == 1 {
				if len(name) < 1 {
					// $$ escape
					res = append(res, '$')
					inName = 0
				} else { // $ at end of name
					apply()
				}
			}
		case '{':
			if inName > 0 {
				inName++
			} else {
				res = append(res, '{')
			}
		case '}':
			if inName > 0 {
				inName--
				if inName == 1 {
					apply()
				}
			} else {
				res = append(res, '}')
			}
		default:
			if inName > 0 {
				if IsNumber(r) || IsLetter(r) {
					name = append(name, r)
				} else {
					apply()
					res = append(res, r)
				}
			} else {
				res = append(res, r)
			}
		}
	}
	if len(name) > 0 {
		apply()
	}
	return string(res)
}

func (env Environment) Printi(msg string, args ...Value) (int, error) {
	msg = env.Interpolate(msg, args...)
	return env.Write(msg)
}

func (env Environment) Write(msg string) (int, error) {
	buf := []byte(msg)
	writer := env.Out
	if len(env.Frames) > 0 {
		writer = env.Frames[len(env.Frames)-1].Out
	}
	if writer == nil {
		return -1, env.ErrorFromString("no writer set in environment.")
	}
	return writer.Write(buf)
}

func (env Environment) ErrorFromString(msg string, args ...Value) *Error {
	msg = env.Interpolate(msg, args...)
	return ErrorFromString(msg)
}

// Complete is for use with liner
func (env Environment) Complete(prefix String) List {
	res := List{}
	for _, frame := range env.Frames {
		for name, _ := range frame.Variables {
			if len(name) >= len(prefix) {
				if String(name[0:len(prefix)]) == prefix {
					res = append(res, String(name))
				}
			}
		}
	}
	if len(res) == 0 {
		res = append(res, prefix)
	}
	return res.SortStrings()
}

func (env *Environment) Overload(name string, target Value, types []Value) Value {
	val := env.Lookup(name)
	cov, ok := val.(Overload)
	if val == nil {
		cov = make(Overload)
	} else if !ok {
		return env.FailString("Not a overload: " + name)
	}
	signature := ""
	for _, arg := range types {
		signature += "_" + arg.String()
	}
	if _, ok := target.(String); ok {
		tarVal := env.Lookup(target.String())
		cov[signature] = tarVal
	} else if _, ok := target.(Word); ok {
		tarVal := env.Lookup(target.String())
		cov[signature] = tarVal
	} else {
		cov[signature] = target
	}

	env.Define(name, cov, -1)
	return cov
}
