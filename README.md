# 🛒 E-commerce Platform

## 📌 Giới thiệu
E-commerce Platform là một hệ thống thương mại điện tử được xây dựng với kiến trúc microservices, hỗ trợ thanh toán trực tuyến, xử lý đơn hàng, quản lý sản phẩm và thông báo thời gian thực. Hệ thống được thiết kế nhằm đảm bảo hiệu suất cao, khả năng mở rộng và tính linh hoạt trong việc tích hợp với các dịch vụ bên thứ ba.

## 🏗 Kiến trúc hệ thống

<img src="design-system/High%20level%20architecture.png" alt="System Design" width="800"/>

Hệ thống sử dụng kiến trúc microservices, giao tiếp qua gRPC và sự kiện Kafka. Các thành phần chính bao gồm:

### 1️⃣ **API Gateway**
- Xử lý authentication, điều hướng request từ frontend đến các service backend.
- Lưu OTP trong Redis để xác thực email, số điện thoại (TTL = 5 phút).
- Cache kiểm tra tồn kho sản phẩm để tối ưu hiệu suất.
- Giao tiếp RESTful API với frontend.

### 2️⃣ **Order Service + Payment Service**
- Xử lý đơn hàng, thanh toán.
- Tích hợp Stripe để thanh toán trực tuyến.
- Lưu thông tin đơn hàng vào **Orders DB**.

### 3️⃣ **Supplier Service + Product Service**
- Quản lý sản phẩm, kho hàng.
- Lưu dữ liệu vào **Supplier DB**.

### 4️⃣ **Notification Service**
- Gửi thông báo khi có sự kiện quan trọng (đơn hàng, khuyến mãi, v.v.).
- Sử dụng Kafka để xử lý thông báo nền.
- Lưu dữ liệu vào **Notifications DB**.

### 5️⃣ **MinIO Storage**
- Lưu trữ ảnh sản phẩm và tài nguyên khác (tương tự S3).

### 6️⃣ **Stripe Adapter & Notification Adapter**
- **Stripe Adapter**: Giao tiếp với Stripe để xử lý thanh toán.
- **Notification Adapter**: Tích hợp với các nhà cung cấp dịch vụ thông báo bên thứ ba.

## 🛠 Công nghệ sử dụng
- **Backend**: Golang, gRPC, REST API
- **Frontend**: ReactJS, JS, CharkaUI
- **Message Queue**: Kafka
- **Database**: PostgreSQL, Redis (cache & OTP)
- **Storage**: MinIO (tương thích S3)
- **Authentication**: OAuth2 / JWT
- **Deployment**: Docker, Docker Compose

## 🔧 Cài đặt & Chạy ứng dụng

### 1️⃣ Clone repository
```bash
git clone https://github.com/TienMinh25/ecommerce-platform.git
cd ecommerce-platform
```

### 2️⃣ Cấu hình môi trường
- Tạo file `.env.prod` trong thư mục `configs/` với các biến môi trường cần thiết.
- Cấu hình Postgres, Kafka, Redis, MinIO, Stripe theo yêu cầu của hệ thống.

### 3️⃣ Chạy ứng dụng bằng Docker Compose
```bash
docker-compose up -d
```

### 4️⃣ Truy cập API Gateway
- API Gateway mặc định chạy tại:
  ```
  http://localhost:3000
  ```

## 📄 Tài liệu API
Xem tài liệu API đầy đủ trong thư mục `docs`.

## 📄 Tài liệu thiết kế chi tiết cho từng service
Xem tài liệu đầy đủ tại thư mục desgin-system

## 📝 License
Dự án này được phát hành dưới giấy phép MIT.

## 🚀 Contributors
- **Tiến Minh** - *Backend Developer*
- Và các contributors khác...

---  
🔥 Chúc bạn triển khai thành công! 🚀