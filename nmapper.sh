#!/bin/bash

echo "Analizando la lista de IPs proporcionada"
for i in $@
do
        total=$#
        counter=1
        perc=$((counter*100/total))
        echo "Analizando la IP: $i\n\n\n $counter de $total ($perc %)"
        nmap --top-ports 10000 --open -A -T3 -v -sT -Pn "$i" -oX nmap_"$i".xml && xsltproc nmap_"$i".xml -o nmap
_"$i".html 2>>errors.log
        echo "Análisis de $i completo\n\n\n"
        counter=$((counter+1))

done

echo "Ya se han analizado todas la IPs de la lista"
