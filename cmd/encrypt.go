package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	var password string

	encryptCmd := &cobra.Command{
		Use:   "encrypt [profile]",
		Short: "Encrypt a profile file with a password",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if password == "" {
				fatalf("password is required (--password)")
			}
			if err := encryptProfile(args[0], password); err != nil {
				fatalf("%v", err)
			}
		},
	}

	decryptCmd := &cobra.Command{
		Use:   "decrypt [profile]",
		Short: "Decrypt an encrypted profile file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if password == "" {
				fatalf("password is required (--password)")
			}
			if err := decryptProfile(args[0], password); err != nil {
				fatalf("%v", err)
			}
		},
	}

	encryptCmd.Flags().StringVarP(&password, "password", "p", "", "Encryption password")
	decryptCmd.Flags().StringVarP(&password, "password", "p", "", "Decryption password")

	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
}

func deriveKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func encryptProfile(profile, password string) error {
	src := profilePath(profile)
	plaintext, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("profile %q not found", profile)
	}

	block, err := aes.NewCipher(deriveKey(password))
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	dest := filepath.Join(projectDir(), "profiles", profile+".enc")
	if err := os.WriteFile(dest, []byte(hex.EncodeToString(ciphertext)), 0600); err != nil {
		return err
	}
	fmt.Printf("Profile %q encrypted to %s\n", profile, dest)
	return nil
}

func decryptProfile(profile, password string) error {
	src := filepath.Join(projectDir(), "profiles", profile+".enc")
	hexData, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("encrypted profile %q not found", profile)
	}
	ciphertext, err := hex.DecodeString(string(hexData))
	if err != nil {
		return fmt.Errorf("invalid encrypted file: %v", err)
	}

	block, err := aes.NewCipher(deriveKey(password))
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return fmt.Errorf("decryption failed: wrong password?")
	}

	dest := profilePath(profile)
	if err := os.WriteFile(dest, plaintext, 0600); err != nil {
		return err
	}
	fmt.Printf("Profile %q decrypted to %s\n", profile, dest)
	return nil
}
