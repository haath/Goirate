![Logo](assets/logo.png)

This tool aims to become an all-in-one suite for automating your every pirate-y need.

### Progress

- [x] Pirate Bay scraping for booty
- [x] Command-line plunderin'
- [ ] Defining series seasons and episodes
- [ ] Defining [sea shanty](https://en.wikipedia.org/wiki/Sea_shanty) albums
- [ ] Support for your proxy or VPNt o avoid the flogging
- [ ] Torrent client integration
- [ ] Crontab scanner
    - [ ] Watchlist for single torrents
    - [ ] New series episodes
    - [ ] E-mail notifications
    - [ ] Automatic downloads

## Command line tool

The primary source of this tool's torrents is The Pirate Bay.

### Mirrors

To get a list of The Pirate Bay mirrors, use the `goirate mirrors` command.

```shell
$ goirate mirrors
|   | Country |                   URL                    |
|---|---------|------------------------------------------|
| x |   UK    | https://pirateproxy.sh                   |
| x |   NL    | https://thepbproxy.com                   |
| x |   US    | https://thepiratebay.red                 |
| x |   US    | https://thepiratebay-org.prox.space      |
| x |   US    | https://cruzing.xyz                      |
| x |   US    | https://tpbproxy.nl                      |
| x |   US    | https://thepiratebay.rocks               |
| x |   US    | https://proxydl.cf                       |
| x |   US    | https://piratebayblocked.com             |
| x |   US    | https://tpb.crushus.com/thepiratebay.org |
| x |   US    | https://ikwilthepiratebay.org            |
| x |   GB    | https://bay.maik.rocks                   |
|   |   FR    | https://www.piratenbaai.ch               |
|   |   US    | https://tpbproxy.gdn                     |
|   |   US    | https://tpb.network                      |
| x |   FR    | https://thepiratebay.freeproxy.fun       |
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


### Torrents

```shell
$ goirate search "debian"
                                                  Title                                                      Size    Seeds/Peers
---------------------------------------------------------------------------------------------------------------------------------
 Debian GNU/Linux Bible [ENG] [.pdf]                                                                        7.5 MB   10 / 12
 https://pirateproxy.sh/torrent/5468273/Debian_GNU_Linux_Bible_[ENG]_[.pdf]
---------------------------------------------------------------------------------------------------------------------------------
 Debian 7- System Administration Best Practices, 2013 [PDF]~Storm                                           2.0 MB   9 / 9
 https://pirateproxy.sh/torrent/9499287/Debian_7-_System_Administration_Best_Practices__2013_[PDF]_Storm
---------------------------------------------------------------------------------------------------------------------------------
 Debian 9 Stretch minimal install (VirtualBox VDI image)                                                   187.7 MB  6 / 6
 https://pirateproxy.sh/torrent/20414237/Debian_9_Stretch_minimal_install_(VirtualBox_VDI_image)
---------------------------------------------------------------------------------------------------------------------------------
 Debian GNU Linux Bible.zip                                                                                 6.1 MB   2 / 2
 https://pirateproxy.sh/torrent/4431647/Debian_GNU_Linux_Bible.zip
---------------------------------------------------------------------------------------------------------------------------------
```

Additional available options

| | |
|-|-|
| `-j`, `--json` | Output JSON |
| `--mirror "https://pirateproxy.sh/"` | Use a specific pirate bay mirror |
| `--source "https://proxybay.bz/"` | Override default mirror list |
| `--trusted` | Only return torrents whose uploader is either Trusted or VIP |
| `--magnet` | Only output magnet links, one on each line |
| `--urls` | Only output torrent urls, one on each line |
| `-c 7`, `--count 7` | Limit the number of results |