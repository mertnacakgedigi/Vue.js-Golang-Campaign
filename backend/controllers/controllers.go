package controllers

import (
    "database/sql"
    "encoding/json" 
    "fmt"
    "backend/models" 
    "log"
    "net/http" 
    "github.com/joho/godotenv" 
    "strconv"  // package used to covert string into int type
    "github.com/gorilla/mux" // used to get the params from the route
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
        err = rows.Scan(&campaign.Id, &campaign.Name, &campaign.Status,&campaign.Type, &campaign.Budget, &campaign.Created_at)

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
    fmt.Println(r.Body)

    // decode the json request to campaign
    err := json.NewDecoder(r.Body).Decode(&campaign)

    if err != nil {
        log.Fatalf("Unable to decode the request body.  %v", err)
    }

    fmt.Println(campaign)

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
    sqlStatement := `INSERT INTO campaigns (id,name, status,type,budget,created_at) VALUES ($1, $2,$3,$4,$5,$6) RETURNING id`

    // the inserted id will store in this id
    var id int64

    // execute the sql statement
    // Scan function will save the insert id in the id
    err := db.QueryRow(sqlStatement, campaign.Id,campaign.Name, campaign.Status,campaign.Type, campaign.Budget,campaign.Created_at).Scan(&id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    fmt.Printf("Inserted a single record %v", id)

    // return the inserted id
    return id
}


func DeleteCampaign(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // get the userid from the request params, key is "id"
    params := mux.Vars(r)

    // convert the id in string to int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    // call the deleteUser, convert the int to int64
    deletedRows := deleteCampaign(int64(id))

    // format the message string
    msg := fmt.Sprintf("Campaign deleted successfully. Total rows/record affected %v", deletedRows)

    // format the reponse message
    res := response{
        ID:      int64(id),
        Message: msg,
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}


func deleteCampaign(id int64) int64 {

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create the delete sql query
    sqlStatement := `DELETE FROM campaigns WHERE id=$1`

    // execute the sql statement
    res, err := db.Exec(sqlStatement, id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    // check how many rows affected
    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }

    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}

func UpdateCampaign(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", " application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "PUT")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

    // get the userid from the request params, key is "id"
    params := mux.Vars(r)

    // convert the id type from string to int
    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }
    fmt.Println(err)
    
    fmt.Println(r.Body,"B")
    // create an empty campaign of type models.User
    var campaign models.Campaign
    fmt.Println(campaign,"CAM")
    // decode the json request to campaign
    err = json.NewDecoder(r.Body).Decode(&campaign)

    fmt.Println("HEY")

    if err != nil {
        log.Println(err)
    }
    

    // call update campaign to update the campaign
    updatedRows := updateCampaign(int64(id), campaign)

    // format the message string
    msg := fmt.Sprintf("Campaign updated successfully. Total rows/record affected %v", updatedRows)

    // format the response message
    res := response{
        ID:      int64(id),
        Message: msg,
    }

    // send the response
    json.NewEncoder(w).Encode(res)
}



// update campaign in the DB
func updateCampaign(id int64, campaign models.Campaign) int64 {

    // create the postgres db connection
    db := createConnection()

    // close the db connection
    defer db.Close()

    // create the update sql query
    sqlStatement := `UPDATE campaigns SET name=$2, status=$3, type=$4, budget=$5 WHERE id=$1`
    fmt.Println(campaign)
    // execute the sql statement
    res, err := db.Exec(sqlStatement, id, campaign.Name, campaign.Status, campaign.Type, campaign.Budget)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    // check how many rows affected
    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }

    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}