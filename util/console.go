package util

import (
	"fmt"
	"os"
)

func AnyKeyToQuit() {
	fmt.Printf("Press any key to exit...")
	b := make([]byte, 1)
	_, _ = os.Stdin.Read(b)
}
