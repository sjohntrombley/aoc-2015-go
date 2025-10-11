package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

func part1(input string) uint {
	var nice_count uint = 0
	for line := range strings.Lines(input) {
		r := strings.NewReader(strings.TrimFunc(line, unicode.IsSpace))

		vowel_count := 0
		double_found := false
		illegal_substring := false

		prev_c, _, err := r.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if is_vowel(prev_c) {
			vowel_count++
		}
		for {
			c, _, err := r.ReadRune()
			if err != nil {
				break
			}

			if is_vowel(c) {
				vowel_count++
			}

			if c == prev_c {
				double_found = true
			}

			is_ab := prev_c == 'a' && c == 'b'
			is_cd := prev_c == 'c' && c == 'd'
			is_pq := prev_c == 'p' && c == 'q'
			is_xy := prev_c == 'x' && c == 'y'
			if is_ab || is_cd || is_pq || is_xy {
				illegal_substring = true
				break
			}

			prev_c = c
		}

		if !illegal_substring && vowel_count >= 3 && double_found {
			nice_count += 1
		}
	}

	return nice_count
}

func is_vowel(c rune) bool {
	switch c {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	default:
		return false
	}
}

type runePair struct {
	i int
	r rune
}

func part2(input string) uint {
	var nice_count uint = 0

	for line := range strings.Lines(input) {
		line = strings.TrimFunc(line, unicode.IsSpace)
		r := strings.NewReader(line)

		i := 0
		// TODO: storing the indexes here is better for checking the first rule
		//	but runes are better for the second rule, so maybe we make a struct
		//  and store both
		var prev_cs [2]runePair
		c, n, err := r.ReadRune()
		if err != nil {
			log.Fatalf("Failed to read first character of %#v (%v).", line, err)
		}
		prev_cs[0] = runePair{i, c}
		i += n
		c, n, err = r.ReadRune()
		if err != nil {
			log.Fatalf("Failed to read second character of %#v (%v).", line, err)
		}
		prev_cs[1] = runePair{i, c}
		i += n

		var double_pair bool
		if strings.Contains(line[i:], line[:i]) {
			double_pair = true
		} else {
			double_pair = false
		}
		repeated_with_gap := false
		for {
			c, n, err = r.ReadRune()
			if err != nil {
				break
			}

			if !double_pair && strings.Contains(line[i + n:], line[prev_cs[1].i:i + n]) {
				double_pair = true
			}

			if prev_cs[0].r == c {
				repeated_with_gap = true
			}

			prev_cs[0] = prev_cs[1]
			prev_cs[1] = runePair{i, c}
			i += n

			if double_pair && repeated_with_gap {
				break
			}
		}

		if double_pair && repeated_with_gap {
			nice_count++
		}
	}

	return nice_count
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

	start := time.Now()
	answer := part1(input)
	run_time := time.Since(start)
	fmt.Printf("Part 1: %v (%v)\n", answer, run_time)
	start = time.Now()
	answer = part2(input)
	run_time = time.Since(start)
	fmt.Printf("Part 2: %v (%v)\n", answer, run_time)
}
