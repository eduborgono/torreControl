# Instalaci�n GO
Es necesario descargar go en su versi�n 1.11.2 https://golang.org/dl/ y tener instalado git
Luego se deben seguir las instrucciones de instalaci�n adecuadas al sistema operativo https://golang.org/doc/install
Es absolutamente necesario seguir las instrucciones al pie de la letra, esto quiere decir instalar en los directorios
especificados y a�adir las variables de entorno correspondientes.

En windows es necesario crear la carpeta %USERPROFILE%\go y 
a�adir una variable de entorno llamada GOPATH con valor 
%USERPROFILE%\go
si es que no fuera creada con el instalador. 

# Instalaci�n Torre Control
Se debe correr el comando
```
go get github.com/eduborgono/torreControl
```
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

