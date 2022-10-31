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

// Initialize custom values here

const (
	accountFile = "./data/account.json"
)

func main() {
	add := flag.Bool("add", false, "Add a new transaction")
	credit := flag.Bool("credit", false, "Add a new credit transaction")
	debit := flag.Bool("debit", false, "Add a new debit transaction")
	correct := flag.Int("correct", 0, "Correct a transaction")
	remove := flag.Int("remove", 0, "delete a transaction")
	refresh := flag.Bool("refresh", false, "Refreshes the account ledger")
	list := flag.Bool("list", false, "List all transactions")

	flag.Parse()

	transactions := &ATM.Transactions{}

	if err := transactions.Init(accountFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := transactions.Load(accountFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	switch {
	// Example command -add [name]
	case *add:
		name, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		transactions.Add(name, 12.99)
		err = transactions.Store(accountFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Printf("Transaction Added \n")

	// Example Command -debit amount, name
	case *debit:
		name, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Print(name)
		fmt.Print(err)
		fmt.Printf("Hit the debit action \n")

	case *credit:
		name, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Print(name)
		fmt.Print(err)
		fmt.Printf("Hit the credit action \n")

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
		fmt.Printf("Transaction Corrected \n")

	case *remove > 0:
		err := transactions.Remove(*remove)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = transactions.Store(accountFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Printf("Transaction Removed \n")

	case *refresh:
		err := transactions.Refresh(accountFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Printf("Account List Refreshed \n")

	case *list:
		transactions.Print()

	default:
		fmt.Fprintln(os.Stdout, "invalid Command")
		os.Exit(0)
	}
	fmt.Printf("Done.\n")
}

func getInput(r io.Reader, args ...string) (string, error) {

	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	fmt.Printf("Empty Input, please type in a Command;\n")
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()

	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed")
	}

	return text, nil
}
