# Open Source UoM — Static Site

Ο ιστότοπος της Κοινότητας Ανοιχτού Λογισμικού του Πανεπιστημίου Μακεδονίας.  
Χτισμένος με έναν **custom static site generator γραμμένο σε Go**, χωρίς εξωτερικές εξαρτήσεις framework.

🌐 **Live site:** https://open-source-uom.github.io/

---

## Πίνακας Περιεχομένων

- [Πώς λειτουργεί](#πώς-λειτουργεί)
- [Δομή directories](#δομή-directories)
- [Προαπαιτούμενα](#προαπαιτούμενα)
- [Τοπική εκτέλεση (run locally)](#τοπική-εκτέλεση-run-locally)
- [Build του site](#build-του-site)
- [Deployment (GitHub Actions)](#deployment-github-actions)
- [Αλλαγές στη σχεδίαση (CSS / HTML)](#αλλαγές-στη-σχεδίαση-css--html)
- [Δημιουργία νέου άρθρου (blog post)](#δημιουργία-νέου-άρθρου-blog-post)
- [Δημιουργία νέας σελίδας (page)](#δημιουργία-νέας-σελίδας-page)
- [Εικόνες και media](#εικόνες-και-media)

---

## Πώς λειτουργεί

Το σύστημα αποτελείται από:

1. **Content** (`content/`): Αρχεία Markdown (`.md`) που περιέχουν τα άρθρα και τις σελίδες.
2. **Generator** (`generator/`): Ένα πρόγραμμα Go που διαβάζει τα Markdown αρχεία και παράγει στατική HTML.
3. **Templates** (`generator/templates/`): HTML templates που ορίζουν τη δομή κάθε σελίδας.
4. **Static assets** (`generator/static/`): CSS και άλλα στατικά αρχεία.
5. **Storage** (`storage/`): Εικόνες και media αρχεία των άρθρων.
6. **Public** (`public/`): Ο φάκελος εξόδου — αυτό που ανεβαίνει στο GitHub Pages.

### Ροή εκτέλεσης

```
content/*.md  +  storage/  +  templates/  +  static/
         │
         ▼
    go run .   (μέσα από τον φάκελο generator/)
         │
         ▼
      public/   ←── αυτό σερβίρεται στο GitHub Pages
```

Ο generator:
- Διαβάζει όλα τα `.md` αρχεία από `content/blog/` και `content/pages/`
- Κάνει parse το **YAML frontmatter** (τίτλος, ημερομηνία, συγγραφέας, tags, εικόνα)
- Μετατρέπει το **Markdown body** σε HTML (χρησιμοποιεί τη βιβλιοθήκη `goldmark`)
- Εφαρμόζει τα HTML templates και παράγει τα αρχεία στον φάκελο `public/`
- Αντιγράφει τα `storage/` και `generator/static/` μέσα στο `public/`

---

## Δομή directories

```
opensource.uom.gr/
│
├── content/                  # Περιεχόμενο σε Markdown
│   ├── blog/                 # Άρθρα blog
│   │   ├── fosscomm-2024.md
│   │   ├── git-advice-for-freshmen.md
│   │   └── ...
│   └── pages/                # Στατικές σελίδες
│       ├── register.md
│       ├── xorigoi.md
│       ├── videos.md
│       └── ...
│
├── generator/                # Ο static site generator (Go)
│   ├── main.go               # Entry point, ορίζει paths & BASE_URL
│   ├── generator.go          # Κύρια λογική παραγωγής HTML
│   ├── parser.go             # Ανάγνωση & parsing Markdown + frontmatter
│   ├── go.mod                # Go module dependencies
│   ├── go.sum
│   │
│   ├── templates/            # HTML templates
│   │   ├── base.html         # Βασικό layout (header, footer, fonts, nav)
│   │   ├── index.html        # Αρχική σελίδα (homepage)
│   │   ├── post.html         # Σελίδα μεμονωμένου άρθρου
│   │   ├── blog-index.html   # Λίστα όλων των άρθρων
│   │   └── page.html         # Στατική σελίδα (register, xorigoi, κτλ.)
│   │
│   └── static/               # Στατικά assets (αντιγράφονται στο public/)
│       └── css/
│           └── style.css     # Το κύριο stylesheet
│
├── storage/                  # Εικόνες και media αρχεία
│   ├── 2022/
│   ├── 2023/
│   ├── 2024/
│   └── 2025/
│
├── public/                   # Έξοδος (auto-generated, μην το επεξεργαστείς)
│   └── ...
│
└── .github/
    └── workflows/
        └── deploy.yml        # GitHub Actions: auto-build & deploy
```

> ⚠️ **Ποτέ μην επεξεργάζεσαι τον φάκελο `public/` χειροκίνητα.**  
> Δημιουργείται αυτόματα σε κάθε build και τα manual αλλαγές χάνονται.

---

## Προαπαιτούμενα

Για τοπική ανάπτυξη χρειάζεσαι:

- **Go** >= 1.22 → https://go.dev/dl/
- **Python 3** (για τον local web server) — προεγκατεστημένο στο Linux/macOS
- **Git**

Για να ελέγξεις αν έχεις Go:
```bash
go version
```

---

## Τοπική εκτέλεση (run locally)

### Βήμα 1: Κλωνοποίηση του repository
```bash
git clone https://github.com/open-source-uom/open-source-uom.github.io.git
cd open-source-uom.github.io
```

### Βήμα 2: Build του site
```bash
cd generator
go mod tidy   # Εγκατάσταση dependencies (μόνο την πρώτη φορά)
go run .
```

Θα δεις κάτι σαν:
```
Starting Open Source UoM Static Site Generator...
Found 68 posts and 6 pages.
Generating static HTML...
Copying static assets...
Site generation complete! Output is in public/
```

### Βήμα 3: Εκκίνηση local web server
```bash
cd ../public
python3 -m http.server 8080
```

Άνοιξε τον browser σου και πήγαινε στο:  
👉 **http://localhost:8080/**

> **Σημαντικό:** Άνοιγε πάντα μέσω `http://localhost:8080` και **όχι** με διπλό κλικ στο αρχείο `index.html`. Αν το ανοίξεις απευθείας ως αρχείο (`file://...`), το CSS δεν θα φορτώσει γιατί τα links ξεκινάνε με `/`.

---

## Build του site

Κάθε φορά που κάνεις αλλαγές (νέο άρθρο, αλλαγές CSS, κτλ.), τρέξε ξανά:

```bash
cd generator
go run .
```

Το `public/` ενημερώνεται αυτόματα.

---

## Deployment (GitHub Actions)

Το deployment γίνεται **αυτόματα** κάθε φορά που κάνεις `push` στον κλάδο `master`.

### Πώς λειτουργεί το CI/CD

Το αρχείο `.github/workflows/deploy.yml`:
1. Ενεργοποιείται σε κάθε push στο `master`
2. Εγκαθιστά Go
3. Τρέχει `go run .` μέσα στον φάκελο `generator/`
4. Ανεβάζει τον φάκελο `public/` στο GitHub Pages

**Το `BASE_URL` γίνεται inject αυτόματα** από το GitHub Actions, οπότε όλα τα links (CSS, εικόνες, σελίδες) θα είναι σωστά για το συγκεκριμένο repository — δεν χρειάζεται καμία χειροκίνητη αλλαγή.

### Χειροκίνητο deployment

Αν θες να τρέξεις το deployment χειροκίνητα:
1. Πήγαινε στο GitHub repo → καρτέλα **Actions**
2. Επέλεξε **Deploy Static Site**
3. Πάτα **Run workflow** → **Run workflow**

### Push αλλαγών

```bash
git add .
git commit -m "Προσθήκη νέου άρθρου: τίτλος"
git push origin master
```

---

## Αλλαγές στη σχεδίαση (CSS / HTML)

### Αλλαγή στυλ (CSS)

Το μοναδικό αρχείο CSS βρίσκεται στο:
```
generator/static/css/style.css
```

Εδώ μπορείς να αλλάξεις:
- **Χρώματα**: Τα CSS variables στο `:root {}` στην αρχή του αρχείου  
  (π.χ. `--accent-color`, `--bg-color`, `--btn-primary`)
- **Typography**: Font family, font sizes, line-height
- **Layout**: Grid, flexbox, max-widths, padding/margin
- **Κάρτες άρθρων**: `.post-card`, `.feature-card`
- **Κουμπιά**: `.btn`, `.btn-primary`, `.btn-outline`
- **Responsive**: Το `@media (max-width: 768px)` στο τέλος

### Αλλαγή HTML templates

| Template | Τι ελέγχει |
|---|---|
| `generator/templates/base.html` | Το global layout: header, nav, footer, Google Fonts |
| `generator/templates/index.html` | Η αρχική σελίδα (Hero, sections, articles) |
| `generator/templates/post.html` | Η σελίδα κάθε μεμονωμένου άρθρου |
| `generator/templates/blog-index.html` | Η λίστα με όλα τα άρθρα |
| `generator/templates/page.html` | Οποιαδήποτε στατική σελίδα (register, videos, κτλ.) |

Τα templates χρησιμοποιούν το [Go template syntax](https://pkg.go.dev/text/template):
- `{{ .Title }}` → εκτυπώνει τον τίτλο
- `{{ .BaseURL }}` → το base URL του site (π.χ. `/`)
- `{{ range .RecentPosts }} ... {{ end }}` → loop
- `{{ define "content" }} ... {{ end }}` → ορίζει block που εισάγεται στο `base.html`

### Αλλαγή πλήθους άρθρων στην αρχική

Στο αρχείο `generator/generator.go`, γύρω στη γραμμή 153:
```go
recentPosts := g.Posts
if len(recentPosts) > 6 {   // ← άλλαξε το 6 στον αριθμό που θέλεις
    recentPosts = recentPosts[:6]
}
```

---

## Δημιουργία νέου άρθρου (blog post)

### 1. Δημιούργησε ένα νέο `.md` αρχείο

Τα αρχεία blog βρίσκονται στο:
```
content/blog/
```

Δημιούργησε ένα νέο αρχείο με slugified όνομα (πεζά, παύλες αντί για κενά):
```bash
touch content/blog/titlos-tou-arthrou.md
```

### 2. Γράψε το frontmatter και το περιεχόμενο

```markdown
---
title: Τίτλος του άρθρου
date: "2025-09-15T10:00:00+03:00"
author: Όνομα Συγγραφέα
tags:
    - open source
    - linux
featured_image: /storage/2025/09/photo.jpg
---

Το κείμενο του άρθρου σε Markdown.

## Υποτίτλος

Μπορείς να γράψεις **bold**, *italic* και να βάλεις [links](https://example.com).

![Λεζάντα εικόνας](/storage/2025/09/another-photo.jpg)
```

### Frontmatter — Εξήγηση πεδίων

| Πεδίο | Απαραίτητο | Περιγραφή |
|---|---|---|
| `title` | ✅ Ναι | Ο τίτλος του άρθρου |
| `date` | ✅ Ναι | Ημερομηνία σε RFC3339 format: `"YYYY-MM-DDTHH:MM:SS+03:00"` |
| `author` | Προαιρετικό | Το όνομα του συγγραφέα |
| `tags` | Προαιρετικό | Λίστα tags (το κάθε tag σε νέα γραμμή με `-`) |
| `featured_image` | Προαιρετικό | Path της κύριας εικόνας (ξεκινά με `/storage/...`) |

### 3. Build & preview

```bash
cd generator
go run .
cd ../public
python3 -m http.server 8080
```

Το νέο άρθρο θα εμφανιστεί στη διεύθυνση:  
`http://localhost:8080/blog/titlos-tou-arthrou/`

### 4. Push

```bash
git add content/blog/titlos-tou-arthrou.md
git commit -m "Νέο άρθρο: Τίτλος του άρθρου"
git push origin master
```

---

## Δημιουργία νέας σελίδας (page)

Οι σελίδες λειτουργούν ακριβώς όπως τα άρθρα, αλλά βρίσκονται στο:
```
content/pages/
```

Παράδειγμα νέας σελίδας `content/pages/about.md`:
```markdown
---
title: Σχετικά με εμάς
---

Κείμενο της σελίδας σε Markdown.
```

Η σελίδα θα δημιουργηθεί στη διεύθυνση:  
`https://open-source-uom.github.io/about/`

> Οι σελίδες **δεν εμφανίζονται** στη λίστα blog ή στα πρόσφατα άρθρα.  
> Αν θέλεις να εμφανίζεται στο navigation, θα πρέπει να επεξεργαστείς το `generator/templates/base.html`.

---

## Εικόνες και media

### Πού αποθηκεύεις τις εικόνες

Οι εικόνες αποθηκεύονται στον φάκελο:
```
storage/
├── 2022/
├── 2023/
├── 2024/
│   ├── 11/
│   │   ├── FOSSCOMM24_GROUP.jpg
│   │   └── ...
│   └── ...
└── 2025/
    └── ...
```

**Σύμβαση ονομασίας:** Χρησιμοποίησε την δομή `storage/ΧΡΟΝΙΑ/ΜΗΝΑΣ/photo.jpg`  
Π.χ.: `storage/2025/09/workshop-photo.jpg`

### Πώς ανεβάζεις μια εικόνα

1. Αντίγραψε το αρχείο εικόνας στον σωστό φάκελο:
   ```bash
   cp ~/Downloads/photo.jpg storage/2025/09/photo.jpg
   ```

2. Αναφέρσου σε αυτό στο frontmatter ή στο Markdown:
   ```markdown
   # Στο frontmatter (κύρια εικόνα):
   featured_image: /storage/2025/09/photo.jpg
   
   # Μέσα στο κείμενο:
   ![Περιγραφή εικόνας](/storage/2025/09/photo.jpg)
   ```

3. Ο generator αντιγράφει αυτόματα ολόκληρο τον φάκελο `storage/` στο `public/storage/` κατά το build.

### Συμβουλές για εικόνες

- **Format:** Προτίμησε `jpg` για φωτογραφίες και `png` για screenshots.
- **Μέγεθος:** Προτίμησε εικόνες < 1MB. Μπορείς να τις συμπιέσεις με `jpegoptim` ή `optipng`.
- **Διαστάσεις:** Το site εμφανίζει εικόνες σε πλάτος 300px στα post cards. Οι original εικόνες μπορεί να είναι μεγαλύτερες (το CSS φροντίζει για το resizing).

---

## Σύνοψη workflow

```bash
# 1. Δημιουργία νέου άρθρου
nano content/blog/νεο-αρθρο.md

# 2. Ανέβασμα εικόνων (αν χρειάζεται)
cp photo.jpg storage/2025/09/photo.jpg

# 3. Τοπικό build & preview
cd generator && go run . && cd ../public && python3 -m http.server 8080

# 4. Έλεγξε στο http://localhost:8080

# 5. Push για αυτόματο deployment
git add .
git commit -m "Νέο άρθρο: ..."
git push origin master
```

Μετά το push, το GitHub Actions τρέχει αυτόματα (~1-2 λεπτά) και το site ανανεώνεται στο https://open-source-uom.github.io/ 🚀

---

## Άδεια χρήσης

Ο κώδικας του generator είναι open source. Το περιεχόμενο (`content/`) ανήκει στην Κοινότητα Ανοιχτού Λογισμικού του Πανεπιστημίου Μακεδονίας.
