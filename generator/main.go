// Project creator executable

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dlclark/regexp2"
	"github.com/ldez/go-git-cmd-wrapper/clone"
	"github.com/ldez/go-git-cmd-wrapper/git"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Url         string
	ProjectName string `yaml:"project-name"`
	Branch      string `yaml:"branch"`
}

func main() {
	fmt.Println("Reading config file...")
	b, err := os.ReadFile("config.yml")

	if err != nil {
		log.Fatal(err)
	}

	conf := Config{}
	err = yaml.Unmarshal(b, &conf)

	if err != nil {
		log.Fatal(err)
	}

	err = validate(&conf)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
	fmt.Printf("Cloning from url %s of %s branch\n", conf.Url, conf.Branch)
	_, err = git.Clone(
		clone.Branch(conf.Branch),
		clone.Repository(conf.Url),
		clone.Directory(conf.ProjectName))

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
	err = os.RemoveAll(fmt.Sprintf("%s/.git", conf.ProjectName))

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Applying project name (%s) into a source...\n", conf.ProjectName)
	if err = renameProject(conf.ProjectName, conf); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finished successfully")
}

func validate(conf *Config) error {
	re, err := regexp2.Compile("^[^A-Z][^0-9][a-z]*$", regexp2.None)

	if err != nil {
		return err
	}
	ok, err := re.MatchString(conf.ProjectName)

	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	return errors.New("invalid project name: must contain only lowercase characters without spacing")
}

func renameProject(dir string, config Config) error {

	entries, err := os.ReadDir(dir)

	if err != nil {
		return err
	}
	path, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	for _, de := range entries {
		tmpl := template.New(de.Name())
		p := filepath.Join(path, de.Name())
		fmt.Printf("Traversing %s\n", filepath.Join(path, de.Name()))

		if de.IsDir() {
			return renameProject(filepath.Join(path, de.Name()), config)
		} else {
			tmpl.ParseFiles(p)
			output := &strings.Builder{}
			err := tmpl.Execute(output, config)

			if err != nil {
				fmt.Printf("warn: %v\n", err)
				continue
			}

			err = os.WriteFile(p, []byte(fmt.Sprint(output)), 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
