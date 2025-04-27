package device

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type measurementsQuery struct {
	deviceId int
	start    time.Time
	end      time.Time
}

func getMeasurementsQueryParams(r *http.Request) (*measurementsQuery, error) {
	rawQuery := r.URL.RawQuery
	decodedQuery, err := url.QueryUnescape(rawQuery)
	if err != nil {
		return nil, errors.New("error decoding query parameters")
	}
	r.URL.RawQuery = decodedQuery

	startStr := r.URL.Query().Get("start")
	var start time.Time
	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			return nil, errors.New("invalid start date format. Use RFC3339 (YYYY-MM-DDTHH:MM:SSZ)")
		}
	}

	endStr := r.URL.Query().Get("end")
	end := time.Now().UTC()
	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			return nil, errors.New("invalid end date format. Use RFC3339 (YYYY-MM-DDTHH:MM:SSZ)")
		}
	}

	deviceId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return nil, errors.New("invalid device ID")
	}

	return &measurementsQuery{deviceId, start, end}, nil
}
