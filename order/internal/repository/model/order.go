package order

import (
	"github.com/go-faster/errors"
)

// Order – строка таблицы orders в PostgreSQL.
type Order struct {
	OrderUUID       string   `db:"order_uuid"` // первичный ключ
	UserUUID        string   `db:"user_uuid"`
	PartUuids       []string `db:"part_uuids"` // text[] или jsonb в БД
	TotalPrice      float64  `db:"total_price"`
	TransactionUUID *string  `db:"transaction_uuid"` // NULL, если не оплачен
	PaymentMethod   *string  `db:"payment_method"`   // NULL, если не оплачен
	Status          string   `db:"status"`           // PENDING_PAYMENT, PAID, CANCELLED
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
