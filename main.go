package main

import (
	"SlowReverb/lib"
	"fmt"
)

func main() {
	client := lib.Init(60)
	Result := lib.GetSong("Serpentskirt", "Cocteau Twins", client)
	fmt.Println(*Result.Filename)
	file := lib.ModifySpeed(*Result.Filename, 0.85)

	file = lib.Reverberize(file, 10, 10, 12)
	file = lib.Reverberize(file, 10, 5, 14)

	lib.Play(file)
}
