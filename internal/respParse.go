package internal

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

func ReadRespCommand(r *bufio.Reader) ([]string, error) {

	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSpace(line)

	if line[0] != '*' {
		return nil, fmt.Errorf("expected array")
	}

	n, _ := strconv.Atoi(line[1:])

	args := make([]string, 0, n)

	for i := 0; i < n; i++ {

		_, err := r.ReadString('\n') // $len
		if err != nil {
			return nil, err
		}

		value, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		args = append(args, strings.TrimSpace(value))
	}

	return args, nil
}
