package qldt

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const QLDT_LOGIN_URL = "https://qldt.ptit.edu.vn/api/auth/login"

func FetchToken(r *http.Request) (*TokenResponse, error) {
    user, pass, ok := r.BasicAuth()
    if !ok {
        return nil, errors.New("missing basic auth credentials")
    }

    f := url.Values{
        "username": {user},
        "password": {pass},
        "grant_type": {"password"},
    }

    res, err := http.Post(QLDT_LOGIN_URL, "application/x-www-form-urlencoded", strings.NewReader(f.Encode()))
    if err != nil {
        return nil, err
    }

    defer res.Body.Close()

    bytes, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    var baseResp Response
    if err := json.Unmarshal(bytes, &baseResp); err != nil {
		return nil, err
	}

    switch baseResp.Code {
    case "403":
        var errResp TokenErrorResponse
        if err := json.Unmarshal(bytes, &errResp); err != nil {
			return nil, err
		}
        return nil, errors.New(errResp.Message)

    case "200":
        var tokenResp TokenResponse
        if err := json.Unmarshal(bytes, &tokenResp); err != nil {
            return nil, err
        }
        return &tokenResp, nil

    default:
        return nil, errors.New("unexpected response code")
    }
}
