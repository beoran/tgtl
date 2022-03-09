package tgtl

type Comparer func(v1, v2 Value) int

func (data List) Sort(compare Comparer) List {
	if len(data) < 2 {
		return data
	}
	pivot := data[0]
	smaller := make(List, 0, len(data))
	equal := make(List, 1, len(data))
	larger := make(List, 0, len(data))
	equal[0] = pivot
	for i := 1; i < len(data); i++ {
		cmp := compare(data[i], pivot)
		if cmp > 0 {
			larger = append(larger, data[i])
		} else if cmp < 0 {
			smaller = append(smaller, data[i])
		} else {
			equal = append(equal, data[i])
		}
	}
	res := smaller.Sort(compare)
	res = append(res, equal...)
	res = append(res, larger.Sort(compare)...)
	return res
}

func (data List) SortStrings() List {
	return data.Sort(func(v1, v2 Value) int {
		s1 := v1.String()
		s2 := v2.String()
		if s1 > s2 {
			return 1
		} else if s1 < s2 {
			return -1
		}
		return 0
	})
}

// Converts an integer to a string
func Itoa(i int) string {
	if i == 0 {
		return "0"
	}
	digits := "0123456789"
	res := ""
	neg := false
	if i < 0 {
		neg = true
		i = -i
	}
	for d := i % 10; i > 0; { // dumb digit by digit algorithm
		res = string(digits[d]) + res
		i = i / 10
		d = i % 10
	}
	if neg {
		res = "-" + res
	}
	return res
}

func Args(froms []Value, tos ...interface{}) *Error {
	for i, to := range tos {
		if i >= len(froms) {
			return ErrorFromString("Too few arguments: " +
				Itoa(len(froms)) + " in stead of " + Itoa(len(tos)))
		}
		from := froms[i]
		err := Convert(from, to)
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	maxRune      = '\U0010FFFF'
	surrogateMin = 0xD800
	surrogateMax = 0xDFFF
	rune1Max     = 1<<7 - 1
	rune2Max     = 1<<11 - 1
	rune3Max     = 1<<16 - 1
)

func RuneLen(r rune) int {
	switch {
	case r < 0:
		return -1
	case r <= rune1Max:
		return 1
	case r <= rune2Max:
		return 2
	case surrogateMin <= r && r <= surrogateMax:
		return -1
	case r <= rune3Max:
		return 3
	case r <= maxRune:
		return 4
	}
	return -1
}

// WordCompleter is for use with liner
func WordCompleter(env Environment, line string, pos int) (head string, c []string, tail string) {
	end := pos
	begin := pos - 1
	if begin < 0 {
		begin = 0
	}
	for ; begin > 0 && len(line[begin:]) > 0; begin-- {
		r := []rune(line[begin:])[0]
		if RuneLen(r) < 0 {
			continue
		}
		if !(IsLetter(r) || IsNumber(r)) {
			rl := RuneLen(r) // skip to next rune
			begin += rl
			break
		}
	}
	for end = pos; end < len(line); end++ {
		r := []rune(line[end:])[0]
		if RuneLen(r) < 0 {
			continue
		}
		if !(IsLetter(r) || IsNumber(r)) {
			break
		}
	}
	head = line[0:begin]
	tail = line[end:]
	word := line[begin:end]
	clist := env.Complete(String(word))
	return head, clist.ToStrings(), tail
}
