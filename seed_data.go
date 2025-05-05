package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

// Danh sách các APIs hỗ trợ dữ liệu địa giới hành chính Việt Nam
var vietnamGeoAPIs = []string{
	"https://provinces.open-api.vn/api/?depth=3", // API với đầy đủ phường/xã
	"https://vietnam-administrative-divisions.vercel.app/api/",
	"https://vapi.vnappmob.com/api/province/",
}

// Mở rộng danh sách tên sản phẩm mẫu theo danh mục
var categoryProductNames = map[string][]string{
	"Điện thoại thông minh": {
		"iPhone 13 Pro Max", "iPhone 14", "Samsung Galaxy S22 Ultra", "Samsung Galaxy Z Fold 4",
		"Xiaomi Redmi Note 11", "Xiaomi 12T Pro", "OPPO Reno8 Pro", "OPPO Find X5 Pro",
		"Vivo V25 Pro", "Realme GT Neo 3", "Nokia G21", "Huawei Nova 10",
		"iPhone 15 Pro", "Samsung Galaxy S23", "Google Pixel 7", "OnePlus 11",
		"Xiaomi 13", "Nothing Phone", "Sony Xperia 5 IV", "Asus ROG Phone 7",
	},
	"Máy tính xách tay": {
		"MacBook Air M2", "MacBook Pro 14", "Dell XPS 13", "Dell Inspiron 15",
		"HP Spectre x360", "HP Pavilion 15", "Lenovo ThinkPad X1 Carbon", "Lenovo Yoga 7i",
		"Asus ZenBook 14", "Asus ROG Zephyrus G14", "Acer Swift 5", "MSI Prestige 14",
		"Microsoft Surface Laptop 5", "Razer Blade 15", "LG Gram 17", "Gigabyte Aero 16",
		"Framework Laptop", "Huawei MateBook X Pro", "Samsung Galaxy Book 3", "Alienware m18",
	},
	"Thời trang nam": {
		"Áo sơ mi nam dài tay", "Áo thun nam cổ tròn", "Áo thun polo nam", "Áo khoác denim nam",
		"Áo khoác bomber nam", "Quần jeans nam slim fit", "Quần kaki nam", "Quần short nam",
		"Bộ vest nam công sở", "Áo len nam cổ tròn", "Áo hoodie nam", "Quần tây nam công sở",
		"Áo phông nam cổ V", "Áo khoác gió nam", "Áo blazer nam", "Quần jogger nam",
		"Áo gile nam", "Quần shorts thể thao nam", "Áo thun nam oversize", "Áo sơ mi nam caro",
	},
	"Thời trang nữ": {
		"Áo sơ mi nữ công sở", "Áo blouse nữ", "Áo thun nữ cổ tròn", "Áo khoác denim nữ",
		"Đầm suông nữ", "Đầm ôm body nữ", "Chân váy chữ A", "Chân váy tennis",
		"Quần jeans nữ ống rộng", "Quần culottes nữ", "Áo cardigan nữ", "Set đồ nữ hai mảnh",
		"Áo croptop nữ", "Áo kiểu nữ", "Đầm maxi nữ", "Quần short nữ",
		"Áo thun oversize nữ", "Quần baggy nữ", "Áo khoác blazer nữ", "Chân váy midi",
	},
	"Đồ gia dụng": {
		"Nồi cơm điện", "Máy xay sinh tố", "Bếp từ đơn", "Bếp gas đôi",
		"Lò vi sóng", "Ấm đun nước siêu tốc", "Máy lọc không khí", "Quạt điều hòa",
		"Máy hút bụi", "Bàn ủi hơi nước", "Nồi chiên không dầu", "Máy rửa chén",
		"Máy lọc nước", "Máy ép trái cây", "Máy sấy tóc", "Nồi áp suất điện",
		"Máy đánh trứng", "Bàn là hơi nước đứng", "Máy xay thịt", "Máy làm sữa hạt",
	},
	"Sách": {
		"Nhà Giả Kim", "Đắc Nhân Tâm", "Cà Phê Cùng Tony", "Người Giàu Có Nhất Thành Babylon",
		"Hai Số Phận", "Điều Kỳ Diệu Của Tiệm Tạp Hóa Namiya", "Bước Chậm Lại Giữa Thế Gian Vội Vã", "Tuổi Trẻ Đáng Giá Bao Nhiêu",
		"Chúng Ta Rồi Sẽ Hạnh Phúc, Theo Những Cách Khác Nhau", "Khéo Ăn Nói Sẽ Có Được Thiên Hạ", "Tôi Tài Giỏi, Bạn Cũng Thế", "Dám Nghĩ Lớn",
		"Atomic Habits", "Sapiens", "Thinking Fast and Slow", "Đàn Ông Sao Hỏa Đàn Bà Sao Kim",
		"Tiểu thuyết Sherlock Holmes", "Harry Potter", "Nhật ký Anne Frank", "Nghệ thuật tinh tế của việc đếch quan tâm",
	},
	"Thể thao": {
		"Tạ tay 5kg", "Thảm tập yoga", "Dây nhảy thể dục", "Máy chạy bộ điện",
		"Xe đạp tập thể dục", "Găng tay tập gym", "Ghế tập bụng đa năng", "Bóng đá size 5",
		"Vợt cầu lông", "Vợt tennis", "Bộ cờ vua quốc tế", "Bàn bóng bàn",
		"Xe đạp địa hình", "Gậy golf", "Áo bơi nam", "Áo bơi nữ",
		"Giày chạy bộ", "Bóng rổ", "Dụng cụ bơi lội", "Máy tập thể dục đa năng",
	},
	"Máy tính bảng": {
		"iPad Pro 12.9", "iPad Mini", "Samsung Galaxy Tab S8 Ultra", "Samsung Galaxy Tab A8",
		"Xiaomi Pad 5", "Lenovo Tab P11 Pro", "Huawei MatePad Pro", "Amazon Fire HD 10",
		"Microsoft Surface Pro 9", "Realme Pad", "Nokia T20", "OnePlus Pad",
		"iPad Air", "Redmi Pad", "Vivo Pad", "TCL NXTPAPER 11",
	},
	"Máy tính để bàn": {
		"iMac 24", "Mac Mini M2", "Dell OptiPlex", "HP Pavilion Desktop",
		"Lenovo ThinkCentre", "Asus ProArt", "Acer Aspire TC", "MSI MAG Infinite",
		"HP All-in-One", "Lenovo IdeaCentre AIO", "Microsoft Surface Studio", "Dell XPS Desktop",
		"ASUS ROG Strix", "Corsair One", "Alienware Aurora", "HP Omen",
	},
	"Tai nghe & Loa": {
		"AirPods Pro", "Sony WH-1000XM5", "JBL Flip 6", "Bose QuietComfort Earbuds",
		"Samsung Galaxy Buds Pro", "Harman Kardon Onyx Studio", "Marshall Emberton", "Anker Soundcore",
		"Beats Studio Buds", "Jabra Elite 7 Pro", "Logitech G Pro X", "Sennheiser HD 660S",
		"UE Wonderboom", "Bose SoundLink Revolve+", "Edifier R1280DB", "Audio-Technica ATH-M50x",
	},
	"Máy ảnh & Máy quay": {
		"Sony Alpha A7 IV", "Canon EOS R6", "Nikon Z6 II", "Fujifilm X-T4",
		"GoPro HERO11 Black", "DJI Pocket 2", "Canon EOS 90D", "Sony ZV-1",
		"Panasonic Lumix GH6", "Olympus OM-D E-M10", "Leica Q2", "Ricoh GR III",
		"Canon EOS R5", "Sony A6600", "Nikon Z9", "Fujifilm X100V",
	},
	"Thời trang trẻ em": {
		"Bộ quần áo trẻ em hình thú", "Áo thun bé trai họa tiết siêu nhân", "Váy công chúa bé gái", "Đồ bộ thể thao trẻ em", "Đồ bơi trẻ em", "Áo khoác trẻ em",
		"Quần jeans bé trai", "Giày thể thao trẻ em", "Đầm dự tiệc bé gái", "Quần áo sơ sinh",
		"Áo len trẻ em", "Quần short bé trai", "Váy denim bé gái", "Áo phông trẻ em họa tiết hoạt hình",
		"Bộ đồ ngủ trẻ em", "Mũ bảo hiểm trẻ em", "Túi đeo chéo trẻ em", "Giày sandal trẻ em",
	},
	"Giày dép": {
		"Giày thể thao nam", "Giày cao gót nữ", "Giày tây nam", "Giày sandal nữ",
		"Dép quai hậu nam", "Giày búp bê nữ", "Giày lười nam", "Giày boot nữ",
		"Giày thể thao nữ", "Giày đế xuồng nữ", "Dép tông nam", "Giày oxford nam",
		"Giày sneaker nữ", "Dép lê nam", "Giày mọi nữ", "Giày vải nam",
	},
	"Túi xách": {
		"Túi xách tay nữ", "Túi đeo chéo nữ", "Balo nam", "Balo nữ",
		"Túi đeo hông nam", "Túi clutch nữ", "Túi laptop", "Túi du lịch",
		"Ví nam", "Ví nữ", "Túi tote nữ", "Túi đeo chéo nam",
		"Túi xách công sở nữ", "Túi chống sốc laptop", "Túi đựng mỹ phẩm", "Balo du lịch",
	},
	"Đồng hồ & Trang sức": {
		"Đồng hồ nam", "Đồng hồ nữ", "Nhẫn bạc nữ", "Vòng tay nam",
		"Dây chuyền bạc nữ", "Khuyên tai nữ", "Nhẫn cưới", "Vòng cổ nam",
		"Đồng hồ thông minh", "Lắc tay nữ", "Mặt dây chuyền phật", "Đồng hồ đôi",
		"Nhẫn nam", "Vòng cổ choker", "Bông tai nam", "Nhẫn đính hôn",
	},
	"Đồ dùng nhà bếp": {
		"Bộ nồi inox", "Chảo chống dính", "Dao làm bếp", "Thớt gỗ",
		"Hộp đựng thực phẩm", "Bình giữ nhiệt", "Bát đĩa", "Đũa thìa inox",
		"Máy xay cà phê", "Cây đánh trứng", "Ly thủy tinh", "Nồi lẩu điện",
		"Dao gọt hoa quả", "Ấm trà", "Đồ lọc cà phê", "Dụng cụ làm bánh",
	},
	"Đồ dùng phòng ngủ": {
		"Bộ ga gối đệm", "Chăn lông cừu", "Nệm cao su", "Gối ôm",
		"Gối tựa đầu", "Đệm bông ép", "Vỏ chăn", "Gối latex",
		"Đèn ngủ", "Màn chống muỗi", "Đồng hồ báo thức", "Nệm lò xo",
		"Gối mềm", "Tủ đầu giường", "Rèm cửa phòng ngủ", "Thảm trải sàn",
	},
	"Đồ dùng phòng tắm": {
		"Khăn tắm", "Rèm nhà tắm", "Gương phòng tắm", "Vòi sen",
		"Giá treo khăn", "Thảm chân nhà tắm", "Kệ góc nhà tắm", "Bộ phụ kiện nhà tắm",
		"Bồn rửa mặt", "Vòi rửa mặt", "Máy sấy tóc", "Hộp đựng xà phòng",
		"Giá treo đồ", "Bồn tắm", "Khăn mặt", "Giấy vệ sinh",
	},
	"Văn phòng phẩm": {
		"Bút bi", "Bút gel", "Sổ tay", "Giấy note",
		"Kẹp giấy", "Máy tính cầm tay", "Bút highlight", "Thước kẻ",
		"Hộp bút", "Bút chì", "Gôm tẩy", "Ghim bấm",
		"Kéo văn phòng", "Máy bấm kim", "Giá đỡ tài liệu", "Dập ghim",
	},
	"Sách giáo khoa": {
		"Sách tiếng Việt lớp 1", "Sách Toán lớp 5", "Sách Tiếng Anh lớp 9", "Sách Vật lý lớp 12",
		"Sách Hóa học lớp 11", "Sách Sinh học lớp 10", "Sách Lịch sử lớp 7", "Sách Địa lý lớp 8",
		"Sách Giáo dục công dân", "Sách Tin học", "Sách ôn thi THPT Quốc gia", "Sách ôn thi đại học",
		"Sách tham khảo Toán", "Sách tham khảo Ngữ văn", "Từ điển Anh-Việt", "Sách giáo trình đại học",
	},
	"Tạp chí & Báo": {
		"Tạp chí Thời trang", "Báo Thanh niên", "Tạp chí Đẹp", "Báo Tuổi trẻ",
		"Tạp chí Kinh tế", "Báo Nhân dân", "Tạp chí Sức khỏe", "Báo Pháp luật",
		"Tạp chí Du lịch", "Báo Hà Nội mới", "Tạp chí Kiến trúc", "Báo Tiền phong",
		"Tạp chí Ẩm thực", "Báo Lao động", "Tạp chí Khoa học", "Báo Đầu tư",
	},
	"Dụng cụ thể thao": {
		"Quả bóng đá", "Vợt cầu lông", "Bóng rổ", "Quả bóng chuyền",
		"Gậy golf", "Vợt tennis", "Bàn cờ vua", "Túi đựng đồ thể thao",
		"Túi đựng vợt", "Găng tay tập gym", "Băng đô thể thao", "Khăn lau mồ hôi",
		"Găng tay golf", "Balo thể thao", "Băng bảo vệ cổ tay", "Giàn tập tạ",
	},
	"Vali & Túi du lịch": {
		"Vali kéo 20 inch", "Vali kéo 24 inch", "Vali kéo 28 inch", "Balo du lịch",
		"Túi du lịch", "Túi đựng đồ cá nhân", "Túi đựng giày", "Túi đựng mỹ phẩm",
		"Vali kéo trẻ em", "Balo máy ảnh", "Túi chống nước", "Gói hành lý",
		"Khóa vali", "Thẻ hành lý", "Gối cổ du lịch", "Túi treo quần áo",
	},
	"Thiết bị cắm trại": {
		"Lều cắm trại", "Túi ngủ", "Đèn cắm trại", "Bếp dã ngoại",
		"Ghế xếp", "Bàn xếp", "Dao đa năng", "Thảm picnic",
		"Túi giữ lạnh", "Bình nước", "Võng dã ngoại", "Bếp cồn",
		"Mũ cắm trại", "Áo mưa", "Dụng cụ nấu ăn cắm trại", "Đèn pin",
	},
	"Xe đạp & Phụ kiện": {
		"Xe đạp đường phố", "Xe đạp thể thao", "Xe đạp trẻ em", "Xe đạp điện",
		"Mũ bảo hiểm xe đạp", "Găng tay xe đạp", "Đèn xe đạp", "Bơm xe đạp",
		"Khóa xe đạp", "Túi treo xe đạp", "Chuông xe đạp", "Gương xe đạp",
		"Kính xe đạp", "Giỏ xe đạp", "Bình nước xe đạp", "Găng tay xe đạp",
	},
	"Đồ dùng cho bé": {
		"Tã dán", "Tã quần", "Sữa bột", "Bình sữa",
		"Núm ti", "Ghế ăn dặm", "Xe tập đi", "Nôi em bé",
		"Xe đẩy em bé", "Miếng lót thấm sữa", "Máy hâm sữa", "Máy tiệt trùng bình sữa",
		"Khăn ướt", "Phấn rôm", "Dầu gội em bé", "Sữa tắm em bé",
	},
	"Đồ chơi cho bé": {
		"Đồ chơi xếp hình", "Thú nhồi bông", "Đồ chơi gỗ", "Đồ chơi âm nhạc",
		"Đồ chơi nhà tắm", "Xe ô tô đồ chơi", "Búp bê", "Đồ chơi giáo dục",
		"Đồ chơi ngoài trời", "Đồ chơi nấu ăn", "Đồ chơi bác sĩ", "Lego",
		"Đồ chơi cát", "Đồ chơi câu cá", "Đồ chơi nhà bếp", "Đồ chơi vận động",
	},
	"Thời trang cho bé": {
		"Bộ bodysuit", "Áo sơ sinh", "Quần sơ sinh", "Mũ sơ sinh",
		"Áo thun bé trai", "Áo thun bé gái", "Quần jeans bé trai", "Váy bé gái",
		"Bộ đồ ngủ bé trai", "Bộ đồ ngủ bé gái", "Giày tập đi", "Tất em bé",
		"Giày cho bé", "Đồ bơi bé trai", "Đồ bơi bé gái", "Mũ lưỡi trai trẻ em",
	},
	"Đồ dùng cho mẹ": {
		"Áo cho con bú", "Máy hút sữa", "Túi trữ sữa", "Đai đỡ bụng bầu",
		"Kem chống rạn da", "Sữa rửa mặt cho bà bầu", "Gối bầu", "Quần bầu",
		"Áo bầu", "Túi đựng đồ cho mẹ và bé", "Ghế massage cho bà bầu", "Vitamin cho bà bầu",
		"Sữa bầu", "Túi giữ nhiệt sữa", "Kem nứt đầu ti", "Áo lót cho con bú",
	},
	"Mỹ phẩm": {
		"Son môi", "Phấn má hồng", "Kem nền", "Mascara",
		"Phấn phủ", "Kẻ mắt", "Phấn mắt", "Kem lót",
		"Tẩy trang", "Sữa rửa mặt", "Kem chống nắng", "Xịt khoáng",
		"Kem dưỡng ẩm", "Mặt nạ dưỡng da", "Serum dưỡng da", "Nước hoa",
	},
	"Chăm sóc da": {
		"Sữa rửa mặt", "Toner", "Serum", "Kem dưỡng ẩm",
		"Kem chống nắng", "Mặt nạ giấy", "Mặt nạ ngủ", "Tẩy tế bào chết",
		"Kem mắt", "Kem trị mụn", "Kem trị thâm", "Kem dưỡng body",
		"Sữa tắm", "Dầu gội", "Dầu xả", "Kem tẩy lông",
	},
	"Chăm sóc tóc": {
		"Dầu gội", "Dầu xả", "Kem ủ tóc", "Serum dưỡng tóc",
		"Dầu dưỡng tóc", "Xịt dưỡng tóc", "Thuốc nhuộm tóc", "Máy sấy tóc",
		"Máy duỗi tóc", "Máy uốn tóc", "Lược chải tóc", "Kẹp tóc",
		"Mũ trùm tóc", "Khăn quấn tóc", "Dầu gội khô", "Wax tóc",
	},
	"Thực phẩm chức năng": {
		"Viên uống vitamin C", "Viên uống collagen", "Viên uống sáng da", "Viên uống chống nắng",
		"Viên uống giảm cân", "Viên uống bổ gan", "Viên uống canxi", "Viên uống omega 3",
		"Viên uống vitamin tổng hợp", "Viên uống probiotics", "Viên uống sữa ong chúa", "Viên uống tăng cường miễn dịch",
		"Viên uống mọc tóc", "Viên uống ngừa nám", "Viên uống bổ mắt", "Viên uống tăng chiều cao",
	},
	"Thiết bị y tế": {
		"Máy đo huyết áp", "Máy đo đường huyết", "Máy massage", "Nhiệt kế",
		"Máy xông mũi họng", "Máy hút mũi", "Dụng cụ rửa mũi", "Hộp đựng thuốc",
		"Băng gạc y tế", "Cồn y tế", "Bông gòn y tế", "Dung dịch sát khuẩn",
		"Khẩu trang y tế", "Ống nghe y tế", "Máy đo nồng độ oxy", "Găng tay y tế",
	},
	"Thực phẩm khô": {
		"Gạo", "Mì gói", "Hạt nêm", "Dầu ăn",
		"Nước mắm", "Nước tương", "Đường", "Muối",
		"Bột ngọt", "Hạt tiêu", "Bột canh", "Tương ớt",
		"Cà phê hòa tan", "Trà túi lọc", "Thực phẩm ăn liền", "Ngũ cốc ăn sáng",
	},
	"Thực phẩm tươi sống": {
		"Rau xanh", "Trái cây", "Thịt heo", "Thịt bò",
		"Thịt gà", "Cá", "Tôm", "Mực",
		"Trứng gà", "Sữa tươi", "Phô mai", "Bơ",
		"Nấm", "Đậu hũ", "Rau củ quả", "Hải sản tươi sống",
	},
	"Đồ uống": {
		"Nước khoáng", "Nước ngọt", "Trà xanh", "Cà phê",
		"Sữa tươi", "Sữa chua", "Nước ép trái cây", "Bia",
		"Rượu vang", "Rượu mạnh", "Sinh tố", "Nước detox",
		"Nước tăng lực", "Trà sữa", "Nước uống collagen", "Sữa hạt",
	},
	"Bánh kẹo & Đồ ăn vặt": {
		"Bánh quy", "Bánh gạo", "Kẹo", "Chocolate",
		"Snack khoai tây", "Hạt các loại", "Mứt", "Trái cây sấy",
		"Bánh trung thu", "Bánh mì", "Xúc xích", "Pate",
		"Mì ăn liền", "Cháo ăn liền", "Thịt khô", "Chà bông",
	},
}

// Mở rộng danh sách mô tả sản phẩm mẫu theo danh mục
var categoryProductDescriptions = map[string][]string{
	"Điện thoại thông minh": {
		"Sản phẩm công nghệ hiện đại với màn hình Retina sắc nét, camera độ phân giải cao và thời lượng pin dài.",
		"Thiết kế sang trọng với cấu hình mạnh mẽ, camera AI thông minh và khả năng chống nước IP68.",
		"Smartphone cao cấp với chip xử lý mới nhất, màn hình AMOLED 120Hz và sạc nhanh 65W.",
		"Điện thoại thông minh với camera chuyên nghiệp, khả năng quay video 4K và bộ nhớ lớn.",
		"Thiết bị di động mỏng nhẹ, màn hình rộng với tần số quét cao và hỗ trợ công nghệ 5G mới nhất.",
		"Điện thoại thông minh cao cấp với khả năng chụp ảnh đẹp dưới mọi điều kiện ánh sáng và pin siêu bền.",
		"Smartphone thiết kế gập độc đáo với màn hình linh hoạt, hiệu năng mạnh mẽ và khả năng đa nhiệm tuyệt vời.",
		"Điện thoại cao cấp với màn hình tràn viền, cảm biến vân tay dưới màn hình và camera selfie ẩn dưới màn hình.",
	},
	"Máy tính xách tay": {
		"Laptop mỏng nhẹ với hiệu suất mạnh mẽ, thời lượng pin cả ngày và màn hình Retina sắc nét.",
		"Máy tính xách tay cao cấp dành cho công việc sáng tạo, đồ họa với card màn hình rời và SSD tốc độ cao.",
		"Laptop chuyên gaming với card đồ họa mạnh mẽ, tản nhiệt hiệu quả và bàn phím RGB.",
		"Máy tính 2-in-1 linh hoạt với màn hình cảm ứng, bút stylus và thiết kế gập xoay 360 độ.",
		"Laptop siêu mỏng nhẹ, thiết kế tinh tế với hiệu năng ổn định, phù hợp cho công việc và giải trí hàng ngày.",
		"Máy tính xách tay với màn hình HDR sắc nét, âm thanh vòm sống động và thời lượng pin suốt cả ngày.",
		"Laptop chuyên đồ họa với màn hình hiển thị màu chuẩn, cấu hình mạnh mẽ và nhiều cổng kết nối.",
		"Máy tính xách tay business cao cấp với tính năng bảo mật toàn diện, bền bỉ và hiệu năng ổn định.",
	},
	"Thời trang nam": {
		"Áo thời trang nam phong cách Hàn Quốc, chất liệu cao cấp thoáng mát và thấm hút mồ hôi tốt.",
		"Quần nam thiết kế hiện đại, form dáng vừa vặn tôn dáng người mặc và dễ phối đồ.",
		"Sản phẩm thời trang dành cho nam giới công sở với thiết kế lịch lãm, tinh tế và sang trọng.",
		"Trang phục nam thiết kế theo phong cách đường phố, cá tính và năng động dành cho giới trẻ.",
		"Áo nam chất liệu cotton organic cao cấp, không gây kích ứng da, thiết kế basic dễ phối đồ.",
		"Quần nam form slimfit, chất liệu co giãn thoải mái, phù hợp cho cả môi trường công sở và dạo phố.",
		"Áo khoác nam thiết kế thời thượng, chống thấm nước nhẹ và giữ ấm tốt trong mùa lạnh.",
		"Vest nam may đo tỉ mỉ từ chất liệu cao cấp, cắt may tinh tế tôn dáng người mặc.",
	},
	"Thời trang nữ": {
		"Thời trang nữ thiết kế theo xu hướng mới nhất, tôn dáng người mặc và phù hợp nhiều hoàn cảnh.",
		"Trang phục nữ phong cách Hàn Quốc với chất liệu cao cấp, thoáng mát và thấm hút tốt.",
		"Quần áo nữ thiết kế tinh tế với họa tiết độc đáo, phù hợp cho công sở và dạo phố.",
		"Đầm nữ thiết kế sang trọng, quyến rũ phù hợp cho các buổi tiệc và sự kiện quan trọng.",
		"Trang phục nữ phong cách tối giản, thanh lịch với đường may tinh tế và chất liệu bền đẹp.",
		"Đầm nữ thiết kế hiện đại, cut-out tinh tế, tôn dáng và phù hợp với nhiều vóc dáng khác nhau.",
		"Áo nữ thiết kế theo xu hướng Y2K, mang đậm phong cách retro nhưng vẫn hiện đại và trẻ trung.",
		"Quần nữ cạp cao, ôm dáng vừa phải, tôn lên đường cong cơ thể một cách tinh tế và thanh lịch.",
	},
	"Đồ gia dụng": {
		"Thiết bị gia dụng cao cấp với công nghệ hiện đại, tiết kiệm điện và dễ dàng sử dụng.",
		"Sản phẩm gia dụng thông minh với khả năng kết nối điện thoại và điều khiển từ xa.",
		"Thiết bị nhà bếp đa năng với nhiều chức năng, giúp việc nấu nướng trở nên đơn giản và nhanh chóng.",
		"Sản phẩm gia dụng bền bỉ với chất liệu cao cấp và chế độ bảo hành dài hạn.",
		"Thiết bị gia dụng thiết kế nhỏ gọn, tiết kiệm không gian nhưng vẫn đảm bảo đầy đủ tính năng cần thiết.",
		"Sản phẩm gia dụng an toàn với trẻ em, có chế độ tự ngắt khi quá nhiệt và bảo vệ điện áp.",
		"Thiết bị gia dụng cao cấp với thiết kế sang trọng, là điểm nhấn tô điểm cho không gian sống hiện đại.",
		"Đồ gia dụng thông minh với khả năng học hỏi thói quen sử dụng và tự động điều chỉnh cho phù hợp.",
	},
	"Sách": {
		"Cuốn sách best-seller với nội dung sâu sắc, đem lại nhiều bài học giá trị cho người đọc.",
		"Tác phẩm nổi tiếng của tác giả được yêu thích, đã được dịch ra nhiều thứ tiếng trên thế giới.",
		"Sách hay với nội dung bổ ích, ngôn từ cuốn hút và thông điệp ý nghĩa.",
		"Cuốn sách giúp bạn thay đổi tư duy, phát triển bản thân và đạt được thành công trong cuộc sống.",
		"Tác phẩm văn học kinh điển đã được tái bản nhiều lần với bản dịch mới mang tính học thuật cao.",
		"Sách phát triển bản thân với phương pháp thực tế, dễ áp dụng và mang lại hiệu quả rõ rệt.",
		"Tiểu thuyết lãng mạn với cốt truyện cuốn hút, nhân vật sống động và thông điệp nhân văn sâu sắc.",
		"Sách chuyên ngành với nội dung chuyên sâu, cập nhật kiến thức mới nhất trong lĩnh vực.",
	},
	"Thể thao": {
		"Thiết bị tập thể thao cao cấp với chất liệu bền bỉ, an toàn và hiệu quả cao.",
		"Dụng cụ thể thao đa năng giúp bạn tập luyện nhiều nhóm cơ khác nhau.",
		"Sản phẩm thể thao chuyên nghiệp được thiết kế bởi các chuyên gia hàng đầu.",
		"Thiết bị tập luyện tại nhà tiện lợi, tiết kiệm không gian và dễ dàng cất gọn.",
		"Dụng cụ thể thao chuyên nghiệp được sử dụng bởi các vận động viên Olympic và các giải đấu lớn.",
		"Thiết bị tập luyện thông minh với khả năng theo dõi tiến trình, nhịp tim và lượng calo tiêu thụ.",
		"Dụng cụ thể thao ngoài trời bền bỉ trong mọi điều kiện thời tiết, dễ dàng mang theo khi đi du lịch.",
		"Thiết bị thể thao hiện đại, thiết kế tinh tế và công năng vượt trội so với các sản phẩm thông thường.",
		"Dụng cụ tập luyện phù hợp cho mọi đối tượng từ người mới bắt đầu đến vận động viên chuyên nghiệp.",
	},
	"Máy tính bảng": {
		"Máy tính bảng hiện đại với màn hình Retina sắc nét, hiệu năng mạnh mẽ và thời lượng pin dài.",
		"Thiết bị di động đa năng phù hợp cho giải trí, làm việc và học tập với màn hình lớn.",
		"Máy tính bảng cao cấp hỗ trợ bút cảm ứng, phù hợp với các công việc thiết kế và ghi chú.",
		"Thiết bị di động mỏng nhẹ, màn hình sắc nét, âm thanh sống động phù hợp cho giải trí di động.",
		"Máy tính bảng với màn hình hiển thị True Tone, tần số quét cao và khả năng chống chói vượt trội.",
		"Thiết bị di động linh hoạt với khả năng biến đổi thành laptop khi kết nối với bàn phím chuyên dụng.",
		"Máy tính bảng siêu nhẹ, thiết kế tinh tế với khả năng xử lý đa nhiệm mạnh mẽ và ổn định.",
		"Thiết bị giải trí di động cao cấp với màn hình AMOLED, loa stereo và hỗ trợ các ứng dụng giải trí.",
	},
	"Máy tính để bàn": {
		"Máy tính để bàn hiệu năng cao với bộ vi xử lý mới nhất, dung lượng RAM lớn, phù hợp cho gaming và đồ họa.",
		"PC văn phòng nhỏ gọn, thiết kế tối giản với hiệu năng ổn định cho công việc hằng ngày.",
		"Máy tính All-in-One tiện lợi, tiết kiệm không gian với màn hình lớn và âm thanh chất lượng.",
		"PC gaming cao cấp với card đồ họa mạnh mẽ, tản nhiệt hiệu quả và hệ thống LED RGB đẹp mắt.",
		"Máy tính để bàn chuyên dụng cho công việc đồ họa, render video với hiệu năng mạnh mẽ và độ ổn định cao.",
		"PC workstation chuyên nghiệp với khả năng nâng cấp linh hoạt và hiệu năng xử lý đa nhân vượt trội.",
		"Máy tính để bàn mini nhỏ gọn, tiết kiệm không gian nhưng vẫn đảm bảo hiệu năng cho công việc văn phòng.",
		"PC gaming đa nhiệm với khả năng vừa chơi game vừa livestream mượt mà không giật lag.",
	},
	"Tai nghe & Loa": {
		"Tai nghe không dây với công nghệ chống ồn chủ động, âm thanh HD và thời lượng pin lên đến 30 giờ.",
		"Loa bluetooth di động chống nước, âm thanh mạnh mẽ và pin sử dụng liên tục suốt 24 giờ.",
		"Tai nghe gaming với âm thanh vòm 7.1, micro khử tiếng ồn và đèn LED RGB tùy chỉnh.",
		"Loa soundbar cao cấp kết nối không dây, âm thanh vòm Dolby Atmos và thiết kế sang trọng.",
		"Tai nghe true wireless với kết nối Bluetooth 5.2, chống nước IPX7 và hộp sạc không dây tiện lợi.",
		"Loa thông minh tích hợp trợ lý ảo, âm thanh 360 độ và khả năng điều khiển thiết bị nhà thông minh.",
		"Tai nghe audiophile với driver planar magnetic, tái tạo âm thanh chi tiết và không gian âm rộng.",
		"Loa bookshelf cao cấp với củ loa tweeter mềm, âm trầm mạnh mẽ và dải tần số rộng.",
	},
	"Máy ảnh & Máy quay": {
		"Máy ảnh mirrorless full-frame với cảm biến độ phân giải cao, chống rung trong thân máy và khả năng quay video 4K.",
		"Máy quay phim chuyên nghiệp với cảm biến lớn, khả năng quay slow-motion và hệ thống lấy nét nhanh chính xác.",
		"Máy ảnh compact cao cấp nhỏ gọn với zoom quang học lớn, cảm biến 1 inch và khả năng chụp RAW.",
		"Action camera chống nước, chống rung điện tử và khả năng quay video 5.3K với góc nhìn siêu rộng.",
		"Máy ảnh DSLR chuyên nghiệp với hệ thống lấy nét tiên tiến, tốc độ chụp liên tiếp cao và dải ISO rộng.",
		"Máy quay Gimbal tích hợp với khả năng chống rung 3 trục, theo dõi chủ thể và tính năng timelapse.",
		"Máy ảnh medium format với cảm biến lớn, tái tạo màu sắc chính xác và dải tương phản động cao.",
		"Drone quay phim với camera gimbal 3 trục, quay video 4K HDR và khả năng bay ổn định trong nhiều điều kiện.",
	},
	"Phụ tùng ô tô": {
		"Lốp xe ô tô cao cấp với độ bám đường tốt, chống ồn và tuổi thọ cao.",
		"Dầu nhớt động cơ tổng hợp hoàn toàn, bảo vệ động cơ tối ưu và kéo dài thời gian thay dầu.",
		"Ắc quy khô không bảo dưỡng, khởi động mạnh mẽ và tuổi thọ cao trong mọi điều kiện thời tiết.",
		"Bộ lọc không khí, dầu và xăng chính hãng, đảm bảo hiệu suất tối ưu cho động cơ.",
		"Phanh đĩa và má phanh cao cấp với khả năng phanh mạnh mẽ, ổn định và ít tiếng ồn.",
		"Đèn pha LED siêu sáng với tuổi thọ cao và tiêu thụ điện năng thấp.",
		"Phụ tùng điện tử chính hãng với độ bền cao và tương thích hoàn hảo với xe.",
		"Bộ phụ kiện nâng cấp hiệu suất với khả năng tăng mã lực và tiết kiệm nhiên liệu.",
	},
	"Tủ lạnh & Tủ đông": {
		"Tủ lạnh Side-by-Side với ngăn đá lớn, công nghệ làm lạnh đa chiều và tính năng lấy nước, đá tự động.",
		"Tủ lạnh Inverter tiết kiệm điện, vận hành êm ái và duy trì nhiệt độ ổn định.",
		"Tủ lạnh mini nhỏ gọn, phù hợp cho phòng ngủ, văn phòng và những không gian hạn chế.",
		"Tủ đông dung tích lớn với khả năng làm đông nhanh, tiết kiệm điện năng và hoạt động ổn định.",
		"Tủ lạnh French Door sang trọng với ngăn chuyển đổi nhiệt độ linh hoạt và hệ thống khử mùi.",
		"Tủ lạnh thông minh kết nối Wi-Fi, quản lý thực phẩm và điều khiển từ xa qua smartphone.",
		"Tủ đông đứng với nhiều ngăn kéo tiện lợi, dễ dàng sắp xếp và tìm kiếm thực phẩm.",
		"Tủ mát trưng bày đồ uống với cửa kính trong suốt và hệ thống đèn LED trang trí bắt mắt.",
	},
	"Máy giặt & Máy sấy": {
		"Máy giặt cửa trước với công nghệ giặt hơi nước, diệt khuẩn và làm mềm vải hiệu quả.",
		"Máy giặt Inverter tiết kiệm điện, nước với khả năng cân chỉnh tự động lượng nước và chất tẩy.",
		"Máy sấy tụ hơi thông minh với nhiều chương trình sấy chuyên biệt cho từng loại vải.",
		"Máy giặt sấy kết hợp tiết kiệm không gian với công nghệ sấy bằng bơm nhiệt tiết kiệm điện.",
		"Máy giặt cửa trên dung tích lớn phù hợp cho gia đình đông người với khả năng giặt mạnh mẽ.",
		"Máy sấy thông minh với cảm biến độ ẩm, tự động điều chỉnh thời gian sấy phù hợp.",
		"Máy giặt mini nhỏ gọn phù hợp cho căn hộ, nhà trọ với khả năng tiết kiệm điện, nước hiệu quả.",
		"Bộ đôi máy giặt và máy sấy cùng thương hiệu với thiết kế đồng bộ và khả năng kết nối thông minh.",
	},
	"Nội thất phòng khách": {
		"Sofa da cao cấp với khung gỗ tự nhiên, đệm mút D40 êm ái và kiểu dáng hiện đại, sang trọng.",
		"Bàn trà gỗ tự nhiên thiết kế tinh tế, bề mặt chống trầy xước và chân bàn chắc chắn.",
		"Kệ tivi gỗ công nghiệp phủ melamine chống xước, chống ẩm với nhiều ngăn chứa đồ tiện lợi.",
		"Thảm trang trí lông ngắn mềm mại, họa tiết hiện đại và dễ dàng vệ sinh, làm sạch.",
		"Sofa góc L rộng rãi bọc vải cao cấp kháng bẩn, kháng nước và dễ dàng tháo lắp vệ sinh.",
		"Bàn console trang trí phong cách Bắc Âu với thiết kế tối giản và tinh tế.",
		"Ghế bành thư giãn có tính năng ngả lưng và gác chân, bọc da công nghiệp bền đẹp.",
		"Kệ trang trí đa năng với nhiều ngăn kệ phù hợp trưng bày đồ trang trí và sách.",
	},
	"Vật tư nông nghiệp": {
		"Phân bón NPK cân đối dinh dưỡng, tăng cường năng suất cây trồng và cải thiện chất lượng đất.",
		"Hạt giống rau sạch nhập khẩu với tỷ lệ nảy mầm cao và khả năng kháng bệnh tốt.",
		"Thuốc bảo vệ thực vật an toàn, hiệu quả với đa dạng công dụng diệt trừ sâu bệnh.",
		"Màng phủ nông nghiệp chất lượng cao, chống UV, giữ ẩm và kiểm soát cỏ dại hiệu quả.",
		"Hệ thống tưới nhỏ giọt tiết kiệm nước, tưới đúng chỗ và dễ dàng lắp đặt.",
		"Giá thể trồng cây không đất sạch sẽ, thoáng khí và giàu dinh dưỡng cho cây phát triển tốt.",
		"Vật tư làm vườn đồng bộ từ chậu, đất, phân bón đến dụng cụ chăm sóc cây.",
		"Hệ thống nhà kính mini phù hợp cho sân thượng, ban công với khả năng lắp đặt dễ dàng.",
	},
}

// Mở rộng danh sách attribute cho các danh mục mới
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
	// Thêm attribute cho các danh mục mới
	"Máy tính bảng": {
		"Màu sắc":      {"Đen", "Trắng", "Bạc", "Xám", "Vàng Hồng", "Xanh", "Tím"},
		"Dung lượng":   {"32GB", "64GB", "128GB", "256GB", "512GB", "1TB"},
		"Kết nối":      {"WiFi", "4G/LTE", "5G", "WiFi + Cellular"},
		"Kích thước":   {"7.9 inch", "8.3 inch", "10.2 inch", "10.9 inch", "11 inch", "12.9 inch"},
		"Hệ điều hành": {"iPadOS", "Android", "Windows", "HarmonyOS"},
	},
	"Máy tính để bàn": {
		"CPU":          {"Intel Core i3", "Intel Core i5", "Intel Core i7", "Intel Core i9", "AMD Ryzen 5", "AMD Ryzen 7", "AMD Ryzen 9", "AMD Threadripper"},
		"RAM":          {"4GB", "8GB", "16GB", "32GB", "64GB", "128GB"},
		"Ổ cứng":       {"256GB SSD", "512GB SSD", "1TB SSD", "2TB SSD", "1TB HDD", "2TB HDD", "SSD + HDD"},
		"Card đồ họa":  {"NVIDIA GTX 1650", "NVIDIA RTX 3060", "NVIDIA RTX 3070", "NVIDIA RTX 4080", "AMD Radeon RX 6600", "AMD Radeon RX 6800", "Intel Arc"},
		"Hệ điều hành": {"Windows 11", "Windows 10", "macOS", "Linux"},
	},
	"Tai nghe & Loa": {
		"Loại kết nối":   {"Có dây", "Bluetooth", "Wireless 2.4GHz", "Type-C"},
		"Kiểu đeo":       {"Over-ear", "On-ear", "In-ear", "True Wireless"},
		"Tính năng":      {"Chống ồn chủ động", "Chống nước", "Micro đàm thoại", "Âm thanh vòm"},
		"Thời lượng pin": {"8 giờ", "15 giờ", "24 giờ", "30 giờ", "40 giờ", "50 giờ"},
		"Công suất":      {"5W", "10W", "20W", "30W", "50W", "100W", "300W"},
	},
	"Máy ảnh & Máy quay": {
		"Độ phân giải":        {"12MP", "20MP", "24MP", "32MP", "45MP", "61MP"},
		"Cảm biến":            {"MFT", "APS-C", "Full Frame", "Medium Format"},
		"Khả năng quay video": {"Full HD", "4K/30p", "4K/60p", "8K/30p", "RAW Video"},
		"Hãng sản xuất":       {"Sony", "Canon", "Nikon", "Fujifilm", "Panasonic", "Leica", "Olympus"},
		"Loại máy":            {"DSLR", "Mirrorless", "Compact", "Action Camera", "Cinema Camera"},
	},
	"Nội thất phòng khách": {
		"Chất liệu khung": {"Gỗ tự nhiên", "Gỗ công nghiệp", "Kim loại", "Nhựa cao cấp", "Kết hợp"},
		"Chất liệu bọc":   {"Da thật", "Da công nghiệp", "Vải", "Nỉ", "Nhung", "Canvas"},
		"Phong cách":      {"Hiện đại", "Tân cổ điển", "Cổ điển", "Scandinavian", "Industrial", "Minimalist"},
		"Màu sắc":         {"Đen", "Trắng", "Xám", "Nâu", "Be", "Xanh dương", "Xanh lá", "Đỏ"},
		"Kích thước":      {"Nhỏ (1-2 người)", "Vừa (3-4 người)", "Lớn (5-7 người)", "Rất lớn (8+ người)"},
	},
	"Đồ chơi trẻ em": {
		"Độ tuổi phù hợp": {"0-12 tháng", "1-3 tuổi", "3-5 tuổi", "6-9 tuổi", "10-14 tuổi"},
		"Chất liệu":       {"Nhựa an toàn", "Gỗ tự nhiên", "Vải", "Silicone", "Bông"},
		"Loại đồ chơi":    {"Đồ chơi giáo dục", "Đồ chơi vận động", "Đồ chơi sáng tạo", "Đồ chơi nhập vai", "Đồ chơi xây dựng"},
		"Thương hiệu":     {"Lego", "Fisher-Price", "Barbie", "Hot Wheels", "Nerf", "Vtech", "Melissa & Doug"},
		"Xuất xứ":         {"Việt Nam", "Trung Quốc", "Nhật Bản", "Mỹ", "Đức", "Đan Mạch"},
	},
	"Mỹ phẩm": {
		"Loại da":          {"Da dầu", "Da khô", "Da hỗn hợp", "Da nhạy cảm", "Da thường"},
		"Chứng nhận":       {"Organic", "Cruelty-free", "Vegan", "Không paraben", "Hypoallergenic"},
		"Xuất xứ":          {"Hàn Quốc", "Nhật Bản", "Pháp", "Mỹ", "Việt Nam", "Thái Lan"},
		"Hiệu quả":         {"Dưỡng ẩm", "Chống lão hóa", "Trị mụn", "Sáng da", "Chống nắng"},
		"Thành phần chính": {"Vitamin C", "Retinol", "Hyaluronic Acid", "Niacinamide", "AHA/BHA", "Centella Asiatica"},
	},
	"Thực phẩm khô": {
		"Hạn sử dụng":          {"3 tháng", "6 tháng", "1 năm", "2 năm", "3 năm"},
		"Xuất xứ":              {"Việt Nam", "Thái Lan", "Nhật Bản", "Hàn Quốc", "Trung Quốc", "Đài Loan"},
		"Quy cách đóng gói":    {"100g", "250g", "500g", "1kg", "5kg", "Hộp", "Túi", "Lon"},
		"Phương pháp chế biến": {"Sấy khô", "Đông lạnh", "Đóng hộp", "Lên men", "Ướp muối"},
		"Chứng nhận":           {"Organic", "Non-GMO", "Fair Trade", "Halal", "Kosher", "HACCP"},
	},
	"Đồ uống": {
		"Loại đồ uống": {"Nước khoáng", "Nước ngọt", "Cà phê", "Trà", "Nước ép", "Sữa", "Bia", "Rượu"},
		"Dung tích":    {"250ml", "330ml", "500ml", "1L", "1.5L", "2L", "5L"},
		"Vị":           {"Truyền thống", "Trái cây", "Sữa", "Socola", "Vani", "Caramel", "Không đường"},
		"Đóng gói":     {"Chai", "Lon", "Hộp", "Bịch", "Thùng"},
		"Độ cồn":       {"0%", "4.5%", "5%", "12%", "14%", "40%"},
		"Xuất xứ":      {"Việt Nam", "Thái Lan", "Hàn Quốc", "Nhật Bản", "Pháp", "Ý", "Mỹ"},
	},
	"Túi xách": {
		"Chất liệu":  {"Da thật", "Da PU", "Vải Canvas", "Vải Oxford", "Nylon", "Nhựa", "Kim loại"},
		"Kích thước": {"Mini", "Nhỏ", "Trung bình", "Lớn", "Rất lớn"},
		"Màu sắc":    {"Đen", "Trắng", "Nâu", "Be", "Đỏ", "Xanh", "Hồng", "Vàng", "Bạc", "Đa màu"},
		"Phong cách": {"Casual", "Business", "Party", "Vintage", "Sporty", "Minimalist"},
		"Kiểu dáng":  {"Tote", "Crossbody", "Backpack", "Clutch", "Hobo", "Bucket", "Satchel"},
	},
	"Nội thất phòng ngủ": {
		"Chất liệu":         {"Gỗ tự nhiên", "Gỗ công nghiệp", "Kim loại", "Da", "Vải", "Nhựa"},
		"Kích thước giường": {"1.2m x 2m", "1.5m x 2m", "1.6m x 2m", "1.8m x 2m", "2m x 2.2m"},
		"Độ cứng nệm":       {"Cứng", "Trung bình", "Mềm", "Siêu mềm"},
		"Phong cách":        {"Hiện đại", "Cổ điển", "Tân cổ điển", "Vintage", "Minimalist", "Rustic"},
		"Màu sắc":           {"Trắng", "Đen", "Xám", "Nâu", "Be", "Xanh", "Hồng", "Tím"},
	},
	"Cây cảnh & Hoa": {
		"Loại cây":       {"Cây để bàn", "Cây sàn", "Cây treo", "Cây thủy sinh", "Cây nội thất", "Cây ăn quả mini", "Hoa"},
		"Điều kiện sống": {"Ít ánh sáng", "Ánh sáng vừa", "Nhiều ánh sáng", "Ít nước", "Nhiều nước", "Ẩm cao"},
		"Kích thước":     {"Nhỏ (<30cm)", "Trung bình (30-80cm)", "Lớn (80-150cm)", "Rất lớn (>150cm)"},
		"Chậu cây":       {"Nhựa", "Gốm sứ", "Đất nung", "Gỗ", "Kim loại", "Thủy tinh"},
		"Công dụng":      {"Trang trí", "Lọc không khí", "Phong thủy", "Ăn quả", "Làm thuốc"},
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

// Cập nhật hàm seedCategories để thêm nhiều danh mục cha và con hơn
func seedCategories(ctx context.Context, db *pgxpool.Pool) {
	// Danh mục chính (mở rộng từ 5 lên hơn 10)
	mainCategories := []struct{ name, desc, imageUrl string }{
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
		// Thêm các danh mục chính mới
		{
			"Mẹ & Bé",
			"Sản phẩm dành cho mẹ và trẻ em",
			"https://images.unsplash.com/photo-1518531933037-91b2f5f229cc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Làm đẹp & Sức khỏe",
			"Mỹ phẩm, chăm sóc cá nhân và sức khỏe",
			"https://images.unsplash.com/photo-1571875257727-256c39da42af?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Thực phẩm & Đồ uống",
			"Thực phẩm, đồ uống và nguyên liệu nấu ăn",
			"https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Ô tô & Xe máy",
			"Phụ tùng, phụ kiện và sản phẩm chăm sóc xe",
			"https://images.unsplash.com/photo-1577278689329-1914b6814d58?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Điện gia dụng",
			"Các thiết bị điện tử gia dụng và nhà bếp",
			"https://png.pngtree.com/template/20220330/ourmid/pngtree-electrical-appliances-renewal-season-small-appliances-promotion-poster-image_907595.jpg",
		},
		{
			"Nội thất & Trang trí",
			"Đồ nội thất và trang trí không gian sống",
			"https://images.unsplash.com/photo-1565183928294-7063f23ce0f8?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Đồ chơi & Sở thích",
			"Đồ chơi, game và sản phẩm cho sở thích",
			"https://images.unsplash.com/photo-1566576912321-d58ddd7a6088?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Thiết bị công nghiệp",
			"Máy móc, thiết bị và vật tư công nghiệp",
			"https://images.unsplash.com/photo-1581092918056-0c4c3acd3789?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Nông nghiệp & Vườn tược",
			"Thiết bị, phân bón và sản phẩm nông nghiệp",
			"https://images.unsplash.com/photo-1486328228599-85db4443971f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
	}

	mainCategoryIDs := make(map[string]int64)
	for _, cat := range mainCategories {
		var id int64
		err := db.QueryRow(ctx, `SELECT id FROM categories WHERE name = $1`, cat.name).Scan(&id)

		if err != nil {
			err = db.QueryRow(ctx, `INSERT INTO categories (name, description, image_url, is_active)
                VALUES ($1, $2, $3, TRUE)
                RETURNING id;`, cat.name, cat.desc, cat.imageUrl).Scan(&id)

			if err != nil {
				log.Printf("Error inserting main category: %v", err)
				continue
			}
		}

		mainCategoryIDs[cat.name] = id
	}

	// Seed danh mục con - mở rộng danh sách con
	subCategories := []struct{ name, desc, parent, imageUrl string }{
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
		{
			"Máy tính bảng",
			"Máy tính bảng và phụ kiện",
			"Điện tử & Công nghệ",
			"https://images.unsplash.com/photo-1589739900843-a2120b1d8e92?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1074&q=80",
		},
		{
			"Máy tính để bàn",
			"Máy tính để bàn và linh kiện",
			"Điện tử & Công nghệ",
			"https://images.unsplash.com/photo-1593640408182-31c70c8268f5?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1142&q=80",
		},
		{
			"Tai nghe & Loa",
			"Tai nghe, loa và thiết bị âm thanh",
			"Điện tử & Công nghệ",
			"https://images.unsplash.com/photo-1546435770-a3e426bf472b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1165&q=80",
		},
		{
			"Máy ảnh & Máy quay",
			"Máy ảnh, máy quay và phụ kiện",
			"Điện tử & Công nghệ",
			"https://images.unsplash.com/photo-1516724562728-afc824a36e84?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Thời trang trẻ em",
			"Quần áo và phụ kiện cho trẻ em",
			"Thời trang",
			"https://images.unsplash.com/photo-1519457431-44ccd64a579b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Giày dép",
			"Giày dép các loại cho nam và nữ",
			"Thời trang",
			"https://images.unsplash.com/photo-1600269452121-4f2416e55c28?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1965&q=80",
		},
		{
			"Túi xách",
			"Túi xách, ví và balo thời trang",
			"Thời trang",
			"https://images.unsplash.com/photo-1584917865442-de89df76afd3?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Đồng hồ & Trang sức",
			"Đồng hồ, nhẫn, dây chuyền và trang sức",
			"Thời trang",
			"https://images.unsplash.com/photo-1619946794135-5bc917a27793?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1654&q=80",
		},
		{
			"Đồ dùng nhà bếp",
			"Dụng cụ nấu ăn và đồ dùng nhà bếp",
			"Nhà cửa & Đời sống",
			"https://images.unsplash.com/photo-1556911261-6bd341186b2f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Đồ dùng phòng ngủ",
			"Chăn, ga, gối, đệm và đồ dùng phòng ngủ",
			"Nhà cửa & Đời sống",
			"https://images.unsplash.com/photo-1522771739844-6a9f6d5f14af?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Đồ dùng phòng tắm",
			"Khăn tắm, rèm và đồ dùng phòng tắm",
			"Nhà cửa & Đời sống",
			"https://images.unsplash.com/photo-1584622650111-993a426fbf0a?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Văn phòng phẩm",
			"Bút, giấy và dụng cụ văn phòng",
			"Sách & Văn phòng phẩm",
			"https://images.unsplash.com/photo-1574359411659-11a4b689bc48?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1122&q=80",
		},
		{
			"Sách giáo khoa",
			"Sách giáo khoa và tài liệu học tập",
			"Sách & Văn phòng phẩm",
			"https://images.unsplash.com/photo-1503676260728-1c00da094a0b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1122&q=80",
		},
		{
			"Tạp chí & Báo",
			"Tạp chí, báo và ấn phẩm định kỳ",
			"Sách & Văn phòng phẩm",
			"https://images.unsplash.com/photo-1617137984095-74e4e5e3613f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1074&q=80",
		},
		{
			"Dụng cụ thể thao",
			"Dụng cụ, thiết bị và quần áo thể thao",
			"Thể thao & Du lịch",
			"https://images.unsplash.com/photo-1517649763962-0c623066013b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Vali & Túi du lịch",
			"Vali, túi du lịch và phụ kiện",
			"Thể thao & Du lịch",
			"https://images.unsplash.com/photo-1581553680321-4fffae59fccd?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Thiết bị cắm trại",
			"Lều, túi ngủ và thiết bị cắm trại",
			"Thể thao & Du lịch",
			"https://images.unsplash.com/photo-1504851149312-7a075b496cc7?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Xe đạp & Phụ kiện",
			"Xe đạp, phụ tùng và phụ kiện",
			"Thể thao & Du lịch",
			"https://images.unsplash.com/photo-1541625602330-2277a4c46182?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Đồ dùng cho bé",
			"Bỉm, sữa và đồ dùng cho bé",
			"Mẹ & Bé",
			"https://images.unsplash.com/photo-1515488042361-ee00e0ddd4e4?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1175&q=80",
		},
		{
			"Đồ chơi cho bé",
			"Đồ chơi giáo dục và giải trí cho bé",
			"Mẹ & Bé",
			"https://images.unsplash.com/photo-1566140967404-b8b3932483f5?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Thời trang cho bé",
			"Quần áo, giày dép cho bé",
			"Mẹ & Bé",
			"https://images.unsplash.com/photo-1611042553365-9b101441c135?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Đồ dùng cho mẹ",
			"Sản phẩm dành cho mẹ bầu và sau sinh",
			"Mẹ & Bé",
			"https://images.unsplash.com/photo-1519710164239-da123dc03ef4?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
		},
		{
			"Mỹ phẩm",
			"Mỹ phẩm và trang điểm",
			"Làm đẹp & Sức khỏe",
			"https://images.unsplash.com/photo-1596462502278-27bfdc403348?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=880&q=80",
		},
		{
			"Chăm sóc da",
			"Sản phẩm chăm sóc da mặt và cơ thể",
			"Làm đẹp & Sức khỏe",
			"https://images.unsplash.com/photo-1570172619644-dfd03ed5d881?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Chăm sóc tóc",
			"Sản phẩm chăm sóc và tạo kiểu tóc",
			"Làm đẹp & Sức khỏe",
			"https://images.unsplash.com/photo-1562157873-818bc0726f68?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=880&q=80",
		},
		{
			"Thực phẩm chức năng",
			"Vitamin, thực phẩm bổ sung và thảo dược",
			"Làm đẹp & Sức khỏe",
			"https://images.unsplash.com/photo-1577174881658-0f30ed549adc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Thiết bị y tế",
			"Máy đo đường huyết, huyết áp và thiết bị y tế gia đình",
			"Làm đẹp & Sức khỏe",
			"https://images.unsplash.com/photo-1581595219361-c2a3858daa21?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1169&q=80",
		},
		{
			"Thực phẩm khô",
			"Gạo, mì, ngũ cốc và thực phẩm khô",
			"Thực phẩm & Đồ uống",
			"https://images.unsplash.com/photo-1558961363-fa8fdf82db35?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1065&q=80",
		},
		{
			"Thực phẩm tươi sống",
			"Rau củ, trái cây, thịt và thực phẩm tươi sống",
			"Thực phẩm & Đồ uống",
			"https://images.unsplash.com/photo-1488459716781-31db52582fe9?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Đồ uống",
			"Nước giải khát, bia, rượu và đồ uống",
			"Thực phẩm & Đồ uống",
			"https://images.unsplash.com/photo-1581349485608-9469926a8e5e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=764&q=80",
		},
		{
			"Bánh kẹo & Đồ ăn vặt",
			"Bánh kẹo, snack và đồ ăn vặt",
			"Thực phẩm & Đồ uống",
			"https://images.unsplash.com/photo-1582058091505-f87a2e55a40f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Phụ tùng ô tô",
			"Phụ tùng, linh kiện và phụ kiện ô tô",
			"Ô tô & Xe máy",
			"https://images.unsplash.com/photo-1486262715619-67b85e0b08d3?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1172&q=80",
		},
		{
			"Phụ tùng xe máy",
			"Phụ tùng, linh kiện và phụ kiện xe máy",
			"Ô tô & Xe máy",
			"https://images.unsplash.com/photo-1558981001-792f6c0d5068?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Chăm sóc xe",
			"Sản phẩm chăm sóc, vệ sinh và bảo dưỡng xe",
			"Ô tô & Xe máy",
			"https://images.unsplash.com/photo-1520340356584-f9917d1eea6f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Thiết bị định vị & Điện tử ô tô",
			"Thiết bị GPS, camera hành trình và điện tử ô tô",
			"Ô tô & Xe máy",
			"https://images.unsplash.com/photo-1619538419737-edebb2e4af83?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Tủ lạnh & Tủ đông",
			"Tủ lạnh, tủ đông và tủ mát",
			"Điện gia dụng",
			"https://images.unsplash.com/photo-1588854337221-4cf9fa96059c?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Máy giặt & Máy sấy",
			"Máy giặt, máy sấy và thiết bị giặt ủi",
			"Điện gia dụng",
			"https://images.unsplash.com/photo-1626806787461-102c1a7d1155?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Điều hòa & Quạt",
			"Điều hòa không khí, quạt và thiết bị làm mát",
			"Điện gia dụng",
			"https://images.unsplash.com/photo-1553776590-89774c09baeb?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Thiết bị nhà bếp",
			"Lò vi sóng, lò nướng và thiết bị nhà bếp",
			"Điện gia dụng",
			"https://images.unsplash.com/photo-1630459065645-55f3669a92ed?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=880&q=80",
		},
		{
			"Nội thất phòng khách",
			"Sofa, bàn trà và nội thất phòng khách",
			"Nội thất & Trang trí",
			"https://images.unsplash.com/photo-1583847268964-b28dc8f51f92?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
		},
		{
			"Nội thất phòng ngủ",
			"Giường, tủ quần áo và nội thất phòng ngủ",
			"Nội thất & Trang trí",
			"https://images.unsplash.com/photo-1617325247661-675ab4b64ae2?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
		},
		{
			"Nội thất phòng ăn",
			"Bàn ăn, ghế và nội thất phòng ăn",
			"Nội thất & Trang trí",
			"https://images.unsplash.com/photo-1617806118233-18e1de247200?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1032&q=80",
		},
		{
			"Đèn & Thiết bị chiếu sáng",
			"Đèn trần, đèn bàn và thiết bị chiếu sáng",
			"Nội thất & Trang trí",
			"https://images.unsplash.com/photo-1513506003901-1e6a229e2d15?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},

		// Thêm danh mục con cho "Đồ chơi & Sở thích"
		{
			"Đồ chơi trẻ em",
			"Đồ chơi và trò chơi cho trẻ em",
			"Đồ chơi & Sở thích",
			"https://images.unsplash.com/photo-1558060370-8c436e9e5d76?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1176&q=80",
		},
		{
			"Mô hình & Đồ sưu tầm",
			"Mô hình, đồ sưu tầm và đồ chơi cao cấp",
			"Đồ chơi & Sở thích",
			"https://images.unsplash.com/photo-1516562309708-05f3b2b2c238?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1169&q=80",
		},
		{
			"Nhạc cụ",
			"Đàn guitar, piano và nhạc cụ",
			"Đồ chơi & Sở thích",
			"https://images.unsplash.com/photo-1511192336575-5a79af67a629?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1632&q=80",
		},
		{
			"Đồ thủ công & Mỹ nghệ",
			"Vật liệu thủ công và mỹ nghệ",
			"Đồ chơi & Sở thích",
			"https://images.unsplash.com/photo-1499744349893-0c6de53516e6?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1136&q=80",
		},

		// Thêm danh mục con cho "Thiết bị công nghiệp"
		{
			"Máy móc công nghiệp",
			"Máy móc và thiết bị công nghiệp",
			"Thiết bị công nghiệp",
			"https://images.unsplash.com/photo-1566937169390-7be4c63b8a0e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Dụng cụ điện",
			"Máy khoan, máy cắt và dụng cụ điện",
			"Thiết bị công nghiệp",
			"https://images.unsplash.com/photo-1530124566582-a618bc2615dc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Thiết bị an toàn",
			"Mũ bảo hiểm, găng tay và thiết bị an toàn lao động",
			"Thiết bị công nghiệp",
			"https://images.unsplash.com/photo-1601171903232-8663ec287c2e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Thiết bị đo lường",
			"Thước, máy đo và thiết bị đo lường",
			"Thiết bị công nghiệp",
			"https://images.unsplash.com/photo-1572372783017-2b80336200d5?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},

		// Thêm danh mục con cho "Nông nghiệp & Vườn tược"
		{
			"Vật tư nông nghiệp",
			"Phân bón, hạt giống và vật tư nông nghiệp",
			"Nông nghiệp & Vườn tược",
			"https://images.unsplash.com/photo-1589923188651-268a357A047E?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Máy móc nông nghiệp",
			"Máy cắt cỏ, máy bơm và máy móc nông nghiệp",
			"Nông nghiệp & Vườn tược",
			"https://images.unsplash.com/photo-1575379573116-bd5e9c629046?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		{
			"Cây cảnh & Hoa",
			"Cây cảnh, hạt giống hoa và phụ kiện",
			"Nông nghiệp & Vườn tược",
			"https://images.unsplash.com/photo-1501004318641-b39e6451bec6?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1073&q=80",
		},
		{
			"Dụng cụ làm vườn",
			"Xẻng, kéo cắt cành và dụng cụ làm vườn",
			"Nông nghiệp & Vườn tược",
			"https://images.unsplash.com/photo-1598902468171-0f50e32f3e57?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
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
		err := db.QueryRow(ctx, `SELECT id FROM categories WHERE name = $1`, subCat.name).Scan(&existingID)

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

// Cải tiến seedEnhancedProducts để đảm bảo seedTags được gọi trước khi thêm sản phẩm
func seedEnhancedProducts(ctx context.Context, db *pgxpool.Pool, supplierIDs []int64) {
	// Đảm bảo tags và attributes đã được tạo trước
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

	if !exists {
		log.Printf("The product_variants table does not exist. Make sure tables are created in the correct order.")
		log.Printf("Attempting to seed products first before product variants.")
	}

	// Lấy danh sách categories
	type CategoryInfo struct {
		id       int64
		name     string
		parentID sql.NullInt64
	}
	categories := make(map[int64]CategoryInfo) // id -> CategoryInfo
	categoryNameToID := make(map[string]int64) // name -> id

	rows, err := db.Query(ctx, `SELECT id, name, parent_id FROM categories`)
	if err != nil {
		log.Fatalf("Error getting categories: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category CategoryInfo
		if err := rows.Scan(&category.id, &category.name, &category.parentID); err != nil {
			log.Printf("Error scanning category: %v", err)
			continue
		}
		categories[category.id] = category
		categoryNameToID[category.name] = category.id
	}

	// Tổ chức categories thành parent -> []children
	parentToChildren := make(map[int64][]int64)
	parentCategories := make(map[int64]CategoryInfo)
	childCategories := make(map[int64]CategoryInfo)

	for id, category := range categories {
		if !category.parentID.Valid {
			// Đây là category cha
			parentCategories[id] = category
		} else {
			// Đây là category con
			childCategories[id] = category
			parentID := category.parentID.Int64
			parentToChildren[parentID] = append(parentToChildren[parentID], id)
		}
	}

	// Danh sách ảnh sản phẩm chất lượng cao từ Unsplash theo danh mục
	productImages := map[string][]string{
		"Điện thoại thông minh": {
			"https://images.unsplash.com/photo-1585060544812-6b45742d762f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1281&q=80",
			"https://images.unsplash.com/photo-1598327105666-5b89351aff97?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=2227&q=80",
			"https://images.unsplash.com/photo-1529653762956-b0a27278529c?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1605236453806-6ff36851218e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Máy tính xách tay": {
			"https://images.unsplash.com/photo-1496181133206-80ce9b88a853?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
			"https://images.unsplash.com/photo-1593642702821-c8da6771f0c6?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1332&q=80",
			"https://images.unsplash.com/photo-1603302576837-37561b2e2302?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1168&q=80",
		},
		"Máy tính bảng": {
			"https://images.unsplash.com/photo-1561154464-82e9adf32764?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1587&q=80",
			"https://images.unsplash.com/photo-1544244015-0df4b3ffc6b0?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1473&q=80",
		},
		"Tai nghe & Loa": {
			"https://images.unsplash.com/photo-1546435770-a3e426bf472b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1165&q=80",
			"https://images.unsplash.com/photo-1563330232-57114bb0823c?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Máy ảnh & Máy quay": {
			"https://images.unsplash.com/photo-1516035069371-29a1b244cc32?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1164&q=80",
			"https://images.unsplash.com/photo-1510127034890-ba27508e9f1c?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Tủ lạnh & Tủ đông": {
			"https://images.unsplash.com/photo-1584568694244-14fbdf83bd30?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1586455122341-cb7c5a37590a?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Thời trang nam": {
			"https://images.unsplash.com/photo-1490578474895-699cd4e2cf59?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
			"https://images.unsplash.com/photo-1617137984095-74e4e5e3613f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1074&q=80",
		},
		"Thời trang nữ": {
			"https://images.unsplash.com/photo-1567401893414-76b7b1e5a7a5?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1483985988355-763728e1935b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Đồ gia dụng": {
			"https://images.unsplash.com/photo-1556909172-54557c7e4fb7?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1560185893-a55cbc8c57e8?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Sách": {
			"https://images.unsplash.com/photo-1495446815901-a7297e633e8d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1513475382585-d06e58bcb0e0?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
		},
		"Thể thao": {
			"https://images.unsplash.com/photo-1517649763962-0c623066013b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1547919307-1ecb10702e6f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=688&q=80",
		},
	}

	// Images cho parent categories
	parentCategoryImages := map[string][]string{
		"Điện tử & Công nghệ": {
			"https://images.unsplash.com/photo-1468495244123-6c6c332eeece?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1021&q=80",
			"https://images.unsplash.com/photo-1550745165-9bc0b252726f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Thời trang": {
			"https://images.unsplash.com/photo-1445205170230-053b83016050?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
			"https://images.unsplash.com/photo-1589182337642-35f6e9ccbf8d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
		},
		"Nhà cửa & Đời sống": {
			"https://images.unsplash.com/photo-1484101403633-562f891dc89a?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1174&q=80",
			"https://images.unsplash.com/photo-1493663284031-b7e3aefcae8e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Sách & Văn phòng phẩm": {
			"https://images.unsplash.com/photo-1526243741027-444d633d7365?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
			"https://images.unsplash.com/photo-1512903989781-40c28368bd5d?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Thể thao & Du lịch": {
			"https://images.unsplash.com/photo-1517649763962-0c623066013b?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1530266451970-40ded5a4d66a?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Mẹ & Bé": {
			"https://images.unsplash.com/photo-1518531933037-91b2f5f229cc?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1554684652-57e82094ad77?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Làm đẹp & Sức khỏe": {
			"https://images.unsplash.com/photo-1571875257727-256c39da42af?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1171&q=80",
			"https://images.unsplash.com/photo-1526947425960-945c6e72858f?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		},
		"Thực phẩm & Đồ uống": {
			"https://images.unsplash.com/photo-1495474472287-4d71bcdd2085?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
			"https://images.unsplash.com/photo-1540914124281-342587941389?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1074&q=80",
		},
	}

	// Default images cho các category chưa có ảnh chuyên biệt
	defaultImages := []string{
		"https://images.unsplash.com/photo-1523275335684-37898b6bab30?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1099&q=80",
		"https://images.unsplash.com/photo-1505740420928-5e560c06d30e?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		"https://images.unsplash.com/photo-1542291026-7eec264c27ff?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=1170&q=80",
		"https://images.unsplash.com/photo-1553456558-aff63285bdd1?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8&auto=format&fit=crop&w=687&q=80",
	}

	// Đảm bảo mỗi category đều có ảnh
	for _, category := range categories {
		// Kiểm tra danh mục con
		if _, exists := productImages[category.name]; !exists {
			productImages[category.name] = defaultImages
		}

		// Kiểm tra danh mục cha
		if !category.parentID.Valid {
			if _, exists := parentCategoryImages[category.name]; !exists {
				parentCategoryImages[category.name] = defaultImages
			}
		}
	}

	// Số lượng sản phẩm đã tạo
	totalProducts := 0

	// PHẦN 1: TẠO SẢN PHẨM CHO CATEGORY CHA
	log.Println("Creating products for parent categories...")

	// Tạo thêm tên và mô tả cho sản phẩm của danh mục cha
	parentProductNames := map[string][]string{
		"Điện tử & Công nghệ": {
			"Bộ sản phẩm Điện tử cao cấp",
			"Combo Thiết bị Công nghệ hiện đại",
			"Bộ sưu tập Gadget mới nhất",
			"Tech Premium Bundle",
			"Smart Home Combo",
		},
		"Thời trang": {
			"Bộ sưu tập Thời trang cao cấp",
			"Fashion Collection 2024",
			"Combo Quần áo & Phụ kiện",
			"Seasonal Fashion Bundle",
			"Style Essentials Pack",
		},
		"Nhà cửa & Đời sống": {
			"Bộ sản phẩm Gia đình đa năng",
			"Home Essentials Pack",
			"Combo Nội thất & Gia dụng",
			"Living Space Collection",
			"Home Improvement Bundle",
		},
		"Sách & Văn phòng phẩm": {
			"Bộ sách Bestseller",
			"Combo Văn phòng tiện lợi",
			"Bộ sưu tập Sách & Stationery",
			"Office Essentials Pack",
			"Literature & Craft Bundle",
		},
		"Thể thao & Du lịch": {
			"Bộ dụng cụ Thể thao đa năng",
			"Combo Du lịch tiện ích",
			"Travel & Sport Collection",
			"Fitness Essentials Pack",
			"Adventure Gear Bundle",
		},
	}

	// Thêm mô tả cho sản phẩm danh mục cha
	parentProductDescriptions := map[string][]string{
		"Điện tử & Công nghệ": {
			"Bộ sản phẩm điện tử cao cấp với những tính năng hiện đại nhất, kết nối liền mạch và trải nghiệm người dùng tuyệt vời.",
			"Combo thiết bị công nghệ hiện đại mang lại trải nghiệm số hoá toàn diện cho ngôi nhà thông minh của bạn.",
			"Bộ sưu tập gadget mới nhất với thiết kế tinh tế, chức năng vượt trội và công nghệ tiên tiến.",
			"Tech Premium Bundle sẽ nâng tầm trải nghiệm công nghệ của bạn với các sản phẩm chất lượng cao và tích hợp thông minh.",
			"Smart Home Combo giúp biến ngôi nhà của bạn thành không gian sống thông minh, tiện nghi và tiết kiệm năng lượng.",
		},
		"Thời trang": {
			"Bộ sưu tập thời trang cao cấp từ các thương hiệu hàng đầu, mang đến phong cách độc đáo và đẳng cấp.",
			"Fashion Collection 2024 với những thiết kế mới nhất, theo xu hướng thời trang quốc tế và chất liệu cao cấp.",
			"Combo quần áo & phụ kiện giúp bạn tạo nên phong cách riêng, hài hòa và thời thượng cho mọi dịp.",
			"Seasonal Fashion Bundle mang đến những món đồ thời trang phù hợp với mùa, dễ phối và trendy.",
			"Style Essentials Pack với những món đồ cơ bản không thể thiếu, dễ kết hợp và luôn thời trang.",
		},
	}

	// Tạo mô tả mặc định cho những danh mục chưa có mô tả cụ thể
	defaultParentDescriptions := []string{
		"Bộ sản phẩm cao cấp với thiết kế hiện đại, chất lượng vượt trội và đa dạng công năng sử dụng.",
		"Combo sản phẩm chính hãng với đầy đủ phụ kiện và chế độ bảo hành tốt nhất trên thị trường.",
		"Bộ sưu tập mới nhất với thiết kế độc đáo, công nghệ tiên tiến và trải nghiệm người dùng vượt trội.",
		"Bộ sản phẩm đa năng phù hợp cho mọi nhu cầu sử dụng, tiết kiệm chi phí và không gian.",
		"Combo tiết kiệm với giá cả hợp lý, chất lượng đảm bảo và đa dạng tính năng.",
	}

	for parentID, parentCategory := range parentCategories {
		// Lấy tên và mô tả phù hợp cho danh mục cha này
		var names []string
		var descriptions []string

		if specificNames, ok := parentProductNames[parentCategory.name]; ok {
			names = specificNames
		} else {
			// Tạo tên mặc định nếu không có tên cụ thể
			names = []string{
				fmt.Sprintf("Bộ sản phẩm %s cao cấp", parentCategory.name),
				fmt.Sprintf("Combo %s chính hãng", parentCategory.name),
				fmt.Sprintf("Bộ sưu tập %s mới nhất", parentCategory.name),
				fmt.Sprintf("Bộ %s đa năng", parentCategory.name),
				fmt.Sprintf("Combo %s tiết kiệm", parentCategory.name),
			}
		}

		if specificDescs, ok := parentProductDescriptions[parentCategory.name]; ok {
			descriptions = specificDescs
		} else {
			descriptions = defaultParentDescriptions
		}

		// Lấy ảnh cho danh mục cha
		var parentImages []string
		if images, ok := parentCategoryImages[parentCategory.name]; ok {
			parentImages = images
		} else {
			parentImages = defaultImages
		}

		// Mỗi nhà cung cấp tạo ít nhất 1-2 sản phẩm cho mỗi danh mục cha
		for _, supplierID := range supplierIDs {
			numProducts := gofakeit.Number(1, 2)

			for i := 0; i < numProducts; i++ {
				// Chọn ngẫu nhiên tên sản phẩm và mô tả
				productName := names[gofakeit.Number(0, len(names)-1)]
				productDesc := descriptions[gofakeit.Number(0, len(descriptions)-1)]

				// Chọn ngẫu nhiên ảnh sản phẩm
				productImage := parentImages[gofakeit.Number(0, len(parentImages)-1)]

				// Tạo SKU prefix
				skuPrefix := strings.ToUpper(string([]rune(parentCategory.name)[0])) +
					strings.ToUpper(string([]rune(productName)[0])) +
					fmt.Sprintf("%03d", gofakeit.Number(100, 999))

				// Kiểm tra xem sản phẩm đã tồn tại chưa
				var existingID string
				err := db.QueryRow(ctx, `
					SELECT id FROM products WHERE name = $1 AND supplier_id = $2 AND category_id = $3
				`, productName, supplierID, parentID).Scan(&existingID)

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
						supplierID, parentID, productName, productDesc, productImage,
						"active", gofakeit.Bool(), "standard", skuPrefix, float32(gofakeit.Float32Range(3.5, 5)),
					).Scan(&productID)

					if err != nil {
						log.Printf("Error inserting parent category product: %v", err)
						continue
					}
				} else {
					productID = existingID
					// Cập nhật sản phẩm đã tồn tại nếu cần
					_, err := db.Exec(ctx, `
						UPDATE products SET
						description = $1, image_url = $2, status = 'active',
						tax_class = 'standard', sku_prefix = $3
						WHERE id = $4
					`, productDesc, productImage, skuPrefix, productID)

					if err != nil {
						log.Printf("Error updating parent category product: %v", err)
						continue
					}
				}

				// Thêm tags cho sản phẩm của danh mục cha
				numTags := gofakeit.Number(1, 3)
				tagNames := []string{"Mới nhất", "Bán chạy", "Chính hãng", "Giảm giá", "Chất lượng cao", "Cao cấp", "Bộ sản phẩm", "Combo"}

				for j := 0; j < numTags; j++ {
					randomTag := tagNames[gofakeit.Number(0, len(tagNames)-1)]

					// Lấy tag ID
					var tagID string
					err := db.QueryRow(ctx, `
						SELECT id FROM tags WHERE name = $1
					`, randomTag).Scan(&tagID)

					if err != nil {
						log.Printf("Error getting tag ID: %v", err)
						continue
					}

					// Kiểm tra xem liên kết product-tag đã tồn tại chưa
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

				// Tạo thuộc tính mặc định cho sản phẩm cha
				defaultAttrs := map[string][]string{
					"Màu sắc":   {"Đen", "Trắng", "Bạc", "Xanh", "Đỏ", "Vàng", "Nâu"},
					"Chất liệu": {"Cao cấp", "Nhựa", "Kim loại", "Composite", "Hợp kim", "Vải", "Gỗ"},
					"Xuất xứ":   {"Việt Nam", "Trung Quốc", "Nhật Bản", "Hàn Quốc", "Thái Lan", "Mỹ", "Đức"},
				}

				// Kiểm tra xem bảng product_variants đã tồn tại chưa
				if exists {
					// Tạo biến thể sản phẩm cho sản phẩm cha
					createProductVariants(ctx, db, productID, skuPrefix, defaultAttrs, productImage)
				}

				totalProducts++
			}
		}
	}

	// PHẦN 2: TẠO SẢN PHẨM CHO CATEGORY CON
	log.Println("Creating products for child categories...")

	// Tạo sản phẩm cho mỗi danh mục con
	for catID, category := range childCategories {
		categoryName := category.name

		// Check for product names & descriptions
		productNames, ok := categoryProductNames[categoryName]
		if !ok {
			// Nếu không có tên sản phẩm cụ thể cho danh mục, dùng tên mặc định
			productNames = []string{
				fmt.Sprintf("Sản phẩm %s 1", categoryName),
				fmt.Sprintf("Sản phẩm %s 2", categoryName),
				fmt.Sprintf("Sản phẩm %s cao cấp", categoryName),
				fmt.Sprintf("Sản phẩm %s tiết kiệm", categoryName),
				fmt.Sprintf("Sản phẩm %s đặc biệt", categoryName),
			}
		}

		productDescriptions, ok := categoryProductDescriptions[categoryName]
		if !ok {
			// Nếu không có mô tả cụ thể cho danh mục, dùng mô tả mặc định
			productDescriptions = []string{
				"Sản phẩm chất lượng cao, thiết kế hiện đại và công năng vượt trội.",
				"Sản phẩm tiết kiệm, bền bỉ với giá thành hợp lý cho mọi gia đình.",
				"Sản phẩm cao cấp với thiết kế tinh tế, chất lượng vượt trội và nhiều tính năng đặc biệt.",
				"Sản phẩm đáng tin cậy với chất lượng ổn định và dịch vụ hậu mãi chu đáo.",
			}
		}

		// Check for images
		images, ok := productImages[categoryName]
		if !ok {
			// Đã xử lý ở trên, nhưng kiểm tra lại để đảm bảo
			images = defaultImages
		}

		// Check for attributes
		categoryAttrs, ok := categoryAttributes[categoryName]
		if !ok {
			// Nếu không có thuộc tính cụ thể cho danh mục, dùng thuộc tính mặc định
			categoryAttrs = map[string][]string{
				"Màu sắc":    {"Đen", "Trắng", "Xám", "Xanh", "Đỏ", "Vàng", "Nâu", "Bạc"},
				"Kích thước": {"Nhỏ", "Vừa", "Lớn", "XL", "XXL", "Freesize"},
				"Xuất xứ":    {"Việt Nam", "Trung Quốc", "Nhật Bản", "Hàn Quốc", "Thái Lan", "Mỹ", "Đức"},
				"Chất liệu":  {"Nhựa", "Kim loại", "Vải", "Gỗ", "Da", "Thủy tinh", "Cao su"},
			}
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
                    SELECT id FROM products WHERE name = $1 AND supplier_id = $2 AND category_id = $3
                `, productName, supplierID, catID).Scan(&existingID)

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
						supplierID, catID, productName, productDesc, productImage,
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

// Cải tiến seedAttributeDefinitions để đảm bảo rằng các thuộc tính được tạo đúng
func seedAttributeDefinitions(ctx context.Context, db *pgxpool.Pool) {
	attributes := []struct {
		name, desc, inputType    string
		isFilterable, isRequired bool
	}{
		// Thuộc tính hiện có
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
		{"Kích cỡ màn hình", "Kích thước màn hình hiển thị", "select", true, false},

		// Thêm các thuộc tính mới ở đây
		{"Loại kết nối", "Loại kết nối của thiết bị", "select", true, false},
		{"Kiểu đeo", "Kiểu đeo tai nghe", "select", true, false},
		{"Thời lượng pin", "Thời lượng pin của thiết bị", "select", true, false},
		{"Độ phân giải", "Độ phân giải của camera", "select", true, false},
		{"Cảm biến", "Loại cảm biến của camera", "select", true, false},
		{"Khả năng quay video", "Khả năng quay video của camera", "select", true, false},
		{"Loại da", "Loại da phù hợp với sản phẩm", "select", true, false},
		{"Chứng nhận", "Chứng nhận của sản phẩm", "select", false, false},
		{"Hiệu quả", "Công dụng và hiệu quả của sản phẩm", "select", true, false},
		{"Thành phần chính", "Thành phần chính của sản phẩm", "select", true, false},
		{"Hạn sử dụng", "Thời hạn sử dụng sản phẩm", "select", true, false},
		{"Quy cách đóng gói", "Quy cách đóng gói sản phẩm", "select", true, false},
		{"Phương pháp chế biến", "Phương pháp chế biến sản phẩm", "select", false, false},
		{"Loại đồ uống", "Loại đồ uống", "select", true, false},
		{"Dung tích", "Dung tích của sản phẩm", "select", true, false},
		{"Vị", "Hương vị sản phẩm", "select", true, false},
		{"Đóng gói", "Cách đóng gói sản phẩm", "select", false, false},
		{"Độ cồn", "Độ cồn trong đồ uống", "select", true, false},
		{"Chất liệu khung", "Chất liệu khung của sản phẩm", "select", true, false},
		{"Chất liệu bọc", "Chất liệu bọc của sản phẩm", "select", true, false},
		{"Độ tuổi phù hợp", "Độ tuổi phù hợp với sản phẩm", "select", true, false},
		{"Loại cây", "Loại cây cảnh", "select", true, false},
		{"Điều kiện sống", "Điều kiện sống của cây", "select", true, false},
		{"Chậu cây", "Loại chậu cây", "select", true, false},
		{"Công dụng", "Công dụng của cây cảnh", "select", true, false},
		{"Kích thước giường", "Kích thước giường", "select", true, false},
		{"Độ cứng nệm", "Độ cứng của nệm", "select", true, false},
		{"Card đồ họa", "Loại card đồ họa", "select", true, false},
		{"Kết nối", "Loại kết nối của thiết bị", "select", true, false},
		{"Loại máy", "Loại máy ảnh", "select", true, false},
		{"Loại đồ chơi", "Loại đồ chơi", "select", true, false},
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

// Cải tiến createProductVariants để giải quyết các lỗi
func createProductVariants(
	ctx context.Context,
	db *pgxpool.Pool,
	productID string,
	skuPrefix string,
	categoryAttrs map[string][]string,
	productImage string,
) {
	// Kiểm tra xem có thuộc tính nào đã tồn tại trong database không
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

	// Nếu không có thuộc tính hợp lệ, tạo một mặc định
	if len(validAttrNames) == 0 {
		// Tạo thuộc tính mặc định "Kích thước"
		var attrID int
		err := db.QueryRow(ctx, `
			INSERT INTO attribute_definitions 
			(name, description, input_type, is_filterable, is_required)
			VALUES ('Kích thước', 'Kích thước sản phẩm', 'select', true, true)
			ON CONFLICT (name) DO UPDATE
			SET input_type = 'select', is_filterable = true
			RETURNING id;
		`).Scan(&attrID)

		if err != nil {
			log.Printf("Error creating default attribute: %v", err)
			return
		}

		// Thêm các tùy chọn cho kích thước
		sizeOptions := []string{"S", "M", "L", "XL"}
		for _, option := range sizeOptions {
			_, err := db.Exec(ctx, `
				INSERT INTO attribute_options (attribute_definition_id, option_value)
				VALUES ($1, $2)
				ON CONFLICT (attribute_definition_id, option_value) DO NOTHING;
			`, attrID, option)

			if err != nil {
				log.Printf("Error inserting attribute option: %v", err)
			}
		}

		validAttrNames = append(validAttrNames, "Kích thước")
		categoryAttrs["Kích thước"] = sizeOptions
	}

	// Chọn 2 thuộc tính để tạo biến thể
	var variantAttrs []string
	for _, attrName := range validAttrNames {
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

		// Nếu không có tùy chọn, tạo tùy chọn mặc định
		if len(attributeOptions[attrName]) == 0 {
			for _, optionValue := range categoryAttrs[attrName] {
				var optionID int
				err := db.QueryRow(ctx, `
					INSERT INTO attribute_options (attribute_definition_id, option_value)
					VALUES ($1, $2)
					RETURNING id;
				`, attrID, optionValue).Scan(&optionID)

				if err != nil {
					log.Printf("Error creating attribute option: %v", err)
					continue
				}

				attributeOptions[attrName][optionValue] = optionID
			}
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
	variantCount := 0

	for i, attrValue := range attrValues {
		// Check if option exists, create if not
		optionID, ok := attributeOptions[attrName][attrValue]
		if !ok {
			// Create the option if it doesn't exist
			err := db.QueryRow(ctx, `
				INSERT INTO attribute_options (attribute_definition_id, option_value)
				VALUES ($1, $2)
				RETURNING id;
			`, attributeDefs[attrName], attrValue).Scan(&optionID)

			if err != nil {
				log.Printf("Error creating attribute option: %v", err)
				continue
			}
			attributeOptions[attrName][attrValue] = optionID
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
			// FIX: Đảm bảo discountPrice luôn nhỏ hơn basePrice
			// Áp dụng mức giảm giá từ 5% đến 25%
			discountPercent := gofakeit.Float32Range(0.05, 0.25)
			discountAmount := float32(math.Floor(float64(basePrice*discountPercent)/1000) * 1000)
			// Đảm bảo giảm ít nhất 5000 VND và giảm giá luôn nhỏ hơn giá gốc
			if discountAmount < 5000 {
				discountAmount = 5000
			}
			// Nếu mức giảm > 80% giá gốc, giới hạn ở mức 80%
			if discountAmount > basePrice*0.8 {
				discountAmount = float32(math.Floor(float64(basePrice*0.8)/1000) * 1000)
			}
			discountPrice = basePrice - discountAmount

			// Kiểm tra lại một lần nữa để đảm bảo
			if discountPrice >= basePrice || discountPrice <= 0 {
				discountPrice = basePrice * 0.85 // Giảm giá mặc định 15%
			}
		}

		// FIX: Sửa lỗi UTF-8 trong SKU
		// Thay vì dùng các ký tự Unicode, dùng string cố định
		// Tạo một SKU hoàn toàn không chứa ký tự Unicode
		timestamp := time.Now().UnixNano() % 1000000
		uniqueSKU := fmt.Sprintf("%s-%d-%d", skuPrefix, i+1, timestamp)

		// Đảm bảo SKU độc nhất
		for existingSKUs[uniqueSKU] {
			timestamp = time.Now().UnixNano() % 1000000
			uniqueSKU = fmt.Sprintf("%s-%d-%d", skuPrefix, i+1, timestamp)
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
				inventory_quantity, shipping_class, image_url, alt_text, is_default, is_active
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id;
		`,
			productID, uniqueSKU, variantName, basePrice, discountPriceParam,
			gofakeit.Number(5, 100), "standard", productImage, variantName, i == 0, true,
		).Scan(&variantID)

		if err != nil {
			log.Printf("Error inserting product variant: %v", err)
			continue
		}

		// Mark this attribute option as used
		usedAttributeOptions[optionID] = true
		existingSKUs[uniqueSKU] = true

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

		variantCount++
	}

	log.Printf("Created %d variants for product %s", variantCount, productID)
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
