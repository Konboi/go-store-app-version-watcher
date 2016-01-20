package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

type Config struct {
	BotName      string `yaml:"bot_name"`
	IconEmoji    string `yaml:"icon_emoji"`
	MessageText  string `yaml:"message_text"`
	WebHookUri   string `yaml:"web_hook_uri"`
	GooglePlayId string `yaml:"google_play_id"`
	AppStoreId   string `yaml:"app_store_id"`
	DbPath       string `yaml:"db_path"`
}

type DBH struct {
	*sql.DB
}

func NewConfig(path string) (config Config, err error) {
	config = Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return config, err
	}

	if config.GooglePlayId == "" && config.AppStoreId == "" {
		return config, fmt.Errorf("Please Set Your App Id.")
	}

	db, err := sql.Open("sqlite3", config.DbPath)
	if err != nil {
		return config, err
	}

	err = db.Ping()
	if err != nil {
		return config, err
	}

	dbh = &DBH{db}

	if config.GooglePlayId != "" {
		uri := fmt.Sprintf("%s/store/apps/details?id=%s", GOOGLE_PLAY_BASE_URI, config.GooglePlayId)

		res, err := http.Get(uri)
		if err != nil {
			return config, err
		}
		if res.StatusCode == http.StatusNotFound {
			return config, fmt.Errorf("AppID: %s is not exists at Google Play", config.GooglePlayId)
		}
	}

	if config.AppStoreId != "" {
		uri := fmt.Sprintf("%s/%s", APP_STORE_BASE_URI, config.AppStoreId)

		res, err := http.Get(uri)
		if err != nil {
			return config, err
		}
		if res.StatusCode == http.StatusNotFound {
			return config, fmt.Errorf("AppID: %s is not exists at App Store", config.AppStoreId)
		}
	}

	return config, err
}
