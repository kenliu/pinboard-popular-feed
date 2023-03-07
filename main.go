package main

import "fmt"

type bookmark struct {
	id    string
	title string
	url   string
}

func main() {
	popular := ScrapePinboardPopular()
	for i := 0; i < len(popular); i++ {
		fmt.Println(popular[i].id)
		fmt.Println(popular[i].title)
		fmt.Println(popular[i].url)
		println()
	}
}
