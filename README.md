# Go RESTful API Client

Define RESTful API by protobuf. Auto generate go client code.

## Install

```bash
go get -u github.com/ganlvtech/go-rest-client
go install github.com/ganlvtech/go-rest-client/protoc-gen-gorestclient
```

You should copy or link `$GOPATH/src/github.com/ganlvtech/go-rest-client/protoc-gen-gorestclient/rest/rest.proto` to protobuf's include directory. Or you must `-I` or `--proto_path` to specify import path.

## Usage

Write `auth.proto`.

```protobuf
syntax = "proto3";
package kahla;

import "rest.proto";
import "google/protobuf/empty.proto";

service Auth {
    rpc AuthByPassword (Auth_AuthByPasswordRequest) returns (Auth_AuthByPasswordResponse) {
        option (rest.method) = POST;
    };
}
message Auth_AuthByPasswordRequest {
    string email = 1;
    string password = 2;
}
message Auth_AuthByPasswordResponse {
    sint32 code = 1;
    string message = 2;
}
```

Run protobuf compiler with `gorestclient` plugin.

```bash
protoc *.proto --gorestclient_out=.
```

You may need to add some import path 

```plain
-I ./ -I $GOPATH/src/github.com/ganlvtech/go-rest-client/protoc-gen-gorestclient/rest/
```

Generated `auth.go`.

```go
package kahla

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type Auth_AuthByPasswordRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Auth_AuthByPasswordResponse struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type AuthService struct {
	client  *http.Client
	baseUrl string
}

func NewAuthService(client *http.Client, baseUrl string) *AuthService {
	return &AuthService{
		client:  client,
		baseUrl: baseUrl,
	}
}
func (s *AuthService) AuthByPassword(in *Auth_AuthByPasswordRequest) (out *Auth_AuthByPasswordResponse, resp *http.Response, err error) {
	v := url.Values{}
	v.Add("email", in.Email)
	v.Add("password", in.Password)
	req, err := http.NewRequest("POST", s.baseUrl+"/Auth/AuthByPassword", strings.NewReader(v.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err = s.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	out = &Auth_AuthByPasswordResponse{}
	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil {
		return
	}
	return
}
```

## License

MIT License 
