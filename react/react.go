package react

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksarch-saas/cfgServer/react/api"
	"github.com/ksarch-saas/cfgServer/react/command"
)

type React struct {
	C				*command.Controller
	Router         	*gin.Engine
	HttpBindAddr   	string
	ProcessTimeout 	time.Duration
}

func NewReact(httpPort int) *React {
	react := &React{
		C:				&command.Controller{},
		Router:         gin.Default(),
		HttpBindAddr:   fmt.Sprintf(":%d", httpPort),
		ProcessTimeout: 60,
	}
	
	gin.SetMode(gin.ReleaseMode)

	react.Router.POST("/region/updatenodes", react.HandleUpdateNodes)

	return react
}

func (react *React) Run() {
	react.Router.Run(react.HttpBindAddr)
}

func (react *React) HandleUpdateNodes(c *gin.Context) {
	var params api.UpdateNodesParams
	c.Bind(&params)

	cmd := command.UpdateNodesCommand{
		Region:	params.Region, 
		CfgID:	params.CfgID,
		Seeds:	params.Seeds,
	}

	result, err := react.C.ProcessCommand(&cmd, react.ProcessTimeout*time.Second)
	if err != nil {
		c.JSON(200, api.MakeFailureResponse(err.Error()))
		return
	}

	c.JSON(200, api.MakeSuccessResponse(result))
}