package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type instructionType uint8

const (
	turn_on instructionType = iota
	turn_off
	toggle
)

type instruction struct {
	instruction_type instructionType
	start_i          int
	start_j          int
	end_i            int
	end_j            int
}

func loadInput(input_file_name string) ([]instruction, error) {
	input_file, err := os.Open(input_file_name)
	if err != nil {
		return nil, err
	}

	input_bytes, err := io.ReadAll(input_file)
	if err != nil {
		return nil, err
	}

	input := string(input_bytes)

	var instructions []instruction
	for line := range strings.Lines(input) {
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		if len(fields) < 4 ||
			len(fields) > 5 ||
			fields[0] == "toggle" && len(fields) != 4 ||
			fields[0] != "toggle" && len(fields) != 5 {
			return instructions, fmt.Errorf("%#v contains the wrong number of words", line)
		}

		var through_index int
		var start_pair, end_pair string
		var instruction_type instructionType
		switch fields[0] {
		case "turn":
			through_index = 3
			start_pair = fields[2]
			end_pair = fields[4]

			switch fields[1] {
			case "on":
				instruction_type = turn_on
			case "off":
				instruction_type = turn_off
			default:
				return instructions, fmt.Errorf(`%#v starts with "turn", but it is not followed by "on" or "off"`, line)
			}
		case "toggle":
			through_index = 2
			start_pair = fields[1]
			end_pair = fields[3]
			instruction_type = toggle
		default:
			return instructions, fmt.Errorf(`%#v does not start with "turn" or "toggle"`, line)
		}

		if fields[through_index] != "through" {
			var through_ordinal string
			if through_index == 3 {
				through_ordinal = "fourth"
			} else {
				through_ordinal = "third"
			}
			return instructions, fmt.Errorf(`The %v word of %#v should be "through", not %#v`, through_ordinal, line, fields[through_index])
		}

		start_i, start_j, err := parsePair(start_pair)
		if err != nil {
			return instructions, err
		}

		end_i, end_j, err := parsePair(end_pair)
		if err != nil {
			return instructions, err
		}

		instructions = append(instructions, instruction{instruction_type, start_i, start_j, end_i, end_j})
	}

	return instructions, nil
}

func parsePair(pair_string string) (int, int, error) {
	numbers := strings.Split(pair_string, ",")
	if len(numbers) != 2 {
		return -1, -1, fmt.Errorf(`%#v should contain exactly one ","`, pair_string)
	}

	i, err := strconv.Atoi(numbers[0])
	if err != nil {
		return -1, -1, err
	}

	j, err := strconv.Atoi(numbers[1])
	if err != nil {
		return -1, -1, err
	}

	return i, j, nil
}

func part1(instructions []instruction) uint {
	var lights [1000][1000]bool

	for _, instruction := range instructions {
		for i := instruction.start_i; i <= instruction.end_i; i++ {
			for j := instruction.start_j; j <= instruction.end_j; j++ {
				// There's a readablility-performance tradeoff here (maybe).
				// Putting the switch inside the j for loop checks instruction.instruction_type a bunch of extra times
				//  (assuming there isn't an optimization for this type of thing).
				// Putting the switch outsied the i for loop only checks once per instruction, but requires either
				//  duplicating the inner for loops or creating a function variable, which hurts readability.
				// Testing is needed.
				switch instruction.instruction_type {
				case turn_on:
					lights[i][j] = true
				case turn_off:
					lights[i][j] = false
				case toggle:
					lights[i][j] = !lights[i][j]
				}
			}
		}
	}

	var lit_count uint
	for i := range 1000 {
		for j := range 1000 {
			if lights[i][j] {
				lit_count++
			}
		}
	}

	return lit_count
}

func part2(instructions []instruction) uint {
	var lights [1000][1000]uint

	for _, instruction := range instructions {
		for i := instruction.start_i; i <= instruction.end_i; i++ {
			for j := instruction.start_j; j <= instruction.end_j; j++ {
				// There's a readablility-performance tradeoff here (maybe).
				// Putting the switch inside the j for loop checks instruction.instruction_type a bunch of extra times
				//  (assuming there isn't an optimization for this type of thing).
				// Putting the switch outsied the i for loop only checks once per instruction, but requires either
				//  duplicating the inner for loops or creating a function variable, which hurts readability.
				// Testing is needed.
				switch instruction.instruction_type {
				case turn_on:
					lights[i][j]++
				case turn_off:
					if lights[i][j] > 0 {
						lights[i][j]--
					}
				case toggle:
					lights[i][j] += 2
				}
			}
		}
	}

	var total_brightness uint
	for i := range 1000 {
		for j := range 1000 {
			total_brightness += lights[i][j]
		}
	}

	return total_brightness
}

func main() {
	instructions, err := loadInput("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	answer := part1(instructions)
	time_ := time.Since(start)
	fmt.Printf("Part 1: %v (%v)\n", answer, time_)

	start = time.Now()
	answer = part2(instructions)
	time_ = time.Since(start)
	fmt.Printf("Part 2: %v (%v)\n", answer, time_)
}
