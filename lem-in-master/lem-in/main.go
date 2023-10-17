package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Graph struct {
	Rooms     []*Room
	startRoom string
	endRoom   string
	ants      int
}
type Room struct {
	Roomname string
	adjacent []string
	visited  bool
}

func main() {
	list1 := []*Room{}
	// create a new Graph struct with the Rooms field set to "list1"
	roomList := &Graph{Rooms: list1}
	// read data from the file and populate the graph
	if err := SortFiles(roomList); err != nil {
		fmt.Print(err)
		return
	}

	// open the file specified as a command-line argument
	file, _ := os.Open(os.Args[1])
	// create a scanner to read its contents line by line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	// each line is printed to standard output
	for scanner.Scan() {
		x := scanner.Text()
		fmt.Println(x)
	}
	// ------------------------------------------------------
	allPathsDFS := []string{}
	allPathsBFS := []string{}
	var path string
	// call the DFS func with the startroom, endroom, the graph
	DFS(roomList.startRoom, roomList.endRoom, roomList, path, &allPathsDFS)

	// DO THE SAME FOR BFS
	list2 := []*Room{}
	roomList1 := &Graph{Rooms: list2}
	SortFiles(roomList1)
	BFS(roomList1.startRoom, roomList1.endRoom, roomList1, &allPathsBFS, ShortestPath)

	lenSorter(&allPathsBFS)
	lenSorter(&allPathsDFS)

	// set the value of 'antNum' to the .ants field of the roomList 
	// (retrieving the number of ants that will be moving through the graph, storing it in the 'antNum')
	antNum := roomList.ants
	DFSSearch := AntSender(antNum, allPathsDFS)
	BFSSearch := AntSender(antNum, allPathsBFS)

	//if the length of DFS is less than the length of the BFS, 
	Printer := []string{}
	if len(DFSSearch) < len(BFSSearch) {
		// assign the value of DFS too Printer 
		Printer = DFSSearch
	} else {
		// otherwise, assing the value of BFS to Printer
		Printer = BFSSearch
	}
	fmt.Println()
	// print each string in the Printer slice on a new line
	for _, step := range Printer {
		fmt.Println(step) // This logic is used to determine whether the DFS or BFS 
		// found a shorter path, and then print the shortest path found by either algorithm
	}
}

// read and parse the input file and populate the graph
func SortFiles(g *Graph) error {
	// open the file
	file, _ := os.Open(os.Args[1])
	// create a scanner to read its contents line by line
	scanner := bufio.NewScanner(file)
	start := false
	end := false
	i := 0
	firstLine := true
	scanner.Split(bufio.ScanLines)

	// iterate through the file lines, the 1st line is special - it contains the num of ants
	for scanner.Scan() {
		x := scanner.Text()
		// read the first line
		if firstLine {
			// convert the num of ants to an int
			g.ants, _ = strconv.Atoi(x)
			if g.ants == 0 {
				return errors.New("ERROR: invalid data format")
			}
			firstLine = false
		}
		// check each line of the file for rooms, links, start/endroom
		// split each line into words 
		space := strings.Split(scanner.Text(), " ")
		if len(space) > 1 {
			g.AddRoom(space[0])
			i++
		}
		if start {
			g.startRoom = g.Rooms[i-1].Roomname
			start = false
		} else if end {
			g.endRoom = g.Rooms[i-1].Roomname
			end = false
		}
		// if there are 2 words separated by a hyphen
		hyphen := strings.Split(scanner.Text(), "-")
		if len(hyphen) > 1 {
			if hyphen[0] == hyphen[1] {
				return errors.New("ERROR: invalid data format")
			}
			// add a link between the rooms to the 'Links' field of the corresponding rooms in the 'Graph'
			g.AddLinks(hyphen[0], hyphen[1])
		}
		// if the line contains '##start' -> set a flag to indicate that the next room is startrom
		if x == "##start" {
			start = true
		}
		if x == "##end" {
			end = true
		}
	}
	// if the input is successfully read and parsed - return nil
	return nil
}

func (g *Graph) AddRoom(name string) {
	g.Rooms = append(g.Rooms, &Room{Roomname: name, adjacent: []string{}, visited: false})
}
// adds an edge between two rooms. It first checks if the rooms exist, 
// and then checks if a link already exists between them. If not, it adds the link
func (g *Graph) AddLinks(from, to string) {
	fromRoom := g.getRoom(from)
	toRoom := g.getRoom(to)
	if fromRoom == nil || toRoom == nil {
		err := fmt.Errorf("Room doesn't exsist (%v --- %v)", from, to)
		fmt.Println(err.Error())
	} else if contains(fromRoom.adjacent, to) || contains(toRoom.adjacent, from) {
		err := fmt.Errorf(" Existing Link (%v --- %v)", from, to)
		fmt.Println(err.Error())
	} else if fromRoom.Roomname == g.endRoom {
		toRoom.adjacent = append(toRoom.adjacent, fromRoom.Roomname)
	} else if toRoom.Roomname == g.endRoom {
		fromRoom.adjacent = append(fromRoom.adjacent, toRoom.Roomname)
	} else if toRoom.Roomname == g.startRoom {
		toRoom.adjacent = append(toRoom.adjacent, fromRoom.Roomname)
	} else if fromRoom.Roomname == g.startRoom {
		fromRoom.adjacent = append(fromRoom.adjacent, toRoom.Roomname)
	} else if fromRoom.Roomname != g.endRoom && toRoom.Roomname != g.endRoom {
		fromRoom.adjacent = append(fromRoom.adjacent, toRoom.Roomname)
		toRoom.adjacent = append(toRoom.adjacent, fromRoom.Roomname)
	}
}
// searches for a room with a given name and returns a pointer to it
func (g *Graph) getRoom(name string) *Room {
	for i, v := range g.Rooms {
		if v.Roomname == name {
			return g.Rooms[i]
		}
	}
	return nil
}

// checks if a given string is present in a slice of strings
func contains(s []string, name string) bool {
	for _, v := range s {
		if name == v {
			return true
		}
	}
	return false
}

func BFS(start, end string, g *Graph, paths *[]string, f func(graph *Graph, start string, end string, path Array) Array) {
	begin := g.getRoom(start)
	if len(begin.adjacent) == 2 {
		begin.adjacent[0], begin.adjacent[1] = begin.adjacent[1], begin.adjacent[0]
	}
	for i := 0; i < len(begin.adjacent); i++ {
		var shortPath Array
		//  a helper function that finds the shortest path between two rooms using BFS
		ShortestPath(g, g.startRoom, g.endRoom, shortPath)
		var shortStorer string
		if len(pathArray) != 0 {
			shortStorer = pathArray[0]
		}
		for _, v := range pathArray {
			if len(v) < len(shortStorer) {
				shortStorer = v
			}
		}
		if len(pathArray) != 0 {
			shortStorer = shortStorer[1 : len(shortStorer)-1]
		}
		shortStorerSlc := strings.Split(shortStorer, " ")
		shortStorerSlc = shortStorerSlc[1:]
		for z := 0; z < len(shortStorerSlc)-1; z++ {
			g.getRoom(shortStorerSlc[z]).visited = true
		}
		var pathStr string
		if len(shortStorerSlc) != 0 {
			for i := 0; i < len(shortStorerSlc); i++ {
				if i == len(shortStorerSlc)-1 {
					pathStr += shortStorerSlc[i]
				} else {
					pathStr = pathStr + shortStorerSlc[i] + "-"
				}
			}
		}
		if len(pathStr) != 0 {
			containing := false
			for _, v := range *paths {
				if v == pathStr {
					containing = true
				}
			}
			if !containing {
				*paths = append(*paths, pathStr)
			}
		}
		pathArray = []string{}
	}
}

func DFS(current, end string, g *Graph, path string, pathList *[]string) {
	curr := g.getRoom(current)
	if current != end {
		curr.visited = true
	}
	if curr.Roomname == g.endRoom {
		path += current
	} else if !(curr.Roomname == g.startRoom) {
		path += current + "-"
	}
	final := false
	if current == end {
		*pathList = append(*pathList, path)
		path = ""
		final = true
		for i := 0; i < len(g.getRoom(g.startRoom).adjacent); i++ {
			if g.getRoom(g.startRoom).adjacent[i] == g.endRoom {
				g.getRoom(g.startRoom).adjacent[i] = ""
			}
		}
	}
	if final {
		DFS(g.startRoom, end, g, path, pathList)
	}
	for i := 0; i < len(curr.adjacent); i++ {
		if curr.adjacent[i] == g.endRoom {
			curr.adjacent[0], curr.adjacent[i] = curr.adjacent[i], curr.adjacent[0]
		}
	}
	for i := 0; i < len(curr.adjacent); i++ {
		if curr.adjacent[i] == "" {
			continue
		}
		x := g.getRoom(curr.adjacent[i])
		if x.visited {
			continue
		} else {
			DFS(x.Roomname, end, g, path, pathList)
		}
	}
}

type Array []string

var pathArray Array

func (arr Array) hasPropertyOf(str string) bool {
	for _, v := range arr {
		if str == v {
			return true
		}
	}
	return false
}
func ShortestPath(graph *Graph, start string, end string, path Array) Array {
	path = append(path, start)
	if start == end {
		return path
	}
	shortest := make([]string, 0)
	for _, node := range graph.getRoom(start).adjacent {
		if !path.hasPropertyOf(node) && !graph.isVisited(node) {
			newPath := ShortestPath(graph, node, end, path)
			if len(newPath) > 0 {
				if newPath.hasPropertyOf(graph.startRoom) && newPath.hasPropertyOf(end) {
					pathArray = append(pathArray, fmt.Sprint(newPath))
				}
			}
		}
	}
	return shortest
}
func (graph *Graph) isVisited(str string) bool {
	return graph.getRoom(str).visited
}

// sort a slice of strings by length
func lenSorter(paths *[]string) {
	x := *paths
	for i := 0; i < len(x); i++ {
		for j := 0; j < len(x); j++ {
			if len(x[i]) < len(x[j]) {
				x[i], x[j] = x[j], x[i]
			}
		}
	}
	*paths = x
}
func AntSender(n int, pathList []string) []string {
	pathListStore := [][]string{}
	for _, v := range pathList {
		s := strings.Split(v, "-")
		pathListStore = append(pathListStore, s)
	}
	lenP := len(pathList)
	queue := make([][]string, lenP)
	x := 0
	for i := 1; i <= n; i++ {
		ant := strconv.Itoa(i)
		if x == lenP-1 {
			if len(pathListStore[x])+len(queue[x]) <= len(pathListStore[0])+len(queue[0]) {
				queue[x] = append(queue[x], ant)
			} else {
				x = 0
				queue[x] = append(queue[x], ant)
			}
		} else {
			if len(pathListStore[x])+len(queue[x]) <= len(pathListStore[x+1])+len(queue[x+1]) {
				queue[x] = append(queue[x], ant)
			} else {
				x++
				queue[x] = append(queue[x], ant)
			}
		}
	}
	longest := len(queue[0])
	for i := 0; i < len(queue); i++ {
		if len(queue[i]) > longest {
			longest = len(queue[i])
		}
	}
	order := []int{}
	for j := 0; j < longest; j++ {
		for i := 0; i < len(queue); i++ {
			if j < len(queue[i]) {
				x, _ = strconv.Atoi(queue[i][j])
				order = append(order, x)
			}
		}
	}
	container := make([][][]string, len(queue))
	for i := 0; i < len(queue); i++ {
		for _, a := range queue[i] {
			adder := []string{}
			for _, room := range pathListStore[i] {
				str := "L" + a + "-" + room
				adder = append(adder, str)
			}
			container[i] = append(container[i], adder)
		}
	}
	finalMoves := []string{}
	for _, paths := range container {
		for j, moves := range paths {
			for k, room := range moves {
				if j+k > len(finalMoves)-1 {
					finalMoves = append(finalMoves, room+" ")
				} else {
					finalMoves[j+k] += room + " "
				}
			}
		}
	}
	return finalMoves
}
