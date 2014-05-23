package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type GoogleResultEntry struct {
	unescapedUrl      string
	visisbleUrl       string
	titleNoFormatting string
	content           string
}

type GoogleResponseData struct {
	Data []GoogleResultsData `json:"results"`
}

type GoogleResultsData struct {
	URL        string `json:"unescapedURL"`
	VisibleURL string `json:"visibleUrl"`
	Title      string `json:"titleNoFormatting"`
	Content    string `json:"content"`
}

type GoogleResult struct {
	Results GoogleResponseData `json:"responseData"`
	Status  int                `json:"responseStatus"`
	Error   string
}

func Google(searchTerm string, results int) GoogleResult {
	var resp GoogleResult
	response, err := http.Get(fmt.Sprintf("http://ajax.googleapis.com/ajax/services/search/web?v=1.0&rsz=%d&q=%s", results, strings.Replace(searchTerm, " ", "%20", -1)))
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		resp.Error = err.Error()
		return resp
	}
	err = json.Unmarshal(contents, &resp)
	if err != nil {
		resp.Error = err.Error()
	}
	return resp
}
