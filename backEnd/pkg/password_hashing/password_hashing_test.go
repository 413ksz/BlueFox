package passwordHashing_test

import (
	"log"
	"sync"
	"testing"

	"github.com/413ksz/BlueFox/backEnd/pkg/models"
	passwordHashing "github.com/413ksz/BlueFox/backEnd/pkg/password_hashing"
)

// Global instance of the Argon2ID hasher to be used across tests
var argon2ID *passwordHashing.KriptoArgon2ID

func TestMain(m *testing.M) {
	// A mock pepper secret for testing
	pepperSecret := []byte("a-very-long-and-secure-pepper-secret-for-testing")

	// Setup the Argon2ID instance with recommended parameters
	var err *models.CustomError
	argon2ID, err = passwordHashing.NewKriptoArgon2ID(32, 1, 64*1024, 4, 32, pepperSecret)
	if err != nil {
		log.Fatalf("Failed to create KriptoArgon2ID instance for tests: %v", err)
	}
	m.Run()
}

// TestNewKriptoArgon2Id tests the constructor and its parameter validation
func TestNewKriptoArgon2Id(t *testing.T) {
	pepperSecret := []byte("a-very-long-and-secure-pepper-secret-for-testing")
	tests := []struct {
		name         string
		saltLength   uint8
		memoryCostKb uint32
		keyLengthB   uint32
		pepperLen    int
		wantErr      bool
		wantErrCode  models.ErrorCode
	}{
		{
			name:         "Valid Parameters",
			saltLength:   32,
			memoryCostKb: 64 * 1024,
			keyLengthB:   32,
			pepperLen:    len(pepperSecret),
			wantErr:      false,
			wantErrCode:  "",
		},
		{
			name:         "Invalid Salt Length",
			saltLength:   15,
			memoryCostKb: 64 * 1024,
			keyLengthB:   32,
			pepperLen:    len(pepperSecret),
			wantErr:      true,
			wantErrCode:  models.ERROR_CODE_INTERNAL_SERVER,
		},
		{
			name:         "Invalid Memory Cost",
			saltLength:   32,
			memoryCostKb: 1024,
			keyLengthB:   32,
			pepperLen:    len(pepperSecret),
			wantErr:      true,
			wantErrCode:  models.ERROR_CODE_INTERNAL_SERVER,
		},
		{
			name:         "Invalid Key Length",
			saltLength:   32,
			memoryCostKb: 64 * 1024,
			keyLengthB:   31,
			pepperLen:    len(pepperSecret),
			wantErr:      true,
			wantErrCode:  models.ERROR_CODE_INTERNAL_SERVER,
		},
		{
			name:         "Invalid Pepper Secret Length",
			saltLength:   32,
			memoryCostKb: 64 * 1024,
			keyLengthB:   32,
			pepperLen:    15,
			wantErr:      true,
			wantErrCode:  models.ERROR_CODE_INTERNAL_SERVER,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret := make([]byte, tt.pepperLen)
			_, err := passwordHashing.NewKriptoArgon2ID(tt.saltLength, 1, tt.memoryCostKb, 4, tt.keyLengthB, secret)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewKriptoArgon2Id() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if err.Code != tt.wantErrCode {
					t.Errorf("NewKriptoArgon2Id() returned unexpected error code: got %s, want %s", err.Code, tt.wantErrCode)
				}
			}
		})
	}
}

// TestGenerateNew tests the hashing function
func TestGenerateNew(t *testing.T) {
	if argon2ID == nil {
		t.Fatal("Argon2ID instance is nil, check TestMain setup")
	}

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid Password",
			password: "mysecurepassword",
			wantErr:  false,
		},
		{
			name:     "Empty Password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "Long Password",
			password: "thisisareallylongpasswordthatshouldstillbehashedcorrectlybytheargon2idalgo",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := argon2ID.GenerateNew(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateNew() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == "" {
				t.Errorf("GenerateNew() returned an empty hash for password: %s", tt.password)
			}
		})
	}
}

// TestVerifyPassword tests the password verification function
func TestVerifyPassword(t *testing.T) {
	if argon2ID == nil {
		t.Fatal("Argon2ID instance is nil, check TestMain setup")
	}

	validPassword := "testpassword123"
	// Generate a known valid hash to test against
	validHash, err := argon2ID.GenerateNew(validPassword)
	if err != nil {
		t.Fatalf("Failed to generate test hash: %v", err)
	}

	tests := []struct {
		name         string
		password     string
		hash         string
		wantVerified bool
		wantErrCode  models.ErrorCode
	}{
		{
			name:         "Correct Password",
			password:     validPassword,
			hash:         validHash,
			wantVerified: true,
			wantErrCode:  "",
		},
		{
			name:         "Incorrect Password",
			password:     "wrongpassword",
			hash:         validHash,
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNAUTHORIZED,
		},
		{
			name:         "Empty Password (Correct)",
			password:     "",
			hash:         func() string { h, _ := argon2ID.GenerateNew(""); return h }(),
			wantVerified: true,
			wantErrCode:  "",
		},
		{
			name:         "Empty Password (Incorrect)",
			password:     "",
			hash:         validHash,
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNAUTHORIZED,
		},
		{
			name:         "Empty Hash",
			password:     validPassword,
			hash:         "",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Malformed Hash",
			password:     validPassword,
			hash:         "notavalidargon2idhash",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := argon2ID.Verify(tt.password, tt.hash)
			gotVerified := err == nil
			if gotVerified != tt.wantVerified {
				t.Errorf("VerifyPassword() gotVerified = %v, want %v for password '%s'", gotVerified, tt.wantVerified, tt.password)
			}

			if !tt.wantVerified {
				// If an error was expected, check its code
				if err != nil { // Check to ensure err is not nil before accessing .Code
					if err.Code != tt.wantErrCode {
						// Log the actual and expected error codes
						t.Logf("Actual error message: %s", err.Details)
						t.Errorf("VerifyPassword() returned unexpected error code: got %s, want %s", err.Code, tt.wantErrCode)
					}
				} else {
					// This case should not be reached if wantVerified is false
					t.Errorf("Expected an error but got nil")
				}
			}
		})
	}
}

// TestGenerateNew_Concurrency checks if the hashing function behaves correctly under concurrent access.
func TestGenerateNew_Concurrency(t *testing.T) {
	if argon2ID == nil {
		t.Fatal("Argon2ID instance is nil, check TestMain setup")
	}

	var wg sync.WaitGroup
	numGoroutines := 10
	password := "concurrentpassword"
	hashes := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)

	// Use a WaitGroup to ensure all goroutines finish before the main thread exits
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hash, err := argon2ID.GenerateNew(password)
			if err != nil {
				errors <- err
				return
			}
			hashes <- hash
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(hashes)
	close(errors)

	// Collect and verify results from both channels
	var testErrors []error
	for err := range errors {
		testErrors = append(testErrors, err)
	}

	// If any goroutine failed to generate a hash, the test should fail
	if len(testErrors) > 0 {
		for _, err := range testErrors {
			t.Errorf("Concurrency test failed with error: %v", err)
		}
		return
	}

	// Now, verify each successfully generated hash
	for hash := range hashes {
		if err := argon2ID.Verify(password, hash); err != nil {
			t.Errorf("Concurrency test: Failed to verify a concurrently generated hash: %v", err)
		}
	}
}
