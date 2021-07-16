package selector

import (
	"e.coding.net/mmstudio/blade/kvs"
)

const (
	VersionKey = "service_version"
)

// FilterEndpoint is an endpoint based Select Filter which will
// only return services with the endpoint specified.
func FilterEndpoint(id string) Filter {
	return func(old kvs.Entries) kvs.Entries {
		if id == "" {
			return old
		}

		var services kvs.Entries

		for _, service := range old {
			if service.Identifier == id {
				services = append(services, service)
				break
			}
		}

		return services
	}
}

// FilterLabel is a label based Select Filter which will
// only return services with the label specified.
func FilterLabel(key, val string) Filter {
	return func(old kvs.Entries) kvs.Entries {
		var services kvs.Entries

		for _, service := range old {
			if service.Metadata == nil {
				continue
			}

			if v, ok := service.Metadata[key]; ok && v == val {
				services = append(services, service)
			}
		}

		return services
	}
}

// FilterVersion is a version based Select Filter which will
// only return services with the version specified.
func FilterVersion(version string) Filter {
	return func(old kvs.Entries) kvs.Entries {
		var services kvs.Entries

		for _, service := range old {
			if service.Metadata == nil {
				continue
			}

			if v, ok := service.Metadata[VersionKey]; ok && v == version {
				services = append(services, service)
			}
		}

		return services
	}
}
