package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func String(msg string) string {
	var s string
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "%s ", msg)
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}

	return strings.TrimSpace(s)
}

func Int(msg string) (int, error) {
	var i int64
	var err error
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "%s ", msg)
		s, _ := r.ReadString('\n')
		i, err = strconv.ParseInt(s, 10, 0)
		if err != nil {
			return 0, err
		}
		if s != "" {
			break
		}
	}

	return int(i), nil
}