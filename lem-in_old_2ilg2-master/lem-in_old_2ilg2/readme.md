Here's the order of functions that are called in the code to achieve the result:

main() - this is the main function that runs the program and calls all the other functions.

readAntsFile() - this function reads the input file and returns its contents as a string slice.

NumAnts() - this function takes the string slice returned by readAntsFile() as input and returns the number of ants.

StartRoom() - this function takes the string slice returned by readAntsFile() as input and returns the name of the start room.

EndRoom() - this function takes the string slice returned by readAntsFile() as input and returns the name of the end room.

errHandling() - this function checks the input file for errors such as invalid number of ants, missing start or end rooms, and so on.

BFS() - this function implements the breadth-first search algorithm on the graph to find the shortest path from the start room to all other rooms.

Graph.Print() - this function prints the graph with the shortest path from the start room to each room.

DFS() - this function implements the depth-first search algorithm on the graph to find all possible paths from the start room to the end room.

PathSelection() - this function takes the output of the BFS and DFS algorithms as input and returns the optimal path(s) from the start room to the end room.

PathDupeCheck() - this function removes duplicate starting points from the output of PathSelection().

Reassign() - this function reassigns the paths in ascending order based on the number of rooms in each path.

pathSlice() - this function creates a slice of slices where the index 0 represents the number of rooms in each path.

lowestInt() - this function takes the output of pathSlice() and Reassign() as input and returns the path with the fewest number of rooms.

Increment() - this function increments the zero index of a slice of slices by 1.

RemoveAnt() - this function removes an ant from a slice of ants.

main() (continued) - this function creates ants and assigns them to paths in ascending order based on the number of rooms in each path. It then moves the ants along their paths until all ants have reached the end room. Finally, it prints the final output of the program.