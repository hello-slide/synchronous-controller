package database

type Config struct {
	DriveName *DriveName
}

type DriveName int

const (
	PostogreSQL DriveName = iota
	MySQL
	SQLite
)

func (c DriveName) String() string {
	switch c {
	case PostogreSQL:
		return "postgres"
	case MySQL:
		return "mysql"
	case SQLite:
		return "sqlite3"
	default:
		return "unknown"
	}
}
