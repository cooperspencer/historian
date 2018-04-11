# historian

## What's this?
If you are connected to servers with more than one session, you can run into some problems with your bash-history.
This project seeks to offer a solution for this problem.

## What's to do?
First of all you should edit your .bashrc and add the following to the botton:
```
export PROMPT_COMMAND='if [ ! -e ~/.logs ]; then mkdir ~/.logs; fi; echo "[$(date "+%Y-%m-%d.%H:%M:%S")] $(history 1 | cut -c 8-)" >> ~/.logs/bash-history.log'
```

This puts the commands entered into a separated file, so it won't mess up your history, and you even get a proper timestamp.

Now you could use my tool historian.
It is capable of searching for specific dates and entered commands and higlighting those.
If you think that one month of bash-history is enough, you can tell it to delete everything older than 30 days.

## Usage
Basic:
```
usage: historian [<flags>] <command> [<args> ...]

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  search [<flags>]
    Search for commands in the history

  delete [<flags>]
    Delete older than X days
```

Search:
```
usage: historian search [<flags>]

Search for commands in the history

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
  -d, --date=DATE            Dateformat like 2018-04-10
  -c, --command=COMMAND ...  A command like ping
  -f, --from=FROM            Search from f.e.: 13:00
  -t, --to=TO                Search to f.e.: 13:00
```

Delete:
```
usage: historian delete [<flags>]

Delete older than X days

Flags:
      --help       Show context-sensitive help (also try --help-long and --help-man).
  -d, --days=DAYS  Delete X Days of history
```

## How to install?
Just download the file and put it somewhere where your $PATH points to.
e.g.: /usr/bin

## About
Please keep in mind that this is software isn't finished and there can be bugs.
If anyone has ideas to extend this project, you are welcome to tell me.
