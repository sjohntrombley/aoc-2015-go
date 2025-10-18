package main

import (
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type wires map[string]uint16

type gate interface {
	process(ws wires) bool
	getOutput() string
}

type binaryGateType uint8

const (
	andGate binaryGateType = iota
	orGate
)

type binaryGate struct {
	gateType binaryGateType
	input1   string
	input2   string
	output   string
}

func (g *binaryGate) process(ws wires) bool {
	if _, ok := ws[g.output]; ok {
		return false
	}

	input1, ok := ws[g.input1]
	if !ok {
		return false
	}

	input2, ok := ws[g.input2]
	if !ok {
		return false
	}

	switch g.gateType {
	case andGate:
		ws[g.output] = input1 & input2
	case orGate:
		ws[g.output] = input1 | input2
	default:
		panic("bad gateType")
	}
	return true
}

func (g *binaryGate) getOutput() string {
	return g.output
}

type shiftGateType uint8

const (
	lshiftGate shiftGateType = iota
	rshiftGate
	// Not a shift but it has the same layout
	andConstGate
)

type shiftGate struct {
	gateType shiftGateType
	input    string
	amount   uint16
	output   string
}

func (g *shiftGate) process(ws wires) bool {
	input, ok := ws[g.input]
	if !ok {
		return false
	}

	switch g.gateType {
	case lshiftGate:
		ws[g.output] = input << g.amount
	case rshiftGate:
		ws[g.output] = input >> g.amount
	case andConstGate:
		ws[g.output] = input & g.amount
	default:
		panic("bad gateType")
	}
	return true
}

func (g *shiftGate) getOutput() string {
	return g.output
}

type unaryGateType uint8

const (
	notGate unaryGateType = iota
	// Treating wire-to-wire connections as a gate makes the logic easier
	repeatGate
)

type unaryGate struct {
	gateType unaryGateType
	input    string
	output   string
}

func (g *unaryGate) process(ws wires) bool {
	input, ok := ws[g.input]
	if !ok {
		return false
	}

	switch g.gateType {
	case notGate:
		ws[g.output] = ^input
	case repeatGate:
		ws[g.output] = input
	default:
		panic("bad gateType")
	}
	return true
}

func (g *unaryGate) getOutput() string {
	return g.output
}

type gateMap map[string][]gate
type initialWireMap map[string]uint16

func invalidInstructionError(line string) error {
	return fmt.Errorf("%q is an invalid instruction", line)
}

func parseInput(input string) (gateMap, initialWireMap, error) {
	gates := gateMap{}
	initial_wires := initialWireMap{}
	for line := range strings.Lines(input) {
		line = strings.TrimSpace(line)
		fields := strings.Fields(line)
		if len(fields) < 3 {
			return nil, nil, invalidInstructionError(line)
		}
		if fields[0] == "NOT" {
			if len(fields) != 4 {
				return nil, nil, invalidInstructionError(line)
			}

			gate_ := &unaryGate{
				gateType: notGate,
				input:    fields[1],
				output:   fields[3],
			}
			if gs, ok := gates[gate_.input]; ok {
				gates[gate_.input] = append(gs, gate_)
			} else {
				gates[gate_.input] = []gate{gate_}
			}
		} else if n, err := strconv.ParseUint(fields[0], 10, 16); err == nil {
			if fields[1] == "AND" {
				if len(fields) != 5 {
					return nil, nil, invalidInstructionError(line)
				}

				gate_ := &shiftGate{
					gateType: andConstGate,
					input:    fields[2],
					amount:   uint16(n),
					output:   fields[4],
				}
				if gs, ok := gates[gate_.input]; ok {
					gates[gate_.input] = append(gs, gate_)
				} else {
					gates[gate_.input] = []gate{gate_}
				}
			} else {
				initial_wires[fields[2]] = uint16(n)
			}
		} else if fields[1] == "->" {
			gate_ := &unaryGate{
				gateType: repeatGate,
				input:    fields[0],
				output:   fields[2],
			}
			if gs, ok := gates[gate_.input]; ok {
				gates[gate_.input] = append(gs, gate_)
			} else {
				gates[gate_.input] = []gate{gate_}
			}
		} else if len(fields) < 5 {
			return nil, nil, invalidInstructionError(line)
		} else {
			switch fields[1] {
			case "AND", "OR":
				gate_ := &binaryGate{
					input1: fields[0],
					input2: fields[2],
					output: fields[4],
				}
				if fields[1] == "AND" {
					gate_.gateType = andGate
				} else {
					gate_.gateType = orGate
				}

				if gs, ok := gates[gate_.input1]; ok {
					gates[gate_.input1] = append(gs, gate_)
				} else {
					gates[gate_.input1] = []gate{gate_}
				}
				if gs, ok := gates[gate_.input2]; ok {
					gates[gate_.input2] = append(gs, gate_)
				} else {
					gates[gate_.input2] = []gate{gate_}
				}
			case "LSHIFT", "RSHIFT":
				amount, err := strconv.ParseUint(fields[2], 10, 16)
				if err != nil {
					return nil, nil, invalidInstructionError(line)
				}
				gate_ := &shiftGate{
					input:  fields[0],
					amount: uint16(amount),
					output: fields[4],
				}
				if fields[1] == "LSHIFT" {
					gate_.gateType = lshiftGate
				} else {
					gate_.gateType = rshiftGate
				}

				if gs, ok := gates[gate_.input]; ok {
					gates[gate_.input] = append(gs, gate_)
				} else {
					gates[gate_.input] = []gate{gate_}
				}
			default:
				return nil, nil, invalidInstructionError(line)
			}
		}
	}

	return gates, initial_wires, nil
}

func part1(gates gateMap, initial_wires initialWireMap) uint16 {
	wire_values := wires{}
	wire_update_queue := list.New()
	for wire, value := range initial_wires {
		wire_values[wire] = value
		wire_update_queue.PushBack(wire)
	}

	for wire_update_queue.Len() > 0 {
		wire := wire_update_queue.Remove(wire_update_queue.Front()).(string)
		for _, gate_ := range gates[wire] {
			if gate_.process(wire_values) {
				wire_update_queue.PushBack(gate_.getOutput())
			}
		}
	}

	return wire_values["a"]
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
	start := time.Now()
	gates, initial_wires, err := parseInput(string(input_bytes))
	time_ := time.Since(start)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parse time: %v\n", time_)

	start = time.Now()
	answer := part1(gates, initial_wires)
	time_ = time.Since(start)
	fmt.Printf("Part 1: %v (%v)\n", answer, time_)

	start = time.Now()
	initial_wires["b"] = answer
	answer = part1(gates, initial_wires)
	time_ = time.Since(start)
	fmt.Printf("Part 2: %v (%v)\n", answer, time_)
}
