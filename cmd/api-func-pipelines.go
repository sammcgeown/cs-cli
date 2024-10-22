package cmd

import (
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/mitchellh/mapstructure"
)

func getPipelines(id string, name string, project string, export bool, exportPath string) ([]*CodeStreamPipeline, error) {
	var arrResults []*CodeStreamPipeline
	var qParams = make(map[string]string)
	client := resty.New()

	var filters []string
	if id != "" {
		filters = append(filters, "(id eq '"+id+"')")
	}
	if name != "" {
		filters = append(filters, "(name eq '"+name+"')")
	}
	if project != "" {
		filters = append(filters, "(project eq '"+project+"')")
	}
	if len(filters) > 0 {
		qParams["$filter"] = "(" + strings.Join(filters, " and ") + ")"
	}
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&documentsList{}).
		SetAuthToken(targetConfig.accesstoken).
		Get("https://" + targetConfig.server + "/pipeline/api/pipelines")

	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	for _, value := range queryResponse.Result().(*documentsList).Documents {
		c := CodeStreamPipeline{}
		mapstructure.Decode(value, &c)
		if export {
			exportYaml(c.Name, c.Project, exportPath, "piplines")
			arrResults = append(arrResults, &c)
		} else {
			arrResults = append(arrResults, &c)
		}
	}
	return arrResults, err
}

// patchPipeline - Patch Code Stream Pipeline by ID
func patchPipeline(id string, payload string) (*CodeStreamPipeline, error) {
	client := resty.New()
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		SetResult(&CodeStreamPipeline{}).
		SetAuthToken(targetConfig.accesstoken).
		Patch("https://" + targetConfig.server + "/pipeline/api/pipelines/" + id)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*CodeStreamPipeline), err
}

func deletePipeline(id string) (*CodeStreamPipeline, error) {
	client := resty.New()
	queryResponse, err := client.R().
		SetQueryParams(qParams).
		SetHeader("Accept", "application/json").
		SetResult(&CodeStreamPipeline{}).
		SetAuthToken(targetConfig.accesstoken).
		Delete("https://" + targetConfig.server + "/pipeline/api/pipelines/" + id)
	if queryResponse.IsError() {
		return nil, queryResponse.Error().(error)
	}
	return queryResponse.Result().(*CodeStreamPipeline), err
}
