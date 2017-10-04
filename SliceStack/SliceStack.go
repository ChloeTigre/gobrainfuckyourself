package SliceStack

import "errors"

type UIntStack []uint

type Stack interface {
	Peek() (interface{}, error)
	Push(interface{}) (Stack, error)
	Pop() (Stack, interface{}, error)
}

func (s UIntStack) Peek() (value uint, err error) {
	if len(s) == 0 {
		err = errors.New("no value")
		return
	}
	value = s[len(s)-1]
	return
}

func (s UIntStack) Push(value uint) (rs UIntStack, err error) {
	rs = append(s[:], value)
	return
}

func (s UIntStack) Pop() (retstack UIntStack, value uint, err error) {
	value, err = s.Peek()
	retstack = s
	if err == nil {
		retstack = s[:len(s)-1]
	}
	return
}
