package react

import (
	// "fmt"
	"flag"
	"time"
	// "bytes"
	"testing"
	// "net/http"
	// "io/ioutil"
	// "encoding/json"

	"github.com/ksarch-saas/cfgServer/role"
	"github.com/ksarch-saas/cfgServer/meta"
	"github.com/ksarch-saas/cfgServer/cluster"
	// "github.com/ksarch-saas/cfgServer/react/api"
)

func TestUpdatenodes(t *testing.T) {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")

	initCh := make(chan error)
	notifyCh := make(chan int)
	go meta.Run("ssdb-test", "ssdb-test", "tc", "10.94.46.20:2335", notifyCh, initCh)
	for {
		result := 0
		result = <- notifyCh
		if result != 0 {
			break
		}
	}
	cluster.Init()
	go role.Run(initCh)

	fe := NewReact(2335)
	go fe.Run()
	time.Sleep(2*time.Second)

	notifyCh2 := make(chan int)
	initCh2   := make(chan error)
	go cluster.DiscoverCron("/home/users/lichang04/ksarch/gopath/src/github.com/ksarch-saas/cfgServer/seeds.yml", notifyCh2, initCh2)
	
	
	cluster.ProbeCron(notifyCh2)
	


	// nodes := []*meta.Node{
	// 	&meta.Node{
	// 		NodeID:	"10.144.59.41:2700",
	// 	},
	// }
	// mergeSeedsParam := api.MergeSeedsParams{
	// 	Region:	"nj", 
	// 	CfgID:	"10.194.206.34:3100",
	// 	Seeds:	nodes,
	// }
	// reqJson, err := json.Marshal(&mergeSeedsParam)
	// req, err := http.NewRequest("POST", "http://10.94.46.20:2335/region/updatenodes", bytes.NewBuffer(reqJson))
	// req.Header.Set("Content-Type", "application/json")
	
	// resp, err := http.DefaultClient.Do(req)
	// body, err := ioutil.ReadAll(resp.Body)
	// if resp.StatusCode == 200 {
	// 	var rsp api.Response
	// 	d := json.NewDecoder(bytes.NewReader(body))
	// 	d.UseNumber()
	// 	err = d.Decode(&rsp)
	// 	fmt.Println(rsp, err)
	// }
}