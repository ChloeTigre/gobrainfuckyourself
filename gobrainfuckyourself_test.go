package gobrainfuckyourself
import "testing"
func TestOperators(t *testing.T) {
  var a bfOperator
  var err error
  if a, err = getCommandCode('>');a != OPERATOR_RIGHT {
    t.Error("Wrong value for >")
  }
  if a, err = getCommandCode('<');a != OPERATOR_LEFT {
    t.Error("Wrong value for <")
  }
  if a, err = getCommandCode('-');a != OPERATOR_DEC {
    t.Error("Wrong value for -")
  }
  if a, err = getCommandCode('+');a != OPERATOR_INC {
    t.Error("Wrong value for +")
  }
  if a, err = getCommandCode('.');a != OPERATOR_OUTPUT {
    t.Error("Wrong value for .")
  }
  if a, err = getCommandCode(',');a != OPERATOR_INPUT {
    t.Error("Wrong value for ,")
  }
  if a, err = getCommandCode('[');a != OPERATOR_LOOP_TOP {
    t.Error("Wrong value for [")
  }
  if a, err = getCommandCode(']');a != OPERATOR_LOOP_BOTTOM {
    t.Error("Wrong value for ]")
  }
  if err != nil {
    t.Error("There should not be errors")
  }
}
