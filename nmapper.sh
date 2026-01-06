#!/bin/bash

echo "Analizando la lista de IPs proporcionada"
for i in $@
do
        echo "Analizando la IP: $i\n\n\n"
        nmap -p 0-65535 --open -A -T0 -sS "$i" -oX nmap_"$i".xml && xsltproc nmap_"$i".xml -o nmap_"$i".html 2 >
 error.log
        echo "Análisis de $i completo\n\n\n"

done

echo "Ya se han analizado todas la IPs de la lista"
