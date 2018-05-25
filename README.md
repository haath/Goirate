
# Gorrent


## Torrents

The primary source of this tool's torrents is The Pirate Bay.

### Mirrors

To get a list of The Pirate Bay mirrors, use the `gorrent mirrors` command.

```shell
$ gorrent mirrors
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
$ gorrent mirrors -s https://proxybay.bz/
```

You can also integrate the tool with any application by getting the output in JSON format using the `--json` flag.

```shell
$ gorrent mirrors --json
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