# Go Product & Category API

A RESTful API service built with Go for managing products and categories, using PostgreSQL (Supabase) as the database and GORM as the ORM.

## 🚀 Features

- **Product Management**: CRUD operations for products
- **Category Management**: CRUD operations for categories
- **Relationship**: Products belong to categories (foreign key constraint)
- **Database**: PostgreSQL via Supabase
- **Auto Migration**: Database tables are created automatically
- **Health Check**: Endpoint to monitor service status
- **Clean Architecture**: Repository pattern with separation of concerns

## 📋 Tech Stack

- **Language**: Go 1.24.0
- **Database**: PostgreSQL (Supabase)
- **ORM**: GORM
- **Configuration**: Viper
- **Development**: Air (live-reload)

## 📁 Project Structure

```
go-product-supabase/
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection & migration
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   ├── repository/      # Data access layer
│   └── services/        # Business logic layer
├── migrations/          # Database migrations
├── .air.toml           # Air configuration for live-reload
├── .env                # Environment variables
├── go.mod              # Go dependencies
└── main.go             # Application entry point
```

## 🔧 Prerequisites

- Go 1.24.0 or higher
- PostgreSQL database (Supabase account)
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
- Get it from: Supabase Dashboard → Project Settings → Database → Connection String
- Choose: URI → Direct connection

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
```http
GET    /api/categories       # Get all categories
POST   /api/categories       # Create category
GET    /api/categories/{id}  # Get category by ID
PUT    /api/categories/{id}  # Update category
DELETE /api/categories/{id}  # Delete category
```

### Products
```http
GET    /api/products                      # Get all products
GET    /api/products?category_id={id}     # Get products by category
POST   /api/products                      # Create product
GET    /api/products/{id}                 # Get product by ID
PUT    /api/products/{id}                 # Update product
DELETE /api/products/{id}                 # Delete product
```

## 📝 Request Examples

### Create Category
```json
POST /api/categories
{
  "name": "Electronics",
  "description": "Electronic devices and accessories"
}
```

### Create Product
```json
POST /api/products
{
  "name": "Laptop",
  "price": 999.99,
  "stock": 10,
  "category_id": 1
}
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
   - Your API will be available at: `https://go-product-api.onrender.com`
   - Health check: `https://go-product-api.onrender.com/health`

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

### Categories Table
```sql
CREATE TABLE categories (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  description TEXT
);
```

### Products Table
```sql
CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(200) NOT NULL,
  price DECIMAL(10,2) NOT NULL,
  stock INTEGER DEFAULT 0,
  category_id INTEGER NOT NULL,
  FOREIGN KEY (category_id) REFERENCES categories(id)
);
```

## 🧪 Testing

Use the included `request.http` file with REST Client extension in VS Code, or use tools like:
- Postman
- Insomnia
- cURL

Example with cURL:
```bash
# Health check
curl http://localhost:6000/health

# Get all categories
curl http://localhost:6000/api/categories

# Create category
curl -X POST http://localhost:6000/api/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Electronics","description":"Electronic items"}'
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 👤 Author

Your Name - [Your GitHub Profile](https://github.com/yourusername)

## 🙏 Acknowledgments

- [GORM](https://gorm.io/) - ORM library
- [Viper](https://github.com/spf13/viper) - Configuration management
- [Supabase](https://supabase.com/) - PostgreSQL hosting
- [Render](https://render.com/) - Deployment platform
- [Air](https://github.com/air-verse/air) - Live reload for Go apps
