package component

import (
	"strconv"
	"time"

	"github.com/a-h/templ"

	lw "github.com/Design-Machines-Studio/livewires-templ"
)

// calendarDay is one cell in a rendered month grid.
type calendarDay struct {
	Date    time.Time
	Class   string
	Outside bool
	Today   bool
	Pressed bool
}

// calendarWeekday is a column heading in a month grid.
type calendarWeekday struct {
	Short string // "Mo"
	Long  string // "Monday"
}

// dayKey returns a comparable YYYYMMDD value, or 0 for the zero time.
func dayKey(t time.Time) int {
	if t.IsZero() {
		return 0
	}
	y, m, d := t.Date()
	return y*10000 + int(m)*100 + d
}

// calendarWeekStart returns the first weekday column.
func calendarWeekStart(props CalendarProps) time.Weekday {
	if props.SundayFirst {
		return time.Sunday
	}
	return time.Monday
}

// calendarWeekdays returns the seven column headings starting at the week start.
func calendarWeekdays(props CalendarProps) []calendarWeekday {
	start := calendarWeekStart(props)
	days := make([]calendarWeekday, 7)
	for i := range days {
		wd := time.Weekday((int(start) + i) % 7)
		days[i] = calendarWeekday{Short: wd.String()[:2], Long: wd.String()}
	}
	return days
}

// calendarWeeks returns six weeks of cells covering props.Month.
func calendarWeeks(props CalendarProps) [][]calendarDay {
	y, m, _ := props.Month.Date()
	first := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	offset := (int(first.Weekday()) - int(calendarWeekStart(props)) + 7) % 7
	cursor := first.AddDate(0, 0, -offset)

	today := dayKey(props.Today)
	selected := dayKey(props.Selected)
	rangeStart := dayKey(props.RangeStart)
	rangeEnd := dayKey(props.RangeEnd)
	marked := make(map[int]bool, len(props.Marked))
	for _, t := range props.Marked {
		marked[dayKey(t)] = true
	}

	weeks := make([][]calendarDay, 6)
	for w := range weeks {
		weeks[w] = make([]calendarDay, 7)
		for d := range weeks[w] {
			key := dayKey(cursor)
			day := calendarDay{
				Date:    cursor,
				Outside: cursor.Month() != m,
				Today:   key == today,
			}
			isStart := key == rangeStart
			isEnd := key == rangeEnd
			inRange := rangeStart != 0 && rangeEnd != 0 && key > rangeStart && key < rangeEnd
			day.Pressed = key == selected || isStart || isEnd
			day.Class = lw.ClassNames(
				lw.If(day.Outside, "outside"),
				lw.If(day.Today, "today"),
				lw.If(key == selected, "selected"),
				lw.If(isStart, "range-start"),
				lw.If(inRange, "in-range"),
				lw.If(isEnd, "range-end"),
				lw.If(marked[key], "marked"),
			)
			weeks[w][d] = day
			cursor = cursor.AddDate(0, 0, 1)
		}
	}
	return weeks
}

// calendarTitleID returns the heading id that labels the grid.
func calendarTitleID(props CalendarProps) string {
	return lw.DefaultStr(props.ID, "cal-"+props.Month.Format("2006-01"))
}

// calendarTitle returns the visible month heading.
func calendarTitle(props CalendarProps) string {
	return lw.DefaultStr(props.Title, props.Month.Format("January 2006"))
}

// calendarDayAttrs returns consumer attributes for a day button.
func calendarDayAttrs(props CalendarProps, date time.Time) templ.Attributes {
	if props.DayAttrs == nil {
		return nil
	}
	return props.DayAttrs(date)
}

// calendarGroupMonth returns props for month i of n in a CalendarGroup.
func calendarGroupMonth(props CalendarProps, i, n int) CalendarProps {
	p := props
	y, m, _ := props.Month.Date()
	p.Month = time.Date(y, m+time.Month(i), 1, 0, 0, 0, 0, time.UTC)
	p.HidePrev = props.HidePrev || i > 0
	p.HideNext = props.HideNext || i < n-1
	if i > 0 {
		p.Title = ""
		if props.ID != "" {
			p.ID = props.ID + "-" + strconv.Itoa(i+1)
		}
	}
	return p
}
