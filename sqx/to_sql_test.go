package sqx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToSql(t *testing.T) {
	sq := `SELECT * FROM tb_a 
	INNER JOIN tb_b ON tb_a = tb_b 
	WHERE `

	eq := Eq{
		"tb_a.asdf": 3,
		"tb_b.qwer": 2,
	}

	notEq := NotEq{
		"tb_b.wer": "23423",
	}

	like := Like{
		"tb_a.rer": "%asdf",
	}

	lt := Lt{
		"tb_b.nsk": 32,
	}

	sqOrderBy := " ORDER BY tb_a.asdfa ASC "

	sql, args, err := ToSql(
		sq,
		eq,
		notEq,
		like,
		lt,
		sqOrderBy,
	)
	assert.NoError(t, err)
	assert.Equal(t, sql, ` SELECT * FROM tb_a 
	INNER JOIN tb_b ON tb_a = tb_b 
	WHERE  tb_a.asdf = ? AND tb_b.qwer = ? tb_b.wer <> ? tb_a.rer LIKE ? tb_b.nsk < ?  ORDER BY tb_a.asdfa ASC `)
	assert.Equal(t, args, []interface{}([]interface{}{3, 2, "23423", "%asdf", 32}))

}

func TestToSqlWithConj(t *testing.T) {
	sq := `SELECT * FROM tb_a 
	INNER JOIN tb_b ON tb_a = tb_b `

	eq := Eq{
		"tb_a.asdf": 3,
		"tb_b.qwer": 2,
	}

	notEq := NotEq{
		"tb_b.wer": "23423",
	}

	like := Like{
		"tb_a.rer": "%asdf",
	}

	whereConj := Where{eq, notEq, like}

	havingConj := Having{
		Lt{
			"tb_b.nsk": 32,
		},
		Eq{
			"tb_a.abc": "abc",
		},
	}

	sqOrderBy := " ORDER BY tb_a.asdfa ASC "

	sqGroupBy := "GROUP BY tb_a.id"

	sql, args, err := ToSql(
		sq,
		whereConj,
		havingConj,
		sqOrderBy,
		sqGroupBy,
	)

	assert.NoError(t, err)
	assert.Equal(t, sql, ` SELECT * FROM tb_a 
	INNER JOIN tb_b ON tb_a = tb_b   WHERE tb_a.asdf = ? AND tb_b.qwer = ? AND tb_b.wer <> ? AND tb_a.rer LIKE ?   HAVING tb_b.nsk < ? AND tb_a.abc = ?   ORDER BY tb_a.asdfa ASC  GROUP BY tb_a.id`)
	assert.Equal(t, args, []interface{}{3, 2, "23423", "%asdf", 32, "abc"})

}

func TestToSqlUpdate(t *testing.T) {
	sq := `UPDATE tb_a SET `
	eq := Eq{
		"tb_a.name": "kambing",
		"tb_a.col":  1,
	}

	sql, args, err := ToSql(sq, eq)
	assert.NoError(t, err)
	assert.Equal(t, sql, " UPDATE tb_a SET  tb_a.col = ? AND tb_a.name = ?")
	assert.Equal(t, args, []interface{}{1, "kambing"})
}

func TestToSqlInsert(t *testing.T) {
	sq := `INSERT INTO (name, city, country) `

	v := Values{
		{"a", "b", "c"},
	}

	sql, args, err := ToSql(Expr(sq), v)

	assert.NoError(t, err)
	assert.Equal(t, sql, " INSERT INTO (name, city, country)  VALUES (?,?,?)")
	assert.Equal(t, args, []interface{}{"a", "b", "c"})
}

func TestToSqlInsertMultiple(t *testing.T) {
	sq := `INSERT INTO (name, city, country) `

	data := Values{
		{"a", "b", "c"},
		{"aa", "bb", "cc"},
		{"aaa", "bbb", "ccc"},
	}

	sql, args, err := ToSql(Expr(sq), data)

	assert.NoError(t, err)
	assert.Equal(t, sql, " INSERT INTO (name, city, country)  VALUES (?,?,?),(?,?,?),(?,?,?)")
	assert.Equal(t, args, []interface{}{"a", "b", "c", "aa", "bb", "cc", "aaa", "bbb", "ccc"})
}
