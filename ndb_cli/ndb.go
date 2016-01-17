package main

import (
	"ndb"
	"flag"
	"fmt"
	"os"
)

func main() {
	filename := flag.String("f", "" , "Use -f <ndb file path>")
	query := flag.String("q", "", "Use -q <ndb query>")
	help := flag.Bool("h", false, "Use -h")
	
	flag.Usage = PrintHelp
	flag.Parse()
	
	if *help {
		PrintHelp()
		return
	}
	
	// 当执行出现异常(panic)时, 输出错误信息 
	defer func() {
		if err := recover(); err != nil {  
			fmt.Println("Program Error")
			os.Exit(0)
		}
	}()
	
	if *filename != "" {
		node, err := ndb.Read(*filename)
		if err == nil {
			if *query != "" {
				result, found, err := ndb.Execute(node, *query)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(0)
				}
				if found {
					switch result.(type) { 
					case []*ndb.Node:
						list := result.([]*ndb.Node)
						for _, item := range list {
							PrintNode(item, 0)
						}
					case *ndb.Node:
						PrintNode(result.(*ndb.Node), 0)
					}
				} else {
					fmt.Printf("%s Found NOTHING\n", *query)	
				}
			} else {
				PrintNode(node, 0)
			}
		} else {
			fmt.Printf("Read file %s FAIL\n", *filename)
		}
	} else {
		PrintHelp()
	}
}

func PrintNode(node *ndb.Node, deep int) {
	
	tab := ""
	for i := 0 ; i < deep; i++ {
		tab += "   "
	}
	
	nodeName := node.GetName()
	if nodeName != "" {
		fmt.Println(tab + nodeName)
	}

	values := node.GetValues()
	for key, value := range values {
		for _, valueItem := range value {
			fmt.Println(tab + "|--" + key + " : " + valueItem)
		}
	}
	
	children := node.GetChileren()
	for _, child := range children {
		PrintNode(child, deep + 1)
	}
}

func PrintHelp() {
	fmt.Println("Usage of ndb_cli:")
	fmt.Println("-f\tndb file path\texample: /home/ndb/example.ndb")
	fmt.Println("-q\tndb query\texample: select:root->parent->child->name:jim")
	fmt.Println("-h\tshow help message")
}

