package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"math"
	"net/http"
	"os"
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

// Cấu trúc cho dữ liệu địa giới hành chính
type Province struct {
	ID        string     `json:"Id"`
	Name      string     `json:"Name"`
	Districts []District `json:"Districts"`
}

type District struct {
	ID    string `json:"Id"`
	Name  string `json:"Name"`
	Wards []Ward `json:"Wards"`
}

type Ward struct {
	ID   string `json:"Id"`
	Name string `json:"Name"`
}

// Cấu trúc sản phẩm từ Shopee (dùng để crawl dữ liệu)
type ShopeeProduct struct {
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	Categories  []string `json:"categories"`
}

// Danh sách attribute cho sản phẩm theo danh mục
var categoryAttributes = map[string]map[string][]string{
	"Điện thoại thông minh": {
		"Màu sắc":          {"Đen", "Trắng", "Xanh", "Đỏ", "Hồng", "Vàng", "Bạc", "Xám"},
		"Dung lượng":       {"64GB", "128GB", "256GB", "512GB", "1TB"},
		"RAM":              {"4GB", "6GB", "8GB", "12GB", "16GB"},
		"Hệ điều hành":     {"Android", "iOS"},
		"Kích cỡ màn hình": {"5.5 inch", "6.1 inch", "6.4 inch", "6.7 inch", "6.9 inch"},
	},
	"Máy tính xách tay": {
		"Màu sắc":  {"Đen", "Trắng", "Bạc", "Xám", "Xanh"},
		"CPU":      {"Intel Core i3", "Intel Core i5", "Intel Core i7", "Intel Core i9", "AMD Ryzen 5", "AMD Ryzen 7", "AMD Ryzen 9"},
		"RAM":      {"4GB", "8GB", "16GB", "32GB", "64GB"},
		"Ổ cứng":   {"256GB SSD", "512GB SSD", "1TB SSD", "2TB SSD", "512GB SSD + 1TB HDD"},
		"Màn hình": {"13.3 inch", "14 inch", "15.6 inch", "16 inch", "17.3 inch"},
	},
	"Thời trang nam": {
		"Màu sắc":    {"Đen", "Trắng", "Xanh Navy", "Xanh Lá", "Đỏ", "Xám", "Nâu", "Be"},
		"Kích thước": {"S", "M", "L", "XL", "XXL"},
		"Chất liệu":  {"Cotton", "Polyester", "Len", "Lụa", "Vải Lanh", "Vải Jeans", "Vải Thun"},
		"Kiểu dáng":  {"Regular Fit", "Slim Fit", "Loose Fit", "Skinny Fit"},
		"Mùa":        {"Xuân Hè", "Thu Đông", "Bốn mùa"},
		"Xuất xứ":    {"Việt Nam", "Trung Quốc", "Hàn Quốc", "Thái Lan", "Nhật Bản", "Mỹ"},
		"Phong cách": {"Casual", "Formal", "Street Style", "Vintage", "Minimalist"},
	},
	"Thời trang nữ": {
		"Màu sắc":    {"Đen", "Trắng", "Đỏ", "Hồng", "Xanh Navy", "Xanh Lá", "Tím", "Vàng", "Be", "Nâu"},
		"Kích thước": {"S", "M", "L", "XL", "XXL"},
		"Chất liệu":  {"Cotton", "Polyester", "Len", "Lụa", "Vải Lanh", "Vải Jeans", "Vải Thun", "Ren"},
		"Kiểu dáng":  {"Regular Fit", "Slim Fit", "Loose Fit", "Oversize"},
		"Mùa":        {"Xuân Hè", "Thu Đông", "Bốn mùa"},
		"Xuất xứ":    {"Việt Nam", "Trung Quốc", "Hàn Quốc", "Thái Lan", "Nhật Bản", "Mỹ"},
		"Phong cách": {"Casual", "Formal", "Street Style", "Vintage", "Minimalist", "Sexy", "Bohemian"},
	},
	"Đồ gia dụng": {
		"Màu sắc":   {"Đen", "Trắng", "Bạc", "Xám", "Đỏ", "Xanh", "Hồng", "Vàng"},
		"Chất liệu": {"Nhựa", "Kim loại", "Gỗ", "Thủy tinh", "Gốm sứ", "Silicone", "Inox"},
		"Công suất": {"500W", "700W", "1000W", "1200W", "1500W", "2000W"},
		"Xuất xứ":   {"Việt Nam", "Trung Quốc", "Hàn Quốc", "Thái Lan", "Nhật Bản", "Mỹ", "Đức"},
		"Bảo hành":  {"6 tháng", "12 tháng", "24 tháng", "36 tháng", "60 tháng"},
	},
	"Sách": {
		"Thể loại":     {"Tiểu thuyết", "Khoa học viễn tưởng", "Kinh doanh", "Tâm lý học", "Kỹ năng sống", "Lịch sử", "Trinh thám", "Hồi ký"},
		"Ngôn ngữ":     {"Tiếng Việt", "Tiếng Anh", "Song ngữ Anh-Việt"},
		"Tác giả":      {"Nguyễn Nhật Ánh", "Nguyễn Ngọc Tư", "Trang Hạ", "Paulo Coelho", "Haruki Murakami", "J.K. Rowling", "Stephen King"},
		"Nhà xuất bản": {"NXB Kim Đồng", "NXB Trẻ", "NXB Tổng hợp TPHCM", "NXB Hội Nhà văn", "NXB Giáo dục", "NXB Lao động"},
		"Bìa sách":     {"Bìa mềm", "Bìa cứng", "Bìa gập"},
	},
	"Thể thao": {
		"Màu sắc":       {"Đen", "Trắng", "Xanh", "Đỏ", "Xám", "Cam"},
		"Chất liệu":     {"Nhựa", "Kim loại", "Cao su", "Vải", "Da tổng hợp", "Sợi carbon"},
		"Kích thước":    {"S", "M", "L", "XL", "XXL", "Freesize"},
		"Thương hiệu":   {"Nike", "Adidas", "Puma", "Under Armour", "The North Face", "Columbia", "Lining"},
		"Xuất xứ":       {"Việt Nam", "Trung Quốc", "Mỹ", "Đức", "Nhật Bản", "Thái Lan"},
		"Loại thiết bị": {"Tập lực", "Tập cardio", "Đồ bảo hộ", "Phụ kiện", "Quần áo tập"},
	},
}

// Danh sách tên sản phẩm mẫu theo danh mục
var categoryProductNames = map[string][]string{
	"Điện thoại thông minh": {
		"iPhone 13 Pro Max", "iPhone 14", "Samsung Galaxy S22 Ultra", "Samsung Galaxy Z Fold 4",
		"Xiaomi Redmi Note 11", "Xiaomi 12T Pro", "OPPO Reno8 Pro", "OPPO Find X5 Pro",
		"Vivo V25 Pro", "Realme GT Neo 3", "Nokia G21", "Huawei Nova 10",
	},
	"Máy tính xách tay": {
		"MacBook Air M2", "MacBook Pro 14", "Dell XPS 13", "Dell Inspiron 15",
		"HP Spectre x360", "HP Pavilion 15", "Lenovo ThinkPad X1 Carbon", "Lenovo Yoga 7i",
		"Asus ZenBook 14", "Asus ROG Zephyrus G14", "Acer Swift 5", "MSI Prestige 14",
	},
	"Thời trang nam": {
		"Áo sơ mi nam dài tay", "Áo thun nam cổ tròn", "Áo thun polo nam", "Áo khoác denim nam",
		"Áo khoác bomber nam", "Quần jeans nam slim fit", "Quần kaki nam", "Quần short nam",
		"Bộ vest nam công sở", "Áo len nam cổ tròn", "Áo hoodie nam", "Quần tây nam công sở",
	},
	"Thời trang nữ": {
		"Áo sơ mi nữ công sở", "Áo blouse nữ", "Áo thun nữ cổ tròn", "Áo khoác denim nữ",
		"Đầm suông nữ", "Đầm ôm body nữ", "Chân váy chữ A", "Chân váy tennis",
		"Quần jeans nữ ống rộng", "Quần culottes nữ", "Áo cardigan nữ", "Set đồ nữ hai mảnh",
	},
	"Đồ gia dụng": {
		"Nồi cơm điện", "Máy xay sinh tố", "Bếp từ đơn", "Bếp gas đôi",
		"Lò vi sóng", "Ấm đun nước siêu tốc", "Máy lọc không khí", "Quạt điều hòa",
		"Máy hút bụi", "Bàn ủi hơi nước", "Nồi chiên không dầu", "Máy rửa chén",
	},
	"Sách": {
		"Nhà Giả Kim", "Đắc Nhân Tâm", "Cà Phê Cùng Tony", "Người Giàu Có Nhất Thành Babylon",
		"Hai Số Phận", "Điều Kỳ Diệu Của Tiệm Tạp Hóa Namiya", "Bước Chậm Lại Giữa Thế Gian Vội Vã", "Tuổi Trẻ Đáng Giá Bao Nhiêu",
		"Chúng Ta Rồi Sẽ Hạnh Phúc, Theo Những Cách Khác Nhau", "Khéo Ăn Nói Sẽ Có Được Thiên Hạ", "Tôi Tài Giỏi, Bạn Cũng Thế", "Dám Nghĩ Lớn",
	},
	"Thể thao": {
		"Tạ tay 5kg", "Thảm tập yoga", "Dây nhảy thể dục", "Máy chạy bộ điện",
		"Xe đạp tập thể dục", "Găng tay tập gym", "Ghế tập bụng đa năng", "Bóng đá size 5",
		"Vợt cầu lông", "Vợt tennis", "Bộ cờ vua quốc tế", "Bàn bóng bàn",
	},
}

// Danh sách mô tả sản phẩm mẫu theo danh mục
var categoryProductDescriptions = map[string][]string{
	"Điện thoại thông minh": {
		"Sản phẩm công nghệ hiện đại với màn hình Retina sắc nét, camera độ phân giải cao và thời lượng pin dài.",
		"Thiết kế sang trọng với cấu hình mạnh mẽ, camera AI thông minh và khả năng chống nước IP68.",
		"Smartphone cao cấp với chip xử lý mới nhất, màn hình AMOLED 120Hz và sạc nhanh 65W.",
		"Điện thoại thông minh với camera chuyên nghiệp, khả năng quay video 4K và bộ nhớ lớn.",
	},
	"Máy tính xách tay": {
		"Laptop mỏng nhẹ với hiệu suất mạnh mẽ, thời lượng pin cả ngày và màn hình Retina sắc nét.",
		"Máy tính xách tay cao cấp dành cho công việc sáng tạo, đồ họa với card màn hình rời và SSD tốc độ cao.",
		"Laptop chuyên gaming với card đồ họa mạnh mẽ, tản nhiệt hiệu quả và bàn phím RGB.",
		"Máy tính 2-in-1 linh hoạt với màn hình cảm ứng, bút stylus và thiết kế gập xoay 360 độ.",
	},
	"Thời trang nam": {
		"Áo thời trang nam phong cách Hàn Quốc, chất liệu cao cấp thoáng mát và thấm hút mồ hôi tốt.",
		"Quần nam thiết kế hiện đại, form dáng vừa vặn tôn dáng người mặc và dễ phối đồ.",
		"Sản phẩm thời trang dành cho nam giới công sở với thiết kế lịch lãm, tinh tế và sang trọng.",
		"Trang phục nam thiết kế theo phong cách đường phố, cá tính và năng động dành cho giới trẻ.",
	},
	"Thời trang nữ": {
		"Thời trang nữ thiết kế theo xu hướng mới nhất, tôn dáng người mặc và phù hợp nhiều hoàn cảnh.",
		"Trang phục nữ phong cách Hàn Quốc với chất liệu cao cấp, thoáng mát và thấm hút tốt.",
		"Quần áo nữ thiết kế tinh tế với họa tiết độc đáo, phù hợp cho công sở và dạo phố.",
		"Đầm nữ thiết kế sang trọng, quyến rũ phù hợp cho các buổi tiệc và sự kiện quan trọng.",
	},
	"Đồ gia dụng": {
		"Thiết bị gia dụng cao cấp với công nghệ hiện đại, tiết kiệm điện và dễ dàng sử dụng.",
		"Sản phẩm gia dụng thông minh với khả năng kết nối điện thoại và điều khiển từ xa.",
		"Thiết bị nhà bếp đa năng với nhiều chức năng, giúp việc nấu nướng trở nên đơn giản và nhanh chóng.",
		"Sản phẩm gia dụng bền bỉ với chất liệu cao cấp và chế độ bảo hành dài hạn.",
	},
	"Sách": {
		"Cuốn sách best-seller với nội dung sâu sắc, đem lại nhiều bài học giá trị cho người đọc.",
		"Tác phẩm nổi tiếng của tác giả được yêu thích, đã được dịch ra nhiều thứ tiếng trên thế giới.",
		"Sách hay với nội dung bổ ích, ngôn từ cuốn hút và thông điệp ý nghĩa.",
		"Cuốn sách giúp bạn thay đổi tư duy, phát triển bản thân và đạt được thành công trong cuộc sống.",
	},
	"Thể thao": {
		"Thiết bị tập thể thao cao cấp với chất liệu bền bỉ, an toàn và hiệu quả cao.",
		"Dụng cụ thể thao đa năng giúp bạn tập luyện nhiều nhóm cơ khác nhau.",
		"Sản phẩm thể thao chuyên nghiệp được thiết kế bởi các chuyên gia hàng đầu.",
		"Thiết bị tập luyện tại nhà tiện lợi, tiết kiệm không gian và dễ dàng cất gọn.",
	},
}

// Danh sách các APIs hỗ trợ dữ liệu địa giới hành chính Việt Nam
var vietnamGeoAPIs = []string{
	"https://provinces.open-api.vn/api/?depth=3", // API với đầy đủ phường/xã
	"https://vietnam-administrative-divisions.vercel.app/api/",
	"https://vapi.vnappmob.com/api/province/",
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
	return getSampleAdministrativeDivisions()
}

// Dữ liệu mẫu nếu không tải được từ API hoặc file
func getSampleAdministrativeDivisions() []Province {
	provinces := []Province{
		{ID: "01", Name: "Hà Nội", Districts: []District{
			{ID: "001", Name: "Ba Đình", Wards: []Ward{{ID: "00001", Name: "Phúc Xá"}, {ID: "00002", Name: "Trúc Bạch"}}},
			{ID: "002", Name: "Hoàn Kiếm", Wards: []Ward{{ID: "00003", Name: "Hàng Bạc"}, {ID: "00004", Name: "Hàng Bồ"}}},
			{ID: "003", Name: "Tây Hồ", Wards: []Ward{{ID: "00005", Name: "Bưởi"}, {ID: "00006", Name: "Nhật Tân"}}},
			{ID: "004", Name: "Long Biên", Wards: []Ward{{ID: "00007", Name: "Bồ Đề"}, {ID: "00008", Name: "Sài Đồng"}}},
			{ID: "005", Name: "Cầu Giấy", Wards: []Ward{{ID: "00009", Name: "Quan Hoa"}, {ID: "00010", Name: "Nghĩa Đô"}}},
		}},
		{ID: "02", Name: "TP Hồ Chí Minh", Districts: []District{
			{ID: "006", Name: "Quận 1", Wards: []Ward{{ID: "00011", Name: "Bến Nghé"}, {ID: "00012", Name: "Bến Thành"}}},
			{ID: "007", Name: "Quận 3", Wards: []Ward{{ID: "00013", Name: "Võ Thị Sáu"}, {ID: "00014", Name: "Nguyễn Cư Trinh"}}},
			{ID: "008", Name: "Quận 7", Wards: []Ward{{ID: "00015", Name: "Tân Thuận Đông"}, {ID: "00016", Name: "Tân Thuận Tây"}}},
			{ID: "009", Name: "Bình Thạnh", Wards: []Ward{{ID: "00017", Name: "Phường 1"}, {ID: "00018", Name: "Phường 2"}}},
			{ID: "010", Name: "Thủ Đức", Wards: []Ward{{ID: "00019", Name: "Linh Đông"}, {ID: "00020", Name: "Linh Tây"}}},
		}},
		{ID: "03", Name: "Đà Nẵng", Districts: []District{
			{ID: "011", Name: "Hải Châu", Wards: []Ward{{ID: "00021", Name: "Thanh Bình"}, {ID: "00022", Name: "Hải Châu I"}}},
			{ID: "012", Name: "Thanh Khê", Wards: []Ward{{ID: "00023", Name: "Tam Thuận"}, {ID: "00024", Name: "Thanh Khê Đông"}}},
			{ID: "013", Name: "Sơn Trà", Wards: []Ward{{ID: "00025", Name: "An Hải Bắc"}, {ID: "00026", Name: "Mân Thái"}}},
		}},
		{ID: "04", Name: "Hải Phòng", Districts: []District{
			{ID: "014", Name: "Hồng Bàng", Wards: []Ward{{ID: "00027", Name: "Minh Khai"}, {ID: "00028", Name: "Quang Trung"}}},
			{ID: "015", Name: "Ngô Quyền", Wards: []Ward{{ID: "00029", Name: "Lạch Tray"}, {ID: "00030", Name: "Đông Khê"}}},
		}},
		{ID: "05", Name: "Cần Thơ", Districts: []District{
			{ID: "016", Name: "Ninh Kiều", Wards: []Ward{{ID: "00031", Name: "Tân An"}, {ID: "00032", Name: "An Phú"}}},
			{ID: "017", Name: "Bình Thủy", Wards: []Ward{{ID: "00033", Name: "Bình Thủy"}, {ID: "00034", Name: "Trà An"}}},
		}},
		{ID: "06", Name: "Nha Trang", Districts: []District{
			{ID: "018", Name: "Khánh Hòa", Wards: []Ward{{ID: "00035", Name: "Vạn Thạnh"}, {ID: "00036", Name: "Phương Sài"}}},
			{ID: "019", Name: "Vĩnh Trường", Wards: []Ward{{ID: "00037", Name: "Vĩnh Nguyên"}, {ID: "00038", Name: "Vĩnh Hòa"}}},
		}},
		{ID: "07", Name: "Huế", Districts: []District{
			{ID: "020", Name: "Thừa Thiên", Wards: []Ward{{ID: "00039", Name: "Phú Hậu"}, {ID: "00040", Name: "Vĩnh Ninh"}}},
			{ID: "021", Name: "Phú Vang", Wards: []Ward{{ID: "00041", Name: "Thuận An"}, {ID: "00042", Name: "Phú Thuận"}}},
		}},
	}
	return provinces
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
			{ModuleID: 11, Permissions: []int{4}},
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
	count := len(userIDs) / 10 // 10% users as suppliers
	supplierUserIDs := make([]int64, 0, count+1)

	// Đảm bảo admin cũng là supplier (ID=1)
	var adminID int64 = 1 // Thông thường admin là ID đầu tiên
	supplierUserIDs = append(supplierUserIDs, adminID)

	// Trộn ngẫu nhiên để chọn users làm supplier
	gofakeit.ShuffleAnySlice(userIDs)
	for i := 0; i < count && i < len(userIDs); i++ {
		if userIDs[i] != adminID {
			supplierUserIDs = append(supplierUserIDs, userIDs[i])
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
	count := len(userIDs) / 20 // 5% users as deliverers
	delivererUserIDs := make([]int64, 0, count+1)

	// Thêm admin vào danh sách deliverer
	var adminID int64 = 1
	delivererUserIDs = append(delivererUserIDs, adminID)

	// Trộn ngẫu nhiên để chọn users làm deliverer
	gofakeit.ShuffleAnySlice(userIDs)
	for i := 0; i < count && i < len(userIDs); i++ {
		if userIDs[i] != adminID {
			delivererUserIDs = append(delivererUserIDs, userIDs[i])
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

// Order Service Seeding
func seedOrderIndependentTables(ctx context.Context, db *pgxpool.Pool, adminDivisions []Province) {
	seedAreasFromAdminDivisions(ctx, db, adminDivisions)
	seedPaymentMethods(ctx, db)
}

func seedAreasFromAdminDivisions(ctx context.Context, db *pgxpool.Pool, adminDivisions []Province) {
	if len(adminDivisions) == 0 {
		log.Println("⚠️ No administrative divisions data, using fallback data")
		seedAreas(ctx, db)
		return
	}

	// Chọn một số tỉnh/thành phố và quận/huyện để seed
	for _, province := range adminDivisions {
		for _, district := range province.Districts {
			// Chọn ngẫu nhiên một số phường/xã
			for _, ward := range district.Wards {
				areaCode := fmt.Sprintf("%s-%s-%s", province.ID, district.ID, ward.ID)

				_, err := db.Exec(ctx, `
					INSERT INTO areas (city, country, district, ward, area_code)
					VALUES ($1, 'Việt Nam', $2, $3, $4)
					ON CONFLICT (area_code) DO NOTHING;
				`, province.Name, district.Name, ward.Name, areaCode)

				if err != nil {
					log.Printf("Error inserting area: %v", err)
				}
			}
		}
	}
	log.Println("✅ Areas seeded successfully")
}

// Fallback cho seedAreas nếu không có dữ liệu
func seedAreas(ctx context.Context, db *pgxpool.Pool) {
	sampleAreas := []struct {
		city, district, ward, areaCode string
	}{
		{"Hà Nội", "Ba Đình", "Phúc Xá", "01-001-00001"},
		{"Hà Nội", "Ba Đình", "Trúc Bạch", "01-001-00002"},
		{"Hà Nội", "Hoàn Kiếm", "Hàng Bạc", "01-002-00003"},
		{"Hà Nội", "Hoàn Kiếm", "Hàng Bồ", "01-002-00004"},
		{"TP Hồ Chí Minh", "Quận 1", "Bến Nghé", "02-006-00011"},
		{"TP Hồ Chí Minh", "Quận 1", "Bến Thành", "02-006-00012"},
		{"TP Hồ Chí Minh", "Quận 3", "Võ Thị Sáu", "02-007-00013"},
		{"Đà Nẵng", "Hải Châu", "Thanh Bình", "03-011-00021"},
		{"Đà Nẵng", "Hải Châu", "Hải Châu I", "03-011-00022"},
		{"Hải Phòng", "Hồng Bàng", "Minh Khai", "04-014-00027"},
	}

	for _, area := range sampleAreas {
		_, err := db.Exec(ctx, `
			INSERT INTO areas (city, country, district, ward, area_code)
			VALUES ($1, 'Việt Nam', $2, $3, $4)
			ON CONFLICT (area_code) DO NOTHING;
		`, area.city, area.district, area.ward, area.areaCode)

		if err != nil {
			log.Printf("Error inserting area: %v", err)
		}
	}
	log.Println("✅ Sample areas seeded successfully")
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

func seedCategories(ctx context.Context, db *pgxpool.Pool) {
	// Danh mục chính
	mainCategories := []struct {
		name, desc, imageUrl string
	}{
		{
			"Điện tử & Công nghệ",
			"Các sản phẩm điện tử và công nghệ hiện đại",
			"https://images.unsplash.com/photo-1468495244123-6c6c332eeece?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1021&q=80",
		},
		{
			"Thời trang",
			"Quần áo, giày dép và phụ kiện thời trang",
			"https://images.unsplash.com/photo-1445205170230-053b83016050?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Nhà cửa & Đời sống",
			"Đồ gia dụng và vật dụng sinh hoạt hàng ngày",
			"https://images.unsplash.com/photo-1484101403633-562f891dc89a?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1174&q=80",
		},
		{
			"Sách & Văn phòng phẩm",
			"Sách, văn phòng phẩm và học cụ",
			"https://images.unsplash.com/photo-1526243741027-444d633d7365?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Thể thao & Du lịch",
			"Dụng cụ thể thao và đồ dùng du lịch",
			"https://images.unsplash.com/photo-1517649763962-0c623066013b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
	}

	// Tiến hành seed categories - Sửa để bỏ ON CONFLICT
	mainCategoryIDs := make(map[string]int64)
	for _, cat := range mainCategories {
		// Trước tiên kiểm tra xem danh mục đã tồn tại chưa
		var id int64
		err := db.QueryRow(ctx, `
            SELECT id FROM categories WHERE name = $1
        `, cat.name).Scan(&id)

		if err != nil {
			// Nếu không tìm thấy hoặc lỗi khác, thêm mới
			err = db.QueryRow(ctx, `
                INSERT INTO categories (name, description, image_url, is_active)
                VALUES ($1, $2, $3, TRUE)
                RETURNING id;
            `, cat.name, cat.desc, cat.imageUrl).Scan(&id)

			if err != nil {
				log.Printf("Error inserting main category: %v", err)
				continue
			}
		}

		mainCategoryIDs[cat.name] = id
	}

	// Seed danh mục con - cũng sửa để bỏ ON CONFLICT
	subCategories := []struct {
		name, desc, parent, imageUrl string
	}{
		{
			"Điện thoại thông minh",
			"Điện thoại thông minh từ các thương hiệu nổi tiếng",
			"Điện tử & Công nghệ",
			"https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=880&q=80",
		},
		{
			"Máy tính xách tay",
			"Laptop và máy tính xách tay các loại",
			"Điện tử & Công nghệ",
			"https://images.unsplash.com/photo-1496181133206-80ce9b88a853?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Thời trang nam",
			"Quần áo và phụ kiện dành cho nam giới",
			"Thời trang",
			"https://images.unsplash.com/photo-1490578474895-699cd4e2cf59?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Thời trang nữ",
			"Quần áo và phụ kiện dành cho nữ giới",
			"Thời trang",
			"https://images.unsplash.com/photo-1567401893414-76b7b1e5a7a5?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Đồ gia dụng",
			"Thiết bị điện và đồ dùng gia đình",
			"Nhà cửa & Đời sống",
			"https://images.unsplash.com/photo-1556909172-54557c7e4fb7?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Sách",
			"Sách tiếng Việt và ngoại văn các thể loại",
			"Sách & Văn phòng phẩm",
			"https://images.unsplash.com/photo-1495446815901-a7297e633e8d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Thể thao",
			"Dụng cụ tập luyện thể thao",
			"Thể thao & Du lịch",
			"https://images.unsplash.com/photo-1517649763962-0c623066013b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
	}

	for _, subCat := range subCategories {
		parentID, exists := mainCategoryIDs[subCat.parent]
		if !exists {
			log.Printf("Parent category not found: %s", subCat.parent)
			continue
		}

		// Kiểm tra xem danh mục con đã tồn tại chưa
		var existingID int64
		err := db.QueryRow(ctx, `
            SELECT id FROM categories WHERE name = $1
        `, subCat.name).Scan(&existingID)

		if err != nil {
			// Nếu không tìm thấy hoặc lỗi khác, thêm mới
			_, err = db.Exec(ctx, `
                INSERT INTO categories (name, description, parent_id, image_url, is_active)
                VALUES ($1, $2, $3, $4, TRUE);
            `, subCat.name, subCat.desc, parentID, subCat.imageUrl)

			if err != nil {
				log.Printf("Error inserting sub category: %v", err)
			}
		} else {
			// Nếu đã tồn tại, cập nhật
			_, err = db.Exec(ctx, `
                UPDATE categories 
                SET description = $1, parent_id = $2, image_url = $3, is_active = TRUE
                WHERE id = $4;
            `, subCat.desc, parentID, subCat.imageUrl, existingID)

			if err != nil {
				log.Printf("Error updating sub category: %v", err)
			}
		}
	}
	log.Println("✅ Categories seeded successfully")
}

func seedTags(ctx context.Context, db *pgxpool.Pool) {
	tags := []string{
		"Mới nhất", "Bán chạy", "Giảm giá", "Cao cấp", "Giá rẻ",
		"Chính hãng", "Chất lượng cao", "Hàng hiệu", "Thương hiệu", "Nhập khẩu",
		"Xu hướng", "Thịnh hành", "Ưu đãi", "Miễn phí vận chuyển", "Khuyến mãi",
		"Phân phối chính thức", "Hàng độc quyền", "Phiên bản giới hạn",
	}

	for _, tag := range tags {
		// Check if the tag already exists
		var tagExists bool
		err := db.QueryRow(ctx, `
            SELECT EXISTS (SELECT 1 FROM tags WHERE name = $1)
        `, tag).Scan(&tagExists)

		if err != nil {
			log.Printf("Error checking tag existence: %v", err)
			continue
		}

		if !tagExists {
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

func seedAttributeDefinitions(ctx context.Context, db *pgxpool.Pool) {
	attributes := []struct {
		name, desc, inputType    string
		isFilterable, isRequired bool
	}{
		{"Màu sắc", "Màu sắc của sản phẩm", "select", true, true},
		{"Kích thước", "Kích thước của sản phẩm", "select", true, true},
		{"Chất liệu", "Chất liệu của sản phẩm", "select", true, false},
		{"Dung lượng", "Dung lượng lưu trữ", "select", true, false},
		{"RAM", "Dung lượng RAM", "select", true, false},
		{"CPU", "Loại CPU", "select", true, false},
		{"Ổ cứng", "Loại và dung lượng ổ cứng", "select", true, false},
		{"Màn hình", "Kích thước màn hình", "select", true, false},
		{"Kiểu dáng", "Kiểu dáng sản phẩm", "select", true, false},
		{"Thương hiệu", "Thương hiệu sản phẩm", "select", true, false},
		{"Xuất xứ", "Quốc gia xuất xứ", "select", false, false},
		{"Công suất", "Công suất thiết bị", "select", false, false},
		{"Bảo hành", "Thời gian bảo hành", "select", false, false},
		{"Thể loại", "Thể loại sách", "select", true, false},
		{"Ngôn ngữ", "Ngôn ngữ sách", "select", true, false},
		{"Tác giả", "Tác giả sách", "select", true, false},
		{"Nhà xuất bản", "Nhà xuất bản sách", "select", false, false},
		{"Bìa sách", "Loại bìa sách", "select", false, false},
		{"Mùa", "Mùa phù hợp", "select", false, false},
		{"Phong cách", "Phong cách thời trang", "select", true, false},
		{"Loại thiết bị", "Loại thiết bị thể thao", "select", true, false},
		{"Hệ điều hành", "Hệ điều hành thiết bị", "select", true, false},
		// Add the missing attribute
		{"Kích cỡ màn hình", "Kích thước màn hình hiển thị", "select", true, false},
	}

	// Seed attribute definitions
	for _, attr := range attributes {
		var attrID int
		err := db.QueryRow(ctx, `
			INSERT INTO attribute_definitions (name, description, input_type, is_filterable, is_required)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (name) DO UPDATE
			SET description = $2, input_type = $3, is_filterable = $4, is_required = $5
			RETURNING id;
		`, attr.name, attr.desc, attr.inputType, attr.isFilterable, attr.isRequired).Scan(&attrID)

		if err != nil {
			log.Printf("Error inserting attribute definition: %v", err)
			continue
		}

		// Seed attribute options based on category_attributes map
		if options, exists := getAttributeOptions(attr.name); exists {
			for _, option := range options {
				_, err := db.Exec(ctx, `
					INSERT INTO attribute_options (attribute_definition_id, option_value)
					VALUES ($1, $2)
					ON CONFLICT (attribute_definition_id, option_value) DO NOTHING;
				`, attrID, option)

				if err != nil {
					log.Printf("Error inserting attribute option: %v", err)
				}
			}
		}
	}
	log.Println("✅ Attribute definitions and options seeded successfully")
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

		// Tạo supplier document
		_, err = partnerDb.Exec(ctx, `
			INSERT INTO supplier_documents (supplier_id, document_url, verification_status, admin_note)
			VALUES ($1, $2, 'approved', 'Đã xác thực hồ sơ')
			ON CONFLICT DO NOTHING;
		`, supplierID, "https://images.unsplash.com/photo-1600880292203-757bb62b4baf?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80")

		if err != nil {
			log.Printf("Error inserting supplier document: %v", err)
		}
	}

	log.Printf("✅ Created %d supplier profiles", len(supplierIDs))
	return supplierIDs
}

func seedEnhancedProducts(ctx context.Context, db *pgxpool.Pool, supplierIDs []int64) {
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

	if !exists {
		log.Printf("The product_variants table does not exist. Make sure tables are created in the correct order.")
		log.Printf("Attempting to seed products first before product variants.")
	}

	// Lấy danh sách categories
	categories := make(map[string]int64)
	rows, err := db.Query(ctx, `SELECT id, name FROM categories WHERE parent_id IS NOT NULL`)
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
		categories[name] = id
	}

	// Danh sách ảnh sản phẩm chất lượng cao từ Unsplash theo danh mục
	productImages := map[string][]string{
		"Điện thoại thông minh": {
			"https://images.unsplash.com/photo-1585060544812-6b45742d762f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1281&q=80",
			"https://images.unsplash.com/photo-1598327105666-5b89351aff97?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2227&q=80",
			"https://images.unsplash.com/photo-1529653762956-b0a27278529c?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1605236453806-6ff36851218e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1589492477829-5e65395b66cc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
			"https://images.unsplash.com/photo-1616348436168-de43ad0db179?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=781&q=80",
		},
		"Máy tính xách tay": {
			"https://images.unsplash.com/photo-1496181133206-80ce9b88a853?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
			"https://images.unsplash.com/photo-1593642702821-c8da6771f0c6?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1332&q=80",
			"https://images.unsplash.com/photo-1603302576837-37561b2e2302?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1168&q=80",
			"https://images.unsplash.com/photo-1611186871348-b1ce696e52c9?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1541807084-5c52b6b3adef?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
		},
		"Thời trang nam": {
			"https://images.unsplash.com/photo-1617137968427-85924c800a22?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1287&q=80",
			"https://images.unsplash.com/photo-1516257984-b1b4d707412e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
			"https://images.unsplash.com/photo-1553143820-3c5ea7ec8c4e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1611312449408-fcece27cdbb7?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1169&q=80",
			"https://images.unsplash.com/photo-1496345875659-11f7dd282d1d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Thời trang nữ": {
			"https://images.unsplash.com/photo-1552874869-5c39ec9288dc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
			"https://images.unsplash.com/photo-1566206091558-7f218b696731?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=764&q=80",
			"https://images.unsplash.com/photo-1577900232427-18219b9166a0?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1525507119028-ed4c629a60a3?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=735&q=80",
			"https://images.unsplash.com/photo-1554412933-514a83d2f3c8?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1172&q=80",
		},
		"Đồ gia dụng": {
			"https://images.unsplash.com/photo-1587316205943-b15dc52a12e0?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=678&q=80",
			"https://images.unsplash.com/photo-1594225513563-c9eecb233345?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=764&q=80",
			"https://images.unsplash.com/photo-1565065524861-0be4646f450b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=880&q=80",
			"https://images.unsplash.com/photo-1625575499389-0a2003624731?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=880&q=80",
			"https://images.unsplash.com/photo-1556911220-bda9f7f8677e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Sách": {
			"https://images.unsplash.com/photo-1589998059171-988d887df646?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1176&q=80",
			"https://images.unsplash.com/photo-1541963463532-d68292c34b19?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=688&q=80",
			"https://images.unsplash.com/photo-1544947950-fa07a98d237f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
			"https://images.unsplash.com/photo-1512820790803-83ca734da794?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1198&q=80",
			"https://images.unsplash.com/photo-1543002588-bfa74002ed7e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
		},
		"Thể thao": {
			"https://images.unsplash.com/photo-1574680096145-d05b474e2155?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1169&q=80",
			"https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1591311630200-ffa9120a540f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1584735935682-2f2b69dff9d2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
			"https://images.unsplash.com/photo-1517649763962-0c623066013b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
	}

	// Số lượng sản phẩm đã tạo
	totalProducts := 0

	// Tạo sản phẩm cho mỗi danh mục
	for categoryName, categoryID := range categories {
		// Check for product names & descriptions
		productNames, ok := categoryProductNames[categoryName]
		if !ok {
			log.Printf("No product names found for category: %s", categoryName)
			continue
		}

		productDescriptions, ok := categoryProductDescriptions[categoryName]
		if !ok {
			log.Printf("No product descriptions found for category: %s", categoryName)
			continue
		}

		// Check for images
		images, ok := productImages[categoryName]
		if !ok {
			log.Printf("No product images found for category: %s", categoryName)
			continue
		}

		// Check for attributes
		categoryAttrs, ok := categoryAttributes[categoryName]
		if !ok {
			log.Printf("No category attributes found for category: %s", categoryName)
			continue
		}

		// Tạo sản phẩm cho mỗi nhà cung cấp
		for _, supplierID := range supplierIDs {
			// Mỗi nhà cung cấp tạo 1-3 sản phẩm cho mỗi danh mục
			numProducts := gofakeit.Number(1, 3)

			for i := 0; i < numProducts; i++ {
				// Chọn ngẫu nhiên tên sản phẩm
				productName := productNames[gofakeit.Number(0, len(productNames)-1)]

				// Chọn ngẫu nhiên mô tả sản phẩm
				productDesc := productDescriptions[gofakeit.Number(0, len(productDescriptions)-1)]

				// Chọn ngẫu nhiên ảnh sản phẩm
				productImage := images[gofakeit.Number(0, len(images)-1)]

				// Tạo SKU prefix dựa trên tên sản phẩm và tên danh mục
				skuPrefix := strings.ToUpper(string([]rune(categoryName)[0])) +
					strings.ToUpper(string([]rune(productName)[0])) +
					fmt.Sprintf("%03d", gofakeit.Number(100, 999))

				// Kiểm tra xem sản phẩm đã tồn tại chưa
				var existingID string
				err := db.QueryRow(ctx, `
                    SELECT id FROM products WHERE name = $1 AND supplier_id = $2
                `, productName, supplierID).Scan(&existingID)

				var productID string
				if err != nil && err != pgx.ErrNoRows {
					log.Printf("Error checking product existence: %v", err)
					continue
				}

				if err == pgx.ErrNoRows {
					// Nếu chưa tồn tại, tạo mới
					err := db.QueryRow(ctx, `
                        INSERT INTO products (
                            supplier_id, category_id, name, description, image_url, 
                            status, featured, tax_class, sku_prefix, average_rating
                        )
                        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
                        RETURNING id;
                    `,
						supplierID, categoryID, productName, productDesc, productImage,
						"active", gofakeit.Bool(), "standard", skuPrefix, float32(gofakeit.Float32Range(3.5, 5)),
					).Scan(&productID)

					if err != nil {
						log.Printf("Error inserting product: %v", err)
						continue
					}
				} else {
					productID = existingID
					// Update existing product if needed
					_, err := db.Exec(ctx, `
                        UPDATE products SET 
                        description = $1, image_url = $2, status = 'active', 
                        tax_class = 'standard', sku_prefix = $3
                        WHERE id = $4
                    `, productDesc, productImage, skuPrefix, productID)

					if err != nil {
						log.Printf("Error updating product: %v", err)
						continue
					}
				}

				// Add tags to product
				numTags := gofakeit.Number(1, 3)
				tagNames := []string{"Mới nhất", "Bán chạy", "Chính hãng", "Giảm giá", "Chất lượng cao"}
				for j := 0; j < numTags; j++ {
					randomTag := tagNames[gofakeit.Number(0, len(tagNames)-1)]

					// First get the tag ID
					var tagID string
					err := db.QueryRow(ctx, `
                        SELECT id FROM tags WHERE name = $1
                    `, randomTag).Scan(&tagID)

					if err != nil {
						log.Printf("Error getting tag ID: %v", err)
						continue
					}

					// Check if the product-tag relation already exists
					var relationExists bool
					err = db.QueryRow(ctx, `
                        SELECT EXISTS (
                            SELECT 1 FROM products_tags 
                            WHERE product_id = $1 AND tag_id = $2
                        )
                    `, productID, tagID).Scan(&relationExists)

					if err != nil {
						log.Printf("Error checking product_tag existence: %v", err)
						continue
					}

					if !relationExists {
						_, err := db.Exec(ctx, `
                            INSERT INTO products_tags (product_id, tag_id)
                            VALUES ($1, $2);
                        `, productID, tagID)

						if err != nil {
							log.Printf("Error inserting product_tag: %v", err)
						}
					}
				}

				// Now check if product_variants table exists before trying to create variants
				err = db.QueryRow(ctx, `
                    SELECT EXISTS (
                        SELECT FROM information_schema.tables 
                        WHERE table_schema = 'public' 
                        AND table_name = 'product_variants'
                    )
                `).Scan(&exists)

				if err != nil {
					log.Printf("Error checking product_variants table: %v", err)
					continue
				}

				if exists {
					// Create product variants
					createProductVariants(ctx, db, productID, skuPrefix, categoryAttrs, productImage)
				} else {
					log.Printf("Skipping variant creation as product_variants table doesn't exist yet")
				}

				totalProducts++
			}
		}
	}

	log.Printf("✅ Created %d products with variants", totalProducts)
}

func createProductVariants(
	ctx context.Context,
	db *pgxpool.Pool,
	productID string,
	skuPrefix string,
	categoryAttrs map[string][]string,
	productImage string,
) {
	// Chọn 2 thuộc tính để tạo biến thể
	var variantAttrs []string
	for attrName := range categoryAttrs {
		variantAttrs = append(variantAttrs, attrName)
		if len(variantAttrs) >= 2 {
			break
		}
	}

	// Nếu không đủ thuộc tính, thì bỏ qua
	if len(variantAttrs) < 1 {
		log.Printf("Not enough attributes for product: %s", productID)
		return
	}

	// Lấy thông tin định nghĩa thuộc tính và tùy chọn
	attributeDefs := make(map[string]int)               // name -> id
	attributeOptions := make(map[string]map[string]int) // attribute name -> option value -> id

	// Lấy định nghĩa thuộc tính
	for _, attrName := range variantAttrs {
		var attrID int
		err := db.QueryRow(ctx, `
            SELECT id FROM attribute_definitions WHERE name = $1
        `, attrName).Scan(&attrID)

		if err != nil {
			log.Printf("Error getting attribute definition: %v", err)
			continue
		}

		attributeDefs[attrName] = attrID
		attributeOptions[attrName] = make(map[string]int)

		// Lấy các tùy chọn cho thuộc tính này
		rows, err := db.Query(ctx, `
            SELECT id, option_value FROM attribute_options WHERE attribute_definition_id = $1
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

	// Get a list of existing SKUs for this product to avoid duplicates
	existingSKUs := make(map[string]bool)
	rows, err := db.Query(ctx, `
        SELECT sku FROM product_variants WHERE product_id = $1
    `, productID)

	if err != nil {
		log.Printf("Error checking existing SKUs: %v", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var sku string
			if err := rows.Scan(&sku); err != nil {
				log.Printf("Error scanning SKU: %v", err)
				continue
			}
			existingSKUs[sku] = true
		}
	}

	// Track used attribute options to avoid duplicate constraint violations
	usedAttributeOptions := make(map[int]bool)

	// Tạo biến thể sản phẩm dựa trên thuộc tính đầu tiên
	attrName := variantAttrs[0]
	attrValues := categoryAttrs[attrName]

	for i, attrValue := range attrValues {
		// Bỏ qua nếu không có option_id cho giá trị này
		optionID, ok := attributeOptions[attrName][attrValue]
		if !ok {
			continue
		}

		// Skip if this attribute option has already been used
		if usedAttributeOptions[optionID] {
			continue
		}

		// Tính giá và giá giảm
		basePrice := gofakeit.Float32Range(100000, 5000000) // 100k - 5tr VND
		// Làm tròn giá theo 1000 đồng
		basePrice = float32(math.Round(float64(basePrice/1000)) * 1000)

		discountPrice := basePrice
		hasDiscount := gofakeit.Bool()
		if hasDiscount {
			discountPercent := gofakeit.Float32Range(0.05, 0.3) // Giảm 5% - 30%
			discountPrice = float32(math.Round(float64(basePrice*(1-discountPercent)/1000)) * 1000)
		}

		// Tạo SKU với một unique identifier để tránh trùng lặp
		sku := fmt.Sprintf("%s-%03d-%s", skuPrefix, i+1, uuid.New().String()[:4])

		// Skip if this SKU already exists
		if existingSKUs[sku] {
			continue
		}

		// Tạo tên biến thể
		variantName := fmt.Sprintf("%s - %s", attrName, attrValue)

		// Tạo biến thể sản phẩm
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
                inventory_quantity, shipping_class, image_url, alt_text, is_default
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
            RETURNING id;
        `,
			productID, sku, variantName, basePrice, discountPriceParam,
			gofakeit.Number(5, 100), "standard", productImage, variantName, i == 0,
		).Scan(&variantID)

		if err != nil {
			log.Printf("Error inserting product variant: %v", err)
			continue
		}

		// Mark this attribute option as used
		usedAttributeOptions[optionID] = true

		// First check if this variant already has this attribute option
		var attrExists bool
		err = db.QueryRow(ctx, `
            SELECT EXISTS (
                SELECT 1 FROM product_variant_attributes 
                WHERE product_variant_id = $1 AND attribute_option_id = $2
            )
        `, variantID, optionID).Scan(&attrExists)

		if err != nil {
			log.Printf("Error checking variant attribute existence: %v", err)
			continue
		}

		// Skip if attribute already exists for this variant
		if attrExists {
			continue
		}

		// Thêm thuộc tính cho biến thể
		_, err = db.Exec(ctx, `
			INSERT INTO product_variant_attributes (
				product_variant_id, attribute_definition_id, attribute_option_id
			)
			VALUES ($1, $2, $3);
		`, variantID, attributeDefs[attrName], optionID)

		if err != nil {
			log.Printf("Error inserting product variant attribute: %v", err)
		}
	}
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
				user_id, id_card_number, id_card_front_image, id_card_back_iamge,
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
					product_id, user_id, rating, comment, is_verified_purchase, helpful_votes
				)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT DO NOTHING;
			`, product.id, userID, rating, comment, gofakeit.Bool(), gofakeit.Number(0, 20))

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
		code, name, desc, discountType, appliesTo string
		discountValue, maxDiscount, minOrder      float32
		startDate, endDate                        time.Time
		usageLimit                                int
	}{
		{
			code:          "WELCOME10",
			name:          "Chào mừng thành viên mới",
			desc:          "Giảm 10% cho đơn hàng đầu tiên",
			discountType:  "percentage",
			discountValue: 10,
			maxDiscount:   100000,
			minOrder:      0,
			appliesTo:     "order",
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
			appliesTo:     "order",
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
			appliesTo:     "order",
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
			appliesTo:     "category",
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
			appliesTo:     "order",
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
				start_date, end_date, usage_limit, is_active, applies_to
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, 'VND', $8, $9, $10, TRUE, $11)
			ON CONFLICT (code) DO UPDATE
			SET name = $2, description = $3, discount_type = $4, discount_value = $5,
				maximum_discount_amount = $6, minimum_order_amount = $7,
				start_date = $8, end_date = $9, usage_limit = $10, applies_to = $11;
		`,
			coupon.code, coupon.name, coupon.desc, coupon.discountType, coupon.discountValue,
			coupon.maxDiscount, coupon.minOrder, coupon.startDate, coupon.endDate,
			coupon.usageLimit, coupon.appliesTo)

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
	// Số lượng thông báo mỗi người dùng (3-10)
	numPerUser := 5

	// Các loại thông báo
	types := []int{1, 2, 3, 4, 5} // 1: order, 2: payment, 3: product, 4: promotion, 5: system

	// Tiêu đề và nội dung thông báo
	titles := map[int][]string{
		1: { // Order
			"Đơn hàng đã được xác nhận",
			"Đơn hàng đang được xử lý",
			"Đơn hàng đang được giao",
			"Đơn hàng đã được giao thành công",
		},
		2: { // Payment
			"Thanh toán thành công",
			"Thanh toán đang được xử lý",
			"Yêu cầu thanh toán đơn hàng",
		},
		3: { // Product
			"Sản phẩm đang giảm giá",
			"Sản phẩm bạn quan tâm đã có hàng",
			"Đánh giá sản phẩm đã mua",
		},
		4: { // Promotion
			"Khuyến mãi mùa hè",
			"Flash sale cuối tuần",
			"Mã giảm giá cho thành viên",
			"Ưu đãi đặc biệt dành cho bạn",
		},
		5: { // System
			"Cập nhật thông tin tài khoản",
			"Xác thực tài khoản thành công",
			"Bảo mật tài khoản",
		},
	}

	contents := map[int][]string{
		1: { // Order
			"Đơn hàng #ORDER-ID của bạn đã được xác nhận. Chúng tôi sẽ sớm xử lý đơn hàng.",
			"Đơn hàng #ORDER-ID của bạn đang được xử lý. Dự kiến đơn hàng sẽ được giao trong 3-5 ngày tới.",
			"Đơn hàng #ORDER-ID của bạn đang được giao. Vui lòng chuẩn bị nhận hàng.",
			"Đơn hàng #ORDER-ID của bạn đã được giao thành công. Cảm ơn bạn đã mua sắm!",
		},
		2: { // Payment
			"Thanh toán đơn hàng #ORDER-ID của bạn đã thành công. Cảm ơn bạn!",
			"Thanh toán đơn hàng #ORDER-ID của bạn đang được xử lý. Chúng tôi sẽ thông báo cho bạn khi hoàn tất.",
			"Vui lòng thanh toán đơn hàng #ORDER-ID của bạn trong vòng 24 giờ để tránh bị hủy.",
		},
		3: { // Product
			"Sản phẩm [PRODUCT-NAME] bạn đã xem gần đây đang được giảm giá 20%. Mua ngay!",
			"Sản phẩm [PRODUCT-NAME] bạn quan tâm đã có hàng trở lại. Nhanh tay mua ngay!",
			"Bạn đã mua sản phẩm [PRODUCT-NAME] gần đây. Vui lòng đánh giá sản phẩm để nhận voucher!",
		},
		4: { // Promotion
			"Khuyến mãi mùa hè với hàng ngàn sản phẩm giảm giá lên đến 50%. Khám phá ngay!",
			"Flash sale cuối tuần - Giảm giá sốc chỉ trong 2 giờ. Bắt đầu từ 20:00 tối nay.",
			"Tặng bạn mã giảm giá SUMMER10 giảm 10% cho đơn hàng tiếp theo. Hạn sử dụng 7 ngày.",
			"Ưu đãi đặc biệt cho thành viên thân thiết - Giảm 15% cho các sản phẩm thời trang.",
		},
		5: { // System
			"Thông tin tài khoản của bạn đã được cập nhật thành công.",
			"Tài khoản của bạn đã được xác thực thành công. Bạn có thể sử dụng đầy đủ tính năng của hệ thống.",
			"Vì lý do bảo mật, vui lòng cập nhật mật khẩu của bạn định kỳ.",
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

	// Tạo thông báo cho mỗi người dùng
	totalNotifs := 0

	for _, userID := range userIDs {
		for i := 0; i < numPerUser; i++ {
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
			_, err := db.Exec(ctx, `
				INSERT INTO notifications (
					user_id, type, title, content, is_read, image_title
				)
				VALUES ($1, $2, $3, $4, $5, $6)
				ON CONFLICT DO NOTHING;
			`,
				userID, notifType, title, content,
				gofakeit.Bool(), imageURL)

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
