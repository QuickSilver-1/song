package customError

import "song/internal/domain"

// Алиасы типов для различных типов ошибок
type MigratingError = domain.BaseError

type DbConnectionError = domain.BaseError

type DbQueryError = domain.BaseError

type RowsNotFoundError = domain.BaseError

type LoggerBuildError = domain.BaseError

type RedisQueryError = domain.BaseError

type InvalidInputData = domain.BaseError
