package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func initDB() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		fmt.Println("DATABASE_URL missing!")
		return
	}
	var err error
	db, err = sql.Open("pgx", url)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		panic(err)
	}
	// Tengeneza table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS michango (
		id SERIAL PRIMARY KEY,
		jina TEXT,
		kiasi BIGINT,
		aina TEXT,
		tarehe TIMESTAMP DEFAULT NOW()
	)`)
	if err != nil {
		panic(err)
	}
	fmt.Println("DB OK - Postgres Mahida Ready")
}

func main() {
	initDB()
	http.HandleFunc("/", home)
	http.HandleFunc("/ongeza", ongeza)
	http.HandleFunc("/api/ripoti", ripoti)

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	fmt.Println("Mahida running on :" + port)
	http.ListenAndServe(":"+port, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
<!DOCTYPE html><html><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>Mahida System</title>
<style>body{font-family:sans-serif;padding:15px;background:#f5f5f5} .card{background:white;padding:15px;border-radius:12px;box-shadow:0 2px 8px #0001;margin-bottom:15px} input,select,button{width:100%;padding:12px;margin:5px 0;border-radius:8px;border:1px solid #ccc} button{background:#0d7a3a;color:white;font-weight:bold;border:none} .mapato{color:green;font-weight:bold} .matoleo{color:red}</style>
</head><body>
<div class="card"><h2>🕌 Mahida - Ingiza Muamala</h2>
<form action="/ongeza" method="POST">
<input name="jina" placeholder="Jina / Maelezo" required>
<input name="kiasi" type="number" placeholder="Kiasi TZS" required>
<select name="aina"><option value="mapato">Mapato</option><option value="matoleo">Matoleo</option></select>
<button>💾 Hifadhi</button>
</form></div>
<div class="card" id="ripoti">Inapakia ripoti...</div>
<script>
async function load(){let r=await fetch('/api/ripoti');let t=await r.text();document.getElementById('ripoti').innerHTML=t}
load()
</script>
</body></html>`)
}

func ongeza(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	jina := r.FormValue("jina")
	kiasi := r.FormValue("kiasi")
	aina := r.FormValue("aina")
	if jina == "" || kiasi == "" {
		http.Redirect(w, r, "/", 302)
		return
	}
	_, err := db.Exec("INSERT INTO michango (jina,kiasi,aina) VALUES ($1,$2,$3)", jina, kiasi, aina)
	if err != nil {
		fmt.Fprint(w, "Error: "+err.Error())
		return
	}
	http.Redirect(w, r, "/", 302)
}

func ripoti(w http.ResponseWriter, r *http.Request) {
	var mapato, matoleo sql.NullInt64
	db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM michango WHERE aina='mapato' AND date_trunc('month',tarehe)=date_trunc('month',NOW())").Scan(&mapato)
	db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM michango WHERE aina='matoleo' AND date_trunc('month',tarehe)=date_trunc('month',NOW())").Scan(&matoleo)
	baki := mapato.Int64 - matoleo.Int64

	rows, _ := db.Query("SELECT jina,kiasi,aina,to_char(tarehe,'DD Mon HH24:MI') FROM michango ORDER BY tarehe DESC LIMIT 20")
	defer rows.Close()
	var list strings.Builder
	for rows.Next() {
		var j, a, t string
		var k int64
		rows.Scan(&j, &k, &a, &t)
		cls := "mapato"
		if a == "matoleo" {
			cls = "matoleo"
		}
		list.WriteString(fmt.Sprintf("<div><span class='%s'>%s: %d (%s)</span> <small>%s</small></div>", cls, j, k, a, t))
	}
	fmt.Fprintf(w, `<h3>Ripoti ya Mwezi Huu</h3>
	<p class="mapato">Mapato: %d TZS</p>
	<p class="matoleo">Matoleo: %d TZS</p>
	<h3>Baki: %d TZS</h3><hr>%s`, mapato.Int64, matoleo.Int64, baki, list.String())
}