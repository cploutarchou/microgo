/*
 * MIT License
 *
 * Copyright (c) 2022 Christos Ploutarchou
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package MicroGO

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func (m *MicroGo) OpenDB(driverName, dataSourceName string) (*bun.DB, error) {
	var sqldb *sql.DB
	var db *bun.DB
	var err error
	switch {
	case driverName == "postgres" || driverName == "postgresql":
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dataSourceName)))
		db = bun.NewDB(sqldb, pgdialect.New())
	case driverName == "mysql" || driverName == "mariadb":
		driverName = "mysql"
		sqldb, err = sql.Open(driverName, dataSourceName)
	default:
		driverName = "sqlite"
		sqldb, err = sql.Open(driverName, dataSourceName)
		sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
		db = bun.NewDB(sqldb, sqlitedialect.New())
	}

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
