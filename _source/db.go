package main

import (
	"encoding/json"
	"github.com/google/go-github/github"
	"strings"
	"time"
)

const RecordSeparator = "\n"

type DB struct {
	storage Storage
	records RecordCollection
}
type RecordCollection map[int64]Record

type Record struct {
	Id         int64
	CreateTime time.Time
	UpdateTime time.Time
	BaseInfo   PRInfo
	Data       github.PullRequest
	ReviewData []*github.PullRequestReview
}

func (db DB) Init(path string) DB {

	dbNew := DB{
		storage: Storage{}.New(path),
		records: RecordCollection{},
	}
	dbNew.ReadRecords()

	return dbNew
}

func (db *DB) dumpRecords() bool {
	recRaws := []string{}
	for _, record := range db.records {
		recRaw, _ := json.Marshal(record)
		recRaws = append(recRaws, string(recRaw))
	}
	raw := strings.Join(recRaws, RecordSeparator)
	db.storage.ReplaceWith(raw)
	return true
}

func (db *DB) ReadRecords() RecordCollection {
	raw := db.storage.Get()
	recRaws := strings.Split(raw, RecordSeparator)
	for _, recRawSingle := range recRaws {
		var record Record
		_ = json.Unmarshal([]byte(recRawSingle), &record)
		if record.Id == 0 {
			continue
		}
		db.records[record.Id] = record
	}

	return db.records
}

func (db *DB) GetRecord(id int64) Record {
	return db.records[id]
}

func (db *DB) GetAllRecords() RecordCollection {
	return db.records
}

func (db *DB) Add(prInfo PRInfo, pr github.PullRequest, reviewData []*github.PullRequestReview) {

	var rec Record
	if db.Exists(pr.GetID()) {
		rec = db.GetRecord(pr.GetID())
		rec.UpdateTime = time.Now()
		rec.Data = pr
		rec.ReviewData = reviewData
	} else {
		rec = Record{
			Id:         pr.GetID(),
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
			BaseInfo:   prInfo,
			Data:       pr,
			ReviewData: reviewData,
		}
	}
	db.records[pr.GetID()] = rec

	db.dumpRecords()
}

func (db *DB) Remove(id int64) bool {
	if db.Exists(id) {
		delete(db.records, id)
		db.dumpRecords()
	}
	return true
}

func (db *DB) Exists(id int64) bool {
	if _, ok := db.records[id]; ok {
		return true
	}
	return false
}
