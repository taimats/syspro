package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func LoadEnvFile(path string) error {
	envMap := make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if isComment(line) || isBlank(line) {
			continue
		}
		if hasSpace(line) {
			strs := strings.Fields(line)
			for _, s := range strs {
				pair := strings.Split(s, "=")
				envMap[pair[0]] = pair[1]
			}
			continue
		}
		strs := strings.Split(line, "=")
		envMap[strs[0]] = strs[1]
	}

	for k, v := range envMap {
		err := os.Setenv(k, v)
		if err != nil {
			return err
		}
	}
	fmt.Println("環境変数をセット")

	fmt.Println("設定した環境変数を取得...")
	for k := range envMap {
		env, ok := os.LookupEnv(k)
		if !ok {
			return errors.New("no such an environment")
		}
		fmt.Printf("{key:%s, value:%s}\n", k, env)
	}

	return nil
}

func isComment(s string) bool {
	return strings.Contains(s, "#") || strings.Contains(s, "//")
}

func isBlank(s string) bool {
	return len(s) == 0
}

func hasSpace(s string) bool {
	for _, r := range s {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}
