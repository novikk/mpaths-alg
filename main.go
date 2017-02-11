package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"strconv"

	"github.com/novikk/mpaths-alg/algorithm"
	"github.com/novikk/mpaths-alg/models"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	pts := algorithm.RandomPoints([2]models.Point{
		models.Point{41.542137, 2.426475},
		models.Point{41.538419, 2.451129},
	}, 100)

	clusters := algorithm.KmeansMaxDist(pts, 200)
	for _, c := range clusters {
		fmt.Println(c.Centroid, "-->", c.Pts, "(", c.Radius, ")")
	}

	// call algorithm in java
	// build stdin
	var stdin bytes.Buffer
	stdin.WriteString(strconv.Itoa(len(clusters)) + "\n")

	// sum := 0
	for i := range clusters {
		latStr := strconv.FormatFloat(clusters[i].Centroid.Lat, 'f', -1, 64)
		lngStr := strconv.FormatFloat(clusters[i].Centroid.Lng, 'f', -1, 64)

		stdin.WriteString(latStr + " " + lngStr + " " + strconv.Itoa(len(clusters[i].Pts)) + "\n")
		// sum += len(clusters[i].Pts)
	}

	// fmt.Println(stdin.String())
	// println(sum)

	path, _ := exec.LookPath("java")
	cmd := exec.Command(path, "-cp", "vrp-hackathon.jar:.", "hackathon.hackathon.RunVRP")
	cmd.Dir = "."
	cmd.Stdin = strings.NewReader(stdin.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output[:]), err)
		return
	}
}
