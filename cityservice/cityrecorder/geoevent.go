package cityrecorder

import (
	"encoding/json"
	"fmt"
	. "github.com/dimroc/urbanevents/cityservice/utils"
	"log"
	"strings"
	"time"
)

type GeoEvent struct {
	CityKey       string     `json:"city"`
	CreatedAt     time.Time  `json:"createdAt"`
	FullName      string     `json:"fullName"`
	GeoJson       GeoJson    `json:"geojson"`
	Hashtags      []string   `json:"hashtags"`
	Id            string     `json:"id"`
	MediaUrl      string     `json:"mediaUrl"`
	Link          string     `json:"link"`
	LocationType  string     `json:"locationType"`
	MediaType     string     `json:"mediaType"`
	Text          string     `json:"text,omitempty"`
	TextFrench    string     `json:"text_fr,omitempty"`
	Point         [2]float64 `json:"point"`
	Service       string     `json:"service"`
	ThumbnailUrl  string     `json:"thumbnailUrl"`
	Type          string     `json:"type"`
	Username      string     `json:"username"`
	Place         string     `json:"place"`
	Neighborhoods []string   `json:"neighborhoods"`
	ExpandedUrl   string     `json:"-"`
}

type GeoJson struct {
	Type           string           `json:"type"`
	CoordinatesRaw *json.RawMessage `json:"coordinates"` // Coordinate always has to have exactly 2 values
}

func (g GeoJson) Center() [2]float64 {
	return g.GenerateShape().Center()
}

type GeoShape interface {
	Center() [2]float64
}

type Point struct {
	GeoJson
	Coordinates [2]float64 `json:"coordinates"` // Coordinate always has to have exactly 2 values
}

func (p *Point) Center() [2]float64 {
	return p.Coordinates
}

type BoundingBox struct {
	GeoJson
	Coordinates [][][]float64 `json:"coordinates"`
}

func (bb *BoundingBox) Center() [2]float64 {
	// TODO: Get average
	center := [2]float64{bb.Coordinates[0][0][0], bb.Coordinates[0][0][1]}
	return center
}

func GeoJsonFrom(typeValue string, v interface{}) GeoJson {
	b, err := json.Marshal(v)
	Check(err)

	geojson := GeoJson{
		Type:           typeValue,
		CoordinatesRaw: (*json.RawMessage)(&b),
	}

	return geojson
}

func (g *GeoEvent) String() string {
	return fmt.Sprintf(
		"{CreatedAt: %s, GeoJson: %s, Point: %s, Id: %s, CityKey: %s, LocationType: %s, Type: %s, Text: %s}",
		g.CreatedAt,
		g.GeoJson,
		g.Point,
		g.Id,
		g.CityKey,
		g.LocationType,
		g.Type,
		g.Text,
	)
}

func (geojson *GeoJson) GenerateShape() GeoShape {
	var shape GeoShape
	var coordinatesDestination interface{}

	switch strings.ToLower(geojson.Type) {
	case "point":
		point := &Point{GeoJson: *geojson}
		coordinatesDestination = point.Coordinates
		shape = point
	case "polygon":
		box := &BoundingBox{GeoJson: *geojson}
		coordinatesDestination = &box.Coordinates
		shape = box
	}

	err := json.Unmarshal(*geojson.CoordinatesRaw, coordinatesDestination)
	if err != nil {
		log.Fatalln("error:", err)
	}

	return shape
}
