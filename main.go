package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/novikk/mpaths-alg/algorithm"
	"github.com/novikk/mpaths-alg/models"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	pts := algorithm.RandomPoints([2]models.Point{
		models.Point{41.542137, 2.426475},
		models.Point{41.538419, 2.451129},
	}, 100)

	routes := algorithm.GetRoutes(pts)
	fmt.Println(routes)
}
