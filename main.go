package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func main() {
	if u := os.Getenv("DATABASE_URL"); u != "" {
		var err error
		db, err = sql.Open("pgx", u)
		if err == nil {
			db.Exec(`CREATE TABLE IF NOT EXISTS wanachama(id SERIAL PRIMARY KEY,jina TEXT,simu TEXT,namba TEXT)`)
			db.Exec(`CREATE TABLE IF NOT EXISTS miamala(id SERIAL PRIMARY KEY,jina TEXT,kiasi BIGINT,aina TEXT,tarehe TIMESTAMP DEFAULT NOW())`)
		}
	}
	http.HandleFunc("/", home)
	http.HandleFunc("/ongeza", ongeza)
	http.HandleFunc("/mwanachama", mwanachama)
	http.HandleFunc("/export", exportCSV)
	http.HandleFunc("/delete", del)
	http.HandleFunc("/edit", edit)
	http.HandleFunc("/futa-zote", futaZote)

	p := os.Getenv("PORT")
	if p == "" {
		p = "10000"
	}
	log.Println("Starting on port", p)
	http.ListenAndServe("0.0.0.0:"+p, nil)
}

func getTotals() (mapato, dharura, matumizi, mikopo int64) {
	if db == nil {
		return
	}
	db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina IN ('hisa_anzia','ada_mwezi','marejesho','faini','fedha_fomu','michango_mingine')`).Scan(&mapato)
	db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='fedha_dharura'`).Scan(&dharura)
	db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='matumizi'`).Scan(&matumizi)
	db.QueryRow(`SELECT COALESCE(SUM(kiasi),0) FROM miamala WHERE aina='mkopo'`).Scan(&mikopo)
	return
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	mapato, dharura, matumizi, mikopo := getTotals()
	bakaa := (mapato + dharura) - (matumizi + mikopo)
	waText := fmt.Sprintf("UWASHUSEMA MAHIDA RIPOTI BAKAA TZS %d MAPATO %d DHARURA %d MATUMIZI %d MIKOPO %d", bakaa, mapato, dharura, matumizi, mikopo)
	waLink := "https://wa.me/?text=" + url.QueryEscape(waText)

	alert := ""
	if r.URL.Query().Get("futa") == "ok" {
		alert = `<div class="card" style="background:#ffebee;border-left:4px solid red">✅ Data zote zimefutwa!</div>`
	}
	if r.URL.Query().Get("edit") == "ok" {
		alert = `<div class="card" style="background:#e8f5e9;border-left:4px solid green">✅ Muamala umehaririwa!</div>`
	}

	w.Write([]byte(`
<!DOCTYPE html><html><head><meta name="viewport" content="width=device-width,initial-scale=1"><title>Uwashusema Mahida</title>
<style>
body{font-family:sans-serif;background:#f6f7f6;margin:0;padding:0 10px 90px}
.header{padding:18px 6px 8px;font-weight:800;font-size:20px}
.card{background:#fff;border-radius:16px;padding:14px;margin-bottom:12px}
.bakaa{background:#0f7a3a;color:#fff;border-radius:20px;padding:18px}
.big{font-size:28px;font-weight:900}
.grid2{display:grid;grid-template-columns:1fr 1fr;gap:10px}
.small{font-size:11px;opacity:.8}
.big2{font-size:18px;font-weight:800}
.btn{width:100%;padding:12px;border-radius:12px;border:0;font-weight:700;margin:6px 0;cursor:pointer}
.btn-green{background:#0f7a3a;color:#fff}
.btn-white{background:#fff;border:1px solid #ddd}
.btn-wa{background:#25D366;color:#fff}
.btn-red{background:#fff;border:1px solid #ef5350;color:#ef5350}
input,select{width:100%;padding:13px;border-radius:12px;border
