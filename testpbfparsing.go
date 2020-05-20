package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"math"
	"strconv"

	"github.com/qedus/osmpbf"
)

func main(){
	lat, err := strconv.ParseFloat(os.Args[2], 64)
	lon, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		return
	}
	parsePBF(os.Args[1], lat, lon)
}

func parsePBF(filename string, lat, lon float64){
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	var shortestDistance float64
	var cityName string

	for {
			if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				// Process Node v.
				tagMap := v.Tags
				place, placeKeyErr := tagMap["place"]
				if placeKeyErr {
					if place == "muncipality" || place == "city" || place == "town"{
						elat := v.Lat
						elon := v.Lon
						currentDistance := getDistanceFromLatLon(lat, lon, elat, elon)
						//fmt.Println(tagMap["name"], v.Lat, v.Lon, currentDistance)
						if ((shortestDistance == 0) || (shortestDistance > currentDistance)) {
							shortestDistance = currentDistance
							cityName = tagMap["name"]
						}
					}
				}
			}
		}
	}
	fmt.Printf("Nearest City is %v at a distance of %v from %v, %v\n", cityName, shortestDistance, lat, lon)
}

func getDistanceFromLatLon(lat1, lon1, lat2, lon2 float64) float64 {
	R := 6371.0
	dLat := deg2rad(lat2- lat1)
	dLon := deg2rad(lon2 - lon1)
	a := math.Sin(dLat/2) * math.Sin(dLat/2) + math.Cos(deg2rad(lat1)) * math.Cos(deg2rad(lat2)) * math.Sin(dLon/2) * math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	d := R * c
	return d
}

func deg2rad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
