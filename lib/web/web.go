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
	var uri string = fmt.Sprintf("http://ajax.googleapis.com/ajax/services/search/web?v=1.0&rsz=%d&q=%s",
		results, url.QueryEscape(searchTerm))
	body, err := Get(&uri)
	if err != "" {
		resp.Error = err
		return resp
	}
	errr := json.Unmarshal(body, &resp)
	if errr != nil {
		resp.Error = errr.Error()
		return resp
	}
	if len(resp.Results.Data) == 0 {
		resp.Error = fmt.Sprintf("Google couldn't find \"%s\"", searchTerm)
		return resp
	}
	for i, _ := range resp.Results.Data {
		if resp.Results.Data[i].Title != "" {
			resp.Results.Data[i].Title = html.UnescapeString(resp.Results.Data[i].Title)
		}
		if resp.Results.Data[i].Content != "" {
			resp.Results.Data[i].Content = html.UnescapeString(resp.Results.Data[i].Content)
		}
	}
	return resp
}

func Get(rawuri *string) ([]byte, string) {
	_, err := url.Parse(*rawuri)
	if err != nil {
		return nil, err.Error()
	}
	response, err := http.Get(*rawuri)
	if err != nil {
		return nil, err.Error()
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err.Error()
	}
	return body, ""
}
