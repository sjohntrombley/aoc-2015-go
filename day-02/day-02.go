package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func load_input(path string) ([][3]uint, error) {
	input_file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	input_bytes, err := io.ReadAll(input_file)
	if err != nil {
		log.Fatal(err)
	}
	input_string := string(input_bytes)

	var output [][3]uint = nil
	var cur_dims [3]uint
	for line := range strings.Lines(input_string) {
		dim_strings := strings.Split(strings.TrimFunc(line, unicode.IsSpace), "x")
		if len(dim_strings) != 3 {
			return output, fmt.Errorf(
				"A box should have exactly 3 dimensions, but %#v has %v.",
				line,
				len(line))
		}

		for i, dim := range dim_strings {
			dim, err := strconv.ParseUint(dim, 10, strconv.IntSize)
			if err != nil {
				return output, err
			}
			cur_dims[i] = uint(dim)
		}

		slices.Sort(cur_dims[:])
		output = append(output, cur_dims)
	}

	return output, nil
}

func part1(input [][3]uint) uint {
	var total_area uint = 0
	for _, dims := range input {
		a := dims[0]
		b := dims[1]
		c := dims[2]
		total_area += 3*a*b + 2*a*c + 2*b*c
	}
	return total_area
}

func part2(input [][3]uint) uint {
	var total_length uint = 0
	for _, dims := range input {
		a := dims[0]
		b := dims[1]
		c := dims[2]
		total_length += 2*a + 2*b + a*b*c
	}
	return total_length
}

func main() {
	input, err := load_input("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Part 1: %v\n", part1(input))
	fmt.Printf("Part 2: %v\n", part2(input))
}
