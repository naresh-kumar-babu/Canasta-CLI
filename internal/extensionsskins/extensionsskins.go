package extensionsskins

import (
	"fmt"
	"strings"

	"github.com/CanastaWiki/Canasta-CLI-Go/internal/logging"
	"github.com/CanastaWiki/Canasta-CLI-Go/internal/orchestrators"
)

type Item struct {
	Name                     string
	RelativeInstallationPath string
	PhpCommand               string
}

func Contains(list []string, element string) bool {
	for _, item := range list {
		if item == element {
			return true
		}
	}
	return false
}

func List(instance logging.Installation, constants Item) {
	fmt.Printf("Available %s:\n", constants.Name)
	fmt.Print(orchestrators.Exec(instance.Path, instance.Orchestrator, "web", "cd $MW_HOME/"+constants.RelativeInstallationPath+" && find * -maxdepth 0 -type d"))
}

func CheckInstalled(name string, instance logging.Installation, constants Item) (string, error) {
	output := orchestrators.Exec(instance.Path, instance.Orchestrator, "web", "cd $MW_HOME/"+constants.RelativeInstallationPath+" && find * -maxdepth 0 -type d")
	if !Contains(strings.Split(output, "\n"), name) {
		return "", fmt.Errorf("%s %s doesn't exist", name, constants.Name)
	}
	return name, nil
}

func Enable(name string, instance logging.Installation, constants Item) {
	phpScript := fmt.Sprintf("<?php\n// This file was generated by Canasta\n%s( '%s' );", constants.PhpCommand, name)
	filePath := fmt.Sprintf("/mediawiki/config/settings/%s.php", name)
	output, err := orchestrators.ExecWithError(instance.Path, instance.Orchestrator, "web", "ls "+filePath)
	if err == nil {
		logging.Fatal(fmt.Errorf("Extension is already enabled! Skipping overwrite."))
	} else if Contains(strings.Split(output, ":"), " No such file or directory\n") {
		command := fmt.Sprintf(`echo -e "%s" > %s`, phpScript, filePath)
		orchestrators.Exec(instance.Path, instance.Orchestrator, "web", command)
		fmt.Printf("Extension %s enabled\n", name)
	}
}

func CheckEnabled(name string, instance logging.Installation, constants Item) (string, error) {
	output := orchestrators.Exec(instance.Path, instance.Orchestrator, "web", "ls /mediawiki/config/settings/")
	if !Contains(strings.Split(output, "\n"), name+".php") {
		return "", fmt.Errorf("%s %s is not enabled", name, constants.Name)
	}
	output = orchestrators.Exec(instance.Path, instance.Orchestrator, "web", fmt.Sprintf(`cat /mediawiki/config/settings/%s.php`, name))
	if !Contains(strings.Split(output, "\n"), "// This file was generated by Canasta") {
		return "", fmt.Errorf("%s %s was not generated by Canasta cli", name, constants.Name)
	}
	return name, nil
}

func Disable(name string, instance logging.Installation, constants Item) {
	command := fmt.Sprintf(`rm /mediawiki/config/settings/%s.php`, name)
	orchestrators.Exec(instance.Path, instance.Orchestrator, "web", command)
	fmt.Printf("%s %s disabled\n", constants.Name, name)
}
