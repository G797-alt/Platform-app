package main
import (
 "database/sql"
 "encoding/json"
 "fmt"
 "log"
 "net/http"
 "os"
 "strconv"
 _ "github.com/jackc/pgx/v5/stdlib"
)
var db *sql.DB
type Muamala struct{ ID int; Mwezi int; Mwaka int; Jina string; Aina string; Kiasi int64; Maelezo string }

func main(){
 var err error
 db, err = db, err = sql.Open("pgx", os.Getenv("DATABASE_URL"))
 if err!=nil{fmt.Println("sqlite err",err)}
 db.Exec(`CREATE TABLE IF NOT EXISTS wanachama (id INTEGER PRIMARY KEY AUTOINCREMENT, jina TEXT UNIQUE NOT NULL)`)
 db.Exec(`CREATE TABLE IF NOT EXISTS miamala (id INTEGER PRIMARY KEY AUTOINCREMENT, mwezi INT, mwaka INT, jina TEXT, aina TEXT, kiasi BIGINT, maelezo TEXT)`)
db.Exec(`INSERT INTO wanachama(jina) VALUES('Test') ON CONFLICT DO NOTHING`)
 fmt.Println("DB OK - Postgres")
 http.HandleFunc("/api/wanachama",handleWanachama)
 http.HandleFunc("/api/miamala",handleMiamala)
 http.HandleFunc("/api/summary",handleSummary)
 http.Handle("/",http.FileServer(http.Dir(".")))
 p:=os.Getenv("PORT"); if p==""{p="8081"}
 fmt.Printf("Mahida 2023-2035 LIVE on http://localhost:%s\n",p)
 log.Fatal(http.ListenAndServe(":"+p,nil))
}
func handleWanachama(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","application/json"); w.Header().Set("Access-Control-Allow-Origin","*"); w.Header().Set("Access-Control-Allow-Methods","GET,POST,DELETE,OPTIONS"); w.Header().Set("Access-Control-Allow-Headers","Content-Type")
 if r.Method=="OPTIONS"{return}
 if r.Method=="GET"{
  rows,_:=db.Query("SELECT id,jina FROM wanachama ORDER BY jina"); defer rows.Close()
  var l []map[string]interface{};
  if rows!=nil{for rows.Next(){var id int; var j string; rows.Scan(&id,&j); l=append(l,map[string]interface{}{"id":id,"jina":j})}}
  if l==nil{l=[]map[string]interface{}{}}; json.NewEncoder(w).Encode(l);return
 }
 if r.Method=="POST"{var m map[string]string; json.NewDecoder(r.Body).Decode(&m); j:=m["jina"]; db.Exec("INSERT OR IGNORE INTO wanachama(jina) VALUES(?)",j); var id int; db.QueryRow("SELECT id FROM wanachama WHERE jina=?",j).Scan(&id); json.NewEncoder(w).Encode(map[string]interface{}{"id":id,"jina":j});return}
 if r.Method=="DELETE"{id:=r.URL.Query().Get("id"); db.Exec("DELETE FROM wanachama WHERE id=?",id); json.NewEncoder(w).Encode(map[string]string{"status":"deleted"})}
}
func handleMiamala(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","application/json"); w.Header().Set("Access-Control-Allow-Origin","*"); w.Header().Set("Access-Control-Allow-Methods","GET,POST,DELETE,OPTIONS"); w.Header().Set("Access-Control-Allow-Headers","Content-Type")
 if r.Method=="OPTIONS"{return}
 if r.Method=="GET"{
  mwezi,_:=strconv.Atoi(r.URL.Query().Get("mwezi")); mwaka,_:=strconv.Atoi(r.URL.Query().Get("mwaka"))
  var rows *sql.Rows
  if mwezi>0{rows,_=db.Query("SELECT id,mwezi,mwaka,jina,aina,kiasi,COALESCE(maelezo,'') FROM miamala WHERE mwezi=? AND mwaka=? ORDER BY id DESC",mwezi,mwaka)}else{rows,_=db.Query("SELECT id,mwezi,mwaka,jina,aina,kiasi,COALESCE(maelezo,'') FROM miamala ORDER BY id DESC LIMIT 500")}
  if rows!=nil{defer rows.Close()}; var l []Muamala;
  if rows!=nil{for rows.Next(){var m Muamala; rows.Scan(&m.ID,&m.Mwezi,&m.Mwaka,&m.Jina,&m.Aina,&m.Kiasi,&m.Maelezo); l=append(l,m)}}
  if l==nil{l=[]Muamala{}}; json.NewEncoder(w).Encode(l);return
 }
 if r.Method=="POST"{var m Muamala; json.NewDecoder(r.Body).Decode(&m); res,_:=db.Exec("INSERT INTO miamala(mwezi,mwaka,jina,aina,kiasi,maelezo) VALUES(?,?,?,?,?,?)",m.Mwezi,m.Mwaka,m.Jina,m.Aina,m.Kiasi,m.Maelezo); id,_:=res.LastInsertId(); json.NewEncoder(w).Encode(map[string]int64{"id":id});return}
 if r.Method=="DELETE"{id:=r.URL.Query().Get("id"); db.Exec("DELETE FROM miamala WHERE id=?",id); json.NewEncoder(w).Encode(map[string]string{"status":"deleted"})}
}
func handleSummary(w http.ResponseWriter,r *http.Request){
 w.Header().Set("Content-Type","application/json"); w.Header().Set("Access-Control-Allow-Origin","*")
 mwezi,_:=strconv.Atoi(r.URL.Query().Get("mwezi")); mwaka,_:=strconv.Atoi(r.URL.Query().Get("mwaka"))
 var a,b int64
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE mwezi=? AND mwaka=? AND aina='Mapato'",mwezi,mwaka).Scan(&a)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE mwezi=? AND mwaka=? AND aina='Matumizi'",mwezi,mwaka).Scan(&b)
 json.NewEncoder(w).Encode(map[string]int64{"mapato":a,"matumizi":b,"bakaa":a-b})
}
