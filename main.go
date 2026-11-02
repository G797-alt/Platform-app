
package main

import (
	"log"
	"net/http"
	"os"
)

const htmlPage = `<!DOCTYPE html><html lang="sw"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Uwashusema FIXED</title><script src="https://unpkg.com/xlsx@0.18.5/dist/xlsx.full.min.js"></script><style>*{box-sizing:border-box}body{margin:0;font-family:system-ui;background:#f1f5f1;padding-bottom:30px}header{background:#1a6d2f;color:#fff;padding:16px;text-align:center;font-weight:900;font-size:18px}.card{background:#fff;margin:10px;border-radius:12px;padding:12px;box-shadow:0 2px 8px #0001}h3{margin:0 0 10px 0;font-size:14px;color:#1a6d2f;border-bottom:1px solid #eee;padding-bottom:6px}.top3{display:grid;grid-template-columns:1fr 1fr 1fr;gap:6px}.top3 div{padding:10px;border-radius:8px;text-align:center;color:#fff;font-weight:900}.g{background:#2e7d32}.r{background:#c62828}.b{background:#1565c0}select,input{width:100%;padding:12px;border-radius:8px;border:1px solid #bbb;margin-bottom:8px;font-size:14px}.btn{width:100%;padding:12px;border:0;border-radius:8px;font-weight:900;color:#fff;margin-top:6px;font-size:14px}.green{background:#1a6d2f}.blue{background:#1565c0}.row{display:flex;justify-content:space-between;align-items:center;padding:8px 0;border-bottom:1px solid #eee;font-size:13px}</style></head><body>
<header>Uwashusema Mahida</header>
<div class="card"><div style="display:flex;gap:6px"><select id="mwezi"><option value="1">Jan</option><option value="2">Feb</option><option value="3">Mar</option><option value="4">Apr</option><option value="5">Mei</option><option value="6">Jun</option><option value="7">Jul</option><option value="8">Ago</option><option value="9">Sep</option><option value="10">Okt</option><option value="11">Nov</option><option value="12">Des</option></select><input id="mwaka" type="number" value="2026"><button onclick="render()" style="background:#1a6d2f;color:#fff;border:0;border-radius:8px;padding:0 15px">Sawa</button></div><div class="top3" style="margin-top:10px"><div class="g"><small>MAPATO</small><br><b id="tMapato">0</b></div><div class="r"><small>MATUMIZI</small><br><b id="tMatumizi">0</b></div><div class="b"><small>SALIO</small><br><b id="tSalio">0</b></div></div></div>
<div class="card"><h3>Wanachama <button onclick="addW()" style="float:right;background:#1a6d2f;color:#fff;border:0;padding:5px 12px;border-radius:6px">+ Mpya</button></h3><div id="wlist"></div></div>
<div class="card"><h3>Ingiza Data</h3><select id="selW"></select><select id="selAina"><option value="hisa">Hisa</option><option value="ada">Ada ya Mwezi</option><option value="marejesho">Marejesho ya Mkopo</option><option value="faini">Faini / Riba 30%</option><option value="fomu">Fomu 5%</option><option value="michango">Michango Mingine</option><option value="mkopo_toka">Mkopo kwa Mwanachama</option><option value="dharura">Dharura</option><option value="mkono_pole">Mkono wa Pole</option><option value="matumizi_mengine">Matumizi Mengine</option></select><input id="kiasi" type="number" placeholder="Kiasi TZS"><input id="maelezo" placeholder="Maelezo"><button class="btn green" onclick="addM()">Hifadhi</button><div id="calc" style="color:#1a6d2f;font-weight:800;margin-top:6px;font-size:13px"></div></div>
<div class="card"><h3>Muhtasari</h3><div id="muhtasari"></div><button class="btn blue" onclick="backup()">📊 Pakua Excel 4 Sheets</button></div>
<div class="card"><h3>Historia - Gusa ✏️ Kuhariri</h3><div id="hist"></div></div>
<script>
var DATA = JSON.parse(localStorage.getItem("mahida_FINAL2") || "null") || {wanachama:[{jina:"Charles Joseph Mgaya"}], miamala:[]};
function curM(){return parseInt(document.getElementById("mwezi").value)}
function curY(){return parseInt(document.getElementById("mwaka").value)||2026}
function save(){localStorage.setItem("mahida_FINAL2", JSON.stringify(DATA)); render();}
function addW(){
 var j = prompt("Jina la Mwanachama:");
 if(!j) return;
 DATA.wanachama.push({jina:j});
 save();
 alert("Ameongezwa: "+j);
}
function editW(i){
 var oldName = DATA.wanachama[i].jina;
 var n = prompt("Hariri Jina:", oldName);
 if(!n || n===oldName) return;
 DATA.wanachama[i].jina = n;
 for(var k=0;k<DATA.miamala.length;k++){ if(DATA.miamala[k].w===oldName){DATA.miamala[k].w=n;}}
 save();
}
function delW(i){
 if(!confirm("Futa "+DATA.wanachama[i].jina+" ?")) return;
 DATA.wanachama.splice(i,1);
 save();
}
var selAina = document.getElementById("selAina");
var kiasi = document.getElementById("kiasi");
var calc = document.getElementById("calc");
selAina.addEventListener("change", showCalc);
kiasi.addEventListener("input", showCalc);
function showCalc(){
 var k = parseFloat(kiasi.value)||0;
 if(selAina.value==="mkopo_toka" && k>0){
   var f = Math.round(k*0.05);
   calc.innerHTML = "Fomu 5% = TZS "+f.toLocaleString()+" | Kwa 1,500,000 = 75,000";
 } else {
   calc.innerHTML = "";
 }
}
function addM(){
 var w = document.getElementById("selW").value;
 var a = selAina.value;
 var k = parseFloat(kiasi.value)||0;
 var ma = document.getElementById("maelezo").value;
 if(!w || !a || !k){alert("Jaza mwanachama, aina na kiasi"); return;}
 DATA.miamala.push({w:w, a:a, k:k, ma:ma, mwezi:curM(), mwaka:curY(), t:new Date().toLocaleString()});
 if(a==="mkopo_toka"){
   var fomu = Math.round(k*0.05);
   DATA.miamala.push({w:w, a:"fomu", k:fomu, ma:"Fomu 5% ya mkopo "+k, mwezi:curM(), mwaka:curY(), t:new Date().toLocaleString()+" AUTO 5%"});
   alert("Mkopo: "+k.toLocaleString()+"\nFomu 5%: "+fomu.toLocaleString()+"\n\n1,500,000 = 75,000 (sio 150,000)");
 }
 document.getElementById("kiasi").value="";
 document.getElementById("maelezo").value="";
 calc.innerHTML="";
 save();
}
function editK(idx){
 var v = prompt("Kiasi Kipya:", DATA.miamala[idx].k);
 if(v===null) return;
 var nk = parseFloat(v);
 if(isNaN(nk)) return;
 DATA.miamala[idx].k = nk;
 if(DATA.miamala[idx].a==="mkopo_toka"){
   for(var i=0;i<DATA.miamala.length;i++){
     if(DATA.miamala[i].a==="fomu" && DATA.miamala[i].w===DATA.miamala[idx].w && DATA.miamala[i].mwezi===DATA.miamala[idx].mwezi){
       if(DATA.miamala[i].t.indexOf("AUTO")>-1){DATA.miamala[i].k = Math.round(nk*0.05);}
     }
   }
 }
 save();
}
function delK(idx){
 if(!confirm("Futa muamala?")) return;
 DATA.miamala.splice(idx,1);
 save();
}
function render(){
 var m = curM(); var y = curY();
 var mi = [];
 for(var i=0;i<DATA.miamala.length;i++){if(DATA.miamala[i].mwezi==m && DATA.miamala[i].mwaka==y){mi.push(DATA.miamala[i]);}}
 var mapA=["hisa","ada","marejesho","faini","fomu","michango"];
 var matA=["mkopo_toka","dharura","mkono_pole","matumizi_mengine"];
 var map=0,mat=0;
 for(var i=0;i<mi.length;i++){if(mapA.indexOf(mi[i].a)>=0) map+=mi[i].k; if(matA.indexOf(mi[i].a)>=0) mat+=mi[i].k;}
 document.getElementById("tMapato").innerText = map
