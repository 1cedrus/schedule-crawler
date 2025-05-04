package qldt

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
    QLDT_LOGIN_URL = "https://qldt.ptit.edu.vn/api/auth/login"
    QLDT_DS_TKB_URL = "https://qldt.ptit.edu.vn/api/sch/w-locdstkbtuanusertheohocky"
)

var (
    Cache = make(map[string]any)
)

func FetchToken(r *http.Request) (*TokenResponse, error) {
    user, pass, ok := r.BasicAuth()
    if !ok {
        return nil, errors.New("missing basic auth credentials")
    }

    h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s", user, pass)))

    if (Cache[string(h[:])] != nil) {
        cachedTokenResp, ok := (Cache[string(h[:])]).(*TokenResponse)
        if !ok {
            goto StartFetching
        }

        cachedTime, err := time.Parse(time.RFC1123, cachedTokenResp.ExpiresAt)
        if err != nil {
            return nil, err
        }

        if cachedTime.After(time.Now()) {
            return cachedTokenResp, nil
        }
    }

StartFetching:
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

        Cache[string(h[:])] = &tokenResp;

        return &tokenResp, nil

    default:
        return nil, errors.New("unexpected response code")
    }
}
