package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"
)

/*

	groutine status: "idle","runnable","running","syscall","waiting","dead","copystack",

*/

type StackItem struct {
	ID     string
	Status string
	Nodes  StackNode
}

type StackNode struct {
	Name     string
	FileName string
	Parent   *StackNode
}

func Parse1(stacktrace string) error {

	lines, err := getLines(stacktrace)
	if err != nil {
		return err
	}

	for _, line := range lines {

		isroot := strings.HasSuffix(line, "goroutine")
		if isroot {
			id, status, err := parseGoRoutineLine(line)
			if err != nil {
				return err
			}

			root := StackItem{}
			root.ID = id
			root.Status = status
		}
	}

	return nil
}

func parseGoRoutineLine(goroutineline string) (id string, status string, err error) {

	r, _ := regexp.Compile(`^goroutine\W+(\d+)\W+\[(.+)\]\:$`)

	groups := r.FindStringSubmatch(goroutineline)
	if len(groups) != 3 {
		return id, status, fmt.Errorf("cannot parsed goroutine line. submatch %d", len(groups))
	}

	id = string(groups[1])
	status = string(groups[2])

	return id, status, nil
}

func getLines(stacktrace string) (lines []string, err error) {

	scanner := bufio.NewScanner(strings.NewReader(stacktrace))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	}

	return lines, nil
}

func main() {

	for i := 0; i < 4; i++ {
		go gogogo()
	}

	readstack()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	<-c

}

func gogogo() {
	time.Sleep(time.Second * 2)
	fmt.Println("done")
}

func readstack() {
	buf := make([]byte, 1<<20)
	stacklen := runtime.Stack(buf, true)
	fmt.Printf("\n%s", buf[:stacklen])
}
