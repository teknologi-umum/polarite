package resources

import (
	"net/url"
	"strings"
)

// ParseURL converts given database connection string url into a DSN format.
func ParseURL(connstr string) (out string, err error) {
	// Validate it first
	if strings.HasPrefix(connstr, "mysql://") {
		// MySQL DSN format is: username:password@tcp(127.0.0.1:3306)/test
		parsedConn, err := url.Parse(connstr)
		if err != nil {
			return out, err
		}

		out += parsedConn.User.String()
		out += "@"
		out += "tcp("
		out += parsedConn.Hostname()
		out += ":"
		out += parsedConn.Port()
		out += ")"
		out += parsedConn.EscapedPath()
		out += "?"
		out += parsedConn.Query().Encode()

		return out, err
	} else if strings.HasPrefix(connstr, "postgres://") || strings.HasPrefix(connstr, "postgresql://") {
		// PostgreSQL DSN format is: user=bob password=secret host=1.2.3.4 port=5432 dbname=mydb sslmode=verify-full
		parsedConn, err := url.Parse(connstr)
		if err != nil {
			return out, err
		}

		pwd, _ := parsedConn.User.Password()

		out += "user="
		out += parsedConn.User.Username()
		out += " password="
		out += pwd
		out += " host="
		out += parsedConn.Hostname()
		out += " port="
		out += parsedConn.Port()
		out += " dbname="
		out += strings.Replace(parsedConn.EscapedPath(), "/", "", 1)
		out += " "
		out += strings.Join(strings.Split(parsedConn.Query().Encode(), "&"), " ")

		return out, err
	}

	return connstr, nil
}
