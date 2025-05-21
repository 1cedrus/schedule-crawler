package qldt

import (
	"bytes"
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
    QLDT_EXAM_URL = "https://qldt.ptit.edu.vn/api/epm/w-locdslichthisvtheohocky"
)

var (
    Cache = make(map[string]any)
)

func FetchLichThi(accessToken string, name string) (*LichThiResponse, error) {
    h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s", name, "DSLichThi")))

    if Cache[string(h[:])] != nil {
        cachedDSLichThi, ok := Cache[string(h[:])].(*LichThiResponse)
        if !ok {
            goto StartFetching
        }

        return cachedDSLichThi, nil
    }

StartFetching:
    var reqBody LichThiRequestBody
    reqBody.Filter.HocKy = "20242" // Fix this to the current semester
    reqBody.Filter.IsGiuaHocKy = false
    reqBody.Additional.Paging.Limit = 100
    reqBody.Additional.Paging.Page = 1
    reqBody.Additional.Ordering = []Ordering{}

    reqBodyBytes, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", QLDT_EXAM_URL, bytes.NewReader(reqBodyBytes))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Set("Content-Type", "application/json")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    data, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    var lichThiResponse LichThiResponse
    if err := json.Unmarshal(data, &lichThiResponse); err != nil {
        return nil, err
    }

    Cache[string(h[:])] = &lichThiResponse

    return &lichThiResponse, nil
}

func FetchDSTKB(accessToken string, name string) (*ScheduleResponse, error) {
    h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s", name, "DSTKB")))

    if Cache[string(h[:])] != nil {
        cachedDSTKB, ok := Cache[string(h[:])].(*ScheduleResponse)
        if !ok {
            goto StartFetching
        }

        return cachedDSTKB, nil
    }

StartFetching:
    var reqBody ScheduleRequestBody
    reqBody.Filter.HocKy = "20242" // Fix this to the current semester
    reqBody.Filter.TenHocKy = ""

    reqBodyBytes, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", QLDT_DS_TKB_URL, bytes.NewReader(reqBodyBytes))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
    req.Header.Set("Content-Type", "application/json")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    data, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    var scheduleResp ScheduleResponse
    if err := json.Unmarshal(data, &scheduleResp); err != nil {
        return nil, err
    }

    Cache[string(h[:])] = &scheduleResp

    return &scheduleResp, nil
}

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

    data, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    var baseResp Response
    if err := json.Unmarshal(data, &baseResp); err != nil {
		return nil, err
	}

    switch baseResp.Code {
    case "403":
        var errResp TokenErrorResponse
        if err := json.Unmarshal(data, &errResp); err != nil {
			return nil, err
		}
        return nil, errors.New(errResp.Message)

    case "200":
        var tokenResp TokenResponse
        if err := json.Unmarshal(data, &tokenResp); err != nil {
            return nil, err
        }

        Cache[string(h[:])] = &tokenResp;

        return &tokenResp, nil

    default:
        return nil, errors.New("unexpected response code")
    }
}
