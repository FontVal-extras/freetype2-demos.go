package main

// Copyright (c) 2022 Hin-Tak Leung

// Diagnostics with FontVal's backend.
//     - Approx re-write of examples/font-diag.py in the sister freetype-py project.

// Notes: Go have functional comments... The below lines are actually used by cgo.

// #cgo pkg-config: freetype2
// #include <ft2build.h>
// #include FT_FREETYPE_H
// #include <freetype/ftmodapi.h>
// #include <freetype/tttables.h>
/*
   // This forward declaration is necessary.
   int go_diagfunc(FT_Face face, int messcode, char *message,
        char *opcode,
        int range_base,
        int is_composite,
        int IP,
        int callTop,
        int opc,
        int start);
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"
)

func set_and_check_interpreter_version(lib C.FT_Library, x int) {
	version := C.FT_Int(x)
	C.FT_Property_Set(lib, C.CString("truetype"), C.CString("interpreter-version"), unsafe.Pointer(&version))
	fmt.Println("Try setting truetype interpreter-version:", version)

	got_version := C.FT_Int(0)
	C.FT_Property_Get(lib, C.CString("truetype"), C.CString("interpreter-version"), unsafe.Pointer(&got_version))
	fmt.Println("Getting truetype interpreter-version:", got_version)
}

// In the following line: No space before "export"! Important! Cannot append comments either!

//export go_diagfunc
func go_diagfunc(face C.FT_Face, messcode C.int, message *C.char,
	opcode *C.char,
	range_base C.int,
	is_composite C.int,
	IP C.int,
	callTop C.int,
	opc C.int,
	start C.int) C.int {
	fmt.Println(messcode, C.GoString(message), C.GoString(opcode), range_base, is_composite, IP, callTop, opc, start)
	// The below shows nothing - requires FT_CONFIG_OPTION_ERROR_STRINGS to work.
	// Kept here for reference.
	fmt.Println(C.GoString(C.FT_Error_String(messcode)))
	return 0
}

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, os.Args[0], "does not take arguments.")
		os.Exit(1)
		return
	}

	var lib C.FT_Library
	if err := C.FT_Init_FreeType(&lib); err != 0 {
		fmt.Fprintln(os.Stderr, "unable to init freetype")
		// FT_Error needs int() cast.
		os.Exit(int(err))
		return
	}
	defer C.FT_Done_FreeType(lib)

	cs := C.CString(os.Args[1])
	defer C.free(unsafe.Pointer(cs))

	var face C.FT_Face
	if err := C.FT_New_Face(lib, cs, 0, &face); err != 0 {
		// FT_Error needs int() cast.
		os.Exit(int(err))
		return
	}
	defer C.FT_Done_Face(face)

	fmt.Println("family:", C.GoString(face.family_name))
	cmaxp := C.FT_Get_Sfnt_Table(face, C.FT_SFNT_MAXP)
	maxp := (*C.TT_MaxProfile)(unsafe.Pointer(cmaxp))
	fmt.Println(maxp.numGlyphs)
	set_and_check_interpreter_version(lib, 40)
	set_and_check_interpreter_version(lib, 38)
	set_and_check_interpreter_version(lib, 60)
	set_and_check_interpreter_version(lib, 35)

	// int is not FT_F26Dot6
	size := C.FT_F26Dot6(10)
	C.FT_Set_Char_Size(face, size*64, 0, 96, 96)
	lf := C.FT_LOAD_DEFAULT | C.FT_LOAD_NO_AUTOHINT | C.FT_LOAD_MONOCHROME | C.FT_LOAD_COMPUTE_METRICS
	lf |= C.FT_LOAD_TARGET_MONO
	// maxp.numGlyphs is not GO int.
	for ig := 1; ig < int(maxp.numGlyphs); ig++ {
		C.TT_Diagnostics_Set(face, (C.FT_DiagnosticsFunc)(unsafe.Pointer(C.go_diagfunc)))
		C.FT_Load_Glyph(face, C.FT_UInt(ig), C.FT_Int32(lf))
		C.TT_Diagnostics_Unset(face)
	}
}
