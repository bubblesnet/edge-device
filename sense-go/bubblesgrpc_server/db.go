package main

import (
	"bytes"
	"fmt"
	"github.com/go-playground/log"
	bolt "go.etcd.io/bbolt"
	"os"
	"time"
)

var csvx = ""

func openWriteable() {
	if writeableDb != nil {
		return
	}
	xdb, err := bolt.Open(databaseFilename, mode_readwrite, &bolt.Options{Timeout: 1 * time.Second})
	for i := 0; i < 5 && err == nil; i++ {
		if err == nil {
			break
		}
		log.Warn(fmt.Sprintf("writeable open timed out - sleeping 1 second"))
		time.Sleep(time.Second)
		xdb, err = bolt.Open(databaseFilename, mode_readwrite, &bolt.Options{Timeout: 1 * time.Second})

	}
	if err != nil {
		log.Error(fmt.Sprintf("writeable open timed out after 5 attempts in 10 seconds"))
		log.Fatal(err)
	}

	writeableDb = xdb
	//	defer writeableDb.Close()
}

func makeBuckets() {
	created, err := makeBucketIfNotExist(messageBucketName)
	if err != nil {
		return
	}
	if created == true {
		log.Debug("Successfully created msg bucket")
	} else  {
		log.Debug("msg Bucket already existed")
	}
	created, err = makeBucketIfNotExist(stateBucketName)
	if err != nil {
		return
	}
	if created == true {
		log.Debug("Successfully created state bucket")
	} else  {
		log.Debug("state Bucket already existed")
	}
}

func initDb() {
	openWriteable()
	makeBuckets()
}

func makeBucketIfNotExist(bucketName string) (bool, error) {
	// Start a writable transaction.
	tx, err := writeableDb.Begin(true)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	// Use the transaction...
	blah, err := tx.CreateBucket([]byte(bucketName))
	if err != nil {
		log.Error(fmt.Sprintf("Create bucket error %v", err))
		return false, nil
	}
	log.Debug( fmt.Sprintf("bucket create = %v", blah))

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func addRecord(bucketName string, message string) error {
	currentTime := time.Now()
	key := currentTime.Format(time.RFC3339)

//	key := fmt.Sprintf("%20.20d", currentTime.Unix())
	log.Debug(fmt.Sprintf("adding record key %s value %s", key, message))
	writeableDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(key), []byte(message))
		if err != nil {
			log.Error(fmt.Sprintf("error addRecord %v", err ))
		}
		return err
	})
	return nil
}

func deleteFromBucket( bucketName string, key []byte ) error {
	//	log.Debug(fmt.Sprintf("deleting key=%s\n", key )
	tx, err := writeableDb.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	tx.Bucket([]byte(bucketName)).Delete(key)

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func clearDatabase( bucketName string ) {
	var deleteThem []string

	prefix := ""
	log.Info(fmt.Sprintf("Deleting records prefixed %s", prefix ))

	writeableDb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket([]byte(bucketName)).Cursor()
		// Iterate over the 90's.
		log.Debug("Deleting - let's seek")
//		for k, _ := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, _ = c.Next()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
				log.Debug(fmt.Sprintf("found key=%s", k))
				deleteThem = append(deleteThem, string(k))
		}
		log.Debug("Done finding")

		return nil
	})
	for _, element := range deleteThem {
		log.Debug("Deleting key %s", element )
		deleteFromBucket(stateBucketName, []byte(element))
	}
	log.Debug("Done finding")
}


func deletePriorTo( bucketName string, unixtime int64 ) {

	var max = fmt.Sprintf("%20.20d", unixtime )
	min := []byte("0")
	log.Info(fmt.Sprintf("Deleting records between %s and %s", min, max))

	writeableDb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		c := tx.Bucket([]byte(bucketName)).Cursor()
		// Iterate over the 90's.
		log.Debug("Deleting - let's seek")
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, []byte(max)) <= 0; k, v = c.Next() {
			log.Debug(fmt.Sprintf("Deleting key %s: %s", k, v))
		}

		return nil
	})
}

func getStatesAsJson(tx *bolt.Tx) error {

	log.Debug(fmt.Sprintf("pid %d getRecordList getStates", os.Getpid()))
	b := tx.Bucket([]byte(stateBucketName))
	count := 0
	log.Debug(fmt.Sprintf("getRecordList foreach"))
	csvx = csvx + "[\n"
	b.ForEach(func(k, v []byte) error {
		if count == 0 {
			csvx = csvx + "\n" + string(v)
		} else {
			csvx = csvx + ",\n" + string(v)
		}
		count = count + 1
		return nil
	})
	csvx = csvx + "\n]"
	log.Debug(fmt.Sprintf("getStates - got %d records", count ))
	return nil
}

func getStatesAsCsv(tx *bolt.Tx) error {
	log.Debug(fmt.Sprintf("pid %d getRecordList getStates", os.Getpid()))
	b := tx.Bucket([]byte(stateBucketName))
	count := 0
	log.Debug(fmt.Sprintf("getRecordList foreach"))
	b.ForEach(func(k, v []byte) error {
		count = count + 1
		csvx = csvx + "\n" + string(v)
		return nil
	})
	log.Debug(fmt.Sprintf("getStates - got %d records", count ))
	return nil
}

func getStateAsCsv( bucketName string, year int, month int, day int) (string, error) {
	csv := ""
	log.Debug(fmt.Sprintf("pid %d getRecordList getStateAsCsv", os.Getpid()))
	csvx = ""
	err := writeableDb.View(getStatesAsCsv)
	log.Debug(fmt.Sprintf("getStateAsCsv Returning nothing %v", err))
	return csv, nil
}

func getStateAsJson( bucketName string, year int, month int, day int) (string, error) {
	csv := ""
	log.Debug(fmt.Sprintf("pid %d getRecordList getStateAsCsv", os.Getpid()))
	csvx = ""
	err := writeableDb.View(getStatesAsJson)
	log.Debug(fmt.Sprintf("getStateAsCsv Returning nothing %v", err))
	return csv, nil
}