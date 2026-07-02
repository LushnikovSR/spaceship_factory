package order

import "github.com/go-faster/errors"

type Order struct {
	// UUID of order.
	OrderUUID string
	// User`s UUID.
	UserUUID string
	// Array of part uuids.
	PartUuids []string
	// Total coast of all parts in the order.
	TotalPrice float64
	// Transaction`s UUID (if paid).
	TransactionUUID OptNilString
	// PaymentMethod (if paid).
	PaymentMethod *NilOrderDtoPaymentMethod
	// Status of order, default "PENDING_PAYMENT".
	Status OrderDtoStatus
}

// GetOrderUUID returns the value of OrderUUID.
func (o *Order) GetOrderUUID() string {
	return o.OrderUUID
}

// GetUserUUID returns the value of UserUUID.
func (o *Order) GetUserUUID() string {
	return o.UserUUID
}

// GetPartUuids returns the value of PartUuids.
func (o *Order) GetPartUuids() []string {
	return o.PartUuids
}

// GetTotalPrice returns the value of TotalPrice.
func (o *Order) GetTotalPrice() float64 {
	return o.TotalPrice
}

// GetTransactionUUID returns the value of TransactionUUID.
func (o *Order) GetTransactionUUID() OptNilString {
	return o.TransactionUUID
}

// GetPaymentMethod returns the value of PaymentMethod.
func (o *Order) GetPaymentMethod() *NilOrderDtoPaymentMethod {
	return o.PaymentMethod
}

// GetStatus returns the value of Status.
func (o *Order) GetStatus() OrderDtoStatus {
	return o.Status
}

// SetOrderUUID sets the value of OrderUUID.
func (o *Order) SetOrderUUID(val string) {
	o.OrderUUID = val
}

// SetUserUUID sets the value of UserUUID.
func (o *Order) SetUserUUID(val string) {
	o.UserUUID = val
}

// SetPartUuids sets the value of PartUuids.
func (o *Order) SetPartUuids(val []string) {
	o.PartUuids = val
}

// SetTotalPrice sets the value of TotalPrice.
func (o *Order) SetTotalPrice(val float64) {
	o.TotalPrice = val
}

// SetTransactionUUID sets the value of TransactionUUID.
func (o *Order) SetTransactionUUID(val OptNilString) {
	o.TransactionUUID = val
}

// SetPaymentMethod sets the value of PaymentMethod.
func (o *Order) SetPaymentMethod(val *NilOrderDtoPaymentMethod) {
	o.PaymentMethod = val
}

// SetStatus sets the value of Status.
func (o *Order) SetStatus(val OrderDtoStatus) {
	o.Status = val
}

// OptNilString is optional nullable string.
type OptNilString struct {
	Value string
	Set   bool
	Null  bool
}

// IsSet returns true if OptNilString was set.
func (o OptNilString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptNilString) Reset() {
	var v string
	o.Value = v
	o.Set = false
	o.Null = false
}

// SetTo sets value to v.
func (o *OptNilString) SetTo(v string) {
	o.Set = true
	o.Null = false
	o.Value = v
}

// IsNull returns true if value is Null.
func (o OptNilString) IsNull() bool { return o.Null }

// SetToNull sets value to null.
func (o *OptNilString) SetToNull() {
	o.Set = true
	o.Null = true
	var v string
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptNilString) Get() (v string, ok bool) {
	if o.Null {
		return v, false
	}
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// NewOptNilString returns new OptNilString with value set to v.
func NewOptNilString(v string) OptNilString {
	return OptNilString{
		Value: v,
		Set:   true,
	}
}

// NewNilOrderDtoPaymentMethod returns new NilOrderDtoPaymentMethod with value set to v.
func NewNilOrderDtoPaymentMethod(v OrderDtoPaymentMethod) NilOrderDtoPaymentMethod {
	return NilOrderDtoPaymentMethod{
		Value: v,
	}
}

// NilOrderDtoPaymentMethod is nullable OrderDtoPaymentMethod.
type NilOrderDtoPaymentMethod struct {
	Value OrderDtoPaymentMethod
	Null  bool
}

// SetTo sets value to v.
func (o *NilOrderDtoPaymentMethod) SetTo(v OrderDtoPaymentMethod) {
	o.Null = false
	o.Value = v
}

// IsNull returns true if value is Null.
func (o NilOrderDtoPaymentMethod) IsNull() bool { return o.Null }

// SetToNull sets value to null.
func (o *NilOrderDtoPaymentMethod) SetToNull() {
	o.Null = true
	var v OrderDtoPaymentMethod
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o NilOrderDtoPaymentMethod) Get() (v OrderDtoPaymentMethod, ok bool) {
	if o.Null {
		return v, false
	}
	return o.Value, true
}

type OrderDtoPaymentMethod string

const (
	OrderDtoPaymentMethodUNKNOWN       OrderDtoPaymentMethod = "UNKNOWN"
	OrderDtoPaymentMethodCARD          OrderDtoPaymentMethod = "CARD"
	OrderDtoPaymentMethodSBP           OrderDtoPaymentMethod = "SBP"
	OrderDtoPaymentMethodCREDITCARD    OrderDtoPaymentMethod = "CREDIT_CARD"
	OrderDtoPaymentMethodINVESTORMONEY OrderDtoPaymentMethod = "INVESTOR_MONEY"
)

// AllValues returns all OrderDtoPaymentMethod values.
func (OrderDtoPaymentMethod) AllValues() []OrderDtoPaymentMethod {
	return []OrderDtoPaymentMethod{
		OrderDtoPaymentMethodUNKNOWN,
		OrderDtoPaymentMethodCARD,
		OrderDtoPaymentMethodSBP,
		OrderDtoPaymentMethodCREDITCARD,
		OrderDtoPaymentMethodINVESTORMONEY,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s OrderDtoPaymentMethod) MarshalText() ([]byte, error) {
	switch s {
	case OrderDtoPaymentMethodUNKNOWN:
		return []byte(s), nil
	case OrderDtoPaymentMethodCARD:
		return []byte(s), nil
	case OrderDtoPaymentMethodSBP:
		return []byte(s), nil
	case OrderDtoPaymentMethodCREDITCARD:
		return []byte(s), nil
	case OrderDtoPaymentMethodINVESTORMONEY:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *OrderDtoPaymentMethod) UnmarshalText(data []byte) error {
	switch OrderDtoPaymentMethod(data) {
	case OrderDtoPaymentMethodUNKNOWN:
		*s = OrderDtoPaymentMethodUNKNOWN
		return nil
	case OrderDtoPaymentMethodCARD:
		*s = OrderDtoPaymentMethodCARD
		return nil
	case OrderDtoPaymentMethodSBP:
		*s = OrderDtoPaymentMethodSBP
		return nil
	case OrderDtoPaymentMethodCREDITCARD:
		*s = OrderDtoPaymentMethodCREDITCARD
		return nil
	case OrderDtoPaymentMethodINVESTORMONEY:
		*s = OrderDtoPaymentMethodINVESTORMONEY
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Merged schema.
type OrderDtoStatus string

const (
	OrderDtoStatusPENDINGPAYMENT OrderDtoStatus = "PENDING_PAYMENT"
	OrderDtoStatusPAID           OrderDtoStatus = "PAID"
	OrderDtoStatusCANCELLED      OrderDtoStatus = "CANCELLED"
)

// AllValues returns all OrderDtoStatus values.
func (OrderDtoStatus) AllValues() []OrderDtoStatus {
	return []OrderDtoStatus{
		OrderDtoStatusPENDINGPAYMENT,
		OrderDtoStatusPAID,
		OrderDtoStatusCANCELLED,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s OrderDtoStatus) MarshalText() ([]byte, error) {
	switch s {
	case OrderDtoStatusPENDINGPAYMENT:
		return []byte(s), nil
	case OrderDtoStatusPAID:
		return []byte(s), nil
	case OrderDtoStatusCANCELLED:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *OrderDtoStatus) UnmarshalText(data []byte) error {
	switch OrderDtoStatus(data) {
	case OrderDtoStatusPENDINGPAYMENT:
		*s = OrderDtoStatusPENDINGPAYMENT
		return nil
	case OrderDtoStatusPAID:
		*s = OrderDtoStatusPAID
		return nil
	case OrderDtoStatusCANCELLED:
		*s = OrderDtoStatusCANCELLED
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Payment method.
// Ref: #/components/schemas/payment_method
type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

// AllValues returns all PaymentMethod values.
func (PaymentMethod) AllValues() []PaymentMethod {
	return []PaymentMethod{
		PaymentMethodUNKNOWN,
		PaymentMethodCARD,
		PaymentMethodSBP,
		PaymentMethodCREDITCARD,
		PaymentMethodINVESTORMONEY,
	}
}

// MarshalText implements encoding.TextMarshaler.
func (s PaymentMethod) MarshalText() ([]byte, error) {
	switch s {
	case PaymentMethodUNKNOWN:
		return []byte(s), nil
	case PaymentMethodCARD:
		return []byte(s), nil
	case PaymentMethodSBP:
		return []byte(s), nil
	case PaymentMethodCREDITCARD:
		return []byte(s), nil
	case PaymentMethodINVESTORMONEY:
		return []byte(s), nil
	default:
		return nil, errors.Errorf("invalid value: %q", s)
	}
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (s *PaymentMethod) UnmarshalText(data []byte) error {
	switch PaymentMethod(data) {
	case PaymentMethodUNKNOWN:
		*s = PaymentMethodUNKNOWN
		return nil
	case PaymentMethodCARD:
		*s = PaymentMethodCARD
		return nil
	case PaymentMethodSBP:
		*s = PaymentMethodSBP
		return nil
	case PaymentMethodCREDITCARD:
		*s = PaymentMethodCREDITCARD
		return nil
	case PaymentMethodINVESTORMONEY:
		*s = PaymentMethodINVESTORMONEY
		return nil
	default:
		return errors.Errorf("invalid value: %q", data)
	}
}

// Ref: #/components/schemas/pay_order_response
type PayOrderResponse struct {
	TransactionUUID string `json:"transaction_uuid"`
}
