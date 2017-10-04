package gobrainfuckyourself

import (
  "testing"
)
const oobprog = `+[+++++++++++++++++++++++++++++++++.<]`
const helloworld = `
+++++ +++++             (initialize counter (cell #0) to 10)
[                       (use loop to set the next four cells to 70/100/30/10)
    > +++++ ++          (    add  7 to cell #1)
    > +++++ +++++       (    add 10 to cell #2 )
    > +++               (    add  3 to cell #3)
    > +                 (    add  1 to cell #4)
    <<<< -              (    decrement counter (cell #0))
]
> ++ .                  (print 'H')
> + .                   (print 'e')
+++++ ++ .              (print 'l')
.                       (print 'l')
+++ .                   (print 'o')
> ++ .                  (print ' ')
<< +++++ +++++ +++++ .  (print 'W')
> .                     (print 'o')
+++ .                   (print 'r')
----- - .               (print 'l')
----- --- .             (print 'd')
> + .                   (print bang)
> .                     (print '\n')
`

const factorials = `
  ;; factorials
  ;; this program prints out the sequence of factorials
  ;; written by Keymaker
  ;; does not terminate by itself


    ++++++++++>>>+>>>>+>+<[[+++++[>++++
    ++++<-]>.<++++++[>--------<-]+<<]<<
    [<<]<.>>>>+<[->[<+>-[<+>-[<+>-[<+>-
    [<+>-[<+>-[<+>-[<+>-[<+>-[<[-]>-+>[
    <->-]<[->>>[>>]<<[->[>>+<<-]>+<<<<]
    <]>[-]+>+<<]]]]]]]]]]<[>+<-]+>>]<<[
    <<]>>[->>[>>]>>[-<<[<<]<<[<<]>[>[>>
    ]>>[>>]>>[>>]>+>+<<<<[<<]<<[<<]<<[<
    <]>-]>[>>]>>[>>]>>[>>]>[<<<[<<]<<[<
    <]<<[<<]>+>[>>]>>[>>]>>[>>]>-]<<<[<
    <]>[>[>>]>+>>+<<<<<[<<]>-]>[>>]>[<<
    <[<<]>+>[>>]>-]>>[<[<<+>+>-]<[>>>+<
    <<-]<[>>+<<-]>>>-]<[-]>>+[>[>>>>]>[
    >>>>]>[-]+>+<[<<<<]>-]>[>>>>]>[>>>>
    ]>->-[<<<+>>+>-]<[>+<-]>[[<<+>+>-]<
    [<->-[<->-[<->-[<->-[<->-[<->-[<->-
    [<->-[<->-[<-<---------->>[-]>>>>[-
    ]+>+<<<<<]]]]]]]]]]<[>>+<<-]>>]<+<+
    <[>>>+<<<-]<<[<<<<]<<<<<[<<]+>>]>>>
    >>[>>>>]+>[>>>>]<<<<[-<<<<]>>>>>[<<
    <<]<<<<<[<<]<<[<<]+>>]>>[>>]>>>>>[-
    >>>>]<<[<<<<]>>>>[>>>>]<<<<<<<<[>>>
    >[<<+>>->[<<+>>-]>]<<<<[<<]<<]<<<<<
    [->[-]>>>>>>>>[<<+>>->[<<+>>-]>]<<<
    <[<<]<<<<<<<]>>>>>>>>>[<<<<<<<+>>>>
    >>>->[<<<<<<<+>>>>>>>-]>]<<<<<<<<<]
`
const quine = `->++>+++>+>+>++>>+>+>+++>>+>+>++>+++>+++>+>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>+>+>++>>>+++>>>>>+++>+>>>>>>>>>>>>>>>>>>>>>>+++>>>>>>>++>+++>+++>+>>+++>+++>+>+++>+>+++>+>++>+++>>>+>+>+>+>++>+++>+>+>>+++>>>>>>>+>+>>>+>+>++>+++>+++>+>>+++>+++>+>+++>+>++>+++>++>>+>+>++>+++>+>+>>+++>>>+++>+>>>++>+++>+++>+>>+++>>>+++>+>+++>+>>+++>>+++>>+[[>>+[>]+>+[<]<-]>>[>]<+<+++[<]<<+]>>>[>]+++[++++++++++>++[-<++++++++++++++++>]<.<-<]`

func TestOperators(t *testing.T) {
	var a bfOperator
	var err error
	if a, err = getCommandCode('>'); a != OPERATOR_RIGHT {
		t.Error("Wrong value for >")
	}
	if a, err = getCommandCode('<'); a != OPERATOR_LEFT {
		t.Error("Wrong value for <")
	}
	if a, err = getCommandCode('-'); a != OPERATOR_DEC {
		t.Error("Wrong value for -")
	}
	if a, err = getCommandCode('+'); a != OPERATOR_INC {
		t.Error("Wrong value for +")
	}
	if a, err = getCommandCode('.'); a != OPERATOR_OUTPUT {
		t.Error("Wrong value for .")
	}
	if a, err = getCommandCode(','); a != OPERATOR_INPUT {
		t.Error("Wrong value for ,")
	}
	if a, err = getCommandCode('['); a != OPERATOR_LOOP_TOP {
		t.Error("Wrong value for [")
	}
	if a, err = getCommandCode(']'); a != OPERATOR_LOOP_BOTTOM {
		t.Error("Wrong value for ]")
	}
	if err != nil {
		t.Error("There should not be errors")
	}
}

func TestBasicProgram(t *testing.T) {
  bfm, err := CreateBFMachine()
  if err != nil {
    t.Errorf("Could not create brainfuck machine: %+v", err)
    t.Fail()
  }
  var data []byte
  bfm.LoadProgram(".")
  bfm.Step()
  data = bfm.OutputBuffer.Next(1024)
  if data[0] != 0 {
    t.Errorf("Error in trivial program: output is %s, should be [0]\n%s\n", data, bfm.Info())
  }
  bfm.LoadProgram(oobprog)
  for {
    err = bfm.Step()
    if err != nil && err != OutOfBoundsError {
      t.Errorf("ERROR: +%v", err)
      break
    } else {
      break
    }
  }
  bfm.LoadProgram(helloworld)
  for {
    err = bfm.Step()
    if err == EndOfProgramError {
      break
    } else if err != nil {
      t.Errorf("Unexpected error %+v", err)
    }
  }
  data = bfm.OutputBuffer.Next(1024)
  if string(data) != "Hello World!\n" {
    t.Errorf("Hello World test did not work [got %+v]\n%s\n", data, bfm.Info())
    }
  bfm.LoadProgram(quine)
  for {
    err = bfm.Step()
    if err == EndOfProgramError {
      break
    } else if err != nil {
      t.Errorf("Unexpected error %+v", err)
    }
  }
  data = bfm.OutputBuffer.Next(2048)
  if string(data) != quine {
    t.Errorf("Quine not a quine.\nExpected: %s\n-----Got: %s", quine, data)
  }
}

func makeRange(min, max int) []int {
    a := make([]int, max-min+1)
    for i := range a {
        a[i] = min + i
    }
    return a
}
