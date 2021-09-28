package snykTool

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const defaultTimeout = 10
const defaultWorkerSize = 10

type ConfigProperties map[string]string

func GetToken() string {
	conf, err := ReadConf()
	if err != nil {
		log.Fatal(err)
	}
	token, _ := conf["token"]
	return token
}

func GetGroupId() string {
	conf, err := ReadConf()
	if err != nil {
		log.Fatal(err)
	}
	id, _ := conf["id"]
	return id
}

func GetTimeoutFromConf() int {
	conf, err := ReadConf()
	if err != nil {
		log.Fatal(err)
	}

	t, _ := strconv.Atoi(conf["timeout"])
	if t > 0 {
		return t
	}

	return defaultTimeout
}

func GetWorkerSizeFromConf() int {
	conf, err := ReadConf()
	if err != nil {
		log.Fatal(err)
	}

	t, _ := strconv.Atoi(conf["worker_size"])
	if t > 0 {
		return t
	}

	return defaultWorkerSize
}

func WriteConf(token string, group_id string) {
	home, err := os.UserHomeDir()
	filename := home + "/.snykctl.conf"
	d1 := []byte("[DEFAULT]\ntoken = " + token + "\nid = " + group_id)
	err = ioutil.WriteFile(filename, d1, 0644)
	if err != nil {
		panic(err)
	}
}

func ReadConf() (ConfigProperties, error) {
	config := ConfigProperties{}
	home, err := os.UserHomeDir()
	filename := home + "/.snykctl.conf"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return config, nil
}
