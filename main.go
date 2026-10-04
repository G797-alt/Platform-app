package main
import (
"database/sql"
"encoding/csv"
"encoding/json"
"fmt"
"html"
"log"
"net/http"
"os"
"strconv"
"time"
_ "github.com/lib/pq"
)
var db *sql.DB
func main(){
 var err error
 db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
 if err!=nil{log.Fatal(err)}
 db.Ping()
 http.HandleFunc("/", home)
 http.HandleFunc("/save", save)
 http.HandleFunc("/save-mwanachama", saveMwanachama)
 http.HandleFunc("/edit", edit)
 http.HandleFunc("/del", del)
 http.HandleFunc("/backup", backup)
 http.HandleFunc("/csv", csvExport)
 http.HandleFunc("/restore", restore)
 p:=os.Getenv("PORT")
 if p==""{p="10000"}
 log.Println("CJ1V v2 FIXED LIVE")
 http.ListenAndServe(":"+p, nil)
}
func fmtInt(n int64) string {
 if n==0{return "0"}
 s:=fmt.Sprintf("%d", n)
 out:=""
 c:=0
 for i:=len(s)-1;i>=0;i--{
  out=string(s[i])+out
  c++
  if c%3==0 && i!=0{out=","+out}
 }
 return out
}
func home(w http.ResponseWriter, r *http.Request){
 w.Header().Set("Content-Type","text/html; charset=utf-8")
 var mapato, matumizi int64
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mahida").Scan(&mapato)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%mkopo%' AND aina NOT ILIKE '%rejesho%'").Scan(&matumizi)
 baki:=mapato-matumizi
 q:=r.URL.Query().Get("q")
 opt:=""
 rows,_:=db.Query("SELECT id,jina FROM wanachama ORDER BY jina ASC")
 if rows!=nil{
  for rows.Next(){
   var id int
   var j sql.NullString
   rows.Scan(&id,&j)
   jv:=""
   if j.Valid{jv=j.String}
   opt+=fmt.Sprintf("<option value='%d|%s'>%s</option>", id, html.EscapeString(jv), html.EscapeString(jv))
  }
  rows.Close()
 }
 hist:=""
 sqlQ:="SELECT id,jina,aina,kiasi FROM mahida "
 var r2 *sql.Rows
 var err error
 if q!=""{
  r2,err=db.Query(sqlQ+" WHERE jina ILIKE $1 OR aina ILIKE $1 ORDER BY id DESC LIMIT 300", "%"+q+"%")
 }else{
  r2,err=db.Query(sqlQ+" ORDER BY id DESC LIMIT 300")
 }
 if err!=nil{
  hist=fmt.Sprintf("<tr><td colspan=6 style='color:red'>DB ERROR: %s</td></tr>", html.EscapeString(err.Error()))
 }else if r2!=nil{
  cnt:=0
  for r2.Next(){
   var id sql.NullInt64
   var j,a sql.NullString
   var k sql.NullInt64
   if err:=r2.Scan(&id,&j,&a,&k); err!=nil{log.Println(err); continue}
   cnt++
   idv:="0"; if id.Valid{idv=strconv.FormatInt(id.Int64,10)}
   jv:=""; if j.Valid{jv=j.String}
   av:=""; if a.Valid{av=a.String}
   kv:=int64(0); if k.Valid{kv=k.Int64}
   hist+=fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td><button onclick=\"editD(%s,%d)\">Edit</button></td><td><a href='/del?id=%s' onclick=\"return confirm('Futa?')\">X</a></td></tr>", idv, html.EscapeString(jv), html.EscapeString(av), fmtInt(kv), idv, kv, idv)
  }
  r2.Close()
  if cnt==0{
   hist=fmt.Sprintf("<tr><td colspan=6 style='text-align:center;padding:20px'>Hakuna data (q=%s) - MAPATO %s</td></tr>", html.EscapeString(q), fmtInt(mapato))
  }
 }
 totals:=""
 rt,_:=db.Query("SELECT aina, SUM(kiasi) FROM mahida GROUP BY aina ORDER BY SUM(kiasi) DESC")
 if rt!=nil{
  for rt.Next(){
   var a sql.NullString
   var s sql.NullInt64
   rt.Scan(&a,&s)
   av:=""; if a.Valid{av=a.String}
   sv:=int64(0); if s.Valid{sv=s.Int64}
   totals+=fmt.Sprintf("<span style='background:#e8f5e9;padding:6px 10px;border-radius:20px;margin:2px;display:inline-block;font-size:12px'>%s: %s</span>", html.EscapeString(av), fmtInt(sv))
  }
  rt.Close()
 }
 fmt.Fprint(w, "<html><head><meta charset='utf-8'><meta name='viewport' content='width=device-width,initial-scale=1'><title>CJ1V v2 FIXED</title>")
 fmt.Fprint(w, "<style>*{box-sizing:border-box}body{font-family:sans-serif;background:#f0f4f0;margin:0;padding:8px}.top{display:grid;grid-template-columns:1fr 1fr 1fr;gap:8px;margin-bottom:12px}.card{padding:14px;border-radius:16px;color:#fff;text-align:center;font-weight:bold}.g{background:#2e7d32}.r{background:#c62828}.b{background:#1565c0}.box{background:#fff;border-radius:12px;padding:12px;margin-bottom:10px}input,select{padding:10px;border:1px solid #ccc;border-radius:8px;margin:2px}.btn{padding:12px;border:none;border-radius:10px;color:#fff;font-weight:bold;width:100%;margin-top:6px}table{width:100%;border-collapse:collapse;font-size:13px}th{background:#e8f5e9;padding:8px;text-align:left}td{padding:8px;border-bottom:1px solid #eee}</style></head><body>")
 fmt.Fprintf(w, "<div class='top'><div class='card g'>MAPATO<br>TZS %s</div><div class='card r'>MATUMIZI<br>TZS %s</div><div class='card b'>BAKAA<br>TZS %s</div></div>", fmtInt(mapato), fmtInt(matumizi), fmtInt(baki))
 fmt.Fprintf(w, "<div class='box'><div style='display:flex;flex-wrap:wrap'>%s</div></div>", totals)
 fmt.Fprintf(w, "<div class='box'><h3>➕ Ongeza Muamala</h3><form action='/save' method='POST' style='display:flex;flex-wrap:wrap;gap:4px'><select name='mid' required style='flex:1'>%s</select><select name='aina'><option>Hisa</option><option>Ada</option><option>Faini</option><option>Fomu</option><option>Mkopo</option><option>Marejesho</option><option>Riba</option></select><input name='kiasi' type='number' placeholder='Kiasi' required style='width:120px'><button class='btn' style='background:#2e7d32;width:auto'>SAVE</button></form><form action='/save-mwanachama' method='POST' style='margin-top:8px;display:flex;gap:4px'><input name='jina' placeholder='Mwanachama mpya' required style='flex:1'><button style='background:#1565c0;color:#fff;border:none;border-radius:8px;padding:10px'>+ Mwanachama</button></form></div>", opt)
 fmt.Fprintf(w, "<div class='box'><form method='GET' style='display:flex;gap:4px'><input name='q' value='%s' placeholder='Search...' style='flex:1'><button style='background:#616161;color:#fff;border:none;border-radius:8px;padding:10px'>Search</button></form></div>", html.EscapeString(q))
 fmt.Fprintf(w, "<div class='box'><h3>Historia LIVE (%s TZS)</h3><table><tr><th>ID</th><th>Jina</th><th>Aina</th><th>Kiasi</th><th>Edit</th><th>Del</th></tr>%s</table></div>", fmtInt(mapato), hist)
 fmt.Fprint(w, "<div class='box'><h3>Backup</h3><a href='/backup'><button class='btn' style='background:#ef6c00'>Pakua JSON</button></a><a href='/csv'><button class='btn' style='background:#616161'>Pakua CSV</button></a><form action='/restore' method='POST' enctype='multipart/form-data' style='margin-top:8px'><input type='file' name='file' required><button class='btn' style='background:#2e7d32'>Restore</button></form></div>")
 fmt.Fprint(w, "<script>function editD(id,oldK){var k=prompt('Edit kiasi ID '+id, oldK); if(k) fetch('/edit?id='+id+'&k='+k,{method:'POST'}).then(()=>location.reload());}</script></body></html>")
}
func save(w http.ResponseWriter, r *http.Request){
 r.ParseForm()
 midFull:=r.FormValue("mid")
 k,_:=strconv.ParseInt(r.FormValue("kiasi"),10,64)
 aina:=r.FormValue("aina")
 jina:=midFull; id:=0
 for i,c:=range midFull{if c=='|' {id,_=strconv.Atoi(midFull[:i]); jina=midFull[i+1:]; break}}
 if jina==""{jina="Haijulikani"}
 db.Exec("INSERT INTO mahida(mwanachama_id,jina,kiasi,aina,mwezi,mwaka) VALUES($1,$2,$3,$4,$5,$6)", id, jina, k, aina, int(time.Now().Month()), time.Now().Year())
 http.Redirect(w,r,"/",302)
}
func saveMwanachama(w http.ResponseWriter, r *http.Request){
 r.ParseForm()
 j:=r.FormValue("jina")
 if j!=""{db.Exec("INSERT INTO wanachama(jina) VALUES($1)", j)}
 http.Redirect(w,r,"/",302)
}
func edit(w http.ResponseWriter, r *http.Request){
 id:=r.URL.Query().Get("id")
 k:=r.URL.Query().Get("k")
 kv,_:=strconv.ParseInt(k,10,64)
 db.Exec("UPDATE mahida SET kiasi=$1 WHERE id=$2", kv, id)
 w.Write([]byte("ok"))
}
func del(w http.ResponseWriter, r *http.Request){
 id:=r.URL.Query().Get("id")
 db.Exec("DELETE FROM mahida WHERE id=$1", id)
 http.Redirect(w,r,"/",302)
}
func backup(w http.ResponseWriter, r *http.Request){
 rows,_:=db.Query("SELECT id,mwanachama_id,jina,aina,kiasi,mwezi,mwaka FROM mahida ORDER BY id")
 var arr []map[string]interface{}
 if rows!=nil{
  for rows.Next(){
   var id,mid,mw,yr sql.NullInt64
   var j,a sql.NullString
   var k sql.NullInt64
   rows.Scan(&id,&mid,&j,&a,&k,&mw,&yr)
   arr=append(arr, map[string]interface{}{"id":id.Int64,"mid":mid.Int64,"jina":j.String,"aina":a.String,"kiasi":k.Int64,"mwezi":mw.Int64,"mwaka":yr.Int64})
  }
  rows.Close()
 }
 w.Header().Set("Content-Disposition","attachment; filename=cj1v-backup.json")
 json.NewEncoder(w).Encode(arr)
}
func csvExport(w http.ResponseWriter, r *http.Request){
 w.Header().Set("Content-Disposition","attachment; filename=cj1v.csv")
 w.Header().Set("Content-Type","text/csv")
 cw:=csv.NewWriter(w)
 cw.Write([]string{"ID","Jina","Aina","Kiasi"})
 rows,_:=db.Query("SELECT id,jina,aina,kiasi FROM mahida ORDER BY id DESC")
 if rows!=nil{
  for rows.Next(){
   var id sql.NullInt64
   var j,a sql.NullString
   var k sql.NullInt64
   rows.Scan(&id,&j,&a,&k)
   cw.Write([]string{strconv.FormatInt(id.Int64,10), j.String, a.String, strconv.FormatInt(k.Int64,10)})
  }
  rows.Close()
 }
 cw.Flush()
}
func restore(w http.ResponseWriter, r *http.Request){
 r.ParseMultipartForm(10<<20)
 f,_,_:=r.FormFile("file")
 if f==nil{http.Redirect(w,r,"/",302); return}
 var arr []map[string]interface{}
 json.NewDecoder(f).Decode(&arr)
 for _,x:=range arr{
  jina:=""; if v,ok:=x["jina"].(string); ok{jina=v}
  aina:=""; if v,ok:=x["aina"].(string); ok{aina=v}
  kiasi:=int64(0); if v,ok:=x["kiasi"].(float64); ok{kiasi=int64(v)}
  mid:=0; if v,ok:=x["mid"].(float64); ok{mid=int(v)}
  db.Exec("INSERT INTO mahida(mwanachama_id,jina,kiasi,aina,mwezi,mwaka) VALUES($1,$2,$3,$4,$5,$6)", mid, jina, kiasi, aina, int(time.Now().Month()), time.Now().Year())
 }
 http.Redirect(w,r,"/",302)
}
