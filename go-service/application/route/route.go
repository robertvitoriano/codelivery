package route

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"
)

// Route represents a request of new delivery request
type Route struct {
	ID        string     `json:"routeId"`
	ClientID  string     `json:"clientId"`
	Positions []Position `json:"position"`
}

// Position is a type which contains the lat and long
type Position struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

// PartialRoutePosition is the actual response which the system will return
type PartialRoutePosition struct {
	ID       string    `json:"routeId"`
	ClientID string    `json:"clientId"`
	Position []float64 `json:"position"`
	Finished bool      `json:"finished"`
}

func(route *Route) LoadPositions() error {
	if route.ID == ""{
		return errors.New("Route Id not informed")
	}

	file, error := os.Open("destinations/" + route.ID +".txt")

	if error != nil {
		return error
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		data := strings.Split(scanner.Text(), ",")

		lat, error := strconv.ParseFloat(data[0],64)

		if error != nil {
			return error
		}

		long, error := strconv.ParseFloat(data[1],64)

		if error != nil {
			return error
		}

		route.Positions = append(route.Positions, Position{
			Lat:  lat,
			Long: long,
		})
	}

	return nil
}

func (r *Route) ExportJsonPosition() ([]string, error){
	var route PartialRoutePosition
	var result []string
	total := len(r.Positions)

	for index, value := range r.Positions {
		route.ID = r.ID
		route.ClientID = r.ClientID
		route.Position = []float64 {value.Lat, value.Long}
		route.Finished = false

		if index == total {
			route.Finished = true
		}
		jsonRoute, error := json.Marshal(route)
		if error != nil {
			return nil, error
		}
		result = append(result, string(jsonRoute))
	}

	return result, nil

}