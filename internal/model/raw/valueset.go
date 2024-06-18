package raw

type ValueSet struct {
	ResourceType string         `json:"resourceType"`
	ID           string         `json:"id"`
	Text         Text           `json:"text"`
	Extension    []Extension    `json:"extension"`
	URL          string         `json:"url"`
	Version      string         `json:"version"`
	Name         string         `json:"name"`
	Title        string         `json:"title"`
	Status       string         `json:"status"`
	Publisher    string         `json:"publisher"`
	Contact      []Contact      `json:"contact"`
	Description  string         `json:"description"`
	Jurisdiction []Jurisdiction `json:"jurisdiction"`
	Copyright    string         `json:"copyright"`
	Compose      Compose        `json:"compose"`
}

type Text struct {
	Status string `json:"status"`
	Div    string `json:"div"`
}

type Jurisdiction struct {
	Coding []Coding `json:"coding"`
}

type Coding struct {
	System string `json:"system"`
	Code   string `json:"code"`
}

type Compose struct {
	Include []Include `json:"include"`
}

type Include struct {
	System string   `json:"system"`
	Filter []Filter `json:"filter"`
}

type Filter struct {
	Property string `json:"property"`
	Op       string `json:"op"`
	Value    string `json:"value"`
}
