# Badminton Tournament Manager

Full-stack badminton application với backend Go/Gin và frontend Next.js 15.

## Kiến trúc dự án

- **Backend**: Go với Gin framework
- **Frontend**: Next.js 15 với JavaScript
- **Database**: SQLite (development) / PostgreSQL (production)
- **Container**: Docker & Docker Compose
- **Hot Reload**: Air cho Go backend

## Cấu trúc thư mục

```
/
├── backend/
│   ├── cmd/                    # Entry points
│   ├── internal/
│   │   ├── models/            # GORM models (Player, Match, Tournament)
│   │   ├── controllers/       # HTTP handlers
│   │   └── views/            # API response structures
│   ├── pkg/                   # Public packages
│   ├── .air.toml             # Air hot reload config
│   ├── Dockerfile            # Backend container
│   └── go.mod               # Go dependencies
├── frontend/
│   ├── app/                  # Next.js 15 App Router
│   ├── components/           # React components
│   ├── public/              # Static assets
│   ├── Dockerfile           # Frontend container
│   └── package.json         # Node.js dependencies
├── docker-compose.yml       # Production setup
└── docker-compose.dev.yml   # Development setup
```

## Chạy dự án

### Development (khuyến nghị)
```bash
# Clone và cd vào project
git clone <repo-url>
cd badminton

# Chạy với hot reload
docker-compose -f docker-compose.dev.yml up
```

### Alternative - Chạy riêng từng service
```bash
# Backend only (cần Go 1.21+)
cd backend
go mod download
go install github.com/cosmtrek/air@latest
air

# Frontend only (cần Node.js 20+)
cd frontend
npm install
npm run dev
```

### Production
```bash
docker-compose up --build
```

## API Endpoints

Tất cả API có prefix `/api/v1/`:

- **Players**: GET/POST/PUT/DELETE `/api/v1/players`
- **Matches**: GET/POST/PUT/DELETE `/api/v1/matches` 
- **Tournaments**: GET/POST/PUT/DELETE `/api/v1/tournaments`

## Truy cập

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

## Phát triển

- Backend tự động reload khi thay đổi file `.go`
- Frontend tự động reload với Next.js Fast Refresh
- Database SQLite được tạo tự động ở `/backend/badminton.db`
- CORS đã được cấu hình cho local development