# Tarea 2 Sistemas Distribuidos
|Nombre| Rol|
|-----|-----|
|Eduardo Borgoño| 201373525-8 |    
|Francisca Ramírez| 201373607-6 |

[Avión](https://github.com/eduborgono/distribuidos-avion) 

[Pantalla de información](https://github.com/eduborgono/pantalla-informacion) 


# Instalación GO
Es necesario descargar go en su versión 1.11.2 https://golang.org/dl/ y tener instalado git.

Luego se deben seguir las instrucciones de instalación adecuadas al sistema operativo https://golang.org/doc/install

Es absolutamente necesario seguir las instrucciones al pie de la letra, es decir, **instalar en los directorios
indicados y añadir las variables de entorno correspondientes**.

En windows puede ser necesario crear la carpeta 
```
%USERPROFILE%\go 
```
y añadir una variable de entorno llamada ```GOPATH``` con valor 
```
%USERPROFILE%\go
```
si es que no fuera creada con el instalador. 

# Instalación Torre de Control
Se debe correr el comando
```
go get github.com/eduborgono/torreControl
```
# Ejecución Torre de Control
### Windows
```
cd %USERPROFILE%\go\bin
.\torreControl.exe
```
### Linux
```
cd $HOME/go/bin
./torreControl
```

