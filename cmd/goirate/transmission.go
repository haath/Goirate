package main

import "github.com/hekmon/transmissionrpc"

// RPCConfig holds the configuration and credentials for communicating with the
// transmission daemon RPC service.
type RPCConfig struct {
	Host     string `toml:"host"`
	Port     uint16 `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	UseSSL   bool   `toml:"ssl"`
}

// RPCClient is a simple wrapper of the underlying transmissionrpc.Client.
type RPCClient struct {
	*transmissionrpc.Client
}

// GetClient returns a transmission RPC client connected with the given configuration.
func (cfg *RPCConfig) GetClient() (*RPCClient, error) {

	client, err := transmissionrpc.New(cfg.Host, cfg.Username, cfg.Password, &transmissionrpc.AdvancedConfig{
		HTTPS: cfg.UseSSL,
		Port:  cfg.Port,
	})

	return &RPCClient{Client: client}, err
}

// AddTorrent sends the given magnet link to the transmission daemon and begins its download,
// configuring the downloaded files to be placed at the specified output directory.
func (client *RPCClient) AddTorrent(magnetLink, downloadDir string) error {

	payload := transmissionrpc.TorrentAddPayload{
		Filename:    &magnetLink,
		DownloadDir: &downloadDir,
	}

	_, err := client.TorrentAdd(&payload)

	return err
}
