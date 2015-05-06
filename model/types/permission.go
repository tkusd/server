package types

import "encoding/json"

type Permission int

const (
	PermissionRead Permission = 1 << iota
	PermissionWrite
	PermissionAdmin
)

type permissionJSON struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
	Admin bool `json:"admin"`
}

func (p Permission) MarshalJSON() ([]byte, error) {
	return json.Marshal(&permissionJSON{
		Read:  p.Readable(),
		Write: p.Writable(),
		Admin: p.Adminable(),
	})
}

func (p *Permission) UnmarshalJSON(val []byte) error {
	var data permissionJSON

	if err := json.Unmarshal(val, &data); err != nil {
		return err
	}

	if data.Read {
		*p |= PermissionRead
	}

	if data.Write {
		*p |= PermissionWrite
	}

	if data.Admin {
		*p |= PermissionAdmin
	}

	return nil
}

func (p Permission) Readable() bool {
	return p&PermissionRead != 0
}

func (p Permission) Writable() bool {
	return p&PermissionWrite != 0
}

func (p Permission) Adminable() bool {
	return p&PermissionAdmin != 0
}
