# Go Product API — Supabase

A RESTful API service built with Go for managing products, categories, and transactions, using PostgreSQL (Supabase) as the database and GORM as the ORM.

**Live Demo**: [https://go-product-supabase.onrender.com/health](https://go-product-supabase.onrender.com/health)

## 🚀 Features

- **Category Management**: Full CRUD operations for product categories
- **Product Management**: Full CRUD operations for products with category relationship
- **Checkout / Transactions**: Process checkout with stock validation and automatic stock deduction
- **Sales Reports**: Today's sales summary and date-range sales reports with best-selling product info
- **Database**: PostgreSQL via Supabase (with PgBouncer connection pooler support)
- **Auto Migration**: Database tables are created automatically on startup
- **Health Check**: Endpoint to monitor service status
- **Clean Architecture**: Repository → Service → Handler pattern with separation of concerns

## 📋 Tech Stack

- **Language**: Go 1.24
- **Database**: PostgreSQL (Supabase)
- **ORM**: GORM
- **DB Driver**: pgx/v5 (with simple protocol for PgBouncer compatibility)
- **Configuration**: Viper
- **Development**: Air (live-reload)

## 📁 Project Structure

```
go-product-supabase/
├── internal/
│   ├── config/          # Configuration management (Viper)
│   ├── database/        # Database connection, migration & health check
│   ├── handlers/        # HTTP handlers (category, product, transaction)
│   ├── models/          # Data models (Category, Product, Transaction, TransactionDetail)
│   ├── repository/      # Data access layer
│   └── services/        # Business logic layer
├── migrations/          # Auto-migration runner
├── .air.toml            # Air configuration for live-reload
├── .env                 # Environment variables (not committed)
├── .env.example         # Environment variables template
├── request.http         # HTTP request examples (localhost)
├── request-production.http # HTTP request examples (production)
├── go.mod               # Go dependencies
└── main.go              # Application entry point
```

## 🔧 Prerequisites

- Go 1.24 or higher
- PostgreSQL database (Supabase account recommended)
- Git

## ⚙️ Installation

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-product-supabase
```

2. **Install dependencies**
```bash
go mod download
```

3. **Install Air for development (optional)**
```bash
go install github.com/air-verse/air@latest
```

4. **Configure environment variables**

Create a `.env` file in the root directory:

```env
# Server Config
SERVER_PORT=6000
SERVER_HOST=localhost

# Supabase DB URL
DATABASE_URL=postgresql://postgres.[PROJECT-REF]:[PASSWORD]@aws-0-[region].pooler.supabase.com:6543/postgres
```

Replace with your Supabase connection string:
- Get it from: **Supabase Dashboard → Project Settings → Database → Connection String**
- Use **Transaction Mode** (port `6543`) — the app handles PgBouncer compatibility automatically

## 🏃 Running Locally

### With Air (live-reload)
```bash
air
```

### Without Air
```bash
go run main.go
```

The server will start on `http://localhost:6000`

## 📡 API Endpoints

### Health Check
```http
GET /health
```

### Categories
| Method   | Endpoint               | Description           |
|----------|------------------------|-----------------------|
| `GET`    | `/api/categories`      | Get all categories    |
| `POST`   | `/api/categories`      | Create a category     |
| `GET`    | `/api/categories/{id}` | Get category by ID    |
| `PUT`    | `/api/categories/{id}` | Update a category     |
| `DELETE` | `/api/categories/{id}` | Delete a category     |

### Products
| Method   | Endpoint                              | Description               |
|----------|---------------------------------------|---------------------------|
| `GET`    | `/api/products`                       | Get all products          |
| `GET`    | `/api/products?name={name}`           | Filter products by name   |
| `GET`    | `/api/products?category_id={id}`      | Filter products by category |
| `POST`   | `/api/products`                       | Create a product          |
| `GET`    | `/api/products/{id}`                  | Get product by ID         |
| `PUT`    | `/api/products/{id}`                  | Update a product          |
| `DELETE` | `/api/products/{id}`                  | Delete a product          |

### Transactions & Checkout
| Method | Endpoint             | Description                         |
|--------|---------------------|-----------------------------------------|
| `POST` | `/api/checkout`     | Process a checkout (creates transaction, deducts stock) |

### Sales Reports
| Method | Endpoint                                              | Description                    |
|--------|-------------------------------------------------------|--------------------------------|
| `GET`  | `/api/report/today`                                   | Today's sales summary          |
| `GET`  | `/api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` | Sales summary by date range |

## 📝 Request Examples

### Create Category
```bash
curl -X POST http://localhost:6000/api/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Elektronik",
    "description": "Perangkat elektronik seperti ponsel, laptop, dan televisi."
  }'
```

### Create Product
```bash
curl -X POST http://localhost:6000/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 17 Pro",
    "description": "Smartphone premium dari Apple",
    "price": 15999000,
    "stock": 50,
    "category_id": 1
  }'
```

### Checkout
```bash
curl -X POST http://localhost:6000/api/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      { "product_id": 1, "quantity": 2 }
    ]
  }'
```

### Sales Report (Today)
```bash
curl http://localhost:6000/api/report/today
```

### Sales Report (Date Range)
```bash
curl "http://localhost:6000/api/report?start_date=2026-01-01&end_date=2026-02-09"
```

## 🌐 Deployment on Render.com

### Prerequisites
- GitHub account
- Render account (sign up at [render.com](https://render.com))
- Supabase database

### Deployment Steps

1. **Push code to GitHub**
```bash
git add .
git commit -m "Initial commit"
git push origin main
```

2. **Create Web Service on Render**
   - Go to [Render Dashboard](https://dashboard.render.com)
   - Click **"New +"** → **"Web Service"**
   - Connect your GitHub repository
   - Configure the service:

3. **Service Configuration**
   - **Name**: `go-product-api` (or your preferred name)
   - **Environment**: `Go`
   - **Region**: Choose closest to your users
   - **Branch**: `main`
   - **Build Command**: 
     ```bash
     go build -o main .
     ```
   - **Start Command**: 
     ```bash
     ./main
     ```

4. **Environment Variables**

   Add these environment variables in Render:
   
   | Key | Value | Description |
   |-----|-------|-------------|
   | `DATABASE_URL` | Your Supabase connection string | **Required** - PostgreSQL connection string |
   | `AUTO_MIGRATE` | `false` | **Recommended** - Set to `false` to skip auto-migration on deploy |

   > **Note**: 
   > - Get your Supabase connection string from: Supabase Dashboard → Settings → Database → Connection String (URI format)
   > - Render automatically sets `PORT` environment variable (no need to set manually)
   > - Migration should be run manually in Supabase dashboard or use migration tool
   > - The app is configured to work with Supabase connection pooler (handles prepared statement caching automatically)

5. **Deploy**
   - Click **"Create Web Service"**
   - Render will automatically build and deploy your application
   - Wait for deployment to complete (~2-5 minutes)

6. **Access Your API**
   - Your API will be available at: `https://<your-service-name>.onrender.com`
   - Health check: `https://<your-service-name>.onrender.com/health`

### Render Configuration Tips

- **Auto-Deploy**: Enable automatic deployments from GitHub
- **Health Check Path**: Set to `/health` in Render settings
- **Instance Type**: Start with free tier, upgrade if needed
- **Logs**: Monitor logs in Render dashboard for debugging

### Important Notes

⚠️ **Free Tier Limitations**:
- Service spins down after 15 minutes of inactivity
- First request after spin down may take 30-50 seconds
- Upgrade to paid plan for always-on service

## 🗄️ Database Schema

### Categories
| Column        | Type          | Constraints          |
|---------------|---------------|----------------------|
| `id`          | `SERIAL`      | PRIMARY KEY          |
| `name`        | `VARCHAR(100)` | NOT NULL, UNIQUE    |
| `description` | `TEXT`        |                      |

### Products
| Column        | Type           | Constraints                        |
|---------------|----------------|------------------------------------|
| `id`          | `SERIAL`       | PRIMARY KEY                        |
| `name`        | `VARCHAR(200)` | NOT NULL                           |
| `price`       | `DECIMAL(10,2)` | NOT NULL                          |
| `stock`       | `INTEGER`      | DEFAULT 0                          |
| `category_id` | `INTEGER`      | NOT NULL, FK → categories(id)      |

### Transactions
| Column         | Type            | Constraints  |
|----------------|-----------------|--------------|
| `id`           | `BIGSERIAL`     | PRIMARY KEY  |
| `total_amount` | `DECIMAL(10,2)` | NOT NULL     |
| `created_at`   | `TIMESTAMPTZ`   | AUTO         |

### Transaction Details
| Column           | Type            | Constraints                              |
|------------------|-----------------|------------------------------------------|
| `id`             | `BIGSERIAL`     | PRIMARY KEY                              |
| `transaction_id` | `BIGINT`        | NOT NULL, FK → transactions(id) ON DELETE CASCADE |
| `product_id`     | `BIGINT`        | NOT NULL, FK → products(id)              |
| `quantity`        | `BIGINT`       | NOT NULL                                 |
| `subtotal`       | `DECIMAL(10,2)` | NOT NULL                                |

## 🧪 Testing

Use the included HTTP request files with the [REST Client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) extension in VS Code:

- `request.http` — local endpoints (`http://localhost:6000`)
- `request-production.http` — production endpoints (`https://go-product-supabase.onrender.com`)

Or use tools like **Postman**, **Insomnia**, or **cURL**.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 👤 Author

Vektor Muhammad Lutfi — [GitHub](https://github.com/vektormuhammadlutfi)

## 🙏 Acknowledgments

- [GORM](https://gorm.io/) - ORM library
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Supabase](https://supabase.com/) - PostgreSQL hosting
- [Render](https://render.com/) - Deployment platform
- [Air](https://github.com/air-verse/air) - Live reload for Go apps
