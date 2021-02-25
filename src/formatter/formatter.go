package formatter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/murata100/hakagi/src/constraint"
)

const (
	baseSql = "ALTER TABLE %s ADD CONSTRAINT FOREIGN KEY (%s) REFERENCES %s(%s);"
)

type Refered struct {
	referedTable, referedColumn string
}

func FormatSql(constraints []constraint.Constraint) string {
	var queries []string

	for _, c := range constraints {
		q := fmt.Sprintf(baseSql, c.Table, c.Column, c.ReferedTable, c.ReferedColumn)
		queries = append(queries, q)
	}

	return strings.Join(queries, "\n")
}

func FormatXML(constraints []constraint.Constraint) string {
	var x []string

	x = append(x, "<schemaMeta xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" "+
		"xsi:noNamespaceSchemaLocation=\"http://schemaspy.org/xsd/6/schemameta.xsd\" >")
	x = append(x, "    <tables>")

	m := map[string]map[string][]*Refered{}
	for _, c := range constraints {
		_, ok := m[c.Table]
		if !ok {
			m[c.Table] = map[string][]*Refered{}
		}
		_, ok = m[c.Table][c.Column]
		if !ok {
			m[c.Table][c.Column] = []*Refered{}
		}
		m[c.Table][c.Column] = append(m[c.Table][c.Column],
			&Refered{referedTable: c.ReferedTable, referedColumn: c.ReferedColumn})
	}

	var mKeys []string
	for k := range m {
		mKeys = append(mKeys, k)
	}
	sort.Strings(mKeys)

	for _, table := range mKeys {
		columns := m[table]

		x = append(x, fmt.Sprintf("        <table name=\"%s\">", table))

		var columnsKeys []string
		for k := range columns {
			columnsKeys = append(columnsKeys, k)
		}
		sort.Strings(columnsKeys)

		for _, column := range columnsKeys {
			refereds := columns[column]

			x = append(x, fmt.Sprintf("            <column name=\"%s\">", column))

			for _, refered := range refereds {
				x = append(x, fmt.Sprintf(
					"                <foreignKey table=\"%s\" column=\"%s\" />",
					refered.referedTable, refered.referedColumn))
			}
			x = append(x, fmt.Sprintf("            </column>"))
		}
		x = append(x, "        </table>")
	}

	x = append(x, "    </tables>")
	x = append(x, "</schemaMeta>")

	return strings.Join(x, "\n")
}
