package helpers

import (
	"fmt"
	"runtime"
)

func RecoveryTracePrint(recoverErr interface{}) {

	if recoverErr != nil {

		fmt.Println("recoverErr::", recoverErr)

		stack := make([]byte, 1<<16)
		runtime.Stack(stack, false)
		fmt.Println("stack::", string(stack))
	}

}
