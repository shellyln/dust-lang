package executor

import (
	"encoding/json"
	"errors"

	mnem "github.com/shellyln/dust-lang/scripting/executor/opcode"
	tsys "github.com/shellyln/dust-lang/scripting/typesys"
	. "github.com/shellyln/takenoco/base"
)

//
type ExecutionContext struct {
	RootScope       *ExecutionScope
	ModuleScope     *ExecutionScope
	LastModuleScope *ExecutionScope
	Scope           *ExecutionScope

	// DynPropCallback func(ctx *ExecutionContext, name string) (*VariableInfo, bool) // TODO:
	// execCount uint // TODO:
	// genSymSeq uint // TODO:

	typeInfoMap   map[string]tsys.TypeInfo
	typeInfoIdMap map[uint32]tsys.TypeInfo
	typeIdMax     uint32
}

//
func NewExecutionContext(rootScope *ExecutionScope) *ExecutionContext {
	typeInfoMap := make(map[string]tsys.TypeInfo)
	typeInfoIdMap := make(map[uint32]tsys.TypeInfo)

	for key, value := range tsys.TypeInfoInitialMap {
		typeInfoMap[key] = value
		typeInfoIdMap[value.Id] = value
	}

	xctx := &ExecutionContext{
		RootScope:       rootScope,
		ModuleScope:     nil,
		LastModuleScope: nil,
		Scope:           rootScope,
		typeInfoMap:     typeInfoMap,
		typeInfoIdMap:   typeInfoIdMap,
		typeIdMax:       tsys.TypeIdInitialMax,
	}

	return xctx
}

// TODO: NewThreadContext(*ExecutionContext)

//
func (p *ExecutionContext) ResetScope() {
	p.Scope = p.RootScope
	p.ModuleScope = nil
}

//
func (p *ExecutionContext) GetVariableInfo(name string) (*VariableInfo, bool) {
	scope := p.Scope
	// if p.DynPropCallback != nil {
	// 	vi, ok := p.DynPropCallback(p, name)
	// 	if ok {
	// 		return vi, ok
	// 	}
	// }
	for {
		vi, ok := scope.Dict[name]
		if ok {
			return vi, true
		} else if scope.Next != nil {
			scope = scope.Next
		} else {
			return nil, false
		}
	}
}

// TODO: GenSym()

//
func (p *ExecutionContext) DefineVariable(name string, vi VariableInfo) (*VariableInfo, bool) {
	z, ok := p.Scope.Dict[name]
	if ok {
		// TODO: shadowing
		return z, false
	}

	ret := &vi
	p.Scope.Dict[name] = &vi

	return ret, true
}

//
func (p *ExecutionContext) CaptureVariable(vi *VariableInfo, newName string) (*VariableInfo, bool) {
	_, ok := p.Scope.Dict[newName]
	if ok {
		return nil, false
	}

	ret := vi
	p.Scope.Dict[newName] = vi

	return ret, true
}

//
func (p *ExecutionContext) PushScope(returnType mnem.AstOpCodeType) {
	prev := p.Scope
	p.Scope = &ExecutionScope{
		Type:      ExecutionScope_Normal,
		Dict:      map[string]*VariableInfo{},
		Next:      prev,
		StackNext: prev,
	}
	if p.ModuleScope == nil {
		p.ModuleScope = p.Scope
		p.LastModuleScope = p.Scope
	}
}

//
func (p *ExecutionContext) PushFuncScope(returnType mnem.AstOpCodeType) {
	var next *ExecutionScope
	if p.ModuleScope == nil {
		next = p.RootScope
	} else {
		next = p.ModuleScope
	}
	prev := p.Scope
	p.Scope = &ExecutionScope{
		Type:      ExecutionScope_Function,
		Dict:      map[string]*VariableInfo{},
		Next:      next,
		StackNext: prev,
	}
}

//
func (p *ExecutionContext) PushLambdaScope(returnType mnem.AstOpCodeType) {
	var next *ExecutionScope
	if p.ModuleScope == nil {
		next = p.RootScope
	} else {
		next = p.ModuleScope
	}
	prev := p.Scope
	p.Scope = &ExecutionScope{
		Type:      ExecutionScope_Lambda,
		Dict:      map[string]*VariableInfo{},
		Next:      next,
		StackNext: prev,
	}
}

//
func (p *ExecutionContext) PopScope() bool {
	if p.Scope.StackNext == nil {
		return false
	} else {
		if p.Scope == p.ModuleScope {
			p.ModuleScope = nil
		}
		p.Scope = p.Scope.StackNext
		return true
	}
}

//
func (p *ExecutionContext) NewTypeId() uint32 {
	p.typeIdMax++
	return p.typeIdMax
}

//
func (p *ExecutionContext) DummyTypeId() uint32 {
	return 0
}

//
func (p *ExecutionContext) GetTypeInfoByName(name string) (tsys.TypeInfo, error) {

	// TODO: find the primitive or custom type by symbol name ([0])
	//       (from scope namespace)

	return tsys.GetTypeInfo(name)
}

// TODO:
func (p *ExecutionContext) GetTypeInfoById(id uint32) (tsys.TypeInfo, error) {
	typeInfo, ok := p.typeInfoIdMap[id]
	if !ok {
		return tsys.TypeInfo_Primitive_Any, errors.New("error!!!") // TODO:
	}

	return typeInfo, nil // TODO:
}

// TODO:
func (p *ExecutionContext) GetRegisteredTypeInfo(ty tsys.TypeInfo) (tsys.TypeInfo, error) {
	bytes, err := json.Marshal(ty)
	if err != nil {
		return tsys.TypeInfo_Primitive_Any, err
	}
	key := string(bytes)

	z, ok := p.typeInfoMap[key]
	if !ok {
		// TODO: register type and get id
		z.Id = p.NewTypeId()
		z.Flags &^= mnem.TypeIdMask
		z.Flags |= AstOpCodeType(z.Id << uint32(mnem.TypeIdOffset))
		p.typeInfoMap[key] = z
	}

	return z, nil
}
