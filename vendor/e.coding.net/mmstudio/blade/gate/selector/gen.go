package selector

import (
	"e.coding.net/mmstudio/blade/kvs"
	"path"
)

//go:generate optiongen --option_with_struct_name=false --v=true
func _SelectOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Filters":  []Filter{},
		"Strategy": Strategy(RoundRobin),
	}
}

//go:generate optiongen --option_with_struct_name=false --v=true
func _OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Registry": kvs.RegistryV2(nil),
		"RegistryPathResolver": func(serviceKey string) string {
			return path.Join("/sandwich/service", serviceKey)
		},
	}
}

//go:generate gotemplate -outfmt gen_%v "e.coding.net/mmstudio/blade/golib/container/templates/cmap" "ConcurrentMapStringService(string,*serviceGroup,cmap.KeyHashStr)"
