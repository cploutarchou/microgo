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
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (m *MicroGo) OpenDB(driverName, dataSourceName string) (*gorm.DB, error) {
	switch {
	case driverName == "postgres" || driverName == "postgresql":
		db, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})

		if err != nil {
			return nil, err
		}
		return db, nil
	case driverName == "mysql" || driverName == "mariadb":
		db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		return db, nil
	default:
		driverName = "sqlite"
		db, err := gorm.Open(sqlite.Open(m.AppName), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		return db, nil
	}

}
