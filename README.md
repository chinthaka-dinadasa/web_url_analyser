[![Golang Web Analyser - CI Prod](https://github.com/chinthaka-dinadasa/web_url_analyser/actions/workflows/backend-ci-cd-production.yml/badge.svg)](https://github.com/chinthaka-dinadasa/web_url_analyser/actions/workflows/backend-ci-cd-production.yml)

[![Next JS Frontend - CI Prod](https://github.com/chinthaka-dinadasa/web_url_analyser/actions/workflows/frontend-ci-cd-production.yml/badge.svg)](https://github.com/chinthaka-dinadasa/web_url_analyser/actions/workflows/frontend-ci-cd-production.yml)

## Web Analyser API

![application run demo](ezgif-6de383e0f5eb32bc.gif)

### How to run the application - API

Clone the Repo

Environment Variables

```
PORT=8080
MAX_WORKERS=50
CACHE_TTL=1 //Hours
```

```bash
go run main.go
```

### How to run the application - Next JS UI

Navigate to frontend/web-analyser-web - Node.js version ">=20.9.0" is required

Environment variables

```
NEXT_PUBLIC_APP_API_URL=http://localhost:8080/process-web-url
```

```
npm install
npm run dev

```

### Test Golang Code with Coverage

```bash
go test -cover -v ./...
```
