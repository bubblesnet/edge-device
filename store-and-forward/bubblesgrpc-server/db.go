/*
 * Copyright (c) John Rodley 2022.
 * SPDX-FileCopyrightText:  John Rodley 2022.
 * SPDX-License-Identifier: MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation files (the "Software"), to deal in the
 * Software without restriction, including without limitation the rights to use, copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so, subject to the
 * following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 * INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
 * PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
 * HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
 * CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package main

// copyright and license inspection - no issues 4/13/22

import (
	log "bubblesnet/edge-device/store-and-forward/bubblesgrpc-server/lawg"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"os"
	"time"
)

var csvx = ""

func openWriteable(dbFilename string) {
	log.Debugf("openWriteable %s", dbFilename)
	if writeableDb != nil {
		return
	}
	log.Debugf("bolt.Open %s\n", dbFilename)
	xdb, err := bolt.Open(dbFilename, modeReadwrite, &bolt.Options{Timeout: 1 * time.Second})
	for i := 0; i < 5 && err != nil; i++ {
		log.Warnf("writeable open timed out - sleeping 1 second")
		time.Sleep(time.Second)
		xdb, err = bolt.Open(dbFilename, modeReadwrite, &bolt.Options{Timeout: 1 * time.Second})
	}
	if err != nil {
		log.Errorf("writeable open timed out after 5 attempts in 10 seconds")
		log.Fatalf("%v", err)
	} else {
		log.Infof("Succeeded opening database %s", databaseFilename)
	}
	//	defer func() {
	//		_ = xdb.Close()
	//	}()

	writeableDb = xdb
	log.Debugf(" writeabledb is %v", writeableDb)

	//	defer writeableDb.Close()
}

func makeBuckets(buckets []string) {
	log.Info("makeBuckets\n")
	messageBucketName := buckets[0]
	nodeEnv := os.Getenv("NODE_ENV")
	if nodeEnv == "DEV" {
		err := deleteBucketIfExist(messageBucketName)
		if err != nil {
			log.Warnf("error deleting bucket %s %v ... continuing", messageBucketName, err)
		}
	}
	created, err := makeBucketIfNotExist(messageBucketName)
	if err != nil {
		return
	}
	if created == true {
		log.Debug("Successfully created msg bucket")
	} else {
		log.Debug("msg Bucket already existed")
	}
	stateBucketName := buckets[1]
	if nodeEnv == "DEV" {
		err := deleteBucketIfExist(stateBucketName)
		if err != nil {
			log.Warnf("error deleting bucket %s %v ... continuing", stateBucketName, err)
		}
	}
	created, err = makeBucketIfNotExist(stateBucketName)
	if err != nil {
		return
	}
	if created == true {
		log.Debug("Successfully created state bucket")
	} else {
		log.Debug("state Bucket already existed")
	}
}

func initDb(databaseFilename string) {
	fmt.Printf("initdb %s\n", databaseFilename)
	openWriteable(databaseFilename)
	makeBuckets([]string{messageBucketName, stateBucketName})
}

func deleteBucketIfExist(bucketName string) error {
	log.Debugf("deleteBucketIfExist %s\n", bucketName)
	// Start a writable transaction.
	log.Debugf(" begin writeabledb is %v", writeableDb)
	tx, err := writeableDb.Begin(true)
	if err != nil {
		log.Errorf("begin transaction error %v", err)
		return err
	} else {
		log.Debugf("succeeded transaction start")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	log.Debugf("DeleteBucket %s\n", bucketName)
	// Use the transaction...
	err = tx.DeleteBucket([]byte(bucketName))
	if err != nil {
		log.Errorf("DeleteBucket bucket error %v", err)
		return nil
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func makeBucketIfNotExist(bucketName string) (bool, error) {
	log.Debugf("makeBucketIfNotExist %s\n", bucketName)
	// Start a writable transaction.
	log.Debugf(" begin writeabledb is %v", writeableDb)
	tx, err := writeableDb.Begin(true)
	if err != nil {
		log.Errorf("begin transaction error %v", err)
		return false, err
	} else {
		log.Debugf("succeeded transaction start")
	}
	defer func() {
		_ = tx.Rollback()
	}()

	log.Debugf("CreateBucket %s\n", bucketName)
	// Use the transaction...
	blah, err := tx.CreateBucket([]byte(bucketName))
	if err != nil {
		log.Errorf("Create bucket error %v", err)
		return false, nil
	}
	log.Debugf("bucket create = %v", blah)

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func addRecord(bucketName string, message string, sequence int32) error {
	currentTime := time.Now()
	key := fmt.Sprintf("%s (%6.6d)", currentTime.Format(time.RFC3339), sequence)

	//	key := fmt.Sprintf("%20.20d", currentTime.Unix())
	//	log.Debugf("adding record key %s value %s", key, message)
	_ = writeableDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(key), []byte(message))
		if err != nil {
			log.Errorf("error addRecord %v", err)
		}
		return err
	})
	return nil
}

func deleteFromBucket(bucketName string, key []byte) error {
	//	log.Debugf("deleting key=%s\n", key )
	tx, err := writeableDb.Begin(true)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_ = tx.Bucket([]byte(bucketName)).Delete(key)

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func getStatesAsJson(tx *bolt.Tx) (err error) {

	log.Debugf("pid %d getRecordList getStates", os.Getpid())
	b := tx.Bucket([]byte(stateBucketName))
	count := 0
	log.Debugf("getRecordList foreach")
	csvx = csvx + "[\n"
	_ = b.ForEach(func(k, v []byte) error {
		if count == 0 {
			csvx = csvx + "\n" + string(v)
		} else {
			csvx = csvx + ",\n" + string(v)
		}
		count = count + 1
		return nil
	})
	csvx = csvx + "\n]"
	log.Debugf("getStates - got %d records", count)
	return nil
}

func getStateAsJson(_ string, _ int, _ int, _ int) (string, error) {
	csv := ""
	log.Debugf("pid %d getRecordList getStateAsJson", os.Getpid())
	csvx = ""
	err := writeableDb.View(getStatesAsJson)
	log.Debugf("getStateAsJson Returning nothing %v", err)
	return csv, nil
}
