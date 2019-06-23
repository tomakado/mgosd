/**

  This file is part of mgosd.

  mgosd is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  mgosd is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with mgosd. If not, see <https://www.gnu.org/licenses/>.


  Ildar Karymov (c) 2019
  E-mail: hi@ildarkarymov.ru

*/

// Package mgosd provides types and main functionality of mgosd software.
package mgosd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	LOCALHOST    = "127.0.0.1"
	DEFAULT_PORT = 27017
)

type MongoDBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Auth     bool
	Login    string `json:"login"`
	Password string `json:"password"`
	Database string `json:"database"`
	Timeout  time.Duration
}

type Dumper struct {
	dbClient       mgo.Session
	mutex          *sync.Mutex // Necessary for concurrent work with file system.
	collection     string
	ticker         time.Ticker
	DatabaseConfig MongoDBConfig
	OutputPath     string
}

func NewDumper(dbConfig MongoDBConfig, interval time.Duration, output string) (*Dumper, error) {
	d := Dumper{
		ticker:         *time.NewTicker(interval),
		DatabaseConfig: dbConfig,
		OutputPath:     output,
	}

	return &d, nil
}

func (d *Dumper) Start(mutex *sync.Mutex, collection string) error {
	log.Printf("[%s] Spawned", collection)

	d.mutex = mutex
	d.collection = collection

	err := d.initSession()
	if err != nil {
		return err
	}

	defer d.dbClient.Close()
	c := d.dbClient.DB(d.DatabaseConfig.Database).C(collection)

	for range d.ticker.C {
		log.Printf("[%s] Fetching documents...", collection)

		var docs []bson.M
		err := c.Find(bson.M{}).All(&docs)
		if err != nil {
			errStr := fmt.Sprintf(
				"Failed to fetch documents from collection '%s'\n%v",
				collection, err)
			return errors.New(errStr)
		}

		d.processDocuments(&docs)

		log.Printf("[%s] Saving %d elements...", collection, len(docs))

		err = d.saveDocuments(&docs)
		if err != nil {
			return err
		}

		log.Printf("[%s] Saved.", collection)
	}
	return nil
}

func (d *Dumper) initSession() error {
	url := d.makeConnectionUrl()

	session, err := mgo.Dial(url)
	if err != nil {
		return err
	}

	d.dbClient = *session
	return nil
}

// makeConnectionUrl() builds connection url string depending on auth enabled or not
func (d *Dumper) makeConnectionUrl() string {
	var url string

	if d.DatabaseConfig.Auth {
		url = fmt.Sprintf("mondodb://%s:%s@%s:%d/%s",
			d.DatabaseConfig.Login,
			d.DatabaseConfig.Password,
			d.DatabaseConfig.Host,
			d.DatabaseConfig.Port,
			d.DatabaseConfig.Database)
	} else {
		url = fmt.Sprintf("mongodb://%s:%d/%s",
			d.DatabaseConfig.Host,
			d.DatabaseConfig.Port,
			d.DatabaseConfig.Database)
	}

	return url
}

// processDocuments transforms field containing ObjectID
// to be properly importable via mongoimport tool
func (d *Dumper) processDocuments(docs *[]bson.M) {
	for _, doc := range *docs {
		keys := reflect.ValueOf(doc).MapKeys()

		for j := range keys {
			key := keys[j].String()
			formattedKey := strings.ToLower(key)

			if strings.Contains(formattedKey, "id") {
				doc[key] = bson.M{
					"$oid": doc[key],
				}
			}
		}
	}
}

func (d *Dumper) saveDocuments(docs *[]bson.M) error {
	docsJson, err := json.Marshal(docs)
	if err != nil {
		return err
	}

	filePath, err := d.resolvePath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, docsJson, os.FileMode.Perm(os.ModePerm))
	if err != nil {
		return err
	}

	return nil
}

func (d *Dumper) resolvePath() (string, error) {
	now := time.Now()
	timeSignature := fmt.Sprintf("%d-%s-%d-%d-%d-%d",
		now.Year(),
		now.Month().String(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)

	dirPath := d.OutputPath + "/" + timeSignature
	filePath := dirPath + "/" + d.collection + ".json"

	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Create OutputPath directory if not exists
	if _, err := os.Stat(d.OutputPath); os.IsNotExist(err) {
		err = os.MkdirAll(d.OutputPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Create directory for dump if not exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	return filePath, nil
}

func (d *Dumper) Stop() {
	d.ticker.Stop()
	d.dbClient.Close()
}
