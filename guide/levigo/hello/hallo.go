package main

const (
	deutschPrefix    = "Hallo, "
	castellanoPrefix = "Hola, "
	frencaiPrefix    = "Bonjour, "
	espanish         = "Castellano"
	french           = "Français"
)

func Hallo(name string, language string) string {
	if name == "" {
		name = "Welt"
	}

	return greetingPrefix(language) + name
}

func greetingPrefix(language string) (prefix string) {
	switch language {
	case french:
		prefix = frencaiPrefix
	case espanish:
		prefix = castellanoPrefix
	default:
		prefix = deutschPrefix
	}
	return
}

func main() {
	println(Hallo("Levi", "Français"))
}
