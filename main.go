package main

import (
	"SlowReverb/lib"
	"fmt"
)

func main() {
	client := lib.Init(60)
	Result := lib.GetSong("Serpentskirt", "Cocteau Twins", client)
	fmt.Println(*Result.Filename)
	file := lib.ModifySpeed(*Result.Filename, 0.9)

	file = lib.Reverberize(file, 10, 10, 12, lib.ReverbTypes().Hall.Large_Hall)
	file = lib.Reverberize(file, 10, 5, 14, lib.ReverbTypes().Hall.Medium_Hall)

	lib.Play(file)
}
