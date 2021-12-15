package synthizer

/*
#include <synthizer.h>
#include <synthizer_constants.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
	"errors"
)

func GOCHECK(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func CHECKED(x C.syz_ErrorCode) error {
	var ret int = x
	if ret {
		return errors.New("SynthizerError: ErrorCode " + ret + ", " + C.syz_getLastErrorMessage())
	}
	return nil
}

func Initialize() error {
	err := CHECKED(C.syz_initialize())
	return GOCHECK(err)
}

func Shutdown() error {
	err := CHECKED(C.syz_shutdown())
	return GOCHECK(err)
}