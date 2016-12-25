package errors
import (
	"errors"
	"fmt"
	"runtime"
	"bytes"
)
type ErrCoder interface {
	GetErrCode() ErrCode
}

type ErrCode int32

const (
	ErrNoCode               ErrCode = -2
	ErrNoError              ErrCode = 0
	ErrUnknown              ErrCode = -1
	ErrDuplicatedTx         ErrCode = 1
	ErrDuplicateInput       ErrCode = 45003
	ErrAssetPrecision       ErrCode = 45004
	ErrTransactionBalance   ErrCode = 45005
	ErrAttributeProgram     ErrCode = 45006
	ErrTransactionContracts ErrCode = 45007
	ErrTransactionPayload   ErrCode = 45008
	ErrDoubleSpend          ErrCode = 45009
	ErrTxHashDuplicate      ErrCode = 45010
	ErrStateUpdaterVaild    ErrCode = 45011
	ErrSummaryAsset         ErrCode = 45012
	ErrLockedAsset          ErrCode = 45013
	ErrDuplicateCharacter   ErrCode = 45014
	ErrInvalidOutput        ErrCode = 45015
	ErrXmitFail             ErrCode = 45016
)

type CallStacker interface {
	GetCallStack() *CallStack
}

type CallStack struct {
	Stacks []uintptr
}

func GetCallStacks(err error) *CallStack {
	if err, ok := err.(CallStacker); ok {
		return err.GetCallStack()
	}
	return nil
}

func CallStacksString(call *CallStack) string {
	buf := bytes.Buffer{}
	if call == nil {
		return fmt.Sprintf("No call stack available")
	}

	for _, stack := range call.Stacks {
		f := runtime.FuncForPC(stack)
		file, line := f.FileLine(stack)
		buf.WriteString(fmt.Sprintf("%s:%d - %s\n", file, line, f.Name()))
	}

	return fmt.Sprintf("%s", buf.Bytes())
}

func getCallStack(skip int, depth int) *CallStack {
	stacks := make([]uintptr, depth)
	stacklen := runtime.Callers(skip, stacks)

	return &CallStack{
		Stacks: stacks[:stacklen],
	}
}

const callStackDepth = 10

type DetailError interface {
	error
	ErrCoder
	CallStacker
	GetRoot() error
}

func NewErr(errmsg string) error {
	return errors.New(errmsg)
}

func NewDetailErr(err error, errcode ErrCode, errmsg string) DetailError {
	if err == nil {
		return nil
	}
	return nil
}

func RootErr(err error) error {
	if err, ok := err.(DetailError); ok {
		return err.GetRoot()
	}
	return err
}