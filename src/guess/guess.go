package guess

import (
	"strings"

	"github.com/jinzhu/inflection"
	"github.com/murata100/hakagi/src/constraint"
	"github.com/murata100/hakagi/src/database"
)

const (
	idColumn           = "id"
	targetColumnSuffix = "_id"
)

type GuessOption func(database.Column, string, database.Column, []string) bool

func isAcceptableAsIndex(left, right string, compatibleTypes []string) bool {
	is := left == right &&
		!(strings.Index(left, "text") != -1 || strings.Index(left, "blob") != -1) &&
		!(strings.Index(right, "text") != -1 || strings.Index(right, "blob") != -1)
	if left == compatibleTypes[0] && right == compatibleTypes[1] {
		is = true
	}
	return is
}

// Recongnize a column thats same name of other table's primary key is a foreign key
// This base idea refers to SchemaSpy DbAnalyzer:
//   https://github.com/schemaspy/schemaspy/blob/master/src/main/java/org/schemaspy/DbAnalyzer.java
func GuessByPrimaryKey() GuessOption {
	return func(i database.Column, table string, pk database.Column, compatibleTypes []string) bool {
		return isAcceptableAsIndex(i.Type, pk.Type, compatibleTypes) && i.Name == pk.Name && pk.Name != idColumn
	}
}

// Recongnize a column thats same name without '_id' suffix of other table  name is a foreign key
func GuessByTableAndColumn() GuessOption {
	return func(i database.Column, table string, pk database.Column, compatibleTypes []string) bool {
		if !isAcceptableAsIndex(i.Type, pk.Type, compatibleTypes) {
			return false
		}

		cLen := len(i.Name)
		tLen := len(targetColumnSuffix)
		if !(cLen >= tLen && i.Name[cLen-tLen:] == targetColumnSuffix) {
			return false
		}

		return inflection.Plural(i.Name[:cLen-tLen]) == table && pk.Name == idColumn
	}
}

// GuessConstraints guesses foreign key constraints from primary keys and indexes.
// NOTE composite primary keys are not supported.
func GuessConstraints(indexes database.Indexes, primaryKeys database.PrimaryKeys,
	compatibleTypes []string, guessOptions ...GuessOption) []constraint.Constraint {
	var constraints []constraint.Constraint

	for indexTable, indexMaps := range indexes {
		for _, indexCols := range indexMaps {
			for pkTable, pk := range primaryKeys {
				if indexTable != pkTable && len(indexCols) == 1 && len(pk) == 1 {
					singleIndex := indexCols[0]
					singlePk := pk[0]

					for _, guesser := range guessOptions {
						if guesser(singleIndex, pkTable, singlePk, compatibleTypes) {
							constraints = append(constraints,
								constraint.Constraint{indexTable, singleIndex.Name, pkTable, singlePk.Name})
						}
						/* else {
							fmt.Println(indexTable, singleIndex, pkTable, singlePk)
						} */
					}
				}
			}
		}
	}

	return constraints
}
