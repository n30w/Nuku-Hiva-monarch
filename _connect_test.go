package main

import "testing"

type MockPlanetscaleDB struct{}

func (m *MockPlanetscaleDB) InsertToSQL(table *Table[Row[id, text]]) error {
	return nil
}

func (m *MockPlanetscaleDB) UpdateSQL(table *Table[Row[id, text]]) error {
	return nil
}

func TestInsertToSQL(t *testing.T) {
	m := &MockPlanetscaleDB{}

}
