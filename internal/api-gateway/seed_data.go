package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TienMinh25/ecommerce-platform/internal/utils"
	"github.com/brianvoe/gofakeit/v7"
	"log"
	"math/rand"
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
	seedUsers(ctx, pool, 2000000) // Số lượng user có thể thay đổi ở đây

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
		"Store Management", "Onboarding",
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
	for i := 1; i <= 10; i++ {
		permissions = append(permissions, PermissionDetail{
			ModuleID:    i,
			Permissions: []int{1, 2, 3, 4, 5, 6},
		})
	}
	bytes, _ := json.Marshal(permissions)
	_, _ = db.Exec(ctx, `INSERT INTO role_user_permissions (role_id, user_id, permission_detail) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`, roleID, userID, bytes)
}

func seedUsers(ctx context.Context, db *pgxpool.Pool, total int) {
	roleMap := map[string]int{}
	rows, _ := db.Query(ctx, `SELECT id, role_name FROM roles`)
	for rows.Next() {
		var id int
		var role string
		_ = rows.Scan(&id, &role)
		roleMap[role] = id
	}
	emails := map[string]bool{"admin@admin.com": true}
	for i := 0; i < total; i++ {
		name := gofakeit.Name()
		email := gofakeit.Email()
		for emails[email] {
			email = gofakeit.Email()
		}
		emails[email] = true
		birth := gofakeit.DateRange(time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC))
		avatar := fmt.Sprintf("https://ui-avatars.com/api/?name=%s", name)

		var userID int64
		db.QueryRow(ctx, `INSERT INTO users (fullname, email, avatar_url, phone, birthdate, email_verified, phone_verified)
			VALUES ($1, $2, $3, $4, $5, TRUE, TRUE) RETURNING id`,
			name, email, avatar, gofakeit.Phone(), birth).Scan(&userID)

		hash, _ := utils.HashPassword("123456")
		_, _ = db.Exec(ctx, `INSERT INTO user_password (id, password) VALUES ($1, $2)`, userID, hash)

		// Chọn role mặc định là customer, nhưng vẫn có thể có supplier, deliverer
		roles := []string{"customer", "supplier", "deliverer"}
		chosen := roles[rand.Intn(len(roles))]
		roleID := roleMap[chosen]

		// Lấy quyền của customer (các quyền tối thiểu của người dùng)
		customerRoleID := roleMap["customer"]
		var permDetails []PermissionDetail
		rows, _ := db.Query(ctx, `SELECT permission_detail FROM role_user_permissions WHERE role_id = $1`, customerRoleID)
		for rows.Next() {
			var permissionDetail []byte
			_ = rows.Scan(&permissionDetail)
			var permission []PermissionDetail
			_ = json.Unmarshal(permissionDetail, &permission)
			permDetails = append(permDetails, permission...)
		}

		// Gán quyền cho user
		bytes, _ := json.Marshal(permDetails)
		_, _ = db.Exec(ctx, `INSERT INTO role_user_permissions (role_id, user_id, permission_detail) VALUES ($1, $2, $3)`, roleID, userID, bytes)
	}
}
