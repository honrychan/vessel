package handler

import (
	"net/http"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module"
	"github.com/containerops/vessel/utils"

	"gopkg.in/macaron.v1"
)

type PipelinePOSTJSON struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	MetaData   struct {
		Name        string `json:"name"`
		Namespace   string `json:"namespace"`
		SelfLink    string `json:"selfLink"`
		Labels      string `json:"labels"`
		Annotations string `json:"annotations"`
	} `json:"metadata"`
	Spec []struct {
		Name                string `json:"name"`
		Dependence          string `json:"dependence"`
		Kind                string `json:"kind"`
		StatusCheckUrl      string `json:"statusCheckUrl"`
		StatusCheckInterval int64  `json:"statusCheckInterval"`
		StatusCheckCount    int64  `json:"statusCheckCount"`
	} `json:"spec"`
}

func V1POSTPipelineHandler(ctx *macaron.Context, reqData PipelinePOSTJSON) (int, []byte) {
	// 收json
	reqStr, _ := ctx.Req.Body().String()
	// 生成新的pipeline
	projectInfo := module.GetProjectInfoByName(ctx.Params(":project"))
	pipeline := new(models.Pipeline)
	pipeline.WorkspaceId = projectInfo.WorkspaceId
	pipeline.ProjectId = projectInfo.Id
	pipeline.Name = reqData.MetaData.Name
	pipeline.SelfLink = ""
	pipeline.Labels = reqData.MetaData.Labels
	pipeline.Annotations = reqData.MetaData.Annotations
	pipeline.Created = time.Now().Unix()
	pipeline.Updated = time.Now().Unix()
	pipeline.Detail = reqStr

	module.CreatePipeline(pipeline)
	// 生成stage和point
	// 根据json判断合法性,并返回stage 和 point的对应关系map
	isLegal, reason, dependenceMap := utils.GenerateDependenceMap(reqStr)
	if !isLegal {
		return http.StatusOK, []byte(reason)
	}

	for name, define := range dependenceMap {
		point := new(models.Point)
		point.PipelineId = pipeline.Id
		point.Created = time.Now().Unix()
		point.Updated = time.Now().Unix()
		point.Name = name
		point.From = define[0]
		point.To = define[1]

		module.CreatePoint(point)

		stage := new(models.Stage)
		stage.PipelineId = pipeline.Id
		stage.Created = time.Now().Unix()
		stage.Updated = time.Now().Unix()
		stage.Name = name
		stage.Detail = define[2]

		module.CreateStage(stage)
	}

	return http.StatusOK, []byte("ok")
}

func V1PUTPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1GETPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1DELETEPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1RunPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1StatusGETHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
