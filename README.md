# 📈 TradeLog Trading Journal API

Backend de alto rendimiento diseñado para el registro, seguimiento y análisis profesional de operaciones bursátiles. Este proyecto implementa una arquitectura segura con **Tenant Isolation**, cálculos financieros de precisión (decimales) y un motor de estadísticas avanzado similar a MetaTrader 5.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Gin Framework](https://img.shields.io/badge/Gin-Framework-ff5a5f?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Supabase-336791?style=flat&logo=postgresql)
![Architecture](https://img.shields.io/badge/Architecture-Clean-green)

## 🚀 Tecnologías

* **Lenguaje:** Golang
* **Web Framework:** Gin Gonic
* **ORM:** GORM (Driver Postgres)
* **Base de Datos & Auth:** Supabase (PostgreSQL + Auth Helpers)
* **Storage:** Supabase Storage (Buckets Públicos para Screenshots)
* **Matemática:** `shopspring/decimal` para cálculos financieros sin errores de punto flotante.
* **Seguridad:** Validación de JWT asimétrica (JWKS) y Middleware de Roles.

## 🏗 Arquitectura

El proyecto sigue una estructura modular para separar la lógica de negocio, los controladores y los servicios de cálculo:

```text
├── api
│   ├── config       # Configuración de variables de entorno
│   ├── database     # Conexión Singleton a BD (GORM)
│   ├── domains      # Modelos de datos (User, Trade)
│   ├── handlers     # Controladores HTTP
│   │   ├── admin       # Gestión de Usuarios (Backoffice)
│   │   ├── dashboard   # Métricas y Calendario
│   │   ├── health      # Health Checks
│   │   ├── trades      # CRUD de Operaciones
│   │   └── users       # Perfil personal
│   ├── middleware   # Auth (JWT) y Roles (AdminOnly)
│   └── services     # Lógica de Negocio Pura
│       └── analytics   # Motor Estadístico (Sharpe, Profit Factor, etc.)
│
└── main.go          # Punto de entrada y Router

```

Funcionalidades Principales

🧠 Motor de Inteligencia Financiera (Analytics)

Servicio dedicado (analytics/stats.go) que procesa el historial de operaciones para generar reportes nivel institucional.

    • KPIs Avanzados: Profit Factor, Recovery Factor, Sharpe Ratio, Expectancy (Expected Payoff).

    • Análisis de Rachas: Detección de Consecutive Wins/Losses y su impacto monetario.

    • Precisión: Uso estricto de tipos Decimal para evitar errores de redondeo en PnL y Comisiones.

📅 Dashboard & Calendario

    • Heatmap Optimizado: Endpoint dedicado que utiliza SQL Aggregations (GROUP BY, SUM) para obtener métricas diarias sin cargar todos los trades en memoria.

    • Context Aware: Filtrado por rangos de fecha personalizados.

🔐 Seguridad y Gestión de Usuarios

    • Tenant Isolation: Middleware y consultas diseñadas para asegurar que un usuario solo acceda a sus propias operaciones.

    • Whitelist Strategy: Los usuarios nuevos nacen con estado INACTIVO y requieren aprobación manual.

    • Panel de Admin: Endpoints protegidos para listar, aprobar (KYC/Whitelist), editar roles y banear usuarios.

    • Perfil: Sincronización automática de Avatar y Nombre desde Google OAuth.

📈 Gestión de Trades (Journal)

    • CRUD Completo: Registro detallado (Entrada, Salida, Lotes, Comisiones, Notas).

    • Soporte Multimedia: Integración con Supabase Storage para guardar capturas de pantalla (URLs públicas).

    • Normalización: Manejo automático de Enums (LONG/SHORT, OPEN/CLOSED) y Arrays (Tags) compatibles con PostgreSQL.

🛠️ Instalación y Configuración

    • Clonar el repositorio:

        git clone [https://github.com/GalahadDev/TradeLog-Backend/tree/main]
        cd samll-trading-back
    
    • Configurar Variables de Entorno: Crea un archivo .env en la raíz del proyecto:

        PORT=
        DB_HOST=
        DB_USER=
        DB_PASSWORD=
        DB_NAME=
        DB_PORT=
    
    • Instalar Dependencias:

        go mod tidy

    • Ejecutar el Servidor:

        go run main.go

## 📡 Endpoints API

A continuación se detallan las rutas disponibles. **Todas las rutas `/api` requieren Header `Authorization: Bearer <TOKEN>`**.

### 🩺 Health Check

| Método | Endpoint | Descripción | Nivel de Acceso |
| --- | --- | --- | --- |
| `GET` | `/health` | Verificar estado del servidor y BD | 🟢 Público |

### 📊 Dashboard & Analytics

| Método | Endpoint | Descripción | Nivel de Acceso |
| --- | --- | --- | --- |
| `GET` | `/api/dashboard/stats` | **KPIs Financieros** (WinRate, Profit Factor, Expectancy, etc.) | 🔵 Usuario |
| `GET` | `/api/dashboard/calendar` | Resumen diario (PnL y Count) para Heatmap | 🔵 Usuario |

### 📈 Operaciones (Trades)

| Método | Endpoint | Descripción | Nivel de Acceso |
| --- | --- | --- | --- |
| `GET` | `/api/trades` | Listar operaciones (Paginado `?page=1&limit=10`) | 🔵 Usuario (Tenant) |
| `GET` | `/api/trades/:id` | Ver detalle de una operación | 🔵 Usuario (Tenant) |
| `POST` | `/api/trades` | Registrar nueva operación | 🔵 Usuario |
| `PATCH` | `/api/trades/:id` | Editar operación (Notas, Precios, Status) | 🔵 Usuario (Tenant) |
| `DELETE` | `/api/trades/:id` | Eliminar operación | 🔵 Usuario (Tenant) |

### 👤 Usuario (Perfil)

| Método | Endpoint | Descripción | Nivel de Acceso |
| --- | --- | --- | --- |
| `GET` | `/api/users/me` | Obtener mi perfil y estado | 🔵 Usuario |
| `PATCH` | `/api/users/me` | Actualizar bio, teléfono, foto, experiencia | 🔵 Usuario |

### 🛡️ Panel de Administración

| Método | Endpoint | Descripción | Nivel de Acceso |
| --- | --- | --- | --- |
| `GET` | `/api/admin/users` | Listar todos los usuarios | 🔴 Admin |
| `GET` | `/api/admin/users/:id` | Ver detalle completo de usuario | 🔴 Admin |
| `PUT` | `/api/admin/users/:id` | **Aprobar/Banear**, cambiar Roles y editar perfil | 🔴 Admin |
| `DELETE` | `/api/admin/users/:id` | Eliminar usuario (Soft Ban) | 🔴 Admin |

Desarrollado para los traders disciplinados.
