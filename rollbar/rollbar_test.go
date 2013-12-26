package rollbar_test

import (
	. "."
	"bufio"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func readJsonFile(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		l := strings.TrimSpace(s.Text())
		if l == "" || strings.HasPrefix(l, "//") {
			continue
		}
		lines = append(lines, s.Text())
	}
	if s.Err() != nil {
		panic(s.Err())
	}

	return []byte(strings.Join(lines, "\n"))
}

func TestUnmarshalItemEvent(t *testing.T) {
	b := readJsonFile("new_item.json")

	var event Event
	err := json.Unmarshal(b, &event)
	if err != nil {
		t.Fatal(err)
	}

	var actual ItemEvent
	err = json.Unmarshal(event.Data["item"], &actual)
	if err != nil {
		t.Fatal(err)
	}

	expected := ItemEvent{
		Environment: "production",
		Title:       "testing aobg98wrwe",
		LastOccurrence: LastOccurrence{
			Level:     "error",
			Timestamp: 1382655421,
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("\nActual:   %+v\nExpected: %+v", actual, expected)
	}
}
