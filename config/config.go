package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type Config struct {
	APIBaseURL string `json:"api_base_url,omitempty"`
	Version    string `json:"-"`
}

var (
	GitTag string = "dev"
	CFG    *Config
)

func castLuaInt(ls *lua.LState, vGlobal string) int {
	v, ok := ls.GetGlobal(vGlobal).(lua.LNumber)
	if !ok {
		log.Fatalf("Erro ao converter %q para int", vGlobal)
	}
	return int(v)
}

func castLuaString(ls *lua.LState, vGlobal string) string {
	v, ok := ls.GetGlobal(vGlobal).(lua.LString)
	if !ok {
		log.Fatalf("Erro ao converter %q para string", vGlobal)
	}
	return string(v)
}

func runLua(luaScript string) error {
	execName, err := os.Executable()
	if err != nil {
		return err
	}
	execName = filepath.Base(execName)
	appName := execName

	ls := lua.NewState()
	ls.PreloadModule("math", lua.OpenMath)

	ls.SetGlobal("AppName", lua.LString(appName))

	ls.SetGlobal("APIBaseURL", lua.LString(CFG.APIBaseURL))
	ls.SetGlobal("Version", lua.LString(CFG.Version))

	err = ls.DoString(luaScript)
	if err != nil {
		return err
	}

	CFG.APIBaseURL = castLuaString(ls, "APIBaseURL")
	// CFG.Version = castLuaString(ls, "Version")

	return nil
}

func CreateKey() string {
	const (
		length  = 32
		charset = "-_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	)
	lenCharset := byte(len(charset))
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = charset[b[i]%lenCharset]
	}
	return string(b)
}

func Encript(key, value string) string {
	if len(key) != 32 {
		log.Fatalf("Key deve ter 32 caracteres")
	}
	if len(value) == 0 {
		log.Fatalf("Value não pode ser vazio")
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatalf("Erro ao criar cipher: %s", err)
	}

	b := make([]byte, 2)
	_, err = rand.Read(b)
	if err != nil {
		log.Fatalf("Erro ao criar byte aleatório: %s", err)
	}

	value = string(b) + value

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("Erro ao criar GCM: %s", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("Erro ao criar nonce: %s", err)
	}

	value = string(gcm.Seal(nonce, nonce, []byte(value), nil))
	value = base64.StdEncoding.EncodeToString([]byte(value))

	return value
}

func Decrypt(key, value string) string {
	if len(key) != 32 {
		log.Fatalf("Key deve ter 32 caracteres")
	}
	if len(value) == 0 {
		log.Fatalf("Value não pode ser vazio")
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatalf("Erro ao criar cipher: %s", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatalf("Erro ao criar GCM: %s", err)
	}

	valueDecoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		log.Fatalf("Erro ao decodificar value: %s", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := valueDecoded[:nonceSize], valueDecoded[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatalf("Erro ao decodificar value: %s", err)
	}

	plaintext = plaintext[2:]

	return string(plaintext)
}

func getEnvInt(key string, defaultValue int) int {
	value, ok := os.LookupEnv(key)
	if ok {
		v, err := strconv.Atoi(value)
		if err != nil {
			log.Fatalf("Erro ao converter %s para int: %s", key, err)
		}
		return v
	}
	return defaultValue
}

func getEnvString(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if ok {
		v, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("Erro ao converter %s para bool: %s", key, err)
		}
		return v
	}
	return defaultValue
}

func processDefaultInt(value int, environmentVar string, defaultValue int) int {
	if value == 0 {
		return getEnvInt(environmentVar, defaultValue)
	}
	return value
}

func processDefaultString(value string, environmentVar string, defaultValue string) string {
	if value == "" {
		return getEnvString(environmentVar, defaultValue)
	}
	return value
}

func processDefaultBool(value bool, environmentVar string, defaultValue bool) bool {
	if value == false {
		return getEnvBool(environmentVar, defaultValue)
	}
	return value
}

func setDefaultStr(value string, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func setDefaultInt(value int, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}

func Load() error {
	CFG = &Config{}
	CFG.Version = GitTag
	if CFG.Version == "" {
		CFG.Version = "dev"
	}

	configFile := "config.lua"
	log.Printf("Using config file %q", configFile)

	b, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Config file %q not found", configFile)
			return err
		}

		log.Panicf("Error reading config file: %v", err)
	}

	err = runLua(string(b))
	if err != nil {
		log.Printf("Error running lua script: %v", err)
		return err
	}

	return err
}

func ListEnvVariables() {
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "GASGO_") {
			fmt.Println(e)
		}
	}
}

func ShowConfig() string {
	b, err := json.MarshalIndent(CFG, "", "  ")
	if err != nil {
		log.Printf("Error marshalling config: %v", err)
		return ""
	}
	return string(b)
}
