package common

import (
	"fmt"
	"slices"
	"strings"
)

const (
	API_GATEWAY_DB   = "api_gateway_db"
	NOTIFICATIONS_DB = "notifications_db"
	ORDERS_DB        = "orders_db"
	PARTNERS_DB      = "partners_db"
)

type AddressType string

const (
	AddressTypeHome      AddressType = "Home"
	AddressTypeOffice    AddressType = "Office"
	AddressTypeWarehouse AddressType = "Warehouse"
	AddressTypStorefront AddressType = "Storefront"
	AddressTypeOther     AddressType = "Other"
)

type RoleName string

const (
	RoleAdmin     RoleName = "admin"
	RoleCustomer  RoleName = "customer"
	RoleDeliverer RoleName = "deliverer"
	RoleSupplier  RoleName = "supplier"
)

type PermissionName string

const (
	Create  PermissionName = "create"
	Update  PermissionName = "update"
	Delete  PermissionName = "delete"
	Read    PermissionName = "read"
	Approve PermissionName = "approve"
	Reject  PermissionName = "reject"
)

type ModuleName string

const (
	UserManagement        ModuleName = "User Management"
	RolePermission        ModuleName = "Role & Permission"
	ProductManagement     ModuleName = "Product Management"
	Cart                  ModuleName = "Cart"
	OrderManagement       ModuleName = "Order Management"
	Payment               ModuleName = "Payment"
	ShippingManagement    ModuleName = "Shipping Management"
	ReviewRating          ModuleName = "Review & Rating"
	StoreManagement       ModuleName = "Store Management"
	Onboarding            ModuleName = "Onboarding"
	AddressTypeManagement ModuleName = "Address Type Management"
	ModuleManagement      ModuleName = "Module Management"
	CouponManagement      ModuleName = "Coupon Management"
)

type MethodType string

const (
	Momo MethodType = "momo"
	Cod  MethodType = "cod"
)

type StatusOrder string

const (
	PendingPayment StatusOrder = "pending_payment"  // Chờ thanh toán
	Pending        StatusOrder = "pending"          // Chờ supplier xác nhận
	Confirmed      StatusOrder = "confirmed"        // Supplier đã xác nhận
	Processing     StatusOrder = "processing"       // Đang chuẩn bị hàng
	ReadyToShip    StatusOrder = "ready_to_ship"    // Sẵn sàng giao hàng
	InTransit      StatusOrder = "in_transit"       // Đang vận chuyển (đang ship)
	OutForDelivery StatusOrder = "out_for_delivery" // Sắp giao (shipper đang trên đường)
	Delivered      StatusOrder = "delivered"        // Đã giao thành công
	Cancelled      StatusOrder = "cancelled"
	PaymentFailed  StatusOrder = "payment_failed"
	Refunded       StatusOrder = "refunded"
)

type Enum interface {
	IsValid() bool
	ErrorMessage() string
}

func (s StatusOrder) IsValid() bool {
	validArray := []StatusOrder{PendingPayment, Pending, Confirmed, Processing, ReadyToShip, InTransit, OutForDelivery, Delivered, Cancelled,
		PaymentFailed, Refunded}

	if slices.Contains(validArray, s) {
		return true
	}

	return false
}

func (s StatusOrder) ErrorMessage() string {
	validArray := []string{
		string(PendingPayment), string(Pending), string(Confirmed),
		string(Processing), string(ReadyToShip), string(InTransit),
		string(OutForDelivery), string(Delivered), string(Cancelled),
		string(PaymentFailed), string(Refunded),
	}

	return fmt.Sprintf("Status must be in the one of: [%v]", strings.Join(validArray, ", "))
}
