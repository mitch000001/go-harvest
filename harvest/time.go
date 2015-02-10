package harvest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ShortDate struct {
	time.Time
}

func NewShortDate(date time.Time) ShortDate {
	return Date(date.Year(), date.Month(), date.Day(), date.Location())
}

func Date(year int, month time.Month, day int, location *time.Location) ShortDate {
	return ShortDate{time.Date(year, month, day, 0, 0, 0, 0, time.UTC)}
}

func (date *ShortDate) MarshalJSON() ([]byte, error) {
	if date.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(date.Format("2006-01-02"))
}

func (date *ShortDate) UnmarshalJSON(data []byte) error {
	unquotedData, _ := strconv.Unquote(string(data))
	time, err := time.Parse("2006-01-02", unquotedData)
	date.Time = time
	return err
}

type Timeframe struct {
	StartDate ShortDate
	EndDate   ShortDate
}

// From returns a Timeframe with the StartDate set to date and the EndDate set to today.
// The EndDate will use the same timezone location as provided in StartDate
func From(date ShortDate) Timeframe {
	endDate := NewShortDate(time.Now().In(date.Location()))
	return Timeframe{date, endDate}
}

func (tf Timeframe) MarshalJSON() ([]byte, error) {
	if tf.StartDate.IsZero() || tf.EndDate.IsZero() {
		return json.Marshal("")
	}
	return json.Marshal(fmt.Sprintf("%s,%s", tf.StartDate.Format("2006-01-02"), tf.EndDate.Format("2006-01-02")))
}

func (tf *Timeframe) UnmarshalJSON(data []byte) error {
	unquotedData, _ := strconv.Unquote(string(data))
	dates := strings.Split(unquotedData, ",")
	if len(dates) != 2 {
		*tf = Timeframe{}
		return nil
	}
	startTime, err1 := time.Parse("2006-01-02", dates[0])
	startDate := ShortDate{startTime}
	endTime, err2 := time.Parse("2006-01-02", dates[1])
	endDate := ShortDate{endTime}
	if err1 != nil || err2 != nil {
		*tf = Timeframe{}
		return nil
	}
	*tf = Timeframe{StartDate: startDate, EndDate: endDate}
	return nil
}
