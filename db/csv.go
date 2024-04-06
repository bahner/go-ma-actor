package db

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/bahner/go-ma-actor/config"
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// Save saves a sync.Map to a CSV file.
func Save(syncMap *sync.Map, filename string) error {

	file, err := os.Create(filename)
	if err != nil {
		log.Error("db.Save: ", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	syncMap.Range(func(key, value interface{}) bool {
		writer.Write([]string{key.(string), value.(string)})
		return true
	})
	writer.Flush()

	return writer.Error()
}

// Watch monitors a CSV file for changes and reloads it
// and updates the sync.Map when changed.
func Watch(filename string, syncMap *sync.Map) error {
	log.Debugf("db.Watch: watching %s", filename)

	err := load(filename, syncMap) // Load the file initially
	if err != nil {
		log.Errorf("db.Watch: %v", err)
		return fmt.Errorf("db.Watch: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorf("db.Watch: %v", err)
		return fmt.Errorf("db.Watch: %w", err)
	}
	defer watcher.Close()

	// Function to re-add the file or directory to the watcher list
	reAddWatcher := func(path string) error {
		if err := watcher.Remove(path); err != nil {
			log.Debugf("db.Watch: error removing watcher for %s: %v", path, err)
		}
		if err := watcher.Add(path); err != nil {
			log.Errorf("db.Watch: error re-adding watcher for %s: %v", path, err)
			return fmt.Errorf("db.Watch: %w", err)
		}
		return nil
	}

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Debugf("db.Watch: event: %v", event)
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Infof("db.Watch: %s modified. Reloading.", filename)
					load(filename, syncMap) // Reload the file when it changes
					// Attempt to re-add the file to the watcher, to continue monitoring
					if err := reAddWatcher(filename); err != nil {
						log.Errorf("db.Watch: %v", err)
						return
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				// Handle errors
				log.Errorf("db.Watch: %v", err)
			}
		}
	}()

	if err := watcher.Add(filename); err != nil {
		log.Errorf("db.Watch: %v", err)
		return fmt.Errorf("db.Watch: %w", err)
	}

	// Block until the watcher is done to keep the watcher alive.
	// You may need a mechanism to signal `done` based on your application's needs (not shown here).
	<-done
	return nil
}

// load reads peers from a CSV file and updates the sync.Map.
func load(filename string, nicks *sync.Map) error {

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, config.CSVMode)
	if err != nil {
		log.Errorf("db.load: %v", err)
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Errorf("db.load: %v", err)
		return err
	}

	tempMap := make(map[string]string)
	for _, record := range records {
		tempMap[record[0]] = record[1]
	}

	// Upsert entries from tempPeers into the global peers map
	for key, value := range tempMap {
		nicks.Store(key, value)
	}

	// Optionally, remove entries in peers that are not in tempPeers
	nicks.Range(func(key, value interface{}) bool {
		if _, exists := tempMap[key.(string)]; !exists {
			nicks.Delete(key)
		}
		return true
	})

	return nil
}
