package koanfx

// SmartMerge keeps both sides when scalar ↔ map conflict happens.
func SmartMerge(src, dest map[string]any) error {
	for k, v := range src {
		if dv, ok := dest[k]; ok {
			switch dmap := dv.(type) {
			case map[string]any:
				// dest is map …
				if smap, ok := v.(map[string]any); ok {
					// … src also map  → recurse
					if err := SmartMerge(smap, dmap); err != nil {
						return err
					}
				} else {
					// … src scalar     → tuck scalar inside the map
					dmap[k] = v
				}
			default:
				// dest scalar …
				if smap, ok := v.(map[string]any); ok {
					// … src map        → move scalar into the incoming map
					smap[k] = dv
					dest[k] = smap
				} else {
					// … src scalar     → last one wins (koanf default)
					dest[k] = v
				}
			}
		} else {
			dest[k] = v
		}
	}
	return nil
}
