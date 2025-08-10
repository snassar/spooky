package age

import (
	"bytes"
	"io"
	"strings"
	"time"

	"filippo.io/age"
	"filippo.io/age/armor"

	spookylogging "spooky/internal/logging"
	spookytypeslogging "spooky/internal/types/logging"
	spookytypessecrets "spooky/internal/types/secrets"
)

// AgeClient interface defines the operations for age encryption/decryption
type AgeClient interface {
	// Encrypt encrypts data with the given recipients
	Encrypt(data []byte, recipients []string) (*spookytypessecrets.EncryptedValue, error)

	// Decrypt decrypts data with the given identity
	Decrypt(encrypted *spookytypessecrets.EncryptedValue, identity string) ([]byte, error)

	// GenerateKey generates a new age key pair
	GenerateKey() (identity string, recipient string, err error)

	// ParseRecipient validates a recipient string
	ParseRecipient(recipient string) error

	// ParseIdentity validates an identity string
	ParseIdentity(identity string) error

	// EncryptStream encrypts data from a reader to a writer with streaming support
	EncryptStream(input io.Reader, output io.Writer, recipients []string) error

	// DecryptStream decrypts data from a reader to a writer with streaming support
	DecryptStream(input io.Reader, output io.Writer, identity string) error
}

// Client implements the AgeClient interface using the filippo.io/age library
type Client struct {
	logger spookytypeslogging.Logger
}

// NewClient creates a new age client
func NewClient() *Client {
	return &Client{
		logger: spookylogging.GetLogger(),
	}
}

// Encrypt encrypts data with the given recipients
func (c *Client) Encrypt(data []byte, recipients []string) (*spookytypessecrets.EncryptedValue, error) {
	c.logger.Debug("Encrypting data", spookylogging.Int("data_size", len(data)), spookylogging.Int("recipients", len(recipients)))

	// Parse recipients
	ageRecipients := make([]age.Recipient, 0, len(recipients))
	for _, recipientStr := range recipients {
		recipient, err := age.ParseX25519Recipient(recipientStr)
		if err != nil {
			return nil, &spookytypessecrets.SecretsError{
				Operation: "parse_recipient",
				Cause:     err,
				Context: map[string]interface{}{
					"recipient": recipientStr,
				},
			}
		}
		ageRecipients = append(ageRecipients, recipient)
	}

	// Create encrypted value
	encrypted := &spookytypessecrets.EncryptedValue{
		Recipients: recipients,
		Created:    time.Now(),
		Metadata:   make(map[string]string),
	}

	// Encrypt the data
	var buf bytes.Buffer
	armorWriter := armor.NewWriter(&buf)
	encryptWriter, err := age.Encrypt(armorWriter, ageRecipients...)
	if err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "create_encrypt_writer",
			Cause:     err,
		}
	}

	if _, err := encryptWriter.Write(data); err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "write_encrypted_data",
			Cause:     err,
		}
	}

	if err := encryptWriter.Close(); err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "close_encrypt_writer",
			Cause:     err,
		}
	}

	if err := armorWriter.Close(); err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "close_armor_writer",
			Cause:     err,
		}
	}

	encrypted.Data = buf.Bytes()
	c.logger.Debug("Data encrypted successfully", spookylogging.Int("encrypted_size", len(encrypted.Data)))

	return encrypted, nil
}

// Decrypt decrypts data with the given identity
func (c *Client) Decrypt(encrypted *spookytypessecrets.EncryptedValue, identityStr string) ([]byte, error) {
	c.logger.Debug("Decrypting data", spookylogging.Int("data_size", len(encrypted.Data)))

	// Parse identity
	identity, err := age.ParseX25519Identity(identityStr)
	if err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "parse_identity",
			Cause:     err,
			Context: map[string]interface{}{
				"identity": identityStr,
			},
		}
	}

	// Create armor reader
	armorReader := armor.NewReader(bytes.NewReader(encrypted.Data))

	// Create decrypt reader
	decryptReader, err := age.Decrypt(armorReader, identity)
	if err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "create_decrypt_reader",
			Cause:     err,
		}
	}

	// Read decrypted data
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, decryptReader); err != nil {
		return nil, &spookytypessecrets.SecretsError{
			Operation: "read_decrypted_data",
			Cause:     err,
		}
	}

	decrypted := buf.Bytes()
	c.logger.Debug("Data decrypted successfully", spookylogging.Int("decrypted_size", len(decrypted)))

	return decrypted, nil
}

// GenerateKey generates a new age key pair
func (c *Client) GenerateKey() (identity string, recipient string, err error) {
	c.logger.Debug("Generating new age key pair")

	// Generate X25519 identity
	ageIdentity, err := age.GenerateX25519Identity()
	if err != nil {
		return "", "", &spookytypessecrets.SecretsError{
			Operation: "generate_identity",
			Cause:     err,
		}
	}

	// Get the recipient from the identity
	ageRecipient := ageIdentity.Recipient()

	// Format the keys
	identity = strings.TrimSpace(ageIdentity.String())
	recipient = strings.TrimSpace(ageRecipient.String())

	c.logger.Debug("Age key pair generated successfully")

	return identity, recipient, nil
}

// ParseRecipient parses a recipient string
func (c *Client) ParseRecipient(recipient string) error {
	_, err := age.ParseX25519Recipient(recipient)
	if err != nil {
		return &spookytypessecrets.SecretsError{
			Operation: "parse_recipient",
			Cause:     err,
			Context: map[string]interface{}{
				"recipient": recipient,
			},
		}
	}
	return nil
}

// ParseIdentity parses an identity string
func (c *Client) ParseIdentity(identity string) error {
	_, err := age.ParseX25519Identity(identity)
	if err != nil {
		return &spookytypessecrets.SecretsError{
			Operation: "parse_identity",
			Cause:     err,
			Context: map[string]interface{}{
				"identity": identity,
			},
		}
	}
	return nil
}

// EncryptStream encrypts data from a reader to a writer with streaming support
func (c *Client) EncryptStream(input io.Reader, output io.Writer, recipients []string) error {
	c.logger.Debug("Encrypting stream", spookylogging.Int("recipients", len(recipients)))

	// Parse recipients
	ageRecipients := make([]age.Recipient, 0, len(recipients))
	for _, recipientStr := range recipients {
		recipient, err := age.ParseX25519Recipient(recipientStr)
		if err != nil {
			return &spookytypessecrets.SecretsError{
				Operation: "parse_recipient",
				Cause:     err,
				Context: map[string]interface{}{
					"recipient": recipientStr,
				},
			}
		}
		ageRecipients = append(ageRecipients, recipient)
	}

	// Create armor writer
	armorWriter := armor.NewWriter(output)
	defer armorWriter.Close()

	// Create encrypt writer
	encryptWriter, err := age.Encrypt(armorWriter, ageRecipients...)
	if err != nil {
		return &spookytypessecrets.SecretsError{
			Operation: "create_encrypt_writer",
			Cause:     err,
		}
	}
	defer encryptWriter.Close()

	// Copy data from input to encrypt writer
	if _, err := io.Copy(encryptWriter, input); err != nil {
		return &spookytypessecrets.SecretsError{
			Operation: "copy_to_encrypt_writer",
			Cause:     err,
		}
	}

	c.logger.Debug("Stream encrypted successfully")
	return nil
}

// DecryptStream decrypts data from a reader to a writer with streaming support
func (c *Client) DecryptStream(input io.Reader, output io.Writer, identityStr string) error {
	c.logger.Debug("Decrypting stream")

	// Parse identity
	identity, err := age.ParseX25519Identity(identityStr)
	if err != nil {
		return &spookytypessecrets.SecretsError{
			Operation: "parse_identity",
			Cause:     err,
			Context: map[string]interface{}{
				"identity": identityStr,
			},
		}
	}

	// Create armor reader
	armorReader := armor.NewReader(input)

	// Create decrypt reader
	decryptReader, err := age.Decrypt(armorReader, identity)
	if err != nil {
		return &spookytypessecrets.SecretsError{
			Operation: "create_decrypt_reader",
			Cause:     err,
		}
	}

	// Copy data from decrypt reader to output
	if _, err := io.Copy(output, decryptReader); err != nil {
		return &spookytypessecrets.SecretsError{
			Operation: "copy_from_decrypt_reader",
			Cause:     err,
		}
	}

	c.logger.Debug("Stream decrypted successfully")
	return nil
}
