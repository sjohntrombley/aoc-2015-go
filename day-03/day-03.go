package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type coordinate struct {
	x int
	y int
}

func part1(input string) (int, error) {
	x, y := 0, 0
	visited := map[coordinate]struct{}{coordinate{x, y}: struct{}{}}
	for _, ch := range input {
		switch ch {
		case '<':
			x--
		case '>':
			x++
		case 'v':
			y--
		case '^':
			y++
		default:
			return len(visited), fmt.Errorf("Unexpected character in input. Expected '<', '>', 'v', or '^', not %#v", ch)
		}

		coord := coordinate{x, y}
		if _, ok := visited[coord]; !ok {
			visited[coord] = struct{}{}
		}
	}

	return len(visited), nil
}

func part2(input string) (int, error) {
	x, y, rx, ry := 0, 0, 0, 0
	r_turn := false
	visited := map[coordinate]struct{}{coordinate{x, y}: struct{}{}}
	for _, ch := range input {
		switch ch {
		case '<':
			if r_turn {
				rx--
			} else {
				x--
			}
		case '>':
			if r_turn {
				rx++
			} else {
				x++
			}
		case 'v':
			if r_turn {
				ry--
			} else {
				y--
			}
		case '^':
			if r_turn {
				ry++
			} else {
				y++
			}
		default:
			return len(visited),
				fmt.Errorf(
					"Unexpected character in input. Expected '<', '>', 'v', "+
						"or '^', not %#v",
					ch,
				)
		}

		var coord coordinate
		if r_turn {
			coord = coordinate{rx, ry}
		} else {
			coord = coordinate{x, y}
		}
		if _, ok := visited[coord]; !ok {
			visited[coord] = struct{}{}
		}

		r_turn = !r_turn
	}

	return len(visited), nil
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
	input := string(input_bytes)

	part1_answer, err := part1(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 1:", part1_answer)

	part2_answer, err := part2(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Part 2:", part2_answer)
}
