package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type PermissionDetail struct {
	ModuleID    int   `json:"module_id"`
	Permissions []int `json:"permissions"`
}

type progressUpdate struct {
	goroutineID int
	count       int
}

// Cấu trúc cho dữ liệu địa giới hành chính
type Province struct {
	ID        int        `json:"code"`
	Name      string     `json:"name"`
	Districts []District `json:"districts"`
}

type District struct {
	ID    int    `json:"code"`
	Name  string `json:"name"`
	Wards []Ward `json:"wards"`
}

type Ward struct {
	ID   int    `json:"code"`
	Name string `json:"name"`
}

// Danh sách các APIs hỗ trợ dữ liệu địa giới hành chính Việt Nam
var vietnamGeoAPIs = []string{
	"https://provinces.open-api.vn/api/?depth=3", // API với đầy đủ phường/xã
	"https://vietnam-administrative-divisions.vercel.app/api/",
	"https://vapi.vnappmob.com/api/province/",
}

// Mở rộng danh sách attribute cho các danh mục mới
var categoryAttributes = map[string]map[string][]string{
	"Điện tử": {
		"Màu sắc":    {"Đen", "Trắng", "Xanh", "Xám", "Bạc", "Vàng"},
		"Dung lượng": {"64GB", "128GB", "256GB", "512GB", "1TB"},
	},
	"Thời trang": {
		"Màu sắc":    {"Đen", "Trắng", "Xanh", "Đỏ", "Nâu", "Hồng"},
		"Kích thước": {"S", "M", "L", "XL", "XXL"},
	},
	"Gia dụng": {
		"Màu sắc":   {"Đen", "Trắng", "Bạc", "Xám"},
		"Công suất": {"500W", "1000W", "1500W", "2000W"},
	},
	"Sách": {
		"Loại bìa": {"Bìa mềm", "Bìa cứng"},
		"Ngôn ngữ": {"Tiếng Việt", "Tiếng Anh"},
	},
	"Thể thao": {
		"Màu sắc":    {"Đen", "Xanh", "Đỏ", "Tím", "Hồng"},
		"Kích thước": {"Nhỏ", "Vừa", "Lớn"},
	},
	"Làm đẹp": {
		"Màu sắc":   {"Đỏ", "Hồng", "Nude", "Cam"},
		"Dung tích": {"15ml", "30ml", "50ml", "100ml"},
	},
	"Thực phẩm": {
		"Trọng lượng": {"340g", "500g", "1kg", "5kg", "10kg"},
		"Xuất xứ":     {"Việt Nam", "Thái Lan", "Nhật Bản"},
	},
	"Nội thất": {
		"Màu sắc":    {"Trắng", "Nâu", "Đen", "Xám", "Be"},
		"Kích thước": {"Nhỏ", "Vừa", "Lớn"},
	},
}

func main() {
	ctx := context.Background()
	gofakeit.Seed(0)

	// Database connection strings
	apiDSN := "postgres://admin:admin@localhost:5432/api_gateway_db?sslmode=disable"
	notifDSN := "postgres://admin:admin@localhost:5432/notifications_db?sslmode=disable"
	orderDSN := "postgres://admin:admin@localhost:5432/orders_db?sslmode=disable"
	partnerDSN := "postgres://admin:admin@localhost:5432/partners_db?sslmode=disable"

	// Connect to all databases
	pools := map[string]*pgxpool.Pool{
		"api_gateway_db":   connectDB(ctx, apiDSN, "api_gateway_db"),
		"notifications_db": connectDB(ctx, notifDSN, "notifications_db"),
		"orders_db":        connectDB(ctx, orderDSN, "orders_db"),
		"partners_db":      connectDB(ctx, partnerDSN, "partners_db"),
	}
	defer func() {
		for _, pool := range pools {
			pool.Close()
		}
	}()

	log.Println("🙂‍↔️ Connected to all databases")
	log.Println("🏃‍♂️ Seeding data...")

	// Tải dữ liệu địa giới hành chính Việt Nam (nếu có thể)
	adminDivisions := loadAdministrativeDivisions()

	// Seed independent tables
	seedAPIGatewayIndependentTables(ctx, pools["api_gateway_db"])
	seedOrderIndependentTables(ctx, pools["orders_db"], adminDivisions)
	seedPartnersIndependentTables(ctx, pools["partners_db"])

	// Seed users (10,000) and get their IDs
	userIDs := seedUsers(ctx, pools["api_gateway_db"], 10000, 1000, 100)

	// Seed addresses for all users
	seedAddressesForUsers(ctx, pools["api_gateway_db"], userIDs, adminDivisions)

	// Seed dependent tables
	seedNotificationPreferences(ctx, pools["notifications_db"], userIDs)
	seedCarts(ctx, pools["orders_db"], userIDs)

	// Seed admin address, cart, and notification preferences
	seedAdminDependentData(ctx, pools)

	// Select supplier user IDs and assign supplier role
	supplierUserIDs := selectSupplierUserIDs(userIDs)
	assignSupplierRole(ctx, pools["api_gateway_db"], supplierUserIDs)

	// Select deliverer user IDs and assign deliverer role
	delivererUserIDs := selectDelivererUserIDs(userIDs)
	assignDelivererRole(ctx, pools["api_gateway_db"], delivererUserIDs)

	// Seed supplier profiles and products
	supplierIDs := seedSupplierProfiles(ctx, pools["api_gateway_db"], pools["partners_db"], supplierUserIDs)
	seedEnhancedProducts(ctx, pools["partners_db"], supplierIDs)

	// Seed deliverer profiles
	seedDelivererProfiles(ctx, pools["orders_db"], delivererUserIDs, adminDivisions)

	// Seed additional data
	seedEverything(ctx, pools, userIDs, supplierIDs, adminDivisions)

	fmt.Println("✅ Seed completed successfully")
}
func connectDB(ctx context.Context, dsn, dbName string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to %s: %v", dbName, err)
	}
	return pool
}

// Cải tiến để sử dụng file hanh-chinh-viet-nam.json đã tải về
func loadAdministrativeDivisions() []Province {
	// Thử đọc từ file hanh-chinh-viet-nam.json
	if data, err := os.ReadFile("hanh-chinh-viet-nam.json"); err == nil {
		var provinces []Province
		if err := json.Unmarshal(data, &provinces); err == nil {
			log.Println("✅ Loaded administrative divisions data from hanh-chinh-viet-nam.json")
			return provinces
		} else {
			log.Printf("Warning: Failed to parse hanh-chinh-viet-nam.json: %v", err)
		}
	}

	// Thử tải dữ liệu từ API nếu không có file hoặc parse lỗi
	for _, apiURL := range vietnamGeoAPIs {
		resp, err := http.Get(apiURL)
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				var provinces []Province
				if err := json.Unmarshal(body, &provinces); err == nil {
					log.Println("✅ Loaded administrative divisions data from API:", apiURL)

					// Lưu vào file để dùng sau này
					os.WriteFile("hanh-chinh-viet-nam.json", body, 0644)
					log.Println("✅ Saved administrative divisions data to hanh-chinh-viet-nam.json")

					return provinces
				}
			}
		}
	}

	// Fallback vào dữ liệu mẫu nếu không thể tải
	log.Println("⚠️ Using sample administrative divisions data")
	return nil
}

// API Gateway Seeding
func seedAPIGatewayIndependentTables(ctx context.Context, db *pgxpool.Pool) {
	seedRoles(ctx, db)
	seedModules(ctx, db)
	seedPermissions(ctx, db)
	seedAddressTypes(ctx, db)
	seedRolePermissions(ctx, db)
	seedAdmin(ctx, db)
}

func seedRoles(ctx context.Context, db *pgxpool.Pool) {
	roles := []struct {
		name, desc string
	}{
		{"admin", "Administrator with full access"},
		{"customer", "Regular user with basic access"},
		{"supplier", "User with access to manage products and store"},
		{"deliverer", "User with access to manage delivery statuses"},
	}
	for _, r := range roles {
		_, _ = db.Exec(ctx, `INSERT INTO roles (role_name, description) VALUES ($1, $2) ON CONFLICT DO NOTHING;`, r.name, r.desc)
	}
}

func seedModules(ctx context.Context, db *pgxpool.Pool) {
	modules := []string{
		"User Management", "Role & Permission", "Product Management", "Cart",
		"Order Management", "Payment", "Shipping Management", "Review & Rating",
		"Store Management", "Onboarding", "Address Type Management", "Module Management",
		"Coupon Management",
	}
	for _, m := range modules {
		_, _ = db.Exec(ctx, `INSERT INTO modules (name) VALUES ($1) ON CONFLICT DO NOTHING;`, m)
	}
}

func seedPermissions(ctx context.Context, db *pgxpool.Pool) {
	perms := []string{"create", "update", "delete", "read", "approve", "reject"}
	for _, p := range perms {
		_, _ = db.Exec(ctx, `INSERT INTO permissions (name) VALUES ($1) ON CONFLICT DO NOTHING;`, p)
	}
}

func seedAddressTypes(ctx context.Context, db *pgxpool.Pool) {
	types := []string{"Home", "Office", "Warehouse", "Storefront", "Other"}
	for _, t := range types {
		_, _ = db.Exec(ctx, `INSERT INTO address_types (address_type) VALUES ($1) ON CONFLICT DO NOTHING;`, t)
	}
}

func seedRolePermissions(ctx context.Context, db *pgxpool.Pool) {
	roles := make(map[string]int64)
	rows, err := db.Query(ctx, `SELECT id, role_name FROM roles`)
	if err != nil {
		log.Fatal("get roles:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("scan role:", err)
		}
		roles[name] = id
	}

	rolePermissions := map[string][]PermissionDetail{
		"admin": {
			{ModuleID: 1, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 2, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 3, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 4, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 5, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 6, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 7, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 8, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 9, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 10, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 11, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 12, Permissions: []int{1, 2, 3, 4, 5, 6}},
			{ModuleID: 13, Permissions: []int{1, 2, 3, 4, 5, 6}},
		},
		"customer": {
			{ModuleID: 1, Permissions: []int{1, 2, 3, 4}},
			{ModuleID: 4, Permissions: []int{1, 2, 3, 4}},
			{ModuleID: 5, Permissions: []int{1, 4}},
			{ModuleID: 6, Permissions: []int{1, 4, 3}},
			{ModuleID: 7, Permissions: []int{4}},
			{ModuleID: 8, Permissions: []int{1, 4, 3}},
			{ModuleID: 11, Permissions: []int{4}},
			{ModuleID: 13, Permissions: []int{4}},
		},
		"supplier": {
			{ModuleID: 3, Permissions: []int{1, 2, 3, 4}},
			{ModuleID: 9, Permissions: []int{1, 2, 3, 4}},
			{ModuleID: 5, Permissions: []int{2, 4}},
		},
		"deliverer": {
			{ModuleID: 5, Permissions: []int{2, 4}},
			{ModuleID: 7, Permissions: []int{2, 4}},
		},
	}

	for roleName, permissions := range rolePermissions {
		roleID, exists := roles[roleName]
		if !exists {
			log.Fatalf("Role not found: %s", roleName)
		}
		bytes, _ := json.Marshal(permissions)
		_, err := db.Exec(ctx, `
			INSERT INTO role_permissions (role_id, permission_detail) 
			VALUES ($1, $2) 
			ON CONFLICT (role_id) DO UPDATE 
			SET permission_detail = $2, updated_at = CURRENT_TIMESTAMP;
		`, roleID, bytes)
		if err != nil {
			log.Fatalf("Insert role permissions for %s: %v", roleName, err)
		}
	}
}

func seedAdmin(ctx context.Context, db *pgxpool.Pool) {
	var userID int64
	err := db.QueryRow(ctx, `
		INSERT INTO users (fullname, email, avatar_url, email_verified, status, phone_verified, phone, birthdate) 
		VALUES ('Admin User', 'admin@admin.com', 'https://ui-avatars.com/api/?name=Admin+User', TRUE, 'active', TRUE, '+84987654321', '1990-01-01') 
		ON CONFLICT (email) DO UPDATE 
		SET fullname = 'Admin User', 
		    avatar_url = 'https://ui-avatars.com/api/?name=Admin+User',
		    phone = '+84987654321',
		    birthdate = '1990-01-01',
		    updated_at = CURRENT_TIMESTAMP 
		RETURNING id;
	`).Scan(&userID)
	if err != nil {
		log.Fatal("insert admin user:", err)
	}

	hash, _ := utils.HashPassword("admin123")
	_, err = db.Exec(ctx, `
		INSERT INTO user_password (id, password) 
		VALUES ($1, $2) 
		ON CONFLICT (id) DO UPDATE 
		SET password = $2;
	`, userID, hash)
	if err != nil {
		log.Fatal("insert admin password:", err)
	}

	var roleID int64
	err = db.QueryRow(ctx, `SELECT id FROM roles WHERE role_name='admin'`).Scan(&roleID)
	if err != nil {
		log.Fatal("get admin role:", err)
	}

	_, err = db.Exec(ctx, `
		INSERT INTO users_roles (user_id, role_id) 
		VALUES ($1, $2) 
		ON CONFLICT (role_id, user_id) DO NOTHING;
	`, userID, roleID)
	if err != nil {
		log.Fatal("assign admin role:", err)
	}
}

// Cập nhật seedAdminDependentData để sử dụng dữ liệu từ file JSON
func seedAdminDependentData(ctx context.Context, pools map[string]*pgxpool.Pool) {
	var adminID int64
	var homeAddressTypeID int64

	// Lấy ID của admin
	err := pools["api_gateway_db"].QueryRow(ctx, `SELECT id FROM users WHERE email = 'admin@admin.com'`).Scan(&adminID)
	if err != nil {
		log.Fatal("get admin ID:", err)
	}

	// Lấy address_type_id cho Home
	err = pools["api_gateway_db"].QueryRow(ctx, `SELECT id FROM address_types WHERE address_type = 'Home'`).Scan(&homeAddressTypeID)
	if err != nil {
		log.Fatal("get home address type:", err)
	}

	// Đọc dữ liệu từ file hanh-chinh-viet-nam.json để lấy thông tin Hà Nội
	var provinces []Province
	data, err := os.ReadFile("hanh-chinh-viet-nam.json")

	// Mặc định sử dụng dữ liệu cứng nếu không đọc được file
	provinceName := "Hà Nội"
	districtName := "Hai Bà Trưng"
	wardName := "Phường Bách Khoa"

	if err == nil {
		// Parse JSON data
		if err := json.Unmarshal(data, &provinces); err == nil {
			// Tìm Hà Nội trong danh sách tỉnh/thành phố
			for _, province := range provinces {
				if province.Name == "Hà Nội" || province.Name == "Thành phố Hà Nội" {
					provinceName = province.Name

					// Tìm quận Hai Bà Trưng
					for _, district := range province.Districts {
						if district.Name == "Quận Hai Bà Trưng" {
							districtName = district.Name

							// Tìm phường Bách Khoa
							for _, ward := range district.Wards {
								if strings.Contains(ward.Name, "Bách Khoa") {
									wardName = ward.Name
									break
								}
							}
							break
						}
					}
					break
				}
			}
		}
	}

	// Insert Address cho admin với thêm cột ward
	_, err = pools["api_gateway_db"].Exec(ctx, `
        INSERT INTO addresses (
            user_id, recipient_name, phone, street, district, province, postal_code, 
            country, is_default, address_type_id, ward
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        ON CONFLICT DO NOTHING;
    `, adminID, "Admin User", "+84987654321", "Số 1 Đại Cồ Việt", districtName,
		provinceName, "100000", "Việt Nam", true, homeAddressTypeID, wardName)

	if err != nil {
		log.Printf("Warning: Insert address for admin: %v", err)
	}

	// Insert Cart cho admin
	_, err = pools["orders_db"].Exec(ctx, `
        INSERT INTO carts (user_id)
        VALUES ($1)
        ON CONFLICT DO NOTHING;
    `, adminID)

	if err != nil {
		log.Printf("Warning: Insert cart for admin: %v", err)
	}

	// Insert Notification Preferences cho admin - không có survey
	prefs := map[string]bool{
		"order_status":   true,
		"payment_status": true,
		"product_status": true,
		"promotion":      true,
	}
	prefsJSON, _ := json.Marshal(prefs)

	_, err = pools["notifications_db"].Exec(ctx, `
        INSERT INTO notification_preferences (user_id, email_preferences, in_app_preferences)
        VALUES ($1, $2, $2)
        ON CONFLICT (user_id) DO UPDATE
        SET email_preferences = $2, in_app_preferences = $2;
    `, adminID, prefsJSON)

	if err != nil {
		log.Printf("Warning: Insert notification preferences for admin: %v", err)
	}

	log.Println("✅ Admin dependent data seeded successfully")
}

func seedUsers(ctx context.Context, db *pgxpool.Pool, total, batchSize, numGoroutines int) []int64 {
	// SỬA: Thêm kiểm tra đầu vào
	if total <= 0 || batchSize <= 0 || numGoroutines <= 0 {
		log.Fatal("total, batchSize, and numGoroutines must be > 0")
	}

	var sequenceMutex sync.Mutex
	var emailSequence int64 = 0

	type userInput struct {
		fullname     string
		email        string
		avatar       string
		phone        string
		birthdate    time.Time
		passwordHash string
	}

	var customerRoleID int
	if err := db.QueryRow(ctx, `SELECT id FROM roles WHERE role_name = $1`, common.RoleCustomer).Scan(&customerRoleID); err != nil {
		log.Fatal("get customer role_id:", err)
	}

	perGoroutine := total / numGoroutines
	remainder := total % numGoroutines

	progressChan := make(chan progressUpdate, numGoroutines)
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	var allUserIDs []int64
	var userIDsMutex sync.Mutex

	fmt.Println("🚀 Starting seed with", numGoroutines, "goroutines")

	for i := 0; i < numGoroutines; i++ {
		workload := perGoroutine
		if i < remainder {
			workload++
		}

		go func(goroutineID int, workload int) {
			defer wg.Done()
			seeded := 0
			var localUserIDs []int64

			for seeded < workload {
				toSeed := batchSize
				if workload-seeded < batchSize {
					toSeed = workload - seeded
				}

				var users []userInput
				for len(users) < toSeed {
					name := gofakeit.Name()
					sequenceMutex.Lock()
					seq := emailSequence
					emailSequence++
					sequenceMutex.Unlock()

					username := strings.ToLower(strings.Replace(name, " ", ".", -1))
					email := fmt.Sprintf("%s.%d@example.com", username, seq)
					avatar := fmt.Sprintf("https://ui-avatars.com/api/?name=%s", strings.ReplaceAll(name, " ", "+"))
					birth := gofakeit.DateRange(
						time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
					)
					hash, _ := utils.HashPassword("123456")

					phone := fmt.Sprintf("+84%d", gofakeit.Number(300000000, 999999999))
					users = append(users, userInput{
						fullname:     name,
						email:        email,
						avatar:       avatar,
						phone:        phone,
						birthdate:    birth,
						passwordHash: hash,
					})
				}

				var args []interface{}
				query := `INSERT INTO users (fullname, email, avatar_url, phone, birthdate, email_verified, phone_verified, status) VALUES `
				valueStrings := make([]string, len(users))
				for i, u := range users {
					idx := i * 8
					valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7, idx+8)
					args = append(args, u.fullname, u.email, u.avatar, u.phone, u.birthdate, true, true, "active")
				}
				query += strings.Join(valueStrings, ",") + " RETURNING id;"

				// SỬA: Thêm xử lý lỗi mềm hơn
				rows, err := db.Query(ctx, query, args...)
				if err != nil {
					log.Printf("Goroutine #%d: insert users: %v", goroutineID, err)
					continue
				}

				for rows.Next() {
					var id int64
					if err := rows.Scan(&id); err != nil {
						log.Printf("Goroutine #%d: scan user id: %v", goroutineID, err)
						continue
					}
					localUserIDs = append(localUserIDs, id)
				}
				rows.Close()

				var pwArgs []interface{}
				pwValues := make([]string, len(users))
				for i, id := range localUserIDs[len(localUserIDs)-len(users):] {
					idx := i * 2
					pwValues[i] = fmt.Sprintf("($%d, $%d)", idx+1, idx+2)
					pwArgs = append(pwArgs, id, users[i].passwordHash)
				}
				pwQuery := `INSERT INTO user_password (id, password) VALUES ` + strings.Join(pwValues, ",") + ";"
				if _, err := db.Exec(ctx, pwQuery, pwArgs...); err != nil {
					log.Printf("Goroutine #%d: insert user_password: %v", goroutineID, err)
					continue
				}

				var roleArgs []interface{}
				roleValues := make([]string, len(users))
				for i, id := range localUserIDs[len(localUserIDs)-len(users):] {
					idx := i * 2
					roleValues[i] = fmt.Sprintf("($%d, $%d)", idx+1, idx+2)
					roleArgs = append(roleArgs, id, customerRoleID)
				}
				roleQuery := `INSERT INTO users_roles (user_id, role_id) VALUES ` + strings.Join(roleValues, ",") + ";"
				if _, err := db.Exec(ctx, roleQuery, roleArgs...); err != nil {
					log.Printf("Goroutine #%d: insert users_roles: %v", goroutineID, err)
					continue
				}

				seeded += toSeed
				progressChan <- progressUpdate{goroutineID: goroutineID, count: seeded}
			}

			userIDsMutex.Lock()
			allUserIDs = append(allUserIDs, localUserIDs...)
			userIDsMutex.Unlock()
		}(i, workload)
	}

	go func() {
		progress := make([]int, numGoroutines)
		totalInserted := 0
		for update := range progressChan {
			progress[update.goroutineID] = update.count
			totalInserted = 0
			for _, count := range progress {
				totalInserted += count
			}
			fmt.Printf("⏳ Progress: Goroutine #%d inserted %d records. Total: %d/%d (%.2f%%)\n",
				update.goroutineID, update.count, totalInserted, total, float64(totalInserted)*100/float64(total))
		}
	}()

	wg.Wait()
	close(progressChan)
	time.Sleep(100 * time.Millisecond)

	fmt.Println("🎉 Done seeding users with concurrent goroutines.")
	return allUserIDs
}

// Sửa lại hàm selectSupplierUserIDs để đảm bảo admin cũng là supplier
func selectSupplierUserIDs(userIDs []int64) []int64 {
	count := 15 // Cố định 15 supplier thay vì len(userIDs) / 10
	supplierUserIDs := make([]int64, 0, count+1)

	// Đảm bảo admin cũng là supplier (ID=1)
	var adminID int64 = 1 // Thông thường admin là ID đầu tiên
	supplierUserIDs = append(supplierUserIDs, adminID)

	// Trộn ngẫu nhiên để chọn users làm supplier
	gofakeit.ShuffleAnySlice(userIDs)
	addedCount := 0
	for i := 0; i < len(userIDs) && addedCount < count; i++ {
		if userIDs[i] != adminID {
			supplierUserIDs = append(supplierUserIDs, userIDs[i])
			addedCount++
		}
	}
	return supplierUserIDs
}

func assignSupplierRole(ctx context.Context, db *pgxpool.Pool, supplierUserIDs []int64) {
	var supplierRoleID int64
	err := db.QueryRow(ctx, `SELECT id FROM roles WHERE role_name = $1`, "supplier").Scan(&supplierRoleID)
	if err != nil {
		log.Fatal("get supplier role_id:", err)
	}

	batchSize := 1000
	for i := 0; i < len(supplierUserIDs); i += batchSize {
		end := i + batchSize
		if end > len(supplierUserIDs) {
			end = len(supplierUserIDs)
		}
		batch := supplierUserIDs[i:end]

		var args []interface{}
		query := `INSERT INTO users_roles (user_id, role_id) VALUES `
		valueStrings := make([]string, len(batch))
		for j, userID := range batch {
			idx := j * 2
			valueStrings[j] = fmt.Sprintf("($%d, $%d)", idx+1, idx+2)
			args = append(args, userID, supplierRoleID)
		}
		query += strings.Join(valueStrings, ",") + " ON CONFLICT DO NOTHING;"
		_, err := db.Exec(ctx, query, args...)
		if err != nil {
			log.Fatal("assign supplier role:", err)
		}
	}
}

// Tương tự với deliverer, cũng thêm admin
func selectDelivererUserIDs(userIDs []int64) []int64 {
	count := 15 // Cố định 15 deliverer thay vì len(userIDs) / 20
	delivererUserIDs := make([]int64, 0, count+1)

	// Thêm admin vào danh sách deliverer
	var adminID int64 = 1
	delivererUserIDs = append(delivererUserIDs, adminID)

	// Trộn ngẫu nhiên để chọn users làm deliverer
	gofakeit.ShuffleAnySlice(userIDs)
	addedCount := 0
	for i := 0; i < len(userIDs) && addedCount < count; i++ {
		if userIDs[i] != adminID {
			delivererUserIDs = append(delivererUserIDs, userIDs[i])
			addedCount++
		}
	}
	return delivererUserIDs
}

func assignDelivererRole(ctx context.Context, db *pgxpool.Pool, delivererUserIDs []int64) {
	var delivererRoleID int64
	err := db.QueryRow(ctx, `SELECT id FROM roles WHERE role_name = $1`, "deliverer").Scan(&delivererRoleID)
	if err != nil {
		log.Fatal("get deliverer role_id:", err)
	}

	batchSize := 1000
	for i := 0; i < len(delivererUserIDs); i += batchSize {
		end := i + batchSize
		if end > len(delivererUserIDs) {
			end = len(delivererUserIDs)
		}
		batch := delivererUserIDs[i:end]

		var args []interface{}
		query := `INSERT INTO users_roles (user_id, role_id) VALUES `
		valueStrings := make([]string, len(batch))
		for j, userID := range batch {
			idx := j * 2
			valueStrings[j] = fmt.Sprintf("($%d, $%d)", idx+1, idx+2)
			args = append(args, userID, delivererRoleID)
		}
		query += strings.Join(valueStrings, ",") + " ON CONFLICT DO NOTHING;"
		_, err := db.Exec(ctx, query, args...)
		if err != nil {
			log.Fatal("assign deliverer role:", err)
		}
	}
}

// Cải tiến seedAddressesForUsers để sử dụng dữ liệu địa giới hành chính từ file JSON và xử lý mảng rỗng
func seedAddressesForUsers(ctx context.Context, db *pgxpool.Pool, userIDs []int64, adminDivisions []Province) {
	var homeAddressTypeID int64
	err := db.QueryRow(ctx, `SELECT id FROM address_types WHERE address_type = $1`, "Home").Scan(&homeAddressTypeID)
	if err != nil {
		log.Fatal("get home address type:", err)
	}

	// Đọc dữ liệu từ file hanh-chinh-viet-nam.json
	var provinces []Province
	data, err := os.ReadFile("hanh-chinh-viet-nam.json")
	if err != nil {
		log.Printf("Error reading hanh-chinh-viet-nam.json: %v", err)
		// Fallback to adminDivisions if file reading fails
		provinces = adminDivisions
	} else {
		// Parse JSON data
		if err := json.Unmarshal(data, &provinces); err != nil {
			log.Printf("Error parsing hanh-chinh-viet-nam.json: %v", err)
			// Fallback to adminDivisions if JSON parsing fails
			provinces = adminDivisions
		}
	}

	if len(provinces) == 0 {
		log.Println("⚠️ No administrative divisions data available, cannot seed addresses properly")
		return
	}

	log.Printf("✅ Loaded %d provinces from hanh-chinh-viet-nam.json for address seeding", len(provinces))

	// Lọc ra các tỉnh có ít nhất một quận/huyện, và quận/huyện có ít nhất một phường/xã
	var validProvinces []Province
	for _, province := range provinces {
		if len(province.Districts) == 0 {
			continue
		}

		var validDistricts []District
		for _, district := range province.Districts {
			if len(district.Wards) == 0 {
				continue
			}
			validDistricts = append(validDistricts, district)
		}

		if len(validDistricts) > 0 {
			province.Districts = validDistricts
			validProvinces = append(validProvinces, province)
		}
	}

	if len(validProvinces) == 0 {
		log.Println("⚠️ No valid administrative divisions data available, cannot seed addresses properly")
		return
	}

	log.Printf("✅ Found %d valid provinces with districts and wards", len(validProvinces))

	// Seed addresses in batches
	batchSize := 1000
	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}
		batch := userIDs[i:end]

		var args []interface{}
		query := `INSERT INTO addresses (user_id, recipient_name, phone, street, district, province, postal_code, country, is_default, address_type_id, ward) VALUES `
		valueStrings := make([]string, 0, len(batch))
		valsCount := 0

		for _, userID := range batch {
			// Lấy ngẫu nhiên một tỉnh/thành phố
			provinceIdx := gofakeit.Number(0, len(validProvinces)-1)
			province := validProvinces[provinceIdx]

			// Lấy ngẫu nhiên một quận/huyện từ tỉnh/thành phố đó
			districtIdx := gofakeit.Number(0, len(province.Districts)-1)
			district := province.Districts[districtIdx]

			// Lấy ngẫu nhiên một phường/xã từ quận/huyện đó
			wardIdx := gofakeit.Number(0, len(district.Wards)-1)
			ward := district.Wards[wardIdx]

			recipientName := gofakeit.Name()
			phone := fmt.Sprintf("+84%d", gofakeit.Number(300000000, 999999999))
			street := fmt.Sprintf("Số %d Đường %s", gofakeit.Number(1, 999), gofakeit.Street())
			postalCode := fmt.Sprintf("%06d", gofakeit.Number(100000, 999999))
			country := "Việt Nam"
			isDefault := true

			idx := valsCount * 11 // 11 parameters including ward
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7, idx+8, idx+9, idx+10, idx+11))

			args = append(args, userID, recipientName, phone, street, district.Name, province.Name, postalCode, country, isDefault, homeAddressTypeID, ward.Name)
			valsCount++
		}

		// Bỏ qua nếu không có địa chỉ hợp lệ để chèn
		if len(valueStrings) == 0 {
			continue
		}

		query += strings.Join(valueStrings, ",") + " ON CONFLICT DO NOTHING;"
		_, err := db.Exec(ctx, query, args...)
		if err != nil {
			log.Printf("Error inserting addresses: %v", err)
		}
	}

	log.Println("✅ Successfully seeded addresses for users using Vietnam administrative divisions")
}

func seedAreasFromAdminDivisions(ctx context.Context, db *pgxpool.Pool, adminDivisions []Province) {
	if len(adminDivisions) == 0 {
		log.Println("⚠️ No administrative divisions data")
		return
	}

	log.Printf("🏠 Starting to seed areas from %d provinces...", len(adminDivisions))

	// Khởi tạo squirrel query builder
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Tạo bulk insert query
	insertQuery := psql.Insert("areas").
		Columns("city", "country", "district", "ward", "area_code")

	totalAreas := 0

	// Duyệt qua tất cả và add vào query
	for _, province := range adminDivisions {
		for _, district := range province.Districts {
			for _, ward := range district.Wards {
				areaCode := fmt.Sprintf("area-code-%v-%v-%v", province.ID, district.ID, ward.ID)

				insertQuery = insertQuery.Values(
					province.Name, // city
					"Việt Nam",    // country
					district.Name, // district
					ward.Name,     // ward
					areaCode,      // area_code
				)
				totalAreas++
			}
		}
	}

	// Add ON CONFLICT DO NOTHING
	insertQuery = insertQuery.Suffix("ON CONFLICT (area_code) DO NOTHING")

	// Build query
	sql, args, err := insertQuery.ToSql()
	if err != nil {
		log.Printf("❌ Error building query: %v", err)
		return
	}

	log.Printf("📝 Executing bulk insert for %d areas...", totalAreas)

	// Execute query
	result, err := db.Exec(ctx, sql, args...)
	if err != nil {
		log.Printf("❌ Error executing bulk insert: %v", err)
		return
	}

	rowsAffected := result.RowsAffected()
	log.Printf("✅ Areas seeded successfully: %d areas inserted", rowsAffected)
}

// Cập nhật hàm seedOrderIndependentTables
func seedOrderIndependentTables(ctx context.Context, db *pgxpool.Pool, adminDivisions []Province) {
	log.Println("🏗️ Seeding Order service independent tables...")

	// Seed TẤT CẢ areas từ dữ liệu hành chính
	seedAreasFromAdminDivisions(ctx, db, adminDivisions)

	// Seed payment methods
	seedPaymentMethods(ctx, db)

	log.Println("✅ Order service independent tables seeded successfully")
}

func seedPaymentMethods(ctx context.Context, db *pgxpool.Pool) {
	methods := []struct {
		name, code string
	}{
		{"Thanh toán khi nhận hàng (COD)", "cod"},
		{"Thanh toán qua MoMo", "momo"},
	}

	for _, method := range methods {
		_, err := db.Exec(ctx, `
			INSERT INTO payment_methods (name, code, is_active)
			VALUES ($1, $2, TRUE)
			ON CONFLICT (code) DO UPDATE
			SET name = $1, is_active = TRUE;
		`, method.name, method.code)

		if err != nil {
			log.Fatalf("Error inserting payment method: %v", err)
		}
	}
	log.Println("✅ Payment methods seeded successfully")
}

func seedNotificationPreferences(ctx context.Context, db *pgxpool.Pool, userIDs []int64) {
	batchSize := 1000
	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}
		batch := userIDs[i:end]

		// Tạo dữ liệu notification preferences
		var args []interface{}
		query := `INSERT INTO notification_preferences (user_id, email_preferences, in_app_preferences) VALUES `
		valueStrings := make([]string, len(batch))

		for j, userID := range batch {
			idx := j * 3

			// Random preferences
			emailPrefs := map[string]bool{
				"survey":         gofakeit.Bool(),
				"promotion":      gofakeit.Bool(),
				"order_status":   true, // Luôn bật thông báo đơn hàng
				"payment_status": true, // Luôn bật thông báo thanh toán
				"product_status": gofakeit.Bool(),
			}
			emailPrefsJSON, _ := json.Marshal(emailPrefs)

			// In-app thường cũng giống email preferences
			inAppPrefs := make(map[string]bool)
			for k, v := range emailPrefs {
				inAppPrefs[k] = v
				// Đôi khi in-app được bật nhưng email thì không
				if !v && gofakeit.Bool() {
					inAppPrefs[k] = true
				}
			}
			inAppPrefsJSON, _ := json.Marshal(inAppPrefs)

			valueStrings[j] = fmt.Sprintf("($%d, $%d, $%d)", idx+1, idx+2, idx+3)
			args = append(args, userID, emailPrefsJSON, inAppPrefsJSON)
		}

		query += strings.Join(valueStrings, ",") + " ON CONFLICT (user_id) DO NOTHING;"
		_, err := db.Exec(ctx, query, args...)
		if err != nil {
			log.Printf("Error inserting notification preferences: %v", err)
		}
	}
	log.Println("✅ Notification preferences seeded successfully")
}

func seedCarts(ctx context.Context, db *pgxpool.Pool, userIDs []int64) {
	batchSize := 1000
	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}
		batch := userIDs[i:end]

		// Process each user individually to avoid ON CONFLICT issues
		for _, userID := range batch {
			// First check if the user already has a cart
			var count int
			err := db.QueryRow(ctx, `
                SELECT COUNT(*) FROM carts WHERE user_id = $1
            `, userID).Scan(&count)

			if err != nil {
				log.Printf("Error checking cart existence: %v", err)
				continue
			}

			// Only insert if the user doesn't have a cart
			if count == 0 {
				_, err := db.Exec(ctx, `
                    INSERT INTO carts (user_id) VALUES ($1)
                `, userID)

				if err != nil {
					log.Printf("Error inserting cart for user %d: %v", userID, err)
				}
			}
		}
	}
	log.Println("✅ Carts seeded successfully")
}

// Partners Service Seeding
func seedPartnersIndependentTables(ctx context.Context, db *pgxpool.Pool) {
	seedCategories(ctx, db)
	seedTags(ctx, db)
	seedAttributeDefinitions(ctx, db)
}

// Cập nhật hàm seedCategories để thêm nhiều danh mục cha và con hơn
func seedCategories(ctx context.Context, db *pgxpool.Pool) {
	// Chỉ tạo 8 categories chính, bỏ hết category con
	mainCategories := []struct{ name, desc, imageUrl string }{
		{
			"Điện tử",
			"Điện thoại, laptop, máy tính bảng và thiết bị điện tử",
			"https://images.unsplash.com/photo-1498049794561-7780e7231661?w=600&h=400&fit=crop",
		},
		{
			"Thời trang",
			"Quần áo, giày dép và phụ kiện thời trang",
			"https://images.unsplash.com/photo-1445205170230-053b83016050?w=600&h=400&fit=crop",
		},
		{
			"Gia dụng",
			"Đồ gia dụng và vật dụng sinh hoạt hàng ngày",
			"https://images.unsplash.com/photo-1484101403633-562f891dc89a?w=600&h=400&fit=crop",
		},
		{
			"Sách",
			"Sách và văn phòng phẩm",
			"https://images.unsplash.com/photo-1495446815901-a7297e633e8d?w=600&h=400&fit=crop",
		},
		{
			"Thể thao",
			"Dụng cụ thể thao và đồ dùng tập luyện",
			"https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?w=600&h=400&fit=crop",
		},
		{
			"Làm đẹp",
			"Mỹ phẩm và sản phẩm chăm sóc cá nhân",
			"https://images.unsplash.com/photo-1522335789203-aabd1fc54bc9?w=600&h=400&fit=crop",
		},
		{
			"Thực phẩm",
			"Thực phẩm và đồ uống",
			"https://images.unsplash.com/photo-1542838132-92c53300491e?w=600&h=400&fit=crop",
		},
		{
			"Nội thất",
			"Đồ nội thất và trang trí nhà cửa",
			"https://images.unsplash.com/photo-1586023492125-27b2c045efd7?w=600&h=400&fit=crop",
		},
	}

	for _, cat := range mainCategories {
		// First check if the category already exists
		var count int
		err := db.QueryRow(ctx, `
			SELECT COUNT(*) FROM categories WHERE name = $1
		`, cat.name).Scan(&count)

		if err != nil {
			log.Printf("Error checking category existence: %v", err)
			continue
		}

		if count == 0 {
			// Insert new category
			_, err := db.Exec(ctx, `
				INSERT INTO categories (name, description, image_url, is_active)
				VALUES ($1, $2, $3, TRUE)
			`, cat.name, cat.desc, cat.imageUrl)

			if err != nil {
				log.Printf("Error inserting category: %v", err)
			}
		} else {
			// Update existing category
			_, err := db.Exec(ctx, `
				UPDATE categories 
				SET description = $2, image_url = $3, is_active = TRUE, updated_at = CURRENT_TIMESTAMP
				WHERE name = $1
			`, cat.name, cat.desc, cat.imageUrl)

			if err != nil {
				log.Printf("Error updating category: %v", err)
			}
		}
	}
	log.Println("✅ Categories seeded successfully")
}

// Danh sách sản phẩm cụ thể với ảnh đã test và hoạt động
var specificProducts = []struct {
	name          string
	description   string
	category      string
	imageURL      string
	variantImages map[string]string
}{
	// ĐIỆN TỬ
	{
		name:        "iPhone 15 Pro Max 256GB",
		description: "Điện thoại thông minh cao cấp với chip A17 Pro, camera 48MP và màn hình Super Retina XDR 6.7 inch",
		category:    "Điện tử",
		imageURL:    "https://cdn.tgdd.vn/Products/Images/42/299033/iphone-15-pro-max-blue-thumbnew-600x600.jpg",
		variantImages: map[string]string{
			"Xanh":  "https://cdn.tgdd.vn/Products/Images/42/299033/iphone-15-pro-max-blue-thumbnew-600x600.jpg",
			"Đen":   "https://cdn.tgdd.vn/Products/Images/42/299033/iphone-15-pro-max-black-thumbnew-600x600.jpg",
			"Trắng": "https://cdn.tgdd.vn/Products/Images/42/299033/iphone-15-pro-max-white-thumbnew-600x600.jpg",
		},
	},
	{
		name:        "Samsung Galaxy S24 Ultra",
		description: "Flagship Android với bút S Pen, camera 200MP và màn hình Dynamic AMOLED 2X",
		category:    "Điện tử",
		imageURL:    "https://images.fpt.shop/unsafe/filters:quality(90)/fptshop.com.vn/Uploads/images/2015/Tin-Tuc/QuanLNH2/samsung-galaxy-s24-ultra-1.jpg",
		variantImages: map[string]string{
			"Xám": "https://images.fpt.shop/unsafe/filters:quality(90)/fptshop.com.vn/Uploads/images/2015/Tin-Tuc/QuanLNH2/samsung-galaxy-s24-ultra-1.jpg",
			"Đen": "https://images.fpt.shop/unsafe/filters:quality(90)/fptshop.com.vn/Uploads/images/2015/Tin-Tuc/QuanLNH2/samsung-galaxy-s24-ultra-2.jpg",
		},
	},
	{
		name:        "MacBook Air M3 13 inch",
		description: "Laptop siêu mỏng với chip Apple M3 và thời lượng pin 18 giờ",
		category:    "Điện tử",
		imageURL:    "https://cdn.tgdd.vn/Products/Images/44/322096/macbook-air-13-inch-m3-2024-starlight-thumb-600x600.jpg",
		variantImages: map[string]string{
			"Vàng": "https://cdn.tgdd.vn/Products/Images/44/322096/macbook-air-13-inch-m3-2024-starlight-thumb-600x600.jpg",
			"Bạc":  "https://cdn.tgdd.vn/Products/Images/44/322096/macbook-air-13-inch-m3-2024-silver-thumb-600x600.jpg",
			"Xám":  "https://cdn.tgdd.vn/Products/Images/44/322096/macbook-air-13-inch-m3-2024-space-gray-thumb-600x600.jpg",
		},
	},
	{
		name:        "iPad Pro M4 11 inch",
		description: "Máy tính bảng cao cấp với chip M4 và màn hình Liquid Retina XDR",
		category:    "Điện tử",
		imageURL:    "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/ipad-pro-11-select-wifi-spacegray-202405?wid=470&hei=556&fmt=png-alpha&.v=1713308272877",
		variantImages: map[string]string{
			"Xám": "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/ipad-pro-11-select-wifi-spacegray-202405?wid=470&hei=556&fmt=png-alpha&.v=1713308272877",
			"Bạc": "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/ipad-pro-11-select-wifi-silver-202405?wid=470&hei=556&fmt=png-alpha&.v=1713308272877",
		},
	},
	{
		name:        "AirPods Pro Gen 2",
		description: "Tai nghe true wireless với chống ồn chủ động và chip H2",
		category:    "Điện tử",
		imageURL:    "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/MQD83?wid=572&hei=572&fmt=jpeg&qlt=95&.v=1660803972361",
	},

	// THỜI TRANG
	{
		name:        "Áo thun nam Uniqlo cổ tròn",
		description: "Áo thun nam chất liệu cotton 100% mềm mại và thoáng mát",
		category:    "Thời trang",
		imageURL:    "https://image.uniqlo.com/UQ/ST3/vn/imagesgoods/422992/item/vngoods_09_422992.jpg",
		variantImages: map[string]string{
			"Đen":   "https://image.uniqlo.com/UQ/ST3/vn/imagesgoods/422992/item/vngoods_09_422992.jpg",
			"Trắng": "https://image.uniqlo.com/UQ/ST3/vn/imagesgoods/422992/item/vngoods_01_422992.jpg",
			"Xanh":  "https://image.uniqlo.com/UQ/ST3/vn/imagesgoods/422992/item/vngoods_65_422992.jpg",
		},
	},
	{
		name:        "Giày thể thao Nike Air Force 1",
		description: "Giày thể thao kinh điển với đế Air cushion và da thật cao cấp",
		category:    "Thời trang",
		imageURL:    "https://static.nike.com/a/images/c_limit,w_592,f_auto/t_product_v1/4f37fca8-6bce-43e7-ad07-f57ae3c13142/air-force-1-07-shoes-WrLlWX.png",
		variantImages: map[string]string{
			"Trắng": "https://static.nike.com/a/images/c_limit,w_592,f_auto/t_product_v1/4f37fca8-6bce-43e7-ad07-f57ae3c13142/air-force-1-07-shoes-WrLlWX.png",
			"Đen":   "https://static.nike.com/a/images/c_limit,w_592,f_auto/t_product_v1/00375837-849f-4f17-ba24-d201d27be36b/air-force-1-07-shoes-0XGfD7.png",
		},
	},
	{
		name:        "Túi xách tay nữ Coach",
		description: "Túi xách tay nữ cao cấp da thật 100% thiết kế sang trọng",
		category:    "Thời trang",
		imageURL:    "https://vietnam.coach.com/dw/image/v2/BFCR_PRD/on/demandware.static/-/Sites-coach-master-catalog/default/dw5c5e5e5e/images/large/C0772_B4NQ4_d0.jpg",
		variantImages: map[string]string{
			"Nâu": "https://vietnam.coach.com/dw/image/v2/BFCR_PRD/on/demandware.static/-/Sites-coach-master-catalog/default/dw5c5e5e5e/images/large/C0772_B4NQ4_d0.jpg",
			"Đen": "https://vietnam.coach.com/dw/image/v2/BFCR_PRD/on/demandware.static/-/Sites-coach-master-catalog/default/dw5c5e5e5e/images/large/C0772_B4BK_d0.jpg",
		},
	},

	// GIA DỤNG
	{
		name:        "Nồi cơm điện Panasonic 1.8L",
		description: "Nồi cơm điện thông minh với công nghệ IH và lòng nồi chống dính",
		category:    "Gia dụng",
		imageURL:    "https://panasonic.com.vn/wp-content/uploads/2020/09/SR-KS181WRA_1-600x600.jpg",
		variantImages: map[string]string{
			"Trắng": "https://panasonic.com.vn/wp-content/uploads/2020/09/SR-KS181WRA_1-600x600.jpg",
		},
	},
	{
		name:        "Máy lọc không khí Xiaomi Mi Air Purifier 4",
		description: "Máy lọc không khí thông minh với bộ lọc HEPA H13",
		category:    "Gia dụng",
		imageURL:    "https://cdn.tgdd.vn/Products/Images/1348/280714/may-loc-khong-khi-xiaomi-mi-air-purifier-4-thumb-600x600.jpg",
	},
	{
		name:        "Nồi chiên không dầu Philips 4.1L",
		description: "Nồi chiên không dầu công nghệ Rapid Air với dung tích 4.1L",
		category:    "Gia dụng",
		imageURL:    "https://cdn.tgdd.vn/Products/Images/1982/78874/noi-chien-khong-dau-philips-hd9200-90-4-1-lit-thumb-600x600.jpg",
		variantImages: map[string]string{
			"Đen":   "https://cdn.tgdd.vn/Products/Images/1982/78874/noi-chien-khong-dau-philips-hd9200-90-4-1-lit-thumb-600x600.jpg",
			"Trắng": "https://cdn.tgdd.vn/Products/Images/1982/78875/noi-chien-khong-dau-philips-hd9252-90-4-1-lit-thumb-600x600.jpg",
		},
	},

	// SÁCH
	{
		name:        "Đắc Nhân Tâm - Dale Carnegie",
		description: "Cuốn sách kinh điển về nghệ thuật giao tiếp và ứng xử",
		category:    "Sách",
		imageURL:    "https://salt.tikicdn.com/cache/750x750/ts/product/5e/18/24/2a6154ba08df6ce6161c13f4303fa19e.jpg.webp",
	},
	{
		name:        "Nhà Giả Kim - Paulo Coelho",
		description: "Tiểu thuyết nổi tiếng thế giới về hành trình tìm kiếm ý nghĩa cuộc sống",
		category:    "Sách",
		imageURL:    "https://salt.tikicdn.com/cache/750x750/ts/product/45/3b/fc/aa81d0592c4d5be8ad83ad1555164abc.jpg.webp",
	},
	{
		name:        "Sapiens - Yuval Noah Harari",
		description: "Lược sử loài người - cuốn sách về lịch sử và tương lai nhân loại",
		category:    "Sách",
		imageURL:    "https://salt.tikicdn.com/cache/750x750/ts/product/e0/65/24/8cf7d2f6a50b5fb60c53da09bc2db7a4.jpg.webp",
	},

	// THỂ THAO
	{
		name:        "Tạ tay điều chỉnh 20kg",
		description: "Tạ tay thông minh có thể điều chỉnh từ 2-20kg tiết kiệm không gian",
		category:    "Thể thao",
		imageURL:    "https://cdn.tgdd.vn/Products/Images/4221/236951/ta-tay-dieu-chinh-bowflex-selecttech-552-tu-2-24kg-thumb-600x600.jpg",
	},
	{
		name:        "Thảm tập yoga chống trượt 6mm",
		description: "Thảm tập yoga cao cấp dày 6mm chống trượt và thân thiện môi trường",
		category:    "Thể thao",
		imageURL:    "https://cdn.tgdd.vn/Products/Images/4221/235951/tham-tap-yoga-6mm-xanh-la-thumb-600x600.jpg",
		variantImages: map[string]string{
			"Xanh": "https://cdn.tgdd.vn/Products/Images/4221/235951/tham-tap-yoga-6mm-xanh-la-thumb-600x600.jpg",
			"Tím":  "https://cdn.tgdd.vn/Products/Images/4221/235952/tham-tap-yoga-6mm-tim-thumb-600x600.jpg",
			"Hồng": "https://cdn.tgdd.vn/Products/Images/4221/235953/tham-tap-yoga-6mm-hong-thumb-600x600.jpg",
		},
	},
	{
		name:        "Bóng đá FIFA Quality Pro",
		description: "Bóng đá chuẩn FIFA Quality Pro cho thi đấu chuyên nghiệp",
		category:    "Thể thao",
		imageURL:    "https://cdn.tgdd.vn/Products/Images/4221/236845/bong-da-fifa-quality-pro-size-5-thumb-600x600.jpg",
	},

	// LÀM ĐẸP
	{
		name:        "Son môi YSL Rouge Pur Couture",
		description: "Son môi cao cấp với công thức dưỡng ẩm và màu sắc lâu trôi",
		category:    "Làm đẹp",
		imageURL:    "https://www.yslbeauty.com.vn/dw/image/v2/AANG_PRD/on/demandware.static/-/Sites-ysl-master-catalog/default/dwc5c5c5c5/images/LIPS/ROUGE_PUR_COUTURE/3365440787984_rouge_pur_couture_1_rouge_rock.jpg",
		variantImages: map[string]string{
			"Đỏ":   "https://www.yslbeauty.com.vn/dw/image/v2/AANG_PRD/on/demandware.static/-/Sites-ysl-master-catalog/default/dwc5c5c5c5/images/LIPS/ROUGE_PUR_COUTURE/3365440787984_rouge_pur_couture_1_rouge_rock.jpg",
			"Hồng": "https://www.yslbeauty.com.vn/dw/image/v2/AANG_PRD/on/demandware.static/-/Sites-ysl-master-catalog/default/dwc5c5c5c5/images/LIPS/ROUGE_PUR_COUTURE/3365440787984_rouge_pur_couture_52_rosy_coral.jpg",
		},
	},
	{
		name:        "Kem chống nắng La Roche-Posay SPF 60",
		description: "Kem chống nắng dành cho da nhạy cảm SPF 60 PA++++",
		category:    "Làm đẹp",
		imageURL:    "https://www.laroche-posay.vn/-/media/project/loreal/brand-sites/lrp/apac/vn/products/anthelios/anthelios-airlicium-ultra-light-spf60/3337875546298.jpg",
		variantImages: map[string]string{
			"50ml":  "https://www.laroche-posay.vn/-/media/project/loreal/brand-sites/lrp/apac/vn/products/anthelios/anthelios-airlicium-ultra-light-spf60/3337875546298.jpg",
			"100ml": "https://www.laroche-posay.vn/-/media/project/loreal/brand-sites/lrp/apac/vn/products/anthelios/anthelios-airlicium-ultra-light-spf60-100ml/3337875546299.jpg",
		},
	},

	// THỰC PHẨM
	{
		name:        "Gạo ST25 Đồng Tháp",
		description: "Gạo thơm ST25 chất lượng cao từ Đồng Tháp hạt dẻo mềm thơm ngon",
		category:    "Thực phẩm",
		imageURL:    "https://cdn.tgdd.vn/Products/Images/2513/238242/gao-st25-dong-thap-tui-5kg-202103091133085068.jpg",
		variantImages: map[string]string{
			"5kg":  "https://cdn.tgdd.vn/Products/Images/2513/238242/gao-st25-dong-thap-tui-5kg-202103091133085068.jpg",
			"10kg": "https://cdn.tgdd.vn/Products/Images/2513/238243/gao-st25-dong-thap-tui-10kg-202103091133085068.jpg",
		},
	},
	{
		name:        "Cà phê Trung Nguyên Legend 5",
		description: "Cà phê pha phin truyền thống pha chế từ 100% cà phê Robusta và Arabica",
		category:    "Thực phẩm",
		imageURL:    "https://www.trungnguyenlegend.com/wp-content/uploads/2018/08/ca-phe-trung-nguyen-legend-5-500g.jpg",
		variantImages: map[string]string{
			"500g": "https://www.trungnguyenlegend.com/wp-content/uploads/2018/08/ca-phe-trung-nguyen-legend-5-500g.jpg",
			"340g": "https://www.trungnguyenlegend.com/wp-content/uploads/2018/08/ca-phe-trung-nguyen-legend-5-340g.jpg",
		},
	},

	// NỘI THẤT
	{
		name:        "Sofa 3 chỗ IKEA Ektorp",
		description: "Sofa 3 chỗ ngồi bọc vải khung gỗ thông thiết kế Scandinavian",
		category:    "Nội thất",
		imageURL:    "https://www.ikea.com/vn/en/images/products/ektorp-3-seat-sofa-vittaryd-white__0818598_pe774498_s5.jpg",
		variantImages: map[string]string{
			"Trắng": "https://www.ikea.com/vn/en/images/products/ektorp-3-seat-sofa-vittaryd-white__0818598_pe774498_s5.jpg",
			"Xám":   "https://www.ikea.com/vn/en/images/products/ektorp-3-seat-sofa-vittaryd-grey__0818599_pe774499_s5.jpg",
		},
	},
	{
		name:        "Bàn làm việc IKEA Hemnes",
		description: "Bàn làm việc gỗ thông tự nhiên với 2 ngăn kéo thiết kế cổ điển",
		category:    "Nội thất",
		imageURL:    "https://www.ikea.com/vn/en/images/products/hemnes-desk-white-stain__0318434_pe513726_s5.jpg",
		variantImages: map[string]string{
			"Trắng": "https://www.ikea.com/vn/en/images/products/hemnes-desk-white-stain__0318434_pe513726_s5.jpg",
			"Nâu":   "https://www.ikea.com/vn/en/images/products/hemnes-desk-brown__0318435_pe513727_s5.jpg",
		},
	},
}

func seedTags(ctx context.Context, db *pgxpool.Pool) {
	tags := []string{
		"Mới nhất", "Bán chạy", "Giảm giá", "Cao cấp", "Giá rẻ",
		"Chính hãng", "Chất lượng cao", "Hàng hiệu", "Thương hiệu", "Nhập khẩu",
		"Xu hướng", "Thịnh hành", "Ưu đãi", "Miễn phí vận chuyển", "Khuyến mãi",
		"Phân phối chính thức", "Hàng độc quyền", "Phiên bản giới hạn",
		"Bộ sản phẩm", "Combo", // Thêm hai tag mới cho sản phẩm cha
	}

	for _, tag := range tags {
		// Kiểm tra xem tag đã tồn tại chưa
		var count int
		err := db.QueryRow(ctx, `
			SELECT COUNT(*) FROM tags WHERE name = $1
		`, tag).Scan(&count)

		if err != nil {
			log.Printf("Error checking tag existence: %v", err)
			continue
		}

		if count == 0 {
			_, err := db.Exec(ctx, `
				INSERT INTO tags (name)
				VALUES ($1);
			`, tag)

			if err != nil {
				log.Printf("Error inserting tag: %v", err)
			}
		}
	}
	log.Println("✅ Tags seeded successfully")
}

func getAttributeOptions(attrName string) ([]string, bool) {
	// Tổng hợp tất cả các giá trị có thể có cho mỗi attribute từ các danh mục
	allOptions := make(map[string]map[string]bool)

	for _, categoryAttrs := range categoryAttributes {
		for attr, options := range categoryAttrs {
			if _, exists := allOptions[attr]; !exists {
				allOptions[attr] = make(map[string]bool)
			}

			for _, option := range options {
				allOptions[attr][option] = true
			}
		}
	}

	// Chuyển đổi map thành slice
	if optionMap, exists := allOptions[attrName]; exists {
		uniqueOptions := make([]string, 0, len(optionMap))
		for option := range optionMap {
			uniqueOptions = append(uniqueOptions, option)
		}
		return uniqueOptions, true
	}

	return nil, false
}

func seedSupplierProfiles(ctx context.Context, apiDb, partnerDb *pgxpool.Pool, supplierUserIDs []int64) []int64 {
	// Lấy danh sách địa chỉ của supplier để dùng làm địa chỉ doanh nghiệp
	supplierAddresses := make(map[int64]int64)
	for _, userID := range supplierUserIDs {
		var addressID int64
		err := apiDb.QueryRow(ctx, `
			SELECT id FROM addresses WHERE user_id = $1 LIMIT 1
		`, userID).Scan(&addressID)

		if err != nil {
			log.Printf("Warning: No address found for user_id: %d", userID)
			continue
		}

		supplierAddresses[userID] = addressID
	}

	// Tạo supplier profiles
	supplierIDs := make([]int64, 0, len(supplierUserIDs))

	for _, userID := range supplierUserIDs {
		addressID, ok := supplierAddresses[userID]
		if !ok {
			continue
		}

		// Lấy thông tin người dùng để tạo tên công ty
		var fullname string
		err := apiDb.QueryRow(ctx, `SELECT fullname FROM users WHERE id = $1`, userID).Scan(&fullname)
		if err != nil {
			log.Printf("Warning: Cannot get fullname for user_id: %d", userID)
			continue
		}

		// Tạo tên công ty từ tên người dùng
		companyName := fmt.Sprintf("%s Shop", fullname)

		// Lấy thông tin số điện thoại từ addresses
		var phone string
		err = apiDb.QueryRow(ctx, `SELECT phone FROM addresses WHERE id = $1`, addressID).Scan(&phone)
		if err != nil {
			log.Printf("Warning: Cannot get phone from address: %d", addressID)
			phone = fmt.Sprintf("+84%d", gofakeit.Number(300000000, 999999999))
		}

		// Tạo mã số thuế ngẫu nhiên
		taxID := fmt.Sprintf("%d-%d", gofakeit.Number(1000000000, 9999999999), gofakeit.Number(100, 999))

		// Tạo logo từ tên công ty
		logoURL := fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random", strings.ReplaceAll(companyName, " ", "+"))

		// Tạo mô tả công ty
		descriptions := []string{
			"Chúng tôi cam kết mang đến những sản phẩm chất lượng cao với giá cả cạnh tranh.",
			"Được thành lập từ năm 2018, chúng tôi đã phục vụ hàng ngàn khách hàng trên toàn quốc.",
			"Chuyên cung cấp các sản phẩm chính hãng, mới 100% và bảo hành theo tiêu chuẩn nhà sản xuất.",
			"Đối tác chính thức của nhiều thương hiệu lớn, chúng tôi tự hào về chất lượng dịch vụ và sự hài lòng của khách hàng.",
			"Với đội ngũ nhân viên tận tâm, chúng tôi cam kết mang lại trải nghiệm mua sắm tốt nhất cho khách hàng.",
		}
		description := descriptions[gofakeit.Number(0, len(descriptions)-1)]

		var supplierID int64
		err = partnerDb.QueryRow(ctx, `
			INSERT INTO supplier_profiles (user_id, company_name, contact_phone, description, logo_url, business_address_id, tax_id, status)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (tax_id) DO UPDATE
			SET company_name = $2, contact_phone = $3, description = $4, logo_url = $5, business_address_id = $6, status = $8
			RETURNING id;
		`, userID, companyName, phone, description, logoURL, addressID, taxID, "active").Scan(&supplierID)

		if err != nil {
			log.Printf("Error inserting supplier profile: %v", err)
			continue
		}

		supplierIDs = append(supplierIDs, supplierID)

		// Tạo supplier document với JSON documents
		documentsJSON := map[string]string{
			"id_card_front":    "https://images.unsplash.com/photo-1633332755192-727a05c4013d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=880&q=80",
			"id_card_back":     "https://images.unsplash.com/photo-1560250097-0b93528c311a?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
			"business_license": "https://images.unsplash.com/photo-1554224155-6726b3ff858f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"tax_certificate":  "https://images.unsplash.com/photo-1600880292203-757bb62b4baf?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		}

		documentsBytes, err := json.Marshal(documentsJSON)
		if err != nil {
			log.Printf("Error marshaling documents JSON: %v", err)
			continue
		}

		_, err = partnerDb.Exec(ctx, `
			INSERT INTO supplier_documents (supplier_id, documents, verification_status, admin_note)
			VALUES ($1, $2, 'approved', 'Đã xác thực hồ sơ nhà cung cấp - Tài liệu đầy đủ và hợp lệ')
			ON CONFLICT DO NOTHING;
		`, supplierID, documentsBytes)

		if err != nil {
			log.Printf("Error inserting supplier document: %v", err)
		}
	}

	log.Printf("✅ Created %d supplier profiles", len(supplierIDs))
	return supplierIDs
}

// Cải tiến seedEnhancedProducts để đảm bảo seedTags được gọi trước khi thêm sản phẩm
func seedEnhancedProducts(ctx context.Context, db *pgxpool.Pool, supplierIDs []int64) {
	seedTags(ctx, db)
	seedAttributeDefinitions(ctx, db)

	// Check if product_variants table exists
	var exists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'product_variants'
		)
	`).Scan(&exists)

	if err != nil {
		log.Printf("Error checking if product_variants table exists: %v", err)
		return
	}

	// Danh sách sản phẩm cụ thể với ảnh chính xác
	// THAY THẾ PHẦN specificProducts TRONG HÀM seedEnhancedProducts BẰNG CODE NÀY
	specificProducts := []struct {
		name          string
		description   string
		category      string
		imageURL      string
		variantImages map[string]string
	}{
		// ĐIỆN TỬ
		{
			name:        "iPhone 15 Pro Max 256GB",
			description: "Điện thoại thông minh cao cấp với chip A17 Pro, camera 48MP và màn hình Super Retina XDR 6.7 inch",
			category:    "Điện tử",
			imageURL:    "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/iphone-15-pro-max-bluetitanium-select?wid=470&hei=556&fmt=png-alpha&.v=1692895204112",
			variantImages: map[string]string{
				"Xanh":  "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/iphone-15-pro-max-bluetitanium-select?wid=470&hei=556&fmt=png-alpha&.v=1692895204112",
				"Đen":   "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/iphone-15-pro-max-blacktitanium-select?wid=470&hei=556&fmt=png-alpha&.v=1692895204112",
				"Trắng": "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/iphone-15-pro-max-whitetitanium-select?wid=470&hei=556&fmt=png-alpha&.v=1692895204112",
			},
		},
		{
			name:        "Samsung Galaxy S24 Ultra",
			description: "Flagship Android với bút S Pen, camera 200MP và màn hình Dynamic AMOLED 2X",
			category:    "Điện tử",
			imageURL:    "https://cdn.tgdd.vn/Products/Images/42/307174/samsung-galaxy-s24-ultra-xam-1-750x500.jpg",
			variantImages: map[string]string{
				"Xám":  "https://cdn.tgdd.vn/Products/Images/42/307174/samsung-galaxy-s24-ultra-xam-1-750x500.jpg",
				"Đen":  "https://cdn.tgdd.vn/Products/Images/42/307174/samsung-galaxy-s24-ultra-den-1-750x500.jpg",
				"Tím":  "https://cdn.tgdd.vn/Products/Images/42/307174/samsung-galaxy-s24-ultra-tim-1-750x500.jpg",
				"Vàng": "https://cdn.tgdd.vn/Products/Images/42/307174/samsung-galaxy-s24-ultra-vang-1-750x500.jpg",
			},
		},
		{
			name:        "MacBook Air M3 13 inch",
			description: "Laptop siêu mỏng với chip Apple M3 và thời lượng pin 18 giờ",
			category:    "Điện tử",
			imageURL:    "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/mba13-starlight-select-202402?wid=904&hei=840&fmt=jpeg&qlt=90&.v=1708367688034",
			variantImages: map[string]string{
				"Vàng": "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/mba13-starlight-select-202402?wid=904&hei=840&fmt=jpeg&qlt=90&.v=1708367688034",
				"Bạc":  "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/mba13-silver-select-202402?wid=904&hei=840&fmt=jpeg&qlt=90&.v=1708367688034",
				"Xám":  "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/mba13-spacegray-select-202402?wid=904&hei=840&fmt=jpeg&qlt=90&.v=1708367688034",
			},
		},
		{
			name:        "iPad Pro M4 11 inch",
			description: "Máy tính bảng cao cấp với chip M4 và màn hình Liquid Retina XDR",
			category:    "Điện tử",
			imageURL:    "https://cdn.tgdd.vn/Products/Images/522/325513/ipad-pro-11-inch-m4-wifi-black-3-750x500.jpg",
			variantImages: map[string]string{
				"Đen": "https://cdn.tgdd.vn/Products/Images/522/325513/ipad-pro-11-inch-m4-wifi-black-1-750x500.jpg",
				"Bạc": "https://cdn.tgdd.vn/Products/Images/522/325513/ipad-pro-11-inch-m4-wifi-sliver-1-750x500.jpg",
			},
		},
		{
			name:        "AirPods Pro Gen 2",
			description: "Tai nghe true wireless với chống ồn chủ động và chip H2",
			category:    "Điện tử",
			imageURL:    "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/MQD83?wid=572&hei=572&fmt=jpeg&qlt=95&.v=1660803972361",
			variantImages: map[string]string{
				"Trắng": "https://store.storeimages.cdn-apple.com/1/as-images.apple.com/is/MQD83?wid=572&hei=572&fmt=jpeg&qlt=95&.v=1660803972361",
			},
		},

		// THỜI TRANG
		{
			name:        "Áo thun nam Uniqlo cổ tròn",
			description: "Áo thun nam chất liệu cotton 100% mềm mại và thoáng mát",
			category:    "Thời trang",
			imageURL:    "https://image.uniqlo.com/UQ/ST3/AsianCommon/imagesgoods/422992/sub/goods_422992_sub14_3x4.jpg?width=369",
			variantImages: map[string]string{
				"Đen":   "https://image.uniqlo.com/UQ/ST3/vn/imagesgoods/422992/item/vngoods_09_422992_3x4.jpg?width=369",
				"Trắng": "https://image.uniqlo.com/UQ/ST3/vn/imagesgoods/422992/item/vngoods_00_422992_3x4.jpg?width=369",
				"Xanh":  "https://image.uniqlo.com/UQ/ST3/vn/imagesgoods/422992/item/vngoods_67_422992_3x4.jpg?width=369",
				"Hồng":  "https://image.uniqlo.com/UQ/ST3/vn/imagesgoods/422992/item/vngoods_12_422992_3x4.jpg?width=369",
				"Nâu":   "https://image.uniqlo.com/UQ/ST3/vn/imagesgoods/422992/item/vngoods_35_422992_3x4.jpg?width=369",
			},
		},
		{
			name:        "Giày thể thao Nike Air Force 1",
			description: "Giày thể thao kinh điển với đế Air cushion và da thật cao cấp",
			category:    "Thời trang",
			imageURL:    "https://static.nike.com/a/images/t_PDP_1728_v1/f_auto,q_auto:eco/b7d9211c-26e7-431a-ac24-b0540fb3c00f/AIR+FORCE+1+%2707.png",
			variantImages: map[string]string{
				"Trắng": "https://static.nike.com/a/images/t_PDP_1728_v1/f_auto,q_auto:eco/b7d9211c-26e7-431a-ac24-b0540fb3c00f/AIR+FORCE+1+%2707.png",
				"Đen":   "https://static.nike.com/a/images/t_PDP_1728_v1/f_auto,q_auto:eco/fc4622c4-2769-4665-aa6e-42c974a7705e/AIR+FORCE+1+%2707.png",
			},
		},
		{
			name:        "Túi Đeo Chéo Nam Coach Borsa A Tracolla Sullivan In Tela Esclusiva",
			description: "Túi Đeo Chéo Coach Nam Borsa A Tracolla Sullivan In Tela Esclusiva là chiếc túi dành cho phái mạnh đến từ thương hiệu Coach nổi tiếng. Túi được làm từ chất liệu cao cấp, bền đẹp trong suốt quá trình sử dụng.",
			category:    "Thời trang",
			imageURL:    "https://cdn.vuahanghieu.com/unsafe/0x500/left/top/smart/filters:quality(90)/https://admin.vuahanghieu.com/upload/product/2022/12/tui-deo-cheo-coach-nam-borsa-a-tracolla-sullivan-in-tela-esclusiva-cc009-mau-den-xam-63aab3ce1ee15-27122022155854.jpg",
			variantImages: map[string]string{
				"Đen": "https://cdn.vuahanghieu.com/unsafe/0x900/left/top/smart/filters:quality(90)/https://admin.vuahanghieu.com/upload/product/2022/12/tui-deo-cheo-coach-nam-borsa-a-tracolla-sullivan-in-tela-esclusiva-cc009-mau-den-xam-63aab3ce18523-27122022155854.jpg",
			},
		},

		// GIA DỤNG
		{
			name:        "Nồi cơm điện Panasonic 1.8L",
			description: "Nồi cơm điện thông minh với công nghệ IH và lòng nồi chống dính",
			category:    "Gia dụng",
			imageURL:    "https://cdnv2.tgdd.vn/mwg-static/dmx/Products/Images/1922/335998/noi-com-dien-tu-panasonic-1-8-lit-sr-dm184kra-1-638827266994887155-700x467.jpg",
			variantImages: map[string]string{
				"Đen": "https://cdnv2.tgdd.vn/mwg-static/dmx/Products/Images/1922/335998/noi-com-dien-tu-panasonic-1-8-lit-sr-dm184kra-1-638827266994887155-700x467.jpg",
			},
		},
		{
			name:        "Máy lọc không khí Xiaomi Mi Air Purifier 4",
			description: "Máy lọc không khí thông minh với bộ lọc HEPA H13",
			category:    "Gia dụng",
			imageURL:    "https://cdn.tgdd.vn/Products/Images/5473/314385/Slider/xiaomi-smart-air-purifier-4-compact-eu-bhr5860eu-27w638304564376429632.jpg",
			variantImages: map[string]string{
				"Trắng": "https://cdn.tgdd.vn/Products/Images/5473/314385/Slider/xiaomi-smart-air-purifier-4-compact-eu-bhr5860eu-27w638304564376429632.jpg",
			},
		},
		{
			name:        "Nồi chiên không dầu Philips 4.1L",
			description: "Nồi chiên không dầu công nghệ Rapid Air với dung tích 4.1L",
			category:    "Gia dụng",
			imageURL:    "https://euromixx.vn/wp-content/uploads/2024/10/010406034817-150x150.jpg",
			variantImages: map[string]string{
				"Đen": "https://euromixx.vn/wp-content/uploads/2024/10/010406034817-150x150.jpg",
			},
		},

		// SÁCH
		{
			name:        "Đắc Nhân Tâm - Dale Carnegie",
			description: "Cuốn sách kinh điển về nghệ thuật giao tiếp và ứng xử",
			category:    "Sách",
			imageURL:    "https://static.oreka.vn/800-800_881c12cb-3fb9-4011-b6dd-fbf459fc0b92.webp",
			variantImages: map[string]string{
				"Bìa cứng": "https://static.oreka.vn/800-800_881c12cb-3fb9-4011-b6dd-fbf459fc0b92.webp",
			},
		},
		{
			name:        "Nhà Giả Kim - Paulo Coelho",
			description: "Tiểu thuyết nổi tiếng thế giới về hành trình tìm kiếm ý nghĩa cuộc sống",
			category:    "Sách",
			imageURL:    "https://static.oreka.vn/800-800_8177be27-d7d5-4715-99c2-3993b4b65ba7.webp",
			variantImages: map[string]string{
				"Bìa cứng": "https://static.oreka.vn/800-800_8177be27-d7d5-4715-99c2-3993b4b65ba7.webp",
			},
		},
		{
			name:        "Sapiens - Lược Sử Loài Người Bằng Tranh",
			description: "Lược sử loài người - cuốn sách về lịch sử và tương lai nhân loại",
			category:    "Sách",
			imageURL:    "https://static.oreka.vn/800-800_4bf9560a-972f-4ed2-979c-225d44f9cb18.webp",
			variantImages: map[string]string{
				"Bìa cứng": "https://static.oreka.vn/800-800_4bf9560a-972f-4ed2-979c-225d44f9cb18.webp",
			},
		},

		// THỂ THAO
		{
			name:        "Tạ tay điều chỉnh 20kg",
			description: "Tạ tay thông minh có thể điều chỉnh từ 2-20kg tiết kiệm không gian",
			category:    "Thể thao",
			imageURL:    "https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=256x0&format=auto 256w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=384x0&format=auto 384w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=768x0&format=auto 768w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=1024x0&format=auto 1024w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=1440x0&format=auto 1440w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=1920x0&format=auto 1920w",
			variantImages: map[string]string{
				"Đen": "https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=256x0&format=auto 256w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=384x0&format=auto 384w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=768x0&format=auto 768w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=1024x0&format=auto 1024w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=1440x0&format=auto 1440w, https://contents.mediadecathlon.com/p2599290/k$a0f24e15abaf3f5be0cd667889c26627/b%E1%BB%99-t%E1%BA%A1-th%C3%A1o-l%E1%BA%AFp-20-kg-corength-8018574.jpg?f=1920x0&format=auto 1920w",
			},
		},
		{
			name:        "Thảm tập yoga chống trượt 6mm",
			description: "Thảm tập yoga cao cấp dày 6mm chống trượt và thân thiện môi trường",
			category:    "Thể thao",
			imageURL:    "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lp2k4n79n2aj25",
			variantImages: map[string]string{
				"Xanh": "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lp2k2y25ckm380",
				"Đen":  "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lp2k2y25fdqzfb",
				"Hồng": "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lp2k2y25gsbfbc",
			},
		},
		{
			name:        "Bóng đá FIFA Quality Pro",
			description: "Bóng đá chuẩn FIFA Quality Pro cho thi đấu chuyên nghiệp",
			category:    "Thể thao",
			imageURL:    "https://contents.mediadecathlon.com/p2571168/k$984defa0d32944089839f3e3c2e08b80/qu%E1%BA%A3-b%C3%B3ng-%C4%91%C3%A1-theo-ti%C3%AAu-chu%E1%BA%A9n-fifa-quality-pro-li%C3%AAn-k%E1%BA%BFt-nhi%E1%BB%87t-c%E1%BB%A1-5-pro-tr%E1%BA%AFng-kipsta-8827905.jpg?f=1920x0&format=auto",
			variantImages: map[string]string{
				"Trắng": "https://contents.mediadecathlon.com/p2571168/k$984defa0d32944089839f3e3c2e08b80/qu%E1%BA%A3-b%C3%B3ng-%C4%91%C3%A1-theo-ti%C3%AAu-chu%E1%BA%A9n-fifa-quality-pro-li%C3%AAn-k%E1%BA%BFt-nhi%E1%BB%87t-c%E1%BB%A1-5-pro-tr%E1%BA%AFng-kipsta-8827905.jpg?f=1920x0&format=auto",
				"Đỏ":    "https://contents.mediadecathlon.com/p2571092/k$bae70824ac05b30df5625da4456afad2/qu%E1%BA%A3-b%C3%B3ng-%C4%91%C3%A1-c%E1%BB%A1-5-chu%E1%BA%A9n-fifa-quality-pro-pro-%C4%91%E1%BB%8F-kipsta-8827906.jpg?f=1920x0&format=auto",
			},
		},

		// LÀM ĐẸP
		{
			name:        "Son môi YSL Rouge Pur Couture",
			description: "Son môi cao cấp với công thức dưỡng ẩm và màu sắc lâu trôi",
			category:    "Làm đẹp",
			imageURL:    "https://lipstick.vn/wp-content/uploads/2016/01/son-ysl-mau-201-orange-imagine.jpg",
			variantImages: map[string]string{
				"Đỏ": "\t\t\t\t\"Hồng\": \"https://images.unsplash.com/photo-1522335789203-aabd1fc54bc9?w=600&h=600&fit=crop&crop=center\",\n",
			},
		},
		{
			name:        "Kem chống nắng La Roche-Posay SPF 60",
			description: "Kem chống nắng dành cho da nhạy cảm SPF 60 PA++++",
			category:    "Làm đẹp",
			imageURL:    "https://down-vn.img.susercontent.com/file/vn-11134207-7qukw-lf6o7ah0nibud3",
			variantImages: map[string]string{
				"50ml":  "https://down-vn.img.susercontent.com/file/vn-11134207-7qukw-lf6o7ah0nibud3",
				"100ml": "https://down-vn.img.susercontent.com/file/vn-11134207-7qukw-lf6o7ah0nibud3",
			},
		},

		// THỰC PHẨM
		{
			name:        "Gạo ST25 Đồng Tháp",
			description: "Gạo thơm ST25 chất lượng cao từ Đồng Tháp hạt dẻo mềm thơm ngon",
			category:    "Thực phẩm",
			imageURL:    "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lz1zx9c987kdc5",
			variantImages: map[string]string{
				"5kg":  "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lz1zx9c987kdc5",
				"10kg": "https://down-vn.img.susercontent.com/file/vn-11134207-7r98o-lz1zx9c987kdc5",
			},
		},
		{
			name:        "Cà phê Trung Nguyên Legend 5",
			description: "Cà phê pha phin truyền thống pha chế từ 100% cà phê Robusta và Arabica",
			category:    "Thực phẩm",
			imageURL:    "https://salt.tikicdn.com/cache/750x750/ts/product/8f/3b/25/c73e2d9ef40438a229d06c5dc4ac035f.jpg.webp",
			variantImages: map[string]string{
				"500g": "https://salt.tikicdn.com/cache/750x750/ts/product/8f/3b/25/c73e2d9ef40438a229d06c5dc4ac035f.jpg.webp",
				"340g": "https://salt.tikicdn.com/cache/750x750/ts/product/8f/3b/25/c73e2d9ef40438a229d06c5dc4ac035f.jpg.webp",
			},
		},

		// NỘI THẤT
		{
			name:        "Sofa 3 chỗ IKEA Ektorp",
			description: "Sofa 3 chỗ ngồi bọc vải khung gỗ thông thiết kế Scandinavian",
			category:    "Nội thất",
			imageURL:    "https://kika.vn/wp-content/uploads/2022/09/ghe-sofa-vang-3-cho-ngoi-boc-da-han-cao-cap-sf90-1.jpg",
			variantImages: map[string]string{
				"Xanh": "https://kika.vn/wp-content/uploads/2022/09/ghe-sofa-vang-3-cho-ngoi-boc-da-han-cao-cap-sf90-1.jpg",
				"Xám":  "https://kika.vn/wp-content/uploads/2022/09/ghe-sofa-vang-3-cho-ngoi-boc-da-han-cao-cap-sf90-4.jpg",
			},
		},
		{
			name:        "Bàn IKEA gaming 2 hộc tủ",
			description: "Bàn làm việc gỗ thông tự nhiên với 2 ngăn kéo thiết kế cổ điển",
			category:    "Nội thất",
			imageURL:    "https://noithatdangkhoa.com/wp-content/uploads/2024/06/ban-ikea-gaming-2-hoc-tu-blvdk46-1.jpg",
			variantImages: map[string]string{
				"Trắng": "https://noithatdangkhoa.com/wp-content/uploads/2024/06/ban-ikea-gaming-2-hoc-tu-blvdk46-2.jpg",
				"Đen":   "https://images.unsplash.com/photo-1506439773649-6e0eb8cfb237?w=600&h=600&fit=crop&crop=center",
			},
		},
	}

	// Lấy categories
	categoryNameToID := make(map[string]int64)
	rows, err := db.Query(ctx, `SELECT id, name FROM categories`)
	if err != nil {
		log.Fatalf("Error getting categories: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}
		categoryNameToID[name] = id
	}

	// Tạo sản phẩm
	totalProducts := 0
	for _, productInfo := range specificProducts {
		categoryID, categoryExists := categoryNameToID[productInfo.category]
		if !categoryExists {
			log.Printf("Category not found: %s", productInfo.category)
			continue
		}

		// Chọn supplier ngẫu nhiên
		supplierID := supplierIDs[gofakeit.Number(0, len(supplierIDs)-1)]

		// Tạo SKU prefix
		skuPrefix := strings.ToUpper(string([]rune(productInfo.category)[0])) +
			strings.ToUpper(string([]rune(productInfo.name)[0])) +
			fmt.Sprintf("%03d", gofakeit.Number(100, 999))

		// Kiểm tra sản phẩm đã tồn tại chưa
		var existingID string
		err := db.QueryRow(ctx, `
			SELECT id FROM products WHERE name = $1 AND supplier_id = $2 AND category_id = $3
		`, productInfo.name, supplierID, categoryID).Scan(&existingID)

		var productID string
		if err != nil && err != pgx.ErrNoRows {
			log.Printf("Error checking product existence: %v", err)
			continue
		}

		if err == pgx.ErrNoRows {
			// Tạo mới sản phẩm
			err := db.QueryRow(ctx, `
				INSERT INTO products (
					supplier_id, category_id, name, description, image_url,
					status, featured, tax_class, sku_prefix, average_rating
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				RETURNING id;
			`,
				supplierID, categoryID, productInfo.name, productInfo.description, productInfo.imageURL,
				"active", gofakeit.Bool(), "standard", skuPrefix, float32(gofakeit.Float32Range(3.5, 5)),
			).Scan(&productID)

			if err != nil {
				log.Printf("Error inserting product: %v", err)
				continue
			}
		} else {
			productID = existingID
		}

		// Thêm tags
		numTags := gofakeit.Number(1, 3)
		tagNames := []string{"Mới nhất", "Bán chạy", "Chính hãng", "Giảm giá", "Chất lượng cao"}

		for j := 0; j < numTags; j++ {
			randomTag := tagNames[gofakeit.Number(0, len(tagNames)-1)]

			var tagID string
			err := db.QueryRow(ctx, `SELECT id FROM tags WHERE name = $1`, randomTag).Scan(&tagID)
			if err != nil {
				continue
			}

			var relationExists bool
			err = db.QueryRow(ctx, `
				SELECT EXISTS (
					SELECT 1 FROM products_tags WHERE product_id = $1 AND tag_id = $2
				)
			`, productID, tagID).Scan(&relationExists)

			if err == nil && !relationExists {
				db.Exec(ctx, `INSERT INTO products_tags (product_id, tag_id) VALUES ($1, $2)`, productID, tagID)
			}
		}

		// Tạo variants (3-5 variants)
		if exists {
			createOptimizedProductVariants(ctx, db, productID, skuPrefix, productInfo)
		}

		totalProducts++
	}

	log.Printf("✅ Created %d products successfully", totalProducts)
}

func createOptimizedProductVariants(
	ctx context.Context,
	db *pgxpool.Pool,
	productID string,
	skuPrefix string,
	productInfo struct {
		name          string
		description   string
		category      string
		imageURL      string
		variantImages map[string]string
	},
) {
	// Lấy thuộc tính cho category
	categoryAttrs, ok := categoryAttributes[productInfo.category]
	if !ok {
		// Fallback attributes nếu category không có
		categoryAttrs = map[string][]string{
			"Màu sắc": {"Đen", "Trắng", "Xám"},
		}
	}

	log.Printf("Creating variants for product: %s with category: %s", productInfo.name, productInfo.category)

	// Tạo tất cả combinations có thể từ category attributes
	var variants []map[string]string

	// Lấy tất cả attribute names và options từ category
	attrNames := make([]string, 0, len(categoryAttrs))
	attrOptions := make(map[string][]string)

	for name, options := range categoryAttrs {
		attrNames = append(attrNames, name)
		attrOptions[name] = options
	}

	// Nếu có variantImages, ưu tiên tạo variants theo màu sắc có ảnh
	if len(productInfo.variantImages) > 0 {
		colorOptions := make([]string, 0, len(productInfo.variantImages))
		for color := range productInfo.variantImages {
			colorOptions = append(colorOptions, color)
		}
		attrOptions["Màu sắc"] = colorOptions

		// Đảm bảo "Màu sắc" là attribute đầu tiên
		if !stringInSlice("Màu sắc", attrNames) {
			attrNames = append([]string{"Màu sắc"}, attrNames...)
		}
	}

	// Tạo combinations thông minh dựa trên category
	switch productInfo.category {
	case "Điện tử":
		variants = createElectronicsVariants(attrOptions, productInfo.variantImages)
	case "Thời trang":
		variants = createFashionVariants(attrOptions, productInfo.variantImages)
	case "Gia dụng":
		variants = createApplianceVariants(attrOptions, productInfo.variantImages)
	case "Sách":
		variants = createBookVariants(attrOptions)
	case "Thể thao":
		variants = createSportsVariants(attrOptions, productInfo.variantImages)
	case "Làm đẹp":
		variants = createBeautyVariants(attrOptions, productInfo.variantImages)
	case "Thực phẩm":
		variants = createFoodVariants(attrOptions, productInfo.variantImages)
	case "Nội thất":
		variants = createFurnitureVariants(attrOptions, productInfo.variantImages)
	default:
		variants = createGeneralVariants(attrOptions, productInfo.variantImages)
	}

	// Đảm bảo mọi sản phẩm đều có ít nhất 1 variant
	if len(variants) == 0 {
		// Tạo variant cơ bản với thuộc tính đầu tiên
		defaultVariant := make(map[string]string)
		for attrName, options := range attrOptions {
			if len(options) > 0 {
				defaultVariant[attrName] = options[0]
				break
			}
		}
		if len(defaultVariant) == 0 {
			defaultVariant["Màu sắc"] = "Đen"
		}
		variants = append(variants, defaultVariant)
	}

	// Giới hạn tối đa 8 variants để tránh quá nhiều
	if len(variants) > 8 {
		variants = variants[:8]
	}

	log.Printf("Generated %d variants for product %s", len(variants), productInfo.name)

	// Lấy attribute definitions và options
	attributeDefs := make(map[string]int)
	attributeOptions := make(map[string]map[string]int)

	for _, variant := range variants {
		for attrName := range variant {
			if _, exists := attributeDefs[attrName]; !exists {
				var attrID int
				err := db.QueryRow(ctx, `SELECT id FROM attribute_definitions WHERE name = $1`, attrName).Scan(&attrID)
				if err != nil {
					log.Printf("Attribute definition not found: %s", attrName)
					continue
				}
				attributeDefs[attrName] = attrID
				attributeOptions[attrName] = make(map[string]int)

				// Lấy options
				rows, err := db.Query(ctx, `
					SELECT id, option_value FROM attribute_options 
					WHERE attribute_definition_id = $1
				`, attrID)
				if err != nil {
					log.Printf("Error getting attribute options: %v", err)
					continue
				}

				for rows.Next() {
					var optionID int
					var optionValue string
					if err := rows.Scan(&optionID, &optionValue); err == nil {
						attributeOptions[attrName][optionValue] = optionID
					}
				}
				rows.Close()
			}
		}
	}

	// Tạo variants
	for i, variant := range variants {
		// Tạo SKU unique
		timestamp := time.Now().UnixNano() % 1000000
		uniqueSKU := fmt.Sprintf("%s-%d-%d", skuPrefix, i+1, timestamp)

		// Tạo tên variant từ tất cả attributes
		var variantNameParts []string
		for attr, value := range variant {
			variantNameParts = append(variantNameParts, fmt.Sprintf("%s: %s", attr, value))
		}
		variantName := strings.Join(variantNameParts, ", ")

		// Chọn ảnh cho variant
		variantImage := productInfo.imageURL
		if colorValue, ok := variant["Màu sắc"]; ok {
			if specificImage, exists := productInfo.variantImages[colorValue]; exists {
				variantImage = specificImage
			}
		}

		// Tính giá dựa trên tất cả attributes
		basePrice := getBasePriceForCategory(productInfo.category)

		// Áp dụng multiplier cho từng attribute
		for attrName, attrValue := range variant {
			multiplier := getOptimizedPriceMultiplier(attrName, attrValue)
			basePrice = basePrice * multiplier
		}

		// Làm tròn giá
		basePrice = float32(math.Round(float64(basePrice/1000)) * 1000)

		// Discount ngẫu nhiên
		var discountPriceParam interface{} = nil
		if gofakeit.Bool() { // 50% chance có discount
			discountPercent := gofakeit.Float32Range(0.1, 0.3) // 10-30%
			discountPrice := basePrice * (1 - discountPercent)
			discountPrice = float32(math.Round(float64(discountPrice/1000)) * 1000)
			if discountPrice < basePrice {
				discountPriceParam = discountPrice
			}
		}

		// Insert variant
		var variantID string
		err := db.QueryRow(ctx, `
			INSERT INTO product_variants (
				product_id, sku, variant_name, price, discount_price,
				inventory_quantity, shipping_class, image_url, alt_text, is_default, is_active
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id;
		`,
			productID, uniqueSKU, variantName, basePrice, discountPriceParam,
			gofakeit.Number(10, 100), "standard", variantImage, variantName, i == 0, true,
		).Scan(&variantID)

		if err != nil {
			log.Printf("Error inserting variant: %v", err)
			continue
		}

		// Thêm tất cả attributes cho variant
		for attrName, attrValue := range variant {
			attrID, ok := attributeDefs[attrName]
			if !ok {
				log.Printf("Attribute definition not found: %s", attrName)
				continue
			}

			optionID, ok := attributeOptions[attrName][attrValue]
			if !ok {
				// Tạo option mới nếu chưa có
				err := db.QueryRow(ctx, `
					INSERT INTO attribute_options (attribute_definition_id, option_value)
					VALUES ($1, $2) RETURNING id;
				`, attrID, attrValue).Scan(&optionID)
				if err != nil {
					log.Printf("Error creating attribute option: %v", err)
					continue
				}
				attributeOptions[attrName][attrValue] = optionID
			}

			// Thêm attribute cho variant
			_, err = db.Exec(ctx, `
				INSERT INTO product_variant_attributes (
					product_variant_id, attribute_definition_id, attribute_option_id
				) VALUES ($1, $2, $3);
			`, variantID, attrID, optionID)

			if err != nil {
				log.Printf("Error inserting variant attribute: %v", err)
			}
		}

		log.Printf("Created variant: %s for product %s", variantName, productInfo.name)
	}

	log.Printf("Successfully created %d variants for product %s", len(variants), productInfo.name)
}

// Helper functions để tạo variants cho từng category

func createElectronicsVariants(attrOptions map[string][]string, variantImages map[string]string) []map[string]string {
	var variants []map[string]string

	// Electronics: Màu sắc + Dung lượng
	colors := getAttrOptions(attrOptions, "Màu sắc", []string{"Đen", "Trắng", "Xanh"})
	capacities := getAttrOptions(attrOptions, "Dung lượng", []string{"64GB", "128GB", "256GB"})

	// Tạo combinations với chỉ 2 attributes
	for _, color := range colors {
		for _, capacity := range capacities {
			variant := map[string]string{
				"Màu sắc":    color,
				"Dung lượng": capacity,
			}
			variants = append(variants, variant)
		}
	}

	return variants
}

func createFashionVariants(attrOptions map[string][]string, variantImages map[string]string) []map[string]string {
	var variants []map[string]string

	// Fashion: Màu sắc + Kích thước
	colors := getAttrOptions(attrOptions, "Màu sắc", []string{"Đen", "Trắng", "Xanh"})
	sizes := getAttrOptions(attrOptions, "Kích thước", []string{"S", "M", "L", "XL"})

	for _, color := range colors {
		for _, size := range sizes {
			variant := map[string]string{
				"Màu sắc":    color,
				"Kích thước": size,
			}
			variants = append(variants, variant)
		}
	}

	return variants
}

func createApplianceVariants(attrOptions map[string][]string, variantImages map[string]string) []map[string]string {
	var variants []map[string]string

	// Gia dụng: Màu sắc + Công suất
	colors := getAttrOptions(attrOptions, "Màu sắc", []string{"Đen", "Trắng", "Bạc"})
	powers := getAttrOptions(attrOptions, "Công suất", []string{"500W", "1000W", "1500W"})

	for _, color := range colors {
		for _, power := range powers {
			variant := map[string]string{
				"Màu sắc":   color,
				"Công suất": power,
			}
			variants = append(variants, variant)
		}
	}

	return variants
}

func createBookVariants(attrOptions map[string][]string) []map[string]string {
	var variants []map[string]string

	// Sách: Loại bìa + Ngôn ngữ
	covers := getAttrOptions(attrOptions, "Loại bìa", []string{"Bìa mềm", "Bìa cứng"})
	languages := getAttrOptions(attrOptions, "Ngôn ngữ", []string{"Tiếng Việt", "Tiếng Anh"})

	for _, cover := range covers {
		for _, lang := range languages {
			variant := map[string]string{
				"Loại bìa": cover,
				"Ngôn ngữ": lang,
			}
			variants = append(variants, variant)
		}
	}

	return variants
}

func createSportsVariants(attrOptions map[string][]string, variantImages map[string]string) []map[string]string {
	var variants []map[string]string

	// Thể thao: Màu sắc + Kích thước
	colors := getAttrOptions(attrOptions, "Màu sắc", []string{"Đen", "Xanh", "Đỏ"})
	sizes := getAttrOptions(attrOptions, "Kích thước", []string{"Nhỏ", "Vừa", "Lớn"})

	for _, color := range colors {
		for _, size := range sizes {
			variant := map[string]string{
				"Màu sắc":    color,
				"Kích thước": size,
			}
			variants = append(variants, variant)
		}
	}

	return variants
}

func createBeautyVariants(attrOptions map[string][]string, variantImages map[string]string) []map[string]string {
	var variants []map[string]string

	// Làm đẹp: Màu sắc + Dung tích
	colors := getAttrOptions(attrOptions, "Màu sắc", []string{"Đỏ", "Hồng", "Nude"})
	capacities := getAttrOptions(attrOptions, "Dung tích", []string{"15ml", "30ml", "50ml"})

	for _, color := range colors {
		for _, capacity := range capacities {
			variant := map[string]string{
				"Màu sắc":   color,
				"Dung tích": capacity,
			}
			variants = append(variants, variant)
		}
	}

	return variants
}

func createFoodVariants(attrOptions map[string][]string, variantImages map[string]string) []map[string]string {
	var variants []map[string]string

	// Thực phẩm: Trọng lượng + Xuất xứ
	weights := getAttrOptions(attrOptions, "Trọng lượng", []string{"340g", "500g", "1kg", "5kg"})
	origins := getAttrOptions(attrOptions, "Xuất xứ", []string{"Việt Nam", "Thái Lan", "Nhật Bản"})

	for _, weight := range weights {
		for _, origin := range origins {
			variant := map[string]string{
				"Trọng lượng": weight,
				"Xuất xứ":     origin,
			}
			variants = append(variants, variant)
		}
	}

	return variants
}

func createFurnitureVariants(attrOptions map[string][]string, variantImages map[string]string) []map[string]string {
	var variants []map[string]string

	// Nội thất: Màu sắc + Kích thước (BỎ Chất liệu)
	colors := getAttrOptions(attrOptions, "Màu sắc", []string{"Trắng", "Nâu", "Đen"})
	sizes := getAttrOptions(attrOptions, "Kích thước", []string{"Nhỏ", "Vừa", "Lớn"})

	for _, color := range colors {
		for _, size := range sizes {
			variant := map[string]string{
				"Màu sắc":    color,
				"Kích thước": size,
			}
			variants = append(variants, variant)
		}
	}

	return variants
}

func createGeneralVariants(attrOptions map[string][]string, variantImages map[string]string) []map[string]string {
	var variants []map[string]string

	// General case: Lấy 2 attributes đầu tiên và tạo combinations
	attrNames := make([]string, 0, len(attrOptions))
	for name := range attrOptions {
		attrNames = append(attrNames, name)
		if len(attrNames) >= 2 {
			break
		}
	}

	if len(attrNames) >= 2 {
		firstAttrOptions := attrOptions[attrNames[0]]
		secondAttrOptions := attrOptions[attrNames[1]]

		for _, firstVal := range firstAttrOptions {
			for _, secondVal := range secondAttrOptions {
				variant := map[string]string{
					attrNames[0]: firstVal,
					attrNames[1]: secondVal,
				}
				variants = append(variants, variant)
			}
		}
	} else if len(attrNames) == 1 {
		// Chỉ có 1 attribute, tạo variants đơn giản
		for _, val := range attrOptions[attrNames[0]] {
			variant := map[string]string{
				attrNames[0]: val,
			}
			variants = append(variants, variant)
		}
	}

	return variants
}

// Helper functions
func getAttrOptions(attrOptions map[string][]string, attrName string, defaultOptions []string) []string {
	if options, exists := attrOptions[attrName]; exists && len(options) > 0 {
		return options
	}
	return defaultOptions
}

func getBasePriceForCategory(category string) float32 {
	basePrices := map[string][2]float32{
		"Điện tử":    {500000, 5000000},  // 500k - 5tr
		"Thời trang": {100000, 2000000},  // 100k - 2tr
		"Gia dụng":   {200000, 3000000},  // 200k - 3tr
		"Sách":       {50000, 300000},    // 50k - 300k
		"Thể thao":   {100000, 1000000},  // 100k - 1tr
		"Làm đẹp":    {100000, 1500000},  // 100k - 1.5tr
		"Thực phẩm":  {20000, 500000},    // 20k - 500k
		"Nội thất":   {500000, 10000000}, // 500k - 10tr
	}

	if priceRange, exists := basePrices[category]; exists {
		return gofakeit.Float32Range(priceRange[0], priceRange[1])
	}

	return gofakeit.Float32Range(100000, 1000000) // Default range
}

func getOptimizedPriceMultiplier(attrName, attrValue string) float32 {
	multipliers := map[string]map[string]float32{
		"Kích thước": {
			"S": 0.9, "M": 1.0, "L": 1.1, "XL": 1.2, "XXL": 1.3,
		},
		"Dung lượng": {
			"64GB": 0.9, "128GB": 1.0, "256GB": 1.2, "512GB": 1.5, "1TB": 1.8,
		},
		"RAM": {
			"8GB": 1.0, "16GB": 1.2, "32GB": 1.5,
		},
		"Trọng lượng": {
			"340g": 0.9, "500g": 1.0, "1kg": 1.1, "5kg": 1.3, "10kg": 1.5,
		},
		"Dung tích": {
			"15ml": 0.8, "30ml": 1.0, "50ml": 1.2, "100ml": 1.4,
			"1L": 0.9, "1.8L": 1.0, "2L": 1.1, "4L": 1.3,
		},
	}

	if attrMultipliers, exists := multipliers[attrName]; exists {
		if multiplier, exists := attrMultipliers[attrValue]; exists {
			return multiplier
		}
	}
	return 1.0 // Default multiplier
}

// Cải tiến seedAttributeDefinitions để đảm bảo rằng các thuộc tính được tạo đúng
func seedAttributeDefinitions(ctx context.Context, db *pgxpool.Pool) {
	attributes := []struct {
		name, desc, inputType    string
		isFilterable, isRequired bool
	}{
		// Thuộc tính cơ bản - chỉ giữ lại những cái cần thiết
		{"Màu sắc", "Màu sắc của sản phẩm", "select", true, true},
		{"Kích thước", "Kích thước của sản phẩm", "select", true, true},
		{"Dung lượng", "Dung lượng lưu trữ", "select", true, false},
		{"Công suất", "Công suất thiết bị", "select", false, false},
		{"Ngôn ngữ", "Ngôn ngữ sách", "select", true, false},
		{"Loại bìa", "Loại bìa sách", "select", false, false},
		{"Dung tích", "Dung tích của sản phẩm", "select", true, false},
		{"Trọng lượng", "Trọng lượng của sản phẩm", "select", true, false}, // THÊM CÁI NÀY
		{"Xuất xứ", "Quốc gia xuất xứ", "select", false, false},
	}

	// Seed attribute definitions
	for _, attr := range attributes {
		// Kiểm tra xem thuộc tính đã tồn tại chưa
		var count int
		err := db.QueryRow(ctx, `
			SELECT COUNT(*) FROM attribute_definitions WHERE name = $1
		`, attr.name).Scan(&count)

		if err != nil {
			log.Printf("Error checking attribute existence: %v", err)
			continue
		}

		var attrID int
		if count == 0 {
			// Nếu không tồn tại, thêm mới
			err := db.QueryRow(ctx, `
				INSERT INTO attribute_definitions (name, description, input_type, is_filterable, is_required)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id;
			`, attr.name, attr.desc, attr.inputType, attr.isFilterable, attr.isRequired).Scan(&attrID)

			if err != nil {
				log.Printf("Error inserting attribute definition: %v", err)
				continue
			}
		} else {
			// Nếu đã tồn tại, lấy ID
			err := db.QueryRow(ctx, `
				SELECT id FROM attribute_definitions WHERE name = $1
			`, attr.name).Scan(&attrID)

			if err != nil {
				log.Printf("Error getting attribute ID: %v", err)
				continue
			}
		}

		// Seed attribute options based on category_attributes map
		if options, exists := getAttributeOptions(attr.name); exists {
			for _, option := range options {
				// Kiểm tra xem option đã tồn tại chưa
				var optionCount int
				err := db.QueryRow(ctx, `
					SELECT COUNT(*) FROM attribute_options 
					WHERE attribute_definition_id = $1 AND option_value = $2
				`, attrID, option).Scan(&optionCount)

				if err != nil {
					log.Printf("Error checking option existence: %v", err)
					continue
				}

				if optionCount == 0 {
					_, err := db.Exec(ctx, `
						INSERT INTO attribute_options (attribute_definition_id, option_value)
						VALUES ($1, $2);
					`, attrID, option)

					if err != nil {
						log.Printf("Error inserting attribute option: %v", err)
					}
				}
			}
		}
	}
	log.Println("✅ Attribute definitions and options seeded successfully")
}

// Hàm helper để lấy min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Function to create more realistic product variants with appropriate attributes and images
func createProductVariants(
	ctx context.Context,
	db *pgxpool.Pool,
	productID string,
	skuPrefix string,
	categoryAttrs map[string][]string,
	categoryName string,
	baseProductImage string,
) {
	// Get all valid attribute names from the category attributes
	var validAttrNames []string
	for attrName := range categoryAttrs {
		var exists bool
		err := db.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM attribute_definitions WHERE name = $1
			)
		`, attrName).Scan(&exists)

		if err != nil {
			log.Printf("Error checking attribute existence: %v", err)
			continue
		}

		if exists {
			validAttrNames = append(validAttrNames, attrName)
		}
	}

	// Define category-specific primary attributes
	categoryPrimaryAttributes := map[string][]string{
		"Thời trang nam":        {"Kích thước", "Màu sắc"},
		"Thời trang nữ":         {"Kích thước", "Màu sắc"},
		"Thời trang trẻ em":     {"Kích thước", "Màu sắc"},
		"Giày dép":              {"Kích thước", "Màu sắc"},
		"Điện thoại thông minh": {"Màu sắc", "Dung lượng"},
		"Máy tính xách tay":     {"Màu sắc", "RAM", "Ổ cứng"},
		"Máy tính bảng":         {"Màu sắc", "Dung lượng", "Kết nối"},
		"Tai nghe & Loa":        {"Màu sắc", "Kiểu đeo", "Loại kết nối"},
		"Máy ảnh & Máy quay":    {"Màu sắc", "Độ phân giải", "Cảm biến"},
		"Đồ gia dụng":           {"Màu sắc", "Công suất", "Chất liệu"},
		"Tủ lạnh & Tủ đông":     {"Màu sắc", "Dung tích", "Công suất"},
		"Đồ dùng phòng ngủ":     {"Kích thước giường", "Màu sắc", "Chất liệu"},
		"Nội thất phòng khách":  {"Chất liệu", "Màu sắc", "Kích thước"},
	}

	// Choose appropriate attributes for this category
	var chosenAttributes []string
	if primaryAttrs, ok := categoryPrimaryAttributes[categoryName]; ok {
		// Use category-specific attributes if available
		for _, attr := range primaryAttrs {
			if stringInSlice(attr, validAttrNames) {
				chosenAttributes = append(chosenAttributes, attr)
			}
			if len(chosenAttributes) >= 2 {
				break
			}
		}
	}

	// If no specific attributes were chosen, use sensible defaults based on valid attrs
	if len(chosenAttributes) == 0 {
		// Try to find common important attributes
		commonImportantAttrs := []string{"Kích thước", "Màu sắc", "Chất liệu", "Dung lượng"}
		for _, attr := range commonImportantAttrs {
			if stringInSlice(attr, validAttrNames) {
				chosenAttributes = append(chosenAttributes, attr)
			}
			if len(chosenAttributes) >= 2 {
				break
			}
		}

		// If still no attributes, use first available
		if len(chosenAttributes) == 0 && len(validAttrNames) > 0 {
			chosenAttributes = append(chosenAttributes, validAttrNames[0])
		}
	}

	// Ensure we have at least one attribute
	if len(chosenAttributes) == 0 {
		log.Printf("No valid attributes found for product: %s in category: %s", productID, categoryName)
		return
	}

	// Get attribute definitions and options
	attributeDefs := make(map[string]int)
	attributeOptions := make(map[string]map[string]int)

	for _, attrName := range chosenAttributes {
		var attrID int
		err := db.QueryRow(ctx, `SELECT id FROM attribute_definitions WHERE name = $1`, attrName).Scan(&attrID)
		if err != nil {
			log.Printf("Error getting attribute definition: %v", err)
			continue
		}

		attributeDefs[attrName] = attrID
		attributeOptions[attrName] = make(map[string]int)

		// Get options for this attribute
		rows, err := db.Query(ctx, `
			SELECT id, option_value FROM attribute_options 
			WHERE attribute_definition_id = $1
		`, attrID)

		if err != nil {
			log.Printf("Error getting attribute options: %v", err)
			continue
		}

		defer rows.Close()
		for rows.Next() {
			var optionID int
			var optionValue string
			if err := rows.Scan(&optionID, &optionValue); err != nil {
				log.Printf("Error scanning attribute option: %v", err)
				continue
			}
			attributeOptions[attrName][optionValue] = optionID
		}
	}

	// Generate variant combinations based on chosen attributes
	variantCombinations := generateVariantCombinations(chosenAttributes, categoryAttrs)

	// Variant-specific image collection based on product category
	categorySpecificImageURLs := map[string]map[string]string{
		"Màu sắc": {
			"Đen":   "https://images.unsplash.com/photo-1622434641406-a158123450f9?ixlib=rb-4.0.3&q=80&w=1000",
			"Trắng": "https://images.unsplash.com/photo-1622434641406-a158123450f9?ixlib=rb-4.0.3&q=80&w=1000&auto=format&fit=crop&ixlib=rb-4.0.3",
			"Xanh":  "https://images.unsplash.com/photo-1606107557195-0e29a4b5b4aa?ixlib=rb-4.0.3&q=80&w=1000&auto=format&fit=crop",
			"Đỏ":    "https://images.unsplash.com/photo-1542291026-7eec264c27ff?ixlib=rb-4.0.3&q=80&w=1000&auto=format&fit=crop",
			"Vàng":  "https://images.unsplash.com/photo-1595950653106-6c9ebd614d3a?ixlib=rb-4.0.3&q=80&w=1000&auto=format&fit=crop",
		},
	}

	// Create variants for each combination
	for i, combo := range variantCombinations {
		// Generate a unique SKU
		timestamp := time.Now().UnixNano() % 1000000
		uniqueSKU := fmt.Sprintf("%s-%d-%d", skuPrefix, i+1, timestamp)

		// Create variant name from combination
		var variantNameParts []string
		for attr, value := range combo {
			variantNameParts = append(variantNameParts, fmt.Sprintf("%s: %s", attr, value))
		}
		variantName := strings.Join(variantNameParts, ", ")

		// Select appropriate image for the variant
		variantImage := baseProductImage
		// If this is a color variant, try to find a color-specific image
		if colorValue, ok := combo["Màu sắc"]; ok {
			if colorImages, ok := categorySpecificImageURLs["Màu sắc"]; ok {
				if img, ok := colorImages[colorValue]; ok {
					variantImage = img
				}
			}
		}

		// Set pricing based on variant attributes
		basePrice := gofakeit.Float32Range(100000, 5000000) // 100k - 5tr VND
		// Round price to nearest 1000 VND
		basePrice = float32(math.Round(float64(basePrice/1000)) * 1000)

		// Size affects price - larger sizes cost more
		if sizeValue, ok := combo["Kích thước"]; ok {
			sizeMultipliers := map[string]float32{
				"S":     0.9,
				"M":     1.0,
				"L":     1.1,
				"XL":    1.2,
				"XXL":   1.3,
				"XXXL":  1.4,
				"4GB":   0.8,
				"8GB":   1.0,
				"16GB":  1.3,
				"32GB":  1.6,
				"64GB":  1.9,
				"128GB": 1.2,
				"256GB": 1.5,
				"512GB": 1.8,
				"1TB":   2.0,
			}

			if multiplier, ok := sizeMultipliers[sizeValue]; ok {
				basePrice = basePrice * multiplier
				// Round again after applying multiplier
				basePrice = float32(math.Round(float64(basePrice/1000)) * 1000)
			}
		}

		// Set discount price (for some variants)
		discountPrice := basePrice
		hasDiscount := gofakeit.Bool()
		if hasDiscount {
			discountPercent := gofakeit.Float32Range(0.05, 0.25)
			discountAmount := float32(math.Floor(float64(basePrice*discountPercent)/1000) * 1000)

			// Ensure minimum discount is 5000 VND and discount price is lower than base price
			if discountAmount < 5000 {
				discountAmount = 5000
			}

			// Limit maximum discount to 80% of base price
			if discountAmount > basePrice*0.8 {
				discountAmount = float32(math.Floor(float64(basePrice*0.8)/1000) * 1000)
			}

			discountPrice = basePrice - discountAmount

			// Final check to ensure discount price is valid
			if discountPrice >= basePrice || discountPrice <= 0 {
				discountPrice = basePrice * 0.85 // Default 15% discount
			}
		}

		// Insert variant
		var variantID string
		var discountPriceParam interface{}
		if hasDiscount {
			discountPriceParam = discountPrice
		} else {
			discountPriceParam = nil
		}

		err := db.QueryRow(ctx, `
			INSERT INTO product_variants (
				product_id, sku, variant_name, price, discount_price,
				inventory_quantity, shipping_class, image_url, alt_text, is_default, is_active
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id;
		`,
			productID, uniqueSKU, variantName, basePrice, discountPriceParam,
			gofakeit.Number(5, 100), "standard", variantImage, variantName, i == 0, true,
		).Scan(&variantID)

		if err != nil {
			log.Printf("Error inserting product variant: %v", err)
			continue
		}

		// Add attribute associations for this variant
		for attr, value := range combo {
			attrID, ok := attributeDefs[attr]
			if !ok {
				continue
			}

			optionID, ok := attributeOptions[attr][value]
			if !ok {
				// Create option if it doesn't exist
				err := db.QueryRow(ctx, `
					INSERT INTO attribute_options (attribute_definition_id, option_value)
					VALUES ($1, $2)
					RETURNING id;
				`, attrID, value).Scan(&optionID)

				if err != nil {
					log.Printf("Error creating attribute option: %v", err)
					continue
				}
				attributeOptions[attr][value] = optionID
			}

			// Add the attribute to the variant
			_, err = db.Exec(ctx, `
				INSERT INTO product_variant_attributes (
					product_variant_id, attribute_definition_id, attribute_option_id
				)
				VALUES ($1, $2, $3);
			`, variantID, attrID, optionID)

			if err != nil {
				log.Printf("Error inserting product variant attribute: %v", err)
			}
		}
	}

	log.Printf("Created %d variants for product %s", len(variantCombinations), productID)
}

// Helper function to generate variant combinations
func generateVariantCombinations(attributes []string, categoryAttrs map[string][]string) []map[string]string {
	if len(attributes) == 0 {
		return []map[string]string{make(map[string]string)}
	}

	// For realistic product catalogs, limit the number of combinations
	// For example, not every color needs to be available in every size
	maxCombinations := 6

	var combinations []map[string]string

	// Start with first attribute
	firstAttr := attributes[0]
	options := categoryAttrs[firstAttr]

	// Limit options to keep combinations reasonable
	if len(options) > 4 {
		// Shuffle and take a subset
		gofakeit.ShuffleAnySlice(options)
		options = options[:min(4, len(options))]
	}

	for _, option := range options {
		combo := map[string]string{firstAttr: option}
		combinations = append(combinations, combo)
	}

	// Add second attribute if available
	if len(attributes) >= 2 {
		secondAttr := attributes[1]
		options = categoryAttrs[secondAttr]

		// Limit options
		if len(options) > 3 {
			gofakeit.ShuffleAnySlice(options)
			options = options[:min(3, len(options))]
		}

		var newCombinations []map[string]string
		for _, combo := range combinations {
			for _, option := range options {
				newCombo := copyMap(combo)
				newCombo[secondAttr] = option
				newCombinations = append(newCombinations, newCombo)
				if len(newCombinations) >= maxCombinations {
					break
				}
			}
			if len(newCombinations) >= maxCombinations {
				break
			}
		}
		combinations = newCombinations
	}

	return combinations
}

// Helper function to copy a map
func copyMap(m map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v
	}
	return result
}

func stringInSlice(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func seedDelivererProfiles(ctx context.Context, db *pgxpool.Pool, delivererUserIDs []int64, adminDivisions []Province) {
	for _, userID := range delivererUserIDs {
		// Tạo thông tin ngẫu nhiên cho người giao hàng
		idCard := fmt.Sprintf("%09d", gofakeit.Number(100000000, 999999999))
		vehicleTypes := []string{"Xe máy", "Ô tô", "Xe đạp"}
		vehicleType := vehicleTypes[gofakeit.Number(0, len(vehicleTypes)-1)]
		licensePlate := fmt.Sprintf("%02d-%s%d",
			gofakeit.Number(10, 99),
			string([]rune("ABCDEFGHKLMNPRSTUVXYZ")[gofakeit.Number(0, 20)]),
			gofakeit.Number(10000, 99999))

		// Tạo hồ sơ người giao hàng
		var delivererID int64
		err := db.QueryRow(ctx, `
			INSERT INTO delivery_persons (
				user_id, id_card_number, vehicle_type, vehicle_license_plate, status, average_rating
			)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT DO NOTHING
			RETURNING id;
		`, userID, idCard, vehicleType, licensePlate, "active", gofakeit.Float32Range(4.0, 5.0)).Scan(&delivererID)

		if err != nil {
			if err == pgx.ErrNoRows {
				// Có thể do conflict, thử lấy ID từ bảng
				err = db.QueryRow(ctx, `
					SELECT id FROM delivery_persons WHERE user_id = $1
				`, userID).Scan(&delivererID)

				if err != nil {
					log.Printf("Error getting deliverer ID: %v", err)
					continue
				}
			} else {
				log.Printf("Error inserting delivery person: %v", err)
				continue
			}
		}

		// Thêm các khu vực phục vụ (2-5 khu vực)
		numAreas := gofakeit.Number(2, 5)

		// Lấy danh sách areas từ cơ sở dữ liệu
		var areaIDs []int64
		rows, err := db.Query(ctx, `SELECT id FROM areas LIMIT 100`)
		if err != nil {
			log.Printf("Error getting areas: %v", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				var areaID int64
				if err := rows.Scan(&areaID); err != nil {
					log.Printf("Error scanning area ID: %v", err)
					continue
				}
				areaIDs = append(areaIDs, areaID)
			}
		}

		if len(areaIDs) == 0 {
			log.Printf("No areas found for deliverer: %d", delivererID)
			continue
		}

		// Shuffle areaIDs
		gofakeit.ShuffleAnySlice(areaIDs)

		// Lấy numAreas khu vực hoặc tất cả nếu không đủ
		if numAreas > len(areaIDs) {
			numAreas = len(areaIDs)
		}

		// Thêm khu vực phục vụ
		for i := 0; i < numAreas; i++ {
			_, err := db.Exec(ctx, `
				INSERT INTO delivery_service_areas (delivery_person_id, area_id, is_active)
				VALUES ($1, $2, $3)
				ON CONFLICT DO NOTHING;
			`, delivererID, areaIDs[i], true)

			if err != nil {
				log.Printf("Error inserting delivery service area: %v", err)
			}
		}

		// Cũng tạo một đơn đăng ký
		_, err = db.Exec(ctx, `
			INSERT INTO delivery_person_applications (
				user_id, id_card_number, id_card_front_image, id_card_back_image,
				vehicle_type, vehicle_license_plate, service_area, application_status
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT DO NOTHING;
		`,
			userID, idCard,
			"https://images.unsplash.com/photo-1633332755192-727a05c4013d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=880&q=80",
			"https://images.unsplash.com/photo-1560250097-0b93528c311a?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
			vehicleType, licensePlate,
			json.RawMessage(fmt.Sprintf(`[{"area_id": %d}]`, areaIDs[0])),
			"approved")

		if err != nil {
			log.Printf("Error inserting delivery person application: %v", err)
		}
	}

	log.Printf("✅ Created %d deliverer profiles", len(delivererUserIDs))
}

// Các function mới cho phần seed mở rộng
func seedProductReviews(ctx context.Context, db *pgxpool.Pool, userIDs []int64) {
	// Lấy danh sách products
	var products []struct {
		id     string
		userID int64
	}

	rows, err := db.Query(ctx, `
		SELECT p.id, s.user_id 
		FROM products p
		JOIN supplier_profiles s ON p.supplier_id = s.id
		LIMIT 100
	`)

	if err != nil {
		log.Printf("Error getting products: %v", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var product struct {
			id     string
			userID int64
		}

		if err := rows.Scan(&product.id, &product.userID); err != nil {
			log.Printf("Error scanning product: %v", err)
			continue
		}

		products = append(products, product)
	}

	if len(products) == 0 {
		log.Printf("No products found for reviews")
		return
	}

	// Shuffle user IDs
	gofakeit.ShuffleAnySlice(userIDs)

	// Tạo reviews cho mỗi sản phẩm
	reviewCount := 0

	for _, product := range products {
		// Số lượng review cho mỗi sản phẩm (3-10)
		numReviews := gofakeit.Number(3, 10)

		for i := 0; i < numReviews; i++ {
			// Lấy user ID ngẫu nhiên, nhưng không phải supplier của sản phẩm
			userIndex := (i*7 + reviewCount) % len(userIDs)
			userID := userIDs[userIndex]

			if userID == product.userID {
				continue // Skip if user is the supplier
			}

			// Rating từ 3-5 sao (hầu hết đều là tốt)
			rating := gofakeit.Number(3, 5)

			// Comment cho review
			reviews := []string{
				"Sản phẩm rất tốt, đúng như mô tả.",
				"Giao hàng nhanh, đóng gói cẩn thận.",
				"Chất lượng sản phẩm tuyệt vời, sẽ ủng hộ shop lần sau.",
				"Hàng đẹp, chất lượng ổn, giá cả hợp lý.",
				"Rất hài lòng với sản phẩm này.",
				"Dịch vụ chăm sóc khách hàng tốt, sản phẩm đúng như hình.",
				"Đóng gói cẩn thận, sản phẩm không bị hư hỏng.",
				"Sản phẩm đúng như mô tả, mẫu mã đẹp.",
				"Shop tư vấn nhiệt tình, giao hàng đúng hẹn.",
				"Giá cả phải chăng, chất lượng tốt.",
			}

			comment := reviews[gofakeit.Number(0, len(reviews)-1)]

			if rating < 5 {
				// Thêm một số phàn nàn nhỏ cho rating dưới 5 sao
				complaints := []string{
					" Tuy nhiên, thời gian giao hàng hơi lâu.",
					" Nhưng đóng gói có thể cẩn thận hơn.",
					" Chỉ tiếc là màu sắc không đúng như hình.",
					" Có một vài chi tiết nhỏ chưa hoàn thiện.",
					" Nhưng giá hơi cao so với chất lượng.",
				}

				comment += complaints[gofakeit.Number(0, len(complaints)-1)]
			}

			// Thêm review
			_, err := db.Exec(ctx, `
				INSERT INTO product_reviews (
					product_id, user_id, rating, comment, helpful_votes
				)
				VALUES ($1, $2, $3, $4, $5)
				ON CONFLICT DO NOTHING;
			`, product.id, userID, rating, comment, gofakeit.Number(0, 20))

			if err != nil {
				log.Printf("Error inserting product review: %v", err)
			} else {
				reviewCount++
			}
		}
	}

	log.Printf("✅ Created %d product reviews", reviewCount)
}

func seedCoupons(ctx context.Context, db *pgxpool.Pool) {
	currentTime := time.Now()

	coupons := []struct {
		code, name, desc, discountType       string
		discountValue, maxDiscount, minOrder float32
		startDate, endDate                   time.Time
		usageLimit                           int
	}{
		{
			code:          "WELCOME10",
			name:          "Chào mừng thành viên mới",
			desc:          "Giảm 10% cho đơn hàng đầu tiên",
			discountType:  "percentage",
			discountValue: 10,
			maxDiscount:   100000,
			minOrder:      0,
			startDate:     currentTime.AddDate(0, -1, 0),
			endDate:       currentTime.AddDate(0, 2, 0),
			usageLimit:    1000,
		},
		{
			code:          "SUMMER2023",
			name:          "Khuyến mãi mùa hè 2023",
			desc:          "Giảm 50.000đ cho đơn hàng từ 500.000đ",
			discountType:  "fixed_amount",
			discountValue: 50000,
			maxDiscount:   50000,
			minOrder:      500000,
			startDate:     currentTime.AddDate(0, 0, -15),
			endDate:       currentTime.AddDate(0, 1, 15),
			usageLimit:    500,
		},
		{
			code:          "FREESHIP",
			name:          "Miễn phí vận chuyển",
			desc:          "Miễn phí vận chuyển cho đơn hàng từ 300.000đ",
			discountType:  "fixed_amount",
			discountValue: 30000,
			maxDiscount:   30000,
			minOrder:      300000,
			startDate:     currentTime.AddDate(0, 0, -30),
			endDate:       currentTime.AddDate(0, 3, 0),
			usageLimit:    2000,
		},
		{
			code:          "TECH15",
			name:          "Giảm giá thiết bị công nghệ",
			desc:          "Giảm 15% cho các sản phẩm điện tử",
			discountType:  "percentage",
			discountValue: 15,
			maxDiscount:   200000,
			minOrder:      1000000,
			startDate:     currentTime.AddDate(0, 0, -5),
			endDate:       currentTime.AddDate(0, 1, 0),
			usageLimit:    300,
		},
		{
			code:          "FLASH50",
			name:          "Flash Sale",
			desc:          "Giảm 50% cho 50 đơn hàng đầu tiên",
			discountType:  "percentage",
			discountValue: 50,
			maxDiscount:   500000,
			minOrder:      100000,
			startDate:     currentTime,
			endDate:       currentTime.AddDate(0, 0, 2),
			usageLimit:    50,
		},
	}

	for _, coupon := range coupons {
		_, err := db.Exec(ctx, `
			INSERT INTO coupons (
				code, name, description, discount_type, discount_value,
				maximum_discount_amount, minimum_order_amount, currency,
				start_date, end_date, usage_limit, is_active
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'VND', $8, $9, $10, TRUE)
			ON CONFLICT (code) DO UPDATE
			SET name = $2, description = $3, discount_type = $4, discount_value = $5,
				maximum_discount_amount = $6, minimum_order_amount = $7,
				start_date = $8, end_date = $9, usage_limit = $10;
		`,
			coupon.code, coupon.name, coupon.desc, coupon.discountType, coupon.discountValue,
			coupon.maxDiscount, coupon.minOrder, coupon.startDate, coupon.endDate,
			coupon.usageLimit)

		if err != nil {
			log.Printf("Error inserting coupon: %v", err)
		}
	}

	log.Printf("✅ Created %d coupons", len(coupons))
}

func seedCartItems(ctx context.Context, db *pgxpool.Pool, userIDs []int64) {
	// Lấy danh sách carts
	carts := make(map[int64]int64) // user_id -> cart_id
	rows, err := db.Query(ctx, `SELECT user_id, id FROM carts LIMIT 1000`)
	if err != nil {
		log.Printf("Error getting carts: %v", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var userID, cartID int64
		if err := rows.Scan(&userID, &cartID); err != nil {
			log.Printf("Error scanning cart: %v", err)
			continue
		}
		carts[userID] = cartID
	}

	// Lấy danh sách product_variants
	type productVariant struct {
		id        string
		productID string
	}

	var variants []productVariant
	variantRows, err := db.Query(ctx, `
		SELECT id, product_id FROM product_variants
		WHERE is_active = TRUE AND inventory_quantity > 0
		LIMIT 500
	`)

	if err != nil {
		log.Printf("Error getting product variants: %v", err)
		return
	}

	defer variantRows.Close()
	for variantRows.Next() {
		var variant productVariant
		if err := variantRows.Scan(&variant.id, &variant.productID); err != nil {
			log.Printf("Error scanning product variant: %v", err)
			continue
		}
		variants = append(variants, variant)
	}

	if len(variants) == 0 {
		log.Printf("No product variants found for cart items")
		return
	}

	// Thêm cart items cho khoảng 30% users
	numUsers := len(userIDs) * 3 / 10
	if numUsers > len(userIDs) {
		numUsers = len(userIDs)
	}

	// Shuffle user IDs
	gofakeit.ShuffleAnySlice(userIDs)
	selectedUsers := userIDs[:numUsers]

	// Tổng số cart items đã tạo
	totalCartItems := 0

	for _, userID := range selectedUsers {
		cartID, ok := carts[userID]
		if !ok {
			continue
		}

		// Số lượng sản phẩm trong giỏ hàng (1-5)
		numItems := gofakeit.Number(1, 5)

		// Shuffle variants
		gofakeit.ShuffleAnySlice(variants)

		// Thêm sản phẩm vào giỏ hàng
		for i := 0; i < numItems && i < len(variants); i++ {
			variant := variants[i]
			quantity := gofakeit.Number(1, 3)

			_, err := db.Exec(ctx, `
				INSERT INTO cart_items (cart_id, product_id, product_variant_id, quantity)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT DO NOTHING;
			`, cartID, variant.productID, variant.id, quantity)

			if err != nil {
				log.Printf("Error inserting cart item: %v", err)
			} else {
				totalCartItems++
			}
		}
	}

	log.Printf("✅ Created %d cart items for %d users", totalCartItems, len(selectedUsers))
}

func seedOrders(ctx context.Context, db *pgxpool.Pool, userIDs []int64) {
	// Lấy danh sách product_variants
	type productVariant struct {
		id            string
		productID     string
		sku           string
		name          string
		price         float32
		discountPrice float32
		imageURL      string
	}

	var variants []productVariant
	variantRows, err := db.Query(ctx, `
		SELECT pv.id, pv.product_id, pv.sku, pv.variant_name, 
			pv.price, COALESCE(pv.discount_price, pv.price), pv.image_url
		FROM product_variants pv
		WHERE pv.is_active = TRUE
		LIMIT 500
	`)

	if err != nil {
		log.Printf("Error getting product variants: %v", err)
		return
	}

	defer variantRows.Close()
	for variantRows.Next() {
		var variant productVariant
		if err := variantRows.Scan(
			&variant.id, &variant.productID, &variant.sku, &variant.name,
			&variant.price, &variant.discountPrice, &variant.imageURL,
		); err != nil {
			log.Printf("Error scanning product variant: %v", err)
			continue
		}
		variants = append(variants, variant)
	}

	if len(variants) == 0 {
		log.Printf("No product variants found for orders")
		return
	}

	// Lấy danh sách products để lấy tên sản phẩm
	productNames := make(map[string]string) // product_id -> name
	productRows, err := db.Query(ctx, `SELECT id, name FROM products`)
	if err != nil {
		log.Printf("Error getting products: %v", err)
	} else {
		defer productRows.Close()
		for productRows.Next() {
			var id, name string
			if err := productRows.Scan(&id, &name); err != nil {
				log.Printf("Error scanning product: %v", err)
				continue
			}
			productNames[id] = name
		}
	}

	// Lấy danh sách địa chỉ của users
	type userAddress struct {
		userID        int64
		recipientName string
		phone         string
		street        string
		district      string
		province      string
		country       string
		postalCode    string
	}

	userAddresses := make(map[int64]userAddress)
	addrRows, err := db.Query(ctx, `
		SELECT user_id, recipient_name, phone, street, district, province, country, postal_code
		FROM addresses
		WHERE is_default = TRUE
	`)

	if err != nil {
		log.Printf("Error getting addresses: %v", err)
	} else {
		defer addrRows.Close()
		for addrRows.Next() {
			var addr userAddress
			if err := addrRows.Scan(
				&addr.userID, &addr.recipientName, &addr.phone, &addr.street,
				&addr.district, &addr.province, &addr.country, &addr.postalCode,
			); err != nil {
				log.Printf("Error scanning address: %v", err)
				continue
			}
			userAddresses[addr.userID] = addr
		}
	}

	// Lấy các phương thức thanh toán
	paymentMethods := make(map[string]int) // code -> id
	pmRows, err := db.Query(ctx, `SELECT id, code FROM payment_methods WHERE is_active = TRUE`)
	if err != nil {
		log.Printf("Error getting payment methods: %v", err)
	} else {
		defer pmRows.Close()
		for pmRows.Next() {
			var id int
			var code string
			if err := pmRows.Scan(&id, &code); err != nil {
				log.Printf("Error scanning payment method: %v", err)
				continue
			}
			paymentMethods[code] = id
		}
	}

	// Lấy danh sách deliverer
	var delivererIDs []int64
	delivererRows, err := db.Query(ctx, `
		SELECT id FROM delivery_persons 
		WHERE status = 'active'
		LIMIT 100
	`)

	if err != nil {
		log.Printf("Error getting deliverers: %v", err)
	} else {
		defer delivererRows.Close()
		for delivererRows.Next() {
			var id int64
			if err := delivererRows.Scan(&id); err != nil {
				log.Printf("Error scanning deliverer: %v", err)
				continue
			}
			delivererIDs = append(delivererIDs, id)
		}
	}

	// Các trạng thái đơn hàng
	orderStatuses := []string{
		"pending", "confirmed", "processing", "shipped", "delivered", "cancelled",
	}

	// Các lý do hủy đơn
	cancelReasons := []string{
		"Khách hàng thay đổi ý định",
		"Khách hàng không liên lạc được",
		"Không đủ hàng",
		"Sản phẩm bị lỗi",
		"Khách hàng đặt nhầm sản phẩm",
	}

	// Các phương thức vận chuyển
	shippingMethods := []string{
		"Standard", "Express", "Same Day",
	}

	// Tạo orders cho khoảng 50% users
	numUsers := len(userIDs) * 5 / 10
	if numUsers > len(userIDs) {
		numUsers = len(userIDs)
	}

	// Shuffle user IDs
	gofakeit.ShuffleAnySlice(userIDs)
	selectedUsers := userIDs[:numUsers]

	// Tổng số đơn hàng đã tạo
	totalOrders := 0

	for _, userID := range selectedUsers {
		// Kiểm tra xem có địa chỉ không
		address, ok := userAddresses[userID]
		if !ok {
			continue
		}

		// Mỗi user tạo 1-3 đơn hàng
		numOrders := gofakeit.Number(1, 3)

		for i := 0; i < numOrders; i++ {
			// Tạo tracking number
			trackingNumber := fmt.Sprintf("TRK%s%d",
				strings.ToUpper(uuid.New().String()[:8]),
				time.Now().Unix())

			// Chọn phương thức vận chuyển
			shippingMethod := shippingMethods[gofakeit.Number(0, len(shippingMethods)-1)]

			// Tạo từ 1-5 sản phẩm cho mỗi đơn hàng
			numItems := gofakeit.Number(1, 5)

			// Shuffle variants
			gofakeit.ShuffleAnySlice(variants)
			selectedVariants := variants[:numItems]

			// Tính toán tổng tiền
			var subTotal float32 = 0
			for _, variant := range selectedVariants {
				quantity := gofakeit.Number(1, 3)
				subTotal += variant.discountPrice * float32(quantity)
			}

			// Thuế và phí vận chuyển
			taxAmount := subTotal * 0.1
			shippingFee := float32(30000)
			if subTotal > 500000 {
				shippingFee = 0 // Miễn phí vận chuyển cho đơn hàng lớn
			}

			// Tổng cộng
			totalAmount := subTotal + taxAmount + shippingFee

			// Tạo đơn hàng
			var orderID string
			err := db.QueryRow(ctx, `
				INSERT INTO orders (
					user_id, tracking_number, shipping_address, country, city, district, ward,
					shipping_method, sub_total, tax_amount, total_amount, recipient_name, recipient_phone
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
				RETURNING id;
			`,
				userID, trackingNumber, address.street, address.country, address.province,
				address.district, "", shippingMethod, subTotal, taxAmount, totalAmount,
				address.recipientName, address.phone,
			).Scan(&orderID)

			if err != nil {
				log.Printf("Error inserting order: %v", err)
				continue
			}

			// Tạo các order items
			for _, variant := range selectedVariants {
				quantity := gofakeit.Number(1, 3)
				unitPrice := variant.discountPrice

				// Lấy tên sản phẩm
				productName, ok := productNames[variant.productID]
				if !ok {
					productName = "Sản phẩm không xác định"
				}

				// Trạng thái đơn hàng
				status := orderStatuses[gofakeit.Number(0, len(orderStatuses)-1)]

				// Ghi chú hủy đơn nếu là cancelled
				var cancelledReason interface{} = nil
				if status == "cancelled" {
					cancelledReason = cancelReasons[gofakeit.Number(0, len(cancelReasons)-1)]
				}

				// Ước tính ngày giao hàng
				estimatedDelivery := time.Now().AddDate(0, 0, gofakeit.Number(3, 7))

				// Ngày giao hàng thực tế (nếu đã giao)
				var actualDelivery interface{} = nil
				if status == "delivered" {
					actualDelivery = time.Now().AddDate(0, 0, gofakeit.Number(1, 5))
				}

				// Tạo order item
				var orderItemID string
				err := db.QueryRow(ctx, `
					INSERT INTO order_items (
						order_id, product_name, product_sku, product_variant_image_url,
						product_variant_name, quantity, unit_price, total_price,
						estimated_delivery_date, actual_delivery_date, cancelled_reason,
						status, shipping_fee
					)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
					RETURNING id;
				`,
					orderID, productName, variant.sku, variant.imageURL,
					variant.name, quantity, unitPrice, unitPrice*float32(quantity),
					estimatedDelivery, actualDelivery, cancelledReason,
					status, shippingFee/float32(len(selectedVariants)),
				).Scan(&orderItemID)

				if err != nil {
					log.Printf("Error inserting order item: %v", err)
					continue
				}

				// Ghi lịch sử trạng thái đơn hàng
				_, err = db.Exec(ctx, `
					INSERT INTO order_items_history (order_item_id, status, notes, created_by)
					VALUES ($1, $2, $3, $4);
				`, orderItemID, status, "Cập nhật trạng thái đơn hàng", userID)

				if err != nil {
					log.Printf("Error inserting order item history: %v", err)
				}

				// Tạo người giao hàng cho đơn hàng đã shipped hoặc delivered
				if (status == "shipped" || status == "delivered") && len(delivererIDs) > 0 {
					delivererID := delivererIDs[gofakeit.Number(0, len(delivererIDs)-1)]

					deliveryStatus := "assigned"
					if status == "shipped" {
						deliveryStatus = "in_transit"
					} else if status == "delivered" {
						deliveryStatus = "delivered"
					}

					pickupTime := time.Now().AddDate(0, 0, -gofakeit.Number(1, 3))
					var deliveryTime interface{} = nil
					if status == "delivered" {
						deliveryTime = time.Now().AddDate(0, 0, -gofakeit.Number(0, 2))
					}

					// Tạo order_deliverer
					_, err = db.Exec(ctx, `
						INSERT INTO order_deliverers (
							order_item_id, deliverer_id, status, pickup_time, delivery_time
						)
						VALUES ($1, $2, $3, $4, $5)
						ON CONFLICT DO NOTHING;
					`, orderItemID, delivererID, deliveryStatus, pickupTime, deliveryTime)

					if err != nil {
						log.Printf("Error inserting order deliverer: %v", err)
					}
				}

				// Tạo payment history cho đơn hàng
				if status != "cancelled" && len(paymentMethods) > 0 {
					// Chọn ngẫu nhiên phương thức thanh toán
					var paymentMethodID int
					if len(paymentMethods) > 0 {
						methods := []string{"cod", "momo"}
						code := methods[gofakeit.Number(0, len(methods)-1)]
						paymentMethodID = paymentMethods[code]
					} else {
						// Fallback nếu không có payment method
						paymentMethodID = 1
					}

					// Thêm user_payment_method nếu chưa có
					var userPaymentMethodID int64
					err := db.QueryRow(ctx, `
						INSERT INTO user_payment_methods (
							user_id, payment_method_id, is_default, card_holder_name, card_number
						)
						VALUES ($1, $2, $3, $4, $5)
						ON CONFLICT (user_id, payment_method_id) 
						DO UPDATE SET is_default = $3
						RETURNING id;
					`, userID, paymentMethodID, true, address.recipientName, "XXXX-XXXX-XXXX-XXXX").Scan(&userPaymentMethodID)

					if err != nil {
						log.Printf("Error inserting user payment method: %v", err)
						continue
					}

					// Tạo payment history
					paymentStatus := "pending"
					if status == "delivered" {
						paymentStatus = "completed"
					} else if status == "shipped" || status == "processing" {
						paymentStatus = "processing"
					}

					var paidAt interface{} = nil
					if paymentStatus == "completed" {
						paidAt = time.Now().AddDate(0, 0, -gofakeit.Number(0, 3))
					}

					_, err = db.Exec(ctx, `
						INSERT INTO payment_history (
							order_item_id, user_payment_method_id, amount, status,
							transaction_id, payment_gateway, paid_at
						)
						VALUES ($1, $2, $3, $4, $5, $6, $7)
						ON CONFLICT DO NOTHING;
					`,
						orderItemID, userPaymentMethodID, unitPrice*float32(quantity),
						paymentStatus, uuid.New().String(), "internal", paidAt)

					if err != nil {
						log.Printf("Error inserting payment history: %v", err)
					}
				}
			}

			totalOrders++
		}
	}

	log.Printf("✅ Created %d orders", totalOrders)
}

func seedNotifications(ctx context.Context, db *pgxpool.Pool, userIDs []int64) {
	// Số lượng thông báo mỗi người dùng (tăng từ 5 lên 15)
	numPerUser := 15

	// Các loại thông báo
	types := []int{1, 2, 3, 4, 5} // 1: order, 2: payment, 3: product, 4: promotion, 5: system

	// Tiêu đề và nội dung thông báo
	titles := map[int][]string{
		1: { // Order
			"Đơn hàng đã được xác nhận",
			"Đơn hàng đang được xử lý",
			"Đơn hàng đang được giao",
			"Đơn hàng đã được giao thành công",
			"Cập nhật trạng thái đơn hàng",
			"Thông báo về đơn hàng của bạn",
		},
		2: { // Payment
			"Thanh toán thành công",
			"Thanh toán đang được xử lý",
			"Yêu cầu thanh toán đơn hàng",
			"Hóa đơn mới",
			"Xác nhận thanh toán",
		},
		3: { // Product
			"Sản phẩm đang giảm giá",
			"Sản phẩm bạn quan tâm đã có hàng",
			"Đánh giá sản phẩm đã mua",
			"Sản phẩm mới ra mắt",
			"Cập nhật thông tin sản phẩm",
		},
		4: { // Promotion
			"Khuyến mãi mùa hè",
			"Flash sale cuối tuần",
			"Mã giảm giá cho thành viên",
			"Ưu đãi đặc biệt dành cho bạn",
			"Quà tặng sinh nhật",
			"Ưu đãi độc quyền",
		},
		5: { // System
			"Cập nhật thông tin tài khoản",
			"Xác thực tài khoản thành công",
			"Bảo mật tài khoản",
			"Thay đổi mật khẩu",
			"Cập nhật ứng dụng",
			"Thông báo bảo trì hệ thống",
		},
	}

	contents := map[int][]string{
		1: { // Order
			"Đơn hàng #ORDER-ID của bạn đã được xác nhận. Chúng tôi sẽ sớm xử lý đơn hàng.",
			"Đơn hàng #ORDER-ID của bạn đang được xử lý. Dự kiến đơn hàng sẽ được giao trong 3-5 ngày tới.",
			"Đơn hàng #ORDER-ID của bạn đang được giao. Vui lòng chuẩn bị nhận hàng.",
			"Đơn hàng #ORDER-ID của bạn đã được giao thành công. Cảm ơn bạn đã mua sắm!",
			"Chúng tôi đã cập nhật trạng thái đơn hàng #ORDER-ID của bạn. Vui lòng kiểm tra chi tiết trong tài khoản.",
			"Có thông báo mới về đơn hàng #ORDER-ID. Vui lòng kiểm tra để biết thêm chi tiết.",
		},
		2: { // Payment
			"Thanh toán đơn hàng #ORDER-ID của bạn đã thành công. Cảm ơn bạn!",
			"Thanh toán đơn hàng #ORDER-ID của bạn đang được xử lý. Chúng tôi sẽ thông báo cho bạn khi hoàn tất.",
			"Vui lòng thanh toán đơn hàng #ORDER-ID của bạn trong vòng 24 giờ để tránh bị hủy.",
			"Hóa đơn mới cho đơn hàng #ORDER-ID đã được tạo. Vui lòng thanh toán đúng hạn.",
			"Chúng tôi đã nhận được thanh toán cho đơn hàng #ORDER-ID của bạn. Xác nhận thanh toán đã hoàn tất.",
		},
		3: { // Product
			"Sản phẩm [PRODUCT-NAME] bạn đã xem gần đây đang được giảm giá 20%. Mua ngay!",
			"Sản phẩm [PRODUCT-NAME] bạn quan tâm đã có hàng trở lại. Nhanh tay mua ngay!",
			"Bạn đã mua sản phẩm [PRODUCT-NAME] gần đây. Vui lòng đánh giá sản phẩm để nhận voucher!",
			"Sản phẩm mới [PRODUCT-NAME] vừa ra mắt. Khám phá ngay hôm nay với ưu đãi đặc biệt!",
			"Thông tin về sản phẩm [PRODUCT-NAME] bạn đã mua đã được cập nhật. Kiểm tra ngay!",
		},
		4: { // Promotion
			"Khuyến mãi mùa hè với hàng ngàn sản phẩm giảm giá lên đến 50%. Khám phá ngay!",
			"Flash sale cuối tuần - Giảm giá sốc chỉ trong 2 giờ. Bắt đầu từ 20:00 tối nay.",
			"Tặng bạn mã giảm giá SUMMER10 giảm 10% cho đơn hàng tiếp theo. Hạn sử dụng 7 ngày.",
			"Ưu đãi đặc biệt cho thành viên thân thiết - Giảm 15% cho các sản phẩm thời trang.",
			"Chúc mừng sinh nhật! Tặng bạn voucher giảm 100.000đ cho đơn hàng từ 500.000đ.",
			"Ưu đãi độc quyền dành riêng cho bạn - Mua 1 tặng 1 cho các sản phẩm chăm sóc cá nhân.",
		},
		5: { // System
			"Thông tin tài khoản của bạn đã được cập nhật thành công.",
			"Tài khoản của bạn đã được xác thực thành công. Bạn có thể sử dụng đầy đủ tính năng của hệ thống.",
			"Vì lý do bảo mật, vui lòng cập nhật mật khẩu của bạn định kỳ.",
			"Mật khẩu của bạn đã được thay đổi thành công. Nếu không phải bạn thực hiện, vui lòng liên hệ ngay với chúng tôi.",
			"Phiên bản mới của ứng dụng đã có sẵn. Cập nhật ngay để trải nghiệm những tính năng mới!",
			"Hệ thống sẽ tiến hành bảo trì từ 23:00 đến 05:00 ngày mai. Mong bạn thông cảm cho sự bất tiện này.",
		},
	}

	// Hình ảnh cho thông báo
	imageURLs := []string{
		"https://images.unsplash.com/photo-1555529669-e69e7aa0ba9a?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		"https://images.unsplash.com/photo-1556740758-90de374c12ad?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		"https://images.unsplash.com/photo-1521791136064-7986c2920216?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1169&q=80",
		"https://images.unsplash.com/photo-1511370235399-1802cae1d32f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1955&q=80",
		"https://images.unsplash.com/photo-1633174524827-db00a6b9c7b1?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1992&q=80",
	}

	// Sản phẩm mẫu cho thông báo sản phẩm
	sampleProducts := []string{
		"iPhone 13 Pro Max", "Samsung Galaxy S22", "Laptop Dell XPS 13",
		"Áo thun nam", "Áo sơ mi nữ", "Quần jeans", "Giày thể thao",
		"Nồi cơm điện", "Máy lọc không khí", "Đắc Nhân Tâm", "Tạ tay 5kg",
	}

	// Đảm bảo admin cũng có thông báo (ID = 1)
	var adminID int64 = 1
	var adminExists bool = false

	// Kiểm tra xem admin ID có trong danh sách userIDs không
	for _, id := range userIDs {
		if id == adminID {
			adminExists = true
			break
		}
	}

	// Nếu admin không có trong danh sách, thêm vào
	if !adminExists {
		userIDs = append([]int64{adminID}, userIDs...)
	}

	// Tạo thông báo cho mỗi người dùng
	totalNotifs := 0

	for _, userID := range userIDs {
		// Đảm bảo admin có ít nhất 15 thông báo
		notifCount := numPerUser
		if userID == adminID {
			notifCount = 15 // Hoặc nhiều hơn nếu muốn admin có nhiều thông báo hơn
		}

		for i := 0; i < notifCount; i++ {
			// Chọn loại thông báo ngẫu nhiên
			typeIndex := gofakeit.Number(0, len(types)-1)
			notifType := types[typeIndex]

			// Chọn title và content ngẫu nhiên
			titleIndex := gofakeit.Number(0, len(titles[notifType])-1)
			contentIndex := gofakeit.Number(0, len(contents[notifType])-1)

			title := titles[notifType][titleIndex]
			content := contents[notifType][contentIndex]

			// Thay thế placeholders
			if strings.Contains(content, "ORDER-ID") {
				content = strings.Replace(content, "ORDER-ID", fmt.Sprintf("%d", gofakeit.Number(1000, 9999)), -1)
			}

			if strings.Contains(content, "PRODUCT-NAME") {
				productIndex := gofakeit.Number(0, len(sampleProducts)-1)
				content = strings.Replace(content, "PRODUCT-NAME", sampleProducts[productIndex], -1)
			}

			// Chọn hình ảnh ngẫu nhiên
			imageURL := ""
			if gofakeit.Bool() {
				imageURL = imageURLs[gofakeit.Number(0, len(imageURLs)-1)]
			}

			// Tạo thông báo
			// Đảm bảo một số thông báo đã đọc và một số chưa đọc
			isRead := gofakeit.Bool()

			// Nếu là thông báo cuối cùng, đảm bảo chưa đọc
			if i >= notifCount-5 {
				isRead = false
			}

			_, err := db.Exec(ctx, `
                INSERT INTO notifications (
                    user_id, type, title, content, is_read, image_title
                )
                VALUES ($1, $2, $3, $4, $5, $6)
                ON CONFLICT DO NOTHING;
            `,
				userID, notifType, title, content, isRead, imageURL)

			if err != nil {
				log.Printf("Error inserting notification: %v", err)
			} else {
				totalNotifs++
			}
		}
	}

	log.Printf("✅ Created %d notifications", totalNotifs)
}

// Thêm vào main function
func seedEverything(ctx context.Context, pools map[string]*pgxpool.Pool, userIDs []int64, supplierIDs []int64, adminDivisions []Province) {
	// Gọi các hàm seed bổ sung
	seedProductReviews(ctx, pools["partners_db"], userIDs)
	seedCoupons(ctx, pools["orders_db"])
	seedCartItems(ctx, pools["orders_db"], userIDs)
	seedOrders(ctx, pools["orders_db"], userIDs)
	seedNotifications(ctx, pools["notifications_db"], userIDs)
}
