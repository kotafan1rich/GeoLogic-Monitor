package maps

import "time"

type (
	Response struct {
		Version   float64   `json:"version"`
		Generator string    `json:"generator"`
		Osm3s     Osm3s     `json:"osm3s"`
		Elements  []Element `json:"elements"`
		Remark    string    `json:"remark,omitempty"`
	}

	Osm3s struct {
		TimestampOSMBase   time.Time `json:"timestamp_osm_base"`
		TimestampAreasBase time.Time `json:"timestamp_areas_base"`
		Copyright          string    `json:"copyright"`
	}

	Element struct {
		Type      string            `json:"type"`
		ID        int64             `json:"id"`
		Lat       *float64          `json:"lat,omitempty"`
		Lon       *float64          `json:"lon,omitempty"`
		Center    *LatLon           `json:"center,omitempty"`
		Bounds    *Bounds           `json:"bounds,omitempty"`
		Geometry  []LatLon          `json:"geometry,omitempty"`
		Nodes     []int64           `json:"nodes,omitempty"`
		Members   []Member          `json:"members,omitempty"`
		Tags      map[string]string `json:"tags,omitempty"`
		Timestamp *time.Time        `json:"timestamp,omitempty"`
		Version   *int              `json:"version,omitempty"`
		Changeset *int64            `json:"changeset,omitempty"`
		User      string            `json:"user,omitempty"`
		UID       *int64            `json:"uid,omitempty"`
	}

	LatLon struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}

	Bounds struct {
		MinLat float64 `json:"minlat"`
		MinLon float64 `json:"minlon"`
		MaxLat float64 `json:"maxlat"`
		MaxLon float64 `json:"maxlon"`
	}

	Member struct {
		Type     string   `json:"type"`
		Ref      int64    `json:"ref"`
		Role     string   `json:"role"`
		Lat      *float64 `json:"lat,omitempty"`
		Lon      *float64 `json:"lon,omitempty"`
		Geometry []LatLon `json:"geometry,omitempty"`
	}
)
