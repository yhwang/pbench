package gen_config

import (
	"embed"
	"encoding/json"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"path/filepath"
	"presto-benchmark/log"
)

const configJson = "config.json"

var (
	TemplatePath  = ""
	TemplateFS    fs.FS
	ParameterPath = ""
	//go:embed template
	builtinTemplate embed.FS
)

func Run(_ *cobra.Command, args []string) {
	gParams := DefaultGeneratorParameters
	if ParameterPath != "" {
		if paramsByte, ioErr := os.ReadFile(ParameterPath); ioErr != nil {
			log.Error().Err(ioErr).Str("parameter_path", ParameterPath).
				Msg("failed to read generator parameter file")
			ParameterPath = ""
		} else {
			params := &GeneratorParameters{}
			if unmarshalErr := json.Unmarshal(paramsByte, params); unmarshalErr != nil {
				log.Error().Err(unmarshalErr).Str("parameter_path", ParameterPath).
					Msg("failed to unmarshal generator parameter file")
			} else {
				gParams = params
			}
		}
	}
	if TemplatePath == "" {
		TemplateFS = builtinTemplate
		TemplatePath = "template"
	} else {
		TemplateFS = os.DirFS(TemplatePath)
		TemplatePath = "."
	}
	configs := make([]*ClusterConfig, 0, 3)
	// We swallow all the errors (file/directory does not exist or parsing error) to continue as much as we can.
	_ = filepath.Walk(args[0], func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() || info.Name() != configJson {
			return nil
		}
		bytes, ioErr := os.ReadFile(path)
		if ioErr != nil {
			log.Error().Err(ioErr).Str("path", path).Msg("failed to read config file.")
			return nil
		}
		cfg := &ClusterConfig{
			GeneratorParameters: gParams,
		}
		if ioErr = json.Unmarshal(bytes, cfg); ioErr != nil {
			log.Error().Err(ioErr).Str("path", path).Msg("failed to parse config file.")
		} else {
			cfg.Path = filepath.Dir(path)
			// Calculate the variables based on the spec in the config.json
			log.Info().Str("path", path).Msg("parsed configuration")
			cfg.Calculate()
			configs = append(configs, cfg)
		}
		return nil
	})
	// Generate config files for each config we read.
	GenerateFiles(configs)
}