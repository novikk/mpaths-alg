package algorithm

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/novikk/mpaths-alg/models"
)

func GetRoutes(pts models.Points) models.Routes {
	clusters := KmeansMaxDist(pts, 200)
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

	fmt.Println(stdin.String())

	// println(sum)

	path, _ := exec.LookPath("java")
	cmd := exec.Command(path, "-cp", "hackathon-vrp.jar:.", "hackathon.hackathon.RunVRP")
	cmd.Dir = "."
	cmd.Stdin = strings.NewReader(stdin.String())

	reNVehicles := regexp.MustCompile(`noVehicles\s+\|\s+(\d+)`)
	reServices := regexp.MustCompile(`\|\s+(\d+)\s+\|\s+vehicle\s+\|\s+service\s+\|\s+(\d+)\s+\|`)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(output[:]), err)
		return nil
	}

	numVehicles := 0
	//for _, line := range strings.Split(string(output[:]), "\n") {
	reResultNVehicles := reNVehicles.FindAllStringSubmatch(string(output[:]), -1)
	//fmt.Println(reResultNVehicles)
	if len(reResultNVehicles) > 0 {
		numVehicles, _ = strconv.Atoi(reResultNVehicles[0][1])
	}

	routes := make(models.Routes, numVehicles)
	reResultServices := reServices.FindAllStringSubmatch(string(output[:]), -1)
	// println(len(reResultServices))
	// for i := range reResultServices {
	// 	fmt.Println(reResultServices[i])
	// }
	if len(reResultServices) > 0 {
		for i := range reResultServices {
			vId, _ := strconv.Atoi(reResultServices[i][1])
			cId, _ := strconv.Atoi(reResultServices[i][2])

			routes[vId-1] = append(routes[vId-1], clusters[cId].Centroid)
		}
	}

	return routes
}
