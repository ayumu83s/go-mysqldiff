package mysqldiff

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/tanimutomo/sqlfile"
	"regexp"
	"strings"
	"time"
)

var pkRegexp = regexp.MustCompile("^\\s*PRIMARY KEY\\s+\\((.*)\\)")
var ukRegexp = regexp.MustCompile("^\\s*UNIQUE KEY\\s+`(.*)`\\s+\\((.*)\\)")
var kRegexp = regexp.MustCompile("^\\s*KEY\\s+`(.*)`\\s+\\((.*)\\)")
var columnRegexp = regexp.MustCompile("^\\s*`(.*?)`\\s+(.+?)[\\n,]?$")

type ColumnInfo struct {
	name       string
	definition string
}

type KeyInfo struct {
	name   string
	column string
}

type TableInfo struct {
	TableName  string
	PrimaryKey string
	UniqueKeys []KeyInfo
	Keys       []KeyInfo
	Columns    []ColumnInfo
	Content    string
}

func getTableNames(db *sql.DB) ([]string, error) {
	tableNames := make([]string, 0)

	query := "SHOW TABLES"
	rows, err := db.Query(query)
	if err != nil {
		return tableNames, fmt.Errorf("Error %s: %s", query, err)
	}
	defer rows.Close()

	for rows.Next() {
		var name sql.NullString
		err = rows.Scan(&name)
		if err != nil {
			return tableNames, err
		}
		tableNames = append(tableNames, name.String)
	}
	return tableNames, rows.Err()
}

func GetTables(localConfig, targetConfig Config) ([]TableInfo, error) {

	local := mysql.NewConfig()
	local.User = localConfig.DBUser
	local.Passwd = localConfig.DBPassword
	local.Net = "tcp"
	local.Addr = fmt.Sprintf("%s:%s", localConfig.DBHost, localConfig.DBPort)
	localDb, err := sql.Open("mysql", local.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("connect error: %s@%s:%s", localConfig.DBUser, localConfig.DBHost, localConfig.DBPort)
	}
	defer localDb.Close()

	target := mysql.NewConfig()
	var tables []TableInfo
	var targetDb *sql.DB
	if targetConfig.DBName != "" {
		if targetConfig.DBHost != "" {
			target.User = targetConfig.DBUser
			target.Passwd = targetConfig.DBPassword
			target.DBName = targetConfig.DBName
			target.Net = "tcp"
			target.Addr = fmt.Sprintf("%s:%s", targetConfig.DBHost, targetConfig.DBPort)
		} else {
			target.User = localConfig.DBUser
			target.Passwd = localConfig.DBPassword
			target.DBName = targetConfig.DBName
			target.Net = "tcp"
			target.Addr = fmt.Sprintf("%s:%s", localConfig.DBHost, localConfig.DBPort)
		}

		targetDb, err = sql.Open("mysql", target.FormatDSN())
		if err != nil {
			return nil, fmt.Errorf("connect error: %s", target.FormatDSN())
		}
		defer targetDb.Close()

		tables, err = getTables(targetDb)

	} else if targetConfig.FilePath != "" {
		target.User = localConfig.DBUser
		target.Passwd = localConfig.DBPassword
		target.Net = "tcp"
		target.Addr = fmt.Sprintf("%s:%s", localConfig.DBHost, localConfig.DBPort)

		targetDb, err = sql.Open("mysql", target.FormatDSN())
		if err != nil {
			return nil, fmt.Errorf("connect error: %s", target.FormatDSN())
		}
		defer targetDb.Close()

		tmpDbName := fmt.Sprintf("mysqldiff_tmp_%d", time.Now().UnixNano())
		createTempSchemaByFile(targetDb, tmpDbName, targetConfig.FilePath)
		tables, err = getTables(targetDb)
		cleanupTmpSchema(targetDb, tmpDbName)
	}

	return tables, err
}

func createTempSchemaByFile(db *sql.DB, dbName string, filePath string) {
	query := fmt.Sprintf("CREATE DATABASE `%s`", dbName)
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}

	query = fmt.Sprintf("USE `%s`", dbName)
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}

	s := sqlfile.New()
	s.File(filePath)
	_, err = s.Exec(db)
	if err != nil {
		panic(err)
	}
}

func cleanupTmpSchema(db *sql.DB, dbName string) {
	query := fmt.Sprintf("DROP DATABASE `%s`", dbName)
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func getTables(db *sql.DB) ([]TableInfo, error) {
	tables := make([]TableInfo, 0)

	tableNames, err := getTableNames(db)
	if err != nil {
		return tables, err
	}

	for _, tableName := range tableNames {
		var _tableName, tableSQL sql.NullString
		query := fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)
		if err := db.QueryRow(query).Scan(&_tableName, &tableSQL); err != nil {
			return tables, fmt.Errorf("error %s: %s", query, err)
		}

		primaryKey := ""
		uniqueKeys := make([]KeyInfo, 0)
		keys := make([]KeyInfo, 0)
		columns := make([]ColumnInfo, 0)

		for _, line := range strings.Split(tableSQL.String, "\n") {

			// parse primary key
			pk := pkRegexp.FindAllStringSubmatch(line, -1)
			if len(pk) == 1 && len(pk[0]) == 2 {
				primaryKey = pk[0][1]

			} else {
				// parse unique key
				uk := ukRegexp.FindAllStringSubmatch(line, -1)
				if len(uk) == 1 && len(uk[0]) == 3 {
					uniqueKeys = append(uniqueKeys, KeyInfo{
						name:   uk[0][1],
						column: uk[0][2],
					})

				} else {
					// parse key
					k := kRegexp.FindAllStringSubmatch(line, -1)
					if len(k) == 1 && len(k[0]) == 3 {
						keys = append(keys, KeyInfo{
							name:   k[0][1],
							column: k[0][2],
						})

					} else {
						// parse column
						column := columnRegexp.FindAllStringSubmatch(line, -1)
						if len(column) == 1 && len(column[0]) == 3 {
							columns = append(columns, ColumnInfo{
								name:       column[0][1],
								definition: column[0][2],
							})
						}
					}
				}
			}
		}
		table := TableInfo{
			TableName:  _tableName.String,
			PrimaryKey: primaryKey,
			UniqueKeys: uniqueKeys,
			Keys:       keys,
			Columns:    columns,
			Content:    tableSQL.String,
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func Diff(src, dst []TableInfo) {
	mappedDstInfo := map[string]TableInfo{}
	for _, info := range dst {
		mappedDstInfo[info.TableName] = info
	}

	for _, srcTableInfo := range src {
		tableName := srcTableInfo.TableName
		if _, ok := mappedDstInfo[tableName]; ok {
			diffColumn(tableName, srcTableInfo.Columns, mappedDstInfo[tableName].Columns)
			diffPrimaryKey(tableName, srcTableInfo.PrimaryKey, mappedDstInfo[tableName].PrimaryKey)
			diffKey(tableName, "UNIQUE INDEX", srcTableInfo.UniqueKeys, mappedDstInfo[tableName].UniqueKeys)
			diffKey(tableName, "INDEX", srcTableInfo.Keys, mappedDstInfo[tableName].Keys)
		} else {
			fmt.Printf("%s;\n\n\n", srcTableInfo.Content)
		}
	}

	mappedSrcInfo := map[string]TableInfo{}
	for _, info := range src {
		mappedSrcInfo[info.TableName] = info
	}
	for _, dstTableInfo := range dst {
		dstTableName := dstTableInfo.TableName
		if _, ok := mappedSrcInfo[dstTableName]; !ok {
			fmt.Printf("DROP TABLE `%s`;\n\n\n", dstTableName)
		}
	}
}

func diffColumn(tableName string, src, dst []ColumnInfo) {
	mappedDstInfo := map[string]ColumnInfo{}
	for _, info := range dst {
		mappedDstInfo[info.name] = info
	}

	afterSpecify := ""
	for _, srcInfo := range src {
		if _, ok := mappedDstInfo[srcInfo.name]; ok {
			if srcInfo.definition != mappedDstInfo[srcInfo.name].definition {
				fmt.Printf("ALTER TABLE `%s` MODIFY `%s` %s;\n\n",
					tableName, srcInfo.name, srcInfo.definition,
				)
			}
		} else {
			fmt.Printf("ALTER TABLE `%s` ADD COLUMN `%s` %s%s;\n\n",
				tableName, srcInfo.name, srcInfo.definition, afterSpecify,
			)
		}
		afterSpecify = fmt.Sprintf(" AFTER `%s`", srcInfo.name)
	}

	mappedSrcInfo := map[string]ColumnInfo{}
	for _, info := range src {
		mappedSrcInfo[info.name] = info
	}
	for _, dstInfo := range dst {
		if _, ok := mappedSrcInfo[dstInfo.name]; !ok {
			fmt.Printf("ALTER TABLE `%s` DROP COLUMN `%s`;\n\n",
				tableName, dstInfo.name,
			)
		}
	}
}

func diffPrimaryKey(tableName string, src, dst string) {
	if src != dst {
		if src != "" && dst == "" {
			fmt.Printf("ALTER TABLE `%s` ADD PRIMARY KEY (%s);\n\n",
				tableName, src,
			)
		}
		if src == "" && dst != "" {
			fmt.Printf("ALTER TABLE `%s` DROP PRIMARY KEY;\n\n",
				tableName,
			)
		}
	}
}

func diffKey(tableName string, keyType string, src, dst []KeyInfo) {
	mappedDstInfo := map[string]KeyInfo{}
	for _, info := range dst {
		mappedDstInfo[info.name] = info
	}

	for _, srcInfo := range src {
		if _, ok := mappedDstInfo[srcInfo.name]; !ok {
			fmt.Printf("ALTER TABLE `%s` ADD %s `%s` (%s);\n\n",
				tableName, keyType, srcInfo.name, srcInfo.column,
			)
		}
	}

	mappedSrcInfo := map[string]KeyInfo{}
	for _, info := range src {
		mappedSrcInfo[info.name] = info
	}
	for _, dstInfo := range dst {
		if _, ok := mappedSrcInfo[dstInfo.name]; !ok {
			fmt.Printf("ALTER TABLE `%s` DROP INDEX `%s`;\n\n",
				tableName, dstInfo.name,
			)
		}
	}
}
