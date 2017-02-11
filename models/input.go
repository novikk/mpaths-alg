package models

type Point struct {
	Lat, Lng float64
}

type Points []Point
type Cluster struct {
	Pts      []Point
	Centroid Point
}
type Clusters []Cluster
