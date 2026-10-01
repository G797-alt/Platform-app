package main

import (
 "database/sql"
 "fmt"
 "html"
 "log"
 "net/http"
 "os"
 "strconv"
 "strings"
 _ "github.com/lib/pq"
)

var db *sql.DB

func main(){
 dsn := os.Getenv("DATABASE_URL")
 if dsn==""{ log.Fatal("DATABASE_URL missing") }
 var err error
 db, err = sql.Open("postgres", dsn)
 if err!=nil{ log.Fatal(err) }
 db.Exec(`CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY, jina TEXT)`)
 db.Exec(`CREATE TABLE IF NOT EXISTS mahida(id SERIAL PRIMARY KEY, mwanachama_id INT, jina TEXT, kiasi BIGINT, aina TEXT, maelezo TEXT, mwezi INT, mwaka INT)`)
 db.Exec(`CREATE TABLE IF NOT EXISTS mikopo_chukua(id SERIAL PRIMARY KEY, mwanachama_id INT, kiasi BIGINT, mwezi INT, mwaka INT)`)
 db.Exec(`CREATE TABLE IF NOT EXISTS mikopo_rudisha(id SERIAL PRIMARY KEY, mwanachama_id INT, kiasi BIGINT, mwezi INT, mwaka INT)`)

 http.HandleFunc("/", home)
 http.HandleFunc("/save", save)
 http.HandleFunc("/addM", addM)
 http.HandleFunc("/delM", delM)
 http.HandleFunc("/chukua", chukua)
 http.HandleFunc("/rudisha", rudisha)

 p := os.Getenv("PORT")
 if p==""{ p="10000" }
 http.ListenAndServe(":"+p, nil)
}

func home(w http.ResponseWriter, r *http.Request){
 w.Header().Set("Content-Type","text/html; charset=utf-8")
 mz:=r.URL.Query().Get("mwezi")
 if mz==""{ mz="Okt" }
 mwS:=r.URL.Query().Get("mwaka")
 if mwS==""{ mwS="2026" }
 mw,_:=strconv.Atoi(mwS)
 mp:=map[string]int{"Jan":1,"Feb":2,"Mar":3,"Apr":4,"Mei":5,"Jun":6,"Jul":7,"Ago":8,"Sep":9,"Okt":10,"Nov":11,"Des":12}
 mNum:=mp[mz]
 if mNum==0{ mNum=10 }

 var mapato, matumiziMikopo, mapatoMikopo int64
 type W struct{ID int; Jina string}
 var list []W
 if db!=nil{
  db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE mwezi=$1 AND mwaka=$2`, mNum, mw).Scan(&mapato)
  db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mikopo_rudisha WHERE mwezi=$1 AND mwaka=$2`, mNum, mw).Scan(&mapatoMikopo)
  db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mikopo_chukua WHERE mwezi=$1 AND mwaka=$2`, mNum, mw).Scan(&matumiziMikopo)
  rows,_:=db.Query("SELECT id,jina FROM wanachama ORDER BY id")
  if rows!=nil{ defer rows.Close(); for rows.Next(){ var id int; var j string; rows.Scan(&id,&j); list=append(list,W{id,j}) } }
 }
 totalMapato := mapato + mapatoMikopo
 totalMatumizi := matumiziMikopo
 baki:= totalMapato - totalMatumizi

 mo:=""; for _,o:=range []string{"Jan","Feb","Mar","Apr","Mei","Jun","Jul","Ago","Sep","Okt","Nov","Des"}{s:=""; if o==mz{ s="selected" }; mo+=fmt.Sprintf("<option %s>%s</option>",s,o)}
 yo:=""; for y:=2023;y<=2035;y++{s:=""; if strconv.Itoa(y)==mwS{ s="selected" }; yo+=fmt.Sprintf("<option %s>%d</option>",s,y)}
 mop:=""; for _,wm:=range list{ mop+=fmt.Sprintf("<option value=%d>%s</option>",wm.ID, html.EscapeString(wm.Jina)) }

 rudishaHTML:=""; deniHTML:=""; histHTML:=""
 if db!=nil{
  rows,_:=db.Query(`SELECT w.jina, SUM(r.kiasi) FROM mikopo_rudisha r JOIN wanachama w ON w.id=r.mwanachama_id WHERE r.mwezi=$1 AND r.mwaka=$2 GROUP BY w.jina`, mNum, mw)
  if rows!=nil{ defer rows.Close(); for rows.Next(){ var jina string; var k int64; rows.Scan(&jina,&k); rudishaHTML+=fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>", html.EscapeString(jina), k) } }
  rows2,_:=db.Query(`SELECT w.jina, COALESCE(SUM(c.kiasi),0)-COALESCE(SUM(r.kiasi),0) FROM wanachama w LEFT JOIN mikopo_chukua c ON c.mwanachama_id=w.id AND c.mwezi=$1 AND c.mwaka=$2 LEFT JOIN mikopo_rudisha r ON r.mwanachama_id=w.id AND r.mwezi=$1 AND r.mwaka=$2 GROUP BY w.jina HAVING COALESCE(SUM(c.kiasi),0)-COALESCE(SUM(r.kiasi),0)>0`, mNum, mw)
  if rows2!=nil{ defer rows2.Close(); for rows2.Next(){ var jina string; var d int64; rows2.Scan(&jina,&d); deniHTML+=fmt.Sprintf("<tr><td>%s</td><td style='color:red'>%d</td></tr>", html.EscapeString(jina), d) } }
  rows3,_:=db.Query(`SELECT w.jina,m.aina,m.kiasi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id WHERE m.mwezi=$1 AND m.mwaka=$2 UNION ALL SELECT w.jina,'Chukua',c.kiasi FROM mikopo_chukua c JOIN wanachama w ON w.id=c.mwanachama_id WHERE c.mwezi=$1 AND c.mwaka=$2 UNION ALL SELECT w.jina,'Rudisha',r.kiasi FROM mikopo_rudisha r JOIN wanachama w ON w.id=r.mwanachama_id WHERE r.mwezi=$1 AND r.mwaka=$2`, mNum, mw, mNum, mw, mNum, mw)
  if rows3!=nil{ defer rows3.Close(); for rows3.Next(){ var j,a string; var k int64; rows3.Scan(&j,&a,&k); histHTML+=fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%d</td></tr>", html.EscapeString(j), html.EscapeString(a), k) } }
 }

 fmt.Fprintf(w, `<!DOCTYPE html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><title>Kikoba</title><style>body{font-family:sans-serif;padding:10px}.card{border:1px solid #ddd;padding:12px;margin:8px 0;border-radius:8px} table{width:100%%;border-collapse:collapse} th,td{border:1px solid #ccc;padding:6px} input{width:90px}</style></head><body>
<h2>Mapato: %d | Matumizi: %d | Bakaa: %d</h2>
<p>IN=Mapato yote, OUT=Mikopo, Bakaa=IN-OUT</p>
<form> Mwezi <select name=mwezi onchange=this.form.submit()>%s</select> Mwaka <select name=mwaka onchange=this.form.submit()>%s</select></form>
<div class=card><h3>1. Walio Rudisha (LIVE)</h3><table><tr><th>Jina</th><th>Kiasi</th></tr>%s</table></div>
<div class=card><h3>2. Wenye Deni (LIVE)</h3><table><tr><th>Jina</th><th>Deni</th></tr>%s</table></div>
<div class=card><h3>3. Historia ya Mwezi (LIVE)</h3><table><tr><th>Jina</th><th>Aina</th><th>Kiasi</th></tr>%s</table></div>
<div class=card>
<form action=/save method=POST><select name=mwanachama_id>%s</select><select name=aina><option>Hisa</option><option>Ada</option><option>Faini</option><option>Fomu</option></select><input name=kiasi type=number><input type=hidden name=mwezi value=%s><input type=hidden name=mwaka value=%s><input type=hidden name=mweziNum value=%d><input type=hidden name=mwakaNum value=%d><button>Save</button></form>
<form action=/chukua method=POST><select name=mwanachama_id>%s</select><input name=kiasi type=number><input type=hidden name=mwezi value=%s><input type=hidden name=mwaka value=%s><input type=hidden name=mweziNum value=%d><input type=hidden name=mwakaNum value=%d><button>Chukua Mkopo</button></form>
<form action=/rudisha method=POST><select name=mwanachama_id>%s</select><input name=kiasi type=number><input type=hidden name=mwezi value=%s><input type=hidden name=mwaka value=%s><input type=hidden name=mweziNum value=%d><input type=hidden name=mwakaNum value=%d><button>Rudisha</button></form>
</div>
</body></html>`, totalMapato, totalMatumizi, baki, mo, yo, rudishaHTML, deniHTML, histHTML, mop, mz, mwS, mNum, mw, mop, mz, mwS, mNum, mw, mop, mz, mwS, mNum, mw)
}

func addM(w http.ResponseWriter, r *http.Request){
 r.ParseForm()
 j:=strings.TrimSpace(r.FormValue("jina"))
 mz:=r.FormValue("mwezi"); mw:=r.FormValue("mwaka")
 if j!=""{ db.Exec("INSERT INTO wanachama(jina) VALUES($1)", j) }
 http.Redirect(w,r, "/?mwezi="+mz+"&mwaka="+mw, 302)
}
func delM(w http.ResponseWriter, r *http.Request){
 r.ParseForm()
 id:=r.FormValue("id")
 if id==""{ id=r.URL.Query().Get("id") }
 db.Exec("DELETE FROM mahida WHERE id=$1", id)
 w.Write([]byte("ok"))
}
func save(w http.ResponseWriter, r *http.Request){
 r.ParseForm()
 mID,_:=strconv.Atoi(r.FormValue("mwanachama_id"))
 aina:=r.FormValue("aina")
 kiasi,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64)
 mNum,_:=strconv.Atoi(r.FormValue("mweziNum"))
 mw,_:=strconv.Atoi(r.FormValue("mwakaNum"))
 mz:=r.FormValue("mwezi")
 mwS:=r.FormValue("mwaka")
 var jina string
 db.QueryRow("SELECT jina FROM wanachama WHERE id=$1", mID).Scan(&jina)
 db.Exec("INSERT INTO mahida(mwanachama_id,jina,kiasi,aina,mwezi,mwaka) VALUES($1,$2,$3,$4,$5,$6)", mID, jina, kiasi, aina, mNum, mw)
 http.Redirect(w,r, "/?mwezi="+mz+"&mwaka="+mwS, 302)
}
func chukua(w http.ResponseWriter, r *http.Request){
 r.ParseForm()
 mID,_:=strconv.Atoi(r.FormValue("mwanachama_id"))
 kiasi,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64)
 mNum,_:=strconv.Atoi(r.FormValue("mweziNum"))
 mw,_:=strconv.Atoi(r.FormValue("mwakaNum"))
 mz:=r.FormValue("mwezi"); mwS:=r.FormValue("mwaka")
 db.Exec("INSERT INTO mikopo_chukua(mwanachama_id,kiasi,mwezi,mwaka) VALUES($1,$2,$3,$4)", mID, kiasi, mNum, mw)
 http.Redirect(w,r, "/?mwezi="+mz+"&mwaka="+mwS, 302)
}
func rudisha(w http.ResponseWriter, r *http.Request){
 r.ParseForm()
 mID,_:=strconv.Atoi(r.FormValue("mwanachama_id"))
 kiasi,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64)
 mNum,_:=strconv.Atoi(r.FormValue("mweziNum"))
 mw,_:=strconv.Atoi(r.FormValue("mwakaNum"))
 mz:=r.FormValue("mwezi"); mwS:=r.FormValue("mwaka")
 db.Exec("INSERT INTO mikopo_rudisha(mwanachama_id,kiasi,mwezi,mwaka) VALUES($1,$2,$3,$4)", mID, kiasi, mNum, mw)
 http.Redirect(w,r, "/?mwezi="+mz+"&mwaka="+mwS, 302)
}