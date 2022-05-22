package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"sykesdev.ca/gog/internal/common"
)

func String(msg string) string {
	var b []byte
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Fprintf(os.Stderr, "%s ", msg)
		b, _ = r.ReadBytes('\n')
		if len(b) > 0 {
			break
		}
	}

	return common.CleanStdoutSingleline(b)
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