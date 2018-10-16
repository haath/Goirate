![Logo](assets/logo.png)

[![](https://git.gmantaos.com/haath/Goirate/badges/master/pipeline.svg)](https://git.gmantaos.com/haath/Goirate/pipelines)
[![](https://git.gmantaos.com/haath/Goirate/badges/master/coverage.svg)](https://git.gmantaos.com/haath/Goirate/-/jobs/artifacts/master/browse?job=test)
[![](https://goreportcard.com/badge/git.gmantaos.com/haath/Goirate)](https://goreportcard.com/report/git.gmantaos.com/haath/Goirate)
[![](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
[![](https://img.shields.io/github/release/gmantaos/Goirate.svg)](https://github.com/gmantaos/Goirate/releases)

Watching a lot of movies and series, it quickly became difficult to keep track of what was coming out and when.
Not to call managing a few torrents tiresome, but when there's 5-10 new episodes for things you watch coming out
every week, you begin to wonder if all of these extra clicks are really necessary. This also refers to wasted clicks,
for when you know that an episode aired but you don't quite know yet if there's a torrent out for the 1080p version you prefer.
With all this in mind, I first attempted automate this with a simple python script, which would run as a cron job, crawl the Pirate Bay 
for torrents of new episodes, send me the ones it finds via e-mail and update the list of series so that it would begin watching out
for the next episode. And the funny thing that script **worked like clockwork, monitoring at least 40 different series over a period of two years**.

Then the point came, when I wanted more features and automation, and that old python script was never written to be particularly scalable.
So I began development of Goirate.
This tool aims to become an all-in-one suite for automating your every piratey need.
It works as a CLI program, which is designed to go through the internet searching for torrents much like a human would.
Expanding upon the original idea of scanning for torrents as part of the cron job, this tool operates on a more robust foundation,
which is able to detect and go through multiple Pirate Bay mirrors.
It also expands upon dealing with media, by utilizing APIs, crawling through IMDb and more.

### 🗺️ P️rogress 

- [x] Pirate Bay scraping for booty
- [x] Command-line plunderin'
- [x] IMDB Scraper
- [x] Robust execution - searching multiple mirrors
- [x] Global configuration management
- [x] Defining series seasons and episodes
- [x] TVDB integration
- [x] Scanning for new series episodes
- [ ] Defining [sea shanties](https://en.wikipedia.org/wiki/Sea_shanty) and their albums
- [x] Torrent client integration ([Transmission](https://transmissionbt.com/))
- [ ] Kodi-friendly download storage
- [ ] Crontab scanner
    - [x] Defining handlers for torrents found
        - [x] E-mail notifications
        - [x] Automatic downloads
    - [ ] Watchlist for single torrents
    - [x] New series episodes
    - [ ] RSS Feeds (?)
- [ ] Support for a proxy or VPN to avoid getting flogged
- [ ] Docker image with Transmission and OpenVPN.

## Installation

You can find compiled binaries for your architecture on the [GitHub releases page](https://github.com/gmantaos/Goirate/releases).

### Updating

The tool comes with a self-updater.

```sh
$ goirate update
Updating to version: 0.9.1

$ goirate --version
Goirate build: 0.9.1
```

### Build from source

By default [dep](https://github.com/golang/dep) is used for dependency management.

The `Makefile` has a shortcut to running `dep` and `go install`.

```sh
$ make install
```

Using `go get` to fetch dependencies is theoretically possible but it is not recommended.
Also, attempting to install the tool with `go get -u` will not work as it uses [packr](https://github.com/gobuffalo/packr)
for building. To build yourself use the `Makefile` or have a look at it.

## ⚓ Usage

### Torrents

The primary source of this tool's torrents is The Pirate Bay.

Commands that search for torrents support the following options.

| | |
|-|-|
| `-j`, `--json` | Output JSON |
| `--mirror "https://pirateproxy.sh/"` | Use a specific pirate bay mirror |
| `--source "https://proxybay.bz/"` | Override default mirror list |
| `--trusted` | Only return torrents whose uploader is either Trusted or VIP |
| `--only-magnet` | Only output magnet links, one on each line |
| `--only-url` | Only output torrent urls, one on each line |
| `-c 7`, `--count 7` | Limit the number of results |

#### Pirate Bay mirrors

To get a list of The Pirate Bay mirrors, use the `mirrors` command.

```shell
$ goirate mirrors
|   | Country |                   URL                    |
|---|---------|------------------------------------------|
| x |   UK    | https://pirateproxy.sh                   |
| x |   NL    | https://thepbproxy.com                   |
| x |   US    | https://thetorrents.red                  |
| x |   US    | https://thetorrents-org.prox.space       |
| x |   US    | https://cruzing.xyz                      |
| x |   US    | https://tpbproxy.nl                      |
| x |   US    | https://thetorrents.rocks                |
| x |   US    | https://proxydl.cf                       |
| x |   US    | https://torrentsblocked.com              |
| x |   US    | https://tpb.crushus.com/thetorrents.org  |
| x |   US    | https://ikwilthetorrents.org             |
| x |   GB    | https://bay.maik.rocks                   |
|   |   FR    | https://www.piratenbaai.ch               |
|   |   US    | https://tpbproxy.gdn                     |
|   |   US    | https://tpb.network                      |
| x |   FR    | https://thetorrents.freeproxy.fun        |
```

By default, the tool will attempt to fetch them from [proxybay.github.io](https://proxybay.github.io). To override that set the `-s` option.

```shell
$ goirate mirrors -s https://proxybay.bz/
```

You can also integrate the tool with any application by getting the output in JSON format using the `--json` flag.

```shell
$ goirate mirrors --json
[
   {
      "url": "https://pirateproxy.sh",
      "country": "uk",
      "status": true
   },
   {
      "url": "https://thepbproxy.com",
      "country": "nl",
      "status": true
   },
   ...
}
```

The mirror-picking process for the entire tool can be configured in your `~/.goirate/config.toml`.
Both the blacklist and the whitelist can contain partial mirror URLs or country codes.

```toml
# Only allow mirrors in the US, but do not allow thepiratebay.vin or *.biz domains
[tpb_mirrors]
  whitelist = ["US"]
  blacklist = ["thepiratebay.vin", ".biz"]
```

#### Searching for torrents

The `search` command can be used to find torrents given a specific query and filters.

```shell
$ goirate search "debian"
                                                  Title                                                      Size    Seeds/Peers
---------------------------------------------------------------------------------------------------------------------------------
 Debian GNU/Linux Bible [ENG] [.pdf]                                                                        7.5 MB   10 / 12
 https://pirateproxy.sh/torrent/** omitted **
---------------------------------------------------------------------------------------------------------------------------------
 Debian 7- System Administration Best Practices, 2013 [PDF]~Storm                                           2.0 MB   9 / 9
 https://pirateproxy.sh/torrent/** omitted **
---------------------------------------------------------------------------------------------------------------------------------
 Debian 9 Stretch minimal install (VirtualBox VDI image)                                                   187.7 MB  6 / 6
 https://pirateproxy.sh/torrent/** omitted **
---------------------------------------------------------------------------------------------------------------------------------
 Debian GNU Linux Bible.zip                                                                                 6.1 MB   2 / 2
 https://pirateproxy.sh/torrent/** omitted **
---------------------------------------------------------------------------------------------------------------------------------
```

Much like any other command, you can use the `--help` flag the retrieve the 
list of available options.

```shell
$ goirate search --help
```

## Movies

This tool scrapes [IMDb.com](https://www.imdb.com/) for info on movies.

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

Using the `-d` or `--download` options will also send the torrent to the running transmission daemon for download.

## Series

For this tool to manage series, you need to obtain an API key from [TheTVDB.com](https://www.thetvdb.com/)
and include it in Goirate's configuration at `~/.goirate/config.toml`.
Once logged in, the following can be found [here](https://www.thetvdb.com/member/api).

```toml
[tvdb]
  api_key = "< API Key >"
  user_key = "< Unique ID >"
  username = "< Username >"
```

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

The tool can also be configured to automatically send new torrents to a running [Transmission](https://transmissionbt.com/)
daemon for download. To enable this edit the configuration file at `~/.goirate/config.toml` to include the necessary RPC
configuration.

```toml
[transmission_rpc]
  host = "localhost"
  port = 9091
  username = ""
  password = ""
  ssl = false

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

With this enabled, any torrents found during scanning will have their magnet links added to the [Transmission](https://transmissionbt.com/)
daemon. Whether or not they begin downloading immediately once they are added depends on the configuration of daemon.

## Environment Variables

These variables are used to configured Goirate when editing the configuration file 

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
| GOIRATE_DOWNLOADS_DIR | The directory used to store torrent downloads this tool initiates using [Transmission](https://transmissionbt.com/). | `~/Downloads` |
| GOIRATE_DOWNLOADS_MOVIES | The directory used to store movie torrent downloads this tool initiates using [Transmission](https://transmissionbt.com/). | `~/Downloads` |
| GOIRATE_DOWNLOADS_SERIES | The directory used to store series torrent downloads this tool initiates using [Transmission](https://transmissionbt.com/). | `~/Downloads` |
| GOIRATE_DOWNLOADS_MUSIC | The directory used to store music torrent downloads this tool initiates using [Transmission](https://transmissionbt.com/). | `~/Downloads` |
| GOIRATE_RPC_HOST | The host of the [Transmission](https://transmissionbt.com/) RPC daemon. | `localhost` |
| GOIRATE_RPC_PORT | The port of the [Transmission](https://transmissionbt.com/) RPC daemon. | 9091 |
| GOIRATE_RPC_USERNAME | The username used to authenticate to the RPC daemon. | |
| GOIRATE_PRC_PASSWORD | The password used to authenticate to the RPC daemon. | |
| GOIRATE_RPC_SSL | Set to `true` to indicate that the [Transmission](https://transmissionbt.com/) RPC should accessed over SSL. | `false` |
| GOIRATE_SMTP_HOST | The address of the SMTP server used for sending out e-mails. | `smtp.gmail.com` |
| GOIRATE_SMTP_PORT | The port of the SMTP server. | 587 |
| GOIRATE_SMTP_USERNAME | The username used to authenticate with the SMTP server. | |
| GOIRATE_SMTP_PASSWORD | The password used to authenticate with the SMTP server. | |
| GOIRATE_ACTIONS_EMAIL | Enable e-mail notifications for torrents found when scanning. Requires a valid SMTP configuration. | `false` |
| GOIRATE_ACTIONS_NOTIFY | A comma-separated list of the e-mails to send torrents to. | |
| GOIRATE_ACTIONS_DOWNLOAD | Enable automatic torrent downloads with [Transmission](https://transmissionbt.com/). Requires a valid RPC configuration. | `false` |

## Known Issues

- The scanner will read and write to the `~/.goirate/series.toml` file multiple times while running without actually locking the file. So editting the file manually while the `series scan` command is running may cause your changes to be overwritten.