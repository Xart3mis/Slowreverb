package main

import (
	"SlowReverb/lib"
	"fmt"
)

func main() {
	client := lib.Init(60)
	Result := lib.GetSong("Disgusted With Myself", "Negative XP", client)
	fmt.Println(*Result.Filename)
	lib.ModifySpeed(*Result.Filename, 0.86)
}
