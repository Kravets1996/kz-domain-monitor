package api

import (
	"strconv"
	"time"
)

type PsResponse struct {
	Result string `json:"result"`
	Answer struct {
		Expire struct {
			Unix string `json:"unix"`
		} `json:"expire"`
	} `json:"answer"`
}

func (response PsResponse) GetExpireTime() time.Time {
	sec, _ := strconv.ParseInt(response.Answer.Expire.Unix, 10, 64)
	return time.Unix(sec, 0)
}
