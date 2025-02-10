package models

type APIKey struct {
	Key				string
	CanAccessSecret	bool
	CanAddAlbum		bool
	CanViewAlbum	bool
}
