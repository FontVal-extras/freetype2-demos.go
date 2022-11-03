package main

// Copyright (c) 2022 Hin-Tak Leung

// Basic FreeType test - just print Version.

// Notes: Go have functional comments... The below lines are actually used by cgo.

// #cgo pkg-config: freetype2
// #include <ft2build.h>
// #include FT_FREETYPE_H
import "C"
import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, os.Args[0], "does not take arguments.")
		os.Exit(1)
		return
	}

	var lib C.FT_Library
	if err := C.FT_Init_FreeType(&lib); err != 0 {
		fmt.Fprintln(os.Stderr, "unable to init freetype")
		os.Exit(int(err))
		return
	}

	defer C.FT_Done_FreeType(lib)

	var major, minor, patch C.FT_Int
	C.FT_Library_Version(lib, &major, &minor, &patch)
	// cast to int() as FT_Int is not necessarily Go's int.
	// Go's int is 64-bit on 64-bit platforms.
	fmt.Printf("OK (%d.%d.%d)\n", int(major), int(minor), int(patch))
}	
