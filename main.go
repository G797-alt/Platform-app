package main
import("database/sql";"fmt";"net/http";"os";_ "github.com/jackc/pgx/v5/stdlib")
var db *sql.DB
func main(){
db,_=sql.Open("pgx",os.Getenv("DATABASE_URL"))
db.Exec(`CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT UNIQUE)`)
db.Exec(`CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,mwezi INT,mwaka INT,jina TEXT,aina TEXT,kiasi BIGINT,maelezo TEXT)`)
http.HandleFunc("/",func(w http.ResponseWriter,r *http.Request){w.Write([]byte("Mahida OK - Postgres"))})
p:=os.Getenv("PORT");if p==""{p="10000"};fmt.Println("DB OK - Postgres");http.ListenAndServe(":"+p,nil)
}