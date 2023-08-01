package sqlite

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"gorm.io/gorm"
)

type person struct {
	gorm.Model
	Like string
	Age  int64
	Type personType
}

const (
	normal personType = iota
	vip
)

type personType int

func (p personType) Value() (driver.Value, error) {
	return int64(p), nil
}

func (p *personType) Scan(value interface{}) error {
	switch v := value.(type) {
	case int64:
		*p = personType(v)
		return nil
	}

	return errors.New("unsupported type")
}

func TestMigrate(t *testing.T) {
	db, err := gorm.Open(Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&person{}); err != nil {
		t.Fatal(err)
	}
}

func TestInsert(t *testing.T) {
	db, err := gorm.Open(Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	s := []person{
		{
			Like: "apple",
			Age:  18,
			Type: normal,
		},
		{
			Like: "banana",
			Age:  20,
			Type: vip,
		},
	}
	res := db.Create(s)
	if res.Error != nil {
		t.Fatal(res.Error)
	}

	for _, v := range s {
		fmt.Println(v.ID)
	}
}

func TestConcurInsert(t *testing.T) {
	db, err := gorm.Open(Open("test.db?_pragma=journal_mode(wal)"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	begin := time.Now()
	var w sync.WaitGroup
	for i := 0; i < 100; i++ {
		w.Add(1)
		go func(n int) {
			defer w.Done()

			err := db.Create(&person{
				Like: "thing" + strconv.Itoa(n),
				Age:  int64(n),
				Type: normal,
			}).Error
			if err != nil {
				fmt.Println(err.Error())
			}
		}(i)
	}

	w.Wait()
	fmt.Printf("cost: %d ms", time.Since(begin).Milliseconds())
}

func TestQuery(t *testing.T) {
	db, err := gorm.Open(Open("test.db"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	var persons []person
	if err := db.Find(&persons).Error; err != nil {
		t.Fatal(err)
	}

	for _, v := range persons {
		fmt.Printf("%+v \n", v)
	}
}
