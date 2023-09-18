package clause

import "strings"

// AtomType sql拆分后的子句类型
type AtomType int

const (
	INSERT = iota
	DELETE
	UPDATE
	SELECT
	WHERE
	LIMIT
	ORDER
	VALUE
	COUNT
)

// Clause sql子句
type Clause struct {
	sql     map[AtomType]string
	sqlVars map[AtomType][]interface{}
}

func (c *Clause) Set(typ AtomType, vars ...interface{}) {
	if c.sqlVars == nil || c.sql == nil {
		c.sql = make(map[AtomType]string)
		c.sqlVars = make(map[AtomType][]interface{})
	}
	sql, sqlVars := generatorMap[typ](vars...)
	c.sql[typ] = sql
	c.sqlVars[typ] = sqlVars
}

func (c *Clause) Build(names ...AtomType) (string, []interface{}) {
	var resSqlVars []interface{}
	var sqlArr []string

	for _, name := range names {
		if _, ok := c.sql[name]; !ok {
			continue
		}
		sql := c.sql[name]
		sqlVars := c.sqlVars[name]
		sqlArr = append(sqlArr, sql)
		resSqlVars = append(resSqlVars, sqlVars...)
	}

	resSql := strings.Join(sqlArr, " ")
	return resSql, resSqlVars
}
