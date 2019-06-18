# Go RESTful API Client

Define RESTful API by protobuf. Auto generate go client code.

## Usage

Write `auth.proto`.

```protobuf
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

```plain
-I $GO_PATH/src/github.com/ganlvtech/go-rest-client/rest/
```

Generated `auth.go`.

```go
type Auth_AuthByPasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Auth_AuthByPasswordResponse struct {
	Code     int32  `json:"code"`
	Message  string `json:"message"`
}
type AuthService struct {
	client  *http.Client
	baseUrl string
}
func (s *AuthService) AuthByPassword(in *Auth_AuthByPasswordRequest) (*Auth_AuthByPasswordResponse, *http.Response, error) {
	v := url.Values{}
	v.Add("email", in.Email)
	v.Add("password", in.Password)
	req, err := http.NewRequest("POST", s.baseUrl+"/Auth/AuthByPassword", strings.NewReader(v.Encode()))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	response := &Auth_AuthByPasswordResponse{}
	err = json.NewDecoder(resp.Body).Decode(response)
	if err != nil {
		return response, resp, err
	}
	return response, resp, nil
}
```

## Installation

```bash
go install github.com/ganlvtech/go-rest-client
```

## License

MIT License 
