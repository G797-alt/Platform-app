package main
import ("database/sql"; "encoding/csv"; "encoding/json"; "fmt"; "html"; "log"; "net/http"; "os"; "strconv"; "time"; _ "github.com/lib/pq")
var db *sql.DB
func main(){
db,_=sql.Open("postgres",os.Getenv("DATABASE_URL"))
http.HandleFunc("/", home)
http.HandleFunc("/save", save)
http.HandleFunc("/edit", edit)
http.HandleFunc("/del", del)
http.HandleFunc("/backup", backup)
http.HandleFunc("/restore", restore)
http.HandleFunc("/csv", csvExport)
p:=os.Getenv("PORT"); if p==""{p="10000"}
log.Println("Live",p)
http.ListenAndServe(":"+p,nil)
}
func fmtInt(n int64) string{
s:=fmt.Sprintf("%d",n)
if len(s)<=3{return s}
r:=""; c:=0
for i:=len(s)-1;i>=0;i--{
r=string(s[i])+r; c++
if c%3==0 && i!=0{r=","+r}
}
return r
}
func home(w http.ResponseWriter, r *http.Request){
w.Header().Set("Content-Type","text/html; charset=utf-8")
var mapato,matumizi int64
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina NOT ILIKE '%mkopo%' OR aina ILIKE '%rejesho%'`).Scan(&mapato)
if mapato==0{db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida`).Scan(&mapato)}
db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM mahida WHERE aina ILIKE '%mkopo%' AND aina NOT ILIKE '%rejesho%'`).Scan(&matumizi)
baki:=mapato-matumizi
opt:=""; r1,_:=db.Query(`SELECT id,jina FROM wanachama ORDER BY jina`)
if r1!=nil{for r1.Next(){var id int; var j string; r1.Scan(&id,&j); opt+=fmt.Sprintf(`<option value=%d>%s</option>`,id,html.EscapeString(j))}; r1.Close()}
rud:=""; r2,_:=db.Query(`SELECT m.id,w.jina,m.kiasi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id WHERE m.aina ILIKE '%rejesho%' ORDER BY m.id DESC`)
if r2!=nil{for r2.Next(){var id int; var j string; var k int64; r2.Scan(&id,&j,&k); rud+=fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td><button onclick="editD(%d,%d)">✏️</button> <button onclick="delD(%d)">X</button></td></tr>`,html.EscapeString(j),fmtInt(k),id,k,id)}; r2.Close()}
if rud==""{rud=`<tr><td colspan=3>Hakuna</td></tr>`}
deni:=""; r3,_:=db.Query(`SELECT w.jina, COALESCE(SUM(CASE WHEN m.aina ILIKE '%%mkopo%%' AND m.aina NOT ILIKE '%%rejesho%%' THEN m.kiasi ELSE 0 END),0) as ch, COALESCE(SUM(CASE WHEN m.aina ILIKE '%%rejesho%%' THEN m.kiasi ELSE 0 END),0) as ru FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id GROUP BY w.jina HAVING SUM(CASE WHEN m.aina ILIKE '%%mkopo%%' AND m.aina NOT ILIKE '%%rejesho%%' THEN m.kiasi ELSE 0 END)>0`)
if r3!=nil{for r3.Next(){var j string; var ch,ru int64; r3.Scan(&j,&ch,&ru); b:=ch-ru; deni+=fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td style='color:#c62828;font-weight:bold'>%s</td><td>✏️</td></tr>`,html.EscapeString(j),fmtInt(ch),fmtInt(ru),fmtInt(b))}; r3.Close()}
if deni==""{deni=`<tr><td colspan=5>Hakuna deni</td></tr>`}
hist:=""; curM:=int(time.Now().Month()); r4,_:=db.Query(`SELECT m.id,w.jina,m.aina,m.kiasi FROM mahida m JOIN wanachama w ON w.id=m.mwanachama_id WHERE m.mwezi=$1 ORDER BY m.id DESC`,curM)
if r4!=nil{for r4.Next(){var id int; var j,a string; var k int64; r4.Scan(&id,&j,&a,&k); hist+=fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%s</td><td><button onclick="editD(%d,%d)">✏️</button> <button onclick="delD(%d)">❌</button></td></tr>`,html.EscapeString(j),html.EscapeString(a),fmtInt(k),id,k,id)}; r4.Close()}
if hist==""{hist=`<tr><td colspan=4>Hakuna mwezi huu</td></tr>`}
fmt.Fprintf(w,`<!DOCTYPE html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><style>body{font-family:sans-serif;margin:0;background:#f1f5f1;padding:8px}.card{background:#fff;border-radius:12px;padding:12px;margin-bottom:10px}.top{display:flex;gap:6px;margin-bottom:10px}.c{flex:1;padding:12px;border-radius:10px;color:#fff;text-align:center;font-weight:bold;font-size:13px}.g{background:#2e7d32}.r{background:#c62828}.b{background:#1565c0} h3{color:#2e7d32;margin:8px 0;font-size:15px} table{width:100%%;font-size:13px} th{background:#e8f5e9;padding:8px} td{padding:8px;border-bottom:1px solid #eee}.btn{width:100%%;padding:14px;border:none;border-radius:10px;color:#fff;font-weight:bold;margin-bottom:8px}</style></head><body>
<div class=top><div class='c g'>MAPATO<br>TZS %s</div><div class='c r'>MATUMIZI<br>TZS %s</div><div class='c b'>Bakaa<br>TZS %s</div></div>
<div class=card><form action=/save method=POST style="display:flex;gap:4px">