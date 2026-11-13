package quota

import "golang.org/x/exp/collection/set"

type WhiteList struct {
	allowType set.Set[string]
	allowList set.Set[string]
}

func NewWhiteList(types []string, codes []string) WhiteList {
	return WhiteList{
		allowType: set.Of(types...),
		allowList: set.Of(codes...),
	}
}

func (x WhiteList) Allow(md Meter) bool {
	return x.allowT(md.Type) && x.allowC(md.Code)
}

func (x WhiteList) AllowLite(code, typ string) bool {
	return x.allowT(typ) && x.allowC(code)
}

func (x WhiteList) allowT(typ string) bool {
	if x.allowType.Len() == 0 {
		return true
	}
	if x.allowType.Contains("*") {
		return true
	}
	return x.allowType.Contains(typ)
}

func (x WhiteList) allowC(code string) bool {
	if x.allowList.Len() == 0 {
		return true
	}

	if x.allowList.Contains("*") {
		return true
	}

	return x.allowList.Contains(code)
}
