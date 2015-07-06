package schedule

import (
	"testing"
	"time"
)

func XTestBuildCommitScheduleFullWeeks(t *testing.T) {
	var tests = []struct {
		startDay time.Time
		numDays  int
	}{
		{time.Date(2009, time.November, 9, 0, 0, 0, 0, time.UTC), 7},
		{time.Date(2015, time.November, 9, 0, 0, 0, 0, time.UTC), 21},
		{time.Date(2014, time.July, 7, 0, 0, 0, 0, time.UTC), 365},
	}
	for _, test := range tests {
		days := getTestDays(test.startDay, test.numDays)
		schedule := BuildCommitSchedule(days)
		for r, row := range schedule {
			for c, col := range row {
				if col != EMPTY {
					fmt := "Expected only EMPTY values, but got %d at (%d,%d)"
					t.Errorf(fmt, col, r, c)
				}
			}
		}
	}
}

func XTestBuildCommitScheduleWednesdayStart(t *testing.T) {
	var tests = []struct {
		startDay time.Time
		numDays  int
	}{
		{time.Date(2009, time.November, 11, 0, 0, 0, 0, time.UTC), 9},
		{time.Date(2015, time.November, 11, 0, 0, 0, 0, time.UTC), 23},
		{time.Date(2014, time.July, 9, 0, 0, 0, 0, time.UTC), 367},
	}
	for _, test := range tests {
		days := getTestDays(test.startDay, test.numDays)
		schedule := BuildCommitSchedule(days)
		for r, row := range schedule {
			for c, col := range row {
				firstWeekMondayOrTuesDay := (c == 0 && r < 2)
				if firstWeekMondayOrTuesDay {
					if col != NOT_A_FIELD {
						fmt := "Expected NOT_A_FIELD values, but got %d at (%d,%d)"
						t.Errorf(fmt, col, r, c)
					}
				} else if col != EMPTY {
					fmt := "Expected only EMPTY values, but got %d at (%d,%d)"
					t.Errorf(fmt, col, r, c)
				}
			}
		}
	}
}

func XTestBuildCommitScheduleThrusdayEnd(t *testing.T) {
	var tests = []struct {
		startDay time.Time
		numDays  int
	}{
		{time.Date(2009, time.November, 9, 0, 0, 0, 0, time.UTC), 4},
		{time.Date(2015, time.November, 9, 0, 0, 0, 0, time.UTC), 19},
		{time.Date(2014, time.July, 7, 0, 0, 0, 0, time.UTC), 362},
	}
	for _, test := range tests {
		days := getTestDays(test.startDay, test.numDays)
		schedule := BuildCommitSchedule(days)
		for r, row := range schedule {
			for c, col := range row {
				lastWeekFridayOrSaturdayOrSunday := (c == len(row) && r > 4)
				if lastWeekFridayOrSaturdayOrSunday {
					if col != NOT_A_FIELD {
						fmt := "Expected NOT_A_FIELD values, but got %d at (%d,%d)"
						t.Errorf(fmt, col, r, c)
					}
				} else if col != EMPTY {
					fmt := "Expected only EMPTY values, but got %d at (%d,%d)"
					t.Errorf(fmt, col, r, c)
				}
			}
		}
	}
}

func getTestDays(startDay time.Time, numDays int) []time.Time {
	var resultingDays []time.Time
	for i := 0; i < numDays; i++ {
		resultingDays = append(resultingDays, startDay.AddDate(0, 0, i))
	}
	return resultingDays
}

func TestBuildFirstWeek(t *testing.T) {
	var tests = []struct {
		day              time.Weekday
		expectedSchedule []ScheduleEntries
	}{
		{time.Sunday, []ScheduleEntries{0, 0, 0, 0, 0, 0, 0}},
		{time.Monday, []ScheduleEntries{-1, 0, 0, 0, 0, 0, 0}},
		{time.Tuesday, []ScheduleEntries{-1, -1, 0, 0, 0, 0, 0}},
		{time.Wednesday, []ScheduleEntries{-1, -1, -1, 0, 0, 0, 0}},
		{time.Thursday, []ScheduleEntries{-1, -1, -1, -1, 0, 0, 0}},
		{time.Friday, []ScheduleEntries{-1, -1, -1, -1, -1, 0, 0}},
		{time.Saturday, []ScheduleEntries{-1, -1, -1, -1, -1, -1, 0}},
	}
	for _, test := range tests {
		actualSchedule := buildFirstWeek(test.day)
		if !sliceEqual(actualSchedule, test.expectedSchedule) {
			fmt := "Expected %v as schedule, but got %v"
			t.Errorf(fmt, test.expectedSchedule, actualSchedule)
		}
	}
}

func TestBuildLastWeek(t *testing.T) {
	var tests = []struct {
		day              time.Weekday
		expectedSchedule []ScheduleEntries
	}{
		{time.Sunday, []ScheduleEntries{0, -1, -1, -1, -1, -1, -1}},
		{time.Monday, []ScheduleEntries{0, 0, -1, -1, -1, -1, -1}},
		{time.Tuesday, []ScheduleEntries{0, 0, 0, -1, -1, -1, -1}},
		{time.Wednesday, []ScheduleEntries{0, 0, 0, 0, -1, -1, -1}},
		{time.Thursday, []ScheduleEntries{0, 0, 0, 0, 0, -1, -1}},
		{time.Friday, []ScheduleEntries{0, 0, 0, 0, 0, 0, -1}},
		{time.Saturday, []ScheduleEntries{0, 0, 0, 0, 0, 0, 0}},
	}
	for _, test := range tests {
		actualSchedule := buildLastWeek(test.day)
		if !sliceEqual(actualSchedule, test.expectedSchedule) {
			fmt := "Expected %v as schedule, but got %v"
			t.Errorf(fmt, test.expectedSchedule, actualSchedule)
		}
	}
}

func sliceEqual(sliceA, sliceB []ScheduleEntries) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}
	for i := 0; i < len(sliceA); i++ {
		if sliceA[i] != sliceB[i] {
			return false
		}
	}
	return true
}

func TestConnectWeeksToSchedule(t *testing.T) {
	var tests = []struct {
		firstDay   time.Weekday
		lastDay    time.Weekday
		numEntries int
	}{
		{time.Sunday, time.Saturday, 371},
	}
	for _, test := range tests {
		firstWeek := buildFirstWeek(test.firstDay)
		lastWeek := buildLastWeek(test.lastDay)
		schedule := connectWeeksToSchedule(firstWeek, lastWeek)
		length := 0
		for _, row := range schedule {
			for _, entry := range row {
				if entry != EMPTY || entry != NOT_A_FIELD {
					t.Errorf("Entry should be EMPTY or NOT_A_FIELD, but was %v", entry)
				}
				length++
			}
		}
		if length != test.numEntries {
			t.Errorf("Expected length was %d, but got %d", test.numEntries, length)
		}
	}
}