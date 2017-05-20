package cluster

import (
	"time"
	"io/ioutil"

	ymal "gopkg.in/yaml.v2"
	"github.com/golang/glog"
)

const (
	DISCOVER_NEWNODE_TIMEOUT 		= 1
	UPDATE_NEW_SEEDS				= 1
)

var Seeds = make(map[string]string)

func UpdateDiscovery(path string, notifyCh chan int) error {
	buf, err := ioutil.ReadFile(path)
	if err != nil{
		glog.Error("Read seeds path error:", err)
		return err
	}

	newSeeds := make(map[string]string)
	err = ymal.Unmarshal(buf, newSeeds)
	if err != nil {
		glog.Error("Parser ymal error:", err)
		return err
	}

	if len(newSeeds) != len(Seeds) {
		goto update
	}

	for host, _ :=range newSeeds {
		_, ok := Seeds[host]
		if !ok {
			goto update
		}
	}

	return nil

update:
	glog.Info("Seeds change, begion update, old seeds: ", Seeds)
	for old, _ := range Seeds {
		delete(Seeds, old)
	}
	for host, tag := range newSeeds {
		Seeds[host] = tag
	}
	glog.Info("New seeds: ", newSeeds)
	notifyCh <- UPDATE_NEW_SEEDS
	return nil
}

func DiscoverCron(seedsPath string, notifyCh chan int, initCh chan error) {
	err := UpdateDiscovery(seedsPath, notifyCh)
	if err != nil {
		glog.Error("Start up discover error:", err)
		initCh <- err
	}

	tickChan := time.NewTicker(time.Second * DISCOVER_NEWNODE_TIMEOUT).C
	for {
		select {
		case <- tickChan:
			err := UpdateDiscovery(seedsPath, notifyCh)
			if err != nil {
				glog.Error("Period discovery seeds error:", err)
				initCh <- err
			}
		}
	}
}