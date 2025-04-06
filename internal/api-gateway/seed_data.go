package main

import (
	"context"
	"encoding/json"
	"fmt"
	api_gateway_models "github.com/TienMinh25/ecommerce-platform/internal/api-gateway/models"
	"github.com/TienMinh25/ecommerce-platform/internal/common"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/brianvoe/gofakeit/v7"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionDetail struct {
	ModuleID    int   `json:"module_id"`
	Permissions []int `json:"permissions"`
}

func main() {
	ctx := context.Background()
	gofakeit.Seed(0)

	dsn := "postgres://admin:admin@localhost:5432/api_gateway_db?sslmode=disable"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Unable to connect to DB:", err)
	}
	defer pool.Close()
	log.Println("🙂‍↔️Connected to DB api_gateway_db")
	log.Println("🏃‍♂️Seeding data...")

	seedRoles(ctx, pool)
	seedModules(ctx, pool)
	seedPermissions(ctx, pool)
	seedAddressTypes(ctx, pool)
	seedAdmin(ctx, pool)
	seedUsers(ctx, pool, 1000000, 2000) // Số lượng user có thể thay đổi ở đây

	fmt.Println("✅ Seed completed successfully")
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
		_, _ = db.Exec(ctx, `INSERT INTO roles (role_name, description) VALUES ($1, $2) ON CONFLICT DO NOTHING`, r.name, r.desc)
	}
}

func seedModules(ctx context.Context, db *pgxpool.Pool) {
	modules := []string{
		"User Management", "Role & Permission", "Product Management", "Cart",
		"Order Management", "Payment", "Shipping Management", "Review & Rating",
		"Store Management", "Onboarding", "Address Type Management", "Module Management",
	}
	for _, m := range modules {
		_, _ = db.Exec(ctx, `INSERT INTO modules (name) VALUES ($1) ON CONFLICT DO NOTHING`, m)
	}
}

func seedPermissions(ctx context.Context, db *pgxpool.Pool) {
	perms := []string{"create", "update", "delete", "read", "approve", "reject"}
	for _, p := range perms {
		_, _ = db.Exec(ctx, `INSERT INTO permissions (name) VALUES ($1) ON CONFLICT DO NOTHING`, p)
	}
}

func seedAddressTypes(ctx context.Context, db *pgxpool.Pool) {
	types := []string{"Home", "Office", "Warehouse", "Storefront", "Other"}
	for _, t := range types {
		_, _ = db.Exec(ctx, `INSERT INTO address_types (address_type) VALUES ($1) ON CONFLICT DO NOTHING`, t)
	}
}

func seedAdmin(ctx context.Context, db *pgxpool.Pool) {
	_, _ = db.Exec(ctx, `INSERT INTO users (fullname, email, avatar_url, email_verified, status, phone_verified) 
		VALUES ('Admin User', 'admin@admin.com', 'https://ui-avatars.com/api/?name=Admin User', TRUE, 'active', TRUE) ON CONFLICT DO NOTHING`)

	var userID int64
	db.QueryRow(ctx, `SELECT id FROM users WHERE email='admin@admin.com'`).Scan(&userID)

	hash, _ := utils.HashPassword("admin123")
	_, _ = db.Exec(ctx, `INSERT INTO user_password (id, password) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, hash)

	var roleID int64
	db.QueryRow(ctx, `SELECT id FROM roles WHERE role_name='admin'`).Scan(&roleID)

	permissions := []PermissionDetail{}
	for i := 1; i <= 12; i++ {
		permissions = append(permissions, PermissionDetail{
			ModuleID:    i,
			Permissions: []int{1, 2, 3, 4, 5, 6},
		})
	}
	bytes, _ := json.Marshal(permissions)
	_, _ = db.Exec(ctx, `INSERT INTO role_user_permissions (role_id, user_id, permission_detail) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, roleID, userID, bytes)
}

func seedUsers(ctx context.Context, db *pgxpool.Pool, total int, batchSize int) {
	if batchSize <= 0 {
		log.Fatal("batchSize must be > 0")
	}

	// Sử dụng bộ đếm với mutex để tạo email độc nhất
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

	// Lấy roleID trước
	var roleID int
	if err := db.QueryRow(ctx, `SELECT id FROM roles WHERE role_name = $1`, common.RoleCustomer).Scan(&roleID); err != nil {
		log.Fatal("get role_id:", err)
	}

	// Chuẩn bị permission
	permDetail := []api_gateway_models.PermissionDetailType{
		{ModuleID: 1, Permissions: []int{4, 2}},
		{ModuleID: 4, Permissions: []int{1, 2, 3, 4}},
		{ModuleID: 5, Permissions: []int{1, 4}},
		{ModuleID: 6, Permissions: []int{1, 4, 3}},
		{ModuleID: 7, Permissions: []int{4}},
		{ModuleID: 8, Permissions: []int{1, 4, 3}},
	}
	permBytes, _ := json.Marshal(permDetail)

	// Số lượng goroutine
	numGoroutines := 15

	// Chia đều công việc cho các goroutine
	perGoroutine := total / numGoroutines
	remainder := total % numGoroutines

	// Channel để theo dõi tiến độ
	type progressUpdate struct {
		goroutineID int
		count       int
	}
	progressChan := make(chan progressUpdate, numGoroutines)

	// WaitGroup để đợi tất cả goroutine hoàn thành
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	fmt.Println("🚀 Starting seed with", numGoroutines, "goroutines")

	// Khởi chạy các goroutine
	for i := 0; i < numGoroutines; i++ {
		// Tính toán số lượng record mỗi goroutine cần xử lý
		workload := perGoroutine
		if i < remainder {
			workload++
		}

		go func(goroutineID int, workload int) {
			defer wg.Done()

			seeded := 0
			goroutineTotal := workload

			// Mỗi goroutine sẽ thực hiện công việc của mình theo batch
			for seeded < goroutineTotal {
				toSeed := batchSize
				if goroutineTotal-seeded < batchSize {
					toSeed = goroutineTotal - seeded
				}

				var users []userInput
				for len(users) < toSeed {
					name := gofakeit.Name()

					// Lấy một số thứ tự duy nhất để thêm vào email
					sequenceMutex.Lock()
					seq := emailSequence
					emailSequence++
					sequenceMutex.Unlock()

					// Tạo email độc nhất bằng cách thêm số thứ tự
					username := strings.ToLower(strings.Replace(name, " ", ".", -1))
					email := fmt.Sprintf("%s.%d@example.com", username, seq)

					avatar := fmt.Sprintf("https://ui-avatars.com/api/?name=%s", name)
					birth := gofakeit.DateRange(
						time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
						time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
					)
					hash, _ := utils.HashPassword("123456")

					users = append(users, userInput{
						fullname:     name,
						email:        email,
						avatar:       avatar,
						phone:        gofakeit.Phone(),
						birthdate:    birth,
						passwordHash: hash,
					})
				}

				// Insert batch users
				var args []interface{}
				query := `INSERT INTO users (fullname, email, avatar_url, phone, birthdate, email_verified, phone_verified) VALUES `
				valueStrings := make([]string, len(users))
				for i, u := range users {
					idx := i * 7
					valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7)
					args = append(args, u.fullname, u.email, u.avatar, u.phone, u.birthdate, true, true)
				}
				query += strings.Join(valueStrings, ",") + " RETURNING id"

				rows, err := db.Query(ctx, query, args...)
				if err != nil {
					log.Fatal("insert users:", err)
				}

				var userIDs []int64
				for rows.Next() {
					var id int64
					if err := rows.Scan(&id); err != nil {
						log.Fatal("scan user id:", err)
					}
					userIDs = append(userIDs, id)
				}
				if err := rows.Err(); err != nil {
					log.Fatal("row error:", err)
				}
				rows.Close()

				// Insert user_password
				var pwArgs []interface{}
				pwValues := make([]string, len(userIDs))
				for i, id := range userIDs {
					idx := i * 2
					pwValues[i] = fmt.Sprintf("($%d, $%d)", idx+1, idx+2)
					pwArgs = append(pwArgs, id, users[i].passwordHash)
				}
				pwQuery := `INSERT INTO user_password (id, password) VALUES ` + strings.Join(pwValues, ",")
				if _, err := db.Exec(ctx, pwQuery, pwArgs...); err != nil {
					log.Fatal("insert user_password:", err)
				}

				// Insert role_user_permissions
				var permArgs []interface{}
				permValues := make([]string, len(userIDs))
				for i, id := range userIDs {
					idx := i * 3
					permValues[i] = fmt.Sprintf("($%d, $%d, $%d::jsonb)", idx+1, idx+2, idx+3)
					permArgs = append(permArgs, roleID, id, string(permBytes))
				}
				permQuery := `INSERT INTO role_user_permissions (role_id, user_id, permission_detail) VALUES ` + strings.Join(permValues, ",")
				if _, err := db.Exec(ctx, permQuery, permArgs...); err != nil {
					log.Fatal("insert role_user_permissions:", err)
				}

				seeded += toSeed

				// Báo cáo tiến độ
				progressChan <- progressUpdate{goroutineID: goroutineID, count: seeded}
			}
		}(i, workload)
	}

	// Goroutine để theo dõi và in tiến độ
	go func() {
		progress := make([]int, numGoroutines)
		totalInserted := 0

		for update := range progressChan {
			progress[update.goroutineID] = update.count

			// Tính tổng số đã insert
			totalInserted = 0
			for _, count := range progress {
				totalInserted += count
			}

			fmt.Printf("⏳ Progress: Goroutine #%d inserted %d records. Total: %d/%d (%.2f%%)\n",
				update.goroutineID,
				update.count,
				totalInserted,
				total,
				float64(totalInserted)*100/float64(total))
		}
	}()

	// Đợi tất cả goroutine hoàn thành
	wg.Wait()
	close(progressChan)

	// Đảm bảo goroutine theo dõi tiến độ kết thúc
	time.Sleep(100 * time.Millisecond)

	fmt.Println("🎉 Done seeding users with concurrent goroutines.")
}
