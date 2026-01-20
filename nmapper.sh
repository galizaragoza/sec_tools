#!/bin/bash
## Este script toma una lista de IPs y ejecuta el comando de nmap que está en el bucle for e
## itera por cada IP de la lista ejecutando ese comando.
## La idea es modificar los parámetros del mismo en función de las necesidades del escaneo.
## Una vez ha terminado de escanear, borra todos los archivos .xml generados por nmap y
## comprime todos los archivos .html en un archivo.
## El script debe ejecutarse en el directorio en el que se quieren almacenar estos archivos
## de la siguiente manera: sh ruta_al_script $(cat ruta_a_la_lista_de_IPs)

echo "Que parámetros de nmap quieres?"
read params

echo "Creando directorio para meter todos los archivos\n\n\n"
mkdir nmapper_scan
sleep 1
cd nmapper_scan && echo "Carpeta creada"

echo "Analizando la lista de IPs proporcionada"
counter=1
for i in $@; do
  total=$#
  perc=$((counter * 100 / total))
  echo "Analizando la IP: $i\n\n $counter de $total ($perc %)"
  nmap $params "$i" -oX nmap_"$i".xml 2>>errors.log 1>/dev/null &&
    xsltproc nmap_"$i".xml -o nmap_"$i".html

  echo "Análisis de $i completo\n\n\n"
  counter=$((counter + 1))
done

echo "Ya se han analizado todas la IPs de la lista\n\n\n"
if [ -s errors.log ]; then
  echo "El script se ha completado sin generar errores"
else
  echo "Ha habido al menos 1 error durante el escaneo, revisa errors.log"
fi

echo "Eliminando todos los archivos .xml...\n\n\n"
rm *.xml

if [$? >0]; then
  echo "Comprimiendo todos los archivos html"
  tar -c -z -f parsed.tar.gz *.html
fi

echo "Compresión completa\n\n\n"

echo "Trabajo completado"

echo "Saliendo"
