package main

import "errors"
import "fmt"
import "io"
import "log"
import "os"

func part1(input string) (int, error) {
	floor := 0
	for _, ch := range input {
		switch ch {
		case '(':
			floor++
		case ')':
			floor--
		default:
			return floor, fmt.Errorf("Expected '(' or ')', not %#v.", ch)
		}
	}

	return floor, nil
}

func part2(input string) (int, error) {
	floor := 0
	for i, ch := range input {
		position := i + 1
		switch ch {
		case '(':
			floor++
		case ')':
			floor--
			if floor < 0 {
				return position, nil
			}
		default:
			return position, fmt.Errorf("Expected '(' or ')', not %#v.", ch)
		}
	}

	return 0, errors.New("These instructions never enter the basement.")
}

func main() {
	input_file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	input_bytes, err := io.ReadAll(input_file)
	if err != nil {
		log.Fatal(err)
	}
	input_string := string(input_bytes)

	part1_answer, err := part1(input_string)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 1: %v\n", part1_answer)

	part2_answer, err := part2(input_string)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Part 2: %v\n", part2_answer)
}
