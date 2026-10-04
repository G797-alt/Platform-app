var mapato, matumizi int64
db.QueryRow("SELECT COALESCE(SUM(kiasi),0) FROM mahida").Scan(&mapato)
