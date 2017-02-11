package algorithm

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/novikk/mpaths-alg/models"
)

func boundaries(pts models.Points) [2]models.Point {
	minLat, minLng := math.MaxFloat64, math.MaxFloat64
	maxLat, maxLng := -math.MaxFloat64, -math.MaxFloat64

	for _, pt := range pts {
		minLat = math.Min(minLat, pt.Lat)
		minLng = math.Min(minLng, pt.Lng)
		maxLat = math.Max(maxLat, pt.Lat)
		maxLng = math.Max(maxLng, pt.Lng)
	}

	return [2]models.Point{
		models.Point{minLat, minLng},
		models.Point{maxLat, maxLng},
	}
}

func RandomPoints(bounds [2]models.Point, k int) models.Points {
	res := make(models.Points, k)
	for i := 0; i < k; i++ {
		res[i] = models.Point{
			Lat: bounds[0].Lat + rand.Float64()*(bounds[1].Lat-bounds[0].Lat),
			Lng: bounds[0].Lng + rand.Float64()*(bounds[1].Lng-bounds[0].Lng),
		}
	}

	return res
}

func distance(pt1, pt2 models.Point) float64 {
	return math.Sqrt(math.Pow(pt1.Lat-pt2.Lat, 2) + math.Pow(pt1.Lng-pt2.Lng, 2))
}

func distanceInMeters(pt1, pt2 models.Point) float64 {
	// haversine formula
	p := math.Pi / 180
	a := 0.5 - math.Cos((pt2.Lat-pt1.Lat)*p)/2 + math.Cos(pt1.Lat*p)*math.Cos(pt2.Lat*p)*(1-math.Cos((pt2.Lng-pt1.Lng)*p))/2
	return 12742 * math.Asin(math.Sqrt(a)) * 1000
}

func Kmeans(pts models.Points, k int) models.Clusters {
	centroids := RandomPoints(boundaries(pts), k)
	var finalClusters models.Clusters

	repeated := 0

	for repeated < 50 {
		clusters := make(models.Clusters, k)
		for i, ct := range centroids {
			clusters[i].Centroid = ct
		}

		for _, pt := range pts {
			centroidForPt := 0
			dist := distance(pt, centroids[0])

			for i, ct := range centroids {
				distToCt := distance(ct, pt)
				if distToCt < dist {
					dist = distToCt
					centroidForPt = i
				}
			}

			clusters[centroidForPt].Pts = append(clusters[centroidForPt].Pts, pt)
		}

		finalClusters = clusters

		// calculate new centroids
		newCentroids := make(models.Points, k)
		changed := false
		for i := 0; i < k; i++ {
			totalLat := 0.0
			totalLng := 0.0
			for _, pt := range clusters[i].Pts {
				totalLat += pt.Lat
				totalLng += pt.Lng
			}

			newCentroids[i] = models.Point{totalLat / float64(len(clusters[i].Pts)), totalLng / float64(len(clusters[i].Pts))}
			if math.Abs(newCentroids[i].Lat-centroids[i].Lat) > 0.00001 || math.Abs(newCentroids[i].Lng-centroids[i].Lng) > 0.00001 {
				changed = true
			}
		}

		centroids = newCentroids
		if !changed {
			repeated = repeated + 1
		} else {
			repeated = 0
		}
	}

	return finalClusters
}

func KmeansMaxDist(pts models.Points, maxDistMeters float64) models.Clusters {
	// start with k=1 cluster until max dist is satisfied
	k := 1

	for true {
		fmt.Println("Trying k =", k)
		clusters := Kmeans(pts, k)

		currentMaxDist := 0.0
		for _, clust := range clusters {
			for _, pt := range clust.Pts {
				currentMaxDist = math.Max(currentMaxDist, distanceInMeters(pt, clust.Centroid))
			}
		}

		if currentMaxDist <= maxDistMeters {
			fmt.Println(currentMaxDist)
			return clusters
		}

		k = k + 1
	}

	return nil
}
