package user

import (
	"fmt"
	"sync"
	"time"
)

type Memory struct {
	locker sync.Mutex
	sync.Map
	cond      *sync.Cond
	ispresent *bool
}

func NewMemory() *Memory {
	mem := &Memory{
		cond:      sync.NewCond(nil),
		ispresent: new(bool),
	}

	go mem.Loop()
	return mem
}
func (m *Memory) AddClient(code string, detail *Details) {
	m.Store(code, detail)
	//n := *detail
	//fmt.Printf("user  %s added in memory\n", n.Name)
}
func (m *Memory) Update(isok *bool) {
	m.locker.Lock()
	defer m.locker.Unlock()

	m.ispresent = isok
}
func (m *Memory) Isidle() bool {
	m.locker.Lock()
	defer m.locker.Unlock()

	return *m.ispresent
}
func (m *Memory) Loop() {
	for {

		m.Range(func(key, value interface{}) bool {

			myval := value.(*Details)
			currtime := time.Now()
			if currtime.After(myval.TTL) {
				fmt.Printf("Deleting the user %v\n", myval)
				m.Delete(key)
			}
			return true
		})
		m.HasValue()
		if !m.Isidle() {
			var mu sync.Mutex
			m.cond.L = &mu
			mu.Lock()
			m.cond.Wait()
			mu.Unlock()
			fmt.Println("unlocking idle")
		}

	}
}

func (m *Memory) Check(code string) (Details, bool) {
	var found Details
	var isok bool
	m.Range(func(key, value interface{}) bool {
		mykey := key.(string)
		if mykey == code {
			found = *(value.(*Details))
			isok = true
			m.Delete(key)

			return false
		}
		isok = false
		return true
	})
	fmt.Println(found, isok)

	return found, isok
}
func (m *Memory) HasValue() {
	var has bool = false
	m.Range(func(key, value interface{}) bool {
		has = true

		return false
	})
	m.Update(&has)

}
func (m *Memory) Awake() {
	isok := m.Isidle()
	fmt.Println(isok)
	if !m.Isidle() {
		//fmt.Println("awakeing")
		m.cond.Signal()
	}
}
