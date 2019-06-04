package data

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type (
	queryStmt struct {
		from       string
		whereConds string
		whereArgs  []interface{}
	}

	createStmt struct {
		into   string
		fields []string // value fields
	}

	upsertStmt struct {
		into         string
		keyFields    []string // on conflict fields
		updateFields []string // update fields
	}
)

const (
	insertFmt = "insert into %s (%s) values %s"
	upsertFmt = `insert into %s (%s) values %s
		on conflict (%s) do update set %s`
)

func (s *queryStmt) generate(selectField string) (string, []interface{}) {
	stmt := fmt.Sprintf(`select %s from %s`, selectField, s.from)
	if len(s.whereConds) > 0 {
		numbers := make([]interface{}, len(s.whereArgs))
		for i := range s.whereArgs {
			numbers[i] = i + 1
		}
		stmt = strings.Join([]string{
			stmt,
			fmt.Sprintf(s.whereConds, numbers...),
		}, " where ")
	}
	return stmt, s.whereArgs
}

func (s *queryStmt) query(tx *sql.Tx, selectField string) (*sql.Rows, error) {
	stmt, args := s.generate(selectField)
	sqlLog(stmt, args)
	return tx.Query(stmt, args...)
}

func (s *queryStmt) where(params map[string]interface{}, sep string) *queryStmt {
	whereClause, values := parseParams(params, sep)
	whereClause = trimSpace(whereClause)
	if len(whereClause) > 0 {
		var args []interface{}
		if s.whereArgs == nil {
			args = values
		} else {
			args = append(s.whereArgs, values...)
		}
		if len(s.whereConds) > 0 {
			whereClause = fmt.Sprintf("(%s) and (%s)",
				s.whereConds, whereClause)
		}
		return &queryStmt{
			from:       s.from,
			whereConds: whereClause,
			whereArgs:  args,
		}
	}
	return s
}

func (s *createStmt) modify(tx *sql.Tx, obj interface{}) (sql.Result, error) {
	if s.fields == nil || len(s.fields) == 0 {
		return nil, fmt.Errorf("cannot create without fields")
	}
	if obj == nil {
		return nil, fmt.Errorf(
			"cannot exec nil object when statement requires %d fields",
			len(s.fields))
	}
	valueStmts, args, err := processInsertValues(obj, s.fields)
	if err != nil {
		return nil, err
	}
	stmt := fmt.Sprintf(insertFmt, s.into, commaJoin(s.fields), commaJoin(valueStmts))
	sqlLog(stmt, args)
	return tx.Exec(stmt, args...)
}

func (s *upsertStmt) modify(tx *sql.Tx, obj interface{}) (sql.Result, error) {
	if s.keyFields == nil || len(s.keyFields) == 0 {
		return nil, fmt.Errorf("cannot create without keyFields")
	}
	fields := s.keyFields
	if s.updateFields != nil {
		fields = append(fields, s.updateFields...)
	}
	if obj == nil {
		return nil, fmt.Errorf(
			"cannot exec nil object when statement requires %d fields",
			len(fields))
	}
	valueStmts, args, err := processInsertValues(obj, fields)
	if err != nil {
		return nil, err
	}
	updateStmts := make([]string, len(s.updateFields))
	for i, update := range s.updateFields {
		updateStmts[i] = fmt.Sprintf("%s = excluded.%s", update, update)
	}
	stmt := fmt.Sprintf(upsertFmt, s.into, commaJoin(fields), commaJoin(valueStmts),
		commaJoin(s.keyFields), commaJoin(updateStmts))
	sqlLog(stmt, args)
	return tx.Exec(stmt, args...)
}

func processInsertValues(obj interface{}, fields []string) ([]string, []interface{}, error) {
	var (
		valueStmts []string
		arr        []interface{}
		index      int
	)
	switch v := obj.(type) {
	case []interface{}:
		arr = v
	default:
		arr = []interface{}{v}
	}
	args := make([]interface{}, len(arr)*len(fields))
	for i, entry := range arr {
		values := mapValue(entry)
		stmts := make([]string, len(fields))
		for j, field := range fields {
			val, ok := values[field]
			if !ok {
				return nil, nil, fmt.Errorf("object %v missing field %s",
					obj, field)
			}
			index = i*len(fields) + j
			args[index] = val
			stmts[j] = fmt.Sprintf("$%d", index+1)
		}
		valueStmts = append(valueStmts, "("+commaJoin(stmts)+")")
	}
	return valueStmts, args, nil
}

func mapValue(obj interface{}) map[string]interface{} {
	mappedVal := make(map[string]interface{})
	val := reflect.ValueOf(obj).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if tval, ok := tag.Lookup("gorm"); ok {
			valueField := val.Field(i)
			switch v := valueField.Interface().(type) {
			case []int64, []float64:
				mappedVal[tval] = pq.Array(v)
			default:
				mappedVal[tval] = v
			}
		}
	}
	return mappedVal
}

func parseParams(params map[string]interface{}, sep string) (
	whereClause string, args []interface{}) {
	wheres := make([]string, 0, len(params))
	args = make([]interface{}, 0, len(params))

	addArgs := func(key string, v []interface{}) {
		sarr := make([]string, len(v))
		for i := range v {
			sarr[i] = "$%d"
		}
		wheres = append(wheres, fmt.Sprintf("%s in (%s)",
			key, strings.Join(sarr, ",")))
		args = append(args, v...)
	}

	for key, value := range params {
		switch v := value.(type) {
		case []int:
			iarr := make([]interface{}, len(v))
			for i, e := range v {
				iarr[i] = e
			}
			addArgs(key, iarr)
		case []string:
			iarr := make([]interface{}, len(v))
			for i, e := range v {
				iarr[i] = e
			}
			addArgs(key, iarr)
		default:
			wheres = append(wheres, fmt.Sprintf("%s = $%%d", key))
			args = append(args, value)
		}
	}
	whereClause = strings.Join(wheres, fmt.Sprintf(" %s ", sep))
	return
}

func sqlLog(queryStmt string, args []interface{}) {
	queryStmt = trimSpace(queryStmt)
	if len(args) <= 0 {
		log.Infof("[sql]: %s", queryStmt)
	} else {
		strArgs := make([]string, len(args))
		for i, arg := range args {
			strArgs[i] = fmt.Sprintf("%d='%+v'", i+1, arg)
		}
		log.Infof("[sql]: %s, (%s)", queryStmt, commaJoin(strArgs))
	}
}

func trimSpace(s string) string {
	return strings.Trim(s, " \t\n")
}

func commaJoin(sarr []string) string {
	return strings.Join(sarr, ",")
}
