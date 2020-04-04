package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gookit/color"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	db             = &sql.DB{}
	historian_path = ""
	hdb            = ""
	home           = ""

	// KINGPIN
	app    = kingpin.New("historian", "I store your history and search it for you")
	save   = app.Command("save", "Save a command")
	search = app.Command("search", "search for a command")
	svals  = search.Arg("criteria", "criteria you want to search for").Strings()
	conf = &Config{}
	searchcolor = color.Render
	datecolor = color.Render
)

type Config struct {
	DateColor 	string
	SearchColor string
	Secret      bool
	Dateformat  string
}

func ReadConfigfile(configfile string) *Config {
	cfgdata, err := ioutil.ReadFile(configfile)

	if err != nil {
		conf := &Config{"lightblue", "lightgreen", true, "2006.01.02:15:04:05"}
		confpath := filepath.Dir(configfile)
		if _, err := os.Stat(confpath); os.IsNotExist(err) {
			err := os.MkdirAll(confpath, 0700)
			if err != nil {
				panic(err.Error())
			}
			confbytes, err := yaml.Marshal(conf)
			if err != nil {
				panic(err.Error())
			}
			err = ioutil.WriteFile(configfile, confbytes, 0600)
			if err != nil {
				panic(err.Error())
			}
		}
		return conf
	}

	t := Config{}

	err = yaml.Unmarshal([]byte(cfgdata), &t)

	if err != nil {
		return &Config{"lightblue", "lightgreen", true, "2006.01.02:15:04:05"}
	}

	return &t
}

func getColor(colorString string) func(a ...interface{}) string {
	c := color.Color(0)

	if strings.Contains(colorString, "blue") {
		c = color.Blue
	} else if strings.Contains(colorString, "green") {
		c = color.Green
	} else if strings.Contains(colorString, "red") {
		c = color.Red
	} else if strings.Contains(colorString, "cyan") {
		c = color.Cyan
	} else if strings.Contains(colorString, "magenta") {
		c = color.Magenta
	} else if strings.Contains(colorString, "yellow") {
		c = color.Yellow
	} else {
		c = color.White
	}

	if strings.HasPrefix(colorString, "light") && runtime.GOOS != "windows" {
		c = c.Light()
	}

	return c.Render
}

func w2db(command *bytes.Buffer, timestamp time.Time) {
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
		r := regexp.MustCompile(`\[(\d{4})-(\d{2})-(\d{2})Did .(\d{2}):(\d{2}):(\d{2})\]`)
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
	if len(*svals) == 0 {
		fmt.Println("Please add search parameters")
		os.Exit(1)
	}

	squery := "SELECT * FROM history WHERE "
	lquery := ""

	for i, v := range *svals {
		lquery += fmt.Sprintf("command LIKE '%%%s%%'", v)
		if i < len(*svals)-1 {
			lquery += " and "
		}
	}

	squery += lquery

	rows, err := db.Query(squery)

	if err != nil {
		help()
	}

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
					command = strings.Replace(command, y, searchcolor(y), -1)
				}
			}
		}
		fmt.Printf("[%s] %s", datecolor(t.Format("2006.01.02:15:04:05")), command)
	}
}

func getAll() {
	squery := "SELECT * FROM history"

	rows, err := db.Query(squery)

	if err != nil {
		help()
	}

	var id int
	var timestamp int
	var command string

	for rows.Next() {
		rows.Scan(&id, &timestamp, &command)
		t := time.Unix(int64(timestamp), 0)
		fmt.Printf("[%s] %s", datecolor(t.Format("2006.01.02:15:04:05")), command)
	}
}

func help() {
	b := color.FgBlue
	if runtime.GOOS != "windows" {
		b = b.Light()
	}
	blue := b.Render

	fmt.Println("Did you extend your shell?")
	fmt.Println("For " + blue("zsh") + ": ")
	fmt.Println("\texport PROMPT_COMMAND='history | tail -n 1 | cut -c 8- | historian save'" +
		"\n\tprecmd() {eval \"$PROMPT_COMMAND\"}")
	fmt.Println("For " + blue("bash") + ": ")
	fmt.Println("\texport PROMPT_COMMAND='history 1 | cut -c 8- | historian save'")
	os.Exit(1)
}

func main() {
	app.Version("1.0.2")
	usr, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	home = usr.HomeDir
	historian_path = fmt.Sprintf("%s/.historian", home)
	conf = ReadConfigfile(fmt.Sprintf("%s/.config/historian/config.yml", home))

	if conf.Dateformat == "" {
		conf.Dateformat = "2006.01.02:15:04:05"
	}
	datecolor = getColor(conf.DateColor)
	searchcolor = getColor(conf.SearchColor)

	if _, err := os.Stat(historian_path); os.IsNotExist(err) {
		os.Mkdir(historian_path, 0700)
	}

	hdb = fmt.Sprintf("%s/historian.db", historian_path)

	db, err = sql.Open("sqlite3", hdb)
	if err != nil {
		help()
	}
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

		if !bytes.HasPrefix(buf.Bytes(), []byte(os.Args[0])) {
			if !conf.Secret {
				w2db(buf, t)
			} else {
				if !bytes.HasPrefix(buf.Bytes(), []byte(" ")) {
					w2db(buf, t)
				}
			}
		}
	case search.FullCommand():
		sidb()
	}
}
