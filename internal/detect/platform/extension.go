package platform

type Extension string

func (e Extension) Default(os OS) Extension {
	switch os {
	case Windows:
		return Extension(".exe")
	default:
		return Extension("")
	}
}

func (e Extension) String() string {
	return string(e)
}
