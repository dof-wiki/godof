package parser

type ErrEmptyValue struct{}

func (e *ErrEmptyValue) Error() string {
	return "empty value"
}

type ErrValueType struct{}

func (e *ErrValueType) Error() string {
	return "err value type"
}
