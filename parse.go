package tgtl

// ParseFunc is a parser function.
// It parses the input input starting from *index, which must be
// guaranteed by the caller to be non-nil.
// It should return as follows:
// * If the parse function matched what it is intended to parse
//   it should return the parsed value, nil, and index should be moved to
//   point right after the parsed part of te string.
// * If the parse function did not match what it is intended to parse
//   it should retirn nil, nil, and index  should be unchanged.
// * If the parse function did match what it is intended to parse
//   but there is a parse error, it should return nil, *Error,
//   and index should be set to the error location.
type ParseFunc func(input []rune, index *int) (Value, *Error)

var Debug = false

func debug(msg string) {
	if Debug {
		print(msg)
	}
}

func ParseAlternative(input []rune, index *int, funcs ...ParseFunc) (Value, *Error) {
	for _, fun := range funcs {
		val, err := fun(input, index)
		if err != nil || val != nil {
			return val, err
		}
	}
	return nil, nil
}

func ParseWhileRuneOk(input []rune, index *int, ok func(r rune) bool) (Value, *Error) {
	length := len(input)
	start, now := *index, 0
	for ; *index < length; *index++ {
		r := input[*index]
		if !ok(r) {
			if now == 0 {
				return nil, nil
			}
			return String(input[start:*index]), nil
		}
		now++
	}
	return nil, ErrorFromString("unexpected EOF: >" + string(input[start:*index]) + "<")
}

type LineInfo struct {
	Line int
	From int
	To   int
}

type LineIndex []LineInfo

func PhysicalLineIndex(input []rune) LineIndex {
	res := LineIndex{}
	line, last, index := 0, 0, 0
	for ; index < len(input); index++ {
		ch := input[index]
		if ch == '\n' {
			line++
			li := LineInfo{line, last, index}
			last = index
			res = append(res, li)
		}
	}
	li := LineInfo{line, last, index}
	res = append(res, li)
	return res
}

func (li LineIndex) Lookup(index int) (row, col int) {
	for _, info := range li {
		if index >= info.From && index < info.To {
			return info.Line, index - info.From
		}
	}
	return -1, -1
}

func Parse(input string) (value Value, rerr *Error) {
	index := 0
	return ParseScript([]rune(input), &index)
}

func ParseScript(input []rune, index *int) (value Value, rerr *Error) {
	defer func() {
		val := recover()
		err, ok := val.(*Error)
		if ok {
			rerr.Children = append(rerr.Children, err)
		}
	}()
	value, rerr = ParseStatements([]rune(input), index)
	if value != nil {
		value = Block{value.(List)}
	}
	return value, rerr
}

func IsEof(input []rune, index *int) bool {
	return *index >= len(input)
}

func ParseStatements(input []rune, index *int) (Value, *Error) {
	debug("ParseStatements")
	statements := List{}
	for {
		val, err := ParseStatement(input, index)
		if err != nil {
			debug("error in statement")
			return nil, err
		}
		if val != nil {
			statements = append(statements, val)
		}
		sep, err := ParseRs(input, index)
		if IsEof(input, index) {
			return statements, nil
		}
		if err != nil {
			debug("error in rs")
			return nil, err
		}
		if sep == nil {
			return statements, nil
		}
	}
}

func ParseRs(input []rune, index *int) (Value, *Error) {
	debug("ParseRs")
	SkipWs(input, index)
	return ParseWhileRuneOk(input, index, func(r rune) bool {
		return r == '\n' || r == '\r' || r == ';'
	})
}

func ParseWs(input []rune, index *int) (Value, *Error) {
	debug("ParseWs")
	return ParseWhileRuneOk(input, index, func(r rune) bool {
		return r == ' ' || r == '\t'
	})
}

func ParseWsRs(input []rune, index *int) (Value, *Error) {
	debug("ParseRs")
	SkipWs(input, index)
	return ParseWhileRuneOk(input, index, func(r rune) bool {
		return r == '\n' || r == '\r' || r == ';' || r == ' ' || r == '\t'
	})
}

func SkipWs(input []rune, index *int) {
	ParseWs(input, index)
}

func SkipRs(input []rune, index *int) {
	ParseRs(input, index)
}

func SkipWsRs(input []rune, index *int) {
	ParseWsRs(input, index)
}

func ParseComment(input []rune, index *int) (Value, *Error) {
	debug("ParseComment")
	start := *index
	if !RequireRune(input, index, '#') {
		return nil, nil
	}
	for ; *index < len(input); *index++ {
		r := input[*index]
		if r == '\n' || r == '\r' {
			end := *index
			return Comment(string(input[start:end])), nil
		}
	}
	return nil, ErrorFromString("unexpected EOF in comment")
}

func ParseStatement(input []rune, index *int) (Value, *Error) {
	debug("ParseStatement")
	SkipWs(input, index)
	return ParseAlternative(input, index, ParseCommand, ParseBlock, ParseComment)
}

func ParseParameters(input []rune, index *int) (Value, *Error) {
	debug("ParseParameters")
	params := List{}
	for {
		sep, err := ParseWs(input, index)
		if err != nil {
			return nil, err
		}
		if sep == nil {
			return params, nil
		}
		val, err := ParseParameter(input, index)
		if err != nil {
			return nil, err
		}
		if val == nil {
			return params, nil
		}
		params = append(params, val)
	}
}

func ParseParameter(input []rune, index *int) (Value, *Error) {
	debug("ParseParameter")
	funcs := []ParseFunc{ParseLiteral, ParseEvaluation, ParseBlock, ParseGetter}
	return ParseAlternative(input, index, funcs...)
}

func ParseOrder(input []rune, index *int) (Value, *Error) {
	debug("ParseOrder")
	return ParseAlternative(input, index, ParseLiteral, ParseEvaluation)
}

func ParseCommand(input []rune, index *int) (Value, *Error) {
	debug("ParseCommand")
	order, err := ParseOrder(input, index)
	if err != nil || order == nil {
		return order, err
	}
	params, err := ParseParameters(input, index)
	if err != nil {
		return params, err
	}
	if params == nil {
		params = List{}
	}
	return Command{order, params.(List)}, nil
}

// RequireRune requires a single rune to be present,
// and skips it, however that rune is discared.
// Returns true if the rune was found, false if not
func RequireRune(input []rune, index *int, req rune) bool {
	if input[*index] == req {
		*index++
		return true
	}
	return false
}

func ParseEvaluation(input []rune, index *int) (Value, *Error) {
	debug("ParseEvaluation")
	if !RequireRune(input, index, '[') {
		return nil, nil
	}
	res, err := ParseCommand(input, index)
	if err != nil {
		return nil, err
	}
	if !RequireRune(input, index, ']') {
		print(input[*index])
		return nil, ErrorFromString("Expected end of evaluation ]")
	}
	if res != nil {
		res = Evaluation{Command: res.(Command)}
	}
	return res, nil
}

func ParseBlock(input []rune, index *int) (Value, *Error) {
	debug("ParseBlock")
	if !RequireRune(input, index, '{') {
		return nil, nil
	}
	res, err := ParseStatements(input, index)
	if err != nil {
		return nil, err
	}
	SkipWsRs(input, index)
	if !RequireRune(input, index, '}') {
		return nil, ErrorFromString("Expected end of block }")
	}
	return Block{Statements: res.(List)}, nil
	return nil, nil
}

func ParseGetter(input []rune, index *int) (Value, *Error) {
	debug("ParseGetter")
	if RequireRune(input, index, '$') {
		if input[*index] == '$' { // recusively parse double getters
			val, err := ParseGetter(input, index)
			if err == nil { // Getter with a getter inside.
				return Getter{val}, err
			} else {
				return nil, err
			}
		} else { // integer, sring or getter name
			key, err := ParseLiteral(input, index)
			if key == nil {
				return nil, ErrorFromString("Expected literal after getter $")
			}
			if err == nil {
				return Getter{key}, nil
			}
			return nil, err
		}
	}
	return nil, nil
}

func ParseLiteral(input []rune, index *int) (Value, *Error) {
	debug("ParseLiteral")
	return ParseAlternative(input, index, ParseWord, ParseString, ParseInteger,
		ParseRawString)
}

func IsLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r > rune(128)) ||
		r == '_' || r == '/'
}

func IsNumber(r rune) bool {
	return (r >= '0' && r <= '9')
}

func ParseWord(input []rune, index *int) (Value, *Error) {
	debug("ParseWord")
	// a word consists of an ascii letter or non asci characters, or underscore
	// followed by an ascii letter or number, or non ascii characters, or underscore
	start := *index
	r := input[*index]
	if !IsLetter(r) {
		return nil, nil
	}
	for *index++; *index < len(input); *index++ {
		r := input[*index]
		if !(IsLetter(r) || IsNumber(r)) {
			return Word(string(input[start:*index])), nil
		}
	}
	return nil, ErrorFromString("unexpected EOF in string")
}

func next(input []rune, index *int) {
	*index++
	if *index >= len(input) {
		panic(ErrorFromString("Unexpected end of input."))
	}
}

func ParseEscape(input []rune, index *int) (Value, *Error) {
	res := ""
	if input[*index] != '\\' {
		return nil, nil
	}
	next(input, index)
	switch input[*index] {
	case 'a':
		res += "\a"
	case 'b':
		res += "\b"
	case 'e':
		res += "\033"
	case 'f':
		res += "\f"
	case 'n':
		res += "\n"
	case 'r':
		res += "\r"
	case 't':
		res += "\t"
	case '\\':
		res += "\\"
	case '"':
		res += "\""
	default:
		return nil, ErrorFromString("Unknown escape sequence character")
	}

	return String(res), nil
}

func ParseString(input []rune, index *int) (Value, *Error) {
	debug("ParseString")
	res := ""
	ch := input[*index]
	if ch != '"' {
		return nil, nil
	}
	*index++
	for *index < len(input) {
		ch = input[*index]
		esc, err := ParseEscape(input, index)
		if err != nil {
			return nil, err
		}
		if esc != nil {
			res += string(esc.(String))
		} else if ch == '"' {
			*index++
			return String(res), nil
		} else {
			res += string(ch)
		}
		*index++
	}
	return nil, ErrorFromString("Unexpected end of input.")
}

func ParseRawString(input []rune, index *int) (Value, *Error) {
	debug("ParseRawString")
	res := ""
	ch := input[*index]
	if ch != '`' {
		return nil, nil
	}
	*index++
	for *index < len(input) {
		ch = input[*index]
		if ch == '`' {
			*index++
			return String(res), nil
		} else {
			res += string(ch)
		}
		*index++
	}
	return nil, ErrorFromString("Unexpected end of input.")
}

func ParseInteger(input []rune, index *int) (Value, *Error) {
	debug("ParseInteger")
	ch := input[*index]
	neg := 1
	res := 0
	if ch == '-' {
		neg = -1
	} else if ch == '+' {
		// do nothing, ignore + as an integer prefix
	} else {
		res = int(ch - '0')
		if res < 0 || res > 9 { // Not a digit, no integer
			return nil, nil
		}
	}
	*index++
	for *index < len(input) {
		ch = input[*index]
		ch -= '0'
		if ch < 0 || ch > 9 { // Not a digit, finished
			return Int(neg * res), nil
		}
		res = res * 10
		res = res + int(ch)
		*index++
	}
	return nil, ErrorFromString("unexpected EOF in number")
}
