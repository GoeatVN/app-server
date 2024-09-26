
Cấu trúc dự án GO-ARC
app-server/
├── cmd/
│   ├── api/
│   │   └── main.go               # Entry point của ứng dụng
│   ├── migrate/
│   │   └── main.go               # Entry point cho database migrations
│   └── worker/
│       └── main.go               # Entry point cho background workers (nếu cần)
├── internal/
│   ├── domain/
│   │   ├── entity/               # Entity (User, Product, Order)
│   │   │   ├── user.go
│   │   │   ├── product.go
│   │   │   └── order.go
│   │   └── repository/           # Interface repo
│   │       ├── user_repository.go
│   │       ├── product_repository.go
│   │       └── order_repository.go
│   ├── usecase/                  # Business logic
│   │   ├── user/
│   │   │   ├── interface.go      # Interface cho use case
│   │   │   └── service.go        # Implement use case (Create, Update...)
│   │   ├── product/
│   │   └── order/
│   ├── interface/
│   │   └── api/
│   │       └── handler/
│   │           ├── v1/           # API version 1
│   │           │   ├── user_handler.go
│   │           │   ├── product_handler.go
│   │           │   └── order_handler.go
│   │           └── v2/           # API version 2 (nếu cần)
│   ├── persistence/              # Data access layer triển khai repo
│   │   └── repository/
│   │       ├── generic_repo.go   # GenericBaseRepo[users] (Create, Update, Delete...)
│   │       └── postgres/
│   │           ├── user_repository.go
│   │           ├── product_repository.go
│   │           └── order_repository.go
│   ├── shared/                   # Shared model request, response, DTOs
│   │   └── model/
│   │       ├── user_model.go
│   │       ├── product_model.go
│   │       └── order_model.go
│   └── infrastructure/
│       ├── database/
│       │   └── postgres.go       # Database connection
│       ├── server/
│       │   └── http.go           # HTTP server setup
│       └── middleware/           # HTTP middleware
│           ├── auth.go           # Xác thực
│           ├── authorization.go  # Xác thực role permission
│           ├── cors.go           # Cors
│           ├── logger.go         # Logging
│           ├── validate.go       # Validate input
│           └── response_handler.go # Xử lý kết quả success và error
├── pkg/                          # Shared packages
│   ├── errors/
│   │   └── errors.go             # Custom error types
│   ├── cache/
│   │   └── redis.go              # Cache with redis
│   ├── email/
│   │   └── email.go              # Send email
│   ├── logger/
│   │   └── logger.go             # Logging utility
│   └── response/
│       └── response.go           # API response structures
├── configs/
│   ├── config.yaml               # Base configuration
│   ├── config.development.yaml   # Development-specific config
│   ├── config.staging.yaml       # Staging-specific config
│   └── config.production.yaml    # Production-specific config
├── migrations/                   # Database migrations
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   ├── 002_create_products_table.up.sql
│   └── 002_create_products_table.down.sql
├── docs/
│   ├── api/
│   │   └── swagger.yaml          # API documentation
│   └── README.md                 # Project documentation
├── scripts/
│   ├── setup.sh                  # Setup script
│   ├── run_tests.sh              # Script to run tests
│   └── build.sh                  # Build script
├── tests/                        # Integration and e2e tests
├── go.mod                        # Go module file
├── go.sum
├── Dockerfile                    # Run with docker
├── docker-compose.yml            # Run with docker compose
├── .gitignore
├── README.md
├── Makefile                      # For common tasks
└── ci/
    └── .github/
        └── workflows/
            └── main.yml          # Run CI/CD with GitHub
Tổng quan tính năng dự án:
Dự án bao gồm các tính năng chính sau:
Middleware:
auth: Xác thực Bearer token từ header HTTP để kiểm tra người dùng hợp lệ.
authorization: Kiểm tra quyền truy cập của người dùng dựa trên vai trò (role-based access control).
response_handler: Xử lý và định dạng kết quả trả về từ API, bao gồm thành công, lỗi, và validate.
errorHandler: Xử lý lỗi phát sinh trong quá trình xử lý và trả về dưới dạng chuẩn APIResponse.
validate: Middleware kiểm tra dữ liệu đầu vào và trả về lỗi nếu dữ liệu không hợp lệ.
logger: Middleware ghi log chi tiết về mỗi yêu cầu HTTP (phương thức, đường dẫn, trạng thái, thời gian xử lý).
API Response:
Chuẩn hóa phản hồi từ API bằng struct APIResponse với các hàm:
Success: Trả về phản hồi thành công.
Error: Trả về lỗi hệ thống hoặc lỗi truy cập.
ValidationError: Trả về lỗi validate dữ liệu đầu vào.
Repository Pattern:
GenericBaseRepository: Một repository tổng quát (generic) để quản lý các thao tác CRUD cơ bản cho mọi entity trong hệ thống.
UserRepository, ProductRepository, OrderRepository: Các repository cụ thể kế thừa từ GenericBaseRepository và có thể mở rộng thêm các phương thức riêng nếu cần (ví dụ: tìm kiếm người dùng theo email).
Use Case:
Business Logic: Được triển khai theo từng entity (User, Product, Order), bao gồm các chức năng tạo, lấy, cập nhật và xóa đối tượng. Mỗi use case đều tuân theo nguyên tắc single responsibility.
Entity:
Các entity như User, Product, Order được định nghĩa để phản ánh các bảng trong cơ sở dữ liệu, sử dụng GORM.
Database:
GORM: Sử dụng GORM để kết nối và tương tác với cơ sở dữ liệu PostgreSQL. Có lớp kết nối cơ sở dữ liệu và các phương thức quản lý phiên bản (migrations).
Kết nối với PostgreSQL thông qua lớp database adapter.
API Handlers:
HTTP Handlers xử lý các yêu cầu từ client và thực hiện các tác vụ tương ứng, như tạo mới, cập nhật, tìm kiếm và xóa đối tượng User, Product, Order. Sử dụng framework Gin.
HTTP Server:
Gin Framework: Sử dụng Gin để khởi tạo server HTTP, áp dụng các middleware và định nghĩa các route API cho các phiên bản (v1, v2 nếu cần).
Database Migrations:
Hệ thống quản lý phiên bản cho database thông qua các file SQL migration được lưu trữ trong thư mục migrations/.
Cấu hình (Configs):
Cấu hình chung và môi trường (development, staging, production) được lưu trữ trong các file cấu hình YAML, bao gồm cấu hình kết nối cơ sở dữ liệu, port server và các thiết lập môi trường khác.
Wire (Dependency Injection):
Sử dụng Wire để quản lý và tự động inject các thành phần như repository, service, và handler.
Chi tiết các tính năng:
1. Middleware:
a. Xác thực (auth)
Mục đích: Xác thực người dùng thông qua Bearer token.
Chức năng: Middleware kiểm tra token từ header Authorization, nếu không có hoặc không hợp lệ sẽ trả về lỗi 401 Unauthorized.
b. Kiểm tra quyền truy cập (authorization)
Mục đích: Xác minh vai trò của người dùng.
Chức năng: Middleware kiểm tra vai trò của người dùng và xác thực quyền truy cập vào tài nguyên. Nếu quyền hạn không đủ, trả về lỗi 403 Forbidden.
c. Xử lý phản hồi API (response_handler)
Mục đích: Chuẩn hóa định dạng phản hồi từ API.
Chức năng: Mọi phản hồi đều được xử lý và trả về theo định dạng APIResponse, bao gồm HttpStatus, Errors, và Data nếu thành công.
d. Xử lý lỗi API (errorHandler)
Mục đích: Quản lý lỗi phát sinh.
Chức năng: Mọi lỗi trong quá trình xử lý request được trả về theo định dạng chuẩn APIResponse với mã lỗi HTTP và thông báo lỗi chi tiết.
e. Kiểm tra dữ liệu đầu vào (validate)
Mục đích: Kiểm tra và xác thực dữ liệu đầu vào từ client.
Chức năng: Middleware kiểm tra các trường dữ liệu bắt buộc và định dạng, trả về lỗi 400 Bad Request nếu không hợp lệ.
f. Ghi log (logger)
Mục đích: Ghi lại thông tin về các yêu cầu HTTP.
Chức năng: Middleware ghi log các thông tin quan trọng như phương thức HTTP, đường dẫn, trạng thái HTTP, và thời gian xử lý.
2. API Response:
a. APIResponse
Mục đích: Chuẩn hóa phản hồi API.
Chức năng:
Success: Trả về phản hồi thành công với dữ liệu.
Error: Trả về lỗi hệ thống hoặc lỗi xác thực quyền truy cập.
ValidationError: Trả về lỗi kiểm tra dữ liệu đầu vào.
3. Repository Pattern:
a. GenericBaseRepository
Mục đích: Cung cấp các chức năng CRUD chung.
Chức năng:
Create: Tạo mới một bản ghi.
Update: Cập nhật bản ghi.
Delete: Xóa bản ghi theo ID.
FindAll: Lấy tất cả bản ghi.
FindByID: Tìm bản ghi theo ID.
b. Repository cụ thể
Mục đích: Cung cấp các chức năng CRUD cho từng entity (User, Product, Order).
Chức năng: Kế thừa từ GenericBaseRepository và mở rộng nếu cần.
4. Use Case:
a. Business Logic
Mục đích: Xử lý các nghiệp vụ liên quan đến các đối tượng User, Product, Order.
Chức năng:
CreateUser, GetUserByID, UpdateUser, DeleteUser: Các thao tác CRUD với User.
Tương tự cho Product và Order.
5. Entity:
a. Entity
Mục đích: Phản ánh các bảng trong cơ sở dữ liệu.
Chức năng: Các entity như User, Product, Order được định nghĩa với các trường tương ứng với bảng trong database.
6. Database:
a. GORM
Mục đích: Kết nối và thao tác với cơ sở dữ liệu PostgreSQL.
Chức năng: Sử dụng GORM để quản lý các phiên làm việc với database.
b. Database Migration
Mục đích: Quản lý schema database.
Chức năng: Các file SQL trong thư mục migrations/ để tạo và xóa bảng.
7. API Handlers:
a. Handlers
Mục đích: Xử lý các yêu cầu HTTP cho các endpoint API.
Chức năng:
GetUsers, CreateUser, UpdateUser, DeleteUser: Tương tác với entity User qua HTTP.
Tương tự cho Product và Order.
8. HTTP Server:
a. Gin Framework
Mục đích: Khởi chạy server HTTP.
Chức năng: Khởi tạo server với các route API và áp dụng các middleware.
