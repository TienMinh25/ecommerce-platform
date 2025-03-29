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
