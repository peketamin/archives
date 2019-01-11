package main

import (
	"strconv"
	"testing"
	"time"
)

func dbtest() {
	page, err := pickupWasDisplayedMarkPage()
	if err != nil {
		page, err = pickupUndisplayedPageRandom()
		if err != nil {
			resetWasDisplayedMarks()
			page, _ = pickupUndisplayedPageRandom()
		}
		page.updateWasDisplayedMarkAndFirstDisplayedTime()
		//log.Fatal("error")
	}
}

func dbtest_insert_data() {
	for i := 1; i < 6; i++ {
		page := Page{
			Title:  "t" + strconv.Itoa(i),
			Body:   "b" + strconv.Itoa(i),
			Note:   "n" + strconv.Itoa(i),
			Source: "s" + strconv.Itoa(i),
			//CurrentDisplayMark: false,
			WasDisplayedMark: false,
			FirstDisplayedAt: time.Now().AddDate(0, 0, -15),

			//Image: Binary{},
		}
		if i == 3 {
			page.WasDisplayedMark = true
			page.FirstDisplayedAt = time.Now()
		}
		db.Create(&page)
		//fmt.Printf("page: %v\n", page)
	}
}

func Test_pickupWasDisplayedMarkPage(t *testing.T) {
	dbconnect("db_test.dat")
	truncateTables()
	dbtest_insert_data()

	page, err := pickupWasDisplayedMarkPage()
	if err != nil {
		t.Error(err)
		t.Errorf("\nFailed: pickupWasDisplayedMarkPage().\n")
	}

	in := page.WasDisplayedMark
	const want = true
	if in != want {
		t.Errorf("page.WasDisplayedMark: %v, want: %v\n", in, want)
	}
}

func Test_resetWasDisplayedMarks(t *testing.T) {
	dbconnect("db_test.dat")
	truncateTables()
	dbtest_insert_data()

	err := resetWasDisplayedMarks()
	if err != nil {
		t.Errorf("Failed: resetWasDisplayedMarks(). Error: %v\n", err)
	}

	var count int
	db.Model(Page{}).Where("was_displayed_mark = ?", 1).Count(&count)

	in := count
	const want = 0
	if in != want {
		t.Errorf("Count rows where WasDisplayedMark is true after resetWasDisplayedMarks(): %v, want: %v\n", in, want)
	}
}
