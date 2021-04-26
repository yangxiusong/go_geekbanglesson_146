package ch02

import (
	"fmt"
	xerrors "github.com/pkg/errors"
)

func DoFind(id int32) {
	InitDb()
	ret := &Person{}
	dao := &Dao{}
	query := "SELECT * FROM person WHERE id=?"
	dao.Find(query, ret, id)
	if dao.err != nil {
		fmt.Printf("original error:%T %v\n", xerrors.Cause(dao.err), xerrors.Cause(dao.err))
		fmt.Printf("stack trace:\n%+v\n", dao.err)
		return
	}
	fmt.Println(dao.data)
}

func DoFind2(id int32) {
	InitDb()
	ret := Person{}
	dao := &Dao{}
	query := "SELECT * FROM person WHERE id=?"
	dao.FindPerson(query, &ret, id)
	if dao.err != nil {
		fmt.Printf("original error:%T %v\n", xerrors.Cause(dao.err), xerrors.Cause(dao.err))
		fmt.Printf("stack trace:\n%+v\n", dao.err)
		return
	}
	fmt.Println("dao",dao)
	fmt.Println("Persondata:", ret)
}

func DoList() {
	InitDb()
	dao := &Dao{}
	query := "SELECT * FROM person"
	dao.List(query)
	if dao.err != nil {
		fmt.Printf("original error:%T %v\n", xerrors.Cause(dao.err), xerrors.Cause(dao.err))
		fmt.Printf("stack trace:\n%+v\n", dao.err)
		return
	}
	fmt.Println("dao",dao)
}