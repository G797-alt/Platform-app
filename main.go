package main

import (
	"log"
	"net/http"
	"os"
)

const htmlPage = `<!DOCTYPE html><html lang="sw"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Uwashusema Mahida - na Fomu 5% & Riba 30%</title><style>*{box-sizing:border-box}body{margin:0;font-family:system-ui;background:#f1f5f1;padding-bottom:20px}header{background:#1a6d2f;color:#fff;padding:16px;text-align:center;font-weight:900;font-size:18px}.card{background:#fff;margin:10px;border-radius:12px;padding:12px;box-shadow:0 2px 8px #0001}h3{margin:0 0 10px 0;font-size:14px;color:#1a6d2f;border-bottom:1px solid #eee;padding-bottom:6px}.top3{display:grid;grid-template-columns:1fr 1fr 1fr;gap:6px}.top3 div{padding:10px;border-radius:8px;text-align:center;color:#fff;font-size:13px}.mapato{background:#2e7d32}.matumizi{background:#c62828}.salio{background:#1565c0}select,input{width:100%;padding:11px;border-radius:8px;border:1px solid #ccc;margin-bottom:8px}.btn{width:100%;padding:12px;border:0;border-radius:8px;font-weight:900;color:#fff;margin-top:6px}.green{background:#1a6d2f}.blue{background:#1565c0}.red{background:#c62828}.orange{background:#ef6c00}table{width:100%;border-collapse:collapse;font-size:12px}th{background:#e8f5e9;padding:6px;text-align:left}td{padding:7px;border-bottom:1px solid #eee}[contenteditable]{background:#fffde7;outline:1px solid #ffca28;border-radius:4px}.row{display:flex;justify-content:space-between;padding:8px 0;border-bottom:1px solid #eee;font-size:13px}.badge{background:#ffecb3;color:#7a4a00;padding:2px 6px;border-radius:10px;font-size:10px;font-weight:800}</style></head><body>
<header>Uwashusema Mahida<br><small style="font-weight:400;font-size:12px">Fomu 5% & Riba 30% Automatic</small></header>
<div class="card"><div style="display:flex;gap:6px"><select id="mwezi"><option value="1">Jan</option><option value="2">Feb</option><option value="3">Mar</option><option value="4">Apr</option><option value="5">Mei</option><option value="6">Jun</option><option value="7">Jul</option><option value="8">Ago</option><option value="9">Sep</option><option value="10">Okt</option><option value="11">Nov</option><option value="12">Des</option></select><select id="mwaka"><option>2025</option><option selected>2026</option></select><button onclick="render()" style="background:#1a6d2f;color:#fff;border:0;border-radius:8px;padding:0 15px">Sawa</button></div><div class="top3" style="margin-top:10px"><div class="mapato"><small>MAPATO</small><b id="tMapato">0</b><br><small id="tFomu" style="font-size:10px"></small></div><div class="matumizi"><small>MATUMIZI</small><b id="tMatumizi">0</b></div><div class="salio"><small>SALIO</small><b id="tSalio">0</b></div></div></div>
<div class="card"><h3>Wanachama <button onclick="addW()" style="float:right;background:#1a6d2f;color:#fff;border:0;padding:4px 10px;border-radius:6px">+ Mpya</button></h3><div id="wlist"></div></div>
<div class="card"><h3>Ingiza Data</h3><select id="selW"></select><select id="selAina"><optgroup label="MAPATO"><option value="hisa">Hisa</option><option value="ada">Ada ya Mwezi</option><option value="marejesho">Marejesho ya Mkopo</option><option value="faini">Faini / Riba 30%</option><option value="fomu">Fomu (5% auto)</option><option value="michango">Michango Mingine</option></optgroup><optgroup label="MATUMIZI"><option value="mkopo_toka">Mkopo kwa Mwanachama</option><option value="dharura">Dharura</option><option value="mkono_pole">Mkono wa Pole</option><option value="matumizi_mengine">Matumizi Mengine</option></optgroup></select><input id="kiasi" type="number" placeholder="Kiasi TZS"><input id="rejeshoTarajiwa" type="number" placeholder="Rejesho la Mwezi Analotakiwa (kwa Riba)"><input id="maelezo" placeholder="Maelezo"><button class="btn green" onclick="addM()">Hifadhi</button><small style="color:#666">Ukiselecta Mkopo kwa Mwanachama, Fomu 5% itajiongeza yenyewe kwenye Mapato</small></div>
<div class="card"><h3>Wenye Deni & Wasio Lipa (Piga Riba 30%)</h3><table><thead><tr><th>Jina</th><th>Deni</th><th>Rejesho Tarajiwa</th><th>Riba 30%</th></tr></thead><tbody id="tblDeni"></tbody></table></div>
<div class="card"><h3>Muhtasari + Tuma</h3><div id="muhtasari"></div><button class="btn blue" onclick="backup()">📊 Pakua Excel (Fomu 5% & Riba Ndani)</button><button class="btn green" onclick="wa()">📱 Tuma WhatsApp Ripoti</button></div>
<div class="card"><h3>Historia</h3><div id="hist"></div></div>
<script>
let DATA=JSON.parse(localStorage.getItem('mahida_v3')||'null')||{wanachama:[{jina:"Charles Joseph Mgaya"}],miamala:[]};
const curM=()=>+mwezi.value,curY=()=>+mwaka.value;
const save=()=>{localStorage.setItem('mahida_v3',JSON.stringify(DATA));render()};
function addW(){let j=prompt('Jina:');if(!j)return;DATA.wanachama.push({jina:j});save()}
function addM(){
 let w=selW.value,a=selAina.value,k=+kiasi.value,rt=+rejeshoTarajiwa.value||0,ma=maelezo.value;
 if(!w||!a||!k)return alert('Jaza jina, aina na kiasi');
 // 1. Hifadhi muamala mkuu
 DATA.miamala.push({w,a,k,rt,ma,mwezi:curM(),mwaka:curY(),t:new Date().toLocaleString()});
 // 2. LOGIC YA FOMU 5% - kama ni mkopo, ongeza fomu auto
 if(a==='mkopo_toka'){
   let fomu = Math.round(k*0.05);
   DATA.miamala.push({w,a:'fomu',k:fomu,rt:0,ma:'Fomu 5% ya mkopo TZS '+k.toLocaleString(),mwezi:curM(),mwaka:curY(),t:new Date().toLocaleString()+' - AUTO'});
   alert('Mkopo TZS '+k.toLocaleString()+' umehifadhiwa + Fomu 5% = TZS '+fomu.toLocaleString()+' imeongezwa kwenye Mapato');
 }
 kiasi.value='';rejeshoTarajiwa.value='';maelezo.value='';save();
}
function pigaRiba(jina, rejesho){
 let riba = Math.round(rejesho*0.30);
 if(!confirm('Kumpiga '+jina+' Riba 30% ya TZS '+rejesho.toLocaleString()+' = TZS '+riba.toLocaleString()+' ?'))return;
 DATA.miamala.push({w:jina,a:'faini',k:riba,rt:0,ma:'Riba 30% ya kutolipa rejesho TZS '+rejesho.toLocaleString(),mwezi:curM(),mwaka:curY(),t:new Date().toLocaleString()+' - RIBA AUTO'});
 save();
}
function render(){
 let m=curM(),y=curY(),mi=DATA.miamala.filter(x=>x.mwezi==m&&x.mwaka==y);
 let mapA=['hisa','ada','marejesho','faini','fomu','michango'],matA=['mkopo_toka','dharura','mkono_pole','matumizi_mengine'];
 let map=0,mat=0,fomuTotal=0; mi.forEach(t=>{if(mapA.includes(t.a))map+=t.k;if(t.a=='fomu')fomuTotal+=t.k;if(matA.includes(t.a))mat+=t.k;});
 tMapato.textContent='TZS '+map.toLocaleString();tFomu.textContent='Fomu 5%: '+fomuTotal.toLocaleString();tMatumizi.textContent='TZS '+mat.toLocaleString();tSalio.textContent='TZS '+(map-mat).toLocaleString();
 wlist.innerHTML='';selW.innerHTML='<option value="">Chagua Mwanachama</option>';
 DATA.wanachama.forEach((w,i)=>{wlist.innerHTML+='<div class="row"><b>'+w.jina+'</b><span><button onclick="let n=prompt(\'Jina jipya:\',DATA.wanachama['+i+'].jina);if(n){let old=DATA.wanachama['+i+'].jina;DATA.wanachama['+i+'].jina=n;DATA.miamala.forEach(x=>{if(x.w==old)x.w=n});save()}" style="border:0;background:#e3f2fd;border-radius:4px">✏️</button> <button onclick="if(confirm(\'Futa?\')){DATA.wanachama.splice('+i+',1);save()}" style="border:0;background:#ffebee;color:#c00;border-radius:4px">X</button></span></div>';selW.innerHTML+='<option>'+w.jina+'</option>'});
 // DENI
 let dmap={}; DATA.miamala.forEach(t=>{if(!dmap[t.w])dmap[t.w]={c:0,r:0,lastRT:0}; if(t.a=='mkopo_toka')dmap[t.w].c+=t.k; if(t.a=='marejesho')dmap[t.w].r+=t.k; if(t.rt>0)dmap[t.w].lastRT=t.rt;});
 let td=tblDeni;td.innerHTML='';
 for(let j in dmap){let b=dmap[j].c-dmap[j].r; if(b>0||dmap[j].c>0){let rt=dmap[j].lastRT||0; let riba=Math.round(rt*0.30); td.innerHTML+='<tr><td>'+j+'</td><td style="color:#c00;font-weight:700">'+b.toLocaleString()+'</td><td contenteditable="true" onblur="let v=this.innerText.replace(/,/g,\'\');for(let i=DATA.miamala.length-1;i>=0;i--){if(DATA.miamala[i].w===\''+j+'\'&&DATA.miamala[i].rt>0){DATA.miamala[i].rt=+v;break;}}if(v>0){let found=false;for(let x of DATA.miamala){if(x.w===\''+j+'\'&&x.mwezi==='+m+'&&x.mwaka==='+y+'&&x.rt==v){found=true}}if(!found)DATA.miamala.push({w:\''+j+'\',a:\'marejesho\',k:0,rt:+v,ma:\'Rejesho tarajiwa\',mwezi:'+m+',mwaka:'+y+',t:new Date().toLocaleString()});}save()">'+(rt||'<small style=color:#999>weka</small>')+'</td><td>'+(rt?'<button onclick="pigaRiba(\''+j+'\','+rt+')" class="red" style="border:0;padding:5px 8px;border-radius:5px;color:#fff;font-size:11px">Piga 30% = '+riba.toLocaleString()+'</button>':'')+'</td></tr>';}}
 if(!td.innerHTML)td.innerHTML='<tr><td colspan=4 style="color:#999">Hakuna deni mwezi huu</td></tr>';
 let mu=muhtasari;mu.innerHTML='';let labels={hisa:'Hisa',ada:'Ada',marejesho:'Marejesho',faini:'Faini / Riba 30%',fomu:'Fomu (5%)',michango:'Michango',mkopo_toka:'Mkopo Toka',dharura:'Dharura',mkono_pole:'Mkono wa Pole',matumizi_mengine:'Matumizi Mengine'};let sum={};document.querySelectorAll('#selAina option').forEach(o=>sum[o.value]=0);mi.forEach(t=>sum[t.a]+=t.k);for(let k in labels){if(sum[k]>0)mu.innerHTML+='<div class="row"><span>'+labels[k]+'</span><b>'+sum[k].toLocaleString()+'</b></div>'}
 let hi=hist;hi.innerHTML='';[...mi].reverse().forEach(t=>{let idx=DATA.miamala.indexOf(t);hi.innerHTML+='<div class="row"><div><b>'+t.w+'</b> - '+t.a+' '+(t.a=='fomu'?'<span class=badge>5% AUTO</span>':'')+(t.a=='faini'?'<span class=badge style=background:#ffcdd2;color:#b71c1c>30% RIBA</span>':'')+'<br><small>'+t.t+' - TZS '+t.k.toLocaleString()+' '+(t.ma?' - '+t.ma:'')+'</small></div><div><button onclick="let v=prompt(\'Kiasi jipya:\',DATA.miamala['+idx+'].k);if(v!==null){DATA.miamala['+idx+'].k=+v;save()}" style="border:0;background:#e3f2fd;border-radius:4px">✏️</button> <button onclick="if(confirm(\'Futa?\')){DATA.miamala.splice('+idx+',1);save()}" style="border:0;background:#ffebee;border-radius:4px;color:#c00">X</button></div></div>'});
}
function backup(){let m=curM(),y=curY();let dmap={};DATA.miamala.filter(x=>x.mwezi==m&&x.mwaka==y).forEach(t=>{if(!dmap[t.w])dmap[t.w]={c:0,r:0,fomu:0,faini:0};if(t.a=='mkopo_toka')dmap[t.w].c+=t.k;if(t.a=='marejesho')dmap[t.w].r+=t.k;if(t.a=='fomu')dmap[t.w].fomu+=t.k;if(t.a=='faini')dmap[t.w].faini+=t.k;});let csv='Jina,Mkopo Aliouchukua,Jumla Marejesho,Fomu 5%,Riba 30%,Baki\n';for(let j in dmap){let b=dmap[j].c-dmap[j].r;csv+='"'+j+'",'+dmap[j].c+','+dmap[j].r+','+dmap[j].fomu+','+dmap[j].faini+','+b+'\n';}let a=document.createElement('a');a.href=URL.createObjectURL(new Blob([csv],{type:'text/csv'}));a.download='Uwashusema_Fomu5_Riba30_'+m+'_'+y+'.csv';a.click();}
function wa(){let txt='*Uwashusema Mahida* '+curM()+'/'+curY()+'\nMAPATO: '+tMapato.textContent+' ('+tFomu.textContent+')\nMATUMIZI: '+tMatumizi.textContent+'\nSALIO: '+tSalio.textContent+'\n\nWenye Deni:\n';for(let r of tblDeni.rows){txt+=r.innerText+'\n'}window.open('https://wa.me/?text='+encodeURIComponent(txt),'_blank')}
mwezi.value=new Date().getMonth()+1;render();
</script></body></html>`

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlPage))
	})
	log.Println("Running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
