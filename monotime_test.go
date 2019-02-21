// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package monotime

import (
	"testing"
	"io/ioutil"
	"fmt"
	"encoding/binary"
	"time"
	"os"
)

func TestNow(t *testing.T) {
	for i := 0; i < 100; i++ {
		t1 := Now()
		t2 := Now()
		// I honestly thought that we needed >= here, but in some environments
		// two consecutive calls can return the same value!
		if t1 > t2 {
			t.Fatalf("t1=%d should have been less than or equal to t2=%d", t1, t2)
		}
	}
}

func TestSince(t *testing.T) {
	for i := 0; i < 100; i++ {
		ts := Now()
		d := Since(ts)
		if d < 0 {
			t.Fatalf("d=%d should be greater than or equal to zero", d)
		}
	}
}

func TestPersist(t *testing.T) {
	// File to persist time
	fileName := "time-file"

	// Create if not exists
	os.OpenFile(fileName, os.O_CREATE, 0666)

	// Read persisted time
	buff, err := ioutil.ReadFile(fileName)
	check(err)

	currTime := Raw()
	fmt.Printf("Current Time\t\t\t= %dns\n", currTime)

	if len(buff) != 8 {
		fmt.Println("No Persisted Time!")
	} else {
		persistedTime := int64(binary.LittleEndian.Uint64(buff))
		fmt.Printf("Persisted Time\t\t\t= %dns\n", persistedTime)

		if persistedTime > currTime {
			defer t.Fatalf("Persisted Time > Current Time!")
		}

		elapsed := Since(time.Duration(persistedTime))
		fmt.Printf("Duration since last persist\t= %s\n", elapsed)
	}

	buff = make([]byte, 8)
	fmt.Println("Writing Current time!")
	binary.LittleEndian.PutUint64(buff, uint64(currTime))

	// Write current time to file
	err = ioutil.WriteFile(fileName, buff, 0666)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
