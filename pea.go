/*
Shoot many peas. THey they are shall, but they are many.
*/
package peashooter

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	/// "github.com/k4s/webrowser"
	"github.com/k4s/webrowser"
	. "github.com/k4s/webrowser"
	"github.com/sirupsen/logrus"
)

const (
	site = "https://prolifewhistleblower.com"
)

func FakeValues() string {
	v := url.Values{}
	v.Add(
		"_ga", "GA1.1.1216678778.1630641978",
	)
	v.Add(
		"_ga_M5WWGNMLR8", "GS1.1.1630711553.4.1.1630711738.0",
	)
	v.Add(
		"hustle_module_show_count-social_sharing-1", "11",
	)
	v.Add(
		"sucuri_cloudproxy_uuid_40466e6c3", "349feba0a3b1f9e0d71b82a857e885f",
	)
	v.Add(
		"sucuricp_tfca_6e453141ae697f9f78b18427b4c54df1", "1",
	)
	return v.Encode()
}

func Fake(ctx context.Context, host string, req *http.Request, concurrency int, log *logrus.Entry) (*http.Response, error) {

	req.Header.Add(
		"Cookie",
		FakeValues(),
	)
	data := &Param{
		Method:       "POST",
		Url:          site,
		Header:       req.Header,
		UsePhantomJS: true,
	}

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
