package main

import (
	"fmt"
	"github.com/jingkang99/ztam-amaas/pkg/global"
	"github.com/jingkang99/ztam-amaas/cmd/devop/cmd"
)

const N = 52

func main() {
	cmd.Execute()

	if(global.Debug){

	// Example: this will give us a 44 byte, base64 encoded output
	token, err := global.GenerateRandomStringURLSafe(32)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)

	// Example: this will give us a 32 byte output
	token, err = global.GenerateRandomString(32)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)

	fmt.Printf("%v", global.Shuffle(52))
	}
}
