package main

import (
	"fmt"
	"math"
	"bytes"
	"github.com/bradleyjkemp/memviz"
	"github.com/davecgh/go-spew/spew"
)
// was very useful for debugging - 


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
	var FLAG_END bool = false

	for {
		ts := string(text[cindex])

		fmt.Println("Iterating - ",cindex, ts)

		if remainder == 0 {
			fmt.Println("\t remainder is 0")
			remainder = 1
			active_string = ""
			active_length = 0
			last_node = nil
			tedge, ok := active_node.edges[ts]
			if ok {
				active_edge = tedge
			}
		} else if active_edge == nil {
			fmt.Println("\t remainder not 0 but active edge is nil")
			if active_length > 0 {
				tedge, ok := active_node.edges[string(active_string[0])]
				if ok {
					active_edge = tedge
				}
			} else {
				tedge, ok := active_node.edges[ts]
				if ok {
					active_edge = tedge
				}	
			}
		}

		fmt.Println("\t active_node - ", active_node)
		fmt.Println("\t after remainder")
		fmt.Println("\t active edge - ", active_edge)

		if active_edge == nil {
			fmt.Println("\t no active edge")
			new_edge := &edge{start_index: cindex, end_index: text_length, end_node: nil}
			active_node.edges[ts] = new_edge
			remainder--

			fmt.Println("\t add a new edge")
			last_node = nil

			if active_length > 0 {
				active_length--
			}

			if active_node.ntype != "root" {
				fmt.Println("\t Follow suffix link if not root and exists")
				if active_node.suffix_link != nil {
					active_node = active_node.suffix_link
				} else {
					active_node = root
				}
			} else {
				fmt.Println("\t Root node, check next char")
				cindex++
			}

			active_edge = nil
		} else {
			fmt.Println("\t edge FOUNDDDDD")
			fmt.Println("\t active edge - ", active_edge)
			fmt.Println("\t active edge - ", *active_edge)
			fmt.Println("\t active_string - ", active_string)
			fmt.Println("\t active_length - ", active_length)

			fmt.Println("\t skip all the way to the last edge ?")
			for (active_edge.end_index - active_edge.start_index) <= active_length {
				fmt.Println("\t skipping - length of edge < active_length")
				active_length = active_length - (active_edge.end_index - active_edge.start_index)
				// if active_length > 0 {
				// 	active_string = text[cindex - active_length: active_length]
				// } else {
				// 	active_string = ts
				// }
				// active_string = text[cindex - active_length: active_length]
				active_node = active_edge.end_node
				tedge, ok := active_node.edges[string(text[cindex - active_length])]
				if ok {
					active_edge = tedge
				} else {
					// active_edge = nil
					break
				}
			}

			fmt.Println("\t no more skips length of edge > active_length ")
			fmt.Println("\t active edge - ", active_edge)
			fmt.Println("\t active edge - ", *active_edge)
			fmt.Println("\t active_string - ", active_string)
			fmt.Println("\t active_length - ", active_length)

			if active_edge != nil {
				fmt.Println("\t next character in edge - ", string(text[active_edge.start_index + active_length]))
				if ts == string(text[active_edge.start_index + active_length]) {
					fmt.Println("\t char matched")
					cindex++
					remainder++
					active_length++
					active_string += ts
					last_node = nil

					if active_edge.end_index != text_length {
						if active_edge.end_index - active_edge.start_index <= active_length {
							fmt.Println("\t length of edge <= active_length ")

							active_node = active_edge.end_node
							fmt.Println("\t active edge - ", active_edge)
							fmt.Println("\t active_string - ", active_string)
							fmt.Println("\t active_length - ", active_length)

							// if active_length > 0 {
							// 	active_string = text[(cindex - active_length): active_length]
							// } else {
							// 	active_string = ts
							// }
							active_length = active_length - (active_edge.end_index - active_edge.start_index)

							// active_string = ts
							if cindex != text_length {
								tedge, ok := active_node.edges[string(text[cindex - active_length])]
								if ok {
									active_edge = tedge
								} else {
									active_edge = nil
								}
							}
						}
					}	
				} else {
					fmt.Println("\t no active edge, so no char did not match")
	
					new_node := new(node)
					new_node.edges = make(map[string]*edge)
					new_node.ntype = "inner"
					fmt.Println("\t make sure there's no existing edge already")
		
					// add two edges to this node, one for the original, one at the split
					fmt.Println("\t string mismatched at edge - ", string(text[active_edge.start_index + active_length]))
					fmt.Println("\t current string - ", ts )
					new_node.edges[string(text[active_edge.start_index + active_length])] = &edge{start_index: active_edge.start_index + active_length, end_index: active_edge.end_index, end_node: active_edge.end_node}
					new_node.edges[ts] = &edge{start_index: cindex, end_index: text_length, end_node: nil}
					fmt.Println("\t new node - ", new_node)
					fmt.Println("\t also update active edge - ")
					active_edge.end_node = new_node
					active_edge.end_index = active_edge.start_index + active_length
	
					if last_node != nil && remainder != 0{
						fmt.Println("\t last node not nil, add suffix link")
						last_node.suffix_link = new_node
					}
					fmt.Println("\t before last_node - ", last_node)
					last_node = new_node
					fmt.Println("\t active_string - ", active_string)			
					
					remainder--
					fmt.Println("\t active_node - ", active_node)
					fmt.Println("\t active edge - ", active_edge)
					fmt.Println("\t active_string - ", active_string)
					fmt.Println("\t active_length - ", active_length)
					fmt.Println("\t remainder - ", remainder)
					fmt.Println("\t last_node - ", last_node)
					if active_node.ntype == "root" {
						active_length--
						// active_string = string(text[cindex - active_length])
						active_string = active_string[1:]
						active_edge = nil
					} else {
						fmt.Println("\t Follow suffix link if not root and exists")
						active_string = active_string[1:]
						// active_length--
						if active_node.suffix_link != nil {
							fmt.Println("\t suffix link exists")
							active_node = active_node.suffix_link

							tedge, ok := active_node.edges[string(text[cindex-1])]
							if ok {
								active_edge = tedge
							} else {
								active_edge = nil
							}	
						} else {
							active_node = root

							if active_length > 0 {
								tedge, ok := active_node.edges[string(active_string[0])]
								if ok {
									active_edge = tedge
								}
							} else {
								tedge, ok := active_node.edges[ts]
								if ok {
									active_edge = tedge
								}	
							}

							// tedge, ok := active_node.edges[string(active_string[0])]
							// if ok {
							// 	active_edge = tedge
							// } else {
							// 	active_edge = nil
							// }	
						}

					}
				}
			} 
		}

		fmt.Println("\t ALMOST END")
		fmt.Println("\t active_node - ", active_node)
		fmt.Println("\t active edge - ", active_edge)
		fmt.Println("\t active_string - ", active_string)
		fmt.Println("\t active_length - ", active_length)
		fmt.Println("\t remainder - ", remainder)
		fmt.Println("\t last_node - ", last_node)
		// fmt.Println("\t temp_edge - ", temp_edge)
		fmt.Println("\t ROOT - ", root)
		spew.Dump(root)
		buf := &bytes.Buffer{}
		memviz.Map(buf, root)
		fmt.Println(buf.String())
		fmt.Println("END ITERATION")

		if cindex >= text_length {

			// if remainder != 0 {
			// 	fmt.Println("remainder is not 0, cleanup!")
			// 	FLAG_END = true
			// } else {
				break
			// }

			cindex = text_length - 1
			// cindex = text_length-1
		}

		fmt.Println("END ITERATION ?", FLAG_END)

        if FLAG_END {
            fmt.Println("END ITERATION ?", FLAG_END)
        }
	}
	return root
}

func main() {
	// var text string = "gattaca$"
	// var text string = "abcabxabcd$"
	// fmt.Println("start building tree")
	// var tree *node = build_suffix_tree(text)
	// fmt.Println("finished building tree")

	// fmt.Println("search tree")
	// fmt.Println(tree.search("abce", text, 3))
	// var tree *node = build_suffix_tree("gattaca")
	// var text string = "bananasbanananananananananabananas"
	// var text string = "abcdefabxybcdmnabcdex"
	// var text string = "mississippi"
	// var text string = "almasamolmaz"
	var text string = "cdddcdc"
	// var text string = "dedododeeodoeodooedeeododooodoede"
	// var text string = "panamabananas"
	// var text string = "GAGACCTCGATCGGCCAACTCATCTGTGAAACGTCAGTCATTGTAAGACTGGACATTTAGGAAGTAAGCCTTTTTCTTATAGCCAATCCCGCTTTCAATTGAACGGCTAAACGAAGGTCGTTTGCCACTGATTAGCAATTGGTCCGGTGAAAAATTGTGTATTTTGGAAAGATGTAATCCTGCGAGACCTCGATCGGC$"

	// fmt.Println("start building tree")
	var tree *node = build_suffix_tree(text)
	// fmt.Println("finished building tree")
	
	// spew.Dump(tree)

	// buf := &bytes.Buffer{}
	// memviz.Map(buf, tree)
	// fmt.Println(buf.String())

	fmt.Println(tree)

	// fmt.Println("search tree")
	// fmt.Println(tree.search("AGACTGG", text, 3))

}