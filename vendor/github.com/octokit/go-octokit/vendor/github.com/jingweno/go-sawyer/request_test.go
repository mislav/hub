package sawyer

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bmizerany/assert"
	"github.com/jingweno/go-sawyer/mediatype"
)

func TestSuccessfulGet(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		head := w.Header()
		head.Set("Content-Type", "application/json")
		link := `<https://api.github.com/user/repos?page=3&per_page=100>; rel="next", <https://api.github.com/user/repos?page=50&per_page=100>; rel="last"`
		head.Set("Link", link)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "login": "sawyer"}`))
	})

	client := setup.Client
	user := &TestUser{}

	req, err := client.NewRequest("user")
	assert.Equal(t, nil, err)

	res := req.Get()
	assert.Equal(t, false, res.IsError())
	assert.Equal(t, false, res.IsApiError())

	assert.Equal(t, nil, res.Decode(user))
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "sawyer", user.Login)

	mheader := res.MediaHeader
	assert.Equal(t, "https://api.github.com/user/repos?page=3&per_page=100", string(mheader.Relations["next"]))
	assert.Equal(t, "https://api.github.com/user/repos?page=50&per_page=100", string(mheader.Relations["last"]))
}

func TestSuccessfulGetWithoutOutput(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "login": "sawyer"}`))
	})

	client := setup.Client
	user := &TestUser{}

	req, err := client.NewRequest("user")
	assert.Equal(t, nil, err)

	res := req.Get()
	assert.Equal(t, false, res.IsError())
	assert.Equal(t, false, res.IsApiError())

	assert.Tf(t, !res.IsError(), "Response shouldn't have error")
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, false, res.BodyClosed)
	assert.Equal(t, 0, user.Id)
	assert.Equal(t, "", user.Login)

	dec := json.NewDecoder(res.Body)
	dec.Decode(user)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "sawyer", user.Login)
}

func TestSuccessfulGetWithoutDecoder(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		head := w.Header()
		head.Set("Content-Type", "application/booya+booya")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id": 1, "login": "sawyer"}`))
	})

	client := setup.Client
	user := &TestUser{}

	req, err := client.NewRequest("user")
	assert.Equal(t, nil, err)

	res := req.Get()
	assert.NotEqual(t, nil, res.Decode(user), "response should have decoder error")
	assert.Tf(t, strings.HasPrefix(res.Error(), "No decoder found for format booya"), "Bad error: %s", res.Error())
}

func TestSuccessfulPost(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	mtype, err := mediatype.Parse("application/json")

	setup.Mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, mtype.String(), r.Header.Get("Content-Type"))

		user := &TestUser{}
		mtype.Decode(user, r.Body)
		assert.Equal(t, "sawyer", user.Login)

		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"login": "sawyer2"}`))
	})

	client := setup.Client
	user := &TestUser{}

	req, err := client.NewRequest("users")
	assert.Equal(t, nil, err)

	user.Login = "sawyer"
	req.SetBody(mtype, user)
	res := req.Post()
	assert.Equal(t, false, res.IsError())
	assert.Equal(t, false, res.IsApiError())
	assert.Equal(t, nil, res.Decode(user))

	assert.Equal(t, nil, err)
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, "sawyer2", user.Login)
	assert.Equal(t, true, res.BodyClosed)
}

func TestErrorResponse(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/404", func(w http.ResponseWriter, r *http.Request) {
		head := w.Header()
		head.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	})

	client := setup.Client
	user := &TestUser{}
	apierr := &TestError{}

	req, err := client.NewRequest("404")
	if err != nil {
		t.Fatalf("request errored: %s", err)
	}

	res := req.Get()
	assert.Equal(t, true, res.IsApiError())
	assert.Equal(t, false, res.IsError())
	assert.Equal(t, nil, res.Decode(apierr))

	assert.Equal(t, 404, res.StatusCode)
	assert.Equal(t, 0, user.Id)
	assert.Equal(t, "", user.Login)
	assert.Equal(t, "not found", apierr.Message)
	assert.Equal(t, true, res.BodyClosed)
}

func TestResolveRequestQuery(t *testing.T) {
	setup := Setup(t)
	defer setup.Teardown()

	setup.Mux.HandleFunc("/q", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "1", q.Get("a"))
		assert.Equal(t, "4", q.Get("b"))
		assert.Equal(t, "3", q.Get("c"))
		assert.Equal(t, "2", q.Get("d"))
		assert.Equal(t, "1", q.Get("e"))
		w.WriteHeader(123)
		w.Write([]byte("ok"))
	})

	assert.Equal(t, "1", setup.Client.Query.Get("a"))
	assert.Equal(t, "1", setup.Client.Query.Get("b"))

	setup.Client.Query.Set("b", "2")
	setup.Client.Query.Set("c", "3")

	req, err := setup.Client.NewRequest("/q?d=4")
	assert.Equal(t, nil, err)

	req.Query.Set("b", "4")
	req.Query.Set("c", "3")
	req.Query.Set("d", "2")
	req.Query.Set("e", "1")

	res := req.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, 123, res.StatusCode)
}

type TestUser struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
}

type TestError struct {
	Message string `json:"message"`
}

type SetupServer struct {
	Client *Client
	Server *httptest.Server
	Mux    *http.ServeMux
}

func Setup(t *testing.T) *SetupServer {
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	client, err := NewFromString(srv.URL+"?a=1&b=1", nil)
	assert.Equalf(t, nil, err, "Unable to parse %s", srv.URL)

	return &SetupServer{client, srv, mux}
}

func (s *SetupServer) Teardown() {
	s.Server.Close()
}
