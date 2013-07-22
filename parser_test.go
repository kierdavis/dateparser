package dateparser

import (
    "fmt"
    "math/rand"
    "testing"
    "time"
)

var TestTZInfos = map[string]int{"BRST": -10800}
var UTCLoc = time.FixedZone("UTC", 0)
var BRSTLoc = time.FixedZone("BRST", -10800)
var TestDefault = time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
                                        
func check(t *testing.T, parser *Parser, timestr string, expect time.Time) {
    res, err := parser.Parse(timestr)
    if err != nil {
        t.Fail()
        t.Logf("Parse failure: %s", err.Error())
    }
    
    if !res.Equal(expect) {
        t.Fail()
        t.Logf("Expected '%s',", expect)
        t.Logf("  parsed '%s'", res)
        t.Logf("when parsing '%s'", timestr)
    }
}

func TestDateCommandFormat(t *testing.T) {
    parser := &Parser{TZInfos: TestTZInfos}
    timestr := "Thu Sep 25 10:36:28 BRST 2003"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, BRSTLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatReversed(t *testing.T) {
    parser := &Parser{TZInfos: TestTZInfos}
    timestr := "2003 10:36:28 BRST 25 Sep Thu"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, BRSTLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatIgnoreTZ(t *testing.T) {
    parser := &Parser{IgnoreTZ: true}
    timestr := "Thu Sep 25 10:36:28 BRST 2003"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip1(t *testing.T) {
    parser := &Parser{}
    timestr := "Thu Sep 25 10:36:28 2003"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip2(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Thu Sep 25 10:36:28"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip3(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Thu Sep 10:36:28"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, UTCLoc)
    check(t, parser, timestr, expect)
}   

func TestDateCommandFormatStrip4(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Thu 10:36:28"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip5(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Sep 10:36:28"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip6(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:36:28"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip7(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:36"
    expect := time.Date(2003, 9, 25, 10, 36, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip8(t *testing.T) {
    parser := &Parser{}
    timestr := "Thu Sep 25 2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip9(t *testing.T) {
    parser := &Parser{}
    timestr := "Sep 25 2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip10(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Sep 2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip11(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Sep"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateCommandFormatStrip12(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateRCommandFormat(t *testing.T) {
    parser := &Parser{}
    timestr := "Thu, 25 Sep 2003 10:49:41 -0300"
    expect := time.Date(2003, 9, 25, 10, 49, 41, 0, BRSTLoc)
    check(t, parser, timestr, expect)
}

func TestISOFormat(t *testing.T) {
    parser := &Parser{}
    timestr := "2003-09-25T10:49:41.5-03:00"
    expect := time.Date(2003, 9, 25, 10, 49, 41, 500000000, BRSTLoc)
    check(t, parser, timestr, expect)
}

func TestISOFormatStrip1(t *testing.T) {
    parser := &Parser{}
    timestr := "2003-09-25T10:49:41-03:00"
    expect := time.Date(2003, 9, 25, 10, 49, 41, 0, BRSTLoc)
    check(t, parser, timestr, expect)
}

func TestISOFormatStrip2(t *testing.T) {
    parser := &Parser{}
    timestr := "2003-09-25T10:49:41"
    expect := time.Date(2003, 9, 25, 10, 49, 41, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestISOFormatStrip3(t *testing.T) {
    parser := &Parser{}
    timestr := "2003-09-25T10:49"
    expect := time.Date(2003, 9, 25, 10, 49, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestISOFormatStrip4(t *testing.T) {
    parser := &Parser{}
    timestr := "2003-09-25T10"
    expect := time.Date(2003, 9, 25, 10, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestISOFormatStrip5(t *testing.T) {
    parser := &Parser{}
    timestr := "2003-09-25"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestISOStrippedFormat(t *testing.T) {
    parser := &Parser{}
    timestr := "20030925T104941.5-0300"
    expect := time.Date(2003, 9, 25, 10, 49, 41, 500000000, BRSTLoc)
    check(t, parser, timestr, expect)
}

func TestISOStrippedFormatStrip1(t *testing.T) {
    parser := &Parser{}
    timestr := "20030925T104941-0300"
    expect := time.Date(2003, 9, 25, 10, 49, 41, 0, BRSTLoc)
    check(t, parser, timestr, expect)
}

func TestISOStrippedFormatStrip2(t *testing.T) {
    parser := &Parser{}
    timestr := "20030925T104941"
    expect := time.Date(2003, 9, 25, 10, 49, 41, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestISOStrippedFormatStrip3(t *testing.T) {
    parser := &Parser{}
    timestr := "20030925T1049"
    expect := time.Date(2003, 9, 25, 10, 49, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestISOStrippedFormatStrip4(t *testing.T) {
    parser := &Parser{}
    timestr := "20030925T10"
    expect := time.Date(2003, 9, 25, 10, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestISOStrippedFormatStrip5(t *testing.T) {
    parser := &Parser{}
    timestr := "20030925"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestNoSeparator1(t *testing.T) {
    parser := &Parser{}
    timestr := "199709020908"
    expect := time.Date(1997, 9, 2, 9, 8, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestNoSeparator2(t *testing.T) {
    parser := &Parser{}
    timestr := "19970902090807"
    expect := time.Date(1997, 9, 2, 9, 8, 7, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash1(t *testing.T) {
    parser := &Parser{}
    timestr := "2003-09-25"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash2(t *testing.T) {
    parser := &Parser{}
    timestr := "2003-Sep-25"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash3(t *testing.T) {
    parser := &Parser{}
    timestr := "25-09-2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash4(t *testing.T) {
    parser := &Parser{}
    timestr := "25-Sep-2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash5(t *testing.T) {
    parser := &Parser{}
    timestr := "Sep-25-2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash6(t *testing.T) {
    parser := &Parser{}
    timestr := "09-25-2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash7(t *testing.T) {
    parser := &Parser{}
    timestr := "25-09-2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash8(t *testing.T) {
    parser := &Parser{DayFirst: true}
    timestr := "10-09-2003"
    expect := time.Date(2003, 9, 10, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash9(t *testing.T) {
    parser := &Parser{}
    timestr := "10-09-2003"
    expect := time.Date(2003, 10, 9, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash10(t *testing.T) {
    parser := &Parser{}
    timestr := "10-09-03"
    expect := time.Date(2003, 10, 9, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDash11(t *testing.T) {
    parser := &Parser{YearFirst: true}
    timestr := "10-09-03"
    expect := time.Date(2010, 9, 3, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot1(t *testing.T) {
    parser := &Parser{}
    timestr := "2003.09.25"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot2(t *testing.T) {
    parser := &Parser{}
    timestr := "2003.Sep.25"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot3(t *testing.T) {
    parser := &Parser{}
    timestr := "25.09.2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot4(t *testing.T) {
    parser := &Parser{}
    timestr := "25.Sep.2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot5(t *testing.T) {
    parser := &Parser{}
    timestr := "Sep.25.2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot6(t *testing.T) {
    parser := &Parser{}
    timestr := "09.25.2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot7(t *testing.T) {
    parser := &Parser{}
    timestr := "25.09.2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot8(t *testing.T) {
    parser := &Parser{DayFirst: true}
    timestr := "10.09.2003"
    expect := time.Date(2003, 9, 10, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot9(t *testing.T) {
    parser := &Parser{}
    timestr := "10.09.2003"
    expect := time.Date(2003, 10, 9, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot10(t *testing.T) {
    parser := &Parser{}
    timestr := "10.09.03"
    expect := time.Date(2003, 10, 9, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithDot11(t *testing.T) {
    parser := &Parser{YearFirst: true}
    timestr := "10.09.03"
    expect := time.Date(2010, 9, 3, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash1(t *testing.T) {
    parser := &Parser{}
    timestr := "2003/09/25"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash2(t *testing.T) {
    parser := &Parser{}
    timestr := "2003/Sep/25"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash3(t *testing.T) {
    parser := &Parser{}
    timestr := "25/09/2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash4(t *testing.T) {
    parser := &Parser{}
    timestr := "25/Sep/2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash5(t *testing.T) {
    parser := &Parser{}
    timestr := "Sep/25/2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash6(t *testing.T) {
    parser := &Parser{}
    timestr := "09/25/2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash7(t *testing.T) {
    parser := &Parser{}
    timestr := "25/09/2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash8(t *testing.T) {
    parser := &Parser{DayFirst: true}
    timestr := "10/09/2003"
    expect := time.Date(2003, 9, 10, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash9(t *testing.T) {
    parser := &Parser{}
    timestr := "10/09/2003"
    expect := time.Date(2003, 10, 9, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash10(t *testing.T) {
    parser := &Parser{}
    timestr := "10/09/03"
    expect := time.Date(2003, 10, 9, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSlash11(t *testing.T) {
    parser := &Parser{YearFirst: true}
    timestr := "10/09/03"
    expect := time.Date(2010, 9, 3, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace1(t *testing.T) {
    parser := &Parser{}
    timestr := "2003 09 25"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace2(t *testing.T) {
    parser := &Parser{}
    timestr := "2003 Sep 25"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace3(t *testing.T) {
    parser := &Parser{}
    timestr := "25 09 2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace4(t *testing.T) {
    parser := &Parser{}
    timestr := "25 Sep 2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace5(t *testing.T) {
    parser := &Parser{}
    timestr := "Sep 25 2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace6(t *testing.T) {
    parser := &Parser{}
    timestr := "09 25 2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace7(t *testing.T) {
    parser := &Parser{}
    timestr := "25 09 2003"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace8(t *testing.T) {
    parser := &Parser{DayFirst: true}
    timestr := "10 09 2003"
    expect := time.Date(2003, 9, 10, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace9(t *testing.T) {
    parser := &Parser{}
    timestr := "10 09 2003"
    expect := time.Date(2003, 10, 9, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace10(t *testing.T) {
    parser := &Parser{}
    timestr := "10 09 03"
    expect := time.Date(2003, 10, 9, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace11(t *testing.T) {
    parser := &Parser{YearFirst: true}
    timestr := "10 09 03"
    expect := time.Date(2010, 9, 3, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestDateWithSpace12(t *testing.T) {
    parser := &Parser{}
    timestr := "25 09 03"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestStrangelyOrderedDate1(t *testing.T) {
    parser := &Parser{}
    timestr := "03 25 Sep"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestStrangelyOrderedDate2(t *testing.T) {
    parser := &Parser{}
    timestr := "2003 25 Sep"
    expect := time.Date(2003, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestStrangelyOrderedDate3(t *testing.T) {
    parser := &Parser{}
    timestr := "25 03 Sep"
    expect := time.Date(2025, 9, 3, 0, 0, 0, 0, UTCLoc)   
    check(t, parser, timestr, expect)
}

func TestHourWithLetters(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10h36m28.5s"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 500000000, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourWithLettersStrip1(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10h36m28s"
    expect := time.Date(2003, 9, 25, 10, 36, 28, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourWithLettersStrip2(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10h36m"
    expect := time.Date(2003, 9, 25, 10, 36, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourWithLettersStrip3(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10h"
    expect := time.Date(2003, 9, 25, 10, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourWithLettersStrip4(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10 h 36"
    expect := time.Date(2003, 9, 25, 10, 36, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm1(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10h am"
    expect := time.Date(2003, 9, 25, 10, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm2(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10h pm"
    expect := time.Date(2003, 9, 25, 22, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm3(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10am"
    expect := time.Date(2003, 9, 25, 10, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm4(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10pm"
    expect := time.Date(2003, 9, 25, 22, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm5(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:00 am"
    expect := time.Date(2003, 9, 25, 10, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm6(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:00 pm"
    expect := time.Date(2003, 9, 25, 22, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm7(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:00am"
    expect := time.Date(2003, 9, 25, 10, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm8(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:00pm"
    expect := time.Date(2003, 9, 25, 22, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm9(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:00a.m"
    expect := time.Date(2003, 9, 25, 10, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm10(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:00p.m"
    expect := time.Date(2003, 9, 25, 22, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm11(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:00a.m."
    expect := time.Date(2003, 9, 25, 10, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHourAmPm12(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "10:00p.m."
    expect := time.Date(2003, 9, 25, 22, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestPertain1(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Sep 03"
    expect := time.Date(2003, 9, 3, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestPertain2(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Sep of 05"
    expect := time.Date(2005, 9, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestWeekdayAlone(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Wed"
    expect := time.Date(2003, 10, 1, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestLongWeekday(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "Wednesday"
    expect := time.Date(2003, 10, 1, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestLongMonth(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "October"
    expect := time.Date(2003, 10, 25, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestZeroYear(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "31-Dec-00"
    expect := time.Date(2000, 12, 31, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestFuzzy1(t *testing.T) {
    parser := &Parser{Fuzzy: true}
    timestr := "Today is 25 of September of 2003, exactly at 10:49:41 with timezone -03:00."
    expect := time.Date(2003, 9, 25, 10, 49, 41, 0, BRSTLoc)
    check(t, parser, timestr, expect)
}

func TestFuzzy2(t *testing.T) {
    s := "Today is 25 of September of 2003, exactly at 10:49:41 with timezone -03:00."
    _, err := (&Parser{}).Parse(s)
    if err == nil {
        t.Fatalf("Parse of fuzzy string without fuzzy mode enabled should fail, but did not.")
    }
}

func TestExtraSpace(t *testing.T) {
    parser := &Parser{}
    timestr := "  July   4 ,  1976   12:01:02   am  "
    expect := time.Date(1976, 7, 4, 0, 1, 2, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat1(t *testing.T) {
    parser := &Parser{}
    timestr := "Wed, July 10, '96"
    expect := time.Date(1996, 7, 10, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat2(t *testing.T) {
    parser := &Parser{IgnoreTZ: true}
    timestr := "1996.07.10 AD at 15:08:56 PDT"
    expect := time.Date(1996, 7, 10, 15, 8, 56, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat3(t *testing.T) {
    parser := &Parser{}
    timestr := "1996.July.10 AD 12:08 PM"
    expect := time.Date(1996, 7, 10, 12, 8, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat4(t *testing.T) {
    parser := &Parser{IgnoreTZ: true}
    timestr := "Tuesday, April 12, 1952 AD 3:30:42pm PST"
    expect := time.Date(1952, 4, 12, 15, 30, 42, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat5(t *testing.T) {
    parser := &Parser{IgnoreTZ: true}
    timestr := "November 5, 1994, 8:15:30 am EST"
    expect := time.Date(1994, 11, 5, 8, 15, 30, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat6(t *testing.T) {
    parser := &Parser{IgnoreTZ: true}
    timestr := "1994-11-05T08:15:30-05:00"
    expect := time.Date(1994, 11, 5, 8, 15, 30, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat7(t *testing.T) {
    parser := &Parser{IgnoreTZ: true}
    timestr := "1994-11-05T08:15:30Z"
    expect := time.Date(1994, 11, 5, 8, 15, 30, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat8(t *testing.T) {
    parser := &Parser{}
    timestr := "July 4, 1976"
    expect := time.Date(1976, 7, 4, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat9(t *testing.T) {
    parser := &Parser{}
    timestr := "7 4 1976"
    expect := time.Date(1976, 7, 4, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat10(t *testing.T) {
    parser := &Parser{}
    timestr := "4 jul 1976"
    expect := time.Date(1976, 7, 4, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat11(t *testing.T) {
    parser := &Parser{}
    timestr := "7-4-76"
    expect := time.Date(1976, 7, 4, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat12(t *testing.T) {
    parser := &Parser{}
    timestr := "19760704"
    expect := time.Date(1976, 7, 4, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat13(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "0:01:02"
    expect := time.Date(2003, 9, 25, 0, 1, 2, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat14(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "12h 01m02s am"
    expect := time.Date(2003, 9, 25, 0, 1, 2, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat15(t *testing.T) {
    parser := &Parser{}
    timestr := "0:01:02 on July 4, 1976"
    expect := time.Date(1976, 7, 4, 0, 1, 2, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat16(t *testing.T) {
    parser := &Parser{}
    timestr := "0:01:02 on July 4, 1976"
    expect := time.Date(1976, 7, 4, 0, 1, 2, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat17(t *testing.T) {
    parser := &Parser{IgnoreTZ: true}
    timestr := "1976-07-04T00:01:02Z"
    expect := time.Date(1976, 7, 4, 0, 1, 2, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat18(t *testing.T) {
    parser := &Parser{}
    timestr := "July 4, 1976 12:01:02 am"
    expect := time.Date(1976, 7, 4, 0, 1, 2, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat19(t *testing.T) {
    parser := &Parser{}
    timestr := "Mon Jan  2 04:24:27 1995"
    expect := time.Date(1995, 1, 2, 4, 24, 27, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat20(t *testing.T) {
    parser := &Parser{IgnoreTZ: true}
    timestr := "Tue Apr 4 00:22:12 PDT 1995"
    expect := time.Date(1995, 4, 4, 0, 22, 12, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat21(t *testing.T) {
    parser := &Parser{}
    timestr := "04.04.95 00:22"
    expect := time.Date(1995, 4, 4, 0, 22, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat22(t *testing.T) {
    parser := &Parser{}
    timestr := "Jan 1 1999 11:23:34.578"
    expect := time.Date(1999, 1, 1, 11, 23, 34, 578000000, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat23(t *testing.T) {
    parser := &Parser{}
    timestr := "950404 122212"
    expect := time.Date(1995, 4, 4, 12, 22, 12, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat24(t *testing.T) {
    parser := &Parser{Default: TestDefault, IgnoreTZ: true}
    timestr := "0:00 PM, PST"
    expect := time.Date(2003, 9, 25, 12, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat25(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "12:08 PM"
    expect := time.Date(2003, 9, 25, 12, 8, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat26(t *testing.T) {
    parser := &Parser{}
    timestr := "5:50 AM on June 13, 1990"
    expect := time.Date(1990, 6, 13, 5, 50, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat27(t *testing.T) {
    parser := &Parser{}
    timestr := "3rd of May 2001"
    expect := time.Date(2001, 5, 3, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat28(t *testing.T) {
    parser := &Parser{}
    timestr := "5th of March 2001" 
    expect := time.Date(2001, 3, 5, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat29(t *testing.T) {
    parser := &Parser{}
    timestr := "1st of May 2003"
    expect := time.Date(2003, 5, 1, 0, 0, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat30(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "01h02m03"
    expect := time.Date(2003, 9, 25, 1, 2, 3, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TsetRandomFormat31(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "01h02"
    expect := time.Date(2003, 9, 25, 1, 2, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat32(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "01h02s"
    expect := time.Date(2003, 9, 25, 1, 0, 2, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat33(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "01m02"
    expect := time.Date(2003, 9, 25, 0, 1, 2, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat34(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "01m02h"
    expect := time.Date(2003, 9, 25, 2, 1, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestRandomFormat35(t *testing.T) {
    parser := &Parser{Default: TestDefault}
    timestr := "2004 10 Apr 11h30m"
    expect := time.Date(2004, 4, 10, 11, 30, 0, 0, UTCLoc)
    check(t, parser, timestr, expect)
}

func TestHighPrecisionSeconds(t *testing.T) {
    parser := &Parser{}
    timestr := "20080227T21:26:01.123456789"
    expect := time.Date(2008, 2, 27, 21, 26, 1, 123456000, UTCLoc) // Parser only stores microseconds, not nanoseconds
    check(t, parser, timestr, expect)
}

func testIncreasingInternal(t *testing.T, format string) {
    date := time.Date(1900, 1, 1, 0, 0, 0, 0, UTCLoc)
    
    // delta = approx. 1s + 1m + 1h + 1d + 1m + 1y
    delta := time.Second + time.Minute + time.Hour +
        (24 * time.Hour) +              // Day
        (31 * 24 * time.Hour) +         // Month
        (365 * 24 * time.Hour)          // Year
    
    parser := &Parser{}
        
    for i := 0; i < 200; i++ {
        check(t, parser, date.Format(format), date)
        date = date.Add(delta)
    }
}

func TestIncreasingANSICFormat(t *testing.T) {
    testIncreasingInternal(t, time.ANSIC)
}

// Very similar to ISO format
func TestIncreasingRFC3339Format(t *testing.T) {
    testIncreasingInternal(t, time.RFC3339)
}

func testRandomInternal(t *testing.T, format string) {
    parser := &Parser{}
    
    for i := 0; i < 200; i++ {
        date := time.Date(
            1900 + rand.Intn(200),
            time.Month(1 + rand.Intn(12)),
            1 + rand.Intn(28), // Cap at 28 to avoid potential clashes with the max number of days in given month
            rand.Intn(24),
            rand.Intn(60),
            rand.Intn(60),
            0,
            UTCLoc,
        )
        
        check(t, parser, date.Format(format), date)
    }
}

func TestRandomANSICFormat(t *testing.T) {
    testRandomInternal(t, time.ANSIC)
}

// Very similar to ISO format
func TestRandomRFC3339Format(t *testing.T) {
    testRandomInternal(t, time.RFC3339)
}

func ExampleParser_Parse() {
    parser := &Parser{}
    t, err := parser.Parse("Thu, 25 Sep 2003 10:49:41")
    if err != nil {
        fmt.Printf("Parse error: %s\n", err.Error())
    }
    
    fmt.Printf("year: %d, month: %d, day: %d\n", t.Year(), t.Month(), t.Day())
    fmt.Printf("hour: %d, minute: %d, second: %d\n", t.Hour(), t.Minute(), t.Second())
    
    // Output:
    // year: 2003, month: 9, day: 25
    // hour: 10, minute: 49, second: 41
}
