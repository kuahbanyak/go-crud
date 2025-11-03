package types
import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
)
type MSSQLUUID struct {
	uuid.UUID
}
func NewMSSQLUUID() MSSQLUUID {
	return MSSQLUUID{UUID: uuid.New()}
}
func ParseMSSQLUUID(s string) (MSSQLUUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return MSSQLUUID{}, err
	}
	return MSSQLUUID{UUID: u}, nil
}
func FromUUID(u uuid.UUID) MSSQLUUID {
	return MSSQLUUID{UUID: u}
}
func (u MSSQLUUID) ToUUID() uuid.UUID {
	return u.UUID
}
func (u MSSQLUUID) Value() (driver.Value, error) {
	if u.UUID == uuid.Nil {
		return nil, nil
	}
	b := u.UUID[:]
	result := make([]byte, 16)
	result[0] = b[3]
	result[1] = b[2]
	result[2] = b[1]
	result[3] = b[0]
	result[4] = b[5]
	result[5] = b[4]
	result[6] = b[7]
	result[7] = b[6]
	copy(result[8:], b[8:])
	return result, nil
}
func (u *MSSQLUUID) Scan(value interface{}) error {
	if value == nil {
		u.UUID = uuid.Nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		if len(v) != 16 {
			return fmt.Errorf("invalid UUID length: %d", len(v))
		}
		result := make([]byte, 16)
		result[0] = v[3]
		result[1] = v[2]
		result[2] = v[1]
		result[3] = v[0]
		result[4] = v[5]
		result[5] = v[4]
		result[6] = v[7]
		result[7] = v[6]
		copy(result[8:], v[8:])
		var err error
		u.UUID, err = uuid.FromBytes(result)
		return err
	case string:
		var err error
		u.UUID, err = uuid.Parse(v)
		return err
	default:
		return fmt.Errorf("cannot scan %T into MSSQLUUID", value)
	}
}
func (u MSSQLUUID) String() string {
	return u.UUID.String()
}
func (u MSSQLUUID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + u.UUID.String() + `"`), nil
}
func (u *MSSQLUUID) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid UUID JSON format")
	}
	str := string(data[1 : len(data)-1])
	var err error
	u.UUID, err = uuid.Parse(str)
	return err
}

