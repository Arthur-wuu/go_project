// Copyright 2018 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	sql, args, err := Select("c, d").From("table1").LeftJoin("table2", Eq{"table1.id": 1}.And(Lt{"table2.id": 3})).
		RightJoin("table3", "table2.id = table3.tid").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c, d FROM table1 LEFT JOIN table2 ON table1.id=? AND table2.id<? RIGHT JOIN table3 ON table2.id = table3.tid WHERE a=?",
		sql)
	assert.EqualValues(t, []interface{}{1, 3, 1}, args)

	sql, args, err = Select("c, d").From("table1").LeftJoin("table2", Eq{"table1.id": 1}.And(Lt{"table2.id": 3})).
		FullJoin("table3", "table2.id = table3.tid").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c, d FROM table1 LEFT JOIN table2 ON table1.id=? AND table2.id<? FULL JOIN table3 ON table2.id = table3.tid WHERE a=?",
		sql)
	assert.EqualValues(t, []interface{}{1, 3, 1}, args)

	sql, args, err = Select("c, d").From("table1").LeftJoin("table2", Eq{"table1.id": 1}.And(Lt{"table2.id": 3})).
		CrossJoin("table3", "table2.id = table3.tid").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c, d FROM table1 LEFT JOIN table2 ON table1.id=? AND table2.id<? CROSS JOIN table3 ON table2.id = table3.tid WHERE a=?",
		sql)
	assert.EqualValues(t, []interface{}{1, 3, 1}, args)

	sql, args, err = Select("c, d").From("table1").LeftJoin("table2", Eq{"table1.id": 1}.And(Lt{"table2.id": 3})).
		InnerJoin("table3", "table2.id = table3.tid").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c, d FROM table1 LEFT JOIN table2 ON table1.id=? AND table2.id<? INNER JOIN table3 ON table2.id = table3.tid WHERE a=?",
		sql)
	assert.EqualValues(t, []interface{}{1, 3, 1}, args)

	subQuery2 := Select("e").From("table2").Where(Gt{"e": 1})
	subQuery3 := Select("f").From("table3").Where(Gt{"f": "2"})
	sql, args, err = Select("c, d").From("table1").LeftJoin(subQuery2, Eq{"table1.id": 1}.And(Lt{"table2.id": 3})).
		InnerJoin(subQuery3, "table2.id = table3.tid").Where(Eq{"a": 1}).ToSQL()
	assert.NoError(t, err)
	assert.EqualValues(t, "SELECT c, d FROM table1 LEFT JOIN (SELECT e FROM table2 WHERE e>?) ON table1.id=? AND table2.id<? INNER JOIN (SELECT f FROM table3 WHERE f>?) ON table2.id = table3.tid WHERE a=?",
		sql)
	assert.EqualValues(t, []interface{}{1, 1, 3, "2", 1}, args)
}
