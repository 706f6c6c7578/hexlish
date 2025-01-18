package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var hexToHexlish = map[rune]rune{
	'0': 'A', '1': 'C', '2': 'E', '3': 'H',
	'4': 'I', '5': 'J', '6': 'L', '7': 'M',
	'8': 'N', '9': 'O', 'A': 'P', 'B': 'R',
	'C': 'S', 'D': 'T', 'E': 'U', 'F': 'V',
	'a': 'P', 'b': 'R', 'c': 'S', 'd': 'T',
	'e': 'U', 'f': 'V',
}

var hexlishToHex = map[rune]rune{
	'A': '0', 'C': '1', 'E': '2', 'H': '3',
	'I': '4', 'J': '5', 'L': '6', 'M': '7',
	'N': '8', 'O': '9', 'P': 'A', 'R': 'B',
	'S': 'C', 'T': 'D', 'U': 'E', 'V': 'F',
}

const usage = `Usage: hexlish [-d] < input_file > output_file

Hexlish is a tool to convert between hexadecimal and Hexlish encoding.

Options:
  -d    decode mode (convert Hexlish to hexadecimal)
  -h    show this help message

Examples:
  echo "DEADBEEF" | hexlish
  echo "TUPTRUUV" | hexlish -d
  cat large_hex_file.txt | hexlish > encoded.txt
  cat encoded.txt | hexlish -d > decoded.txt

Hexlish Alphabet Mapping:
  Hexlish:  A C E H I J L M N O P R S T U V
  Hex:      0 1 2 3 4 5 6 7 8 9 A B C D E F`

func encode(input string) string {
	var result strings.Builder
	result.Grow(len(input))
	for _, ch := range input {
		if val, ok := hexToHexlish[ch]; ok {
			result.WriteRune(val)
		} else if !strings.ContainsRune(" \t\n\r", ch) {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func decode(input string) string {
	var result strings.Builder
	result.Grow(len(input))
	for _, ch := range input {
		if val, ok := hexlishToHex[ch]; ok {
			result.WriteRune(val)
		} else if !strings.ContainsRune(" \t\n\r", ch) {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func processStream(reader io.Reader, writer io.Writer, decodeMode bool) error {
    scanner := bufio.NewScanner(reader)
    scanner.Buffer(make([]byte, 64*1024), 1024*1024)

    for scanner.Scan() {
        input := scanner.Text()
        var output string
        if decodeMode {
            output = decode(input)
        } else {
            output = encode(input)
        }
        if _, err := fmt.Fprintln(writer, output); err != nil {
            return fmt.Errorf("error writing output: %v", err)
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("error reading input: %v", err)
    }
    return nil
}

func main() {
	decodeFlag := flag.Bool("d", false, "decode mode")
	helpFlag := flag.Bool("h", false, "show help message")
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", usage)
	}
	
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "Error: This program only accepts input from stdin\n\n")
		flag.Usage()
		os.Exit(1)
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintf(os.Stderr, "Error: No input provided\n\n")
		flag.Usage()
		os.Exit(1)
	}

	if err := processStream(os.Stdin, os.Stdout, *decodeFlag); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
