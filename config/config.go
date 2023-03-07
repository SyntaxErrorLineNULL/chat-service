package config

type Config struct {
	Redis struct {
		Address    []string
		Password   string
		ClientName string
	}
	Database struct {
		URL          string
		DatabaseName string
	}
}
