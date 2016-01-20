package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Konboi/go-store-app-version-watcher/scraper"
)

type SlackPayload struct {
	Text      string `json:"text"`
	UserName  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
}

const (
	GOOGLE_PLAY_BASE_URI = "https://play.google.com"
	APP_STORE_BASE_URI   = "https://itunes.apple.com/jp/app"
	INIT_APP_VERSION     = "0.0.1"
	PLATFORM_APP_STORE   = "app_store"
	PLATFORM_GOOGLE_PLAY = "google_play"
)

var (
	dbh        *DBH
	configFile = flag.String("c", "config.yml", "config file")
)

func main() {
	flag.Parse()

	config, err := NewConfig(*configFile)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("version check start")

	if config.AppStoreId != "" {
		log.Println("app store version check start")
		version, err := scraper.AppStore(config.AppStoreId)
		if err != nil {
			log.Println(err)
			return
		}

		current_version, err := CurrentVersion(PLATFORM_APP_STORE)
		if err != nil {
			log.Println(err)
			return
		}

		if !(strings.Contains(version, current_version)) {
			err := SaveVersion(PLATFORM_APP_STORE, version)
			if err != nil {
				log.Println(err)
				return
			}

			err = PostSlack(config, PLATFORM_APP_STORE, version)
			if err != nil {
				log.Println(err)
				return
			}
		}
		log.Println("app store version check done")
	}

	if config.GooglePlayId != "" {
		log.Println("google play version check start")
		version, err := scraper.GooglePlay(config.GooglePlayId)
		if err != nil {
			log.Println(err)
			return
		}

		current_version, err := CurrentVersion(PLATFORM_GOOGLE_PLAY)
		if err != nil {
			log.Println(err)
			return
		}

		if !(strings.Contains(version, current_version)) {
			err := SaveVersion(PLATFORM_GOOGLE_PLAY, version)
			if err != nil {
				log.Println(err)
				return
			}

			err = PostSlack(config, PLATFORM_GOOGLE_PLAY, version)
			if err != nil {
				log.Println(err)
				return
			}
		}
		log.Println("google play version check done")
	}

	log.Println("version check done")
}

func CurrentVersion(platform string) (current_version string, err error) {
	row := dbh.QueryRow("SELECT version FROM app_version WHERE platform = ? LIMIT 1", platform)
	err = row.Scan(&current_version)

	if err != nil {
		// insert init data
		if err.Error() == "sql: no rows in result set" {
			_, _ = dbh.Exec("INSERT INTO app_version (platform, version) VALUES (?, ?)",
				platform, INIT_APP_VERSION)

			return INIT_APP_VERSION, nil
		}
	}

	return current_version, err
}

func SaveVersion(platform string, version string) (err error) {
	_, err = dbh.Exec("UPDATE app_version SET version = ? WHERE platform = ?",
		version, platform)

	if err != nil {
		return err
	}

	return nil
}

func PostSlack(config Config, platform string, version string) (err error) {
	slackPayload := SlackPayload{
		UserName:  config.BotName,
		IconEmoji: config.IconEmoji,
		Text:      fmt.Sprintf("%s %s:%s", config.MessageText, platform, version),
	}

	payload, err := json.Marshal(slackPayload)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", config.WebHookUri, bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
