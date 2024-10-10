package main

func splitTags(tags []string) ([]string, []string) {
	var with, without []string
	for _, tag := range tags {
		if tag[0] == '!' {
			without = append(without, tag[1:])
		} else {
			with = append(with, tag)
		}
	}
	return with, without
}
