package dateparser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var zeroTime time.Time

func hasFractional(n float64) (r bool) {
	return float64(int(n)) != n
}

func getFractional(n float64) (r float64) {
	return n - float64(int(n))
}

// This error is returned if the input could not be parsed. Unicode decoding
// errors are returned as they are, not converted into this error.
type ParseError struct {
	Timestr string // The whole input string that could not be parsed.
	Why     string // A textual description of why the parse failed.
	Where   string // The token in the input string that caused the failure.
}

// Returns a string representation of the error.
func (e ParseError) Error() string {
	return fmt.Sprintf("Could not parse date %q: %s (at %q)", e.Timestr, e.Why, e.Where)
}

func parseMS(s string) (secs int, ms int, ok bool) {
	var secs64, ms64 int64
	var err1, err2 error

	if findPeriod(s) < 0 {
		secs64, err1 = strconv.ParseInt(s, 10, 0)
	} else {
		i, f := splitPeriod(s)
		f = (f + "00000")[:6]
		secs64, err1 = strconv.ParseInt(i, 10, 0)
		ms64, err2 = strconv.ParseInt(f, 10, 0)
	}

	return int(secs64), int(ms64), (err1 == nil && err2 == nil)
}

func convertYear(year int) (res int) {
	thisYear := time.Now().Year()

	if year < 100 {
		year += (thisYear / 100) * 100

		if thisYear-year >= 50 {
			year += 100
		} else if year-thisYear >= 50 {
			year -= 100
		}
	}

	return year
}

func isUpper(s string) (r bool) {
	for _, c := range s {
		if !unicode.IsUpper(c) {
			return false
		}
	}

	return true
}

const (
	_JUMP_NONE int = -1 + iota

	_JUMP
)

const (
	_WEEKDAY_NONE int = -1 + iota

	_WEEKDAY_SUN
	_WEEKDAY_MON
	_WEEKDAY_TUE
	_WEEKDAY_WED
	_WEEKDAY_THU
	_WEEKDAY_FRI
	_WEEKDAY_SAT
)

const (
	_MONTH_NONE int = -1 + iota

	_MONTH_ZERO
	_MONTH_JAN
	_MONTH_FEB
	_MONTH_MAR
	_MONTH_APR
	_MONTH_MAY
	_MONTH_JUN
	_MONTH_JUL
	_MONTH_AUG
	_MONTH_SEP
	_MONTH_OCT
	_MONTH_NOV
	_MONTH_DEC
)

const (
	_HMS_NONE int = -1 + iota

	_HMS_HOUR
	_HMS_MINUTE
	_HMS_SECOND
)

const (
	_AMPM_NONE int = -1 + iota

	_AMPM_AM
	_AMPM_PM
)

const (
	_UTCZONE_NONE int = -1 + iota

	_UTCZONE
)

const (
	_PERTAIN_NONE int = -1 + iota

	_PERTAIN
)

var jumpST = stBuild([]stinput{
	stinput{" ", _JUMP},
	stinput{".", _JUMP},
	stinput{",", _JUMP},
	stinput{";", _JUMP},
	stinput{"-", _JUMP},
	stinput{"/", _JUMP},
	stinput{"'", _JUMP},
	stinput{"at", _JUMP},
	stinput{"on", _JUMP},
	stinput{"and", _JUMP},
	stinput{"ad", _JUMP},
	stinput{"m", _JUMP},
	stinput{"t", _JUMP},
	stinput{"of", _JUMP},
	stinput{"st", _JUMP},
	stinput{"nd", _JUMP},
	stinput{"rd", _JUMP},
	stinput{"th", _JUMP},
})

var weekdayST = stBuild([]stinput{
	stinput{"mon", _WEEKDAY_MON},
	stinput{"tue", _WEEKDAY_TUE},
	stinput{"wed", _WEEKDAY_WED},
	stinput{"thu", _WEEKDAY_THU},
	stinput{"fri", _WEEKDAY_FRI},
	stinput{"sat", _WEEKDAY_SAT},
	stinput{"sun", _WEEKDAY_SUN},
	stinput{"monday", _WEEKDAY_MON},
	stinput{"tuesday", _WEEKDAY_TUE},
	stinput{"wednesday", _WEEKDAY_WED},
	stinput{"thursday", _WEEKDAY_THU},
	stinput{"friday", _WEEKDAY_FRI},
	stinput{"saturday", _WEEKDAY_SAT},
	stinput{"sunday", _WEEKDAY_SUN},
})

var monthST = stBuild([]stinput{
	stinput{"jan", _MONTH_JAN},
	stinput{"feb", _MONTH_FEB},
	stinput{"mar", _MONTH_MAR},
	stinput{"apr", _MONTH_APR},
	stinput{"may", _MONTH_MAY},
	stinput{"jun", _MONTH_JUN},
	stinput{"jul", _MONTH_JUL},
	stinput{"aug", _MONTH_AUG},
	stinput{"sep", _MONTH_SEP},
	stinput{"sept", _MONTH_SEP},
	stinput{"oct", _MONTH_OCT},
	stinput{"nov", _MONTH_NOV},
	stinput{"dec", _MONTH_DEC},
	stinput{"january", _MONTH_JAN},
	stinput{"february", _MONTH_FEB},
	stinput{"march", _MONTH_MAR},
	stinput{"april", _MONTH_APR},
	stinput{"may", _MONTH_MAY},
	stinput{"june", _MONTH_JUN},
	stinput{"july", _MONTH_JUL},
	stinput{"august", _MONTH_AUG},
	stinput{"september", _MONTH_SEP},
	stinput{"october", _MONTH_OCT},
	stinput{"november", _MONTH_NOV},
	stinput{"december", _MONTH_DEC},
})

var hmsST = stBuild([]stinput{
	stinput{"h", _HMS_HOUR},
	stinput{"m", _HMS_MINUTE},
	stinput{"s", _HMS_SECOND},
	stinput{"hour", _HMS_HOUR},
	stinput{"minute", _HMS_MINUTE},
	stinput{"second", _HMS_SECOND},
	stinput{"hours", _HMS_HOUR},
	stinput{"minutes", _HMS_MINUTE},
	stinput{"seconds", _HMS_SECOND},
})

var ampmST = stBuild([]stinput{
	stinput{"am", _AMPM_AM},
	stinput{"pm", _AMPM_PM},
	stinput{"a", _AMPM_AM},
	stinput{"p", _AMPM_PM},
})

var utczoneST = stBuild([]stinput{
	stinput{"utc", _UTCZONE},
	stinput{"gmt", _UTCZONE},
	stinput{"z", _UTCZONE},
})

func pertain(s string) (p int) {
	if strings.ToLower(s) == "of" {
		return _PERTAIN
	}

	return _PERTAIN_NONE
}

type parseresult struct {
	Day         int
	Hour        int
	Microsecond int
	Minute      int
	Month       int
	Second      int
	TZName      string
	TZOffset    int
	HasTZOffset bool
	Weekday     int
	Year        int
}

// Contains the settings used in parsing. An empty structure (new(Parser) or
// &Parser{}) is suffice to be able to parse dates.
type Parser struct {
    // The default time (from which components not present in the input are
    // taken from). Defaults to the current time truncated to the nearest day
    // (i.e. midnight today).
	Default time.Time
    
    // Whether or not to perform a fuzzy search. Specifically, invalid tokens
    // are ignored if this is true and cause a parse failure if this is false.
	Fuzzy bool
    
    // In cases where the day and month cannot be distinguished (such as
    // 10-09-2003), this setting is checked to determine the preferred order.
    // With DayFirst set to false, "10-09-2003" will be parsed as "mm-dd-yyyy" (9th Oct 2003).
    // With DayFirst set to true, "10-09-2003" will be parsed as "dd-mm-yyyy" (10th Sep 2003).
	DayFirst bool
    
    // In cases where a 2-digit year cannot be distinguished from the day or
    // month (such as 10-09-03), this setting is checked to determine the
    // preferred order.
    // With YearFirst set to false, "10-09-03" will be parsed as "mm-dd-yy" (9th Oct 2003).
    // With YearFirst set to true, "10-09-03" will be parsed as "yy-mm-dd" (3rd Sep 2010).
	YearFirst bool
    
    // Whether or not to ignore timezones. This only affects the post-processing
    // of the result; timezone information in the input will still be parsed and
    // so still has the chance to trigger parsing errors. If this option is true
    // all returned times have a timezone of UTC+0.
	IgnoreTZ bool
    
    // A map of custom timezone names to their offsets in seconds. If a timezone
    // name is not found in this map, it is looked up using the
    // time.LoadLocation function.
	TZInfos map[string]int
}

// Parses the input string and returns either a parsed date or an error. The
// error may be a ParseError, or an error returned by time.LoadLocation or
// bufio.Reader.ReadRune.
func (parser *Parser) Parse(timestr string) (t time.Time, err error) {
	def := parser.Default

	if def.IsZero() {
		def = time.Now().Truncate(time.Hour * 24).Add(-time.Hour)
	}

	res, err := parser.parseInternal(timestr)
	if err != nil {
		return zeroTime, err
	}

	if res.Year != -1 {
		res.Year = convertYear(res.Year)
	}

	if res.TZOffset == 0 && (res.TZName == "" || res.TZName == "Z") {
		res.TZName = "UTC"
	} else if res.TZOffset != 0 && res.TZName != "" && utczoneST.search(res.TZName) != _UTCZONE_NONE {
		res.TZOffset = 0
	}

	year := res.Year
	if year == -1 {
		year = def.Year()
	}
	month := res.Month
	if month == -1 {
		month = int(def.Month())
	}
	day := res.Day
	if day == -1 {
		day = def.Day()
	}
	hour := res.Hour
	if hour == -1 {
		hour = def.Hour()
	}
	minute := res.Minute
	if minute == -1 {
		minute = def.Minute()
	}
	second := res.Second
	if second == -1 {
		second = def.Second()
	}

	nanosecond := def.Nanosecond()
	if res.Microsecond != -1 {
		nanosecond = res.Microsecond * 1000
	}

	loc := def.Location()

	if parser.IgnoreTZ {
		loc = time.FixedZone("UTC", 0)

	} else {
		if res.TZName != "" {
			if res.HasTZOffset {
				loc = time.FixedZone(res.TZName, res.TZOffset)

			} else {
				ok := false
				if parser.TZInfos != nil {
					var offset int
					offset, ok = parser.TZInfos[res.TZName]
					if ok {
						loc = time.FixedZone(res.TZName, offset)
					}
				}

				if !ok {
					loc, err = time.LoadLocation(res.TZName)
					if err != nil {
						return zeroTime, err
					}
				}
			}

		} else if res.HasTZOffset {
			loc = time.FixedZone("", res.TZOffset)
		}
	}

	t = time.Date(year, time.Month(month), day, hour, minute, second, nanosecond, loc)

	if res.Weekday != -1 && res.Day == -1 {
		weekdayOffset := (res.Weekday - int(t.Weekday()) + 7) % 7
		t = t.Add(time.Duration(weekdayOffset*24) * time.Hour)
	}

	return t, nil
}

func (parser *Parser) parseInternal(timestr string) (res parseresult, err error) {
	lex := newLexer(strings.NewReader(timestr))
	tokens, err := lex.lexAll()
	if err != nil {
		return res, err
	}

	var parseIntResult64 int64

	res = parseresult{
		Day:         -1,
		Hour:        -1,
		Microsecond: -1,
		Minute:      -1,
		Month:       -1,
		Second:      -1,
		TZName:      "",
		TZOffset:    0,
		HasTZOffset: false,
		Weekday:     -1,
		Year:        -1,
	}

	i := 0
	numTokens := len(tokens)
	monthNameIndex := -1
	ymd := make([]int, 0, 3)

loop:
	for i < numTokens {

		token := tokens[i]
		value, err := strconv.ParseFloat(tokens[i], 64)
		isNumber := err == nil

		if isNumber {
			tokenLength := len(token)
			i++

			switch {
			case len(ymd) == 3 &&
				(tokenLength == 2 || tokenLength == 4) &&
				(i >= numTokens ||
					((tokens[i] != ":") && hmsST.search(tokens[i]) == _HMS_NONE)):

				parseIntResult64, err = strconv.ParseInt(token[:2], 10, 0)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", token[:2]}
				}

				res.Hour = int(parseIntResult64)

				if tokenLength == 4 {
					parseIntResult64, err = strconv.ParseInt(token[2:], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", token[2:]}
					}

					res.Minute = int(parseIntResult64)
				}

			case tokenLength == 6 || (tokenLength > 6 && findPeriod(token) == 6):

				if len(ymd) == 0 && findPeriod(token) == -1 {
					parseIntResult64, err = strconv.ParseInt(token[:2], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", token[:2]}
					}

					ymd = append(ymd, convertYear(int(parseIntResult64)))
					parseIntResult64, err = strconv.ParseInt(token[2:4], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", token[2:4]}
					}

					ymd = append(ymd, int(parseIntResult64))
					parseIntResult64, err = strconv.ParseInt(token[4:], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", token[4:]}
					}

					ymd = append(ymd, int(parseIntResult64))

				} else {

					parseIntResult64, err = strconv.ParseInt(token[:2], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", token[:2]}
					}

					res.Hour = int(parseIntResult64)
					parseIntResult64, err = strconv.ParseInt(token[2:4], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", token[2:4]}
					}

					res.Minute = int(parseIntResult64)
					parsems_sec, parsems_ms, parsems_ok := parseMS(token[4:])
					if !parsems_ok {
						return res, ParseError{timestr, "Could not parse sec/ms", token[4:]}
					}
					res.Second = parsems_sec
					res.Microsecond = parsems_ms
				}

			case tokenLength == 8:

				parseIntResult64, err = strconv.ParseInt(token[:4], 10, 0)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", token[:4]}
				}

				ymd = append(ymd, int(parseIntResult64))
				parseIntResult64, err = strconv.ParseInt(token[4:6], 10, 0)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", token[4:6]}
				}

				ymd = append(ymd, int(parseIntResult64))
				parseIntResult64, err = strconv.ParseInt(token[6:], 10, 0)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", token[6:]}
				}

				ymd = append(ymd, int(parseIntResult64))

			case tokenLength == 12 || tokenLength == 14:

				parseIntResult64, err = strconv.ParseInt(token[:4], 10, 0)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", token[:4]}
				}

				ymd = append(ymd, int(parseIntResult64))
				parseIntResult64, err = strconv.ParseInt(token[4:6], 10, 0)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", token[4:6]}
				}

				ymd = append(ymd, int(parseIntResult64))
				parseIntResult64, err = strconv.ParseInt(token[6:8], 10, 0)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", token[6:8]}
				}

				ymd = append(ymd, int(parseIntResult64))
				parseIntResult64, err = strconv.ParseInt(token[8:10], 10, 0)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", token[8:10]}
				}

				res.Hour = int(parseIntResult64)
				parseIntResult64, err = strconv.ParseInt(token[10:12], 10, 0)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", token[10:12]}
				}

				res.Minute = int(parseIntResult64)

				if tokenLength == 14 {
					parseIntResult64, err = strconv.ParseInt(token[12:], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", token[12:]}
					}

					res.Second = int(parseIntResult64)
				}

			case (i < numTokens && hmsST.search(tokens[i]) != _HMS_NONE) ||
				(i+1 < numTokens && tokens[i] == " " && hmsST.search(tokens[i+1]) != _HMS_NONE):

				if tokens[i] == " " {
					i++
				}

				hmstype := hmsST.search(tokens[i])
				for {
					switch hmstype {
					case _HMS_HOUR:
						res.Hour = int(value)
						if hasFractional(value) {
							res.Minute = int(60.0 * getFractional(value))
						}

					case _HMS_MINUTE:
						res.Minute = int(value)
						if hasFractional(value) {
							res.Second = int(60.0 * getFractional(value))
						}

					case _HMS_SECOND:
						parsems_sec, parsems_ms, parsems_ok := parseMS(token)
						if !parsems_ok {
							return res, ParseError{timestr, "Could not parse sec/ms", token}
						}
						res.Second = parsems_sec
						res.Microsecond = parsems_ms
					}

					i++
					if i >= numTokens || hmstype == _HMS_SECOND {
						break
					}

					token = tokens[i]
					value, err = strconv.ParseFloat(tokens[i], 64)
					if err != nil {
						break
					} else {
						i++
						hmstype++
						if i < numTokens {
							newhmstype := hmsST.search(tokens[i])
							if newhmstype != _HMS_NONE {
								hmstype = newhmstype
							}
						}
					}
				}

			case i == numTokens && i >= 3 && tokens[i-2] == " " && hmsST.search(tokens[i-3]) != _HMS_NONE:

				hms := hmsST.search(tokens[i-3]) + 1
				if hms == _HMS_MINUTE {
					res.Minute = int(value)
					if hasFractional(value) {
						res.Second = int(60.0 * getFractional(value))
					}

				} else if hms == _HMS_SECOND {
					parsems_sec, parsems_ms, parsems_ok := parseMS(token)
					if !parsems_ok {
						return res, ParseError{timestr, "Could not parse sec/ms", token}
					}
					res.Second = parsems_sec
					res.Microsecond = parsems_ms
				}

				i++

			case i+1 < numTokens && tokens[i] == ":":

				res.Hour = int(value)
				i++
				value, err = strconv.ParseFloat(tokens[i], 64)
				if err != nil {
					return res, ParseError{timestr, "Could not parse number", tokens[i]}
				}
				res.Minute = int(value)
				if hasFractional(value) {
					res.Second = int(60.0 * getFractional(value))
				}
				i++
				if i < numTokens && tokens[i] == ":" {
					parsems_sec, parsems_ms, parsems_ok := parseMS(tokens[i+1])
					if !parsems_ok {
						return res, ParseError{timestr, "Could not parse sec/ms", tokens[i+1]}
					}
					res.Second = parsems_sec
					res.Microsecond = parsems_ms
					i += 2
				}

			case i < numTokens && (tokens[i] == "-" || tokens[i] == "/" || tokens[i] == "."):
				sep := tokens[i]
				ymd = append(ymd, int(value))
				i++

				if i < numTokens && jumpST.search(tokens[i]) == _JUMP_NONE {
					v, err := strconv.ParseInt(tokens[i], 10, 0)
					if err == nil {

						ymd = append(ymd, int(v))
					} else {

						month := monthST.search(tokens[i])
						if month == _MONTH_NONE {
							return res, ParseError{timestr, "Expected month name", tokens[i]}
						}
						ymd = append(ymd, month)
						if monthNameIndex != -1 {
							return res, ParseError{timestr, "Multiple month names found", tokens[i]}
						}
						monthNameIndex = len(ymd) - 1
					}

					i++
					if i < numTokens && tokens[i] == sep {

						i++
						month := monthST.search(tokens[i])
						if month != _MONTH_NONE {
							ymd = append(ymd, month)
							if monthNameIndex != -1 {
								return res, ParseError{timestr, "Multiple month names found", tokens[i]}
							}
							monthNameIndex = len(ymd)
						} else {
							parseIntResult64, err = strconv.ParseInt(tokens[i], 10, 0)
							if err != nil {
								return res, ParseError{timestr, "Could not parse number", tokens[i]}
							}

							ymd = append(ymd, int(parseIntResult64))
						}
						i++
					}
				}

			case i >= numTokens || jumpST.search(tokens[i]) != _JUMP_NONE:
				if i+1 < numTokens && ampmST.search(tokens[i+1]) != _AMPM_NONE {

					ampm := ampmST.search(tokens[i+1])
					res.Hour = int(value)
					if res.Hour < 12 && ampm == _AMPM_PM {
						res.Hour += 12
					} else if res.Hour == 12 && ampm == _AMPM_AM {
						res.Hour = 0
					}
					i++

				} else {

					ymd = append(ymd, int(value))
				}

			case ampmST.search(tokens[i]) != _AMPM_NONE:

				ampm := ampmST.search(tokens[i])
				res.Hour = int(value)
				if res.Hour < 12 && ampm == _AMPM_PM {
					res.Hour += 12
				} else if res.Hour == 12 && ampm == _AMPM_AM {
					res.Hour = 0
				}
				i++

			default:
				if parser.Fuzzy {
					i++
				} else {
					return res, ParseError{timestr, "Unrecognised token", tokens[i]}
				}
			}

		} else {

			weekday := weekdayST.search(tokens[i])
			if weekday != _WEEKDAY_NONE {
				res.Weekday = weekday
				i++
				continue loop
			}

			month := monthST.search(tokens[i])
			if month != _MONTH_NONE {
				ymd = append(ymd, month)
				if monthNameIndex != -1 {
					return res, ParseError{timestr, "Multiple month names found", tokens[i]}
				}
				monthNameIndex = len(ymd) - 1
				i++

				if i < numTokens {
					sep := tokens[i]
					if sep == "-" || sep == "/" {

						i++
						parseIntResult64, err = strconv.ParseInt(tokens[i], 10, 0)
						if err != nil {
							return res, ParseError{timestr, "Could not parse number", tokens[i]}
						}

						ymd = append(ymd, int(parseIntResult64))
						i++
						if i < numTokens && tokens[i] == sep {

							i++
							parseIntResult64, err = strconv.ParseInt(tokens[i], 10, 0)
							if err != nil {
								return res, ParseError{timestr, "Could not parse number", tokens[i]}
							}

							ymd = append(ymd, int(parseIntResult64))
							i++
						}

					} else if i+3 < numTokens && tokens[i] == " " && tokens[i+2] == " " && pertain(tokens[i+1]) != _PERTAIN_NONE {

						year, err := strconv.ParseInt(tokens[i+3], 10, 0)
						if err == nil {
							ymd = append(ymd, convertYear(int(year)))
						}
						i += 4
					}
				}

				continue loop
			}

			ampm := ampmST.search(tokens[i])
			if ampm != _AMPM_NONE {
				if res.Hour < 12 && ampm == _AMPM_PM {
					res.Hour += 12
				} else if res.Hour == 12 && ampm == _AMPM_AM {
					res.Hour = 0
				}

				i++
				continue loop
			}

			if res.Hour != -1 && len(tokens[i]) <= 5 && isUpper(tokens[i]) {
				res.TZName = tokens[i]
				if utczoneST.search(res.TZName) != _UTCZONE_NONE {
					res.TZOffset = 0
					res.HasTZOffset = true
				}

				i++

				if i < numTokens {
					if tokens[i] == "+" {
						tokens[i] = "-"
						res.HasTZOffset = false
						if utczoneST.search(res.TZName) != _UTCZONE_NONE {
							res.TZName = ""
						}

					} else if tokens[i] == "-" {
						tokens[i] = "+"
						res.HasTZOffset = false
						if utczoneST.search(res.TZName) != _UTCZONE_NONE {
							res.TZName = ""
						}
					}
				}

				continue loop
			}

			if res.Hour != -1 && (tokens[i] == "+" || tokens[i] == "-") {
				sign := -1
				if tokens[i] == "+" {
					sign = 1
				}

				i++
				tokenLength := len(tokens[i])

				if tokenLength == 4 {

					parseIntResult64, err = strconv.ParseInt(tokens[i][:2], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", tokens[i][:2]}
					}

					a := int(parseIntResult64)
					parseIntResult64, err = strconv.ParseInt(tokens[i][2:], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", tokens[i][2:]}
					}

					b := int(parseIntResult64)
					res.TZOffset = a*3600 + b*60
					res.HasTZOffset = true

				} else if i+1 < numTokens && tokens[i+1] == ":" {

					parseIntResult64, err = strconv.ParseInt(tokens[i], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", tokens[i]}
					}

					a := int(parseIntResult64)
					parseIntResult64, err = strconv.ParseInt(tokens[i+2], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", tokens[i+2]}
					}

					b := int(parseIntResult64)
					res.TZOffset = a*3600 + b*60
					res.HasTZOffset = true
					i += 2

				} else if tokenLength <= 2 {

					parseIntResult64, err = strconv.ParseInt(tokens[i], 10, 0)
					if err != nil {
						return res, ParseError{timestr, "Could not parse number", tokens[i]}
					}

					res.TZOffset = int(parseIntResult64) * 3600
					res.HasTZOffset = true

				} else {
					return res, ParseError{timestr, "Bad numbered timezone", tokens[i]}
				}

				i++
				res.TZOffset *= sign

				if i+3 < numTokens && jumpST.search(tokens[i]) != _JUMP_NONE && tokens[i+1] == "(" && tokens[i+3] == ")" && len(tokens[i+2]) >= 3 && len(tokens[i+2]) <= 5 && isUpper(tokens[i+2]) {

					res.TZName = tokens[i+2]
					i += 4
				}

				continue loop
			}

			if jumpST.search(tokens[i]) != _JUMP_NONE {
				i++
				continue loop

			}

			if parser.Fuzzy {
				i++
				continue loop
			} else {
				return res, ParseError{timestr, "Unrecognised token", tokens[i]}
			}
		}
	}

	numYMD := len(ymd)

	switch {
	case numYMD > 3:
		return res, ParseError{timestr, "Too many year/month/day components found", "<no-specific-location>"}

	case numYMD == 1 || (monthNameIndex != -1 && numYMD == 2):
		if monthNameIndex != -1 {
			res.Month = ymd[monthNameIndex]
			copy(ymd[monthNameIndex:], ymd[monthNameIndex+1:])
		}

		if numYMD > 1 || monthNameIndex == -1 {
			if ymd[0] > 31 {
				res.Year = ymd[0]
			} else {
				res.Day = ymd[0]
			}
		}

	case numYMD == 2:

		if ymd[0] > 31 {
			res.Year = ymd[0]
			res.Month = ymd[1]
		} else if ymd[1] > 31 {
			res.Month = ymd[0]
			res.Year = ymd[1]
		} else if parser.DayFirst && ymd[1] <= 12 {
			res.Day = ymd[0]
			res.Month = ymd[1]
		} else {
			res.Month = ymd[0]
			res.Day = ymd[1]
		}

	case numYMD == 3:

		switch monthNameIndex {
		case 0:
			res.Month = ymd[0]
			res.Day = ymd[1]
			res.Year = ymd[2]

		case 1:
			if ymd[0] > 31 || (parser.YearFirst && ymd[2] <= 31) {
				res.Year = ymd[0]
				res.Month = ymd[1]
				res.Day = ymd[2]
			} else {
				res.Day = ymd[0]
				res.Month = ymd[1]
				res.Year = ymd[2]
			}

		case 2:
			if ymd[1] > 31 {
				res.Day = ymd[0]
				res.Year = ymd[1]
				res.Month = ymd[2]
			} else {
				res.Year = ymd[0]
				res.Day = ymd[1]
				res.Month = ymd[2]
			}

		default:
			if ymd[0] > 31 || (parser.YearFirst && ymd[1] <= 12 && ymd[2] <= 31) {
				res.Year = ymd[0]
				res.Month = ymd[1]
				res.Day = ymd[2]
			} else if ymd[0] > 12 || (parser.DayFirst && ymd[1] <= 12) {
				res.Day = ymd[0]
				res.Month = ymd[1]
				res.Year = ymd[2]
			} else {
				res.Month = ymd[0]
				res.Day = ymd[1]
				res.Year = ymd[2]
			}
		}
	}

	return res, nil
}

var defaultParser = new(Parser)

// Parses timestr using a parser with all values at their defaults. To get more
// control over the parsing, create an instance of Parser and use its Parse
// method.
func Parse(timestr string) (t time.Time, err error) {
    return defaultParser.Parse(timestr)
}
