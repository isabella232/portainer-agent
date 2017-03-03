package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/rancher/go-rancher/v2"
)

const (
	metadataURL = "http://169.254.169.250/2015-12-19"
)

func Watch(file, accessKey, secretKey, url string) error {
	logrus.Infof("Watching for changes %s %s", accessKey, url)
	m, err := metadata.NewClientAndWait(metadataURL)
	if err != nil {
		return err
	}

	client, err := client.NewRancherClient(&client.ClientOpts{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Url:       url,
	})
	if err != nil {
		return err
	}

	proxy := NewProxy(client)

	for {
		time.Sleep(2 * time.Second)

		hosts, err := m.GetHosts()
		if err != nil {
			logrus.Errorf("Error gettings hosts: %v", err)
			continue
		}

		hostMap := proxy.AddHosts(hosts)
		content, err := ConstructFile(hostMap)
		if err != nil {
			logrus.Errorf("Error constructing hosts %v: %v", hostMap, err)
			continue
		}

		if err := WriteFile(file, content); err != nil {
			logrus.Errorf("Failed to write [%s] to file %s: %v", content, file, err)
		}
	}
}

func ConstructFile(data map[string]string) ([]byte, error) {
	result := []map[string]string{}
	for name, address := range data {
		result = append(result, map[string]string{
			"Name": name,
			"URL":  fmt.Sprintf("tcp://%s", address),
		})
	}
	return json.Marshal(result)
}

func WriteFile(file string, content []byte) error {
	err := ioutil.WriteFile(file+".tmp", []byte(content), 0644)
	if err != nil {
		return err
	}

	return os.Rename(file+".tmp", file)
}
