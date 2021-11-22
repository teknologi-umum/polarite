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

		var out strings.Builder

		out.WriteString(parsedConn.User.String())
		out.WriteString("@")
		out.WriteString("tcp(")
		out.WriteString(parsedConn.Hostname())
		out.WriteString(":")
		out.WriteString(parsedConn.Port())
		out.WriteString(")")
		out.WriteString(parsedConn.EscapedPath())
		out.WriteString("?")
		out.WriteString(parsedConn.Query().Encode())

		return out.String(), err
	} else if strings.HasPrefix(connstr, "postgres://") || strings.HasPrefix(connstr, "postgresql://") {
		// PostgreSQL DSN format is: user=bob password=secret host=1.2.3.4 port=5432 dbname=mydb sslmode=verify-full
		parsedConn, err := url.Parse(connstr)
		if err != nil {
			return out, err
		}

		pwd, _ := parsedConn.User.Password()

		var out strings.Builder

		out.WriteString("user=")
		out.WriteString(parsedConn.User.Username())
		out.WriteString(" password=")
		out.WriteString(pwd)
		out.WriteString(" host=")
		out.WriteString(parsedConn.Hostname())
		out.WriteString(" port=")
		out.WriteString(parsedConn.Port())
		out.WriteString(" dbname=")
		out.WriteString(strings.Replace(parsedConn.EscapedPath(), "/", "", 1))
		out.WriteString(" ")
		out.WriteString(strings.Join(strings.Split(parsedConn.Query().Encode(), "&"), " "))

		return out.String(), err
	}

	return connstr, nil
}
