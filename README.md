# NutriPlan

O aplicație desktop pentru nutriționiști/antrenori personali — gestionează clienți, construiește planuri alimentare (zile → mese → alimente/suplimente) și exportă-le ca PDF-uri gata de trimis către client.

Construită cu [Wails v2](https://wails.io) (Go + WebView nativ), frontend [Svelte 5](https://svelte.dev) + TypeScript, bază de date SQLite locală.

<!--
  TODO: adaugă aici 2-3 screenshot-uri reale din aplicație, de exemplu:

  ![Listă clienți](docs/screenshots/clients.png)
  ![Editor plan](docs/screenshots/plan-editor.png)
  ![Setări](docs/screenshots/settings.png)
-->

## Funcționalități

- **Clienți** — CRUD complet, căutare și filtrare după obiectiv (slăbire, masă musculară, menținere, recompoziție)
- **Planuri alimentare** — structură ierarhică: plan → variante de zi (zi de antrenament/repaus/etc.) → mese → alimente (cu alternativă opțională) și suplimente
- **Duplică plan** — copiază instant un plan existent (cu tot cu zile/mese/alimente) ca punct de plecare pentru un client nou sau o lună nouă
- **Export PDF** — generează un PDF profesional al planului, cu temă complet personalizabilă (culori, font) din Setări; fiecare client primește propriul subfolder de export
- **Previzualizare PDF** — randează planul într-un fișier temporar și îl deschide direct în vizualizatorul de PDF implicit al sistemului, fără să atingă exportul final
- **Setări**:
  - profil (nume + profesie afișate în bara laterală, cerute la prima pornire)
  - temă PDF (10 culori personalizabile, cu resetare la valorile implicite)
  - locul de export (dialog nativ de alegere folder)
  - aspect aplicație — temă deschisă/întunecată/automat (sincronizată cu sistemul)
  - limbă — română/engleză, aplicată atât interfeței cât și etichetelor din PDF-ul generat
  - backup — arhivează baza de date locală într-un `.zip` cu timestamp
- **Logging structurat** — `log/slog`, cu rotație de fișiere (JSON/text, consolă/fișier/ambele)
- **Izolare pe mediu** — `development`/`testing`/`production` folosesc automat baze de date și foldere de export separate, ca să nu se amestece date de test cu date reale de client

## Stack tehnic

| Strat | Tehnologie |
|---|---|
| Backend | Go 1.26, [Wails v2](https://wails.io) |
| Bază de date | SQLite (driver pur Go, [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite), fără CGO) |
| ORM ușor | [`crud-depot`](https://github.com/ciprianiordache/crud-depot) + [`schema-builder`](https://github.com/ciprianiordache/schema-builder) |
| Generare PDF | [`maroto/v2`](https://github.com/johnfercher/maroto) |
| Frontend | Svelte 5 (runes) + TypeScript, Vite, Tailwind CSS v4 |
| Logging | `log/slog` + [`lumberjack`](https://github.com/natefinch/lumberjack) (rotație) |

## Arhitectură

Backend-ul urmează consecvent un pattern pe straturi, per domeniu (`client`, `plan`, `day_variant`, `meal`, `food`, `supplement`, `pdf_export`, `settings`):

```
internals/app/<domeniu>/
  model.go     — struct-uri de domeniu + Request/Response DTO-uri
  repo.go      — acces la bază de date (crud-depot)
  service.go   — logică de business, validare, logging
  handler.go   — expus către frontend prin Wails binding
  errors.go
  *_test.go
```

Pachete transversale:
- `internals/database` — conexiune SQLite, rezolvare cale fișier
- `internals/generator` — toolkit generic de generare PDF (temă, layout), independent de domeniu
- `internals/logger` — construcție logger structurat din config
- `internals/osutil` — deschidere fișiere/foldere în aplicația implicită a sistemului de operare
- `config` — încărcare `config.yaml`/`.env`, izolare pe mediu

## Rulare locală

### Cerințe

- [Go](https://go.dev) 1.26+
- [Node.js](https://nodejs.org) 22+
- [Wails CLI](https://wails.io/docs/gettingstarted/installation) v2

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Configurare

Copiază/editează `config.yaml` din rădăcina proiectului. La prima rulare, dacă `config.yaml`/`.env` nu există, se generează automat cu valori implicite.

Câmpuri relevante:
```yaml
app:
    environment: development   # development / testing / production — izolează DB și exportul per mediu
database:
    driver: sqlite
    name: nutrition-planner
generator:
    # temă implicită PDF — suprascrisă de Setări din aplicație, dacă sunt setate
logger:
    output: both                # stdout / file / both
```

### Dezvoltare

```bash
wails dev
```

Pornește un server Vite cu hot-reload pentru frontend, plus aplicația Wails. Bindings-urile TypeScript pentru metodele Go expuse (`frontend/wailsjs/go/...`) se regenerează automat la fiecare pornire.

### Build de producție

```bash
wails build
```

Rezultă executabilul în `build/bin/`. **Important:** `config.yaml`/`.env` nu sunt copiate automat lângă executabil — vezi `build.ps1` din rădăcina proiectului pentru un script care rulează `wails build` și copiază configul automat.

### Popularea cu date de test

```bash
go run ./cmd/seed
```

Populează baza de date curentă (respectând `app.environment`) cu clienți, planuri, mese și alimente de test.

## Testare

**Backend:**
```bash
go build ./...
go vet ./...
go test ./...
```

**Frontend:**
```bash
cd frontend
npm ci
npm run check   # svelte-check (verificare de tipuri)
npm run build   # build de producție Vite
```

## Structura proiectului

```
├── app.go                    # wiring: config → DB → repo → service → handler
├── main.go                   # entrypoint Wails
├── config/                   # încărcare config.yaml/.env
├── internals/
│   ├── app/<domeniu>/        # câte un pachet per domeniu (vezi Arhitectură)
│   ├── database/
│   ├── generator/            # toolkit generic de PDF
│   ├── logger/
│   ├── osutil/
│   └── testutil/             # bază de date SQLite in-memory pentru teste
├── cmd/seed/                 # populare cu date de test
├── frontend/
│   ├── src/
│   │   ├── lib/               # componente reutilizabile + i18n, api, stores
│   │   └── views/              # ecranele principale (Clienți, Plan, Setări)
│   └── wailsjs/                # bindings generate automat de Wails (comise în git)
└── build.ps1                  # wails build + copiere config.yaml/.env
```
