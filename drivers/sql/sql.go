package sql

import (
	"context"
	"errors"
	"fmt"
	"math"
	reflect "reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
)

const (
	// Boolean is a logical Boolean (true/false)
	Boolean Datatype = "boolean"
	// Text is a variable length character string
	Text Datatype = "text"
	// Integer is a signed 4byte integer
	Integer Datatype = "integer"
	// Bigint is a signed 8byte integer
	Bigint Datatype = "bigint"
	// Real is a single precision floating point number
	Real Datatype = "real"
	// Double is a double precision floating point number
	Double Datatype = "double precision"
	// Timestamp is date and time without timezone
	Timestamp Datatype = "timestamp without time zone"
	// Blob is a binary data ("byte array")
	Blob Datatype = "bytea"
	// JSONB is a json document
	JSONB Datatype = "jsonb"
	//UUID is unique id
	UUID Datatype = "uuid"

	timeFormat = "2006-01-02 15:04:05.000000"
)

// Datatype describes a database data type
type Datatype string

// Driver defines an SQL Driver
type Driver interface {
	// GENERAL PURPOSE APIs
	BuildInsertQuery(ctx context.Context, rows Rows, fields ...Field) string
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryTx(ctx context.Context, transaction *Transaction, query string, args ...interface{}) (Rows, error)
	BatchQuery(ctx context.Context, queries []string) ([]Rows, error)
	BatchQueryTx(ctx context.Context, transaction *Transaction, queries []string) ([]Rows, error)
	//Exec executes a provided SQL query and will return the number of rows effected or error if any
	Exec(ctx context.Context, query string, args ...interface{}) (int64, error)
	BatchExec(ctx context.Context, queries []string, args [][]interface{}) error
	//DeleteByID  is used to delete a row from the table by id
	//and will return the number of rows effected or error if any
	DeleteByID(ctx context.Context, tablename string, value int) (int64, error)
	//DeleteQuery it's creates a delete query for deleting a row from a table
	//and will return the number of rows effected or error if any
	Delete(ctx context.Context, tablename, columnname, value string) (int64, error)
	//Update - It creates an update query based on variadic arguments
	//and will return the number of rows effected or error if any
	Update(ctx context.Context, args ...string) (int64, error)
	GetTablesName(ctx context.Context) ([]string, error)
	GetFields(ctx context.Context, tableName string) ([]Field, error)
	//CreateTable create a new table into the DB, if isTimeSeries is set to true a column
	//If the CommandTag was not for a row affecting command (such as "CREATE TABLE") then it returns 0
	CreateTable(ctx context.Context, rows Rows, isTimeSeries bool, constrains ...string) (int64, error)
	//CreateIndex it create an index on the specified table or error if any
	CreateIndex(ctx context.Context, indexName, tableName string, FieldsName []string) (int64, error)
	//DeleteValuesFromInterval delete values from a particular table older than the number of days specified
	//and will return the number of rows effected or error if any
	DeleteValuesFromInterval(ctx context.Context, tableName string, days int32) (int64, error)
	Insert(ctx context.Context, rows Rows, fields ...Field) (Rows, error)
	Unmarshal(source []interface{}, destination ...interface{}) error
	ExecTx(ctx context.Context, tx *Transaction, query string, args ...interface{}) error
	InsertTx(ctx context.Context, tx *Transaction, rows Rows, fields ...Field) (Rows, error)
	CreateTransaction() (*Transaction, error)
	Rollback(tx *Transaction) error
	Commit(tx *Transaction) error
	//Upsert upserts rows into the database
	//and will return the number of rows effected or error if any
	Upsert(ctx context.Context, rows Rows, constrain string) (int64, error)
	Copy(rows Rows) error
	Close()
	Reset()
	Stat() ConnPoolStat

	// TIMESERIES APIs
	GetLatestValue(ctx context.Context, fieldsName []string, tableName string) ([]Rows, error)
	GetTimeSeriesValues(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time) ([]Rows, error)
	GetDownSampledValues(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, seconds int) ([]Rows, error)
	GetInterpolatedValues(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, seconds int) ([]Rows, error)

	// HISTORIC APIs
	GetHistorianTimeSeriesValues(ctx context.Context, fieldsName string, tableName string, startTime, endTime time.Time, desc bool, where HistoricWhere, limit uint32) (Rows, error)
	GetHistorianAggrMax(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, desc bool, limit uint32, interval int64) ([]Rows, error)
	GetHistorianAggrMin(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, desc bool, limit uint32, interval int64) ([]Rows, error)
	GetHistorianAggrCount(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, desc bool, limit uint32, interval int64) ([]Rows, error)
	GetHistorianAggrAvg(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, desc bool, limit uint32, interval int64) ([]Rows, error)
}

type pgxConnPool interface {
	QueryEx(ctx context.Context, query string, opts *pgx.QueryExOptions, args ...interface{}) (*pgx.Rows, error)
	ExecEx(ctx context.Context, query string, opts *pgx.QueryExOptions, args ...interface{}) (pgx.CommandTag, error)
	BeginBatch() *pgx.Batch
	CopyFrom(tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int, error)
	Close()
	Begin() (*pgx.Tx, error)
	Reset()
	Stat() pgx.ConnPoolStat
}

// ConnPoolStat represents stats for connection pool
type ConnPoolStat struct {
	MaxConnections       int // max simultaneous connections to use
	CurrentConnections   int // current live connections
	AvailableConnections int // unused live connections
}
type pgxDriver struct {
	cp pgxConnPool
}

//Transaction is a transaction
type Transaction struct {
	tx  *pgx.Tx
	err error
}

// CreatePgxConnPool defines the functions to create a new connection pool
type CreatePgxConnPool func(pgx.ConnPoolConfig) (*pgx.ConnPool, error)

// Rows is a set of rows
type Rows struct {
	SchemaName string
	TableName  string
	Fields     []Field
	Values     [][]interface{}
}

// Field is a database field
type Field struct {
	Datatype
	Name string
}

// NewPgxDriver returns a pgx driver, createPoolFunc is a function to create a new connection pool, if it is
// nil the deafult pgx function will be used. connString is the DB connection string
// i.e. host={hostName} port={portNumber} user={userName} password=%s dbname=%s
// maxConn is the max number of connections to be used for the pool
// acqTimeOut is the timeout used in acquiring a connection from the pool
func NewPgxDriver(createPoolFunc CreatePgxConnPool, connString string, maxConn, acqTimeOut int) (Driver, error) {
	if createPoolFunc == nil {
		createPoolFunc = pgx.NewConnPool
	}

	connCfg, e := pgx.ParseDSN(connString)
	if e != nil {
		return nil, e
	}

	connPool, e := createPoolFunc(pgx.ConnPoolConfig{
		ConnConfig:     connCfg,
		MaxConnections: maxConn,
		AcquireTimeout: time.Duration(acqTimeOut) * time.Second,
	})
	if e != nil {
		return nil, e
	}
	return pgxDriver{
		cp: connPool,
	}, nil
}

func (d pgxDriver) Close() {
	d.cp.Close()
}

// CreatePostgresConnection creates connection to postgres database.
func CreatePostgresConnection(dbHost, dbPort, dbUser, dbPassword, dbName string, enableDbConnReset bool, dbResetTime time.Duration, maxDbConnAllowed int) (Driver, error) {
	postgresqlConn := "host=%s port=%s user=%s password=%s dbname=%s"
	connStr := fmt.Sprintf(postgresqlConn, dbHost, dbPort, dbUser, dbPassword, dbName)
	pgxDriver, err := NewPgxDriver(nil, connStr, maxDbConnAllowed, 30)
	if err != nil {
		return nil, err
	}
	if enableDbConnReset {
		go func(pgxDriver Driver) {
			resetTick := time.NewTicker(dbResetTime * time.Minute)
			for {
				select {
				case <-resetTick.C:
					pgxDriver.Reset()
				}
			}
		}(pgxDriver)
	}
	return pgxDriver, nil
}

// Reset closes all open connections, but leaves the pool open. It is intended
// for use when an error is detected that would disrupt all connections (such as
// a network interruption or a server state change).
//
// It is safe to reset a pool while connections are checked out. Those
// connections will be closed when they are returned to the pool.
func (d pgxDriver) Reset() {
	d.cp.Reset()
}

//Begin begins a transaction
func (d pgxDriver) CreateTransaction() (*Transaction, error) {
	tx, err := d.cp.Begin()
	transaction := &Transaction{
		tx:  tx,
		err: err,
	}
	return transaction, nil
}

//Rollback rolls back a transaction
func (d pgxDriver) Rollback(transaction *Transaction) error {
	return transaction.tx.Rollback()
}

//Commit commits a transaction
func (d pgxDriver) Commit(transaction *Transaction) error {
	return transaction.tx.Commit()
}

//ExecTx executes a provided SQL query by transaction
func (d pgxDriver) ExecTx(ctx context.Context, transaction *Transaction, query string, args ...interface{}) error {

	_, e := transaction.tx.ExecEx(ctx, query, nil, args...)
	if e != nil {
		return e
	}

	return nil
}

// QueryTx executes query on the database as a result of a specific query by transaction
func (d pgxDriver) QueryTx(ctx context.Context, transaction *Transaction, query string, args ...interface{}) (Rows, error) {
	rows := Rows{
		Fields: []Field{},
		Values: [][]interface{}{},
	}

	pgxRows, e := transaction.tx.QueryEx(ctx, query, nil, args...)
	if e != nil {
		return rows, e
	}
	defer pgxRows.Close()
	return d.extractRowsFromPgxRows(ctx, pgxRows)
}

//Unmarshal scans the row values into destination fields
func (d pgxDriver) Unmarshal(source []interface{}, dest ...interface{}) error {
	if len(source) != len(dest) {
		return errors.New("source and destination doesn't match")
	}
	for i, s := range source {
		so := reflect.ValueOf(s)
		if !reflect.DeepEqual(so, reflect.Zero(reflect.TypeOf(so)).Interface()) {
			reflect.ValueOf(dest[i]).Elem().Set(so)
		}
	}
	return nil
}

func (d pgxDriver) BuildInsertQuery(ctx context.Context, rows Rows, fields ...Field) string {
	var sb strings.Builder
	if rows.SchemaName == "" {
		rows.SchemaName = "public"
	}
	sb.WriteString(rows.SchemaName)
	sb.WriteString(".")
	sb.WriteString("\"")
	sb.WriteString(rows.TableName)
	sb.WriteString("\"")

	query := fmt.Sprintf(insert,
		sb.String(),
		buildFieldsString(rows.Fields),
		buildValuesString(rows))
	if len(fields) > 0 {
		query += " returning " + buildReturningClause(fields)
	}
	return query
}

// InsertTx inserts rows into the database and returns fields
func (d pgxDriver) InsertTx(ctx context.Context, tx *Transaction, rows Rows, fields ...Field) (Rows, error) {
	var sb strings.Builder
	r := Rows{}
	if rows.SchemaName == "" {
		rows.SchemaName = "public"
	}
	sb.WriteString(rows.SchemaName)
	sb.WriteString(".")
	sb.WriteString("\"")
	sb.WriteString(rows.TableName)
	sb.WriteString("\"")

	var query string
	if len(fields) == 0 {
		query = fmt.Sprintf(insert,
			sb.String(),
			buildFieldsString(rows.Fields),
			buildValuesString(rows))
		return r, d.ExecTx(ctx, tx, query)
	}
	query = fmt.Sprintf(insert,
		sb.String(),
		buildFieldsString(rows.Fields),
		buildValuesString(rows)+" returning "+buildFieldsString(fields),
	)
	return d.QueryTx(ctx, tx, query)
}

// BatchQueryTx retrieves rows from the database as a result of a batch query
func (d pgxDriver) BatchQueryTx(ctx context.Context, tx *Transaction, queries []string) ([]Rows, error) {
	var rows []Rows

	batch := tx.tx.BeginBatch()

	for _, query := range queries {
		batch.Queue(query, nil, nil, nil)
	}

	if e := batch.Send(ctx, nil); e != nil {
		return rows, e
	}

	for range queries {
		queryRow := Rows{
			Fields: []Field{},
			Values: [][]interface{}{},
		}

		r, e := batch.QueryResults()
		if e != nil {
			return rows, e
		}

		for r.Next() {
			values, e := r.Values()
			if e != nil {
				return rows, e
			}

			queryRow.Values = append(queryRow.Values, values)
		}

		for _, fd := range r.FieldDescriptions() {
			tn, e := d.getTableFromOID(ctx, fd.Table)
			if e != nil {
				return rows, e
			}
			queryRow.TableName = tn
			queryRow.Fields = append(queryRow.Fields,
				Field{
					Name:     fd.Name,
					Datatype: Datatype(fd.DataTypeName),
				})
		}

		rows = append(rows, queryRow)
	}

	if e := batch.Close(); e != nil {
		return rows, e
	}

	return rows, nil
}

// Query retrieves rows from the database as a result of a specific query
func (d pgxDriver) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows := Rows{
		Fields: []Field{},
		Values: [][]interface{}{},
	}
	pgxRows, e := d.cp.QueryEx(ctx, query, nil, args...)
	if e != nil {
		return rows, e
	}
	defer pgxRows.Close()
	return d.extractRowsFromPgxRows(ctx, pgxRows)

}

// BatchQuery retrieves rows from the database as a result of a batch query
func (d pgxDriver) BatchQuery(ctx context.Context, queries []string) ([]Rows, error) {
	var rows []Rows

	batch := d.cp.BeginBatch()

	for _, query := range queries {
		batch.Queue(query, nil, nil, nil)
	}

	if e := batch.Send(ctx, nil); e != nil {
		return rows, e
	}

	for range queries {
		queryRow := Rows{
			Fields: []Field{},
			Values: [][]interface{}{},
		}

		r, e := batch.QueryResults()
		if e != nil {
			return rows, e
		}

		for r.Next() {
			values, e := r.Values()
			if e != nil {
				return rows, e
			}

			queryRow.Values = append(queryRow.Values, values)
		}

		for _, fd := range r.FieldDescriptions() {
			tn, e := d.getTableFromOID(ctx, fd.Table)
			if e != nil {
				return rows, e
			}
			queryRow.TableName = tn
			queryRow.Fields = append(queryRow.Fields,
				Field{
					Name:     fd.Name,
					Datatype: Datatype(fd.DataTypeName),
				})
		}

		rows = append(rows, queryRow)
	}

	if e := batch.Close(); e != nil {
		return rows, e
	}

	return rows, nil
}

// getTableFromOID returns table's name based on its OID
func (d pgxDriver) getTableFromOID(ctx context.Context, oid pgtype.OID) (string, error) {
	var tableName string
	pgxRows, e := d.cp.QueryEx(ctx, fmt.Sprintf(getTableFromOID, oid), nil)
	if e != nil {
		return "", e
	}

	defer pgxRows.Close()

	for pgxRows.Next() {
		if e := pgxRows.Scan(&tableName); e != nil {
			return "", e
		}
	}

	return tableName, nil
}

// Exec executes a provided SQL query
func (d pgxDriver) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	var defaultIntVal int64
	CommandTag, e := d.cp.ExecEx(ctx, query, nil, args...)
	if e != nil {
		return defaultIntVal, e
	}

	return CommandTag.RowsAffected(), nil
}

// BatchExec executes a provided SQL query
func (d pgxDriver) BatchExec(ctx context.Context, queries []string, args [][]interface{}) error {
	batch := d.cp.BeginBatch()

	for i, query := range queries {
		if args != nil && len(args[i]) > 0 {
			batch.Queue(query, args[i], nil, nil)
		} else {
			batch.Queue(query, nil, nil, nil)
		}
	}

	if e := batch.Send(ctx, nil); e != nil {
		return e
	}

	if _, e := batch.ExecResults(); e != nil {
		return e
	}

	return batch.Close()
}

// DeleteByID - This function is used to delete a row from the table by id
// This function assumes one of the columns is called id, which is an int value
func (d pgxDriver) DeleteByID(ctx context.Context, tablename string, value int) (int64, error) {
	return d.Delete(ctx, tablename, "id", strconv.Itoa(value))
}

// DeleteQuery - This function creates a delete query for deleting a row from a table
func (d pgxDriver) Delete(ctx context.Context, tablename, columnname, value string) (int64, error) {
	var sb strings.Builder
	sb.WriteString("delete from ")
	sb.WriteString(tablename)
	sb.WriteString(" where ")
	sb.WriteString(columnname)
	sb.WriteString(" = ")
	sb.WriteString(QuoteString(value))
	return d.Exec(ctx, sb.String())
}

// Update - This function creates an update query based on variadic arguments
func (d pgxDriver) Update(ctx context.Context, args ...string) (int64, error) {
	for idx := range args {
		args[idx] = QuoteString(args[idx])
	}
	var sb strings.Builder
	sb.WriteString("update ")
	sb.WriteString(args[0])
	sb.WriteString(" set ")
	for index := range args[1:] {
		if index%2 > 0 {
			sb.WriteString(args[index])
			sb.WriteString(" = ")
			if index == len(args)-4 {
				if isString(args[index+1]) {
					sb.WriteString(args[index+1])
				} else {
					sb.WriteString("'")
					sb.WriteString(args[index+1])
					sb.WriteString("'")
				}
				sb.WriteString(" where ")
				sb.WriteString(args[index+2])
				sb.WriteString(" = ")
				sb.WriteString(args[index+3])
				break
			}
			if isString(args[index+1]) {
				sb.WriteString(args[index+1])
			} else {
				sb.WriteString("'")
				sb.WriteString(args[index+1])
				sb.WriteString("'")
			}
			sb.WriteString(", ")
		}
	}
	return d.Exec(ctx, sb.String())
}

// GetTablesName retrieves the names of the tables currently present in the database
func (d pgxDriver) GetTablesName(ctx context.Context) ([]string, error) {

	pgxRows, e := d.cp.QueryEx(ctx, getTablesName, nil)
	if e != nil {
		return nil, e
	}

	defer pgxRows.Close()

	return unMarshallRowsToTableNames(pgxRows)
}

func (d pgxDriver) getHyperTablesName() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	pgxRows, e := d.cp.QueryEx(ctx, tsdbGetHypertables, nil)
	if e != nil {
		return nil, e
	}

	defer pgxRows.Close()

	return unMarshallRowsToTableNames(pgxRows)
}

// unMarshallRowsToTableNames unmarshall the names of the tables received from the
// database into the relevant struct
func unMarshallRowsToTableNames(rows *pgx.Rows) ([]string, error) {
	tablesName := []string{}

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}

		tablesName = append(tablesName, tableName)
	}
	return tablesName, nil
}

// GetFields retrieves the columns belonging to a particular postgres table
func (d pgxDriver) GetFields(ctx context.Context, tableName string) ([]Field, error) {
	pgxRows, e := d.cp.QueryEx(ctx, fmt.Sprintf(getColumnsMetaData, tableName), nil)
	if e != nil {
		return nil, e
	}

	defer pgxRows.Close()
	return unMarshallRowsToColumnsMeta(pgxRows)
}

// unMarshallRowsToColumnsMeta unmarshall postgres columns received from the DB
// into the relevant struct
func unMarshallRowsToColumnsMeta(rows *pgx.Rows) ([]Field, error) {
	fields := []Field{}

	for rows.Next() {
		f := Field{}
		err := rows.Scan(
			&f.Name,
			&f.Datatype,
		)

		if err != nil {
			return nil, err
		}

		fields = append(fields, f)
	}

	return fields, nil
}

// CreateTable create a new table into the DB, if isTimeSeries is set to true a column
// called time with timestamp as datatype will be automatically created for the table
// and the relevant timescale hypertable will be created as well
func (d pgxDriver) CreateTable(ctx context.Context, rows Rows, isTimeSeries bool, constrains ...string) (int64, error) {
	var (
		sb            strings.Builder
		defaultIntVal int64
	)
	sb.WriteString("\"")
	sb.WriteString(rows.TableName)
	sb.WriteString("\" (")

	for _, f := range rows.Fields {
		if len(constrains) > 0 {
			for _, c := range constrains {
				if f.Name == c {
					sb.WriteString(fmt.Sprintf("\"%s\" %s not null unique, ", f.Name, f.Datatype))
				} else {
					sb.WriteString(fmt.Sprintf("\"%s\" %s, ", f.Name, f.Datatype))
				}
			}
		} else {
			sb.WriteString(fmt.Sprintf("\"%s\" %s, ", f.Name, f.Datatype))
		}
	}

	if _, e := d.Exec(ctx, fmt.Sprintf(createTable, strings.TrimRight(sb.String(), ", ")+")")); e != nil {
		return defaultIntVal, e
	}

	if isTimeSeries {
		if _, e := d.createHyperTable(ctx, rows.TableName); e != nil {
			return defaultIntVal, e
		}
	}

	return defaultIntVal, nil
}

// createHyperTable creates a new timescale DB hyper table
func (d pgxDriver) createHyperTable(ctx context.Context, tableName string) (int64, error) {
	var sb strings.Builder

	sb.WriteString("\"")
	sb.WriteString(tableName)
	sb.WriteString("\"")

	return d.Exec(ctx, fmt.Sprintf(createHyperTable, sb.String()))
}

// CreateIndex create an index on the specified table
func (d pgxDriver) CreateIndex(ctx context.Context, indexName, tableName string, fieldsName []string) (int64, error) {
	var sb strings.Builder

	sb.WriteString("(")

	for _, fn := range fieldsName {
		sb.WriteString("\"")
		sb.WriteString(fn)
		sb.WriteString("\", ")
	}

	queryDetails := strings.TrimRight(sb.String(), ", ") + ")"

	sb.Reset()
	sb.WriteString("\"")
	sb.WriteString(tableName)
	sb.WriteString("\"")

	return d.Exec(ctx, fmt.Sprintf(createIndex, "\""+indexName+"\"", sb.String(), queryDetails))
}

// DeleteTableValues delete values from a particular table older than the number of days specified
func (d pgxDriver) DeleteValuesFromInterval(ctx context.Context, tableName string, days int32) (int64, error) {
	var (
		isHyperTable  bool
		defaultIntVal int64
	)

	pgxRows, e := d.cp.QueryEx(ctx, checkTSDBavailability, nil)
	if e != nil {
		return defaultIntVal, e
	}

	defer pgxRows.Close()

	var values []interface{}
	for pgxRows.Next() {
		values, e = pgxRows.Values()
		if e != nil {
			return defaultIntVal, e
		}
	}

	if len(values) == 0 {
		return d.Exec(ctx, fmt.Sprintf(pgDeleteValuesByInterval, tableName, days))
	}

	hyperTablesName, e := d.getHyperTablesName()
	if e != nil {
		return defaultIntVal, e
	}

	for _, tn := range hyperTablesName {
		if tn == tableName {
			isHyperTable = true
		}
	}

	if isHyperTable {
		return d.Exec(ctx, fmt.Sprintf(tsdbDeleteValuesByInterval, days, tableName))
	}

	return d.Exec(ctx, fmt.Sprintf(pgDeleteValuesByInterval, tableName, days))
}

// Insert inserts rows into the database
func (d pgxDriver) Insert(ctx context.Context, rows Rows, fields ...Field) (Rows, error) {
	var sb strings.Builder
	r := Rows{}
	if rows.SchemaName == "" {
		rows.SchemaName = "public"
	}
	sb.WriteString(rows.SchemaName)
	sb.WriteString(".")
	sb.WriteString("\"")
	sb.WriteString(rows.TableName)
	sb.WriteString("\"")

	var query string
	if len(fields) == 0 {
		query = fmt.Sprintf(insert,
			sb.String(),
			buildFieldsString(rows.Fields),
			buildValuesString(rows))
		_, err := d.Exec(ctx, query)
		return r, err
	}
	query = fmt.Sprintf(insert,
		sb.String(),
		buildFieldsString(rows.Fields),
		buildValuesString(rows)+" returning "+buildFieldsString(fields),
	)
	return d.Query(ctx, query)
}

// Upsert upserts rows into the database
func (d pgxDriver) Upsert(ctx context.Context, rows Rows, constrain string) (int64, error) {
	var sb strings.Builder

	sb.WriteString("\"")
	sb.WriteString(rows.TableName)
	sb.WriteString("\"")

	query := fmt.Sprintf(upsert,
		sb.String(),
		buildFieldsString(rows.Fields),
		buildValuesString(rows),
		constrain,
		buildOnConstrainString(rows.Fields, constrain))

	return d.Exec(ctx, query)
}

// Upsert upserts rows into the database
func (d pgxDriver) Copy(rows Rows) (err error) {
	fieldNames := []string{}

	for _, f := range rows.Fields {
		fieldNames = append(fieldNames, f.Name)
	}

	_, err = d.cp.CopyFrom(
		pgx.Identifier{rows.SchemaName, rows.TableName},
		fieldNames,
		pgx.CopyFromRows(rows.Values),
	)

	return
}

func isString(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return false
	}
	return isJSON(s)
}

// This function identifies whether or not a particular string is actually a JSON
// If it is JSON, the output is false
// If it is not JSON, the output is true
func isJSON(s string) bool {
	return strings.Contains(s, "{")
}

// buildReturningClause builds the columns substring for the insert query for the columns that are returned
// i.e (ColName1, ColName2, ...., ColNameN)
func buildReturningClause(fields []Field) string {
	var sb strings.Builder

	for _, field := range fields {
		sb.WriteString(field.Name)
		sb.WriteString(", ")
	}

	return strings.TrimRight(sb.String(), ", ")
}

// buildFieldsString build the columns substring for the insert query
// i.e (ColName1, ColName2, ...., ColNameN)
func buildFieldsString(fields []Field) string {
	var sb strings.Builder
	sb.WriteString("(")

	for _, field := range fields {
		sb.WriteString("\"")
		sb.WriteString(field.Name)
		sb.WriteString("\", ")
	}

	return strings.TrimRight(sb.String(), ", ") + ")"
}

// buildValuesString build the values substring for the insert query
// i.e (ColVal1, ColVal3, ...., ColValN)
func buildValuesString(rows Rows) string {
	var sb strings.Builder
	sb.WriteString("(")

	for _, rowValue := range rows.Values {
		for i, value := range rowValue {
			switch rows.Fields[i].Datatype {
			case Timestamp, Text, JSONB, UUID:
				if value != nil {
					sb.WriteString(fmt.Sprintf("'%s', ", QuoteString(value.(string))))

				} else {
					sb.WriteString("NULL, ")
				}
			case Double:
				if math.IsInf(value.(float64), 1) || math.IsInf(value.(float64), -1) || math.IsNaN(value.(float64)) {
					sb.WriteString(fmt.Sprintf("'%v', ", value))
				} else {
					sb.WriteString(fmt.Sprintf("%v, ", value))
				}
			case Blob:
				if value != nil {
					sb.WriteString(fmt.Sprintf("E'%s', ", value))

				} else {
					sb.WriteString("NULL, ")
				}
			default:
				if value != nil {
					sb.WriteString(fmt.Sprintf("%v, ", value))
				} else {
					sb.WriteString("NULL, ")
				}
			}
		}
		str := strings.TrimRight(sb.String(), ", ")
		sb.Reset()
		sb.WriteString(str)
		sb.WriteString("), (")
	}

	return strings.TrimRight(sb.String(), ", (")
}

// buildOnConstrainString build the on constrain substring for the
// insert query i.e (ColName1=EXCLUDED.COlName1,...., ColNameN=EXCLUDED.COlNameN)
func buildOnConstrainString(fields []Field, constrain string) string {
	var sb strings.Builder
	var fsb strings.Builder

	exclString := "%s=EXCLUDED.%s, "

	for _, field := range fields {
		if field.Name != constrain {

			fsb.WriteString("\"")
			fsb.WriteString(field.Name)
			fsb.WriteString("\"")
			sb.WriteString(fmt.Sprintf(exclString, fsb.String(), fsb.String()))
			fsb.Reset()
		}
	}

	return strings.TrimRight(sb.String(), ", ")
}

// GetLatestValue returns latest values for a list of tags belonging to the same table
func (d pgxDriver) GetLatestValue(ctx context.Context, fieldsName []string, tableName string) ([]Rows, error) {
	queries := []string{}
	for _, fieldName := range fieldsName {
		queries = append(queries, fmt.Sprintf(getLatestValue, "\""+fieldName+"\"", "\""+tableName+"\"", "\""+fieldName+"\""))
	}

	return d.BatchQuery(ctx, queries)
}

// GetTimeSeriesValue returns time series values for a list of tags belonging to the same table
func (d pgxDriver) GetTimeSeriesValues(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time) ([]Rows, error) {
	queries := []string{}
	st := string(time.Unix(0, startTime.UnixNano()).UTC().AppendFormat(make([]byte, 0, len(timeFormat)), timeFormat))
	et := string(time.Unix(0, endTime.UnixNano()).UTC().AppendFormat(make([]byte, 0, len(timeFormat)), timeFormat))

	for _, fieldName := range fieldsName {
		queries = append(queries, fmt.Sprintf(getTimeSeriesValues, "\""+fieldName+"\"", "\""+tableName+"\"", st, et))
	}

	return d.BatchQuery(ctx, queries)
}

// GetDownSampledValues returns down sampled values for a list of tags belonging to the same table
func (d pgxDriver) GetDownSampledValues(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, seconds int) ([]Rows, error) {
	queries := []string{}

	st := string(time.Unix(0, startTime.UnixNano()).UTC().AppendFormat(make([]byte, 0, len(timeFormat)), timeFormat))
	et := string(time.Unix(0, endTime.UnixNano()).UTC().AppendFormat(make([]byte, 0, len(timeFormat)), timeFormat))

	for _, fieldName := range fieldsName {
		queries = append(queries, fmt.Sprintf(getDownSampledValues, seconds, "\""+fieldName+"\"", "\""+tableName+"\"", st, et))
	}

	return d.BatchQuery(ctx, queries)
}

// GetInterpolatedValues returns interpolated values for a list of tags belonging to the same table
func (d pgxDriver) GetInterpolatedValues(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, seconds int) ([]Rows, error) {
	queries := []string{}

	st := string(time.Unix(0, startTime.UnixNano()).UTC().AppendFormat(make([]byte, 0, len(timeFormat)), timeFormat))
	et := string(time.Unix(0, endTime.UnixNano()).UTC().AppendFormat(make([]byte, 0, len(timeFormat)), timeFormat))
	fmt.Println(st)
	fmt.Println(et)
	for _, fieldName := range fieldsName {
		queries = append(queries, fmt.Sprintf(getInterpolatedValues, seconds, st, et, "\""+fieldName+"\"", "\""+tableName+"\"", st, et))
	}
	return d.BatchQuery(ctx, queries)
}

func (d pgxDriver) extractRowsFromPgxRows(ctx context.Context, pgxRows *pgx.Rows) (Rows, error) {
	rows := Rows{
		Fields: []Field{},
		Values: [][]interface{}{},
	}
	fieldsDesc := pgxRows.FieldDescriptions()
	for _, fieldDesc := range fieldsDesc {

		rows.Fields = append(rows.Fields, Field{Name: fieldDesc.Name, Datatype: Datatype(fieldDesc.DataTypeName)})
	}

	for pgxRows.Next() {
		values, e := pgxRows.Values()
		if e != nil {
			return rows, e
		}

		rows.Values = append(rows.Values, values)
	}
	return rows, nil
}

//GetHistorianTimeSeriesValues
func (d pgxDriver) GetHistorianTimeSeriesValues(ctx context.Context, fieldsName string, tableName string, startTime, endTime time.Time, desc bool, where HistoricWhere, limit uint32) (Rows, error) {
	//Format time
	st := string(time.Unix(0, startTime.UnixNano()).UTC().AppendFormat(make([]byte, 0, len(timeFormat)), timeFormat))
	et := string(time.Unix(0, endTime.UnixNano()).UTC().AppendFormat(make([]byte, 0, len(timeFormat)), timeFormat))
	var query strings.Builder
	//Build the SelectFrom
	query.WriteString(fmt.Sprintf(selectFrom, fieldsName, tableName))
	//Build the Where
	switch where {
	case HistoricWhereTimeGreater:
		query.WriteString(fmt.Sprintf(whereTimetGreater, st))
	case HistoricWhereTimeLess:
		query.WriteString(fmt.Sprintf(whereTimeLess, st))
	case HistoricWhereTimeSlot:
		query.WriteString(fmt.Sprintf(whereTimeSlot, st, et))
	case HistoricWhereTimeSlotEqStart:
		query.WriteString(fmt.Sprintf(whereTimeSlotEqStart, st, et))
	case HistoricWhereTimeSlotEqEnd:
		query.WriteString(fmt.Sprintf(whereTimeSlotEqEnd, st, et))
	}
	//Add the not Null
	query.WriteString(fmt.Sprintf(isNotNull, fieldsName))
	//Build the Asc/Desc
	if desc {
		query.WriteString(orderByDesc)
	} else {
		query.WriteString(orderByAsc)
	}
	//Build the Limit
	if limit > 0 {
		query.WriteString(fmt.Sprintf(setLimit, limit))
	}
	//End
	query.WriteString(end)
	// fmt.Println(query.String())
	//Query
	return d.Query(ctx, query.String())
}

//GetHistorianAggrMax
func (d pgxDriver) GetHistorianAggrMax(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, desc bool, limit uint32, interval int64) ([]Rows, error) {
	return nil, nil
}

//GetHistorianAggrMin
func (d pgxDriver) GetHistorianAggrMin(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, desc bool, limit uint32, interval int64) ([]Rows, error) {
	return nil, nil
}

//GetHistorianAggrCount
func (d pgxDriver) GetHistorianAggrCount(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, desc bool, limit uint32, interval int64) ([]Rows, error) {
	return nil, nil
}

//GetHistorianAggrAvg
func (d pgxDriver) GetHistorianAggrAvg(ctx context.Context, fieldsName []string, tableName string, startTime, endTime time.Time, desc bool, limit uint32, interval int64) ([]Rows, error) {
	return nil, nil
}

// Stat return connection pool statistics
func (d pgxDriver) Stat() ConnPoolStat {
	stats := d.cp.Stat()
	return ConnPoolStat{MaxConnections: stats.MaxConnections, CurrentConnections: stats.CurrentConnections, AvailableConnections: stats.AvailableConnections}
}
