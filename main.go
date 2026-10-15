package main

import (
 "database/sql"
 "encoding/csv"
 "fmt"
 "log"
 "net/http"
 "net/url"
 "os"
 "strconv"
 "time"
 _ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func main() {
 if u := os.Getenv("DATABASE_URL"); u!= "" {
  db, _ = sql.Open("pgx", u)
  db.Exec("CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)")
  db.Exec("CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT)")
  db.Exec("ALTER TABLE miamala ADD COLUMN IF NOT EXISTS tarehe TIMESTAMP DEFAULT NOW()")
  db.Exec("ALTER TABLE miamala ADD COLUMN IF NOT EXISTS riba INT DEFAULT 0")
  db.Exec("ALTER TABLE miamala ADD COLUMN IF NOT EXISTS tarehe_mwisho TIMESTAMP")
  log.Println("V12.1 Riba+Due ready")
 }
 http.HandleFunc("/", home)
 http.HandleFunc("/ongeza", ongeza)
 http.HandleFunc("/mwanachama", mwanachama)
 http.HandleFunc("/export", exportCSV)
 http.HandleFunc("/edit", edit)
 http.HandleFunc("/delete", del)
 p := os.Getenv("PORT")
 if p == "" { p = "10000" }
 http.ListenAndServe("0.0.0.0:"+p, nil)
}

func getTotals() (mapato, dharura, matumizi, mikopo, marejesho, faida_riba int64, count int) {
 if db == nil { return }
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina IN ('hisa_anzia','ada_mwezi','faini','fedha_fomu','michango_mingine')").Scan(&mapato)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='fedha_dharura'").Scan(&dharura)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='matumizi'").Scan(&matumizi)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='mkopo'").Scan(&mikopo)
 db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='marejesho'").Scan(&marejesho)
 db.QueryRow("SELECT COALESCE(SUM(kiasi*riba/100),0) FROM miamala WHERE aina='mkopo'").Scan(&faida_riba)
 db.QueryRow("SELECT COUNT(*) FROM miamala").Scan(&count)
 return
}

func home(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Cache-Control", "no-cache")
 w.Header().Set("Content-Type", "text/html; charset=utf-8")
 mapato, dharura, matumizi, mikopo, marejesho, faida_riba, count := getTotals()
 bakaa := (mapato + dharura + marejesho + faida_riba) - (matumizi + mikopo)
 mwaka := time.Now().Year()
 waLink := "https://wa.me/?text=" + url.QueryEscape(fmt.Sprintf("MAHIDA %d%0ABAKAA %d%0AMikopo %d Riba %d", mwaka, bakaa, mikopo, faida_riba))
 msg := ""
 if r.URL.Query().Get("add") == "ok" { msg = "<div class=card style=background:#e8f5e9;border-left:4px solid green>✅ Imehifadhiwa</div>" }
 if r.URL.Query().Get("edit") == "ok" { msg = "<div class=card style=background:#e8f5e9;border-left:4px solid green>✅ Imehaririwa</div>" }
 today := time.Now().Format("2006-01-02")
 dueDefault := time.Now().AddDate(0, 1, 0).Format("2006-01-02")

 fmt.Fprint(w, "<html><head><meta name=viewport content=\"width=device-width,initial-scale=1\"><style>body{font-family:sans-serif;background:#f6f7f6;margin:0;padding:10px;padding-bottom:90px}.header{padding:16px 6px;font-weight:800;font-size:20px}.card{background:#fff;border-radius:16px;padding:14px;margin-bottom:12px}.bakaa{background:#0f7a3a;color:#fff;border-radius:20px;padding:18px}.big{font-size:26px;font-weight:900}.grid2{display:grid;grid-template-columns:1fr 1fr;gap:10px}.small{font-size:11px;opacity:.9}.big2{font-size:17px;font-weight:800}.btn{width:100%;padding:12px;border-radius:12px;border:0;font-weight:700;margin:6px 0}.btn-green{background:#0f7a3a;color:#fff}.btn-white{background:#fff;border:1px solid #ddd}.btn-wa{background:#25D366;color:#fff}input,select{width:100%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:7px 0}table{width:100%;border-collapse:collapse;font-size:11px}th{padding:8px;background:#222;color:#fff;text-align:left}td{padding:9px;border-bottom:1px solid #eee}.badge{background:#eef2ee;padding:3px 7px;border-radius:20px;font-size:10px}.editbtn{background:#1976d2;color:#fff;padding:6px 12px;border-radius:8px;text-decoration:none;font-weight:700;font-size:11px}.kopaji{border-left:4px solid #ff9800;padding:12px;margin:8px 0;border-radius:10px;background:#fff}.progress{height:8px;background:#eee;border-radius:10px;overflow:hidden;margin:6px 0}.progress-bar{height:100%}.overdue{color:#d32f2f;font-weight:800;font-size:11px}</style><script>function toggleMkopo(){var a=document.getElementById('aina').value;var r=document.getElementById('ribaBox');var d=document.getElementById('dueBox');var j=document.getElementById('jina');if(a=='mkopo'){r.style.display='block';d.style.display='block';j.placeholder='Jina la Mkopaji';}else if(a=='marejesho'){r.style.display='none';d.style.display='none';j.placeholder='Jina lile lile la Mkopaji anayerejesha';}else{r.style.display='none';d.style.display='none';j.placeholder='Jina la Mwanachama';} if(a=='fedha_dharura'||a=='matumizi'){j.value='Washusema Nzima';j.readOnly=true}else{j.readOnly=false;if(j.value=='Washusema Nzima')j.value='';}}</script></head><body>")
 fmt.Fprintf(w, "<div class=header>Uwashusema Mahida %d | %d</div>%s", mwaka, count, msg)
 fmt.Fprintf(w, "<div class=\"card bakaa\"><div class=small>BAKAA %d (Mapato+Dharura+Marejesho+Riba - Matumizi-Mikopo)</div><div class=big>TZS %d</div><div class=small>Mapato %d + Dharura %d + Marejesho %d + Riba %d - Matumizi %d - Mikopo %d</div></div>", mwaka, bakaa, mapato, dharura, marejesho, faida_riba, matumizi, mikopo)
 fmt.Fprintf(w, "<div class=grid2><div class=card><div class=small>MAPATO</div><div class=big2>%d</div></div><div class=card><div class=small>DHARURA</div><div class=big2>%d</div></div><div class=card><div class=small>MATUMIZI</div><div class=big2 style=color:red>%d</div></div><div class=card><div class=small>MIKOPO+Riba</div><div class=big2 style=color:orange>%d + %d</div></div></div>", mapato, dharura, matumizi, mikopo, faida_riba)

 // WAKOPAJI NA RIBA + DUE DATE
 fmt.Fprint(w, "<div class=card><h3>👥 Waliochukua Mkopo + Riba + Due Date</h3>")
 if db!= nil {
  rows, _ := db.Query(`
   SELECT jina,
   COALESCE(SUM(kiasi),0) as kop,
   COALESCE(SUM(kiasi*riba/100),0) as riba_tot,
   COALESCE(MAX(tarehe_mwisho), NOW()) as due,
   COALESCE(SUM(CASE WHEN aina='mkopo' THEN riba ELSE 0 END),0) as riba_perc
   FROM miamala WHERE aina='mkopo' GROUP BY jina ORDER BY due ASC`)
  // marejesho per jina
  type rej struct{ kop, riba, marejesho int64; due time.Time; perc int }
  // query marejesho map
  rejMap := make(map[string]int64)
  rr, _ := db.Query("SELECT jina, COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='marejesho' GROUP BY jina")
  if rr!= nil { for rr.Next() { var j string; var s int64; rr.Scan(&j, &s); rejMap[j] = s } }
  if rows!= nil {
   for rows.Next() {
    var jina string
    var kop, ribaTot int64
    var due time.Time
    var perc int
    rows.Scan(&jina, &kop, &ribaTot, &due, &perc)
    marej := rejMap[jina]
    jumlaDeni := kop + ribaTot
    baki := jumlaDeni - marej
    per := 0
    if jumlaDeni > 0 { per = int(marej * 100 / jumlaDeni) }
    overdue := ""
    col := "#ff9800"
    if baki <= 0 { col = "#0f7a3a"; overdue = "✅ AMEMALIZA" } else if time.Now().After(due) { col = "#d32f2f"; overdue = fmt.Sprintf("🔴 IMECHELEWA! Due %s", due.Format("02/01/2006")) } else { overdue = fmt.Sprintf("Due: %s", due.Format("02/01/2006")) }
    fmt.Fprintf(w, "<div class=kopaji style=border-left-color:%s><div style=display:flex;justify-content:space-between><b>%s</b><span class=overdue>%s</span></div><div style=font-size:12px>Mkopo: <b>%d</b> + Riba %d%% = <b>%d</b> | Rejesha: <b>%d</b> | <b>Baki %d</b></div><div class=progress><div class=progress-bar style=width:%d%%;background:%s></div></div><div style=font-size:11px>%d%%</div></div>", col, jina, overdue, kop, perc, jumlaDeni, marej, baki, per, col, per)
   }
  }
 }
 fmt.Fprint(w, "</div>")

 fmt.Fprintf(w, "<div class=card><a href=/export><button class=\"btn btn-white\">📥 Pakua Excel na Riba</button></a><a href=\"%s\" target=_blank><button class=\"btn btn-wa\">WhatsApp</button></a></div>", waLink)
 fmt.Fprintf(w, "<div class=card><h3>+ Ongeza Muamala</h3><form method=POST action=/ongeza><input id=jina name=jina placeholder=\"Jina la Mkopaji\" required><input name=kiasi type=number placeholder=\"Kiasi TZS\" required><input type=date name=tarehe value=\"%s\" required><select id=aina name=aina required onchange=toggleMkopo()><option value=\"\" disabled selected>Chagua Aina</option><option value=mkopo>MKOPO - na Riba & Due</option><option value=marejesho>MAREJESHO YA MKOPO</option><option value=hisa_anzia>Hisa Anzia</option><option value=ada_mwezi>Ada ya Mwezi</option><option value=faini>Faini</option><option value=fedha_fomu>Fomu</option><option value=michango_mingine>Mingine</option><option value=fedha_dharura>Dharura</option><option value=matumizi>Matumizi</option></select><div id=ribaBox style=display:none><input name=riba type=number placeholder=\"Riba %% mf 10\" value=10><small style=font-size:11px>💰 Riba ya mkopo kwa %%</small></div><div id=dueBox style=display:none><input name=tarehe_mwisho type=date value=\"%s\"><small style=font-size:11px>📅 Tarehe ya mwisho kurudisha</small></div><button class=\"btn btn-green\">HIFADHI</button></form></div>", today, dueDefault)

 fmt.Fprintf(w, "<div class=card><h3>📜 Historia - EDIT</h3><table><tr><th>Tarehe</th><th>Jina</th><th>Aina</th><th>Kiasi</th><th>Riba</th><th>Due</th><th>EDIT</th></tr>", count)
 if db!= nil {
  rows, _ := db.Query("SELECT id,jina,kiasi,aina,tarehe,riba,tarehe_mwisho FROM miamala ORDER BY id DESC LIMIT 300")
  if rows!= nil {
   for rows.Next() {
    var id int
    var k int64
    var j, a string
    var t time.Time
    var riba sql.NullInt32
    var due sql.NullTime
    rows.Scan(&id, &j, &k, &a, &t, &riba, &due)
    rStr := "-"
    if riba.Valid { rStr = fmt.Sprintf("%d%%", riba.Int32) }
    dStr := "-"
    if due.Valid { dStr = due.Time.Format("02/01") }
    fmt.Fprintf(w, "<tr><td>%s</td><td>%s</td><td><span class=badge>%s</span></td><td><b>%d</b></td><td>%s</td><td>%s</td><td><a class=editbtn href=\"/edit?id=%d\">EDIT</a></td></tr>", t.Format("02/01"), j, a, k, rStr, dStr, id)
   }
  }
 }
 fmt.Fprint(w, "</table></div></body></html>")
}

func mwanachama(w http.ResponseWriter, r *http.Request) {
 r.ParseForm()
 if db!= nil { db.Exec("INSERT INTO wanachama(jina,simu,namba)VALUES($1,$2,$3)", r.FormValue("jina"), r.FormValue("simu"), r.FormValue("namba")) }
 http.Redirect(w, r, "/?add=ok", 302)
}
func ongeza(w http.ResponseWriter, r *http.Request) {
 r.ParseForm()
 jina := r.FormValue("jina")
 aina := r.FormValue("aina")
 kiasi, _ := strconv.ParseInt(r.FormValue("kiasi"), 10, 64)
 riba, _ := strconv.Atoi(r.FormValue("riba"))
 if r.FormValue("riba") == "" { riba = 0 }
 t, _ := time.Parse("2006-01-02", r.FormValue("tarehe"))
 if t.IsZero() { t = time.Now() }
 var due sql.NullTime
 if dt := r.FormValue("tarehe_mwisho"); dt!= "" {
  tt, _ := time.Parse("2006-01-02", dt)
  due = sql.NullTime{Time: tt, Valid: true}
 }
 if (aina == "fedha_dharura" || aina == "matumizi") && jina == "" { jina = "Washusema Nzima" }
 if db!= nil {
  if aina == "mkopo" {
   db.Exec("INSERT INTO miamala(jina,kiasi,aina,tarehe,riba,tarehe_mwisho)VALUES($1,$2,$3,$4,$5,$6)", jina, kiasi, aina, t, riba, due)
  } else {
   db.Exec("INSERT INTO miamala(jina,kiasi,aina,tarehe)VALUES($1,$2,$3,$4)", jina, kiasi, aina, t)
  }
 }
 http.Redirect(w, r, "/?add=ok", 302)
}
func del(w http.ResponseWriter, r *http.Request) {
 if db!= nil { db.Exec("DELETE FROM miamala WHERE id=$1", r.URL.Query().Get("id")) }
 http.Redirect(w, r, "/", 302)
}
func edit(w http.ResponseWriter, r *http.Request) {
 id := r.URL.Query().Get("id")
 if id == "" { id = r.FormValue("id") }
 if r.Method == "POST" {
  kiasi, _ := strconv.ParseInt(r.FormValue("kiasi"), 10, 64)
  riba, _ := strconv.Atoi(r.FormValue("riba"))
  t, _ := time.Parse("2006-01-02", r.FormValue("tarehe"))
  if t.IsZero() { t = time.Now() }
  var due sql.NullTime
  if dt := r.FormValue("tarehe_mwisho"); dt!= "" {
   tt, _ := time.Parse("2006-01-02", dt)
   due = sql.NullTime{Time: tt, Valid: true}
  }
  if db!= nil { db.Exec("UPDATE miamala SET jina=$1,kiasi=$2,aina=$3,tarehe=$4,riba=$5,tarehe_mwisho=$6 WHERE id=$7", r.FormValue("jina"), kiasi, r.FormValue("aina"), t, riba, due, r.FormValue("id")) }
  http.Redirect(w, r, "/?edit=ok", 302)
  return
 }
 var jina, aina string
 var kiasi int64
 var t time.Time
 var riba sql.NullInt32
 var due sql.NullTime
 if db!= nil { db.QueryRow("SELECT jina,kiasi,aina,tarehe,riba,tarehe_mwisho FROM miamala WHERE id=$1", id).Scan(&jina, &kiasi, &aina, &t, &riba, &due) }
 rVal := 0
 if riba.Valid { rVal = int(riba.Int32) }
 dVal := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
 if due.Valid { dVal = due.Time.Format("2006-01-02") }
 w.Header().Set("Content-Type", "text/html; charset=utf-8")
 fmt.Fprintf(w, "<html><head><meta name=viewport content=\"width=device-width,initial-scale=1\"><style>body{font-family:sans-serif;background:#f6f7f6;padding:12px}.card{background:#fff;border-radius:16px;padding:18px}input,select{width:100%%;padding:13px;border-radius:12px;border:1px solid #ddd;margin:8px 0}.btn{width:100%%;padding:13px;border-radius:12px;border:0;font-weight:700;background:#0f7a3a;color:#fff}</style></head><body><div class=card><h2>✏️ EDIT #%s</h2><form method=POST action=/edit><input type=hidden name=id value=\"%s\"><label>Jina</label><input name=jina value=\"%s\" required><label>Kiasi</label><input name=kiasi type=number value=\"%d\" required><label>Tarehe</label><input type=date name=tarehe value=\"%s\" required><label>Aina</label><select name=aina><option value=\"%s\" selected>%s</option><option value=mkopo>MKOPO</option><option value=marejesho>MAREJESHO</option><option value=hisa_anzia>Hisa</option><option value=ada_mwezi>Ada</option></select><label>Riba %%</label><input name=riba type=number value=\"%d\"><label>Due Date</label><input name=tarehe_mwisho type=date value=\"%s\"><button class=btn>💾 HIFADHI</button></form><br><a href=\"/\">Rudi</a> | <a href=\"/delete?id=%s\" style=color:red>Futa</a></div></body></html>", id, id, jina, kiasi, t.Format("2006-01-02"), aina, aina, rVal, dVal, id)
}
func exportCSV(w http.ResponseWriter, r *http.Request) {
 w.Header().Set("Content-Type", "text/csv")
 w.Header().Set("Content-Disposition", "attachment; filename=Mikopo_Riba_Due_"+time.Now().Format("2006")+".csv")
 writer := csv.NewWriter(w)
 writer.Write([]string{"Tarehe", "Jina", "Aina", "Kiasi", "Riba %", "Riba TZS", "Due Date"})
 if db!= nil {
  rows, _ := db.Query("SELECT tarehe,jina,aina,kiasi,riba,tarehe_mwisho FROM miamala ORDER BY id DESC")
  if rows!= nil {
   for rows.Next() {
    var t time.Time
    var j, a string
    var k int64
    var r sql.NullInt32
    var due sql.NullTime
    rows.Scan(&t, &j, &a, &k, &r, &due)
    rP := 0
    rT := int64(0)
    if r.Valid { rP = int(r.Int32); rT = k * int64(rP) / 100 }
    dS := ""
    if due.Valid { dS = due.Time.Format("02/01/2006") }
    writer.Write([]string{t.Format("02/01/2006"), j, a, fmt.Sprintf("%d", k), fmt.Sprintf("%d", rP), fmt.Sprintf("%d", rT), dS})
   }
  }
 }
 writer.Flush()
}
