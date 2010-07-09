package main

import(
	"os"
	"fmt"
	"path"
	"strconv"
	"image/jpeg"
	"image/png"
)

func usageE(e string) {
	fmt.Printf("Error: %s\n", e)
	fmt.Printf("---------------------------\n\n")

	usage()
}

func usage() {
	fmt.Printf("\033[0;1mUsage:\033[m\n")
	fmt.Printf("\t%s <in> <out> <part-size>\n\n", os.Args[0])
	fmt.Printf("\t\tin: A Jpeg or PNG file.\n")
	fmt.Printf("\t\tout: A Jpeg or PNG file.\n")
	fmt.Printf("\t\tpart-size: The width of each section of the SIRDS.\n")
}

func main() {
	if len(os.Args) != 4 {
		usage()
		os.Exit(1)
	}

	inN := os.Args[1]
	outN := os.Args[2]
	partSize, _ := strconv.Atoi(os.Args[3])

	inR, err := os.Open(inN)
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}

	switch path.Ext(inN) {
		case "jpg":
			in, _ := jpeg.Decode(inR)
		case "png":
			in, _ := png.Decode(inR)
		default:
			usageE("Format not supported...")
			os.Exit(1)
	}
}
