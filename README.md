# Tarea 2 Sistemas Distribuidos
|Nombre| Rol|
|-----|-----|
|Eduardo Borgo�o| 201373525-8 |    
|Francisca Ram�rez| 201373607-6 |

[Avi�n](https://github.com/eduborgono/distribuidos-avion) 

[Pantalla de informaci�n](https://github.com/eduborgono/pantalla-informacion) 


# Instalaci�n GO
Es necesario descargar go en su versi�n 1.11.2 https://golang.org/dl/ y tener instalado git.

Luego se deben seguir las instrucciones de instalaci�n adecuadas al sistema operativo https://golang.org/doc/install

Es absolutamente necesario seguir las instrucciones al pie de la letra, es decir, **instalar en los directorios
indicados y a�adir las variables de entorno correspondientes**.

En windows puede ser necesario crear la carpeta 
```
%USERPROFILE%\go 
```
y a�adir una variable de entorno llamada ```GOPATH``` con valor 
```
%USERPROFILE%\go
```
si es que no fuera creada con el instalador. 

# Instalaci�n Torre de Control
Se debe correr el comando
```
go get github.com/eduborgono/torreControl
```
# Ejecuci�n Torre de Control
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

