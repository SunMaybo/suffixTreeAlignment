package main

import (
	"fmt"
	"bytes"
	"os"
	"encoding/json"
	"github.com/bradleyjkemp/memviz"
	// "github.com/davecgh/go-spew/spew"
)
// was very useful for debugging - 
// the code follows Ukkonen's algorithm from the lecture notes

// The code is structured as follows (TODO: draw a better ascii diagram)
// 1. Nodes have edges  * -------- () ------- * 
//                              |
//                              |
//							    *
// 2. Edges end in a node -----> ()
// 3. Once an edge is created, we assume it extends all the way to the 
//    length of the text, unless it gets split

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
	suffix_index int
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
	var active_string string = ""
	var active_length, remainder int = 0, 0
	var text_length = len(text)
	var cindex int = 0

	// build the tree for every character in the string
	for cindex < text_length {
		// fmt.Println(" ITERATING ", cindex)
		ts := string(text[cindex])
		// fmt.Println("\t character - ", cindex, ts)
		// fmt.Println("\t remainder - ", remainder)

		if remainder == 0 {
			// fmt.Println("\t remainder - ", remainder)
			// fmt.Println("\t active_node - ", active_node)
			_, ok := active_node.edges[ts]
			if ok {
				active_string = ts
				active_length = 1
				remainder = 1
				
				tedge, _ := active_node.edges[active_string]
				var active_edge *edge = tedge
				var end_index_temp int = active_edge.end_index
				if active_edge.end_index == text_length {
					end_index_temp = cindex
				}

				if end_index_temp - active_edge.start_index + 1 == active_length {
					active_edge, _ := active_node.edges[active_string]
					active_node = active_edge.end_node
					active_string = ""
					active_length = 0
				}
			} else {
				new_node := new(node)
				new_node.edges = make(map[string]*edge)
				new_node.ntype = "inner"

				new_edge := &edge{start_index: cindex, end_index: text_length, end_node: new_node}
				active_node.edges[string(text[cindex])] = new_edge
			}
		} else {
			// fmt.Println("\t remainder not zero - ", remainder)
			// fmt.Println("\t active_string - ", active_string)
			// fmt.Println("\t active_length - ", active_length)	

			if active_string == "" && active_length == 0 {
				_, ok := active_node.edges[ts]
				if ok {
					active_string = ts
					active_length = 1
					remainder++
				} else {
					remainder++
					// there's a mismatch, add this suffix to the tree
					remainder, active_node, active_string, active_length = propagate_suffix(root, text, cindex, remainder, active_node, active_string, active_length)
				}
			} else {
				tedge, _ := active_node.edges[active_string]
				var active_edge *edge = tedge
				// fmt.Println("\t tedge - ", tedge)	

				var end_index_temp int = active_edge.end_index
				if active_edge.end_index == text_length {
					end_index_temp = cindex
				}

				posi := active_edge.start_index + active_length

				// fmt.Println("\t posi - ", posi)	
				// fmt.Println("\t string(text[posi]) - ", string(text[posi]))

				if string(text[posi]) != ts {
					// fmt.Println("\t not same char ")
					remainder++
					// edge doesn't end with the same character, add this suffix to the tree
					remainder, active_node, active_string, active_length = propagate_suffix(root, text, cindex, remainder, active_node, active_string, active_length)
				} else {
					// tedge, _ := active_node.edges[active_string]
					// fmt.Println("\t same char ")
					// fmt.Println("tedge end index, ", active_edge.end_index)

					if posi < end_index_temp {
						active_length++
						remainder++
					} else {
						remainder++
						
						active_node = active_edge.end_node

						// fmt.Println("tedge end index, ", end_index_temp)

						if posi == end_index_temp {
							active_length = 0
							active_string = ""
						} else {
							active_length = 1
							active_string = ts
						}
					}
				}
			}
		}

		// if active_node == nil {
		// 	new_node := new(node)
		// 	new_node.edges = make(map[string]*edge)
		// 	new_node.ntype = "inner"

		// 	active_edge.end_node = new_node
		// 	active_node = new_node
		// }

		cindex++

		// fmt.Println("ALMOST END")
		// fmt.Println("active_node - ", active_node)
		// // fmt.Println("active edge - ", active_edge)
		// fmt.Println("active_string - ", active_string)
		// fmt.Println("active_length - ", active_length)
		// fmt.Println("remainder - ", remainder)
		// // fmt.Println("\t temp_edge - ", temp_edge)
		// fmt.Println("ROOT - ", root)
		// spew.Dump(root)
		// buf := &bytes.Buffer{}
		// memviz.Map(buf, root)
		// fmt.Println(buf.String())
		// fmt.Println("END ITERATION")
	}
	return root
}

// adds the current suffix to the tree (last character)
func propagate_suffix(root *node, text string, cindex int, remainder int, active_node *node, active_string string, active_length int) (int, *node, string, int) {
	var temp_node *node
	var last_node *node
	text_length := len(text)
	ts := string(text[cindex])
	
	// fmt.Println("\t propagate_suffix - ", ts)
	// fmt.Println("\t\t active_node - ", active_node)
	// fmt.Println("\t\t active_string - ", active_string)
	// fmt.Println("\t\t active_length - ", active_length)
	// fmt.Println("\t\t remainder - ", remainder)

	for remainder > 0 {
		// fmt.Println("\t\t\t remainder loop - ", remainder)
		suffixes := text[cindex - remainder + 1 : cindex + 1]
		// fmt.Println("\t\t\t suffixes - ", suffixes)
		active_length_suffixes := len(suffixes) - active_length - 1
		// fmt.Println("\t\t\t active_length_suffixes - ", active_length_suffixes)
        // fmt.Println("before - walk stree")

		active_node, active_string, active_length, active_length_suffixes = walk_stree(text, cindex, active_node, active_string, active_length, suffixes, active_length_suffixes)
		
		// fmt.Println("\t\t\t active_node - ", active_node)
		// fmt.Println("\t\t\t active_string - ", active_string)
		// fmt.Println("\t\t\t active_length - ", active_length)
		// fmt.Println("\t\t\t active_length_suffixes - ", active_length_suffixes)
		// fmt.Println("\t\t\t remainder - ", remainder)
        // fmt.Println("before - update_suffix")

		var ok bool
		ok, active_node, active_string, active_length, active_length_suffixes = update_suffix(text, cindex, active_node, active_string, active_length, suffixes, active_length_suffixes) 
		 
		// fmt.Println("\t\t\t active_node - ", active_node)
		// fmt.Println("\t\t\t active_string - ", active_string)
		// fmt.Println("\t\t\t active_length - ", active_length)
		// fmt.Println("\t\t\t active_length_suffixes - ", active_length_suffixes)
		// fmt.Println("\t\t\t remainder - ", remainder)

		// if a node is added in the current iteration
		// add a suffix link. only if the previous node was not a root
		if ok {
			if active_length == 1 && temp_node != nil && active_node.ntype != "root" {
				temp_node.suffix_link = active_node
			}
			return remainder, active_node, active_string, active_length
		}
		// fmt.Println("\t\t\t active_length - ", active_length)

		if active_length == 0 {
			_, oke := active_node.edges[string(suffixes[active_length_suffixes])]
			if !oke {
				new_node := new(node)
				new_node.edges = make(map[string]*edge)
				new_node.ntype = "inner"

				new_edge := &edge{start_index: cindex, end_index: text_length, end_node: new_node}
				active_node.edges[ts] = new_edge
				last_node = active_node
			}
		} else {
			tedge, _ := active_node.edges[active_string]
			if suffixes[active_length_suffixes + active_length] != text[tedge.start_index + active_length] {
				var active_edge *edge
				tedge, _ := active_node.edges[active_string]
				active_edge = tedge
				new_node := new(node)
				new_node.edges = make(map[string]*edge)
				new_node.ntype = "inner"
				// fmt.Println("\t\t\t make sure there's no existing edge already")
	
				// add two edges to this node, one for the original, one at the split
				// fmt.Println("\t\t\t string mismatched at edge - ", string(text[tedge.start_index + active_length]))
				// fmt.Println("\t\t\t current string - ", ts )
				new_node.edges[string(text[tedge.start_index + active_length])] = &edge{start_index: tedge.start_index + active_length, end_index: tedge.end_index, end_node: tedge.end_node}

				new_node2 := new(node)
				new_node2.edges = make(map[string]*edge)
				new_node2.ntype = "inner"

				new_node.edges[ts] = &edge{start_index: cindex, end_index: text_length, end_node: new_node2}
				// fmt.Println("\t\t\t new node - ", new_node)
				// fmt.Println("\t\t\t also update active edge - ")
				active_edge.end_node = new_node
				active_edge.end_index = active_edge.start_index + active_length - 1
				last_node = new_node
			} else {
				return remainder, active_node, active_string, active_length
			}
		}

		if temp_node != nil && last_node !=nil && last_node.ntype != "root" {
			temp_node.suffix_link = last_node
		}

		if last_node != nil && last_node.ntype != "root" {
			temp_node = last_node
		}

		if active_node.ntype == "root" && remainder > 1 {
			active_string = string(suffixes[1])
			active_length--
		}

		if active_node.suffix_link != nil {
			active_node = active_node.suffix_link
		} else {
			active_node = root
		}

		remainder--
	}
	// fmt.Println("\t\t remainder, active_node, active_string, active_length - ", remainder, active_node, active_string, active_length)

	return remainder, active_node, active_string, active_length
}

// walk the tree using active node, string & length to get to the last character
func walk_stree(text string, index int, active_node *node, active_string string, active_length int, suffixes string, remainder int) (*node, string, int, int) {
	// fmt.Println("\t\t walk_stree - ")
	// fmt.Println("\t\t\t suffixes - ", suffixes)
	text_length := len(text)

	if active_length == 0 || active_string == "" {
		return active_node, active_string, active_length, remainder
	}
	// fmt.Println("\t\t\t active_node - ", active_node)
	// fmt.Println("\t\t\t active_string - ", active_string)

	tedge, _ := active_node.edges[active_string]
	var active_edge *edge = tedge
	var end_index_temp int = active_edge.end_index
	if active_edge.end_index == text_length {
		end_index_temp = index
	}

	edge_length := end_index_temp - tedge.start_index + 1

	for active_length > edge_length {
		tedge, _ := active_node.edges[active_string]
		var active_edge *edge = tedge
		active_node = tedge.end_node
		remainder += edge_length
		active_string = string(suffixes[remainder])
		active_length -= edge_length
		
		var ok bool
		tedge, ok = active_node.edges[active_string]
		active_edge = tedge

		var end_index_temp int = active_edge.end_index
		if active_edge.end_index == text_length {
			end_index_temp = index
		}

		edge_length = end_index_temp - tedge.start_index + 1

		if ok {
			// fmt.Println("\t\t\t oksih exists")
		}
	}

	if active_length == edge_length {
		tedge, _ := active_node.edges[active_string]
		active_node = tedge.end_node
		active_string = ""
		active_length = 0
		remainder += edge_length
	}

	// fmt.Println("\t\t\t active_node, active_string, active_length, remainder - ", active_node, active_string, active_length, remainder)

	return active_node, active_string, active_length, remainder
}

// This adds the suffix to the tree
func update_suffix(text string, cindex int, active_node *node, active_string string, active_length int, suffixes string, remainder int) (bool, *node, string, int, int) {
	suffix_update := suffixes[remainder:]
	text_length := len(text)

	// fmt.Println("\t\t update_suffix - ")
	// fmt.Println("\t\t\t suffix_update - ", suffix_update)

	if active_length > 0 {
		tedge, _ := active_node.edges[active_string]
		var active_edge *edge = tedge
		var end_index_temp int = active_edge.end_index
		if active_edge.end_index == text_length {
			end_index_temp = cindex
		}
		edge_suffix := text[active_edge.start_index: end_index_temp + 1]

		// fmt.Println("\t\t\t edge Suffix", edge_suffix)
		// fmt.Println("\t\t\t edge Suffix starts", edge_suffix[0 : len(suffix_update)])

		if edge_suffix[0 : len(suffix_update)] == suffix_update {
			active_length = len(suffix_update)
			active_string = string(suffix_update[0])
			// fmt.Println("\t\t\t true, active_node, active_string, active_length, remainder - ", true, active_node, active_string, active_length, remainder)
			return true, active_node, active_string, active_length, remainder
		}
	} else {
		if remainder < len(suffixes) {
			tedge, ok := active_node.edges[string(suffixes[remainder])]
			if ok {
				var active_edge *edge = tedge

				var end_index_temp int = active_edge.end_index
				if active_edge.end_index == text_length {
					end_index_temp = cindex
				}

				edge_suffix := text[tedge.start_index: end_index_temp + 1]

				// fmt.Println("\t\t\t edge Suffix", edge_suffix)
				// fmt.Println("\t\t\t edge Suffix starts", edge_suffix[0 : len(suffix_update)])

				if edge_suffix[0 : len(suffix_update)] == suffix_update {
					active_length = len(suffix_update)
					active_string = string(suffix_update[0])
					// fmt.Println("\t\t\t true, active_node, active_string, active_length, remainder - ", true, active_node, active_string, active_length, remainder)
					return true, active_node, active_string, active_length, remainder
				}
			}
		}
	}

	// fmt.Println("\t\t\t false, active_node, active_string, active_length, remainder - ", false, active_node, active_string, active_length, remainder)

	return false, active_node, active_string, active_length, remainder
}

// Walks suffix tree and adds indices to the leaf node
func add_suffix_indices(tree *node, tree_height int, text_length int) { 

	if len(tree.edges) == 0 {
		// fmt.Println("tree len is 0 ")
		tree.suffix_index = text_length - tree_height + 1
		return
	} else {
		for _, child := range tree.edges { 
			add_suffix_indices(child.end_node, tree_height + child.end_index - child.start_index + 1, text_length)
			// fmt.Println("on node ", nuc, child)
		}
	}
} 

// Get all indices from a given node
func get_suffix_indices(tree *node, indices []int) []int {

	if len(tree.edges) == 0 {
		indices = append(indices, tree.suffix_index)
	} else {
		for _, child := range tree.edges { 
			indices = get_suffix_indices(child.end_node, indices)
		}
	}

	return indices
}

// Search for a pattern
// query string : query to search
// text string : original text used in constructing the suffix tree
func (stree *node) search(query, text string) map[string][]int {
	found := make(map[string][]int)
	var active_node *node = stree
	var active_edge * edge
	var active_length int = 0
	// var match bool = true
	var qindex int = 0
	// var text_length = len(text)
	// search the branch where the node happens
	for qindex < len(query) {
		// tq := string(query[qindex])
		// fmt.Println("query - ", qindex, tq)
		// fmt.Println("active_edge - ", active_edge)
		// fmt.Println("active_node - ", active_node)
		// fmt.Println("active_length - ", active_length)

		if active_edge == nil {
			tedge, ok := active_node.edges[string(query[qindex - active_length])]

			if ok {
				// fmt.Println("yo matched!")

				// we have a match
				active_edge = tedge
				active_length++
				qindex++

				if active_length == active_edge.end_index - active_edge.start_index + 1 {
					// fmt.Println("reached edge end!")
					active_node = active_edge.end_node
					active_edge = nil
					active_length = 0
				}
			} else {
				break
			}
		} else {

			// fmt.Println("active edge not nil!")
			// fmt.Println(active_edge)
			// fmt.Println(active_edge.start_index)
			// fmt.Println(active_length)
			// fmt.Println(text)

			// fmt.Println("char in edge - ", string(text[active_edge.start_index + active_length]))
			// fmt.Println("char in query -", string(query[qindex]))

			if string(query[qindex]) == string(text[active_edge.start_index + active_length]) {
				active_length++
				qindex++

				if active_length == active_edge.end_index - active_edge.start_index + 1 {
					active_node = active_edge.end_node
					active_edge = nil
					active_length = 0
				}
			} else {
				break
			}
		}
	}

	if active_edge != nil {
		active_node = active_edge.end_node
	}

	// depth first traverse to get all index positions for this substring
	var indices []int
	indices = get_suffix_indices(active_node, indices)

	found[query[0:qindex]] = indices
	return found
}

// Search for a pattern for all suffixes of the query
// min_match_length int : if entire query doesn't match, report all lengths above
//        this threshold.
func (stree *node) multi_search(query, text string, min_match_length int) map[string]map[string][]int {
	matched_indices := make(map[string]map[string][]int)

	for i := 0; i < len(query); i++ {
		fmt.Println(query[i:])

		tempind := stree.search(query[i:], text)

		fmt.Println(tempind)

		for tk, tv := range tempind {
			// fmt.Println("tk, tv", tk, tv)
			if len(tk) > min_match_length {
				_, ok := matched_indices[tk]
				if ok {
					// for _, vt := range val {
					// 	matched_indices[tk]["query_index"] = append(matched_indices[tk]["query_index"], i)
					// 	matched_indices[tk]["text_index"] = append(matched_indices[tk]["text_index"], vt)
					// }
				} else {
					tmp_val := map[string][]int{"query_index": []int{i}, "text_index": tv }
					matched_indices[tk] = tmp_val
				}
			}
		}
	}

	return matched_indices
}

// main functions
// TODO: accept command line arguments
func main() {

	// first argument is the name of the binary
	args := os.Args[1:]


	// arg1 should be the text to build a suffix tree from
	text := args[0]
	var tree *node = build_suffix_tree(text)
	add_suffix_indices(tree, 0, len(text))

	buf2 := &bytes.Buffer{}
	memviz.Map(buf2, tree)
	fmt.Println(buf2.String())

	// arg2 should be the query to search (optional)
	if len(args) > 1 {
		query := args[1]
		fmt.Println("query is ", query)
		matches := tree.multi_search(query, text, 2)
		match_json, _ := json.Marshal(matches)
		fmt.Println("query matches - ", string(match_json))
	}

	// // var text string = "gattaca$"
	// // var text string = "abcabxabcd$"
	// // fmt.Println("start building tree")
	// // var tree *node = build_suffix_tree(text)
	// // fmt.Println("finished building tree")

	// // fmt.Println("search tree")
	// // fmt.Println(tree.search("abce", text, 3))
	// // var tree *node = build_suffix_tree("gattaca")
	// // var text string = "bananasbanananananananananabananas$"
	// // var text string = "abcdefabxybcdmnabcdex"
	// var text string = "mississippi@tississilli$"
	// // var text string = "almasamolmaz"
	// // var text string = "cdddcdc"
	// // var text string = "dedododeeodoeodooedeeododooodoede"
	// // var text string = "panamabananas"
	// // var text string = "GAGACCTCGATCGGCCAACTCATCTGTGAAACGTCAGTCATTGTAAGACTGGACATTTAGGAAGTAAGCCTTTTTCTTATAGCCAATCCCGCTTTCAATTGAACGGCTAAACGAAGGTCGTTTGCCACTGATTAGCAATTGGTCCGGTGAAAAATTGTGTATTTTGGAAAGATGTAATCCTGCGAGACCTCGATCGGC"

	// // fmt.Println("start building tree")
	// var tree *node = build_suffix_tree(text)
	// // fmt.Println("finished building tree")
	// // spew.Dump(tree)

	// // buf := &bytes.Buffer{}
	// // memviz.Map(buf, tree)
	// // fmt.Println(buf.String())

	// add_suffix_indices(tree, 0, len(text))

	// buf2 := &bytes.Buffer{}
	// memviz.Map(buf2, tree)
	// fmt.Println(buf2.String())

	// // fmt.Println("search tree for get all possible indexes where the suffix matches")
	// // var query string = "issp"
	// // fmt.Println(query)
	// // fmt.Println(query[0:len(query)])
	
	// // matched_indices := tree.multi_search(query, text, 2)

	// // fmt.Println("matches - ", matched_indices)
}