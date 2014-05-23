package web

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
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
	var uri string = fmt.Sprintf("http://ajax.googleapis.com/ajax/services/search/web?v=1.0&rsz=%d&q=%s", results, url.QueryEscape(searchTerm))
	Get(uri, func(error string, body []byte) {
		err := json.Unmarshal(body, &resp)
		if err != nil {
			resp.Error = err.Error()
		} else {
			for i, _ := range resp.Results.Data {
				if resp.Results.Data[i].Title != "" {
					resp.Results.Data[i].Title = html.UnescapeString(resp.Results.Data[i].Title)
				}
				if resp.Results.Data[i].Content != "" {
					resp.Results.Data[i].Content = html.UnescapeString(resp.Results.Data[i].Content)
				}
			}
		}
	})
	return resp
}

func Get(rawuri string, callback func(err string, body []byte)) {
	_, err := url.Parse(rawuri)
	if err != nil {
		callback(err.Error(), []byte{})
		return
	}
	response, err := http.Get(rawuri)
	if err != nil {
		callback(err.Error(), []byte{})
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		callback(err.Error(), []byte{})
		return
	}
	callback("", body)
}
