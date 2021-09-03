/*
	Shoot many peas. THey they are shall, but they are many.
*/
package peashooter

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	/// "github.com/k4s/webrowser"
	"github.com/k4s/webrowser"
	. "github.com/k4s/webrowser"
)

const (
	site = "https://prolifewhistleblower.com"
)

type P struct {
	*webrowser.Param
}
type R struct {
	*http.Request
}

func (r *R) GetConnTimeout() time.Duration {
	return 600 / time.Millisecond
}
func (r *R) GetDialTimeout() time.Duration {
	return 600 / time.Millisecond
}
func (r *R) GetHeader() http.Header {
	return http.Header{}
}
func (r *R) GetMethod() string {
	return r.GetMethod()
}

func Fake(ctx context.Context, host string, req *http.Request) (*http.Response, error) {
	data := &Param{
		Method:       "POST",
		Url:          site,
		Header:       http.Header{"Cookie": []string{"your cookie"}},
		UsePhantomJS: true,
	}
	//data.Set("Header" http.Header{"Cookie": []string{"your cookie"}}

	browser := webrowser.NewWebrowse()
	resp, err := browser.Download(data)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(body))
	fmt.Println(resp.Cookies())
	return resp, err
}
