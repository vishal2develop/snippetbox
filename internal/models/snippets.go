package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet. fields of the struct correspond to the fields in our MySQL snippets table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
// all snippet-related queries go through this model.
// *sql.DB is a connection pool, not a single connection.
type SnippetModel struct {
	DB *sql.DB
}

// insert into snippets table
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// sql insert query. using backquotes to split the query into multiple lines
	statement := `INSERT INTO snippets 
    (title, content, created, expires)
VALUES (?,?,UTC_TIMESTAMP(),DATE_ADD(UTC_TIMESTAMP(),INTERVAL ? DAY))`
	// Use the Exec() method on the embedded connection pool to execute the statement.
	result, err := m.DB.Exec(statement, title, content, expires)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result to get the ID of our newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type before returning.
	return int(id), nil
}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	statement := `SELECT * FROM snippets WHERE expires > UTC_TIMESTAMP() and id = ?`

	// QueryRow - to query one record
	row := m.DB.QueryRow(statement, id)

	// Initialize a new zeroed Snippet struct.
	/* s := Snippet{} vs var s Snippet
	* s := Snippet{} -> creates an empty snippet, even if no records returned. When you call GET on an invalid id, in Addition to 404, u will also see empty struct.
	* var s Snippet -> doesn’t create a structure until there are valid records to populate. When you call GET on an invalid id, you will get only a 404 and not see an empty structure.
	 */
	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		// check if query return no records error
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}
	// if everything went ok, return filled snippet struct
	return s, nil

}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {
	statement := "SELECT * FROM snippets where expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10"

	rows, err := m.DB.Query(statement)

	if err != nil {
		return nil, err
	}

	// defer row.close() to ensure sql.Rows resultset always properly closed
	// defer should come always after error check, otherwise, if Query() returns an error, you'll get a panic trying to close a nil resultset
	defer rows.Close()

	// initialize an empty slice to hold the snippet structs
	var snippets []Snippet

	// rows.Next to iterate through the rows in the resultset lazily
	// any errors whole iterating will not terminate the loop, that is why we have to do a final error check after loop completion
	for rows.Next() {
		// create a new zeroed value snippet struct
		var s Snippet
		// rows.Scan() to copy the values from each field in the row
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets
		snippets = append(snippets, s)

	}

	// Very important - Once rows.next has completed iterating over the result set, check for any error during iteration
	// idiomatic go way:
	/**
	Call rows.Err()
	Assign its result to the existing err variable
	Immediately check it
	*/
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// if everything went ok return the snippets
	return snippets, nil
}
