package main

import (
	"ndb"
	"ndb/data"
	"flag"
	"fmt"
	"strings"
)

func main() {
	filename := flag.String("f", "" , "Use -f <ndb file path>")
	query := flag.String("q", "", "Use -q <ndb query>")
	help := flag.Bool("h", false, "Use -h <help message>")
	
	flag.Usage = PrintHelp
	flag.Parse()
	
	
	if *help {
		PrintHelp()
	}
	
	if *filename != "" && *query != "" {
		node, err := ndb.Read(*filename)
		if err == nil {
			result, found := ndb.Execute(node, *query)
			if found {
				switch result.(type) { 
				case []*data.Node:
					list := result.([]*data.Node)
					for _, item := range list {
						PrintNode(item)
					}
				case *data.Node:
					PrintNode(result.(*data.Node))
				}
			} else {
				fmt.Printf("%s Found NOTHING\n", *query)	
			}
		} else {
			fmt.Printf("Read file %s FAIL\n", *filename)	
		}
	} else {
		PrintHelp()
	}
}

func PrintNode(node *data.Node) {
	tab := "|--"
	
	fmt.Println(node.GetName())
	values := node.GetValues();
	for key, value := range values {
		fmt.Println(tab + key + " : " + strings.Join(value, ","))
	}
}

func PrintHelp() {
	fmt.Println("Usage of ndb_cli:")
	fmt.Println("-f\tndb file path\texample: /home/ndb/example.ndb")
	fmt.Println("-q\tndb query\texample: select:root->parent->child->name:jim")
	fmt.Println("-h\tshow help message")
}
