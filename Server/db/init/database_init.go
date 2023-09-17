package db_init

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

type User struct {
	Id   string `bun:"type:varchar(255),primary" json:"id"` // Note the change here from 'primary_key' to 'primary'
	Name string `bun:"type:varchar(255)" json:"name"`
}

func CreateTable(db *bun.DB) {
	_, err := db.NewCreateTable().Model((*User)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Println("createTable")
}
