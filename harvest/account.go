package harvest

type Account struct {
	Company *Company `json:"company,omitempty"`
	User    *User    `json:"user,omitempty"`
}

type WeekStartDay string

const (
	Sunday   WeekStartDay = "Sunday"
	Saturday WeekStartDay = "Saturday"
	Monday   WeekStartDay = "Monday"
)

type TimeFormat string

const (
	Decimal      TimeFormat = "decimal"
	HoursMinutes TimeFormat = "hours_minutes"
)

type ClockFormat string

const (
	H12 ClockFormat = "12h"
	H24 ClockFormat = "24h"
)

type DecimalSymbol string

const (
	PeriodDS DecimalSymbol = "."
	CommaDS  DecimalSymbol = ","
)

type ColorScheme string

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

type ThousandsSeparator string

const (
	CommaTS    ThousandsSeparator = ","
	PeriodTS   ThousandsSeparator = "."
	Apostrophe ThousandsSeparator = "'"
	Space      ThousandsSeparator = " "
)

type Modules struct {
	Expenses  bool `json:"expenses,omitempty"`
	Invoices  bool `json:"invoices,omitempty"`
	Estimates bool `json:"estimates,omitempty"`
	Approval  bool `json:"approval,omitempty"`
}

type Company struct {
	BaseUri            string             `json:"base_uri,omitempty"`
	FullDomain         string             `json:"full_domain,omitempty"`
	Name               string             `json:"name,omitempty"`
	Active             bool               `json:"active,omitempty"`
	WeekStartDay       WeekStartDay       `json:"week_start_day,omitempty"`
	TimeFormat         TimeFormat         `json:"time_format,omitempty"`
	Clock              ClockFormat        `json:"clock,omitempty"`
	DecimalSymbol      DecimalSymbol      `json:"decimal_symbol,omitempty"`
	ColorScheme        ColorScheme        `json:"color_scheme,omitempty"`
	Modules            *Modules           `json:"modules,omitempty"`
	ThousandsSeparator ThousandsSeparator `json:"thousands_separator,omitempty"`
}
