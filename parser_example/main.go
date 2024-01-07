package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input := "$8\r\nCommand1\r\n"
	reader := bufio.NewReader(strings.NewReader(input))

	b, _ := reader.ReadByte()

	if b != '$' {
		fmt.Println("Invalid type, expecting bulk strings only")
		os.Exit(1)
	}

	size, _ := reader.ReadByte()

	strSize, _ := strconv.ParseInt(string(size), 10, 64)

	// consume /r/n
	consumeRN(reader)

	name := make([]byte, strSize)
	reader.Read(name)

	fmt.Println(string(name))
}

func consumeRN(reader *bufio.Reader) {
	reader.ReadByte()
	reader.ReadByte()
}
