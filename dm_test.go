package dm

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/jiangliuhong/gorm-driver-dm/dmr"
	dmSchema "github.com/jiangliuhong/gorm-driver-dm/schema"
	"gorm.io/gorm"
	"testing"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	//dsn := "dm://sysdba:SYSDBA@local.Leefs.ren:5236?autoCommit=true"
	dsn := "dm://IACSTORE:IACSTORE@114.55.86.238:5236?schema=IACSTORE&autoCommit=true"
	db, err = gorm.Open(Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		fmt.Printf("Error: failed to connect dm server: %v\n", err)
		return
	}

	TB(db)
}

type User struct {
	gorm.Model
	Key      string `gorm:"index:key,unique"`
	Name     string `gorm:"index:name"`
	Age      int
	Content  dmSchema.Clob `gorm:"size:1024000"`
	Birthday time.Time
	From     string `gorm:"size:16"`
	IsDemo   bool   `gorm:"size:16;comment:'创建人';gorm:"default:false""`
}

type PerPel struct {
	gorm.Model
	Key      string `gorm:"index:key,unique"`
	Name     string `gorm:"index:name"`
	Age      int
	Content  dmSchema.Clob `gorm:"size:1024000"`
	Birthday time.Time
}

func TestAutoCreat(t *testing.T) {
	var err error

	if Table().HasTable(&User{}) {
		err := Table().DropTable(&User{})

		if err != nil {
			fmt.Printf("Error: failed to DropTable: %v\n", err)
			return
		}
	}

	err = Table().AutoMigrate(&User{})

	if err != nil {
		fmt.Printf("Error: failed to AutoMigrate: %v\n", err)
		return
	}
}
func TestAutoMigrate(t *testing.T) {
	var err error

	if Table().HasTable(&User{}) {
		err := Table().DropTable(&User{})

		if err != nil {
			fmt.Printf("Error: failed to DropTable: %v\n", err)
			return
		}
	}

	err = Table().AutoMigrate(&User{})
	err = Table().AutoMigrate(&User{})
	err = Table().AutoMigrate(&User{})
	err = Table().AutoMigrate(&PerPel{})
	err = Table().AutoMigrate(&PerPel{})
	err = Table().AutoMigrate(&PerPel{})

	if err != nil {
		fmt.Printf("Error: failed to AutoMigrate: %v\n", err)
		return
	}
}

func TestCreate(t *testing.T) {
	err := Table(&User{Key: "1", Name: "Jinzhu", Age: 18, Content: "asdfdasfasdfasdfj手机卡是点击", Birthday: time.Now()}).Create()
	_ = Table(&User{Key: "2", Name: "Jinzhu1", Age: 19, Content: "bbb", Birthday: time.Now()}).Create()
	_ = Table(&User{Key: "3", Name: "Jinzhu2", Age: 20, Content: "ccc", Birthday: time.Now()}).Create()

	if err != nil {
		fmt.Printf("Error: failed to Create: %v\n", err)
		return
	}
}

func TestGet(t *testing.T) {
	var data User
	err := Table(&User{Name: "Jinzhu"}).Get(&data)

	if err != nil {
		fmt.Printf("Error: failed to Get: %v\n", err)
		return
	}
}

func TestWhere(t *testing.T) {
	var data []User
	err := Table(&User{Name: "Jinzhu"}).GetWhere(&data)

	if err != nil {
		fmt.Printf("Error: failed to Where: %v\n", err)
		return
	}
}

func TestGetAll(t *testing.T) {
	var data []User
	err := Table().GetAll(&data)

	if err != nil {
		fmt.Printf("Error: failed to GetAll: %v\n", err)
		return
	}
}

func TestUpdate(t *testing.T) {
	err := Table(&User{Key: "3"}).Update(&User{Content: "DDDD"})

	if err != nil {
		fmt.Printf("Error: failed to Update: %v\n", err)
		return
	}
}

func TestDelete(t *testing.T) {
	err := Table(&User{Key: "1"}).Delete()

	var data []User
	_ = Table(&User{Key: "1"}).GetWhere(&data)

	if err != nil || len(data) == 1 {
		fmt.Printf("Error: failed to Delete: %v\n", err)
		return
	}
}

// err
func TestClausesAssignmentColumns(t *testing.T) {
	err := Table(&User{Key: "2", Content: "EEE"}).ClausesAssignmentColumns("KEY", []string{"DELETED_AT", "CONTENT"})
	err = Table(&User{Key: "4", Name: "Jinzhu5", Content: "FFF", Birthday: time.Now()}).ClausesAssignmentColumns("KEY", []string{"DELETED_AT", "CONTENT", "BIRTHDAY"})

	if err != nil {
		fmt.Printf("Error: failed to ClausesAssignmentColumns: %v\n", err)
		return
	}
}

type Migration struct {
	ID string `gorm:"primary_key;column:id"`
}

func TestQuery(t *testing.T) {
	var migrationId Migration
	tx := conn.Table("t_db_migrations").Last(&migrationId)
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	t.Log(migrationId.ID)
}

type JSON json.RawMessage

func (v JSON) Value() (driver.Value, error) {
	return []byte(v), nil
}

func (v *JSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch vt := value.(type) {
	case *dmr.DmClob:
		var c dmSchema.Clob
		err := c.Scan(value)
		if err != nil {
			return err
		}
		*v = []byte(c)
		return nil
	case []byte:
		bs := value.([]byte)
		*v = bs
		return nil
	default:
		return fmt.Errorf("invalid type %T, value: %v", vt, value)
	}
}

type TestConfig struct {
	Id      uint   `gorm:"primarykey"`
	Name    string `gorm:"not null;uniqueIndex:idx_cfg_name"`
	Version int64  `gorm:"not null;default:0"`
	Value   JSON   `gorm:"type:clob;"`
}

func (TestConfig) TableName() string {
	return "t_config"
}

func TestQuery2(t *testing.T) {
	rs := &TestConfig{}
	tx := db.Where("name = ?", "task-manager-lock").Find(&rs)
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	t.Log(rs)
}
