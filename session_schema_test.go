// Copyright 2017 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xorm

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStoreEngine(t *testing.T) {
	assert.NoError(t, prepareEngine())

	assert.NoError(t, testEngine.DropTables(context.Background(), "user_store_engine"))

	type UserinfoStoreEngine struct {
		Id   int64
		Name string
	}

	assert.NoError(t, testEngine.StoreEngine("InnoDB").Table("user_store_engine").CreateTable(context.Background(), &UserinfoStoreEngine{}))
}

func TestCreateTable(t *testing.T) {
	assert.NoError(t, prepareEngine())

	assert.NoError(t, testEngine.DropTables(context.Background(), "user_user"))

	type UserinfoCreateTable struct {
		Id   int64
		Name string
	}

	assert.NoError(t, testEngine.Table("user_user").CreateTable(context.Background(), &UserinfoCreateTable{}))
}

func TestCreateMultiTables(t *testing.T) {
	assert.NoError(t, prepareEngine())

	session := testEngine.NewSession()
	defer session.Close()

	type UserinfoMultiTable struct {
		Id   int64
		Name string
	}

	user := &UserinfoMultiTable{}
	assert.NoError(t, session.Begin())

	for i := 0; i < 10; i++ {
		tableName := fmt.Sprintf("user_%v", i)

		assert.NoError(t, session.DropTable(context.Background(), tableName))

		assert.NoError(t, session.Table(tableName).CreateTable(context.Background(), user))
	}

	assert.NoError(t, session.Commit())
}

type SyncTable1 struct {
	Id   int64
	Name string
	Dev  int `xorm:"index"`
}

type SyncTable2 struct {
	Id     int64
	Name   string `xorm:"unique"`
	Number string `xorm:"index"`
	Dev    int
	Age    int
}

func (SyncTable2) TableName() string {
	return "sync_table1"
}

func TestSyncTable(t *testing.T) {
	assert.NoError(t, prepareEngine())

	assert.NoError(t, testEngine.Sync2(context.Background(), new(SyncTable1)))

	assert.NoError(t, testEngine.Sync2(context.Background(), new(SyncTable2)))
}

func TestIsTableExist(t *testing.T) {
	assert.NoError(t, prepareEngine())

	exist, err := testEngine.IsTableExist(context.Background(), new(CustomTableName))
	assert.NoError(t, err)
	assert.False(t, exist)

	assert.NoError(t, testEngine.CreateTables(context.Background(), new(CustomTableName)))

	exist, err = testEngine.IsTableExist(context.Background(), new(CustomTableName))
	assert.NoError(t, err)
	assert.True(t, exist)
}

func TestIsTableEmpty(t *testing.T) {
	assert.NoError(t, prepareEngine())

	type NumericEmpty struct {
		Numeric float64 `xorm:"numeric(26,2)"`
	}

	type PictureEmpty struct {
		Id          int64
		Url         string `xorm:"unique"` //image's url
		Title       string
		Description string
		Created     time.Time `xorm:"created"`
		ILike       int
		PageView    int
		From_url    string
		Pre_url     string `xorm:"unique"` //pre view image's url
		Uid         int64
	}

	assert.NoError(t, testEngine.DropTables(context.Background(), &PictureEmpty{}, &NumericEmpty{}))

	assert.NoError(t, testEngine.Sync2(context.Background(), new(PictureEmpty), new(NumericEmpty)))

	isEmpty, err := testEngine.IsTableEmpty(context.Background(), &PictureEmpty{})
	assert.NoError(t, err)
	assert.True(t, isEmpty)

	tbName := testEngine.GetTableMapper().Obj2Table("PictureEmpty")
	isEmpty, err = testEngine.IsTableEmpty(context.Background(), tbName)
	assert.NoError(t, err)
	assert.True(t, isEmpty)
}

type CustomTableName struct {
	Id   int64
	Name string
}

func (c *CustomTableName) TableName() string {
	return "customtablename"
}

func TestCustomTableName(t *testing.T) {
	assert.NoError(t, prepareEngine())

	c := new(CustomTableName)
	assert.NoError(t, testEngine.DropTables(context.Background(), c))

	assert.NoError(t, testEngine.CreateTables(context.Background(), c))
}

func TestDump(t *testing.T) {
	assert.NoError(t, prepareEngine())

	fp := testEngine.Dialect().URI().DbName + ".sql"
	os.Remove(fp)
	assert.NoError(t, testEngine.DumpAllToFile(context.Background(), fp))
}

type IndexOrUnique struct {
	Id        int64
	Index     int `xorm:"index"`
	Unique    int `xorm:"unique"`
	Group1    int `xorm:"index(ttt)"`
	Group2    int `xorm:"index(ttt)"`
	UniGroup1 int `xorm:"unique(lll)"`
	UniGroup2 int `xorm:"unique(lll)"`
}

func TestIndexAndUnique(t *testing.T) {
	assert.NoError(t, prepareEngine())

	assert.NoError(t, testEngine.CreateTables(context.Background(), &IndexOrUnique{}))

	assert.NoError(t, testEngine.DropTables(context.Background(), &IndexOrUnique{}))

	assert.NoError(t, testEngine.CreateTables(context.Background(), &IndexOrUnique{}))

	assert.NoError(t, testEngine.CreateIndexes(context.Background(), &IndexOrUnique{}))

	assert.NoError(t, testEngine.CreateUniques(context.Background(), &IndexOrUnique{}))

	assert.NoError(t, testEngine.DropIndexes(context.Background(), &IndexOrUnique{}))
}

func TestMetaInfo(t *testing.T) {
	assert.NoError(t, prepareEngine())
	assert.NoError(t, testEngine.Sync2(context.Background(), new(CustomTableName), new(IndexOrUnique)))

	tables, err := testEngine.DBMetas(context.Background())
	assert.NoError(t, err)
	assert.EqualValues(t, 2, len(tables))
	tableNames := []string{tables[0].Name, tables[1].Name}
	assert.Contains(t, tableNames, "customtablename")
	assert.Contains(t, tableNames, "index_or_unique")
}

func TestCharst(t *testing.T) {
	assert.NoError(t, prepareEngine())

	err := testEngine.DropTables(context.Background(), "user_charset")
	if err != nil {
		t.Error(err)
		panic(err)
	}

	err = testEngine.Charset("utf8").Table("user_charset").CreateTable(context.Background(), &Userinfo{})
	if err != nil {
		t.Error(err)
		panic(err)
	}
}
