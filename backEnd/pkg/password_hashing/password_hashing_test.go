package passwordHashing_test

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"
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
	argon2ID, err = passwordHashing.NewKriptoArgon2ID(32, 1, 64, 4, 32, pepperSecret)
	if err != nil {
		log.Fatalf("Failed to create KriptoArgon2ID instance for tests: %v", err)
	}
	m.Run()
}

// TestNewKriptoArgon2Id tests the constructor and its parameter validation
func TestNewKriptoArgon2Id(t *testing.T) {
	pepperSecret := []byte("a-very-long-and-secure-pepper-secret-for-testing")
	tests := []struct {
		name                string
		saltLength          uint8
		memoryCostMegaBytes uint32
		keyLengthB          uint32
		pepperLen           int
		wantErr             bool
		wantErrCode         models.ErrorCode
	}{
		{
			name:                "Valid Parameters",
			saltLength:          32,
			memoryCostMegaBytes: 64,
			keyLengthB:          32,
			pepperLen:           len(pepperSecret),
			wantErr:             false,
			wantErrCode:         "",
		},
		{
			name:                "Invalid Salt Length",
			saltLength:          15,
			memoryCostMegaBytes: 64,
			keyLengthB:          32,
			pepperLen:           len(pepperSecret),
			wantErr:             true,
			wantErrCode:         models.ERROR_CODE_INITIALIZE_ERROR,
		},
		{
			name:                "Invalid Memory Cost",
			saltLength:          32,
			memoryCostMegaBytes: 10,
			keyLengthB:          32,
			pepperLen:           len(pepperSecret),
			wantErr:             true,
			wantErrCode:         models.ERROR_CODE_INITIALIZE_ERROR,
		},
		{
			name:                "Invalid Key Length",
			saltLength:          32,
			memoryCostMegaBytes: 64,
			keyLengthB:          31,
			pepperLen:           len(pepperSecret),
			wantErr:             true,
			wantErrCode:         models.ERROR_CODE_INITIALIZE_ERROR,
		},
		{
			name:                "Invalid Pepper Secret Length",
			saltLength:          32,
			memoryCostMegaBytes: 64,
			keyLengthB:          32,
			pepperLen:           15,
			wantErr:             true,
			wantErrCode:         models.ERROR_CODE_INITIALIZE_ERROR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret := make([]byte, tt.pepperLen)
			_, err := passwordHashing.NewKriptoArgon2ID(tt.saltLength, 1, tt.memoryCostMegaBytes, 4, tt.keyLengthB, secret)

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
		// 1. Success and Expected Failures
		{
			name:         "Correct Password",
			password:     validPassword,
			hash:         validHash,
			wantVerified: true,
			wantErrCode:  "",
		},
		{
			name:         "Empty Password (Correct)",
			password:     "",
			hash:         func() string { h, _ := argon2ID.GenerateNew(""); return h }(),
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
			name:         "Empty Password (Incorrect)",
			password:     "",
			hash:         validHash,
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNAUTHORIZED,
		},

		// 2. Malformed Hash Structure
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
		{
			name:         "Missing Salt/Hash Separator",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=65536,t=1,p=4$c29tZXNhbHQ",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Missing Cost Factor Separator",
			password:     validPassword,
			hash:         "$argon2id$v=19m=65536,t=1,p=4$c29tZXNhbHQ$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},

		// 3. Malformed Hash Parameters
		{
			name:         "Invalid Memory Cost Prefix",
			password:     validPassword,
			hash:         "$argon2id$v=19$x=65536,t=1,p=4$c29tZXNhbHQ$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Invalid Time Cost Prefix",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=65536,x=1,p=4$c29tZXNhbHQ$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Invalid Threads Prefix",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=65536,t=1,x=4$c29tZXNhbHQ$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Version Not an Integer",
			password:     validPassword,
			hash:         "$argon2id$v=string$m=64,t=1,p=4$c29tZXNhbHQ$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Memory Cost Not an Integer",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=string,t=1,p=4$c29tZXNhbHQ$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Time Cost Not an Integer",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=65536,t=string,p=4$c29tZXNhbHQ$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Threads Not an Integer",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=64,t=1,p=string$c29tZXNhbHQ$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},

		// 4. Missing/Empty Hash Components
		{
			name:         "Missing Salt",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=64,t=1,p=4$$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Missing Password Hash Data",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=64,t=1,p=4$c29tZXNhbHQ$",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},

		// 5. Invalid Data Encoding
		{
			name:         "Invalid Base64 Data in Salt",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=64,t=1,p=4$invalid_base64_data$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Invalid Base64 Data in Password Hash",
			password:     validPassword,
			hash:         "$argon2id$v=19$m=64,t=1,p=4$c29tZXNhbHQ$invalid_base64_data",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Hash with Extra Fields",
			password:     validPassword,
			hash:         validHash + "$extrafield",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:         "Hash with Different Argon2 Variant",
			password:     validPassword,
			hash:         "$argon2i$v=19$m=65536,t=1,p=4$c29tZXNhbHQ$c29tZWRhdGE",
			wantVerified: false,
			wantErrCode:  models.ERROR_CODE_UNPROCESSABLE_ENTITY,
		},
		{
			name:     "Very Long Password",
			password: "This is a very long password string created to test the system's ability to handle edge cases where the password input is exceptionally lengthy. It should not cause any issues.",
			hash: func() string {
				h, _ := argon2ID.GenerateNew("This is a very long password string created to test the system's ability to handle edge cases where the password input is exceptionally lengthy. It should not cause any issues.")
				return h
			}(),
			wantVerified: true,
			wantErrCode:  "",
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

func TestVerify_ParameterMismatch(t *testing.T) {
	// A helper function to create a new KriptoArgon2ID instance for testing.
	createArgon2ID := func() *passwordHashing.KriptoArgon2ID {
		argon2Id, err := passwordHashing.NewKriptoArgon2ID(
			16, // saltLength
			3,  // iterations (timeCost)
			64, // memoryCostMegaBytes
			1,  // threads
			32, // keyLengthB
			[]byte("secure_pepper_secret_12345"),
		)
		if err != nil {
			t.Fatalf("Failed to create KriptoArgon2ID instance: %v", err)
		}
		return argon2Id
	}

	argon2Id := createArgon2ID()

	// Use a fixed password and generate a hash with the initial parameters.
	password := "testpassword123"
	fullHash, err := argon2Id.GenerateNew(password)
	if err != nil {
		t.Fatalf("Failed to generate a valid hash for testing: %v", err)
	}

	tests := []struct {
		name     string
		testFunc func() *models.CustomError
		wantErr  bool
	}{
		{
			name: "Mismatched_Salt_Length",
			testFunc: func() *models.CustomError {
				// Change the salt length in the hash string to be different from the
				// one used by the KriptoArgon2ID instance.
				parts := strings.Split(fullHash, "$")
				// A valid hash for a 16-byte salt should have a Base64-encoded string of length 22.
				// We'll change it to 20 to simulate a mismatch.
				mismatchedSaltPart := parts[4][:20]
				mismatchedHash := fmt.Sprintf("%s$%s", strings.Join(parts[:4], "$"), mismatchedSaltPart)

				// This test case would fail gracefully if the code checked the salt length,
				// but because the `Verify` function recomputes a hash with the instance's
				// salt length, this will still return an unauthorized error. It's
				// good to explicitly test this case.
				return argon2Id.Verify(password, mismatchedHash)
			},
			wantErr: true,
		},
		{
			name: "Mismatched_Key_Length_In_Hash",
			testFunc: func() *models.CustomError {
				// This test is to ensure that the `Verify` function correctly
				// recomputes the hash using the parameters from the provided hash string.
				// We'll change the key length in the hash string itself to a different value.
				// The `argon2.IDKey` function will still work because it uses the
				// hash from the string to determine the output length.
				parts := strings.Split(fullHash, "$")
				// Use a valid Base64 string to simulate a mismatched hash part.
				// A 30-byte array encoded with RawStdEncoding gives a 40-character string.
				differentHashPart := base64.RawStdEncoding.EncodeToString(make([]byte, 30))
				mismatchedHash := fmt.Sprintf("%s$%s", strings.Join(parts[:5], "$"), differentHashPart)
				return argon2Id.Verify(password, mismatchedHash)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFunc()
			if (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerify_DifferentParams(t *testing.T) {
	t.Run("Verify with different memory cost", func(t *testing.T) {
		// Generate a hash with a specific memory cost.
		argon2IdGen, _ := passwordHashing.NewKriptoArgon2ID(16, 3, 128, 1, 32, []byte("secure_pepper_secret_12345"))
		password := "another_test_pass"
		fullHash, err := argon2IdGen.GenerateNew(password)
		if err != nil {
			t.Fatalf("Failed to generate hash: %v", err)
		}

		// Create a new instance with the same parameters for verification.
		argon2IdVerify, _ := passwordHashing.NewKriptoArgon2ID(16, 3, 128, 1, 32, []byte("secure_pepper_secret_12345"))

		// The verification should be successful.
		if err := argon2IdVerify.Verify(password, fullHash); err != nil {
			t.Errorf("Verify() failed with different memory cost instance: %v", err)
		}
	})

	t.Run("Verify with different time cost", func(t *testing.T) {
		// Generate a hash with a specific time cost.
		argon2IdGen, _ := passwordHashing.NewKriptoArgon2ID(16, 5, 64, 1, 32, []byte("secure_pepper_secret_12345"))
		password := "different_time_cost_pass"
		fullHash, err := argon2IdGen.GenerateNew(password)
		if err != nil {
			t.Fatalf("Failed to generate hash: %v", err)
		}

		// Create a new instance with the same parameters for verification.
		argon2IdVerify, _ := passwordHashing.NewKriptoArgon2ID(16, 5, 64, 1, 32, []byte("secure_pepper_secret_12345"))

		// The verification should be successful.
		if err := argon2IdVerify.Verify(password, fullHash); err != nil {
			t.Errorf("Verify() failed with different time cost instance: %v", err)
		}
	})
}
