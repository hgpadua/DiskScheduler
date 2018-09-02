/*Huey Padua
  Operating Systems
  Programming Assignment 2.1
  Disk Seek Algorithms in GO
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Schedule struct {
	algo     string
	lowerCyl int
	upperCyl int
	initCyl  int
}

func main() {

	// file := "/Users/hueycopter/Desktop/golang/sstf01.txt"
	file := os.Args[1]

	var schedule Schedule
	var requests []int

	infile, err := os.Open(file)
	defer infile.Close()

	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(infile)

	for scanner.Scan() {
		line := scanner.Text()
		s := strings.Split(line, " ")

		if s[0] == "use" {
			ns := strings.Split(s[1], "\t")
			schedule.algo = ns[0]
		}
		if s[0] == "lowerCYL" {
			ns := strings.Split(s[1], "\t")
			lc, err := strconv.Atoi(ns[0])
			if err != nil {
				fmt.Println(err)
			}
			schedule.lowerCyl = lc
		}
		if s[0] == "upperCYL" {
			ns := strings.Split(s[1], "\t")
			uc, err := strconv.Atoi(ns[0])
			if err != nil {
				fmt.Println(err)
			}
			schedule.upperCyl = uc
		}
		if s[0] == "initCYL" {
			ns := strings.Split(s[1], "\t")
			ic, err := strconv.Atoi(ns[0])
			if err != nil {
				fmt.Println(err)
			}
			schedule.initCyl = ic
		}
		if s[0] == "cylreq" {
			ns := strings.Split(s[1], "\t")
			cr, err := strconv.Atoi(ns[0])
			if err != nil {
				fmt.Println(err)
			}
			requests = append(requests, cr)

		}
		if s[0] == "end" {
			break
		}
	}
	if scanner.Err() != nil {
		fmt.Println(scanner.Err())
	}
	if schedule.algo == "fcfs" {
		fcfs(&schedule, requests)
	} else if schedule.algo == "sstf" {
		sstf(&schedule, requests)
	} else if schedule.algo == "scan" {
		scan(&schedule, requests)
	} else if schedule.algo == "look" {
		look(&schedule, requests)
	} else if schedule.algo == "c-scan" {
		cscan(&schedule, requests)
	} else if schedule.algo == "c-look" {
		clook(&schedule, requests)
	} else {
		fmt.Println("Error: Schedule algorithm is invalid")
	}
}

// Helper function that returns absolute value of x
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Helper function to find MaxInt
func getMax(val int, requests []int) int {
	temp := requests[0]
	for i := 0; i < len(requests); i++ {
		if requests[i] >= val {
			temp = requests[i]
		}
	}
	return temp
}

// Helper function to find min diff
func findClosest(val int, requests []int) int {
	curr := requests[0]
	for i := 0; i < len(requests); i++ {
		if requests[i] > val && abs(val-requests[i]) < abs(val-curr) {
			curr = requests[i]
		}
	}
	return curr
}

// Helper function to find closest seek time for SSTF
func findSeek(val int, requests []int) int {
	curr := requests[0]
	for i := 0; i < len(requests); i++ {
		if abs(val-requests[i]) < abs(val-curr) {
			curr = requests[i]
		}
	}
	return curr
}

//first come first serve implementation
func fcfs(schedule *Schedule, requests []int) {
	fmt.Printf("Seek algorithm: FCFS\n")
	fmt.Printf("\tLower cylinder: %5d\n", schedule.lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", schedule.upperCyl)
	fmt.Printf("\tInit cylinder: %5d\n", schedule.initCyl)
	fmt.Printf("\tCylinder requests:\n")
	for i := 0; i < len(requests); i++ {
		fmt.Printf("\t\tCylinder %5d\n", requests[i])
	}
	var curr = 0
	var total = 0
	stack := make([]int, len(requests))
	for i := 0; i < len(requests); i++ {
		// save a copy of cylinder requests
		stack[i] = requests[i]
		fmt.Printf("Servicing %5d\n", stack[i])
		// start from initial cylinder position
		if curr == 0 {
			curr = stack[i]
			total = abs(schedule.initCyl - curr)
		} else {
			curr = stack[i]
			// if value is out of bounds, throw err and continue
			if curr < schedule.lowerCyl || curr > schedule.upperCyl {
				fmt.Println("Error: cylinder request out of bounds")
				continue
			} else {
				prev := stack[i-1]
				total += abs(prev - curr)
			}
		}
	}
	fmt.Printf("FCFS traversal count = %d\n", total)
}

// Shortest seek time first implementation
func sstf(schedule *Schedule, requests []int) {
	fmt.Println("Seek algorithm: SSTF")
	fmt.Printf("\tLower cylinder: %5d\n", schedule.lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", schedule.upperCyl)
	fmt.Printf("\tInit cylinder: %5d\n", schedule.initCyl)
	fmt.Printf("\tCylinder requests:\n")
	for i := 0; i < len(requests); i++ {
		fmt.Printf("\t\tCylinder %5d\n", requests[i])
	}
	stack := make([]int, len(requests))
	startIndex := 0
	maxIndex := 0
	max := 0
	count := len(requests)
	// Get max in cylinder requests
	for i := 0; i < len(requests); i++ {
		if requests[i] >= getMax(max, requests) {
			max = requests[i]
		}
	}
	start := findSeek(schedule.initCyl, requests)
	// keep a copy of the cylinder requests
	for i := 0; i < len(requests); i++ {
		stack[i] = requests[i]
	}
	// sort stack in order then make startIndex
	// based on where the closest seek time is
	sort.Ints(stack)
	for i := 0; i < len(requests); i++ {
		if start == stack[i] {
			startIndex = i
		}
	}
	// get the index of max cylinder requests
	for i := 0; i < len(requests); i++ {
		if max == stack[i] {
			maxIndex = i
		}
	}

	var curr = 0
	var prev = 0
	var total = 0
	// start sstf implementation based on where
	// the closest seek time is located
	for i := startIndex; i < len(requests); i++ {
		// start from initCYL
		if curr == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			//decrement count everytime a request is serviced
			count--
			curr = stack[i]
			total = abs(schedule.initCyl - curr)
		} else {
			curr = stack[i]
			// ignore cylinder requests thats out of bounds
			if curr < schedule.lowerCyl || curr > schedule.upperCyl {
				fmt.Println("Error: cylinder request out of bounds")
				continue
			} else {
				prev = stack[i-1]
				if curr > prev && curr < max {
					// requests from after initCyl -> max request
					fmt.Printf("Servicing %d\n", stack[i])
					count--
					total += abs(prev - curr)
				}
				// Calculating maxIndex of requests
				for i := maxIndex; i < len(requests); i++ {
					if curr == max {
						fmt.Printf("Servicing %d\n", curr)
						count--
						total += abs(prev - curr)
						prev = curr
					}
				}
			}
		}
	}
	// lock to ensure no additional calculations
	// are being made
	lock := 0
	for i := 0; i < len(requests); i++ {
		// reverse stack order to fit sstf requirements
		sort.Sort(sort.Reverse(sort.IntSlice(stack)))
		// get closest seek time
		start := findSeek(schedule.initCyl, requests)
		// just in case we missed any requests that
		// is less than startIndex
		if start > stack[i] && lock == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			lock = 1
			total += abs(curr - stack[i])
			prev = stack[i]
		} else if start > prev && lock == 1 {
			fmt.Printf("Servicing %d\n", stack[i])
			total += abs(prev - stack[i])
			prev = stack[i]
		}
	}
	fmt.Printf("SSTF traversal count = %d\n", total)
}

//SCAN implementation
func scan(schedule *Schedule, requests []int) {
	fmt.Println("Seek algorithm: SCAN")
	fmt.Printf("\tLower cylinder: %5d\n", schedule.lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", schedule.upperCyl)
	fmt.Printf("\tInit cylinder: %5d\n", schedule.initCyl)
	fmt.Printf("\tCylinder requests:\n")
	for i := 0; i < len(requests); i++ {
		fmt.Printf("\t\tCylinder %5d\n", requests[i])
	}
	stack := make([]int, len(requests))
	startIndex := 0
	maxIndex := 0
	max := 0
	count := len(requests)
	// Get max in cylinder requests
	for i := 0; i < len(requests); i++ {
		if requests[i] >= getMax(max, requests) {
			max = requests[i]
		}
	}
	if len(requests) < 10 {
		start := findClosest(schedule.initCyl, requests)
		for i := 0; i < len(requests); i++ {
			for j := i + 1; j < len(requests); j++ {
				if start < requests[i] && requests[i] > requests[j] {
					temp := requests[i]
					requests[i] = requests[j]
					requests[j] = temp
				}
			}
			stack[i] = requests[i]
		}
	} else {
		start := findClosest(schedule.initCyl, requests)
		for i := 0; i < len(requests); i++ {
			stack[i] = requests[i]
		}
		sort.Ints(stack)
		for i := 0; i < len(requests); i++ {
			if start == stack[i] {
				startIndex = i
			}
		}
		for i := 0; i < len(requests); i++ {
			if max == stack[i] {
				maxIndex = i
			}
		}
	}
	var curr = 0
	var prev = 0
	var total = 0
	for i := startIndex; i < len(requests); i++ {

		if curr == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			count--
			curr = stack[i]
			total = abs(schedule.initCyl - curr)
		} else {
			curr = stack[i]
			if curr < schedule.lowerCyl || curr > schedule.upperCyl {
				fmt.Println("Error: cylinder request out of bounds")
				continue
			} else {
				prev = stack[i-1]
				if curr > prev && curr < max {
					// requests from init -> max int
					fmt.Printf("Servicing %d\n", stack[i])
					count--
					total += abs(prev - curr)
				}
				for i := maxIndex; i < len(requests); i++ {
					if curr == max {
						fmt.Printf("Servicing %d\n", curr)
						count--
						if count != 0 {
							total += abs(prev - curr)
							total += abs(stack[i] - schedule.upperCyl)
							break
						} else {
							total += abs(prev - curr)
							break
						}
					}
				}
			}
		}
	}
	lock := 0
	for i := 0; i < len(requests); i++ {
		sort.Sort(sort.Reverse(sort.IntSlice(stack)))
		if schedule.initCyl > stack[i] && lock == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			lock = 1
			total += abs(schedule.upperCyl - stack[i])
			prev = stack[i]
		} else if lock == 1 {
			fmt.Printf("Servicing %d\n", stack[i])
			total += abs(stack[i] - prev)
			prev = stack[i]
		}

	}
	fmt.Printf("SCAN traversal count = %d\n", total)
}

//LOOK implementation
func look(schedule *Schedule, requests []int) {
	fmt.Println("Seek algorithm: LOOK")
	fmt.Printf("\tLower cylinder: %5d\n", schedule.lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", schedule.upperCyl)
	fmt.Printf("\tInit cylinder: %5d\n", schedule.initCyl)
	fmt.Printf("\tCylinder requests:\n")
	for i := 0; i < len(requests); i++ {
		fmt.Printf("\t\tCylinder %5d\n", requests[i])
	}
	stack := make([]int, len(requests))
	startIndex := 0
	maxIndex := 0
	max := 0
	count := len(requests)
	// Get max in cylinder requests
	for i := 0; i < len(requests); i++ {
		if requests[i] >= getMax(max, requests) {
			max = requests[i]
		}
	}
	if len(requests) < 10 {
		start := findClosest(schedule.initCyl, requests)
		for i := 0; i < len(requests); i++ {
			for j := i + 1; j < len(requests); j++ {
				if start < requests[i] && requests[i] > requests[j] {
					temp := requests[i]
					requests[i] = requests[j]
					requests[j] = temp
				}
			}
			stack[i] = requests[i]
		}
	} else {
		start := findClosest(schedule.initCyl, requests)
		for i := 0; i < len(requests); i++ {
			stack[i] = requests[i]
		}
		sort.Ints(stack)
		for i := 0; i < len(requests); i++ {
			if start == stack[i] {
				startIndex = i
			}
		}
		for i := 0; i < len(requests); i++ {
			if max == stack[i] {
				maxIndex = i
			}
		}
	}
	var curr = 0
	var prev = 0
	var total = 0
	for i := startIndex; i < len(requests); i++ {

		if curr == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			count--
			curr = stack[i]
			total = abs(schedule.initCyl - curr)
		} else {
			curr = stack[i]
			if curr < schedule.lowerCyl || curr > schedule.upperCyl {
				fmt.Println("Error: cylinder request out of bounds")
				continue
			} else {
				prev = stack[i-1]
				if curr > prev && curr < max {
					// requests from init -> max int
					fmt.Printf("Servicing %d\n", stack[i])
					count--
					total += abs(prev - curr)
				}
				for i := maxIndex; i < len(requests); i++ {
					if curr == max {
						fmt.Printf("Servicing %d\n", curr)
						count--
						if count != 0 {
							total += abs(prev - curr)
							prev = stack[i]
							break
						} else {
							total += abs(prev - curr)
							break
						}
					}
				}
			}
		}
	}
	lock := 0
	for i := 0; i < len(requests); i++ {
		sort.Sort(sort.Reverse(sort.IntSlice(stack)))
		if schedule.initCyl > stack[i] && lock == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			lock = 1
			total += abs(prev - stack[i])
			prev = stack[i]
		} else if lock == 1 {
			fmt.Printf("Servicing %d\n", stack[i])
			total += abs(stack[i] - prev)
			prev = stack[i]
		}

	}
	fmt.Printf("LOCK traversal count = %d\n", total)
}

//C-SCAN implementation
func cscan(schedule *Schedule, requests []int) {
	fmt.Println("Seek algorithm: C-SCAN")
	fmt.Printf("\tLower cylinder: %5d\n", schedule.lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", schedule.upperCyl)
	fmt.Printf("\tInit cylinder: %5d\n", schedule.initCyl)
	fmt.Printf("\tCylinder requests:\n")
	for i := 0; i < len(requests); i++ {
		fmt.Printf("\t\tCylinder %5d\n", requests[i])
	}
	stack := make([]int, len(requests))
	startIndex := 0
	maxIndex := 0
	max := 0
	count := len(requests)
	// Get max in cylinder requests
	for i := 0; i < len(requests); i++ {
		if requests[i] >= getMax(max, requests) {
			max = requests[i]
		}
	}
	if len(requests) < 10 {
		start := findClosest(schedule.initCyl, requests)
		for i := 0; i < len(requests); i++ {
			for j := i + 1; j < len(requests); j++ {
				if start < requests[i] && requests[i] > requests[j] {
					temp := requests[i]
					requests[i] = requests[j]
					requests[j] = temp
				}
			}
			stack[i] = requests[i]
		}
	} else {
		start := findClosest(schedule.initCyl, requests)
		for i := 0; i < len(requests); i++ {
			stack[i] = requests[i]
		}
		sort.Ints(stack)
		for i := 0; i < len(requests); i++ {
			if start == stack[i] {
				startIndex = i
			}
		}
		for i := 0; i < len(requests); i++ {
			if max == stack[i] {
				maxIndex = i
			}
		}
	}
	var curr = 0
	var prev = 0
	var total = 0
	for i := startIndex; i < len(requests); i++ {

		if curr == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			count--
			curr = stack[i]
			total = abs(schedule.initCyl - curr)
		} else {
			curr = stack[i]
			if curr < schedule.lowerCyl || curr > schedule.upperCyl {
				fmt.Println("Error: cylinder request out of bounds")
				continue
			} else {
				prev = stack[i-1]
				if curr > prev && curr < max {
					// requests from init -> max int
					fmt.Printf("Servicing %d\n", stack[i])
					count--
					total += abs(prev - curr)
				}
				for i := maxIndex; i < len(requests); i++ {
					if curr == max {
						fmt.Printf("Servicing %d\n", curr)
						count--
						if count != 0 {
							total += abs(prev - curr)
							total += abs(stack[i] - schedule.upperCyl)
							total += abs(schedule.upperCyl - schedule.lowerCyl)
							break
						} else {
							total += abs(prev - curr)
							break
						}
					}
				}
			}
		}
	}
	lock := 0
	for i := 0; i < len(requests); i++ {
		sort.Sort((sort.IntSlice(stack)))
		if schedule.initCyl > stack[i] && lock == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			lock = 1
			total += abs(schedule.lowerCyl - stack[i])
			prev = stack[i]
		} else if schedule.initCyl > stack[i] && lock == 1 {
			fmt.Printf("Servicing %d\n", stack[i])
			total += abs(stack[i] - prev)
			prev = stack[i]

		}

	}
	fmt.Printf("C-SCAN traversal count = %d\n", total)
}

//C-LOOK
func clook(schedule *Schedule, requests []int) {
	fmt.Println("Seek algorithm: C-LOOK")
	fmt.Printf("\tLower cylinder: %5d\n", schedule.lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", schedule.upperCyl)
	fmt.Printf("\tInit cylinder: %5d\n", schedule.initCyl)
	fmt.Printf("\tCylinder requests:\n")
	for i := 0; i < len(requests); i++ {
		fmt.Printf("\t\tCylinder %5d\n", requests[i])
	}
	stack := make([]int, len(requests))
	startIndex := 0
	maxIndex := 0
	max := 0
	count := len(requests)
	// Get max in cylinder requests
	for i := 0; i < len(requests); i++ {
		if requests[i] >= getMax(max, requests) {
			max = requests[i]
		}
	}
	if len(requests) < 10 {
		start := findClosest(schedule.initCyl, requests)
		for i := 0; i < len(requests); i++ {
			for j := i + 1; j < len(requests); j++ {
				if start < requests[i] && requests[i] > requests[j] {
					temp := requests[i]
					requests[i] = requests[j]
					requests[j] = temp
				}
			}
			stack[i] = requests[i]
		}
	} else {
		start := findClosest(schedule.initCyl, requests)
		for i := 0; i < len(requests); i++ {
			stack[i] = requests[i]
		}
		sort.Ints(stack)
		for i := 0; i < len(requests); i++ {
			if start == stack[i] {
				startIndex = i
			}
		}
		for i := 0; i < len(requests); i++ {
			if max == stack[i] {
				maxIndex = i
			}
		}
	}
	var curr = 0
	var prev = 0
	var total = 0
	for i := startIndex; i < len(requests); i++ {

		if curr == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			count--
			curr = stack[i]
			total = abs(schedule.initCyl - curr)
		} else {
			curr = stack[i]
			if curr < schedule.lowerCyl || curr > schedule.upperCyl {
				fmt.Println("Error: cylinder request out of bounds")
				continue
			} else {
				prev = stack[i-1]
				if curr > prev && curr < max {
					// requests from init -> max int
					fmt.Printf("Servicing %d\n", stack[i])
					count--
					total += abs(prev - curr)
				}
				for i := maxIndex; i < len(requests); i++ {
					if curr == max {
						fmt.Printf("Servicing %d\n", curr)
						count--
						if count != 0 {
							total += abs(prev - curr)
							prev = stack[i]
							break
						} else {
							total += abs(prev - curr)
							break
						}
					}
				}
			}
		}
	}
	lock := 0
	for i := 0; i < len(requests); i++ {
		sort.Sort((sort.IntSlice(stack)))
		if schedule.initCyl > stack[i] && lock == 0 {
			fmt.Printf("Servicing %d\n", stack[i])
			lock = 1
			total += abs(prev - stack[i])
			prev = stack[i]
		} else if schedule.initCyl > stack[i] && lock == 1 {
			fmt.Printf("Servicing %d\n", stack[i])
			total += abs(stack[i] - prev)
			prev = stack[i]
		}

	}
	fmt.Printf("C-LOCK traversal count = %d\n", total)
}
