package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionDetail struct {
	ModuleID    int   `json:"module_id"`
	Permissions []int `json:"permissions"`
}

type progressUpdate struct {
	goroutineID int
	count       int
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

	// Seed independent tables
	seedAPIGatewayIndependentTables(ctx, pools["api_gateway_db"])
	seedOrderIndependentTables(ctx, pools["orders_db"])
	seedPartnersIndependentTables(ctx, pools["partners_db"])

	// Seed users (10,000) and get their IDs
	userIDs := seedUsers(ctx, pools["api_gateway_db"], 10000, 1000, 100)

	// Seed addresses for all users
	seedAddressesForUsers(ctx, pools["api_gateway_db"], userIDs)

	// Seed dependent tables
	seedNotificationPreferences(ctx, pools["notifications_db"], userIDs)
	seedCarts(ctx, pools["orders_db"], userIDs)

	// Select supplier user IDs and assign supplier role
	supplierUserIDs := selectSupplierUserIDs(userIDs)
	assignSupplierRole(ctx, pools["api_gateway_db"], supplierUserIDs)

	// Seed supplier profiles and products
	supplierIDs := seedSupplierProfiles(ctx, pools["api_gateway_db"], pools["partners_db"], supplierUserIDs)
	seedProducts(ctx, pools["partners_db"], supplierIDs)

	fmt.Println("✅ Seed completed successfully")
}

func connectDB(ctx context.Context, dsn, dbName string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to %s: %v", dbName, err)
	}
	return pool
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
		},
		"customer": {
			{ModuleID: 1, Permissions: []int{1, 2, 3, 4}},
			{ModuleID: 4, Permissions: []int{1, 2, 3, 4}},
			{ModuleID: 5, Permissions: []int{1, 4}},
			{ModuleID: 6, Permissions: []int{1, 4, 3}},
			{ModuleID: 7, Permissions: []int{4}},
			{ModuleID: 8, Permissions: []int{1, 4, 3}},
		},
		"supplier": {
			{ModuleID: 3, Permissions: []int{1, 2, 3, 4}},
			{ModuleID: 9, Permissions: []int{1, 2, 3, 4}},
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
		INSERT INTO users (fullname, email, avatar_url, email_verified, status, phone_verified) 
		VALUES ('Admin User', 'admin@admin.com', 'https://ui-avatars.com/api/?name=Admin+User', TRUE, 'active', TRUE) 
		ON CONFLICT (email) DO UPDATE 
		SET fullname = 'Admin User', 
		    avatar_url = 'https://ui-avatars.com/api/?name=Admin+User', 
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
					avatar := fmt.Sprintf("https://ui-avatars.com/api/?name=%s", name)
					birth := gofakeit.DateRange(
						time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
					)
					hash, _ := utils.HashPassword("123456")

					phone := gofakeit.Phone()
					if len(phone) > 15 {
						phone = phone[:15]
					}
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
		query += strings.Join(valueStrings, ",") + ";"
		_, err := db.Exec(ctx, query, args...)
		if err != nil {
			log.Fatal("assign supplier role:", err)
		}
	}
}

func selectSupplierUserIDs(userIDs []int64) []int64 {
	count := len(userIDs) / 10 // 10% users as suppliers
	supplierUserIDs := make([]int64, 0, count)
	gofakeit.ShuffleAnySlice(userIDs)
	for i := 0; i < count && i < len(userIDs); i++ {
		supplierUserIDs = append(supplierUserIDs, userIDs[i])
	}
	return supplierUserIDs
}

func seedAddressesForUsers(ctx context.Context, db *pgxpool.Pool, userIDs []int64) {
	var homeAddressTypeID int64
	err := db.QueryRow(ctx, `SELECT id FROM address_types WHERE address_type = $1`, "Home").Scan(&homeAddressTypeID)
	if err != nil {
		log.Fatal("get home address type:", err)
	}

	// SỬA: Mở rộng danh sách provinces và districts để tăng tính đa dạng
	provinces := []string{"Hà Nội", "TP Hồ Chí Minh", "Đà Nẵng", "Hải Phòng", "Cần Thơ", "Nha Trang", "Huế"}
	districts := []string{"Ba Đình", "Quận 1", "Hải Châu", "Hồng Bàng", "Ninh Kiều", "Khánh Hòa", "Thừa Thiên"}

	batchSize := 1000
	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}
		batch := userIDs[i:end]

		var args []interface{}
		query := `INSERT INTO addresses (user_id, recipient_name, phone, street, district, province, postal_code, country, is_default, address_type_id) VALUES `
		valueStrings := make([]string, len(batch))
		for j, userID := range batch {
			idx := j * 10
			recipientName := gofakeit.Name()
			phone := gofakeit.Phone()
			if len(phone) > 20 {
				phone = phone[:20]
			}
			street := gofakeit.Street()
			district := districts[gofakeit.Number(0, len(districts)-1)]
			province := provinces[gofakeit.Number(0, len(provinces)-1)]
			postalCode := gofakeit.Zip()
			country := "Việt Nam"
			isDefault := gofakeit.Bool()

			valueStrings[j] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7, idx+8, idx+9, idx+10)
			args = append(args, userID, recipientName, phone, street, district, province, postalCode, country, isDefault, homeAddressTypeID)
		}
		query += strings.Join(valueStrings, ",") + ";"
		_, err := db.Exec(ctx, query, args...)
		if err != nil {
			log.Fatal("insert addresses:", err)
		}
	}
}

// Notification Service Seeding
func seedNotificationPreferences(ctx context.Context, db *pgxpool.Pool, userIDs []int64) {
	prefs := map[string]bool{
		"order_status":   true,
		"payment_status": true,
		"product_status": true,
		"promotion":      true,
	}
	prefsJSON, _ := json.Marshal(prefs)

	batchSize := 1000
	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}
		batch := userIDs[i:end]

		var args []interface{}
		query := `INSERT INTO notification_preferences (user_id, email_preferences, in_app_preferences) VALUES `
		valueStrings := make([]string, len(batch))
		for j, userID := range batch {
			idx := j * 3
			valueStrings[j] = fmt.Sprintf("($%d, $%d, $%d)", idx+1, idx+2, idx+3)
			args = append(args, userID, prefsJSON, prefsJSON)
		}
		query += strings.Join(valueStrings, ",") + ";"
		_, err := db.Exec(ctx, query, args...)
		if err != nil {
			log.Fatal("insert notification_preferences:", err)
		}
	}
}

// Order Service Seeding
func seedOrderIndependentTables(ctx context.Context, db *pgxpool.Pool) {
	seedAreas(ctx, db)
	seedPaymentMethods(ctx, db)
}

func seedAreas(ctx context.Context, db *pgxpool.Pool) {
	areas := []struct {
		city, country, district, ward, areaCode string
	}{
		{"Hà Nội", "Việt Nam", "Ba Đình", "Phúc Xá", "HN01"},
		{"TP Hồ Chí Minh", "Việt Nam", "Quận 1", "Bến Nghé", "HCM01"},
		{"Đà Nẵng", "Việt Nam", "Hải Châu", "Hải Châu I", "DN01"},
	}
	for _, a := range areas {
		_, _ = db.Exec(ctx, `
			INSERT INTO areas (city, country, district, ward, area_code) 
			VALUES ($1, $2, $3, $4, $5) 
			ON CONFLICT DO NOTHING;
		`, a.city, a.country, a.district, a.ward, a.areaCode)
	}
}

func seedPaymentMethods(ctx context.Context, db *pgxpool.Pool) {
	methods := []struct {
		name, code string
	}{
		{"Thẻ tín dụng", "CREDIT_CARD"},
		{"Chuyển khoản ngân hàng", "BANK_TRANSFER"},
		{"Thanh toán khi nhận hàng", "COD"},
		{"Momo", "MOMO"},
	}
	for _, m := range methods {
		_, _ = db.Exec(ctx, `
			INSERT INTO payment_methods (name, code) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING;
		`, m.name, m.code)
	}
}

func seedCarts(ctx context.Context, db *pgxpool.Pool, userIDs []int64) {
	batchSize := 1000
	for i := 0; i < len(userIDs); i += batchSize {
		end := i + batchSize
		if end > len(userIDs) {
			end = len(userIDs)
		}
		batch := userIDs[i:end]

		var args []interface{}
		query := `INSERT INTO carts (user_id) VALUES `
		valueStrings := make([]string, len(batch))
		for j, userID := range batch {
			idx := j * 1
			valueStrings[j] = fmt.Sprintf("($%d)", idx+1)
			args = append(args, userID)
		}
		query += strings.Join(valueStrings, ",") + ";"
		_, err := db.Exec(ctx, query, args...)
		if err != nil {
			log.Fatal("insert carts:", err)
		}
	}
}

// Partners Service Seeding
func seedPartnersIndependentTables(ctx context.Context, db *pgxpool.Pool) {
	seedCategories(ctx, db)
	seedTags(ctx, db)
	seedAttributeDefinitions(ctx, db)
	seedAttributeOptions(ctx, db)
}

func seedCategories(ctx context.Context, db *pgxpool.Pool) {
	// Define root categories
	rootCategories := []struct {
		name, imageURL string
	}{
		{"Điện tử", "https://picsum.photos/seed/electronics/150/150"},
		{"Thời trang", "https://picsum.photos/seed/fashion/150/150"},
		{"Đồ gia dụng", "https://picsum.photos/seed/home_appliances/150/150"},
		{"Sách", "https://picsum.photos/seed/books/150/150"},
		{"Thể thao", "https://picsum.photos/seed/sports/150/150"},
	}

	// Map to store root category IDs
	rootCategoryIDs := make(map[string]int64)

	// Insert root categories
	for _, rc := range rootCategories {
		var rootID int64
		rows, err := db.Query(ctx, `INSERT INTO categories (name, image_url, is_active) VALUES ($1, $2, TRUE) RETURNING id;`, rc.name, rc.imageURL)
		if err != nil {
			log.Fatalf("insert root category %s: %v", rc.name, err)
		}
		for rows.Next() {
			if err := rows.Scan(&rootID); err != nil {
				log.Fatal("scan root category id:", err)
			}
		}
		rows.Close()
		rootCategoryIDs[rc.name] = rootID
	}

	// Define subcategories
	subCategories := []struct {
		name, imageURL, parentName string
		parentID                   *int64
	}{
		// Điện tử
		{"Điện thoại thông minh", "https://picsum.photos/seed/smartphone/150/150", "Điện tử", nil},
		{"Máy tính xách tay", "https://picsum.photos/seed/laptop/150/150", "Điện tử", nil},
		// Thời trang
		{"Thời trang nam", "https://picsum.photos/seed/mens_fashion/150/150", "Thời trang", nil},
		{"Thời trang nữ", "https://picsum.photos/seed/womens_fashion/150/150", "Thời trang", nil},
		{"Phụ kiện thời trang", "https://picsum.photos/seed/accessories/150/150", "Thời trang", nil},
		// Đồ gia dụng
		{"Đồ dùng nhà bếp", "https://picsum.photos/seed/kitchen_appliances/150/150", "Đồ gia dụng", nil},
		{"Đồ dùng vệ sinh", "https://picsum.photos/seed/cleaning_appliances/150/150", "Đồ gia dụng", nil},
		// Sách
		{"Sách tiểu thuyết", "https://picsum.photos/seed/fiction_books/150/150", "Sách", nil},
		{"Sách phi tiểu thuyết", "https://picsum.photos/seed/nonfiction_books/150/150", "Sách", nil},
		// Thể thao
		{"Thiết bị thể dục", "https://picsum.photos/seed/fitness_equipment/150/150", "Thể thao", nil},
		{"Thể thao ngoài trời", "https://picsum.photos/seed/outdoor_sports/150/150", "Thể thao", nil},
	}

	// Assign parent IDs to subcategories
	for i := range subCategories {
		if parentID, exists := rootCategoryIDs[subCategories[i].parentName]; exists {
			subCategories[i].parentID = new(int64)
			*subCategories[i].parentID = parentID
		} else {
			log.Fatalf("Parent category not found: %s", subCategories[i].parentName)
		}
	}

	// Insert subcategories
	for _, sc := range subCategories {
		_, err := db.Exec(ctx, `
			INSERT INTO categories (name, image_url, parent_id, is_active) 
			VALUES ($1, $2, $3, TRUE) 
			ON CONFLICT DO NOTHING;
		`, sc.name, sc.imageURL, *sc.parentID)
		if err != nil {
			log.Fatalf("insert subcategory %s: %v", sc.name, err)
		}
	}
}

func seedTags(ctx context.Context, db *pgxpool.Pool) {
	tags := []string{"Công nghệ", "Giảm giá", "Mới", "Thời trang", "Gia dụng", "Thể thao", "Sách"}
	for _, t := range tags {
		_, _ = db.Exec(ctx, `INSERT INTO tags (name) VALUES ($1) ON CONFLICT DO NOTHING;`, t)
	}
}

func seedAttributeDefinitions(ctx context.Context, db *pgxpool.Pool) {
	attrs := []struct {
		name, inputType string
	}{
		{"Màu sắc", "select"},
		{"Kích thước", "select"},
		{"Chất liệu", "select"},
		{"Công suất", "select"},     // For appliances
		{"Chất liệu vải", "select"}, // For clothing
		{"Thể loại", "select"},      // For books
		{"Loại thiết bị", "select"}, // For sports equipment
	}
	for _, a := range attrs {
		_, _ = db.Exec(ctx, `
			INSERT INTO attribute_definitions (name, input_type) 
			VALUES ($1, $2) 
			ON CONFLICT DO NOTHING;
		`, a.name, a.inputType)
	}
}

func seedAttributeOptions(ctx context.Context, db *pgxpool.Pool) {
	options := []struct {
		attrName, value string
	}{
		// Màu sắc
		{"Màu sắc", "Đỏ"},
		{"Màu sắc", "Xanh dương"},
		{"Màu sắc", "Xanh lá"},
		{"Màu sắc", "Đen"},
		{"Màu sắc", "Trắng"},
		{"Màu sắc", "Xám"},
		{"Màu sắc", "Nâu"},
		{"Màu sắc", "Vàng"},
		{"Màu sắc", "Cam"},
		{"Màu sắc", "Tím"},
		// Kích thước
		{"Kích thước", "S"},
		{"Kích thước", "M"},
		{"Kích thước", "L"},
		{"Kích thước", "XL"},
		{"Kích thước", "XXL"},
		// Chất liệu
		{"Chất liệu", "Cotton"},
		{"Chất liệu", "Polyester"},
		{"Chất liệu", "Len"},
		{"Chất liệu", "Da"},
		{"Chất liệu", "Vải lanh"},
		{"Chất liệu", "Kim loại"},
		{"Chất liệu", "Nhựa"},
		// Công suất (for appliances)
		{"Công suất", "500W"},
		{"Công suất", "1000W"},
		{"Công suất", "1500W"},
		{"Công suất", "2000W"},
		// Chất liệu vải (for clothing)
		{"Chất liệu vải", "Cotton"},
		{"Chất liệu vải", "Polyester"},
		{"Chất liệu vải", "Len"},
		{"Chất liệu vải", "Lụa"},
		// Thể loại (for books)
		{"Thể loại", "Tiểu thuyết"},
		{"Thể loại", "Khoa học viễn tưởng"},
		{"Thể loại", "Lịch sử"},
		{"Thể loại", "Tự truyện"},
		// Loại thiết bị (for sports equipment)
		{"Loại thiết bị", "Máy chạy bộ"},
		{"Loại thiết bị", "Xe đạp tập"},
		{"Loại thiết bị", "Tạ tay"},
		{"Loại thiết bị", "Bóng rổ"},
	}

	attrIDs := make(map[string]int)
	rows, err := db.Query(ctx, `SELECT id, name FROM attribute_definitions`)
	if err != nil {
		log.Fatal("get attribute definitions:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("scan attribute definition:", err)
		}
		attrIDs[name] = id
	}

	for _, opt := range options {
		attrID, exists := attrIDs[opt.attrName]
		if !exists {
			log.Fatalf("Attribute not found: %s", opt.attrName)
		}
		_, err := db.Exec(ctx, `
            INSERT INTO attribute_options (attribute_definition_id, option_value) 
            VALUES ($1, $2) 
            ON CONFLICT DO NOTHING;
        `, attrID, opt.value)
		if err != nil {
			log.Fatal("insert attribute option:", err)
		}
	}
}

func seedSupplierProfiles(ctx context.Context, apiDB, partnerDB *pgxpool.Pool, supplierUserIDs []int64) []int64 {
	batchSize := 100
	supplierIDs := make([]int64, 0, len(supplierUserIDs))

	for i := 0; i < len(supplierUserIDs); i += batchSize {
		end := i + batchSize
		if end > len(supplierUserIDs) {
			end = len(supplierUserIDs)
		}
		batch := supplierUserIDs[i:end]

		var args []interface{}
		query := `INSERT INTO supplier_profiles (user_id, company_name, contact_phone, logo_url, business_address_id, tax_id, status) VALUES `
		valueStrings := make([]string, len(batch))

		for j, userID := range batch {
			// Get address_id for the user from api_gateway_db
			var addressID int64
			err := apiDB.QueryRow(ctx, `SELECT id FROM addresses WHERE user_id = $1 LIMIT 1`, userID).Scan(&addressID)
			if err != nil {
				log.Fatal("get address for user:", err)
			}

			idx := j * 7
			companyName := gofakeit.Company()
			contactPhone := gofakeit.Phone()
			if len(contactPhone) > 20 {
				contactPhone = contactPhone[:20]
			}
			logoURL := fmt.Sprintf("https://picsum.photos/seed/%s/150/150", companyName)
			taxID := fmt.Sprintf("TAX-%d", gofakeit.Number(1000, 9999))
			status := "active"

			valueStrings[j] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7)
			args = append(args, userID, companyName, contactPhone, logoURL, addressID, taxID, status)
		}
		query += strings.Join(valueStrings, ",") + " RETURNING id;"

		rows, err := partnerDB.Query(ctx, query, args...)
		if err != nil {
			log.Fatal("insert supplier profiles:", err)
		}
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				log.Fatal("scan supplier id:", err)
			}
			supplierIDs = append(supplierIDs, id)
		}
		rows.Close()
	}

	return supplierIDs
}

func seedProducts(ctx context.Context, db *pgxpool.Pool, supplierIDs []int64) {
	// Get all subcategory IDs and names
	subCategories := make(map[int64]string) // Map of category ID to category name
	rows, err := db.Query(ctx, `SELECT id, name FROM categories WHERE parent_id IS NOT NULL`)
	if err != nil {
		log.Fatal("get subcategories:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("scan subcategory:", err)
		}
		subCategories[id] = name
	}

	var subCategoryIDs []int64
	for id := range subCategories {
		subCategoryIDs = append(subCategoryIDs, id)
	}
	if len(subCategoryIDs) == 0 {
		log.Fatal("no subcategories found")
	}

	// Get tag IDs
	tagIDs := make(map[string]uuid.UUID)
	rows, err = db.Query(ctx, `SELECT id, name FROM tags`)
	if err != nil {
		log.Fatal("get tags:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id uuid.UUID
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("scan tag:", err)
		}
		tagIDs[name] = id
	}

	// Get attribute option IDs
	attrOptionIDs := make(map[string][]int)
	attrOptionToDefIDs := make(map[int]int)
	rows, err = db.Query(ctx, `
        SELECT ao.id, ao.option_value, ad.name, ao.attribute_definition_id 
        FROM attribute_options ao 
        JOIN attribute_definitions ad ON ao.attribute_definition_id = ad.id
    `)
	if err != nil {
		log.Fatal("get attribute options:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, defID int
		var value, attrName string
		if err := rows.Scan(&id, &value, &attrName, &defID); err != nil {
			log.Fatal("scan attribute option:", err)
		}
		attrOptionIDs[attrName] = append(attrOptionIDs[attrName], id)
		attrOptionToDefIDs[id] = defID
	}

	// Map categories to relevant attributes
	categoryAttributes := map[string][]struct {
		attrName string
		count    int
	}{
		"Điện thoại thông minh": {
			{"Màu sắc", 2},
			{"Kích thước", 1},
			{"Chất liệu", 1},
		},
		"Máy tính xách tay": {
			{"Màu sắc", 2},
			{"Kích thước", 1},
			{"Chất liệu", 1},
		},
		"Thời trang nam": {
			{"Màu sắc", 2},
			{"Kích thước", 1},
			{"Chất liệu vải", 1},
		},
		"Thời trang nữ": {
			{"Màu sắc", 2},
			{"Kích thước", 1},
			{"Chất liệu vải", 1},
		},
		"Phụ kiện thời trang": {
			{"Màu sắc", 2},
			{"Chất liệu", 1},
		},
		"Đồ dùng nhà bếp": {
			{"Màu sắc", 1},
			{"Công suất", 1},
			{"Chất liệu", 1},
		},
		"Đồ dùng vệ sinh": {
			{"Màu sắc", 1},
			{"Công suất", 1},
			{"Chất liệu", 1},
		},
		"Sách tiểu thuyết": {
			{"Thể loại", 1},
		},
		"Sách phi tiểu thuyết": {
			{"Thể loại", 1},
		},
		"Thiết bị thể dục": {
			{"Loại thiết bị", 1},
			{"Màu sắc", 1},
		},
		"Thể thao ngoài trời": {
			{"Loại thiết bị", 1},
			{"Màu sắc", 1},
		},
	}

	// Map categories to tag names
	categoryTags := map[string]string{
		"Điện thoại thông minh": "Công nghệ",
		"Máy tính xách tay":     "Công nghệ",
		"Thời trang nam":        "Thời trang",
		"Thời trang nữ":         "Thời trang",
		"Phụ kiện thời trang":   "Thời trang",
		"Đồ dùng nhà bếp":       "Gia dụng",
		"Đồ dùng vệ sinh":       "Gia dụng",
		"Sách tiểu thuyết":      "Sách",
		"Sách phi tiểu thuyết":  "Sách",
		"Thiết bị thể dục":      "Thể thao",
		"Thể thao ngoài trời":   "Thể thao",
	}

	// Map categories to product name prefixes
	categoryProductNames := map[string]func() string{
		"Điện thoại thông minh": func() string { return "Điện thoại " + gofakeit.ProductName() },
		"Máy tính xách tay":     func() string { return "Laptop " + gofakeit.ProductName() },
		"Thời trang nam":        func() string { return "Áo " + gofakeit.ProductName() + " Nam" },
		"Thời trang nữ":         func() string { return "Váy " + gofakeit.ProductName() + " Nữ" },
		"Phụ kiện thời trang":   func() string { return "Phụ kiện " + gofakeit.ProductName() },
		"Đồ dùng nhà bếp":       func() string { return "Máy " + gofakeit.ProductName() + " Nhà Bếp" },
		"Đồ dùng vệ sinh":       func() string { return "Máy " + gofakeit.ProductName() + " Vệ Sinh" },
		"Sách tiểu thuyết":      func() string { return "Sách " + gofakeit.BookTitle() },
		"Sách phi tiểu thuyết":  func() string { return "Sách " + gofakeit.BookTitle() },
		"Thiết bị thể dục":      func() string { return "Thiết bị " + gofakeit.ProductName() },
		"Thể thao ngoài trời":   func() string { return "Dụng cụ " + gofakeit.ProductName() },
	}

	batchSize := 100
	for _, supplierID := range supplierIDs {
		for i := 0; i < 10; i++ { // 10 products per supplier
			var products []struct {
				id                          uuid.UUID
				supplierID                  int64
				name, description, imageURL string
				categoryID                  int64
				status, taxClass            string
			}
			var variants []struct {
				id, productID       uuid.UUID
				sku, variantName    string
				price               float64
				inventoryQuantity   int
				shippingClass       string
				imageURL, altText   string
				isActive, isDefault bool
			}
			for j := 0; j < batchSize && (i*batchSize+j) < 10; j++ {
				// Randomly select a subcategory
				categoryID := subCategoryIDs[gofakeit.Number(0, len(subCategoryIDs)-1)]
				categoryName := subCategories[categoryID]

				// Generate product name based on category
				name := categoryProductNames[categoryName]()
				productID := uuid.New()
				products = append(products, struct {
					id                          uuid.UUID
					supplierID                  int64
					name, description, imageURL string
					categoryID                  int64
					status, taxClass            string
				}{
					id:          productID,
					supplierID:  supplierID,
					name:        name,
					description: gofakeit.Sentence(10),
					imageURL:    fmt.Sprintf("https://picsum.photos/seed/%s/150/150", strings.ReplaceAll(name, " ", "_")),
					categoryID:  categoryID,
					status:      "active",
					taxClass:    "standard",
				})

				// Create a variant for the product
				variantID := uuid.New()
				variants = append(variants, struct {
					id, productID       uuid.UUID
					sku, variantName    string
					price               float64
					inventoryQuantity   int
					shippingClass       string
					imageURL, altText   string
					isActive, isDefault bool
				}{
					id:                variantID,
					productID:         productID,
					sku:               fmt.Sprintf("SKU-%s-%d", strings.ReplaceAll(name, " ", "-"), gofakeit.Number(1000, 9999)),
					variantName:       name + " Variant",
					price:             gofakeit.Price(10, 1000),
					inventoryQuantity: gofakeit.Number(1, 100),
					shippingClass:     "standard",
					imageURL:          fmt.Sprintf("https://picsum.photos/seed/%s/150/150", strings.ReplaceAll(name, " ", "_")),
					altText:           name,
					isActive:          true,
					isDefault:         true,
				})
			}

			// Insert products
			if len(products) > 0 {
				var args []interface{}
				query := `INSERT INTO products (id, supplier_id, name, description, image_url, category_id, status, tax_class) VALUES `
				valueStrings := make([]string, len(products))
				for j, p := range products {
					idx := j * 8
					valueStrings[j] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7, idx+8)
					args = append(args, p.id, p.supplierID, p.name, p.description, p.imageURL, p.categoryID, p.status, p.taxClass)
				}
				query += strings.Join(valueStrings, ",") + ";"
				_, err := db.Exec(ctx, query, args...)
				if err != nil {
					log.Fatal("insert products:", err)
				}
			} else {
				log.Printf("Warning: No products to insert for supplier %d, iteration %d", supplierID, i)
			}

			// Insert product variants
			if len(variants) > 0 {
				var varArgs []interface{}
				varQuery := `INSERT INTO product_variants (id, product_id, sku, variant_name, price, inventory_quantity, shipping_class, image_url, alt_text, is_active, is_default, currency) VALUES `
				varValueStrings := make([]string, len(variants))
				for j, v := range variants {
					idx := j * 12
					varValueStrings[j] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7, idx+8, idx+9, idx+10, idx+11, idx+12)
					varArgs = append(varArgs, v.id, v.productID, v.sku, v.variantName, v.price, v.inventoryQuantity, v.shippingClass, v.imageURL, v.altText, v.isActive, v.isDefault, "VND")
				}
				varQuery += strings.Join(varValueStrings, ",") + ";"
				_, err = db.Exec(ctx, varQuery, varArgs...)
				if err != nil {
					log.Fatal("insert product variants:", err)
				}
			} else {
				log.Printf("Warning: No variants to insert for supplier %d, iteration %d", supplierID, i)
			}

			// Insert product tags
			if len(products) > 0 {
				for _, p := range products {
					// Find the category name for this product
					categoryName := subCategories[p.categoryID]
					tagName := categoryTags[categoryName]
					tagID := tagIDs[tagName]
					// Randomly add "Giảm giá" or "Mới" tags
					if gofakeit.Bool() {
						if gofakeit.Bool() {
							tagID = tagIDs["Giảm giá"]
						} else {
							tagID = tagIDs["Mới"]
						}
					}
					_, err := db.Exec(ctx, `
						INSERT INTO products_tags (product_id, tag_id) 
						VALUES ($1, $2) 
						ON CONFLICT DO NOTHING;
					`, p.id, tagID)
					if err != nil {
						log.Fatal("insert product tags:", err)
					}
				}
			} else {
				log.Printf("Warning: No product tags to insert for supplier %d, iteration %d", supplierID, i)
			}

			// Insert product variant attributes
			if len(variants) > 0 {
				for _, v := range variants {
					// Find the product corresponding to this variant
					var categoryName string
					for _, p := range products {
						if p.id == v.productID {
							categoryName = subCategories[p.categoryID]
							break
						}
					}
					attrs, exists := categoryAttributes[categoryName]
					if !exists {
						log.Printf("Warning: No attributes defined for category %s", categoryName)
						continue
					}
					for _, attr := range attrs {
						options := attrOptionIDs[attr.attrName]
						if len(options) == 0 {
							log.Printf("Warning: No options found for attribute %s", attr.attrName)
							continue
						}
						gofakeit.ShuffleAnySlice(options)
						for k := 0; k < attr.count && k < len(options); k++ {
							optionID := options[k]
							defID, exists := attrOptionToDefIDs[optionID]
							if !exists {
								log.Printf("Warning: No attribute definition found for option ID %d", optionID)
								continue
							}
							_, err := db.Exec(ctx, `
								INSERT INTO product_variant_attributes (product_variant_id, attribute_definition_id, attribute_option_id) 
								VALUES ($1, $2, $3) 
								ON CONFLICT DO NOTHING;
							`, v.id, defID, optionID)
							if err != nil {
								log.Fatal("insert product variant attributes:", err)
							}
						}
					}
				}
			} else {
				log.Printf("Warning: No product variant attributes to insert for supplier %d, iteration %d", supplierID, i)
			}
		}
	}
}
