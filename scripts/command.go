package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Peter-Tabarani/PiconexBackend/internal/utils"
)

func main() {
	// Your existing DSN
	dsn := "piconex:pjaplmTabs7!@tcp(178.156.189.138:3306)/piconexdb"

	// Connect to the database using your existing utils.Connect
	db, err := utils.Connect(dsn)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	// --- Put any SQL query you want here ---
	query := `
ALTER TABLE poc_admin CHANGE COLUMN activity_id point_of_contact_id INT NOT NULL;
`
	// Decide if it's a query (returns rows) or command (update/insert)
	// We'll try Query first, then fallback to Exec if no rows returned
	rows, err := db.Query(query)
	if err != nil {
		// Probably not a SELECT — try Exec
		res, execErr := db.Exec(query)
		if execErr != nil {
			log.Fatal("❌ Query/Exec failed:", execErr)
		}
		affected, _ := res.RowsAffected()
		fmt.Printf("✅ Command executed successfully. Rows affected: %d\n", affected)
		return
	}
	defer rows.Close()

	// Print all columns from the query for quick inspection
	cols, _ := rows.Columns()
	allRows := []map[string]interface{}{}

	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal("❌ Row scan failed:", err)
		}

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			val := values[i]
			if b, ok := val.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = val
			}
		}

		allRows = append(allRows, rowMap)
	}

	// Pretty print as JSON
	out, _ := json.MarshalIndent(allRows, "", "  ")
	fmt.Println(string(out))

}
