package main

import (
	""
	"fmt"
)

func main() {
	client := lib.Init(60)
	Result := lib.GetSong("forget not", "ne obliviscaris", client)
	fmt.Println(Result)
}
