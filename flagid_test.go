package saarflagid

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetIDsFromStatus(t *testing.T) {
	for _, tt := range commonTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIDsFromStatus(tt.service, tt.teamIP, []byte(testData))
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIDsFromStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIDsFromStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIDsFromURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(testData))
	}))
	defer ts.Close()

	for _, tt := range commonTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIDsFromURL(tt.service, tt.teamIP, ts.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIDsFromURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetIDsFromURL() got = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("invalid url", func(t *testing.T) {
		_, err := GetIDsFromURL("a", "b", "http://1.2.3.4.5.6.7/")
		if err == nil {
			t.Error("GetIDsFromURL() not returned error on invalid url")
			return
		}
	})

	tsTimeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(11 * time.Second)
		_, _ = w.Write([]byte(testData))
	}))
	defer func() { tsTimeout.CloseClientConnections(); tsTimeout.Close() }()

	t.Run("timeout", func(t *testing.T) {
		_, err := GetIDsFromURL("a", "b", tsTimeout.URL)
		if err == nil {
			t.Error("GetIDsFromURL() not returned error on timeout")
			return
		}
	})
}

const testData = `{
    "teams": [
        {
            "id": 1,
            "name": "NOP",
            "ip": "10.32.1.2"
        },
        {
            "id": 2,
            "name": "saarsec",
            "ip": "10.32.2.2"
        }
    ],
    "flag_ids": {
        "service_1": {
            "10.32.1.2": {
                "15": ["username1", "username1.2"],
                "16": ["username2", "username2.2"]
            },
            "10.32.2.2": {
                "15": ["username3", "username3.2"],
                "16": ["username4", "username4.2"]
            }
        }
    }
}`

var commonTests = []struct {
	name    string
	service string
	teamIP  string
	want    []string
	wantErr bool
}{
	{
		name:    "the example from site for team1",
		service: "service_1",
		teamIP:  "10.32.1.2",
		want:    []string{"username1", "username1.2", "username2", "username2.2"},
		wantErr: false,
	},

	{
		name:    "the example from site for team2",
		service: "service_1",
		teamIP:  "10.32.2.2",
		want:    []string{"username3", "username3.2", "username4", "username4.2"},
		wantErr: false,
	},

	{
		name:    "the example from site for fake team",
		service: "service_1",
		teamIP:  "10.32.1.2.3",
		wantErr: true,
	},

	{
		name:    "the example from site for fake service",
		service: "service_fuck",
		teamIP:  "10.32.1.2",
		wantErr: true,
	},
}
