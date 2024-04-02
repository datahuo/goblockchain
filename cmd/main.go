package main

import (
	"fmt"

	"github.com/erdincmutlu/goblockchain/utils"
)

func main() {
	// fmt.Println("127.0.0.1:5000", utils.IsHostFound("127.0.0.1", 5000))
	fmt.Println(utils.FindNeighbours("127.0.1.1", 5000, 0, 3, 5000, 5003))
	fmt.Println(utils.GetHost())
}
