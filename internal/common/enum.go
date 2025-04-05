package common

const (
	API_GATEWAY_DB   = "api_gateway_db"
	NOTIFICATIONS_DB = "notifications_db"
	ORDERS_DB        = "orders_db"
	PARTNERS_DB      = "partners_db"
)

type AddressType string

const (
	AddressTypeHome   AddressType = "home"
	AddressTypeWork   AddressType = "work"
	AddressTypePickup AddressType = "pickup"
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
)
