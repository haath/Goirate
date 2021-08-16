package main

import (
	"fmt"

	"github.com/imerkle/go-qbittorrent/qbt"
)

// QBittorrentConfig holds the configuration and credentials for communicating with the
// qBittorrent daemon RPC service.
type QBittorrentConfig struct {
	URL      string `toml:"url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// QBittorrentClient is a simple wrapper of the underlying go-qbittorrent client.
type QBittorrentClient struct {
	*qbt.Client
}

// GetClient returns a qBittorrent HTTP client connected with the given configuration.
func (cfg *QBittorrentConfig) GetClient() (client *QBittorrentClient, err error) {

	qb := qbt.NewClient(cfg.URL)

	if cfg.Username != "" {

		loginOpts := qbt.LoginOptions{
			Username: cfg.Username,
			Password: cfg.Password,
		}

		err = qb.Login(loginOpts)
	}

	client = &QBittorrentClient{Client: qb}

	return
}

// AddTorrent sends the given magnet link to the qBittorrent daemon and begins its download,
// configuring the downloaded files to be placed at the specified output directory.
func (client *QBittorrentClient) AddTorrent(magnetLink, downloadDir string) error {

	options := map[string]string{
		"savepath": downloadDir,
	}

	resp, err := client.DownloadFromLink(magnetLink, options)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err != nil {

		return err
	}

	if resp.StatusCode != 200 {

		err = fmt.Errorf("qBittorrent HTTP error: %s", resp.Status)
	}

	return err
}
