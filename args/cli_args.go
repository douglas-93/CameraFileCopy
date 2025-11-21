package args

import (
	"flag"
	"fmt"
)

type CliArgs struct {
	DestDir  string
	OrigDir  string
	Clean    bool
	DaysOld  int
	MaxItens int
}

func ParseArgs() *CliArgs {
	destDir := flag.String("d", "", "Destination directory")
	origDir := flag.String("o", "", "Origin directory")
	cleanDir := flag.Bool("clean", false, "Clean origin directory after copying")
	daysOld := flag.Int("days", 0, "Consider files older than specified days")
	maxItens := flag.Int("max", 4, "Maximum number of items to process")
	flag.Parse()
	return &CliArgs{
		DestDir:  *destDir,
		OrigDir:  *origDir,
		Clean:    *cleanDir,
		DaysOld:  *daysOld,
		MaxItens: *maxItens,
	}
}

func HelpMenu() {
	fmt.Println()
	fmt.Println("Usage: program -d <destination_directory> -o <source_directory> [-clean]")
	fmt.Println("  -d: Destination directory where files will be copied.")
	fmt.Println("  -o: Source directory from which files will be copied.")
	fmt.Println("  -clean: Optional. If specified, the source directory will be cleaned after copying.")
	fmt.Println("  -days: Optional. Consider files older than specified days to clean. Days must be greater than 0.")
	fmt.Println("  -max: Optional. Maximum number of items to process. Default is 4.")
	fmt.Println()
}
