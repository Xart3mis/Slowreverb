package main

import (
	"SlowReverb/lib"
	"fmt"
)

func main() {
	client := lib.Init(60)
	fmt.Println("Downloading Song:")
	Result := lib.GetSong("4 morant", "doja cat", client)
	fmt.Println("Done.")
	// file := *Result.Filename

	fmt.Println("Slowing down audio:")
	file := lib.ModifySpeed(*Result.Filename, 0.75)
	fmt.Println("Done.")

	fmt.Println("\nApplying reverb 1:")
	file = lib.Reverberize(file, 10, 10, 12, lib.ReverbTypes().Hall.Large_Hall)
	fmt.Println("Done.")

	fmt.Println("\nApplying reverb 2:")
	file = lib.Reverberize(file, 10, 5, 10, lib.ReverbTypes().Chamber.Vocal_Chamber)
	fmt.Println("Done.")

	fmt.Println("\nPitching down:")
	file = lib.AlterPitch(file, 0.2, false)
	fmt.Println("Done.")

	fmt.Println("Playing song:")
	finished_playing := make(chan bool)
	go lib.Play(file, finished_playing, false)

	<-finished_playing
	fmt.Println("Done.")
}
