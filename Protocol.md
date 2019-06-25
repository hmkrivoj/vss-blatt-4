# Kommunikation zwischen den Services
## Liste der Services
- `cinema.cinema_hall.service` - Kinosaalverwaltung
- `cinema.movie.service` - Filmverwaltung
- `cinema.cinema_showing.service` - Vorstellungsverwaltung
- `cinema.user.service` - Benutzerverwaltung
- `cinema.reservations.service` - Reservierungsverwaltung
## Liste der Topics
Über die Publish/Subscribe-Topics wird die Eventual Consistency sicher gestellt 
und gleichzeitig eine Umkehr der Abhängigkeiten erzielt. 
So kann die Kinosaalverwaltung beispielsweise über die Löschung der Kinosäle informieren, 
ohne die betroffenen Services kennen zu müssen. 
Stattdessen lauschen die betroffenen Services selbst auf dem jeweiligen Topic.
- `cinema.cinema_hall.deleted` - Topic über gelöschte Kinosäle
- `cinema.movie.deleted` - Topic über gelöschte Filme
- `cinema.cinema_showing.deleted` - Topic über gelöschte Vorstellungen
## Kommunikation zwischen den Services
### `cinema.cinema_hall.service`
- Kommuniziert nicht aktiv mit anderen Services -> Ist unabhängig von den anderen Services
- Sendet beim Löschen von Kinosälen auf `cinema.cinema_hall.deleted`
### `cinema.movie.service`
- Kommuniziert nicht aktiv mit anderen Services -> Ist unabhängig von den anderen Services
- Sendet beim Löschen von Filmen auf `cinema.movie.deleted`
### `cinema.cinema_showing.service`
- Kommuniziert nicht aktiv mit anderen Services, lauscht aber auf 
`cinema.cinema_hall.deleted` und `cinema.movie.deleted`
- Sendet beim Löschen von Vorstellungen auf `cinema.cinema_showing.deleted`
### `cinema.reservations.service`
- Kommuniziert aktiv mit `cinema.cinema_showing.service` und `cinema.cinema_hall.service`
- Lauscht weiterhin auf `cinema.cinema_showing.deleted`
### `cinema.user.service`
- Kommuniziert aktiv mit `cinema.reservations.service`
