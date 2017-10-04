package gobrainfuckyourself

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ChloeTigre/gobrainfuckyourself/SliceStack"
	"github.com/golang/example/stringutil"
	"strings"
	_ "time"
)

type memoryLocation uint
type bfOperator int

// these two definitions are bound (they should be in the same order)
const bfOperators string = "><+-.,[]"
const (
	OPERATOR_RIGHT bfOperator = iota
	OPERATOR_LEFT
	OPERATOR_INC
	OPERATOR_DEC
	OPERATOR_OUTPUT
	OPERATOR_INPUT
	OPERATOR_LOOP_TOP
	OPERATOR_LOOP_BOTTOM
	OPERATOR_INVALID = -1
)

const (
	OPMASK_RIGHT bfOperator = 1 << (2 * iota)
	OPMASK_LEFT
	OPMASK_INC
	OPMASK_DEC
	OPMASK_OUTPUT
	OPMASK_INPUT
	OPMASK_LOOP_TOP
	OPMASK_LOOP_BOTTOM
	OPMASK_INVALID = 0
)

type BFMachine struct {
	InstructionPointer, DataPointer uint
	Memory                          []byte
	Program                         string
	jumpStack                       SliceStack.UIntStack
	OutputBuffer                    bytes.Buffer
	loopCounter                     uint
}

type BFMachineState struct {
	InstructionPointer, DataPointer uint
	Memory                          []byte
	jumpStack                       SliceStack.UIntStack
	waitForInput, waitForOutput     bool
	loopCounter                     uint
}

func CreateBFMachine() (machine *BFMachine, err error) {
	bfm := BFMachine{}
	bfm.Memory = make([]byte, 60000)
	bfm.jumpStack = make(SliceStack.UIntStack, 1024)
	machine = &bfm
	err = nil
	return
}

func (bfm *BFMachine) LoadProgram(prog string) (err error) {
	bfm.Program = prog
	bfm.InstructionPointer = 0
	bfm.DataPointer = 0
	return
}

func (bfm *BFMachine) Step() (err error) {
	var nextState *BFMachineState
	nextState, err = bfm.EvalNextStep()
	if err != nil {
		return
	}
	if len(nextState.Memory) > 0 {
	}
	err = bfm.updateMachineState(nextState)
	if err != nil {
		panic("fuck2")
	}
	if nextState.waitForOutput {
		// fmt.Print(";", bfm.Memory[bfm.DataPointer], ",")
		bfm.OutputBuffer.WriteByte(bfm.Memory[bfm.DataPointer])
	}
	return
}

func (bfm *BFMachine) updateMachineState(nextState *BFMachineState) (err error) {
	bfm.InstructionPointer = nextState.InstructionPointer
	bfm.DataPointer = nextState.DataPointer
	if len(nextState.Memory) > 0 {
		// a := copy(nextState.Memory, bfm.Memory)
		bfm.Memory = []byte{}
		bfm.Memory = append(bfm.Memory, nextState.Memory...)
	}
	if len(nextState.jumpStack) > 0 {
		bfm.jumpStack = SliceStack.UIntStack{}
		bfm.jumpStack = append(bfm.jumpStack, nextState.jumpStack...)
	}
	bfm.loopCounter = nextState.loopCounter

	return
}

func getCommandCode(theChar byte) (operator bfOperator, err error) {
	code := strings.Index(bfOperators, string(theChar))
	if code == -1 {
		operator = OPERATOR_INVALID
		err = errors.New("Not an operator")
	} else {
		operator = bfOperator(code)
		err = nil
	}
	return
}

func (bfm *BFMachine) EvalNextStep() (nextStep *BFMachineState, err error) {
	if int(bfm.InstructionPointer) >= len(bfm.Program) {
		err = errors.New("Execution finished")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("INFO: ", bfm.Info())
			panic(r)
		}
	}()

	var command bfOperator
	var errint error
	nextStep = &BFMachineState{
		InstructionPointer: bfm.InstructionPointer,
		DataPointer:        bfm.DataPointer,
		Memory:             make([]byte, 0),
		loopCounter:        bfm.loopCounter,
	}
	command, errint = getCommandCode(bfm.Program[bfm.InstructionPointer])
	if errint != nil {
		nextStep.InstructionPointer = bfm.InstructionPointer + 1
		return
	}
	// when we're in an operation that changes some memory zones, copy them
	// to the State
	if command == OPERATOR_INC || command == OPERATOR_DEC {
		//nextStep.Memory = make([]byte, len(bfm.Memory))
		//copy(bfm.Memory, nextStep.Memory)
		nextStep.Memory = []byte{}
		nextStep.Memory = append(nextStep.Memory, bfm.Memory...)
	}
	if command == OPERATOR_LOOP_TOP || command == OPERATOR_LOOP_BOTTOM {
		// nextStep.jumpStack = make(SliceStack.UIntStack, 1024)
		nextStep.jumpStack = SliceStack.UIntStack{}
		nextStep.jumpStack = append(nextStep.jumpStack, bfm.jumpStack...)
		// copy(bfm.jumpStack, nextStep.jumpStack)
	}

	switch command {
	case OPERATOR_DEC:
		nextStep.Memory[bfm.DataPointer] = bfm.Memory[bfm.DataPointer] - 1
		break
	case OPERATOR_INC:
		nextStep.Memory[bfm.DataPointer] = bfm.Memory[bfm.DataPointer] + 1
		break
	case OPERATOR_LEFT:
		nextStep.DataPointer = bfm.DataPointer - 1
		break
	case OPERATOR_RIGHT:
		nextStep.DataPointer = bfm.DataPointer + 1
		break
	case OPERATOR_OUTPUT:
		nextStep.waitForOutput = true
		break
	case OPERATOR_INPUT:
		nextStep.waitForInput = true
		panic("input not implemented")
	case OPERATOR_LOOP_TOP:
		// when data pointer has a non-0 value
		// then store next-command location for a jump, and move forward
		// else, set skip bit so we will skip commands until we meet a ]
		nextStep.loopCounter = bfm.loopCounter + 1
		if bfm.Memory[bfm.DataPointer] == 0 {
			//nextStep.jumpStack, nextStep.InstructionPointer, err = nextStep.jumpStack.Pop()
			nextStep.InstructionPointer, err = forwardScan([]rune(bfm.Program), bfm.InstructionPointer)
			if err != nil {
				panic(err)
			}
			return
		}
		break
	case OPERATOR_LOOP_BOTTOM:
		if bfm.Memory[bfm.DataPointer] != 0 {
			nextStep.InstructionPointer, err = backwardScan([]rune(bfm.Program), bfm.InstructionPointer)
			if err != nil {
				panic("bbbbbh")
			}
			return
		} else {
		}
		break
	case OPERATOR_INVALID:
		break
	default:
		panic("unknown operator")
	}
	nextStep.InstructionPointer = bfm.InstructionPointer + 1
	return
}

func forwardScan(program []rune, startPosition uint) (pos uint, err error) {
	occ := 1
	for i, e := range program[startPosition+1:] {
		switch e {
		case '[':
			occ += 1
			break
		case ']':
			occ -= 1
			break
		default:
			break
		}
		if occ == 0 {
			pos = uint(i) + 1
			return
		}
	}
	err = errors.New("Asymmetric brackets")
	return
}

func backwardScan(program []rune, startPosition uint) (pos uint, err error) {
	occ := 1
	sprog := string(program)
	rprog := stringutil.Reverse(sprog)
	actualStartPosition := len(sprog) - int(startPosition)
	for i, e := range rprog[actualStartPosition:] {
		switch e {
		case ']':
			occ += 1
			break
		case '[':
			occ -= 1
			break
		default:
			break
		}
		if occ == 0 {
			pos = startPosition - uint(i)
			return
		}
	}
	err = errors.New("Asymmetric brackets")
	return
}

func (bfm *BFMachine) Info() string {
	//return fmt.Sprint(rune(bfm.Program[bfm.InstructionPointer]))
	return fmt.Sprintf(`
---
Data Pointer: %d
Instruction Pointer: %d
Loop counter: %d
Memory head: %+v
---`,

		bfm.DataPointer, bfm.InstructionPointer,
		bfm.loopCounter,
		bfm.Memory[0:128])

}

func RunProgram(program string) {
	bfm, err := CreateBFMachine()
	if err != nil {
		panic("oops")
	}
	bfm.Program = program
	var i int = 0
	for {
		i += 1
		err = bfm.Step()
		if bfm.OutputBuffer.Len() > 0 {
			data := bfm.OutputBuffer.Next(1024)
			fmt.Print(string(data))
		}
		if err != nil {
			fmt.Print(bfm.Info())
			break
		}
		// time.Sleep(100 * time.Millisecond)
	}
}
