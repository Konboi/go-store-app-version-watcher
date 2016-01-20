# go-store-app-version-watcher

go-store-app-version-watcher is an incoming webhook that posts app version when changed current version.

### Installation

1. `go get github.com/Konboi/go-store-app-version-watcher`

2. Create a SQLite database using `sqlite3 DB_NAME < ./schema.sql`

3. Copy `config_test.yml` and Add your webhook url, Google Play app id, AppStore app id and the path to your database.

```yaml
bot_name: "review_woman"
icon_emoji: ":ok_woman:"
message_text: "Changed your version"
web_hook_uri: "<SET YOUR SLACK WEB HOOK URL>"
google_play_app_id: "<SET YOUR Google Play App Id if you publish Google Play"
app_store_app_id: "<SET YOUR App Store App Id if you publish App Store"
db_path: "/tmp/version_watcher_db"
```

4. `go-store-app-version-watcher -c <EDIT CONFIG FILE PATH`

5. Run periodically using cron or some other job scheduler
