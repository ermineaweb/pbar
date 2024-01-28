package pbar

type color string

const (
	// public colors for custom conf
	BLACK_BRIGHT   color = "\x1B[1;90m"
	RED_BRIGHT     color = "\x1B[1;91m"
	GREEN_BRIGHT   color = "\x1B[1;92m"
	YELLOW_BRIGHT  color = "\x1B[1;93m"
	BLUE_BRIGHT    color = "\x1B[1;94m"
	MAGENTA_BRIGHT color = "\x1B[1;95m"
	CYAN_BRIGHT    color = "\x1B[1;96m"
	WHITE_BRIGHT   color = "\x1B[1;97m"
	BLACK          color = "\x1B[1;30m"
	RED            color = "\x1B[1;31m"
	GREEN          color = "\x1B[1;32m"
	YELLOW         color = "\x1B[1;33m"
	BLUE           color = "\x1B[1;34m"
	MAGENTA        color = "\x1B[1;35m"
	CYAN           color = "\x1B[1;36m"
	WHITE          color = "\x1B[1;37m"

	default_color color = "\x1B[1;0m"

	delete_line = "\033[K"
)
