package models

type Pokemon struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Type           []string `json:"type"`
	DetailPageURL  string   `json:"detailPageUrl"`
	ThumbnailImage string   `json:"ThumbnailImage"`
	Weight         float32  `json:"weight"` // lbs
	Height         float32  `json:"height"` // inch
	Generation     int      `json:"generation"`
	Region         string   `json:"region"`
	RegionLink     string   `json:"regionLink"`

	//Abilities        []string `json:"abilities"`
	//Weakness         []string `json:"weakness"`
	//Number           string   `json:"number"`
	//CollectiblesSlug string   `json:"collectibles_slug"`
	//Featured         string   `json:"featured"`
	//Slug             string   `json:"slug"`
	//ThumbnailAltText string   `json:"ThumbnailAltText"`
}
