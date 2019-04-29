package main

import (
	"fmt"
	"math"
)

// the code follows Ukkonen's algorithm from the lecture notes

// The code is structured as follows (TODO: draw a better ascii diagram)
// Nodes have edges  * -------- () ------- * 
//                              |
//                              |
//							    *
// Edges can end in a node -----> ()
// Edges do not have to have to end in a node
// Edges only have a node if its split/there's a mismatch 
// Future updates
// 1. exend this to build from multiple strings - Generalized Suffix tree
// 2. parallel implementation - https://people.csail.mit.edu/jshun/JDA2017.pdf
// 3. I've noticed, we can remove a lot of nodes that are repeated 
//    and can compress the tree more  when suffix links are added/traversed
//    probably makes the tree a suffix_graph

// Edge Structure
// start_index int : start_index in the original string
// end_index int : end_index in the origin string
// end_node node : only if end_index != len(text)
type edge struct {
	start_index, end_index int
	end_node *node
}

// Node Structure
// edges map : a map of all possible edges from this node
// suffix_link node : if there is a suffix link to a node
// TODO: ntype string : refactor code to remove this param
type node struct {
	edges map[string]*edge
	suffix_link *node
	ntype string
}

// Search for a pattern
// query string : query to search
// text string : original text used in constructing the suffix tree
// min_match_length int : if entire query doesn't match, report all lengths above
//        this threshold
// TODO: all_matches boolean : include all possible matches
func (stree node) search(query, text string, min_match_length int) map[int]string {
	found := make(map[int]string)
	var qindex, findex int = 0, 0
	var active_node *node = &stree
	
	for qindex < len(query) {
		tq := string(query[qindex])

		fmt.Println("tq - ", tq)
		fmt.Println("qindex - ", qindex)

		tedge, ok := active_node.edges[tq]
		if !ok {
			found[findex] = string(query[qindex - 1])
			break
		}	

		fmt.Println("tedge - ", tedge)
		findex = tedge.start_index
		fmt.Println("findex - ", findex)
		// found[findex] = tq
		// fmt.Println("found - ", found)

		var min_ext int = int(math.Min(float64(tedge.end_index - tedge.start_index), float64(len(query) - qindex)))
		fmt.Println("min_ext - ", min_ext)

		if text[findex: findex + min_ext] == query[qindex: qindex + min_ext] {
			// matched edge, skip all the way
			fmt.Println("matched substring qindex - ", qindex)

			// found[findex] = string(query[qindex: qindex + min_ext])

			fmt.Println("found - ", found)

			if tedge.end_node != nil {
				active_node = tedge.end_node
			} else {
				found[findex] = string(query[qindex])
				break
			}
			qindex = qindex +  min_ext

			fmt.Println("active_node - ", active_node)

		} else {
			// did not match, find the substring matched
			// need to do some work here
			var ftindex int = 0
			for ftindex + findex < len(query) {
				fmt.Println("ftindex - ", ftindex)

				if query[qindex + ftindex] == text[findex + ftindex] {
					fmt.Println("hola")
					// found[fi/ndex] = found[findex] + string(query[qindex + ftindex])
				} else {
					found[findex] = string(query[qindex])
					break;
				}

				ftindex++
			}	
		}
	}

	return found
}

// suffix link structure
// not required can just point to a node
// type suffix_link struct {
// 	snode node
// }


// Builds a suffix tree from a text
// text does not have to end in '$'
func build_suffix_tree(text string) *node {
	root := new(node)
	root.edges = make(map[string]*edge)
	root.ntype = "root"
	var active_node *node = root
	var active_edge *edge
	var active_string string
	var active_length, remainder int = 0, 0
	var last_node *node
	var text_length = len(text)
	var cindex int = 0

	for cindex < text_length {
		ts := string(text[cindex])
		// fmt.Println("Iterating - ",cindex, ts)
		if remainder == 0 {
			// fmt.Println("\t remainder is 0")
			remainder++
			tedge, ok := active_node.edges[ts]
			if ok {
				active_edge = tedge
			}

			// fmt.Println("\t ok - ", ok)
			// fmt.Println("\t active_edge - ", active_edge)
		} else if active_edge == nil {
			// fmt.Println("\t active edge is nil -> see if one exists")

			tedge, ok := active_node.edges[ts]
			if ok {
				active_edge = tedge
			}

			// fmt.Println("\t ok - ", ok)
			// fmt.Println("\t active_edge - ", active_edge)
		}
		// fmt.Println("\t active_node - ", active_node)
		// fmt.Println("\t after remainder")
		// fmt.Println("\t active edge - ", active_edge)

		if active_edge == nil {
			// no match or no active edge
			
			// fmt.Println("\t no active edge")
			new_edge := &edge{start_index: cindex, end_index: text_length, end_node: nil}
			active_node.edges[ts] = new_edge
			remainder--
			cindex++
		} else {
			// fmt.Println("\t edge FOUNDDDDD")
			// fmt.Println("\t active edge - ", active_edge)
			// fmt.Println("\t active edge - ", *active_edge)
			// fmt.Println("\t active_string - ", active_string)
			// fmt.Println("\t active_length - ", active_length)

			// fmt.Println("\t next character in edge - ", string(text[active_edge.start_index + active_length]))
			// matches edge - could also continue from current match
			if ts == string(text[active_edge.start_index + active_length]) {
				// fmt.Println("\t char matched")
				active_string += ts
				active_length++
				remainder++
				cindex++

				// does node need to be updated ?
				if active_edge.start_index + active_length >= active_edge.end_index {
					// fmt.Println("\t EDGE END REACHED - update active_node")
					active_node = active_edge.end_node
					active_edge = nil
					active_string = ""
					active_length = 0
				}
			} else {
				//split this edge
				// fmt.Println("\t char did not match")
				new_node := new(node)
				new_node.edges = make(map[string]*edge)
				new_node.ntype = "inner"

				// add two edges to this node, one for the original, one at the split
				// fmt.Println("\t string mismatched at edge - ", string(text[active_edge.start_index + active_length]))
				// fmt.Println("\t current string - ", ts )
				new_node.edges[string(text[active_edge.start_index + active_length])] = &edge{start_index: active_edge.start_index + active_length, end_index: text_length, end_node: nil}
				new_node.edges[ts] = &edge{start_index: cindex, end_index: text_length, end_node: nil}

				// new_node.ntype = "inner"
				// fmt.Println("\t new node - ", new_node)
				active_edge.end_node = new_node
				active_edge.end_index = active_edge.start_index + active_length

				// after adding update active point
				// tedge, ok := active_node.edges[ts]
				// if ok {
				// 	active_edge = tedge
				// }

				if last_node != nil {
					// fmt.Println("\t last node not nil")
					last_node.suffix_link = new_node
				}

				// fmt.Println("\t before last_node - ", last_node)
				last_node = new_node
				// fmt.Println("\t active_string - ", active_string)

				var node_reset int = 0
				if active_node.ntype != "root" && active_node.suffix_link != nil {
					// fmt.Println("\t suffix link exists")
					active_node = active_node.suffix_link
					tedge, ok := active_node.edges[string(active_string[0])]
					if ok {
						active_edge = tedge
					} else {
						active_edge = nil
					}

					remainder--
					node_reset++
				} else if active_node.ntype != "root" {
					// fmt.Println("\t suffix link is not root")
					active_node = root
					tedge, ok := active_node.edges[string(active_string[0])]
					if ok {
						active_edge = tedge
					} else {
						active_edge = nil
					}

					remainder--
					node_reset++
				}

				if node_reset == 0 {
					if len(active_string) > 1 {
						// fmt.Println("\t active length > 1")
						tedge, ok := active_node.edges[string(active_string[1])]
						if ok {
							active_edge = tedge
						} else {
							active_edge = nil
						}
	
						active_string = active_string[1:]
						active_length--
						remainder--
						// if last_node != nil {
						// 	fmt.Println("\t last node not nil")
						// 	last_node.suffix_link = new_node
						// }
					} else {
						// fmt.Println("\t active length <= 1")
						active_string = ""
						active_length = 0
						active_edge = nil
						last_node = nil
						remainder = 0
					}
				}
			}
		}

		// fmt.Println("\t active_node - ", active_node)
		// fmt.Println("\t active edge - ", active_edge)
		// fmt.Println("\t active_string - ", active_string)
		// fmt.Println("\t active_length - ", active_length)
		// fmt.Println("\t remainder - ", remainder)
		// fmt.Println("\t last_node - ", last_node)
		// // fmt.Println("\t temp_edge - ", temp_edge)
		// fmt.Println("\t ROOT - ", root)
		// fmt.Println("END ITERATION")
	}
	return root
}

func main() {
	// tree := build_suffix_tree("gattaca")
	var text string = "abcabxabcd"
	fmt.Println("start building tree")
	var tree *node = build_suffix_tree(text)
	fmt.Println("finished building tree")

	fmt.Println("search tree")
	fmt.Println(tree.search("abcd", text, 3))
	// fmt.Println(build_suffix_tree("gattaca"))
}