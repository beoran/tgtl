package tgtl

//Converter is an interface that Values can optionally implement
// to allow conversion to other arbitrary types at run time.
type Converter interface {
	Convert(to interface{}) *Error
}

func (from Int) Convert(to interface{}) *Error {
	switch toPtr := to.(type) {
	case *string:
		(*toPtr) = from.String()
	case *int8:
		(*toPtr) = int8(from)
	case *int16:
		(*toPtr) = int16(from)
	case *int32:
		(*toPtr) = int32(from)
	case *int64:
		(*toPtr) = int64(from)
	case *int:
		(*toPtr) = int(from)
	case *bool:
		(*toPtr) = (from != 0)
	case *Bool:
		(*toPtr) = (from != 0)
	case *float32:
		(*toPtr) = float32(from)
	case *float64:
		(*toPtr) = float64(from)
	case *Int:
		(*toPtr) = from
	case *Value:
		(*toPtr) = from
	default:
		return ErrorFromString("Cannot convert Int value")
	}
	return nil
}

func (from Bool) Convert(to interface{}) *Error {
	iVal := 0
	if from {
		iVal = -1
	}
	switch toPtr := to.(type) {
	case *string:
		(*toPtr) = from.String()
	case *int8:
		(*toPtr) = int8(iVal)
	case *int16:
		(*toPtr) = int16(iVal)
	case *int32:
		(*toPtr) = int32(iVal)
	case *int64:
		(*toPtr) = int64(iVal)
	case *int:
		(*toPtr) = int(iVal)
	case *bool:
		(*toPtr) = bool(from)
	case *Bool:
		(*toPtr) = from
	case *float32:
		(*toPtr) = float32(iVal)
	case *float64:
		(*toPtr) = float64(iVal)
	case *Int:
		(*toPtr) = Int(iVal)
	case *Value:
		(*toPtr) = from
	default:
		return ErrorFromString("Cannot convert Int value")
	}
	return nil
}

func (from Word) Convert(to interface{}) *Error {
	switch toPtr := to.(type) {
	case *string:
		(*toPtr) = from.String()
	case *bool:
		(*toPtr) = (from.String() != "")
	case *Bool:
		(*toPtr) = (from.String() != "")
	case *Word:
		(*toPtr) = from
	case *Type:
		(*toPtr) = Type(string(from))
	case *Value:
		(*toPtr) = from
	default:
		return ErrorFromString("Cannot convert Word value")
	}
	return nil
}

func (from Type) Convert(to interface{}) *Error {
	switch toPtr := to.(type) {
	case *string:
		(*toPtr) = from.String()
	case *bool:
		(*toPtr) = (from.String() != "")
	case *Bool:
		(*toPtr) = (from.String() != "")
	case *Type:
		(*toPtr) = from
	case *Word:
		(*toPtr) = Word(string(from))
	case *Value:
		(*toPtr) = from
	default:
		return ErrorFromString("Cannot convert Word value")
	}
	return nil
}

func (from String) Convert(to interface{}) *Error {
	switch toPtr := to.(type) {
	case *string:
		(*toPtr) = from.String()
	case *bool:
		(*toPtr) = (from.String() != "")
	case *Bool:
		(*toPtr) = (from.String() != "")
	case **Error:
		(*toPtr) = ErrorFromString(from.String())
	case *String:
		(*toPtr) = from
	case *Value:
		(*toPtr) = from
	default:
		return ErrorFromString("Cannot convert String value")
	}
	return nil
}

func (from *Error) Convert(to interface{}) *Error {
	switch toPtr := to.(type) {
	case *string:
		(*toPtr) = from.String()
	case *bool:
		(*toPtr) = (from == nil)
	case *Bool:
		(*toPtr) = (from == nil)
	case *error:
		(*toPtr) = from
	case **Error:
		(*toPtr) = from
	case *Value:
		(*toPtr) = from
	default:
		return ErrorFromString("Cannot convert Error value")
	}
	return nil
}

func (from Block) Convert(to interface{}) *Error {
	switch toPtr := to.(type) {
	case *bool:
		(*toPtr) = (len(from.Statements) > 0)
	case *Bool:
		(*toPtr) = (len(from.Statements) > 0)
	case *Block:
		(*toPtr) = from
	case *Value:
		(*toPtr) = from
	default:
		return ErrorFromString("Cannot convert block value")
	}
	return nil
}

func (from Map) Convert(to interface{}) *Error {
	switch toPtr := to.(type) {
	case *bool:
		(*toPtr) = (len(from) > 0)
	case *Bool:
		(*toPtr) = (len(from) > 0)
	case *Map:
		(*toPtr) = from
	case *Value:
		(*toPtr) = from
	default:
		return ErrorFromString("Cannot convert map value")
	}
	return nil
}

func (from List) Convert(to interface{}) *Error {
	switch toPtr := to.(type) {
	case *bool:
		(*toPtr) = (len(from) > 0)
	case *Bool:
		(*toPtr) = (len(from) > 0)
	case *List:
		(*toPtr) = from
	case *Value:
		(*toPtr) = from
	default:
		return ErrorFromString("Cannot convert map value")
	}
	return nil
}

func ValToBool(val Value) bool {
	if val == nil {
		return false
	}
	switch check := val.(type) {
	case *Error:
		return check == nil
	case Int:
		return (int(check) != 0)
	case Bool:
		return bool(check)
	default:
		return val.String() != ""
	}
}

func Convert(val Value, to interface{}) *Error {
	if converter, ok := val.(Converter); ok {
		return converter.Convert(to)
	} else if pval, ok := to.(*Value); ok {
		*pval = val
		return nil
	}
	return ErrorFromString("Value cannot be converted")
}

// StringList makes a List from string arguments
func StringList(sa ...string) List {
	list := List{}
	for _, s := range sa {
		list = append(list, String(s))
	}
	return list
}

// Converts a list to raw strings
func (l List) ToStrings() []string {
	res := []string{}
	for _, s := range l {
		res = append(res, s.String())
	}
	return res
}
