package arcgis

type ArcgisResponse[T any] struct {
	Features []ArcgisEvent[T] `json:"features"`
}

type ArcgisEvent[T any] struct {
	Attributes T              `json:"attributes" validate:"required"`
	Geometry   ArcgisGeometry `json:"geometry" validate:"required"`
}

type ArcgisGeometry struct {
	X float64 `json:"x" validate:"required"`
	Y float64 `json:"y" validate:"required"`
}
