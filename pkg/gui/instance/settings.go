package instance

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Settings struct {
	Database *DatabaseSettings
	WebAPI   *WebAPISettings
}

type DatabaseSettings struct {
	Dialect   string
	Parameter string
}

type WebAPISettings struct {
	Listen string
}

func (settings *Settings) Save() error {
	f, err := os.Create("dm-gui.toml")
	if err != nil {
		return err
	}
	return toml.NewEncoder(f).Encode(settings)
}

func LoadSettings() (*Settings, error) {
	var settings Settings
	f, err := os.Open("dm-gui.toml")
	if err != nil {
		return nil, err
	}
	err = toml.NewDecoder(f).Decode(&settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func GetDefaultSettings() *Settings {
	return &Settings{
		Database: &DatabaseSettings{Dialect: "", Parameter: ""},
		WebAPI:   &WebAPISettings{Listen: "0.0.0.0:8899"},
	}
}
