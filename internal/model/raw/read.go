package raw

import (
	"encoding/json"
	"os"
)

func readJSON[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sd T
	if err := json.NewDecoder(file).Decode(&sd); err != nil {
		return nil, err
	}
	return &sd, nil
}

func ReadValueSet(path string) (*ValueSet, error) {
	return readJSON[ValueSet](path)
}

func ReadCodeSystem(path string) (*CodeSystem, error) {
	return readJSON[CodeSystem](path)
}
