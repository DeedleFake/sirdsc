package main

import(
	"os"
	"fmt"
	"strings"
	"strconv"
)

func usage() {
	fmt.Printf("\033[0;1mUsage:\033[m\n")
	fmt.Printf("\t%s <in> <out> <part-size>\n", os.Args[0])
}

func main() {
	if len(os.Args) != 4 {
		usage()
		os.Exit(1)
	}

	inN := os.Args[1]
	outN := os.Args[2]
	partSize, _ := strconv.Atoi(os.Args[3])

	switch (strings.ToLower(inN[len(inN)-3:])) {
		case "jpg":
			fmt.Printf("Jpeg.\n")
		case "png":
			fmt.Printf("PNG.\n")
	}

	fmt.Printf(outN, partSize)
}
