package model

// CodeSystem represents a FHIR code system.
type CodeSystem struct {
	// Package is the name of the package that the code system is defined in.
	Package string

	// Version is the version of the package that the code system is defined in.
	Version string

	// Description is the full description of the code system.
	Description string

	// URL is the URL of the code system.
	URL string

	// Name is the name of the code system.
	Name string

	// Title is the title of the code system.
	Title string

	// Status is the publication status of the code system.
	Status string

	// Codes is a list of all the codes that are defined in the code system.
	Codes []Code
}

// Code represents a code that is defined in a code system.
type Code struct {
	// Value is the actual string code value.
	Value string

	// Display is the human readable display value.
	Display string

	// Definition is the definition of the code.
	Definition string
}
