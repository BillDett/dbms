package main

import (
	"fmt"
	"os"

	"dbms/models"
)

func main() {

	var dm models.DataModel

	if err := dm.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer dm.Close()

	search := "something"

	// load the database if given an arg
	if len(os.Args) > 1 {
		if os.Args[1] == "--create" {
			if err := dm.LoadPersons(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if err := dm.LoadJournals(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if err := dm.LoadPostsFile(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println("Created database")
			os.Exit(0)
		} else {
			search = os.Args[1]
		}
	}

	if journals, err := dm.GetJournals(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		for _, j := range journals {
			fmt.Printf("%s\t%s\n", j.Name, j.Entry)
		}
	}

	if err := dm.SearchPosts(search); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
