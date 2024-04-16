package pvf_source

import (
	"fmt"
	"github.com/dof-wiki/godof/utils"
	"net/url"
	"strings"
)

type PvfUtilitySource struct {
	host string
}

func NewPvfUtilitySource(host string) *PvfUtilitySource {
	return &PvfUtilitySource{
		host: host,
	}
}

func (p *PvfUtilitySource) getQueryUrl(path string, params map[string]string) string {
	paramList := make([]string, 0, len(params))
	for k, v := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
	}

	query := strings.Join(paramList, "&")
	return fmt.Sprintf("%s%s?%s", p.host, path, query)
}

func (p *PvfUtilitySource) GetFileContent(path string) (string, error) {
	u := p.getQueryUrl(PathGetFileContent, map[string]string{"filePath": path})
	type Rsp struct {
		Data    string
		IsError bool
		Msg     string
	}
	rsp := new(Rsp)
	err := utils.HTTPGet(u, rsp)
	return rsp.Data, err
}

func (p *PvfUtilitySource) SaveFileContent(path, content string) error {
	u := p.getQueryUrl(PathImportFile, map[string]string{"filePath": path})

	type Rsp struct {
		IsError bool
		Msg     string
	}
	return utils.HTTPPost(u, &content, new(Rsp))
}

func (p *PvfUtilitySource) GetItemInfos(paths []string) (map[int]string, error) {
	u := p.host + PathGetItemInfos

	type Item struct {
		ItemName string
		ItemCode int
	}

	type Rsp struct {
		Data    map[string]*Item
		IsError bool
		Msg     string
	}

	rsp := new(Rsp)
	err := utils.HTTPPost(u, &paths, rsp)
	if err != nil {
		return nil, err
	}

	ret := make(map[int]string)
	for _, item := range rsp.Data {
		ret[item.ItemCode] = item.ItemName
	}
	return ret, nil
}
