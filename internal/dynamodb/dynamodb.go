package dynamodb

const (
	APIVersion20120810 = "DynamoDB_20120810"
)

type OperationType string

const (
	OperationTypeDeleteItem OperationType = "DeleteItem"
	OperationTypeGetItem    OperationType = "GetItem"
	OperationTypePutItem    OperationType = "PutItem"
)
