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
		device.print()
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
	outputs []*device_S
}

func (device *device_S) print() {
	var outputsDisplay string
	for _, output := range device.outputs {
		outputsDisplay = outputsDisplay + " " + output.id
	}

	fmt.Println(device.id + ":" + outputsDisplay)
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
	var count int
	searchTwo(deviceMap["svr"], &count, make(map[string]bool))
	fmt.Println("Solution 2:", count)
}

func searchTwo(device *device_S, count *int, callerMap map[string]bool) {
	includes := func(map[string]bool) bool {
		_, ok1 := callerMap["dac"]
		_, ok2 := callerMap["fft"]
		return ok1 && ok2
	}

	if device.id == "out" {
		if includes(callerMap) {
			fmt.Println(callerMap)
			*count++
		}
		return
	}

	filteredOutputs := filter(device.outputs, callerMap)
	for i := range filteredOutputs {
		callerMap[filteredOutputs[i].id] = true
		searchTwo(filteredOutputs[i], count, callerMap)
		delete(callerMap, filteredOutputs[i].id)
	}
}
