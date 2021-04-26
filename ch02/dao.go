package ch02

import (
	"database/sql"
	xerrors "github.com/pkg/errors"
	"log"
	"reflect"
	"strconv"
)

type Dao struct {
	data []map[string]string
	err  error
}

func (d *Dao) mapping(m map[string]string, v reflect.Value) error {
	t := v.Type()
	val := v.Elem()
	typ := t.Elem()
	if !val.IsValid() {
		return xerrors.New("数据类型不正确")
	}

	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		kind := value.Kind()
		tag := typ.Field(i).Tag.Get("col")
		if len(tag) == 0 {
			continue
		}
		meta, ok := m[tag]
		if !ok {
			continue
		}
		if !value.CanSet() {
			return xerrors.New("结构体字段没有读写权限")
		}
		if len(meta) == 0 {
			continue
		}
		if kind == reflect.String {
			value.SetString(meta)
		} else if kind == reflect.Float32 {
			f, err := strconv.ParseFloat(meta, 32)
			if err != nil {
				return xerrors.Wrapf(err, "strconv.ParseFloat failed")
			}
			value.SetFloat(f)
		} else if kind == reflect.Int {
			iv, err := strconv.Atoi(meta)
			if err != nil {
				return xerrors.Wrapf(err, "strconv.Atoi failed")
			}
			value.SetInt(int64(iv))
		} else if kind == reflect.Bool {
			b, err := strconv.ParseBool(meta)
			if err != nil {
				return xerrors.Wrapf(err, "strconv.ParseBool failed")
			}
			value.SetBool(b)
		} else {
			return xerrors.New("数据库映射存在不识别的数据类型")
		}
	}
	return nil
}

func (d *Dao) Unique(in interface{}) {
	if len(d.data) == 0 {
		return
	}
	d.err = d.mapping(d.data[0], reflect.ValueOf(in))
}

func (d *Dao) query(rows *sql.Rows) error {
	col, err := rows.Columns()
	if err != nil {
		return xerrors.Wrap(err, "rows.Columns failed")
	}
	vals := make([][]byte, len(col))
	scans := make([]interface{}, len(col))

	for i := range vals {
		scans[i] = &vals[i]
	}

	results := make([]map[string]string, 0)
	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			return xerrors.Wrap(err, "rows.Scan failed")
		}
		row := make(map[string]string)
		for k, v := range vals {
			key := col[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}
	d.data = results
	return nil
}

/**
 可以俘获sql.ErrNoRows
*/
func (d *Dao) FindPerson(q string, ret *Person, args ...interface{}) *Dao {
	err := db.QueryRow(q, args...).Scan(&ret.Id, &ret.FirstName, &ret.LastName)
	if err != nil {
		if err == sql.ErrNoRows {
			d.err = xerrors.Wrap(err, "sql.ErrNoRows error")
		} else {
			d.err = xerrors.Wrap(err, "FindPerson error")
		}
		return d
	}
	return d
}

func (d *Dao) Find(sql string, ret interface{}, args ...interface{}) *Dao {
	rows, err := db.Query(sql, args...)

	log.Println("find-", sql, args, "err:", err)

	if err != nil {
		d.err = err
		return d
	}
	defer rows.Close()
	err = d.query(rows)
	if err != nil {
		d.err = err
		return d
	}

	d.Unique(ret)
	return d
}

func (d *Dao) List(q string)  {
	rows, err := db.Query(q)

	if err != nil {
		d.err = err
		return
	}
	defer rows.Close()
	err = d.query(rows)
	if err != nil {
		d.err = err
		return
	}
	return
}
