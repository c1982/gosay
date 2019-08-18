package main

import "testing"

func TestGetLines(t *testing.T) {

	stackLog := `goroutine 1 [running]:
	main.readgoroutines()
			/Users/oguzhan/go/src/gosay/main.go:44 +0x6d
	main.main()
			/Users/oguzhan/go/src/gosay/main.go:17 +0x55`

	t.Run("get lines from stacktrace", func(t *testing.T) {
		lines, err := getLines(stackLog)

		if err != nil {
			t.Error(err)
		}

		if len(lines) != 5 {
			t.Errorf("got %d want %d", len(lines), 5)
		}
	})

	t.Run("parse goroutine line", func(t *testing.T) {

		id, status, err := parseGoRoutineLine(stackLog)
		if err != nil {
			t.Error(err)
		}

		if id != "1" {
			t.Errorf("got %s want %s", id, "1")
		}

		if status != "running" {
			t.Errorf("got %s want running", status)
		}

		t.Log("ID:", id, "Status:", status)
	})

}
