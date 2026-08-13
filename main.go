package main

import (
	"flag"
	"fluxKV/configuration"
	"fluxKV/server"
	"fluxKV/utils"
	"os"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to the config yaml")
	flag.Parse()

	cfg, err := configuration.Load(*configPath)
	if err != nil {
		utils.Logger.Fatal("%v", err)
	}

	id := os.Getenv("SERVER_ID")
	if id == "" {
		id = cfg.Servers[0].ID
	}

	self, err := cfg.FindByID(id)
	if err != nil {
		utils.Logger.Fatal("%v", err)
	}

	utils.Logger.Info("starting server %s (role=%s)", self.ID, self.Role)

	s := server.NewServer(cfg, self)

	if err = s.Listen(); err != nil {
		utils.Logger.Fatal("%v", err)
	}
}
