package tgtl

import "testing"
import "reflect"

type testCase struct {
	ParseFunc
	input         string
	index         int
	expectedIndex int
	expectError   bool
	expectedValue Value
}

func (tc *testCase) Run(t *testing.T) {
	t.Logf("Test case input: %s", tc.input)
	res, err := tc.ParseFunc([]rune(tc.input), &tc.index)
	if tc.expectError {
		if err == nil {
			t.Errorf("expected parse error")
			return
		}
	} else {
		if err != nil {
			t.Errorf("error: unexpected parse error at %d: %v", tc.index, err)
			return
		}
	}
	if tc.index != tc.expectedIndex {
		if err == nil {
			t.Errorf("error: index not correct: %d <> %d", tc.index, tc.expectedIndex)
		}
	}
	if tc.expectedValue == nil {
		if res != nil {
			t.Errorf("error: expected nil value, got %v", res)
		} else {
			t.Logf("Test case value: nil")
		}
	} else {
		if res == nil {
			t.Errorf("error: expected value %v was nil", tc.expectedValue)
		} else {
			t.Logf("Test case value: %v", res)
			se := tc.expectedValue.String()
			so := res.String()
			if so != se {
				t.Errorf("error: value is not as expected: %s <-> %s", so, se)
			} else {
				if !reflect.DeepEqual(tc.expectedValue, res) {
					t.Errorf("error: values are not deeply equal: %s <-> %s", so, se)
				}
			}
		}
	}
	rs := []rune(tc.input)
	b := string(rs[0:tc.index])
	a := string(rs[tc.index:len(rs)])
	t.Logf("Test case index: %s|%s", b, a)
}

func TestParseComment(t *testing.T) {
	tcs := []testCase{
		testCase{ParseComment, "# comment\nnot comment", 0, 9, false, Comment("# comment")},
		testCase{ParseComment, "not comment\n#comment\n", 0, 0, false, nil},
		testCase{ParseComment, "# comment\nnot comment", 9, 9, false, nil},
	}
	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseWs(t *testing.T) {
	tcs := []testCase{
		testCase{ParseWs, "# comment   ", 0, 0, false, nil},
		testCase{ParseWs, "1234567890  ", 0, 0, false, nil},
		testCase{ParseWs, "-1234567890 ", 0, 0, false, nil},
		testCase{ParseWs, "+1234567890 ", 0, 0, false, nil},
		testCase{ParseWs, "`string`    ", 0, 0, false, nil},
		testCase{ParseWs, `"string"    `, 0, 0, false, nil},
		testCase{ParseWs, "word        ", 0, 0, false, nil},
		testCase{ParseWs, " \tword     ", 0, 2, false, String(" \t")},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseGetter(t *testing.T) {
	tcs := []testCase{
		testCase{ParseGetter, "not a getter ", 0, 0, false, nil},
		testCase{ParseGetter, "$# comment   ", 0, 0, true, nil},
		testCase{ParseGetter, "$ bad space  ", 0, 0, true, nil},
		testCase{ParseGetter, "$1234567890  ", 0, 11, false, Getter{Int(+1234567890)}},
		testCase{ParseGetter, "$-1234567890 ", 0, 12, false, Getter{Int(-1234567890)}},
		testCase{ParseGetter, "$+1234567890 ", 0, 12, false, Getter{Int(+1234567890)}},
		testCase{ParseGetter, "$`string`    ", 0, 9, false, Getter{String("string")}},
		testCase{ParseGetter, `$"string"    `, 0, 9, false, Getter{String("string")}},
		testCase{ParseGetter, "$word        ", 0, 5, false, Getter{Word("word")}},
		testCase{ParseGetter, "$µ_rd        ", 0, 5, false, Getter{Word("µ_rd")}},
		testCase{ParseGetter, "$wo09        ", 0, 5, false, Getter{Word("wo09")}},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseInteger(t *testing.T) {
	tcs := []testCase{
		testCase{ParseInteger, "01234567890 ", 0, 11, false, Int(1234567890)},
		testCase{ParseInteger, "-1234567890 ", 0, 11, false, Int(-1234567890)},
		testCase{ParseInteger, "+1234567890 ", 0, 11, false, Int(1234567890)},
		testCase{ParseInteger, "+1234567890", 0, 11, true, nil},
		testCase{ParseInteger, "not an integer", 0, 0, false, nil},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseString(t *testing.T) {
	tcs := []testCase{
		testCase{ParseString, `"string"`, 0, 8, false, String("string")},
		testCase{ParseString, `"string\"\n"`, 0, 12, false, String("string\"\n")},
		testCase{ParseString, `"bad        `, 0, 12, true, nil},
		testCase{ParseString, `not a string`, 0, 0, false, nil},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseRawString(t *testing.T) {
	tcs := []testCase{
		testCase{ParseRawString, "`string` ", 0, 8, false, String("string")},
		testCase{ParseRawString, "`string\"quote` ", 0, 14, false,
			String("string\"quote")},
		testCase{ParseRawString, "`bad        ", 0, 11, true, nil},
		testCase{ParseRawString, "not a string` ", 0, 0, false, nil},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseWord(t *testing.T) {
	tcs := []testCase{
		testCase{ParseWord, "word ", 0, 4, false, Word("word")},
		testCase{ParseWord, "word\n", 0, 4, false, Word("word")},
		testCase{ParseWord, "word; ", 0, 4, false, Word("word")},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseLiteral(t *testing.T) {
	tcs := []testCase{
		testCase{ParseLiteral, "# comment   ", 0, 0, false, nil},
		testCase{ParseLiteral, "   \t       ", 0, 0, false, nil},
		testCase{ParseLiteral, "1234567890  ", 0, 10, false, Int(1234567890)},
		testCase{ParseLiteral, "-1234567890 ", 0, 11, false, Int(-1234567890)},
		testCase{ParseLiteral, "+1234567890 ", 0, 11, false, Int(1234567890)},
		testCase{ParseLiteral, "`string`    ", 0, 8, false, String("string")},
		testCase{ParseLiteral, `"string"    `, 0, 8, false, String("string")},
		testCase{ParseLiteral, "word        ", 0, 4, false, Word("word")},
		testCase{ParseLiteral, "µ_rd        ", 0, 4, false, Word("µ_rd")},
		testCase{ParseLiteral, "wo09        ", 0, 4, false, Word("wo09")},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseParameters(t *testing.T) {
	tcs := []testCase{
		testCase{ParseParameters, " world\n", 0, 6, false, List{Word("world")}},
		testCase{ParseParameters, ` world 7 "foo"` + "`\n", 0, 14, false,
			List{Word("world"), Int(7), String("foo")},
		},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseOrder(t *testing.T) {
	tcs := []testCase{
		testCase{ParseOrder, "hello ", 0, 5, false, Word("hello")},
		testCase{ParseOrder, "1 ", 0, 1, false, Int(1)},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseCommand(t *testing.T) {
	tcs := []testCase{
		testCase{ParseCommand, "hello world\n", 0, 11, false,
			Command{Word("hello"), List{Word("world")}},
		},
	}

	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseEvaluation(t *testing.T) {
	tcs := []testCase{
		testCase{ParseEvaluation, "[hello world]", 0, 13, false,
			Evaluation{Command{Word("hello"), List{Word("world")}}},
		},
		testCase{ParseEvaluation, "[hello 123 ]", 0, 12, false,
			Evaluation{Command{Word("hello"), List{Int(123)}}},
		},
	}
	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseAStatement(t *testing.T) {
	tcs := []testCase{
		testCase{ParseStatement, " hello world\n", 0, 12, false,
			Command{Word("hello"), List{Word("world")}},
		},
		testCase{ParseStatement, " hello;world;", 0, 6, false,
			Command{Word("hello"), List{}},
		},
	}
	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseStatements(t *testing.T) {
	tcs := []testCase{
		testCase{ParseStatements, "\n\n\n", 0, 3, false,
			List{},
		},
		testCase{ParseStatements, "\n\n \n", 0, 4, false,
			List{},
		},
		testCase{ParseStatements, "hello world\n", 0, 12, false,
			List{Command{Word("hello"), List{Word("world")}}},
		},
		testCase{ParseStatements, "hello;world;", 0, 12, false,
			List{Command{Word("hello"), List{}}, Command{Word("world"), List{}}},
		},
		testCase{ParseStatements, "hello \n world \n\n", 0, 16, false,
			List{Command{Word("hello"), List{}}, Command{Word("world"), List{}}},
		},
	}
	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseBlock(t *testing.T) {
	tcs := []testCase{
		testCase{ParseBlock, "{hello world}", 0, 13, false,
			Block{List{Command{Word("hello"), List{Word("world")}}}},
		},
		testCase{ParseBlock, "{hello;world}", 0, 13, false,
			Block{List{Command{Word("hello"), List{}}, Command{Word("world"), List{}}}},
		},
		testCase{ParseBlock, "{hello\nworld}", 0, 13, false,
			Block{List{Command{Word("hello"), List{}}, Command{Word("world"), List{}}}},
		},
		testCase{ParseBlock, "{ hello\n world\n}", 0, 16, false,
			Block{List{Command{Word("hello"), List{}}, Command{Word("world"), List{}}}},
		},
		testCase{ParseBlock, "{ #Comment\n hello\n world\n}", 0, 26, false,
			Block{List{Comment("#Comment"), Command{Word("hello"), List{}}, Command{Word("world"), List{}}}}},
	}
	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParse(t *testing.T) {
	script1 := `
	# Comment
	print "Hello world!"
`
	res1 := Block{List{Comment("# Comment"), Command{Word("print"), List{String("Hello world!")}}}}

	script2 := `
	# Comment
	print "Hello world!"
`
	res2 := Block{List{Comment("# Comment"), Command{Word("print"), List{String("Hello world!")}}}}

	tcs := []testCase{
		testCase{ParseScript, script1, 0, len(script1), false, res1},
		testCase{ParseScript, script2, 0, len(script1), false, res2},
	}
	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t)
	}
}

func TestParseAndRun(t *testing.T) {
	script1 := `
		# Comment
		print "Hello world!"
`
	parsed, err := Parse(script1)
	if err != nil {
		t.Errorf("Parse error: %v", err)
		return
	}
	if parsed == nil {
		t.Errorf("No parse results error: %v", parsed)
		return
	}
	env := &Environment{}
	env.Push()
	env.Define("print", Proc(func(e *Environment, args ...Value) (Value, Effect) {
		var msg string
		Args(args, &msg)
		print(msg, "\n")
		return nil, nil
	}), 0)
	parsed.Eval(env)
}

type iTestCase struct {
	in          string
	ex          string
	expectError bool
	args        []Value
}

func (tc *iTestCase) Run(t *testing.T, e Environment) {
	t.Logf("Test case input: %s", tc.in)
	res := e.Interpolate(tc.in, tc.args...)
	t.Logf("Test case result: %v", res)
	if res != tc.ex {
		t.Errorf("error: value not expected: %s <-> %s", res, tc.ex)
	}
}

func TestInterpolate(t *testing.T) {
	e := Environment{}
	e.Push()
	e.Define("foo", String("{world}"), 0)

	tcs := []iTestCase{
		iTestCase{`hello {world}`, `hello {world}`, false, []Value{}},
		iTestCase{`hello ${foo}`, `hello {world}`, false, []Value{}},
		iTestCase{`hello $${foo}`, `hello ${foo}`, false, []Value{}},
		iTestCase{`hello ${1}`, `hello {world}`, false, []Value{String("{world}")}},
	}
	for i, tc := range tcs {
		t.Logf("Case: %d", i+1)
		tc.Run(t, e)
	}
}
