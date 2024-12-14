package sql

type InsertOperation struct {
	TableName string
	Data      []ValData
}

func (operation *InsertOperation) Call(database *Database) []byte {
	table, _ := database.Get(operation.TableName)
	table.Insert(operation.Data)
	return nil
}
