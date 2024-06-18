package raw

type CodeSystem struct {
	ResourceType  string       `json:"resourceType"`
	ID            string       `json:"id"`
	Meta          MetaData     `json:"meta"`
	Extension     []Extension  `json:"extension"`
	URL           string       `json:"url"`
	Identifier    []Identifier `json:"identifier"`
	Version       string       `json:"version"`
	Name          string       `json:"name"`
	Title         string       `json:"title"`
	Status        string       `json:"status"`
	Experimental  bool         `json:"experimental"`
	Date          string       `json:"date"`
	Publisher     string       `json:"publisher"`
	Contact       []Contact    `json:"contact"`
	Description   string       `json:"description"`
	CaseSensitive bool         `json:"caseSensitive"`
	ValueSet      string       `json:"valueSet"`
	Content       string       `json:"content"`
	Concept       []Concept    `json:"concept"`
}

type MetaData struct {
	LastUpdated string `json:"lastUpdated"`
}

type Identifier struct {
	System string `json:"system"`
	Value  string `json:"value"`
}

type Contact struct {
	Telecom []Telecom `json:"telecom"`
}

type Telecom struct {
	System string `json:"system"`
	Value  string `json:"value"`
}

type Concept struct {
	Code       string    `json:"code"`
	Display    string    `json:"display"`
	Definition string    `json:"definition"`
	Concept    []Concept `json:"concept,omitempty"`
}
