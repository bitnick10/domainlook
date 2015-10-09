package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"time"
	"math"
)

func main() {
	db, err := sql.Open("sqlite3", "./domain.db")
	check(err)
	defer db.Close()

	//IsDomainRegistered("zhegeyinggaimeiyourenzhucesdf34.com")
	//return

	domains := GenerateDomainName(2, "")

	for _, domain := range domains {
		if !IsDomainNameInDatabase(db, domain) {
			isRegistered := 0
			if IsDomainRegistered(domain) {
				isRegistered = 1
			}
			insertString := fmt.Sprintf("insert into domain(domain_name,is_registered,last_update_time) values('%s',%d,%d)", domain, isRegistered, time.Now().Unix())
			fmt.Println(insertString)
			_, err = db.Exec(insertString)

			if err != nil {
				fmt.Println("db.Exec err")
				fmt.Println(err)
			}
		}
	}

	fmt.Println("main end")
}

func GenerateDomainName(nameLen int, suffix string) []string {
	size := int(math.Pow(36, float64(nameLen)))
	domains := make([]string, size)
	for i := 0; i < len(domains); i++ {
		name := IntToBase36String(i)
		name = Fill0Before(name, nameLen)
		name += suffix + ".com"
		domains[i] = name
	}
	return domains
}

func Fill0Before(str string, totalLen int) string {
	for i := 0; i < totalLen - len(str); i++ {
		str = "0" + str
	}
	return str;
}

func IntToBase36String(i int) string {
	ret := ""
	base36 := "0123456789abcdefghijklmnopqrstuvwxyz"
	for {
		ret += string(base36[i % 36])
		i /= 36
		if i == 0 {
			break;
		}
	}
	return Reverse(ret)
}

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes) - 1; i < j; i, j = i + 1, j - 1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func IsDomainNameInDatabase(db*sql.DB, domainName string) bool {
	rows, err := db.Query(fmt.Sprintf("select domain_name from domain where domain_name = '%s'", domainName))

	check(err)
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}
func check(err error) {
	if err != nil {
		fmt.Println("db.Exec err")
		fmt.Println(err)
	}
}
func IsDomainRegistered(name string) bool {
	resp, err := http.Get("http://whois.oray.com/whois/" + name)
	if err != nil {
		fmt.Println("http.Get err")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err")
	}
	//fmt.Println(string(body))
	if strings.Contains(string(body), "No match for") {
		return false
	}else {
		return true
	}
}