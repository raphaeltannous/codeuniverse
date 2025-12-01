package postgres

type postgresScanner interface {
	Scan(dest ...any) error
}
