package utils // import "github.com/amieldelatorre/notifi/backend/utils"

type ContextKey string

const (
	RequestIdName ContextKey = "RequestId"
	UserId        ContextKey = "UserId"
)
