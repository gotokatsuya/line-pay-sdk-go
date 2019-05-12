package linepay

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	testEndpointBase := "/test-api-pay"

	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(testEndpointBase+"/", http.StripPrefix(testEndpointBase, mux))
	server := httptest.NewServer(apiHandler)

	client, err := New("testid", "testsecret", WithEndpointBase(server.URL+testEndpointBase+"/"))
	if err != nil {
		panic(err)
	}
	return client, mux, server.URL, server.Close
}

func TestNew(t *testing.T) {
	id := "testid"
	secret := "testsecret"
	wantURL, _ := url.Parse(APIEndpointBaseReal)
	client, err := New(id, secret)
	if err != nil {
		t.Fatal(err)
	}
	if client.channelID != id {
		t.Errorf("channelID %s; want %s", client.channelID, id)
	}
	if client.channelSecret != secret {
		t.Errorf("channelSecret %s; want %s", client.channelSecret, secret)
	}
	if !reflect.DeepEqual(client.endpointBase, wantURL) {
		t.Errorf("endpointBase %v; want %v", client.endpointBase, wantURL)
	}
	if client.httpClient != http.DefaultClient {
		t.Errorf("httpClient %p; want %p", client.httpClient, http.DefaultClient)
	}
}

func TestNewWithOptions(t *testing.T) {
	id := "testid"
	secret := "testsecret"
	endpoint := "https://example.test/"
	httpClient := http.Client{}
	wantURL, _ := url.Parse(endpoint)
	client, err := New(
		id,
		secret,
		WithHTTPClient(&httpClient),
		WithEndpointBase(endpoint),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(client.endpointBase, wantURL) {
		t.Errorf("endpointBase %v; want %v", client.endpointBase, wantURL)
	}
	if client.httpClient != &httpClient {
		t.Errorf("httpClient %p; want %p", client.httpClient, &httpClient)
	}
}

func Test_mergeQuery(t *testing.T) {
	inURL, outURL := "foo", "foo?p=q"
	u, err := mergeQuery(inURL, &struct {
		P string `url:"p"`
	}{"q"})
	if err != nil {
		t.Fatal(err)
	}
	// test that url was merged
	if got, want := u, outURL; got != want {
		t.Errorf("TestMergeQuery(%q) URL is %v, want %v", inURL, got, want)
	}
}

func TestClient_NewRequest(t *testing.T) {
	id := "testid"
	secret := "testsecret"
	client, err := New(id, secret)
	if err != nil {
		t.Fatal(err)
	}

	inURL, outURL := "foo", APIEndpointBaseReal+"foo"
	inBody, outBody := &struct{ Login string }{"l"}, `{"Login":"l"}`+"\n"
	req, _ := client.NewRequest("GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%q) Body is %v, want %v", inBody, got, want)
	}
}

func TestClient_Do(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/do", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != http.MethodGet {
			t.Errorf("Request method: %v, want %v", got, http.MethodGet)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "do", nil)
	body := new(foo)
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatal(err)
	}
	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}
