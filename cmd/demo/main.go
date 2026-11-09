package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ilyutsik/pre-trainee-practice/internal/accesslog"
	"github.com/ilyutsik/pre-trainee-practice/internal/booking"
	"github.com/ilyutsik/pre-trainee-practice/internal/fileproc"
	"github.com/ilyutsik/pre-trainee-practice/internal/pipeline"
	"github.com/ilyutsik/pre-trainee-practice/internal/report"
	"github.com/ilyutsik/pre-trainee-practice/internal/wordcount"
)

func main() {

}

func testFileProc() {
	fmt.Println(fileproc.Process([]string{"testdata/fileproc/ok.txt", "testdata/fileproc/missing.txt", "testdata/fileproc/panic.txt", "testdata/fileproc/empty.txt"}))
}

func testReport() {
	recs := []report.Record{{SKU: "A-1", Qty: 2, PriceCents: 199}, {SKU: "B-9", Qty: 1, PriceCents: 50}}
	exp1 := report.JSONExporter{}
	exp2 := report.CSVExporter{Comma: ':'}
	exporters := []report.Exporter{exp2, exp1}
	v, err := report.RenderAll(recs, exporters)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(v["json"]))
	fmt.Println(string(v["csv"]))

	exp1Comm, err1 := report.CommaOf(exp1)
	exp2Comm, err2 := report.CommaOf(exp2)
	fmt.Println(string(exp1Comm), err1)
	fmt.Println(string(exp2Comm), err2)

}

func testBooking() {
	id1 := "1"
	id2 := "2"
	id3 := "3"
	id4 := "4"
	s := booking.NewMemoryService()

	no1 := s.AddRoom(id1, "Kolomna")
	no2 := s.AddRoom(id2, "SportTime")
	no3 := s.AddRoom(id3, "Irkutsk")
	no4 := s.AddRoom(id4, "IrkutskiBaikal")
	if no1 != nil || no2 != nil || no3 != nil || no4 != nil {
		fmt.Println(no1, no2, no3)
	}
	errIdRequired := s.AddRoom("", "Irkutsk")
	if errIdRequired != nil {
		fmt.Println(errors.Is(errIdRequired, booking.ErrInvalid), errIdRequired)
	}
	errNameRequired := s.AddRoom("4", "")
	if errNameRequired != nil {
		fmt.Println(errors.Is(errNameRequired, booking.ErrInvalid), errNameRequired)
	}
	errIdExist := s.AddRoom(id3, "Irkutsk")
	if errIdExist != nil {
		fmt.Println(errors.Is(errIdExist, booking.ErrInvalid), errIdExist)
	}

	r1, no4 := s.Book(id1, 12, 17)
	r2, no5 := s.Book(id2, 17, 18)
	r3, no6 := s.Book(id3, 7, 12)
	if no4 != nil || no5 != nil || no6 != nil {
		fmt.Println(no1, no2, no3)
	}
	_, errRoomIdRequired := s.Book("", 7, 12)
	if errRoomIdRequired != nil {
		fmt.Println(errors.Is(errRoomIdRequired, booking.ErrInvalid), errRoomIdRequired)
	}
	_, errInvalidTimeInterval := s.Book(id2, -1, 8)
	if errInvalidTimeInterval != nil {
		fmt.Println(errors.Is(errInvalidTimeInterval, booking.ErrInvalid), errInvalidTimeInterval)
	}
	_, errorIdNotFound := s.Book("14", 1, 3)
	if errorIdNotFound != nil {
		fmt.Println(errors.Is(errorIdNotFound, booking.ErrNotFound), errorIdNotFound)
	}
	_, errorOverlap := s.Book(id3, 11, 21)
	if errorOverlap != nil {
		fmt.Println(errors.Is(errorOverlap, booking.ErrConflict), errorOverlap)
	}

	errCancelIdRequired := s.Cancel("")
	if errCancelIdRequired != nil {
		fmt.Println(errors.Is(errCancelIdRequired, booking.ErrInvalid), errCancelIdRequired)
	}
	errCancelIdNotFound := s.Cancel("14")
	if errCancelIdNotFound != nil {
		fmt.Println(errors.Is(errCancelIdNotFound, booking.ErrNotFound), errCancelIdNotFound)

	}
	no7 := s.Cancel(r1.ID)
	if no7 != nil {
		fmt.Println(no7)
	}
	errAlreadyCanceled := s.Cancel(r1.ID)
	if errAlreadyCanceled != nil {
		fmt.Println(errors.Is(errAlreadyCanceled, booking.ErrCancelled), errAlreadyCanceled)
	}

	v, no8 := s.Get(r2.ID)
	if no8 != nil {
		fmt.Println(v)
	}
	_, errGetReservationIdRequired := s.Get("")
	if errGetReservationIdRequired != nil {
		fmt.Println(errors.Is(errGetReservationIdRequired, booking.ErrInvalid), errGetReservationIdRequired)
	}
	_, errGetReservationIdNotFound := s.Get("21312312312312")
	if errGetReservationIdNotFound != nil {
		fmt.Println(errors.Is(errGetReservationIdNotFound, booking.ErrNotFound), errGetReservationIdNotFound)
	}

	_, errListRoomIdRequired := s.ListRoom("")
	if errListRoomIdRequired != nil {
		fmt.Println(errors.Is(errListRoomIdRequired, booking.ErrInvalid), errListRoomIdRequired)
	}
	_, errListRoomIdNotFound := s.ListRoom("1231231231231")
	if errListRoomIdNotFound != nil {
		fmt.Println(errors.Is(errListRoomIdNotFound, booking.ErrNotFound), errListRoomIdNotFound)
	}
	_, errNoReservationsFound := s.ListRoom(id4)
	if errNoReservationsFound != nil {
		fmt.Println(errors.Is(errNoReservationsFound, booking.ErrNotFound), errNoReservationsFound)
	}
	fmt.Println("----------------------------------------------------------")
	fmt.Println(s.ListRoom(id1))
	fmt.Println(s.ListRoom(id2))
	fmt.Println(s.ListRoom(id3))
	fmt.Println(s.Book(id3, 11, 13))
	fmt.Println(s.Cancel(r3.ID))
	fmt.Println(s.ListRoom(id3))
	fmt.Println(s.Book(id3, 11, 13))
	fmt.Println(s.ListRoom(id4))
	fmt.Println(r1)
	fmt.Println(r2)
	fmt.Println(r3)

}

func testAccessLog() {
	text, err := os.ReadFile("testdata/access.log")
	if err != nil {
		fmt.Println(err.Error())
	}
	s := accesslog.Summarize(string(text), 3)
	fmt.Println(s)
}

func testWorldCount() {
	stopWordsString, err := os.ReadFile("testdata/stopList.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	stop := wordcount.StopWords(string(stopWordsString))

	str, err := os.ReadFile("testdata/words.txt")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tokens := wordcount.Tokens(string(str))
	tokensWithoutStopWords := wordcount.DropStop(tokens, stop)
	countTokens := wordcount.Count(tokensWithoutStopWords)
	topNTokens := wordcount.TopN(countTokens, 3)

	for _, v := range topNTokens {
		fmt.Printf("%s %d \n", v.Word, v.Count)
	}
}

func testPipeline() {
	primarySlice := []int{12, 34, 78, 1, 3, 9, 0, -1, -75, 54, 12, 0, -1, 8, 9, -1}
	//var primarySlice []int

	fmt.Printf("primarySlice before %v \n", primarySlice)

	doubleSlice := pipeline.Double(primarySlice)
	filterSlice := pipeline.FilterMin(primarySlice, 12)
	uniqueSlice := pipeline.UniqueStable(primarySlice)

	fmt.Printf("primarySlice after %v \n", primarySlice)
	fmt.Printf("doubleSlice %v \n", doubleSlice)
	fmt.Printf("filterSlice %v \n", filterSlice)
	fmt.Printf("uniqueSlice %v \n", uniqueSlice)
}
