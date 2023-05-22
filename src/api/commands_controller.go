package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"kyleroberts.io/src/api/payloads"
	"kyleroberts.io/src/azure"
)

func (env *AppEnvironment) AzureAuthenticate() {
	authRequirements := azure.AzureTokenAuthRequirements{
		ClientID:     env.Config.Azure.App.ClientID,
		ClientSecret: env.Config.Azure.App.ClientSecret,
		TenantID:     env.Config.Azure.TenantID,
		Scope:        env.Config.Azure.App.Scope,
	}
	authRes, loginErr := azure.GetAzureToken(authRequirements)
	if loginErr != nil {
		fmt.Println("failed to authenticate with Azure")
	}
	env.AzureAccessToken = authRes.AccessToken
}

type SubnetDetails struct {
	VNetName      string `json:"vnet_name"`
	SubnetName    string `json:"subnet_name"`
	Subscription  string `json:"subscription"`
	ResourceGroup string `json:"resource_group"`
}

// type CreateContainerGroupInbound struct {
// 	Subscription       string        `json:"subscription"`
// 	ResourceGroup      string        `json:"resource_group"`
// 	ContainerGroupName string        `json:"container_group_name"`
// 	TemplateName       string        `json:"template_name"`
// 	SubnetDetails      SubnetDetails `json:"subnet"`
// }

func (env *AppEnvironment) CreateContainerGroup(context *gin.Context) {
	env.AzureAuthenticate()
	payload := new(payloads.CreateContainerGroup)
	bindErr := context.BindJSON(&payload)
	if bindErr != nil {
		context.AbortWithError(http.StatusBadRequest, bindErr)
		return
	}
	cgManager := azure.ContainerGroupManager{
		AccessToken:   env.AzureAccessToken,
		APIVersion:    "2022-09-01",
		Subscription:  payload.Subscription,
		ResourceGroup: payload.ResourceGroup,
	}
	azErr := cgManager.Create(payload)
	if azErr != nil {
		context.AbortWithError(http.StatusBadRequest, azErr)
		return
	} else {
		context.JSON(
			http.StatusOK,
			gin.H{"message": "Created Azure Container Instance"},
		)
		return
	}
}

// func (env *AppEnvironment) CreateContainerGroup(context *gin.Context) {
// 	env.AzureAuthenticate()
// 	payload := new(CreateContainerGroupInbound)
// 	bindErr := context.BindJSON(&payload)
// 	if bindErr != nil {
// 		context.AbortWithError(http.StatusBadRequest, bindErr)
// 		return
// 	}
// 	templateConfig, templateErr := templates.Parse(payload.TemplateName)
// 	if templateErr != nil {
// 		context.AbortWithError(http.StatusBadRequest, templateErr)
// 		return
// 	}
// 	cg := azure.ContainerGroup{
// 		Subscription:  payload.Subscription,
// 		ResourceGroup: payload.ResourceGroup,
// 		Name:          payload.ContainerGroupName,
// 	}
// 	err := cg.Create("2022-09-01", env.AzureAccessToken, *templateConfig)
// 	if err != nil {
// 		context.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	} else {
// 		context.JSON(
// 			http.StatusOK,
// 			gin.H{"message": "Created Azure Container Instance"},
// 		)
// 		return
// 	}
// }
