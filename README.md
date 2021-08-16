![Logo](assets/logo.png)

[![](https://gitlab.com/haath/goirate/badges/master/pipeline.svg)](https://gitlab.com/haath/goirate/pipelines)
[![](https://gitlab.com/haath/goirate/badges/master/coverage.svg)](https://gitlab.com/haath/goirate/-/jobs/artifacts/master/browse?job=test)
[![](https://goreportcard.com/badge/gitlab.com/haath/goirate)](https://goreportcard.com/report/gitlab.com/haath/goirate)
[![](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![](https://img.shields.io/github/release/gmantaos/Goirate.svg)](https://github.com/gmantaos/Goirate/releases)
[![](https://api.codeclimate.com/v1/badges/40e202ffa0f346797f35/maintainability)](https://codeclimate.com/github/gmantaos/Goirate)

Watching a lot of movies and series, it quickly became difficult to keep track of what was coming out and when.
Not to call managing a few torrents tiresome, but when there's 5-10 new episodes for things you watch coming out
every week, you begin to wonder if all of these extra clicks are really necessary. This also refers to wasted clicks,
for when you know that an episode aired but you don't quite know yet if there's a torrent out for the 1080p version you prefer.
With all this in mind, I first attempted to automate this with a simple python script, which would run as a cron job, crawl the Pirate Bay
for torrents of new episodes, send me the ones it finds via e-mail and update the list of series so that it would begin watching out
for the next episode. And the funny thing that script worked like clockwork, monitoring at least 40 different series over a period of two years.

Then the point came, when I wanted more features and automation, and that old python script was never written to be particularly scalable.
So I began development of Goirate.
This tool aims to become an all-in-one suite for automating your every piratey need.
It works as a CLI program, which is designed to go through the internet searching for torrents much like a human would.
Expanding upon the original idea of scanning for torrents as part of the cron job, this tool operates on a more robust foundation,
which is able to detect and go through multiple Pirate Bay mirrors.
It also expands upon dealing with media, by utilizing APIs, crawling through IMDb and more.

### 🗺️ TODO

- [x] Replace IMDB scraping with OMDB API.
- [x] Replace use of TVDB with free alternative (TVMaze).
- [ ] Replace tables in stdout with a more readable format.
- [ ] Add cache & retry system for torrents whose attempts to add to the designated torrent client fail.
- [ ] Support for a proxy or VPN to avoid getting flogged.
- [ ] Interactive CLI for search results, so that the user can navigate with the keyboard and select which to send to qBittorrent for download.
- [ ] Add more sources than the PirateBay.


## ⚓ Installation

You can find compiled binaries for your architecture on the [GitHub releases page](https://github.com/gmantaos/Goirate/releases).

### Updating

The tool comes with a self-updater.

```sh
$ goirate update
Updating to version: 0.9.1

$ goirate --version
Goirate build: 0.9.1
```


## Movies

This tool retrieves info on movies from the [OMDb API](https://www.omdbapi.com).
It is necessary to obtain an API key, and add it to `~/.goirate/config.toml` or in the `GOIRATE_OMDB_API_KEY` environment variable.


#### Search

You can fetch a movie - and torrents for it - in three ways.

```sh
# Using its IMDb url
$ goirate movie "https://www.imdb.com/title/tt1028576/"
Secretariat
IMDbID:         1028576
Year:           2010
Rating:         7.1
Duration:       2h 3min
Poster:         https://m.media-amazon.com/images/M/MV5BMTgwNDkyMDU3NV5BMl5BanBnXkFtZTcwNjMyNjI4Mw@@._V1_UX182_CR0,0,182,268_AL_.jpg

Secretariat[2010]DvDrip-aXXo
URL:            ** omitted **
Seeds/Peers:    1 / 1
Size:           735.8 MB
Trusted:        true
Magnet:
magnet:?** omitted **
```

```sh
# Equivalently, using its IMDb ID
$ goirate movie "1028576"
```

```sh
# Using a partial name
$ goirate movie "avengers"
The Avengers
IMDbID:         0848228
Year:           2012
Rating:         8.1
Duration:       2h 23min
...
```

```sh
# Using both a partial name and a release year to narrow down the search
$ goirate movie -y 2018 "avengers"
Avengers: Infinity War
IMDbID:         4154756
Year:           2018
Rating:         8.6
Duration:       2h 29min
...
```

If you don't remember a movie's title or release year very accurately, you can also do a search.

```sh
$ goirate movie-search "harry potter" -c 4
| IMDb ID |             Title              | Year |
|---------|--------------------------------|------|
| 0241527 |      Harry Potter and the      | 2001 |
|         |        Sorcerer's Stone        |      |
| 0330373 | Harry Potter and the Goblet of | 2005 |
|         |              Fire              |      |
| 0417741 |      Harry Potter and the      | 2009 |
|         |       Half-Blood Prince        |      |
| 1201607 |  Harry Potter and the Deathly  | 2011 |
|         |        Hallows: Part 2         |      |
```

Using the `-d` or `--download` options will also send the torrent to the running qBittorrent client for download.


## Series

Create a watchlist of series, by using the `series add` command.
This stores a list of your series in your account's configuration, specifically in `~/.goirate/series.toml`,
along with the last episode watched for each one. The names can be partial, as they
will be used to search for the full name on the TVDB API. If the last episode is
not specified, the API will be used to fetch the number of the last episode that
aired for this series.

```sh
$ goirate series add "Strike Back" -e "S02E04"
$ goirate series add "The Walking Dead" -e "Season 3 Episode 1"
$ goirate series add "expanse"
```

You can also add a series by its IMDb ID or URL.

```sh
$ goirate series add "https://www.imdb.com/title/tt1856010/" --ls
|   ID   |       Series        | Season | Last Episode | Min. Quality |
|--------|---------------------|--------|--------------|--------------|
| 262980 | House of Cards (US) |   5    |      13      |              |
```

The `series show` command can be used to display the series currently on the
watchlist. The `-j` flag also applies here, printing out the list in JSON format instead.

```sh
$ goirate series show
|   ID   |      Series      | Season | Last Episode | Min. Quality |
|--------|------------------|--------|--------------|--------------|
| 280619 | The Expanse      |   3    |      13      |              |
| 153021 | The Walking Dead |   5    |      13      |    1080p     |
```

The `series remove` command can be used to remove a series given either a
case-insensitive substring in its name, or its TVDB ID.

```sh
$ goirate series remove expanse
$ goirate series remove 153021
```

With a list of series in the watchlist, use the `series scan` command to search for new episodes.

```sh
$ goirate series scan
Torrent found for: The Americans (2013) S06E05
https://pirateproxy.gdn/** omitted **
magnet:?** omitted **

Torrent found for: The Americans (2013) S06E06
https://pirateproxy.gdn/** omitted **
magnet:?** omitted **
```

When the scanner finds a new episode it will also advance the series' last watched episode number forward.
This way, ideally, the tool can keep the watchlist updated while scanning periodically as part of a cron job.
To perform a scan without updating the watchlist use the `--no-update` flag, and, to perform one without
any other side-effects or actions use the `--dry-run` flag.


### E-mail Notifications

Torrents found when scanning can be sent via e-mail.
To enable this, edit the configuration file at `~/.goirate/config.toml` to enable e-mail notifications,
configure the `smtp` settings and specify the list of recipients.

```toml
[smtp]
  host = "smtp.gmail.com"
  port = 587
  username = "...@gmail.com"
  password = "..."

[actions]
  email = "true"
  notify = ["...@gmail.com"]
  ...
```

Watchlist actions can also be specified individually for each series, by editing the file at `~/.goirate/series.toml`.
Action-related options that are specified for a specific series, will override those of the global configuration file.

```toml
[[series]]
  title = "The Last Ship"
  ...
  [series.actions]
    email = "false"
    notify = []
```

### Automatic Downloads

The tool can also be configured to automatically send new torrents to a running [qBittorent](https://www.qbittorrent.org/)
client for download. To enable this edit the configuration file at `~/.goirate/config.toml` to include the necessary HTTP
configuration.

```toml
[qbittorrent]
  url = "https://localhost:8080"
  username = ""
  password = ""

[actions]
  ...
  download = "true"
```

Same as with e-mails, this feature can also be enabled or disabled for a specific series at `~/.goirate/series.toml`.

```toml
[[series]]
  title = "The Last Ship"
  ...
  [series.actions]
    download = "false"
```

With this enabled, any torrents found during scanning will have their magnet links added to the [qBittorent](https://www.qbittorrent.org/)
client. Whether or not they begin downloading immediately once they are added depends on the configuration on the client itself.


## Environment Variables

These variables are used to configure Goirate, when editing the configuration file is not preferable.
For example inside a Docker container.

| Variable | Description | Default |
| -------- | ----------- | ------- |
| GOIRATE_DEBUG | If set to `true`, it enables additional diagnostic messages. | |
| GOIRATE_DIR | The directory used to store configurations and lists. | `~/.goirate` |
| GOIRATE_VERIFIED_UPLOADER | Whether to only accept torrents from trusted or verified uploaders. | `false` |
| GOIRATE_MIN_QUALITY | The minimal acceptable quality for a torrent. |  |
| GOIRATE_MAX_QUALITY | The maximum acceptable quality for a torrent. |  |
| GOIRATE_MIN_SIZE | The minimum acceptable size for a torrent. |  |
| GOIRATE_MAX_SIZE | The maximum acceptable size for a torrent. |  |
| GOIRATE_MIN_SEEDERS | The minimum acceptable amount of seeders for a torrent. | `0` |
| GOIRATE_KODI_MEDIA_PATHS | Use Kodi-friendly paths when downloading media like movies, music albums and episodes. | `false` |
| GOIRATE_DOWNLOADS_DIR | The directory used to store torrent downloads this tool initiates using [qBittorrent](https://qBittorrentbt.com/). | `~/Downloads` |
| GOIRATE_DOWNLOADS_MOVIES | The directory used to store movie torrent downloads this tool initiates using [qBittorrent](https://qBittorrentbt.com/). | `~/Downloads` |
| GOIRATE_DOWNLOADS_SERIES | The directory used to store series torrent downloads this tool initiates using [qBittorrent](https://qBittorrentbt.com/). | `~/Downloads` |
| GOIRATE_DOWNLOADS_MUSIC | The directory used to store music torrent downloads this tool initiates using [qBittorrent](https://qBittorrentbt.com/). | `~/Downloads` |
| GOIRATE_QBT_URL | The url of the [qBittorent](https://www.qbittorrent.org/) http server. | `http://localhost:8080` |
| GOIRATE_QBT_USERNAME | The username used to authenticate to the qBittorent server. | |
| GOIRATE_QBT_PASSWORD | The password used to authenticate to the qBittorent server. | |
| GOIRATE_SMTP_HOST | The address of the SMTP server used for sending out e-mails. | `smtp.gmail.com` |
| GOIRATE_SMTP_PORT | The port of the SMTP server. | 587 |
| GOIRATE_SMTP_USERNAME | The username used to authenticate with the SMTP server. | |
| GOIRATE_SMTP_PASSWORD | The password used to authenticate with the SMTP server. | |
| GOIRATE_ACTIONS_EMAIL | Enable e-mail notifications for torrents found when scanning. Requires a valid SMTP configuration. | `false` |
| GOIRATE_ACTIONS_NOTIFY | A comma-separated list of the e-mails to send torrents to. | |
| GOIRATE_ACTIONS_DOWNLOAD | Enable automatic torrent downloads with [qBittorrent](https://qBittorrentbt.com/). Requires a valid RPC configuration. | `false` |
| GOIRATE_OMDB_API_KEY | The API key to use for accessing the [OMDb API](https://www.omdbapi.com/). |  |

## Known Issues

- The scanner will read and write to the `~/.goirate/series.toml` file multiple times while running without actually locking the file. So editting the file manually while the `series scan` command is running may cause your changes to be overwritten.
- To fill the configuration file at `~/.goirate/options.toml` with the default options, the file is overwritten every time the tool runs. Meaning that even for operations that do not affect the configuration, the file is opened with write privileges. This is temporary until I can figure out a better way to update the config file with new options whenever there's an update.
