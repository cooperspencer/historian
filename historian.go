package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/buddyspencer/chameleon"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"
)

var (
	db = &sql.DB{}
	historian_path = ""
	hdb = ""
	home = ""

	// KINGPIN
	app = kingpin.New("historian", "I store your history and search it for you")
	save = app.Command("save", "Save a command")
	search = app.Command("search", "search for a command")
	svals = search.Arg("criteria", "criteria you want to search for").Strings()
	integrate = app.Command("integrate", "integrate the old historian version into your database")
)

func w2db(command *bytes.Buffer, timestamp time.Time){
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS history (id INTEGER PRIMARY KEY AUTOINCREMENT, timestamp INTEGER, command VARCHAR)")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	statement.Exec()
	statement, _ = db.Prepare("INSERT INTO history(timestamp, command) VALUES (?, ?)")
	statement.Exec(int32(timestamp.Unix()), command.String())
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func igrate() {
	f := "2006-01-02.15:04:05"
	file, err := os.Open(fmt.Sprintf("%s/.logs/bash-history.log", home))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		r := regexp.MustCompile(`\[(\d{4})-(\d{2})-(\d{2}).(\d{2}):(\d{2}):(\d{2})\]`)
		matches := r.FindAllString(scanner.Text(), -1)
		if len(matches) > 0 {
			ot := strings.Split(strings.Split(scanner.Text(), "[")[1], "]")[0]
			command := strings.Split(scanner.Text(), "] ")[1]
			t, err := time.Parse(f, ot)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			bbuffer := bytes.NewBufferString(command + "\n")
			w2db(bbuffer, t)
		}
	}
	file.Close()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	os.Remove(fmt.Sprintf("%s/.logs/bash-history.log", home))
}

func sidb() {
	squery := "SELECT * FROM history WHERE "
	lquery := ""

	for i, v := range *svals {
		lquery += fmt.Sprintf("command LIKE '%%%s%%'", v)
		if i < len(*svals) - 1 {
			lquery += " and "
		}
	}

	squery += lquery

	rows, _ := db.Query(squery)

	var id int
	var timestamp int
	var command string

	for rows.Next() {
		rows.Scan(&id, &timestamp, &command)
		t := time.Unix(int64(timestamp), 0)
		for _, x := range *svals {
			r := regexp.MustCompile("(?i)" + x)
			matches := r.FindAllString(command, -1)
			umatches := []string{}
			for _, y := range matches {
				if !contains(umatches, y) {
					umatches = append(umatches, y)
					command = strings.Replace(command, y, fmt.Sprint(chameleon.Green(y)), -1)
				}
			}
		}
		fmt.Printf("[%s] %s", chameleon.Lightblue(t.Format("2006.01.02:15:04:05")), command)
	}
}

func getAll() {
	squery := "SELECT * FROM history"

	rows, _ := db.Query(squery)

	var id int
	var timestamp int
	var command string

	for rows.Next() {
		rows.Scan(&id, &timestamp, &command)
		t := time.Unix(int64(timestamp), 0)
		fmt.Printf("[%s] %s", chameleon.Lightblue(t.Format("2006.01.02:15:04:05")), command)
	}
}

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	home = usr.HomeDir
	historian_path = fmt.Sprintf("%s/.historian", home)

	if _, err := os.Stat(historian_path); os.IsNotExist(err) {
		os.Mkdir(historian_path, 0775)
	}

	hdb = fmt.Sprintf("%s/historian.db", historian_path)

	db, err = sql.Open("sqlite3", hdb)
	defer db.Close()

	if len(os.Args[1:]) == 0 {
		getAll()
		os.Exit(0)
	}

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case save.FullCommand():
		buf := &bytes.Buffer{}
		n, err := io.Copy(buf, os.Stdin)
		if err != nil {
			log.Fatalln(err)
		} else if n <= 1 { // buffer always contains '\n'
			log.Fatalln("no input provided")
		}
		t := time.Now()
		w2db(buf, t)
	case search.FullCommand():
		sidb()
	case integrate.FullCommand():
		igrate()
	}

}