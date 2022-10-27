package ATM

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"time"
)

type transaction struct {
	Name    string
	Action  string
	Amount  float32
	Created time.Time
	Edited  time.Time
}

type Transactions []transaction

func (t *Transactions) Init(fileName string) error {

	// Creates Directory if it doesn't exist
	_, err := os.Stat("data")
	if os.IsNotExist(err) {
		err = os.Mkdir("data", os.ModePerm)
		if err != nil {
			return err
		}

	}

	// Creates File if it doesn't exist
	_, err = os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		data := ``
		err = os.WriteFile(fileName, []byte(data), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Transactions) Refresh(fileName string) error {
	/*
		Ask if you want to refresh
		If no
			err = Usr Canceled
			return error
	*/

	const refreshMsg = "Refreshing the account data. \n WARNING THIS WILL DELETE ALL DATA! \n DO YOU WANT TO CONTINUE? \n "
	fmt.Print(refreshMsg)
	prompt := promptui.Select{
		Label: "[Yes/No]",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return err
	}

	if result == "Yes" {
		// Delete File
		if err := os.Remove(fileName); err != nil {
			return err
		}

		//Create File
		data := ``
		if err := os.WriteFile(fileName, []byte(data), 0644); err != nil {
			return err
		}
		return nil
	} else {
		return fmt.Errorf("⚠️ user has canceled the action")
	}

}

func (t *Transactions) Add(name string, amount float32) {

	// Sanity checks here:

	transaction := transaction{
		Name:    name,
		Action:  "-",
		Amount:  amount,
		Created: time.Now(),
		Edited:  time.Now(),
	}

	*t = append(*t, transaction)
}

func (t *Transactions) Debit(name string, amount float32) {

	transaction := transaction{
		Name:    name,
		Action:  "debit",
		Amount:  amount,
		Created: time.Now(),
		Edited:  time.Now(),
	}

	*t = append(*t, transaction)
}

func (t *Transactions) Credit(name string, amount float32) {

	transaction := transaction{
		Name:    name,
		Action:  "credit",
		Amount:  amount,
		Created: time.Now(),
		Edited:  time.Now(),
	}

	*t = append(*t, transaction)
}

func (t *Transactions) Correct(index int, amount float32) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("invalid Index")
	}

	ls[index-1].Amount = amount
	ls[index-1].Edited = time.Now()

	return nil
}

func (t *Transactions) Remove(index int) error {
	ls := *t

	if index <= 0 || index > len(ls) {
		return errors.New("invalid Index")
	}

	*t = append(ls[:index-1], ls[:index]...)

	return nil
}

func (t *Transactions) Load(fileName string) error {
	file, err := os.ReadFile(fileName)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return err
	}

	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

func (t *Transactions) Store(fileName string) error {
	data, err := json.Marshal(t)

	if err != nil {
		return err
	}

	return os.WriteFile(fileName, data, 0644)
}

func (t *Transactions) Print() {
	for i, transaction := range *t {
		i++
		fmt.Printf("%d - %s : %g \n", i, transaction.Name, transaction.Amount)
	}
}
