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
        err = rows.Scan(&campaign.Id, &campaign.Name, &campaign.Status,&campaign.Type, &campaign.Budget, &campaign.Created_on)

        if err != nil {
            log.Fatalf("Unable to scan the row. %v", err)
        }

        // append the campaign in the campaigns slice
        campaigns = append(campaigns, campaign)

	}

    // return empty campaign on error
    return campaigns, err
}


func CreateCampaign(w http.ResponseWriter, r *http.Request) {
    // set the header to content type x-www-form-urlencoded
    // Allow all origin to handle cors issue
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

    // create an empty campaign of type models.User
    var campaign models.Campaign

    // decode the json request to campaign
    err := json.NewDecoder(r.Body).Decode(&campaign)

    if err != nil {
        log.Fatalf("Unable to decode the request body.  %v", err)
    }

    // call insert campaign function and pass the campaign
    insertID := insertCampaign(campaign)

    // format a response object
    res := response{
        ID:      insertID,
        Message: "Campaign created successfully",
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}

func insertCampaign(campaign models.Campaign) int64 {

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create the insert sql query
    // returning userid will return the id of the inserted campaign
    sqlStatement := `INSERT INTO campaigns (name, status,type,budget) VALUES ($1, $2,$3,$4) RETURNING id`

    // the inserted id will store in this id
    var id int64

    // execute the sql statement
    // Scan function will save the insert id in the id
    err := db.QueryRow(sqlStatement, campaign.Name, campaign.Status,campaign.Type, campaign.Budget).Scan(&id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    fmt.Printf("Inserted a single record %v", id)

    // return the inserted id
    return id
}
