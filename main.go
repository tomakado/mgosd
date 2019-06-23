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

package main

import (
	"./mgosd"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"sync"
	"time"
)

type Config struct {
	Collections []string            `json:"collections"`
	Interval    string              `json:"interval"`
	Output      string              `json:"output"`
	DB          mgosd.MongoDBConfig `json:"db"`
}

var config Config

func main() {
	fmt.Println("+----------------------------------------+")
	fmt.Println("| mgosd (c) Ildar Karymov, 2019          |")
	fmt.Println("| https://github.com/ildarkarymoff/mgosd |")
	fmt.Println("| License: GNU GPLv3                     |")
	fmt.Println("| Version: 1.1.0                         |")
	fmt.Println("+----------------------------------------+")

	usr, err := user.Current()
	if err != nil {
		log.Fatalln("Failed to get info about current OS user")
	}

	config.DB = mgosd.MongoDBConfig{}
	config.DB.Database = *flag.String("db", "", "Database name")
	config.DB.Host = *flag.String("host", mgosd.LOCALHOST, "Database server address")
	config.DB.Port = *flag.Int("port", mgosd.DEFAULT_PORT, "Database server port")
	login := *flag.String("login", "<empty>", "Database username")
	password := *flag.String("password", "<empty>", "Database user password")
	config.Output = *flag.String("o", usr.HomeDir+"/mgosd/", "Path to output directory")
	config.Collections = flag.Args()
	config.Interval = *flag.String("i", "12h", "Interval of dumping")

	flag.Parse()

	if len(os.Args) > 1 {
		configPath := os.Args[1]

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalln(err)
		} else {
			rawConfig, err := ioutil.ReadFile(configPath)
			err = json.Unmarshal(rawConfig, &config)
			if err != nil {
				log.Printf("Failed to read config file: %v", err)
			}
		}
	}

	if login != "<empty>" || password != "<empty>" {
		config.DB.Login = login
		config.DB.Password = password
		config.DB.Auth = true
	}

	fmt.Println("\n mgosd will work with following parameters:")
	fmt.Printf(" * Database host: %s\n", config.DB.Host)
	fmt.Printf(" * Database port: %d\n", config.DB.Port)
	fmt.Printf(" * Database authentication: %v\n", config.DB.Auth)
	fmt.Printf(" * Collections (%d): %v\n", len(config.Collections), config.Collections)
	fmt.Printf(" * Dump interval: %s\n", config.Interval)
	fmt.Printf(" * Output path: %s\n\n", config.Output)

	var (
		wg    sync.WaitGroup
		mutex sync.Mutex
	)

	for i := 0; i < len(config.Collections); i++ {
		wg.Add(1)
		go startWorker(&mutex, &config, config.Collections[0], config.Collections[i])
	}
	wg.Wait()
}

func startWorker(mutex *sync.Mutex, config *Config, masterCollection string, collection string) {
	log.Printf("[%s] Spawning worker...", collection)

	var intervalScale time.Duration
	interval, err := strconv.Atoi(config.Interval[:len(config.Interval)-1])
	if err != nil {
		log.Fatalln(err)
	}

	switch config.Interval[len(config.Interval)-2:] {
	case "h":
		intervalScale = time.Hour
	case "m":
		intervalScale = time.Minute
	case "s":
		intervalScale = time.Second
	default:
		intervalScale = time.Second
	}

	dumper, _ := mgosd.NewDumper(
		config.DB,
		time.Duration(interval)*intervalScale,
		config.Output)

	err = dumper.Start(mutex, collection)
	if err != nil {
		log.Fatalln(err)
	}
}
