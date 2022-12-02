# Squirrel/sqx - Simplified Squirrel SQL expression only

```go
import "github.com/kh411d/squirrel/sqx"

```

**This is a simplified version of Squirrel.** For full implementation, check out
[Masterminds/squirrel](https://github.com/Masterminds/squirrel)

Squirrel helps you build SQL queries from composable parts:

## Select clause
```go
import "github.com/kh411d/squirrel/sqx"

sq := `SELECT * FROM tb_a 
	INNER JOIN tb_b ON tb_a = tb_b `

	eq := sqx.Eq{
		"tb_a.asdf": 3,
		"tb_b.qwer": 2,
	}

	notEq := sqx.NotEq{
		"tb_b.wer": "23423",
	}

	like := sqx.Like{
		"tb_a.rer": "%asdf",
	}

	whereConj := sqx.Where{eq, notEq, like}

	havingConj := sqx.Having{
		Lt{
			"tb_b.nsk": 32,
		},
		Eq{
			"tb_a.abc": "abc",
		},
	}

	sqOrderBy := " ORDER BY tb_a.asdfa ASC "

	sqGroupBy := "GROUP BY tb_a.id"

    // Ordered combined statement
	sql, args, err := sqx.ToSql(
		sq, 
		whereConj,
		havingConj,
		sqOrderBy, 
		sqGroupBy, 
	)

	assert.Equal(t, sql, 
        ` SELECT * FROM tb_a 
        INNER JOIN tb_b ON tb_a = tb_b   WHERE tb_a.asdf = ? AND tb_b.qwer = ? AND tb_b.wer <> ? AND tb_a.rer LIKE ?   HAVING tb_b.nsk < ? AND tb_a.abc = ?   ORDER BY tb_a.asdfa ASC  GROUP BY tb_a.id`,
    )
	assert.Equal(t, args, 
        []interface{}{3, 2, "23423", "%asdf", 32, "abc"},
    )
```

## Insert statement

```go
import "github.com/kh411d/squirrel/sqx"

// Single value
	sq := `INSERT INTO (name, city, country) `

	v := Values{
		{"a", "b", "c"},
	}

    // Ordered combined statement   
	sql, args, err := sqx.ToSql(
        sq, 
        v,
    )

	assert.NoError(t, err)
	assert.Equal(t, sql, 
        " INSERT INTO (name, city, country)  VALUES (?,?,?)"
    )
	assert.Equal(t, args, 
        []interface{}{"a", "b", "c"}
    )

// Multiple insert
	sq := `INSERT INTO (name, city, country) `

	data := sqx.Values{
		{"a", "b", "c"},
		{"aa", "bb", "cc"},
		{"aaa", "bbb", "ccc"},
	}

    // Ordered combined statement
	sql, args, err := sqx.ToSql(
        sq, 
        data,
    )

	assert.Equal(t, sql, 
        " INSERT INTO (name, city, country)  VALUES (?,?,?),(?,?,?),(?,?,?)"
    )

	assert.Equal(t, args, 
        []interface{}{"a", "b", "c", "aa", "bb", "cc", "aaa", "bbb", "ccc"}
    )

```

## FAQ

* **How can I build an IN query on composite keys / tuples, e.g. `WHERE (col1, col2) IN ((1,2),(3,4))`? ([#104](https://github.com/Masterminds/squirrel/issues/104))**

    Squirrel does not explicitly support tuples, but you can get the same effect with e.g.:

    ```go
    sq.Or{
      sq.Eq{"col1": 1, "col2": 2},
      sq.Eq{"col1": 3, "col2": 4}}
    ```

    ```sql
    WHERE (col1 = 1 AND col2 = 2) OR (col1 = 3 AND col2 = 4)
    ```

    (which should produce the same query plan as the tuple version)

* **Why doesn't `Eq{"mynumber": []uint8{1,2,3}}` turn into an `IN` query? ([#114](https://github.com/Masterminds/squirrel/issues/114))**

    Values of type `[]byte` are handled specially by `database/sql`. In Go, [`byte` is just an alias of `uint8`](https://golang.org/pkg/builtin/#byte), so there is no way to distinguish `[]uint8` from `[]byte`.

* **Some features are poorly documented!**

    This isn't a frequent complaints section!

* **Some features are poorly documented?**

    Yes. The tests should be considered a part of the documentation; take a look at those for ideas on how to express more complex queries.

## License

Squirrel is released under the
[MIT License](http://www.opensource.org/licenses/MIT).
