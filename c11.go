package main

import (
	"fmt"
	"strings"
)

/*
Graph problem: Find all paths fron a to b
*/

func runC11() {
	deviceMap := readFile("c11.txt", parseDevices)
	//printAllDevices(deviceMap)

	execC11(deviceMap)
	execC11Two(deviceMap)
}

func execC11(deviceMap map[string]*device_S) {
	var count int
	search(deviceMap["you"], &count, make(map[string]bool))
	fmt.Println("Solution 1:", count)
}

func printAllDevices(deviceMap map[string]*device_S) {
	for _, device := range deviceMap {
		device.print(true)
	}
}

func parseDevices(content []byte) map[string]*device_S {
	NEWLINE := "\n"
	SPACE := " "

	lines := strings.Split(string(content), NEWLINE)
	deviceMap := make(map[string]*device_S, len(lines))
	deviceMap["out"] = &device_S{id: "out", outputs: []*device_S{}}
	for _, line := range lines {
		id := line[0:3]
		deviceMap[id] = &device_S{id: id}
	}
	for _, line := range lines {
		id := line[0:3]
		outputs := strings.Split(line[5:], SPACE)
		devices := []*device_S{}
		for _, output := range outputs {
			devices = append(devices, deviceMap[output])
		}
		deviceMap[id].outputs = devices
	}

	return deviceMap
}

type device_S struct {
	id      string
	inputs  map[string]uint64
	outputs []*device_S
}

func (device *device_S) initInputs() {
	if device.inputs == nil {
		device.inputs = make(map[string]uint64, 1)
	}
}

func (device *device_S) updateInputs(other *device_S) {
	device.initInputs()
	device.inputs[other.id] = other.getInputs()
}

func (device *device_S) getInputs() (acc uint64) {
	for _, val := range device.inputs {
		acc += val
	}

	return
}

func (device *device_S) print(linebreak bool) {
	var outputsDisplay string
	for _, output := range device.outputs {
		outputsDisplay = outputsDisplay + " " + output.id
	}

	fmt.Print(device.id + "[" + fmt.Sprint(device.inputs) + "]" + ":" + outputsDisplay)
	if linebreak {
		fmt.Println()
	}
}

func search(device *device_S, count *int, callerMap map[string]bool) {
	if device == nil {
		return
	}

	if device.id == "out" {
		*count++
		return
	}

	filteredOutputs := filter(device.outputs, callerMap)
	for i := range filteredOutputs {
		callerMap[filteredOutputs[i].id] = true
		search(filteredOutputs[i], count, callerMap)
		delete(callerMap, filteredOutputs[i].id)
	}
}

func filter(outputs []*device_S, callerMap map[string]bool) []*device_S {
	filtered := make([]*device_S, 0, len(outputs))
	for i := range outputs {
		if _, ok := callerMap[outputs[i].id]; !ok {
			filtered = append(filtered, outputs[i])
		}
	}

	return filtered
}

func execC11Two(deviceMap map[string]*device_S) {
	pathMap := map[string]uint64{}

	startReplaceStrategy(deviceMap, "svr", "fft")
	pathMap["svr->fft"] = deviceMap["fft"].getInputs()

	resetInputs(deviceMap)
	startReplaceStrategy(deviceMap, "fft", "dac")
	pathMap["fft->dac"] = deviceMap["dac"].getInputs()

	resetInputs(deviceMap)
	startReplaceStrategy(deviceMap, "dac", "out")
	pathMap["dac->out"] = deviceMap["out"].getInputs()

	resetInputs(deviceMap)

	solution := pathMap["svr->fft"] * pathMap["fft->dac"] * pathMap["dac->out"]
	// path via svr->dac->fft->out not needed as there is no path from dac->fft
	fmt.Println("Solution 2:", solution)
}

func startReplaceStrategy(deviceMap map[string]*device_S, start string, end string) {
	replacements := make(map[string]bool)

	for _, device := range deviceMap[start].outputs {
		device.initInputs()
		device.inputs[start] = 1
		replacements[device.id] = true
	}

	replaceStrategy(replacements, deviceMap, end)
}

func replaceStrategy(replacements map[string]bool, deviceMap map[string]*device_S, end string) {
	if len(replacements) == 0 {
		return
	}

	nextReplacements := make(map[string]bool)

	for key := range replacements {
		device := deviceMap[key]
		for _, output := range device.outputs {
			output.updateInputs(device)
			if output.id != end {
				nextReplacements[output.id] = true
			}
		}
	}

	replaceStrategy(nextReplacements, deviceMap, end)
}

func resetInputs(deviceMap map[string]*device_S) {
	for key := range deviceMap {
		deviceMap[key].inputs = nil
	}
}
