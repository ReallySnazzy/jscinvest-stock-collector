package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/reallysnazzy/jscinvest-stock-collector/iexcloud"
)

const (
	databaseName = "stocktracker"
)

type DatabaseContext struct {
	host     string
	port     int32
	user     string
	password string
	db       *sql.DB
	tx       *sql.Tx
}

func CreateDatabaseContext(host string, port int32, user string, password string) (DatabaseContext, error) {
	ctx := DatabaseContext{}
	ctx.host = host
	ctx.port = port
	ctx.user = user
	ctx.password = password
	err := ctx.connect()
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (this *DatabaseContext) connect() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", this.host, this.port, this.user, this.password, databaseName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	this.db = db
	return nil
}

func (this *DatabaseContext) BeginTransaction() error {
	ctx := context.Background()
	tx, err := this.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	this.tx = tx
	return nil
}

func (this *DatabaseContext) CommitTransaction() error {
	err := this.tx.Commit()
	if err != nil {
		return err
	}
	this.tx = nil
	return nil
}

func (this *DatabaseContext) Disconnect() {
	this.db.Close()
}

func (this *DatabaseContext) GetNextSymbolsToTrack() ([]string, error) {
	const query = `
select 
	ticker
from 
	stock_symbols_to_track
order by
	(case when last_checkin is null then 1 else 0 end) desc,
	last_checkin asc
limit 10;
	`
	rows, err := this.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{}
	for rows.Next() {
		var ticker string
		err = rows.Scan(&ticker)
		if err != nil {
			return nil, err
		}
		result = append(result, ticker)
	}
	return result, nil
}

func (this *DatabaseContext) MarkSymbolCompleted(symbol string) error {
	const query = `
update 
	stock_symbols_to_track 
set 
	last_checkin=now()
where 
	ticker=$1;
	`
	_, err := this.db.Exec(query, symbol)
	return err
}

func (this *DatabaseContext) AddStockPriceHistory(history iexcloud.IexcloudHistoricalPriceEntry) error {
	const query = `
insert into
	stock_price_history 
	(close, high, low, open, symbol, volume, iexcloud_key, iexcloud_subkey, date, update_timestamp,
		change_over_time, market_change_over_time, uclose, uhigh, ulow, uopen, uvolume,
		fclose, fhigh, flow, fopen, fvolume, iexcloud_label, change, change_percent)
values
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
		$21, $22, $23, $24, $25);
	`
	var err error
	_, err = this.db.Exec(query,
		history.Close,
		history.High,
		history.Low,
		history.Open,
		history.Symbol,
		history.Volume,
		history.Key,
		history.SubKey,
		history.Date,
		history.UpdateTimestamp,
		history.ChangeOverTime,
		history.MarketChangeOverTime,
		history.UClose,
		history.UHigh,
		history.ULow,
		history.UOpen,
		history.UVolume,
		history.FClose,
		history.FHigh,
		history.FLow,
		history.FOpen,
		history.FVolume,
		history.Label,
		history.Change,
		history.ChangePercent)
	return err
}
