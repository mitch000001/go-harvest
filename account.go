package main

import (
	"strconv"
)

type Account struct {
	Company *Company `json:"company"`
	User    *User    `json:"user"`
}

type WeekStartDay string

func (w *WeekStartDay) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		*w = WeekStartDay(unquotedData)
	} else {
		*w = WeekStartDay(data)
	}
	return nil
}

const (
	Sunday   WeekStartDay = "Sunday"
	Saturday WeekStartDay = "Saturday"
	Monday   WeekStartDay = "Monday"
)

type TimeFormat string

func (t *TimeFormat) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		*t = TimeFormat(unquotedData)
	} else {
		*t = TimeFormat(data)
	}
	return nil
}

const (
	Decimal      TimeFormat = "decimal"
	HoursMinutes TimeFormat = "hours_minutes"
)

type ClockFormat string

func (c *ClockFormat) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		*c = ClockFormat(unquotedData)
	} else {
		*c = ClockFormat(data)
	}
	return nil
}

const (
	H12 ClockFormat = "12h"
	H24 ClockFormat = "24h"
)

type DecimalSymbol rune

func (d *DecimalSymbol) UnmarshalJSON(data []byte) error {
	unquotedData, _, _, err := strconv.UnquoteChar(string(data), byte('"'))
	if err != nil {
		*d = DecimalSymbol(unquotedData)
	} else {
		return err
	}
	return nil
}

const (
	PeriodDS DecimalSymbol = '.'
	CommaDS  DecimalSymbol = ','
)

type ColorScheme string

func (c *ColorScheme) UnmarshalJSON(data []byte) error {
	unquotedData, err := strconv.Unquote(string(data))
	if err != nil {
		*c = ColorScheme(unquotedData)
	} else {
		*c = ColorScheme(data)
	}
	return nil
}

const (
	Orange  ColorScheme = "orange"
	Spring  ColorScheme = "spring"
	Green   ColorScheme = "green"
	Legacy  ColorScheme = "legacy"
	Behance ColorScheme = "behance"
	Blue    ColorScheme = "blue"
	Purple  ColorScheme = "purple"
	Red     ColorScheme = "red"
	LtGrey  ColorScheme = "lt_grey"
	Gray    ColorScheme = "gray"
)

type ThousandsSeparator rune

func (t *ThousandsSeparator) UnmarshalJSON(data []byte) error {
	unquotedData, _, _, err := strconv.UnquoteChar(string(data), byte('"'))
	if err != nil {
		*t = ThousandsSeparator(unquotedData)
	} else {
		return err
	}
	return nil
}

const (
	CommaTS    ThousandsSeparator = ','
	PeriodTS   ThousandsSeparator = '.'
	Apostrophe ThousandsSeparator = '\''
	Space      ThousandsSeparator = ' '
)

type Modules struct {
	Expenses  bool `json:"expenses"`
	Invoices  bool `json:"invoices"`
	Estimates bool `json:"estimates"`
	Approval  bool `json:"approval"`
}

type Company struct {
	BaseUri            string             `json:"base_uri"`
	FullDomain         string             `json:"full_domain"`
	Name               string             `json:"name"`
	Active             bool               `json:"active"`
	WeekStartDay       WeekStartDay       `json:"week_start_day"`
	TimeFormat         TimeFormat         `json:"time_format"`
	Clock              ClockFormat        `json:"clock"`
	DecimalSymbol      DecimalSymbol      `json:"decimal_symbol"`
	ColorScheme        ColorScheme        `json:"color_scheme"`
	Modules            *Modules           `json:"modules"`
	ThousandsSeparator ThousandsSeparator `json:"thousands_separator"`
}
