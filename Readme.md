hub es una herramienta de línea de comandos que se ajusta gitpara ampliarla con funciones y comandos adicionales que facilitan el trabajo con GitHub.

Para obtener una interfaz de línea de comandos oficial y potencialmente más fácil de usar para GitHub, consulte cli.github.com y esta comparación .

Este repositorio y su rastreador de problemas no son para informar problemas con la interfaz web de GitHub.com . Si tiene un problema con el propio GitHub, comuníquese con Soporte .

Uso
$ hub clon rtomayko / tilt
 # => git clone git: //github.com/rtomayko/tilt.git

# si prefiere HTTPS a los protocolos git / SSH:
$ git config --global hub.protocol https
$ hub clon rtomayko / tilt
# => clon de git https://github.com/rtomayko/tilt.git
Consulte los ejemplos de uso o la documentación de referencia completa para ver todos los comandos y banderas disponibles.

hub también se puede utilizar para crear scripts de shell que interactúen directamente con la API de GitHub .

hub se puede asignar de forma segura al alias como git, por lo que puede escribir $ git <command>en el shell y ampliarlo con hubfunciones.

Instalación
El hubejecutable no tiene dependencias, pero dado que fue diseñado para ajustarse git, se recomienda tener al menos git 1.7.3 o más reciente.

plataforma	gerente	comando para ejecutar
macOS, Linux	Homebrew	brew install hub
macOS, Linux	Nada	nix-env -i hub
Ventanas	Cucharón	scoop install hub
Ventanas	Chocolatey	choco install hub
Fedora Linux	DNF	sudo dnf install hub
Arch Linux	pacman	sudo pacman -S hub
FreeBSD	paquete (8)	pkg install hub
Debian	apto (8)	sudo apt install hub
Ubuntu	Chasquido	Ya no recomendamos instalar el complemento.
openSUSE	Zypper	sudo zypper install hub
Linux vacío	xbps	sudo xbps-install -S hub
Gentoo	Porteo	sudo emerge dev-vcs/hub
ninguna	conda	conda install -c conda-forge hub
Los paquetes que no sean Homebrew son mantenidos por la comunidad (¡gracias!) Y no se garantiza que coincidan con la última versión del hub . Verifique hub versiondespués de instalar un paquete comunitario.

Ser único
hubse puede instalar fácilmente como ejecutable. Descargue el binario más reciente para su sistema y colóquelo en cualquier lugar de su ruta ejecutable.

Acciones de GitHub
hub está listo para usarse en los flujos de trabajo de Acciones de GitHub :

pasos :
- usos : acciones / pago @ v2

- nombre : Lista de solicitudes de extracción abiertas 
  ejecutar : hub pr list 
  env :
     GITHUB_TOKEN : $ {{secrets.GITHUB_TOKEN}}
Tenga en cuenta que el valor predeterminado secrets.GITHUB_TOKENsolo funcionará para las operaciones de API en el ámbito del repositorio que ejecuta este flujo de trabajo. Si necesita interactuar con otros repositorios, genere un token de acceso personal con al menos el repoalcance y agréguelo a los secretos de su repositorio .

Fuente
Los requisitos previos para construir desde la fuente son:

make
Ir 1.11+
Clona este repositorio y ejecuta make install:

git clone \
  --config transfer.fsckobjects = false \
  --config receive.fsckobjects = false \
  --config fetch.fsckobjects = false \
  https://github.com/github/hub.git

hub de cd
hacer instalar prefijo = / usr / local
Aliasing
Algunas funciones del concentrador se sienten mejor cuando tienen un alias git. Esto no es peligroso; todos tus comandos de git normales funcionarán . hub simplemente agrega un poco de azúcar.

hub aliasmuestra instrucciones para el shell actual. Con la -sbandera, genera un script adecuado para eval.

Debe colocar este comando en su .bash_profileu otro script de inicio:

eval  " $ ( hub alias -s ) "
Potencia Shell
Si está usando PowerShell, puede establecer un alias para hubcolocando lo siguiente en su perfil de PowerShell (generalmente ~/Documents/WindowsPowerShell/Microsoft.PowerShell_profile.ps1):

Set-Alias ​​git hub
Una forma sencilla de hacer esto es ejecutar lo siguiente desde el indicador de PowerShell:

Agregar contenido $ PROFILE  " ` nSet-Alias ​​git hub "
Nota: Deberá reiniciar su consola de PowerShell para que los cambios se recojan.

Si su perfil de PowerShell no existe, puede crearlo ejecutando lo siguiente:

New-Item -Type file -Force $ PROFILE
Finalización de pestaña de Shell
El repositorio de hub contiene scripts de finalización de tabulaciones para bash, zsh y fish. Estos scripts complementan los scripts de finalización existentes que se envían con git.

Meta
Errores: https://github.com/github/hub/issues
Autores: https://github.com/github/hub/contributors
Nuestro código de conducta
