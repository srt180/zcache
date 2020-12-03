package zcache

import (
	"sync"
	"testing"
	"time"
)

func TestCachemap_StoreWithTimeout(t *testing.T) {
	m := NewCache()

	m.StoreWithTimeout("apple", "red", time.Duration(2)*time.Second)

	time.Sleep(time.Second)

	value, ok := m.load("apple")
	if ok {
		t.Logf("exist: %+v", value)
	} else {
		t.Log("no exist")
	}

	time.Sleep(time.Second * time.Duration(2))
	//m.delete("apple")

	value, ok = m.load("apple")
	if ok {
		t.Logf("exist: %+v", value)
	} else {
		t.Log("no exist")
	}
}

func TestCachemap_Threaded(t *testing.T) {
	m := cachemap{}

	wg := sync.WaitGroup{}

	f := func(index int) {
		//for {
		m.Store("apple", "red")
		_, exist := m.Load("apple")
		if !exist {
			t.Logf("[%d] apple not exist", index)
		} else {
			t.Logf("[%d] apple exist", index)
		}
		m.Delete("apple")
		_, exist = m.Load("apple")
		if !exist {
			t.Logf("[%d] apple not exist after delete", index)
		} else {
			t.Logf("[%d] apple exist after delete", index)
		}
		//}
		wg.Done()
	}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go f(i)
	}

	wg.Wait()

}
