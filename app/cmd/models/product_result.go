package models

type ProductResult struct {
	Product
	StoreID    uint    `json:"storeId"`
	StoreName  string  `json:"storeName"`
	Lat        float64 `json:"lat"`
	Long       float64 `json:"long"`
	GeoHash    string  `json:"geoHash"`
	DistanceKm float64 `json:"distanceKm,omitempty"`
}
