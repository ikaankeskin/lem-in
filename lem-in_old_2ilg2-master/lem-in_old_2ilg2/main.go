package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Ants is a struct that contains a slice of pointers to Ant structs
type Ants struct {
	antz []*Ant
}

// Ant is a struct that contains the key, path and current room of an ant
type Ant struct {
	key         string
	path        []*Room
	currentRoom Room
}

// Graph is a struct that contains a slice of pointers to Room structs
type Graph struct {
	rooms []*Room
}

// Room is a struct that contains the key, adjacent rooms, path, visited and occupied status of a room
type Room struct {
	key      string
	adjacent []*Room
	path     []*Room
	visited  bool
	occupied bool
}

var startRoom string
var endRoom string

// Graph making function -1
// Add Room to a graph
func (g *Graph) AddRoom(k string) {
	if contains(g.rooms, k) {
		err := fmt.Errorf("Room %v not added because it is an existing key", k)
		fmt.Println(err.Error())
	} else {
		g.rooms = append(g.rooms, &Room{key: k})
	}
}

// getRoom returns a pointer to the Room key integer
func (g *Graph) getRoom(k string) *Room {
	for i, v := range g.rooms {
		if v.key == k {
			return g.rooms[i]
		}
	}
	return nil
}

// contains checks if the Room key exists in a given room struct.
func contains(s []*Room, k string) bool {
	for _, v := range s {
		if k == v.key {
			return true
		}
	}
	return false
}

// finds if a s is inside a slice of strings
func doesContain(s string, sl []string) bool {
	for _, word := range sl {
		if s == word {
			return true
		}
	}
	return false
}

// check if a room exists in the queue
func doesContainRoom(queue []*Room, k string) bool {
	roomMap := make(map[string]bool)
	for _, room := range queue {
		roomMap[room.key] = true
	}
	return roomMap[k]
}

// prints all paths within the given graph
func (g *Graph) PrintPath() {
	args := ""

	if len(os.Args) == 2 {

		args = os.Args[1]
	} else {
		err := fmt.Errorf("ERROR: Invalid number of arguments")
		fmt.Println(err.Error())
		os.Exit(0)
	}
	fmt.Println(StartRoom(readAntsFile(args)))
	for _, v := range g.rooms {
		for _, r := range v.path {
			fmt.Println(r)
		}
	}
}

// Add edge to the graph. deals with a directional graph only but condition in the main makes it undirected
func (g *Graph) AddEdge(from, to string) {
	args := ""

	if len(os.Args) == 2 {

		args = os.Args[1]
	} else {
		err := fmt.Errorf("ERROR: Invalid number of arguments")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	//get Room
	fromRoom := g.getRoom(from)
	toRoom := g.getRoom(to)

	//check error
	if fromRoom == nil || toRoom == nil {
		err := fmt.Errorf("ERROR: invalid edge (%v-->%v)", from, to)
		fmt.Println(err.Error())
		os.Exit(0)
	} else if contains(fromRoom.adjacent, to) {
		err := fmt.Errorf("ERROR: existing edge (%v-->%v)", from, to)
		fmt.Println(err.Error())
		os.Exit(0)
	} else if fromRoom == toRoom {
		err := fmt.Errorf("ERROR: Invalid data format. Can not connect room to itself (%v --> %v)", from, to)
		fmt.Println(err.Error())
		os.Exit(0)
	} else if fromRoom.key == EndRoom(readAntsFile(args)) {

	} else {
		fromRoom.adjacent = append(fromRoom.adjacent, toRoom)
	}

}

// Print will print the adjacent list for each Room of the graph
func (g *Graph) Print() {

	args := ""

	if len(os.Args) == 2 {

		args = os.Args[1]
	} else {
		err := fmt.Errorf("ERROR: Invalid number of arguments")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	fmt.Printf("The number of ants is: %v\n", NumAnts(readAntsFile(args)))

	for _, v := range g.rooms {
		if v.key == StartRoom(readAntsFile(args)) {
			fmt.Printf("\nStart Room is %v : ", StartRoom(readAntsFile(args)))

		} else if v.key == EndRoom(readAntsFile(args)) {
			fmt.Printf("\nEnd Room is %v :", EndRoom(readAntsFile(args)))

		} else {
			fmt.Printf("\nRoom %v :", v.key)
		}

		if len(v.adjacent) > 0 {
			fmt.Printf(" %v", v.adjacent[0].key)
			for i := 1; i < len(v.adjacent); i++ {
				fmt.Printf("\n     %v", v.adjacent[i].key)
			}
		}
	}

	fmt.Println()
}

// global variable which will store all of the valid paths in a slice of slices of string.
var validPaths [][]*Room

// Slices to hold paths from both algos for comparison
var dfsPaths [][]*Room
var bfsPaths [][]*Room

// Depth First Search algorithm that operates recursively
func DFS(r *Room, g Graph) {
	args := ""

	if len(os.Args) == 2 {

		args = os.Args[1]
	} else {
		err := fmt.Errorf("ERROR: Invalid number of arguments")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	sRoom := g.getRoom(StartRoom(readAntsFile(args)))

	// set the room being checked visited status to true
	if r.key != EndRoom(readAntsFile(args)) {
		r.visited = true

		// range through the neighbours of the r
		for _, nbr := range r.adjacent {
			if !nbr.visited {
				/* for each neighbour that hasn't been visited,
				- append their key to the visited slice,
				- then apply dfs to them recursively,
				- then append their key to their path value
				*/

				nbr.path = append(r.path, nbr)
				if contains(nbr.path, EndRoom(readAntsFile(args))) {

					dfsPaths = append(dfsPaths, nbr.path)

				}

				DFS(nbr, Graph{g.rooms})

			}

		}

	} else {

		if len(sRoom.adjacent) > 1 && !contains(sRoom.adjacent, EndRoom(readAntsFile(args))) {

			sRoom.adjacent = sRoom.adjacent[1:][:]

			DFS(sRoom, Graph{g.rooms})

		} else {

		}
	}
	dfsPaths = PathDupeCheck(dfsPaths)

}

// Function to remove room using its key
func RemoveRoomIndex(s []*Room, index string) []*Room {
	for i, v := range s {
		if v.key == index {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// function to remove an element from a slice of ints
func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

// removes a string from a slice (unordered)
func remove(s []*Room, k string) []*Room {
	for i, v := range s {
		if v.key == k {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

// delete edge from starting room
func DeleteEdge(r *Room, g Graph) {

	for i := 0; i < len(r.path); i++ {
		for _, room := range g.rooms {
			//	for _ , edge := range room.adjacent
			for j := 0; j < len(room.adjacent); j++ {
				if room.adjacent[j] == r.path[i] {
					room.adjacent = remove(room.adjacent, r.key)
				}
			}
		}
	}
}

// Breadth First Search algorithm
func BFS(r *Room, g Graph) {
	args := ""

	if len(os.Args) == 2 {

		args = os.Args[1]
	} else {
		err := fmt.Errorf("ERROR: Invalid number of arguments")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	var vPaths [][]*Room

	//queue variable, procedurally populated with rooms yet to be visited
	var queue []*Room

	//set start room as visited
	r.visited = true

	//initialise queue with start room
	queue = append(queue, r)

	// }

	// checks if there is a link between start and end directly
	for i, v := range g.getRoom(StartRoom(readAntsFile(args))).adjacent {
		if v.key == g.getRoom(EndRoom(readAntsFile(args))).key {
			g.getRoom(EndRoom(readAntsFile(args))).path = append(g.getRoom(EndRoom(readAntsFile(args))).path, g.getRoom(StartRoom(readAntsFile(args))))
			vPaths = append(vPaths, g.getRoom(StartRoom(readAntsFile(args))).path)
			g.getRoom(StartRoom(readAntsFile(args))).adjacent = append(g.getRoom(StartRoom(readAntsFile(args))).adjacent[:i], g.getRoom(StartRoom(readAntsFile(args))).adjacent[i+1:]...)
		}

	}

	//checks the queue for the end room and if the queue is not empty

	for !contains(queue, g.getRoom(EndRoom(readAntsFile(args))).key) && len(queue) >= 1 {
		qfront := queue[0]

		for _, room := range qfront.adjacent {
			if !room.visited {
				room.visited = true
				room.path = append(qfront.path, room)
				//
				queue = append(queue, room)
			}

		}

		queue = queue[1:]

		if doesContainRoom(queue, g.getRoom(EndRoom(readAntsFile(args))).key) {

			for _, room := range g.rooms {
				room.visited = false
			}
			vPaths = append(vPaths, qfront.path)

			for _, r := range qfront.path {
				DeleteEdge(r, g)

			}
			if len(g.getRoom(StartRoom(readAntsFile(args))).adjacent) == 0 {

				break
			}

			if len(g.getRoom(StartRoom(readAntsFile(args))).adjacent) >= 1 {
				for _, froom := range g.getRoom(StartRoom(readAntsFile(args))).adjacent {
					for _, sroom := range froom.adjacent {
						if sroom.key != g.getRoom(EndRoom(readAntsFile(args))).key {
							break
						} else {
							BFS(g.getRoom(StartRoom(readAntsFile(args))), Graph{g.rooms})
							queue = queue[1:]
						}
					}
				}
			}
			BFS(g.getRoom(StartRoom(readAntsFile(args))), Graph{g.rooms})

		}
	}
	for _, v := range vPaths {
		v = append(v, g.getRoom(EndRoom(readAntsFile(args))))
		bfsPaths = append(bfsPaths, v)
	}
	bfsPaths = PathDupeCheck(bfsPaths)

}

// returns the optimal path between bfs & dfs algos
func PathSelection(bfs [][]*Room, dfs [][]*Room) [][]*Room {

	bfsPathNum := len(bfs)
	dfsPathNum := len(dfs)

	if bfsPathNum > dfsPathNum {
		validPaths = append(validPaths, bfsPaths...)
	} else if dfsPathNum > bfsPathNum {
		validPaths = PathDupeCheck(append(validPaths, dfsPaths...))
	} else {

		bfscounter := 0

		dfscounter := 0

		for _, path := range bfs {

			bfscounter += len(path)

		}

		for _, path := range dfs {
			dfscounter += len(path)
		}

		if bfscounter < dfscounter {
			validPaths = append(validPaths, bfs...)
		} else if dfscounter < bfscounter {
			validPaths = append(validPaths, dfs...)
		} else {
			validPaths = append(validPaths, bfs...)
		}

	}
	return validPaths

}

// error checking for duplicate starting points
func PathDupeCheck(path [][]*Room) [][]*Room {

	dataMap := make(map[*Room][]*Room)

	for _, item := range path {
		if value, ok := dataMap[item[0]]; !ok {
			dataMap[item[0]] = item
		} else {
			if len(item) <= len(value) {
				dataMap[item[0]] = item

			}
		}
	}

	var output [][]*Room

	for _, value := range dataMap {
		output = append(output, value)
	}

	return output
}

// reassigns the slices in ascending (len) order
func Reassign(a [][]*Room) [][]*Room {

	sort.Slice(a, func(i, j int) bool {
		return len(a[i]) < len(a[j])
	})

	return a

}

// returns a slice of slice with index 0 representing the number of rooms within a given path
func pathSlice(a [][]*Room) [][]int {
	var slice [][]int
	var s []int

	for i := range a {
		s = append(s, len(a[i]))
		slice = append(slice, s)
		s = []int{}
	}

	return slice
}

// finds most efficient path
func lowestInt(a [][]int, b [][]*Room) (int, []*Room) {

	min := 10000
	var path []*Room

	for i := 0; i < len(a); i++ {
		if a[i][0] < min {
			min = a[i][0]
			path = b[i]

		}

	}
	return min, path
}

// increments the zero index for the given array
func Increment(a [][]int, b int) [][]int {

	for _, slice := range a {
		if slice[0] == b {
			slice[0] += 1
			break
		}
	}
	return a

}

// Function to remove room using its key
func RemoveAnt(a []*Ant, b *Ant) []*Ant {
	ret := make([]*Ant, 0)
	if len(a) == 1 {
		return []*Ant{}
	}
	for i := 0; i < len(a); i++ {
		if a[i].key == b.key {
			ret = append(ret, a[:i]...)
			ret = append(ret, a[i+1:]...)
		}
	}
	return ret
}

func main() {
	args := ""

	if len(os.Args) == 2 {

		args = os.Args[1]
	} else {
		err := fmt.Errorf("ERROR: Invalid number of arguments")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	errHandling()

	bfsGraph := Graph{}

	//adding all rooms
	for i, line := range readAntsFile(args) {
		if strings.Contains(string(line), " ") {
			bfsGraph.AddRoom(strings.Split(readAntsFile(args)[i], " ")[0])
		}
		// adding all edges from and to rooms

		if strings.Contains(string(line), "-") {
			bfsGraph.AddEdge(strings.Split(readAntsFile(args)[i], "-")[0], strings.Split(readAntsFile(args)[i], "-")[1])
			bfsGraph.AddEdge(strings.Split(readAntsFile(args)[i], "-")[1], strings.Split(readAntsFile(args)[i], "-")[0])
		}

	}

	BFS(bfsGraph.getRoom(StartRoom(readAntsFile(args))), bfsGraph)
	// bfsGraph.Print()

	dfsGraph := Graph{}

	//adding all rooms
	for i, line := range readAntsFile(args) {
		if strings.Contains(string(line), " ") {
			dfsGraph.AddRoom(strings.Split(readAntsFile(args)[i], " ")[0])
		}
		// adding all edges from and to rooms
		// only adding edges in one direction to make the graph directional
		if strings.Contains(string(line), "-") {
			dfsGraph.AddEdge(strings.Split(readAntsFile(args)[i], "-")[0], strings.Split(readAntsFile(args)[i], "-")[1])
		}

	}

	DFS(dfsGraph.getRoom(StartRoom(readAntsFile(args))), dfsGraph)

	a := Ants{}

	arrange := pathSlice(Reassign(PathDupeCheck(PathSelection(bfsPaths, dfsPaths))))
	rooms := Reassign(PathDupeCheck(PathSelection(bfsPaths, dfsPaths)))

	// --------------------------------------------------------

	var unmovedAnts []*Ant
	var movedAnts []*Ant
	counter := 1

	for counter <= NumAnts(readAntsFile(args)) {

		number, _ := lowestInt(arrange, rooms)
		_, route := lowestInt(arrange, rooms)
		a.antz = append(a.antz, &Ant{key: "L" + strconv.Itoa(counter), path: route})
		Increment(arrange, number)

		counter++
	}

	for _, line := range readAntsFile(args) {
		fmt.Println(line)
	}
	fmt.Println()

	unmovedAnts = append(unmovedAnts, a.antz...)

	for len(unmovedAnts) > 0 || len(movedAnts) >= 1 {

		for _, ant := range unmovedAnts {
			if len(ant.path) == 1 {
				fmt.Print(ant.key, "-", ant.path[0].key, " ")
				ant.path[0].occupied = true
				unmovedAnts = RemoveAnt(unmovedAnts, ant)
				break
			}
		}

		for _, ant := range unmovedAnts {

			if !ant.path[0].occupied {
				fmt.Print(ant.key, "-", ant.path[0].key, " ")
				ant.path[0].occupied = true
				movedAnts = append(movedAnts, ant)
				unmovedAnts = RemoveAnt(unmovedAnts, ant)

			}

		}

		fmt.Println()

		for _, ant := range movedAnts {

			if len(ant.path) > 1 {
				ant.path[0].occupied = false

				ant.path = ant.path[1:]
				fmt.Print(ant.key, "-", ant.path[0].key, " ")

			} else {
				movedAnts = RemoveAnt(movedAnts, ant)
				ant.path = []*Room{}
			}
		}

	}

	fmt.Println()

}

// Reads file and returns a string slice
func readAntsFile(filename string) []string {
	//Open the file
	file, _ := os.Open(filename)
	defer file.Close()
	//Create a scanner to read the file line by line
	scanned := bufio.NewScanner(file)
	//Setting the split function for the scanner
	scanned.Split(bufio.ScanLines)

	var lines []string

	//scanning the file line by line and adding the lines to a string slice
	for scanned.Scan() {
		lines = append(lines, scanned.Text())
	}
	return lines
}

// returns the number of ants from the text file
func NumAnts(s []string) int {

	//getting the first line of the file which is the number of ants
	antNum := s[0]

	antNumInt, _ := strconv.Atoi(antNum)
	return antNumInt
}

func StartRoom([]string) string {
	if startRoom == "" {
		args := ""
		if len(os.Args) == 2 {
			args = os.Args[1]
		} else {
			err := fmt.Errorf("ERROR: Invalid number of arguments")
			fmt.Println(err.Error())
			os.Exit(0)
		}

		s := readAntsFile(args)

		for i := 0; i < len(s); i++ {
			if s[i] == "##start" {
				startRoom = strings.Split(string(s[i+1]), " ")[0]
				break
			}
		}
	}
	return startRoom
}

// Gets out the end room and returns it
func EndRoom([]string) string {
	if endRoom == "" {
		args := ""
		if len(os.Args) == 2 {
			args = os.Args[1]
		} else {
			err := fmt.Errorf("ERROR: Invalid number of arguments")
			fmt.Println(err.Error())
			os.Exit(0)
		}

		s := readAntsFile(args)

		/*s[i+1] is used to get the line after the line that contains "##end"
		string(s[i+1]) converts the line from a slice of bytes to a string.
		strings.Split(string(s[i+1]), " ") is used to split the string on the spaces. This returns a slice of strings where each element is a word in the original string.
		[0] is used to get the first element of the returned slice, which is the name of the end room.
		The value of the first element is then assigned to the endRoom variable.*/

		for i := 0; i < len(s); i++ {
			if s[i] == "##end" {
				endRoom = strings.Split(string(s[i+1]), " ")[0]
				break
			}

		}
	}
	return endRoom
}

func errHandling() {

	args := ""

	if len(os.Args) == 2 {

		args = os.Args[1]
	} else {
		err := fmt.Errorf("ERROR: Invalid number of arguments")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if NumAnts(readAntsFile(args)) <= 0 || NumAnts(readAntsFile(args)) > 19990 {
		err := fmt.Errorf("ERROR: Invalid data format")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if !doesContain("##start", readAntsFile(args)) {
		err := fmt.Errorf("ERROR: No start room found")
		fmt.Println(err.Error())
		os.Exit(0)
	}

	if !doesContain("##end", readAntsFile(args)) {
		err := fmt.Errorf("ERROR: No end room found")
		fmt.Println(err.Error())
		os.Exit(0)
	}

}
