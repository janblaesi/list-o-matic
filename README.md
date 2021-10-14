# List-O-Matic #

Das Projekt ist von der Architektur relativ einfach gehalten. Das Backend ist hier mittels Go und Gin-Gonic als HTTP-Library realisiert.
Das Frontend, welches in Angular geschrieben ist (https://github.com/janblaesi/list-o-matic-frontend) greift mit REST-Calls auf die Endpunkte zu.
Diese Endpunkte sind in routes.go implementiert.

Das Datenmodell ist in models.go implementiert und wird zur Persistenz (welche in db.go implementiert ist) in JSON-Dateien exportiert, dessen Pfad vom Nutzer in config.yml gesetzt werden kann. Ebenso kann vom Nutzer der Pfad der users.json Datei angegeben werden, welche die Benutzerdatenbank enthält. Hier werden Benutzername, Passwort-Hash (SHA-256) sowie das Admin-Flag gespeichert. Unter users.example.json liegt ein Beispiel vor, in dem der Benutzername und Passwort des einizigen existenten Benutzers 'admin' sind.

## Entwicklungsumgebung einrichten ##

Im Folgenden ist erklärt, wie eine Umgebung für List-O-Matic eingerichtet werden kann, falls Anpassungen am Code erfolgen sollen.

### Aufsetzen der Entwicklungsumgebung ###

Es wird eine Go-Umgebung benötigt, diese kann über den Paketmanager des lokalen Systems installiert werden.
Mit ```go install``` können die Abhängigkeiten dieses Projektes installiert werden
Mit ```go run .``` kann dann der Webserver gestartet werden (im Debug-Modus)

### Herstellen einer Release-Version ###

In der Entwicklungs-Go-Umgebung kann durch Ausführen folgender zwei Kommandos eine Release-Version erstellt werden.
Aufgrund der Architektur von Go sind alle benötigten Libraries bereits enthalten:
```
export GIN_MODE=release
go build .
```

## Frontend einrichten ##

Ist das Repository des Frontends ausgecheckt, so können mit ```npm install``` alle NodeJS-Dependencies installiert werden, die benötigt werden, um das Frontend zu bauen.

Der Pfad zur API sollte in environment.prod.ts eingetragen werden, damit das Frontend auf das Backend zugreifen kann. Alternativ kann hier der Standardwert verwendet werden, der von einer Installation wie im nächsten Abschnitt beschrieben ausgeht.
Im Folgenden kann mit ```npx ng build --configuration production``` ein Archiv erzeugt werden, mit dem das System auf einem Server installiert werden kann.

## Deployment ##

Diese kurze Anleitung zeigt kurz auf, wie List-O-Matic auf einem Debian-Server mit nginx in Betrieb gesetzt werden kann.
Der nginx-Webserver liefert dann die Angular-Dateien bzw. das Frontend statisch aus und leitet Anfragen an die API (alles unter Pfad /api) an den Gin-Gonic-Webserver weiter (Reverse-Proxy).
Eine Seitenkonfiguration für einen solchen nginx-Server würde in etwa so aussehen:
```
root <Pfad zu Frontend-Dateien>;
index index.html;

server_name <Domain>;

location / {
        try_files $uri$args $uri$args/ /index.html;
}

location /api/ {
        proxy_pass http://localhost:8080/;
}

location = / {
        return 301 https://$host/lists;
}
```
