package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"github.com/fatih/color"
	"time"
	"strconv"
)

var (
	app = kingpin.New("historian", "A command-line too to get your commands")
	search	= app.Command("search", "Search for commands in the history")
	searchDate = search.Flag("date", "Dateformat like 2018-04-10").Short('d').String()
	searchCommand = search.Arg("command", "A command like ping").Strings()
	searchFrom = search.Flag("from", "Search from f.e.: 13:00").Short('f').String()
	searchTo = search.Flag("to", "Search to f.e.: 13:00").Short('t').String()
	delete = app.Command("delete", "Delete older than X days")
	deleteDays = delete.Flag("days", "Delete X Days of history").Short('d').Int()
	version = app.Command("version", "shows the version of historian")
	vnr = "0.0.2"
	yellow  = color.New(color.FgYellow).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	green   = color.New(color.FgGreen).SprintfFunc()
	def     = color.New(color.FgWhite).SprintfFunc()
	lgs     = []string{}
)

func searchfordate(text string) {
	if *searchDate == "" {
		now := time.Now()
		*searchDate = fmt.Sprintf("%d-%02d-%02d", now.Year(), int(now.Month()), now.Day())
	}

	h, m, s, sh, sm, ss := 0, 0, 0, 0, 0, 0
	th, tm, ts := 0, 0, 0
	timeAndDate := strings.Split(text, " ")[0]
	timeAndDate = timeAndDate[1:]
	d := strings.Split(strings.Split(timeAndDate, ".")[0], "-")
	year, err := strconv.Atoi(d[0])
	month, err := strconv.Atoi(d[1])
	day, err := strconv.Atoi(d[2])

	if err != nil {
		fmt.Println("There's a problem with the date in the logs")
		os.Exit(1)
	}

	syear := 0
	smonth := 0
	sday := 0

	searchD := strings.Split(*searchDate, "-")
	if len(searchD) == 3 {
		syear, err = strconv.Atoi(searchD[0])
		smonth, err = strconv.Atoi(searchD[1])
		sday, err = strconv.Atoi(searchD[2])
	} else {
		fmt.Println("There's a problem with the date you searched")
		os.Exit(2)
	}

	if err != nil {
		fmt.Println("There's a problem with the date you searched")
		os.Exit(1)
	}

	if *searchFrom != "" && *searchTo != "" {
		t := strings.Split(strings.Split(timeAndDate[:len(timeAndDate)-1], ".")[1], ":")
		h, err = strconv.Atoi(t[0])
		m, err = strconv.Atoi(t[1])
		s, err = strconv.Atoi(t[2])
		sf := strings.Split(*searchFrom, ":")
		sh, err = strconv.Atoi(sf[0])
		sm, err = strconv.Atoi(sf[1])
		st := strings.Split(*searchTo, ":")
		th, err = strconv.Atoi(st[0])
		tm, err = strconv.Atoi(st[1])

	}

	sD := time.Date(syear, time.Month(smonth), sday, sh, sm, ss, 0, time.UTC)
	logdate := time.Date(year, time.Month(month), day, h, m, s, 0, time.UTC)
	tD := time.Date(syear, time.Month(smonth), sday, th, tm, ts, 0, time.UTC)

	splitted := strings.Split(text, " ")
	sd := splitted[0]
	splitted = append(splitted[:0], splitted[1:]...)

	if *searchFrom != "" && *searchTo != "" {
		if logdate.Before(tD) && logdate.After(sD) {
			if len(*searchCommand) > 0 {
				searchforcommand(text)
				return
			}
			fmt.Printf("%s %s\n", red(sd), strings.Join(splitted, " "))
		}
	}else {
		if logdate.Equal(sD) {
			if len(*searchCommand) > 0 {
				searchforcommand(text)
				return
			}
			fmt.Printf("%s %s\n", red(sd), strings.Join(splitted, " "))
		}
	}
}

func searchforcommand(text string) {
	splitted := strings.Split(text, " ")
	d := splitted[0]
	splitted = append(splitted[:0], splitted[1:]...)
	comm := strings.Join(splitted, " ")
	c_comm := comm
	match := 0
	for _, c := range *searchCommand {
		if strings.Contains(comm, c) {
			match++
			c_comm = strings.Replace(c_comm, c, fmt.Sprintf("%s", green(c)), -1)
		}
	}
	if match == len(*searchCommand) {
		fmt.Printf("%s %s\n", red(d), c_comm)
	}
}

func deleteoldentries(text string) (del bool) {
	timeAndDate := strings.Split(text, " ")[0]
	timeAndDate = timeAndDate[1:]
	d := strings.Split(strings.Split(timeAndDate, ".")[0], "-")
	year, err := strconv.Atoi(d[0])
	month, err := strconv.Atoi(d[1])
	day, err := strconv.Atoi(d[2])
	if err != nil {
		fmt.Println("There's a problem with the date in the logs")
		return false
	}
	logdate := time.Date(year, time.Month(month), day, 0, 0,0, 0, time.UTC)

	timetoDelete := -24 * *deleteDays
	t := time.Now().Add(time.Duration(timetoDelete) * time.Hour)

	if logdate.After(t) {
		return true
	} else {

		return false
	}
}

func Search() {
	if *searchDate != "" {
		fmt.Printf("Searching for commands on %s\n", *searchDate)
	}
	if len(*searchCommand) > 0 {
		fmt.Printf("Searching for the command %s\n", *searchCommand)
	}

	usr, err := user.Current()

	home := usr.HomeDir
	logfile := fmt.Sprintf("%s/.logs/bash-history.log", home)

	file, err := os.Open(logfile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if len(*searchCommand) > 0 && (*searchDate == "" || (*searchTo != "" && *searchFrom != "")) {
			searchforcommand(scanner.Text())
		} else {
			if *searchDate != "" || (*searchTo != "" && *searchFrom != "") {
				searchfordate(scanner.Text())
			} else {
				splitted := strings.Split(scanner.Text(), " ")
				d := splitted[0]
				splitted = append(splitted[:0], splitted[1:]...)
				fmt.Printf("%s %s\n", red(d), strings.Join(splitted, " "))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	if len(lgs) != 0 {
		file, err = os.Create(logfile)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()
		w := bufio.NewWriter(file)
		for _, line := range lgs {
			fmt.Fprintln(w, line)
		}
		w.Flush()
	}
}

func Delete() {
	usr, err := user.Current()

	home := usr.HomeDir
	logfile := fmt.Sprintf("%s/.logs/bash-history.log", home)

	file, err := os.Open(logfile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if *deleteDays != 0 {
			if deleteoldentries(scanner.Text()) {
				lgs = append(lgs, scanner.Text())
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	if len(lgs) != 0 {
		file, err = os.Create(logfile)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()
		w := bufio.NewWriter(file)
		for _, line := range lgs {
			fmt.Fprintln(w, line)
		}
		w.Flush()
	}
}

func ShowVersion()  {
	fmt.Println(vnr)
}

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case search.FullCommand():
		Search()
	case version.FullCommand():
		ShowVersion()
	case delete.FullCommand():
		Delete()
	}
}