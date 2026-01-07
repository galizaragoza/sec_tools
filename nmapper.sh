#!/bin/bash
##Este script toma una lista de IPs y ejecuta el comando de nmap que está en el bucle for e itera por cada IP de la lista ejecutando ese comando.
## La idea es modificar los parámetros del mismo en función de las necesidades del escaneo. Una vez ha terminado de escanear, borra todos los archivos .xml generados por nmap y comprime todos los archivos .html en un archivo.
## El script debe ejecutarse en el directorio en el que se quieren almacenar estos archivos de la siguiente manera: sh ruta_al_script $(cat ruta_a_la_lista_de_IPs)

echo "Este script ha de correrse en la carpeta en la que quieran generarse todos los archivos con lo
s scans, estás seguro de que este es el directorio correcto?\n\n\n"
pwd && echo "\n\n\n" && sleep 1
cooldown=5
while [ $cooldown -gt 0 ]; do
        echo "Comenzando en $cooldown"
        sleep 1
        cooldown=$((cooldown-1))
done


echo "Analizando la lista de IPs proporcionada"
counter=1
for i in $@
do
        total=$#
        perc=$((counter*100/total))
        echo "Analizando la IP: $i\n\n\n $counter de $total ($perc %)"
        nmap --top-ports 10000 --open -A -T4 -v -sT -Pn "$i" -oX nmap_"$i".xml && xsltproc nmap_"$i"
.xml -o nmap_"$i".html 2>>errors.log
        echo "Análisis de $i completo\n\n\n"
        counter=$((counter+1))

done

echo "Ya se han analizado todas la IPs de la lista\n\n\n"

echo "Eliminando todos los archivos .xml...\n\n\n"
rm *.xml

echo "Comprimiendo todos los archivos html"
tar -c -z -f parsed.tar.gz *.html


echo "Compresión completa\n\n\n"

echo "Trabajo completado"

echo "Saliendo"
