package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

const serverURL = "https://github.com/containerum/solution-list/blob/master/containerum-solutions.csv"

func GetList(ctx *gin.Context) {
	log := logrus.WithField("component", "user_manager_client")

	client := resty.New().
		SetLogger(log.WriterLevel(logrus.DebugLevel)).
		SetHostURL(serverURL).
		SetDebug(true)
	client.JSONMarshal = jsoniter.Marshal
	client.JSONUnmarshal = jsoniter.Unmarshal

	client.R().
		SetContext(ctx).
		SetOutput("solutions.csv").
		Get(serverURL)

	ctx.Status(http.StatusOK)
}
