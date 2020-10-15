package controllers

import (
    "database/sql"
    "encoding/json" 
    "fmt"
    "backend/models" 
    "log"
    "net/http" 
    "github.com/joho/godotenv" 
    _ "github.com/lib/pq"      
)


type response struct {
    ID      int64  `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
    // load .env file
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Open the connection
    db, err := sql.Open("postgres", "postgres://yhlmstfe:h2KsNK160WTAztoFYgg5u0FAWIgwcGyt@raja.db.elephantsql.com:5432/yhlmstfe")

    if err != nil {
        panic(err)
    }

    // check the connection
    err = db.Ping()

    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully connected!")
    // return the connection
    return db
}

func GetAllCampaigns(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    // get all the campaigns in the db
    campaigns, err := getAllCampaigns()

    if err != nil {
        log.Fatalf("Unable to get all campaign. %v", err)
    }

	// send all the campaigns as response

    json.NewEncoder(w).Encode(campaigns)
}
// get one campaign from the DB by its userid
func getAllCampaigns() ([]models.Campaign, error) {
    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    var campaigns []models.Campaign

    // create the select sql query
    sqlStatement := `SELECT * FROM campaigns`

    // execute the sql statement
    rows, err := db.Query(sqlStatement)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }
    

    // close the statement
    defer rows.Close()

    // iterate over the rows
    for rows.Next() {
        var campaign models.Campaign

        // unmarshal the row object to campaign
        err = rows.Scan(&campaign.Id, &campaign.Name, &campaign.Status)

        if err != nil {
            log.Fatalf("Unable to scan the row. %v", err)
        }

        // append the campaign in the campaigns slice
        campaigns = append(campaigns, campaign)

	}

    // return empty campaign on error
    return campaigns, err
}