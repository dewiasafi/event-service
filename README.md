```markdown
# Event Service API

Backend REST API for an event management and ticket booking system (**Event Management & Booking System**) built using **Go (Golang)**, **Gin Web Framework**, **GORM**, and **PostgreSQL**. Equipped with JWT-based authentication, IDOR security protection, search, pagination, booking management, and **ImageKit** cloud storage integration for event image handling.

---

## 🚀 Key Features
* **Secure Authentication (JWT & bcrypt)**: Registration, login, and middleware protection for private routes.
* **Event Management (CRUD)**: 
  * Case-insensitive event search using `ILIKE`.
  * Automatic pagination (`page` & `limit`).
  * Ownership protection to prevent Insecure Direct Object References (IDOR) vulnerabilities.
* **Booking Management**: 
  * Automatic unique booking code generation (`BookingCode`).
  * Double-booking validation to prevent users from booking the same event twice.
  * Secure booking deletion based on user ownership.
* **Image Upload (ImageKit)**: Upload and manage event image assets directly to ImageKit cloud.

---

## 🛠️ Tech Stack
* **Language:** Go (Golang)
* **Web Framework:** Gin (`github.com/gin-gonic/gin`)
* **ORM:** GORM (`gorm.io/gorm`)
* **Database:** PostgreSQL
* **Authentication:** JSON Web Tokens (JWT)
* **Cloud Storage:** ImageKit API
* **Environment Configuration:** `go-dotenv`

---

## ⚙️ System Requirements
Ensure you have installed the following tools on your machine:
* Go (version 1.18 or newer)
* PostgreSQL
* ImageKit Account

---

## 📥 Installation & Running Guide

1. **Clone this repository**:
   ```bash
   git clone [https://github.com/dewiasafi/event-service.git](https://github.com/dewiasafi/event-service.git)
   cd event-service

```

2. **Create a `.env` file** in the root directory of the project, and fill it with the following configuration:
```env
DB_URI=
JWT_SECRET=
IMAGEKIT_PUBLIC_KEY=
IMAGEKIT_PRIVATE_KEY=
IMAGEKIT_URL_ENDPOINT=

```

3. **Install Go dependencies**:
```bash
go mod tidy

```

4. **Run the application**:
```bash
go run .

```

---

## --- API ROUTES ---

### 1. Auth (Public Routes)

* `POST /api/auth/register`
* `POST /api/auth/login`

### 2. Protected Routes (Requires Header: `Authorization: Bearer <token>`)

#### User Profile

* `GET /api/auth/profil-user`

#### Events Management

* `GET /api/events` (Query params: `?search=keyword&page=1&limit=10`)
* `POST /api/event` (Create event, supports image upload)
* `GET /api/event/:id` (Detail event by ID)
* `PUT /api/event/:id` (Update event by ID)
* `DELETE /api/event/:id` (Delete event by ID)

#### Booking Management

* `POST /api/booking` (Create booking, Payload: `{ "phone": "...", "eventId": ... }`)
* `GET /api/booking/user` (Get bookings by logged-in user)
* `DELETE /api/booking/:id` (Delete booking by ID)

