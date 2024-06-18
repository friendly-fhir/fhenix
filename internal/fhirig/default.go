package fhirig

import (
	"os"
	"path/filepath"
)

// SystemCacheDir returns the system cache directory.
func SystemCacheDir() string {
	if dir := os.Getenv("FHENIX_CACHE_DIR"); dir != "" {
		return dir
	}
	if dir, err := os.UserHomeDir(); err == nil {
		return filepath.Join(dir, ".fhenix")
	}
	if dir := os.Getenv("HOME"); dir != "" {
		return filepath.Join(dir, ".fhenix")
	}
	// Fallback to temp directory if we absolutely can't find a home directory.
	// This should almost never happen in practice.
	return filepath.Join(os.TempDir(), ".fhenix")
}

// NewSystemCache creates a new PackageCache with the system cache directory.
func NewSystemCache() *PackageCache {
	return &PackageCache{
		Root:     SystemCacheDir(),
		Registry: "https://packages.simplifier.net",
	}
}
