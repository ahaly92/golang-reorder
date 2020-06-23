package sql

const (
	getTablesName              = "SELECT tableName FROM pg_catalog.pg_tables WHERE schemaname='public'"
	getColumnsMetaData         = "SELECT column_name, data_type from INFORMATION_SCHEMA.COLUMNS WHERE table_name = '%s'"
	createTable                = "CREATE TABLE IF NOT EXISTS %s"
	createHyperTable           = "SELECT create_hypertable('%s','Time',chunk_time_interval => interval '1 days',if_not_exists => TRUE)"
	createIndex                = "CREATE INDEX IF NOT EXISTS %s ON %s %s"
	checkTSDBavailability      = "SELECT * FROM pg_extension WHERE \"extname\" ='timescaledb'"
	pgDeleteValuesByInterval   = "DELETE FROM \"%s\" WHERE \"Time\" < now() - interval '%v days'"
	tsdbDeleteValuesByInterval = "SELECT drop_chunks(interval '%v days', '%s')"
	tsdbGetHypertables         = "SELECT table_name from _timescaledb_catalog.hypertable;"
	getLatestValue             = "SELECT \"Time\", %s FROM %s WHERE %s IS NOT NULL ORDER BY TIME DESC LIMIT 1;"
	getTableFromOID            = "SELECT relname from pg_class WHERE \"oid\"='%v'"
	getTimeSeriesValues        = "SELECT \"Time\", %s FROM %s WHERE \"Time\" BETWEEN '%s' AND '%s' ORDER BY \"Time\" DESC"
	getDownSampledValues       = "SELECT time_bucket('%v seconds', \"Time\") AS dateTime,avg(%s) FROM %s WHERE \"Time\" BETWEEN '%s' AND '%s' GROUP BY dateTime ORDER BY dateTime DESC"
	upsert                     = "INSERT INTO %s %s VALUES %s ON CONFLICT(%q) DO UPDATE SET %s"
	insert                     = "INSERT INTO %s %s VALUES %s"
	getInterpolatedValues      = "SELECT time_bucket_gapfill('%v seconds', \"Time\", '%s','%s') AS dateTime, locf(avg(%s)) FROM %s WHERE \"Time\" BETWEEN '%s' and '%s' GROUP BY dateTime ORDER BY dateTime DESC;"
	selectFrom                 = "SELECT \"Time\", \"%s\" FROM \"%s\""
	setLimit                   = " LIMIT %d"
	whereTimeSlot              = " WHERE \"Time\" > '%s' AND \"Time\" < '%s'"
	whereTimeSlotEqStart       = " WHERE \"Time\" >= '%s' AND \"Time\" < '%s'"
	whereTimeSlotEqEnd         = " WHERE \"Time\" > '%s' AND \"Time\" <= '%s'"
	whereTimeLess              = " WHERE \"Time\" > '%s'"
	whereTimetGreater          = " WHERE \"Time\" < '%s'"
	isNotNull                  = " AND \"%s\" IS NOT NULL"
	orderByAsc                 = " ORDER BY \"Time\" ASC"
	orderByDesc                = " ORDER BY \"Time\" DESC"
	end                        = ";"
)

type HistoricWhere uint8

const (
	HistoricWhereTimeGreater HistoricWhere = iota
	HistoricWhereTimeLess
	HistoricWhereTimeSlot
	HistoricWhereTimeSlotEqStart
	HistoricWhereTimeSlotEqEnd
)

// RAW := "SELECT \"Time\", %s FROM %s WHERE \"Time\" > '%s' AND \"Time\" =< %s ORDER BY \"Time\" ASC LIMIT %d;"
// PROCESSED := "SELECT \"Time\", AVG(%s), time_bucket(intevral, start, tablename) WHERE \"Time\" >= '%s' AND \"Time\" < %s ORDER BY \"Time\" DESC;"
