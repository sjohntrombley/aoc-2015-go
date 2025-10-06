package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

func part1(input string) uint {
	var answer uint = 1
	for {
		test_bytes := fmt.Appendf(nil, "%v%v", input, answer)
		hash_bytes := md5.Sum(test_bytes)
		if hash_bytes[0] == 0 && hash_bytes[1] == 0 && hash_bytes[2] & 0xF0 == 0 {
			break
		}
		answer++
	}
	return answer
}

func part2(input string) uint {
	var answer uint = 1
	for {
		test_bytes := fmt.Appendf(nil, "%v%v", input, answer)
		hash_bytes := md5.Sum(test_bytes)
		if hash_bytes[0] == 0 && hash_bytes[1] == 0 && hash_bytes[2] == 0 {
			break
		}
		answer++
	}
	return answer
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
	input := strings.TrimFunc(string(input_bytes), unicode.IsSpace)

	start := time.Now()
	part1_answer := part1(input)
	part1_time := time.Since(start)
	fmt.Printf("Part 1: %v (%v)\n", part1_answer, part1_time)
	
	start = time.Now()
	part2_answer := part2(input)
	part2_time := time.Since(start)
	fmt.Printf("Part 2: %v (%v)\n", part2_answer, part2_time)
}
