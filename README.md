# historian

## What's this?
If you are connected to servers with more than one session, you can run into some problems with your bash-history.
This project seeks to offer a solution for this problem.

## What's to do?
First of all you should edit your .bashrc and add the following to the botton:
```
export PROMPT_COMMAND='history 1 | cut -c 8- | historian save'
```

Or if you use zshell, edit .zshrc and add this to the bottom:
```
export PROMPT_COMMAND='history | tail -n 1 | cut -c 8- | historian save'
precmd() {eval "$PROMPT_COMMAND"}
```

Historian writes your commands into an sqlite db. With historian you can search for commands and the searched term will be highlighted if it was found.

## Usage
Basic:
```
usage: historian [<flags>] <command> [<args> ...]

I store your history and search it for you

Flags:
  --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  save
    Save a command

  search [<criteria>...]
    search for a command

  integrate
    integrate the old historian version into your database
```

## How to install?
Just download the file and put it somewhere where your $PATH points to.
e.g.: /usr/bin

## About
Please keep in mind that this is software isn't finished and there can be bugs.
If anyone has ideas to extend this project, you are welcome to tell me.
