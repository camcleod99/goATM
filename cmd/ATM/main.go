package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	ATM "github.com/camcleod99/TerminalATM"
	"io"
	"os"
	"strings"
)

const (
	accountFile = ".account.json"
)

func main() {
	add := flag.Bool("add", false, "Add a new transaction")
	correct := flag.Int("correct", 0, "Correct a transaction")
	remove := flag.Int("remove", 0, "delete a transaction")
	list := flag.Bool("list", false, "List all transactions")

	flag.Parse()

	transactions := &ATM.Transactions{}

	if err := transactions.Load(accountFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		name, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		transactions.Add(name, "Debt", 12.99)
		err = transactions.Store(accountFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *correct > 0:
		err := transactions.Correct(*correct, 10.99)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = transactions.Store(accountFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *remove > 0:
		err := transactions.Delete(*remove)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = transactions.Store(accountFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

	case *list:
		transactions.Print()

	default:
		fmt.Fprintln(os.Stdout, "invalid Command")
		os.Exit(0)
	}
}

func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("input String not valid")
	}

	// Unwrap here!

	return text, nil
}
