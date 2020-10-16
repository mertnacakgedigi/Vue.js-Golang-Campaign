package controllers

import (
    "database/sql"
    "encoding/json" 
    "fmt"
    "backend/models" 
    "log"
    "os" 
    "net/http" 
    "github.com/joho/godotenv" 
    "strconv"  
    "github.com/gorilla/mux" 
    _ "github.com/lib/pq"      
)


type response struct {
    ID      int64  `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}

func createConnection() *sql.DB {
    
    err := godotenv.Load(".env")

    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

    if err != nil {
        panic(err)
    }

    err = db.Ping()

    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully connected!")
    
    return db
}

func GetAllCampaigns(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
 
    campaigns, err := getAllCampaigns()

    if err != nil {
        log.Fatalf("Unable to get all campaign. %v", err)
    }
    json.NewEncoder(w).Encode(campaigns)
}

func getAllCampaigns() ([]models.Campaign, error) {
    
    db := createConnection()

    defer db.Close()

    var campaigns []models.Campaign

    sqlStatement := `SELECT * FROM campaigns`
  
    rows, err := db.Query(sqlStatement)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }
      
    defer rows.Close()

    for rows.Next() {
        var campaign models.Campaign

        err = rows.Scan(&campaign.Id, &campaign.Name, &campaign.Status,&campaign.Type, &campaign.Budget, &campaign.Created_at)

        if err != nil {
            log.Fatalf("Unable to scan the row. %v", err)
        }

        campaigns = append(campaigns, campaign)

	}  
    return campaigns, err
}


func CreateCampaign(w http.ResponseWriter, r *http.Request) {
  
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST")
    w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

    var campaign models.Campaign
    
    err := json.NewDecoder(r.Body).Decode(&campaign)

    if err != nil {
        log.Fatalf("Unable to decode the request body.  %v", err)
    }
    
    insertID := insertCampaign(campaign)

    
    res := response{
        ID:      insertID,
        Message: "Campaign created successfully",
    }  
    json.NewEncoder(w).Encode(res)
}

func insertCampaign(campaign models.Campaign) int64 {
   
    db := createConnection()
  
    defer db.Close()

    sqlStatement := `INSERT INTO campaigns (id,name, status,type,budget,created_at) VALUES ($1, $2,$3,$4,$5,$6) RETURNING id`

    var id int64

    err := db.QueryRow(sqlStatement, campaign.Id,campaign.Name, campaign.Status,campaign.Type, campaign.Budget,campaign.Created_at).Scan(&id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }

    fmt.Printf("Inserted a single record %v", id)

    return id
}


func DeleteCampaign(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    params := mux.Vars(r)

    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    deletedRows := deleteCampaign(int64(id))

    msg := fmt.Sprintf("Campaign deleted successfully. Total rows/record affected %v", deletedRows)

    res := response{
        ID:      int64(id),
        Message: msg,
    }
    json.NewEncoder(w).Encode(res)
}


func deleteCampaign(id int64) int64 {

    db := createConnection()

    defer db.Close()
 
    sqlStatement := `DELETE FROM campaigns WHERE id=$1`

    res, err := db.Exec(sqlStatement, id)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }
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

    params := mux.Vars(r)

    id, err := strconv.Atoi(params["id"])

    if err != nil {
        log.Fatalf("Unable to convert the string into int.  %v", err)
    }

    var campaign models.Campaign
  
    err = json.NewDecoder(r.Body).Decode(&campaign)

    if err != nil {
        log.Println(err)
    }

    updatedRows := updateCampaign(int64(id), campaign)

    msg := fmt.Sprintf("Campaign updated successfully. Total rows/record affected %v", updatedRows)

    
    res := response{
        ID:      int64(id),
        Message: msg,
    } 
    json.NewEncoder(w).Encode(res)
}


func updateCampaign(id int64, campaign models.Campaign) int64 {

    db := createConnection()

    defer db.Close()

    sqlStatement := `UPDATE campaigns SET name=$2, status=$3, type=$4, budget=$5 WHERE id=$1`
 
    res, err := db.Exec(sqlStatement, id, campaign.Name, campaign.Status, campaign.Type, campaign.Budget)

    if err != nil {
        log.Fatalf("Unable to execute the query. %v", err)
    }
     
    rowsAffected, err := res.RowsAffected()

    if err != nil {
        log.Fatalf("Error while checking the affected rows. %v", err)
    }

    fmt.Printf("Total rows/record affected %v", rowsAffected)

    return rowsAffected
}