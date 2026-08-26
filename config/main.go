package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Main function to load the configuration from YAML and env files, generate them if they don't exist, and validate secrets.
func Load(yamlFile, envFile string) (*Config, error) {
	cfg := &Config{}

	if !fileExists(yamlFile) || !fileExists(envFile) {
		if err := generateConfig(cfg, yamlFile, envFile); err != nil {
			return nil, err
		}
	}

	if err := godotenv.Load(envFile); err != nil {
		return nil, err
	}

	if err := loadYaml(yamlFile, cfg); err != nil {
		return nil, err
	}

	if err := loadEnvOverrides(cfg); err != nil {
		return nil, err
	}

	if err := validateSecrets(cfg, ""); err != nil {
		return nil, err
	}

	return cfg, nil
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// generateConfig generates the YAML and env configuration files with default values.
func generateConfig(cfg *Config, yamlFile, envFile string) error {
	yamlData := map[string]any{}
	envData := map[string]string{}

	if err := populateDefaults(reflect.ValueOf(cfg).Elem(), yamlData, envData, ""); err != nil {
		return err
	}

	yml, _ := yaml.Marshal(yamlData)
	os.MkdirAll(filepath.Dir(yamlFile), 0755)
	os.WriteFile(yamlFile, yml, 0644)

	var lines []string
	for k, v := range envData {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}

	os.WriteFile(envFile, []byte(strings.Join(lines, "\n")), 0600)
	return nil
}

// populateDefaults populates the default values for the configuration struct based on the struct tags.
func populateDefaults(v reflect.Value, yamlData map[string]any, envData map[string]string, prefix string) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)

		yamlTag := f.Tag.Get("yaml")
		envTag := f.Tag.Get("env")
		def := f.Tag.Get("default")
		secret := f.Tag.Get("secret") == "true"

		if fv.Kind() == reflect.Struct {
			child := map[string]any{}
			p := strings.ToUpper(yamlTag)
			if prefix != "" {
				p = prefix + "_" + p
			}
			populateDefaults(fv, child, envData, p)
			yamlData[yamlTag] = child
			continue
		}

		envKey := envTag
		if envKey == "" {
			envKey = strings.ToUpper(prefix + "_" + yamlTag)
		}

		if secret {
			envData[envKey] = "CHANGE_ME"
			continue
		}

		if def == "" {
			continue
		}

		switch fv.Kind() {
		case reflect.Bool:
			yamlData[yamlTag], _ = strconv.ParseBool(def)
		case reflect.Int, reflect.Int64:
			if fv.Type() == reflect.TypeOf(time.Duration(0)) {
				yamlData[yamlTag], _ = time.ParseDuration(def)
			} else {
				yamlData[yamlTag], _ = strconv.Atoi(def)
			}
		case reflect.Slice:
			yamlData[yamlTag] = strings.Split(def, ",")
		default:
			yamlData[yamlTag] = def
		}
	}
	return nil
}

// loadEnvOverrides overrides the configuration values with environment variables if they are set.
func loadEnvOverrides(cfg any) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)

		if fv.Kind() == reflect.Struct {
			loadEnvOverrides(fv.Addr().Interface())
			continue
		}

		key := f.Tag.Get("env")
		if key == "" {
			continue
		}

		val, exists := os.LookupEnv(key)
		if !exists {
			continue
		}

		if f.Tag.Get("secret") == "true" && val == "" {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(val)
		case reflect.Int:
			i, _ := strconv.Atoi(val)
			fv.SetInt(int64(i))
		case reflect.Bool:
			b, _ := strconv.ParseBool(val)
			fv.SetBool(b)
		}
	}
	return nil
}

// validateSecrets checks if all secret fields in the configuration struct are set and not empty.
func validateSecrets(cfg any, path string) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)

		childPath := f.Name
		if path != "" {
			childPath = path + "." + f.Name
		}

		if fv.Kind() == reflect.Struct {
			if err := validateSecrets(fv.Addr().Interface(), childPath); err != nil {
				return err
			}
			continue
		}

		if f.Tag.Get("secret") == "true" {
			if fv.Kind() != reflect.String || fv.String() == "" {
				return fmt.Errorf("missing secret: %s", childPath)
			}
		}
	}
	return nil
}

// loadYaml reads the YAML configuration file and unmarshals its content into the provided configuration struct.
func loadYaml(file string, cfg *Config) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("Failed to read yaml file: %w", err)
	}
	return yaml.Unmarshal(data, cfg)
}

// MaskSecrets replaces the values of secret fields in the configuration struct with "******" to prevent sensitive information from being exposed.
func MaskSecrets(cfg any) {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)

		if fv.Kind() == reflect.Struct {
			MaskSecrets(fv.Addr().Interface())
			continue
		}

		if f.Tag.Get("secret") == "true" && fv.Kind() == reflect.String {
			fv.SetString("******")
		}
	}
}
