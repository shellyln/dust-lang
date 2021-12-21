package typesys

//
type TypeManagerContext interface {
	NewTypeId() uint32
	GetTypeInfoByName(name string) (TypeInfo, error)
	GetTypeInfoById(id uint32) (TypeInfo, error)
	GetRegisteredTypeInfo(ty TypeInfo) (TypeInfo, error)
}
