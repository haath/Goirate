
# Gorrent


## Torrents

The primary source of this tool's torrents is The Pirate Bay.

### Mirrors

To get a list of The Pirate Bay mirrors, use the `gorrent mirrors` command.

```shell
$ gorrent mirrors
[x] uk https://pirateproxy.sh
[x] nl https://thepbproxy.com
[x] us https://thepiratebay.red
[x] us https://thepiratebay-org.prox.space
[x] us https://cruzing.xyz
[x] us https://tpbproxy.nl
[x] us https://thepiratebay.rocks
[ ] us https://proxydl.cf
[x] us https://piratebayblocked.com
[x] us https://tpb.crushus.com/thepiratebay.org
[x] us https://ikwilthepiratebay.org
[x] gb https://bay.maik.rocks
[ ] fr https://www.piratenbaai.ch
[ ] us https://tpbproxy.gdn
[ ] us https://tpb.network
[ ] fr https://thepiratebay.freeproxy.fun
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