package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %s", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	contentLines := bufio.NewReader(r)
	if err != nil {
		return
	}

	var user User
	var line []byte
	var i int
	for {
		line, _, err = contentLines.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
		if err = jsoniter.Unmarshal(line, &user); err != nil {
			return
		}
		result[i] = user
		i++
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	re := regexp.MustCompile("\\." + domain)
	user := User{}

	for _, user = range u {
		if re.MatchString(user.Email) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
