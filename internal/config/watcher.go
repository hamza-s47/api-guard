package config

import (
	"log"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
)

var currnt atomic.Value

func Get() *Config {
	return currnt.Load().(*Config)
}

func Watch(path string) error {
	cfg, err := Load(path)
	if err != nil {
		return err
	}
	currnt.Store(cfg)

	wather, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-wather.Events:
				if event.Op&fsnotify.Write != 0 {
					newCfg, err := Load(path)
					if err != nil {
						log.Println("Config reload failed:", err)
						continue
					}
					currnt.Store(newCfg)
					log.Panicln("Config reloaded successfully")
				}
			case err := <-wather.Errors:
				log.Println("Watch error:", err)
			}
		}
	}()

	return wather.Add(path)
}
