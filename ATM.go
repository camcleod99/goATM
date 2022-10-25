package ATM

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (t *Transactions) Add(name string, action string, amount float32) {

	// Sanity checks here:

	transaction := transaction{
		Name:    name,
		Action:  action,
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

func (t *Transactions) Delete(index int) error {
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
		fmt.Printf("%d - %s\n", i, transaction.Name)
	}
}
