package main

import (
	"errors"
	"flag"
	"fmt"
	ATM "github.com/camcleod99/TerminalATM"
	"os"
	"regexp"
	"strings"
)

// Initialize custom values here

const (
	accountFile = "./data/account.json"
)

func main() {
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

	case *debit:
		inString, err := getInput(flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		workString := strings.Split(inString, ",")

		name, amount, err := washInput(workString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		transactions.Add(name, amount, "debit")
		err = transactions.Store(accountFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Name :" + name + "\n Amount: £" + amount + "\n")
		fmt.Printf("Account Debited \n")

	case *credit:
		inString, err := getInput(flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		workString := strings.Split(inString, ",")

		name, amount, err := washInput(workString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		transactions.Add(name, amount, "credit")
		err = transactions.Store(accountFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Name :" + name + "\n Amount: £" + amount + "\n")
		fmt.Printf("Account Credited \n")

	// TODO: THIS
	case *correct > 0:
		inString, err := getInput(flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		workString := strings.Split(inString, ",")

		name, amount, err := washInput(workString)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf(name + "\n")

		// Get item at *correct
		// Compare name to name in *correct
		// If not --- and different ask to continue

		/*
			lastName := getName(*correct)
			if name != lastname {
				name = lastname
			}
		*/

		err = transactions.Correct(*correct, name, amount)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		err = transactions.Store(accountFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Name :" + name + "\n Amount: £" + amount + "\n")
		fmt.Printf("Transaction Corrected \n")
	// TODO: THAT

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

func getInput(args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	} else {
		return "", errors.New("arguments Required")
	}
}

func washInput(inWashString []string) (string, string, error) {
	var outName string
	var inAmount string
	var outAmount string

	/*	Case 0 Shouldn't happen but is here to catch getInput() screwing up
		Case 1 Supposed one input; a money amount
		Case 2 Supposes two inputs; a name and a money amount	*/

	switch len(inWashString) {

	case 0:
		err := errors.New("⚠️ No Parameters Given")
		return "", "", err

	case 1:
		outName, inAmount = "---", strings.TrimSpace(inWashString[0])

	case 2:
		outName, inAmount = strings.TrimSpace(inWashString[0]), strings.TrimSpace(inWashString[1])

	default:
		outName, inAmount = strings.TrimSpace(inWashString[0]), strings.TrimSpace(inWashString[1])
	}

	/*	These are to match inAmount to how a user would input the amount;
		ether as an int (ie; they had £25, so typing 25.00 would be superfluous
		or as a float (it; they had £24.99)
		Other formats IE: 24.999 or text input, will not be caught and thus will
		cause the program to error out
		Int Regex = ^\d+$ (Start, Any number of digits, end)
		float Regex = ^\d+\.\d{2}$ (Start, any number of digits, a period, two digits, end)	*/

	matchNoCoin, _ := regexp.Match("^\\d+$", []byte(inAmount))
	matchCoin, _ := regexp.Match("^\\d+\\.\\d{2}$", []byte(inAmount))

	if matchNoCoin {
		outAmount = inAmount + ".00"
	} else if matchCoin {
		outAmount = inAmount
	} else {
		err := errors.New("⚠️ Bad money format - XX or XX.XX required️")
		return "", "", err
	}

	return outName, outAmount, nil
}
