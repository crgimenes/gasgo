package database

type Database interface {
	LoadSettings(dest any) error
}
