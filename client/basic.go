package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

///*
//Summary
//URL: https://prolifewhistleblower.com/anonymous-form/
//Status: 200
//Source: Network
//Address: 192.124.249.104:443
//
//Request
//:method: GET
//:scheme: https
//:authority: prolifewhistleblower.com
//:path: /anonymous-form/
//Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
//Cookie: _ga_M5WWGNMLR8=GS1.1.1630641978.1.1.1630643693.0; _ga=GA1.1.1216678778.1630641978; hustle_module_show_count-social_sharing-1=1; sucuri_cloudproxy_uuid_40466e6c3=349feba0a3b1f9e0d71b82a857e885f7
//Referer: https://prolifewhistleblower.com/anonymous-form/
//Cache-Control: max-age=0
//Host: prolifewhistleblower.com
//User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15
//Accept-Language: en-us
//Accept-Encoding: gzip, deflate, br
//Connection: keep-alive
//
//Response
//:status: 200
//Content-Type: text/html; charset=UTF-8
//X-Content-Type-Options: nosniff
//Content-Security-Policy: upgrade-insecure-requests;
//Date: Fri, 03 Sep 2021 05:15:02 GMT
//X-Frame-Options: SAMEORIGIN
//X-XSS-Protection: 1; mode=block
//Link: <https://prolifewhistleblower.com/wp-json/>; rel="https://api.w.org/", <https://prolifewhistleblower.com/wp-json/wp/v2/pages/27>; rel="alternate"; type="application/json", <https://prolifewhistleblower.com/?p=27>; rel=shortlink
//x-sucuri-id: 12004
//Server: nginx
//x-sucuri-cache: HIT
//*/
func New(ctx context.Context, url string, log *logrus.Entry) *C {
	return &C{http.DefaultClient, url}
}

type C struct {
	*http.Client
	URL string
}

func (c *C) Post(formdata *url.Values) (*http.Response, error) {
	return c.PostForm(c.URL, *formdata)
}

func (c *C) Forever(concurrency int) chan *url.Values {
	q := make(chan *url.Values)
	for i := 0; i < concurrency; i++ {
		for data := range q {
			res, err := c.PostForm(c.URL, *data)
			if err != nil {
				continue
			}
			println(res)
		}
	}
	return q
}
