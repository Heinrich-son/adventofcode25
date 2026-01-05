package main

import (
	"fmt"
	"strings"
)

/*
Graph problem: Find all paths fron a to b
*/

func runC11() {
	deviceMap := readFile("c11_test.txt", parseDevices)
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
	inputs  uint64
	outputs []*device_S
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
	pathMap["svr->fft"] = deviceMap["fft"].inputs

	resetInputs(deviceMap)
	startReplaceStrategy(deviceMap, "fft", "dac")
	pathMap["fft->dac"] = deviceMap["dac"].inputs

	resetInputs(deviceMap)
	startReplaceStrategy(deviceMap, "dac", "out")
	pathMap["dac->out"] = deviceMap["out"].inputs

	resetInputs(deviceMap)

	solution := pathMap["svr->fft"] * pathMap["fft->dac"] * pathMap["dac->out"]
	// path via svr->dac->fft->out not needed as there is no path from dac->fft
	fmt.Println("Solution 2:", solution, pathMap)
}

func startReplaceStrategy(deviceMap map[string]*device_S, start string, end string) {
	replacementSlice := make([]string, 0, len(deviceMap[start].outputs))
	for _, device := range deviceMap[start].outputs {
		device.inputs = 1
		replacementSlice = append(replacementSlice, device.id)
	}

	recursiveReplaceStrategy(replacementSlice, deviceMap, end)
}

func recursiveReplaceStrategy(replacements []string, deviceMap map[string]*device_S, end string) {
	if len(replacements) == 0 {
		return
	}

	nextReplacements := make([]string, 0, 10)

	for _, key := range replacements {
		device := deviceMap[key]
		for _, output := range device.outputs {
			output.inputs = output.inputs + device.inputs
			if output.id != end {
				nextReplacements = append(nextReplacements, output.id)
			}
		}
	}

	nextReplacements = filterDuplicates(nextReplacements)
	recursiveReplaceStrategy(nextReplacements, deviceMap, end)
}

func filterDuplicates(slice []string) []string {
	set := make(map[string]uint8, len(slice))
	newSlice := make([]string, 0, len(set))
	for _, val := range slice {
		if _, ok := set[val]; !ok {
			set[val] = 1
			newSlice = append(newSlice, val)
		}
	}

	return newSlice
}

func resetInputs(deviceMap map[string]*device_S) {
	for key := range deviceMap {
		deviceMap[key].inputs = 0
	}
}
