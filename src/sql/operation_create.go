package sql

type CreateOperation struct {
	TableName string
	Data      []ColData
}

func (operation *CreateOperation) Call(database *Database) []byte {
	database.Create(operation.TableName, operation.Data)
	return nil
}
