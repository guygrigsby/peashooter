package peashooter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/guygrigsby/peashooter/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

const (
	VIDEO = "https://www.xvideos.com/video58146793/furry_fuck_in_car_porn_animation"
)

func TestPost(t *testing.T) {
	log := logrus.New()
	tests := []struct {
		name string
		url  string
		data url.Values
	}{
		{
			"hp",
			"https://prolifewhistleblower.com/anonymous-form/",
			map[string][]string{
				"video": {VIDEO},
				"f":     {},
			},
		},
	}
	for _, tc := range tests {
		client := client.New(context.Background(), tc.url, log.WithField("service", "test"))
		res, err := client.Post(&tc.data)
		require.NoError(t, err)
		fmt.Printf("Response %+v\n", res)
		require.Equal(t, res.StatusCode, http.StatusOK)
		b, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		res.Body.Close()
		fmt.Println("body", string(b))
	}
}
