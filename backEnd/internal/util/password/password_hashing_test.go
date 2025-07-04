package password_test

import (
	"testing"

	passwordHashing "github.com/413ksz/BlueFox/backEnd/internal/util/password"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid Password (short)",
			password: "mysecretpassword123",
			wantErr:  false,
		},
		{
			name:     "Empty Password",
			password: "",
			wantErr:  false, // bcrypt can hash empty strings
		},
		{
			name:     "Password Exact Max Length (72 bytes)",
			password: "thisisasecretpasswordofexactlyseventytwocharacterslongabcdefghijklmno", // 72 chars
			wantErr:  false,
		},
		{
			name:     "Password Exceeds Max Length (73 bytes)",
			password: "thisisareallylongpasswordthatshouldstillbehashedcorrectlybythebcryptalgorithmanditsexactly73characterss", // 73 chars
			wantErr:  true,                                                                                                      // EXPECT AN ERROR HERE NOW
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := passwordHashing.HashPassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() for '%s' error = %v, wantErr %v", tt.password, err, tt.wantErr)
				return
			}

			if !tt.wantErr { // If no error was expected, proceed to verify
				if hashedPassword == "" {
					t.Errorf("HashPassword() returned an empty hash for password: %s", tt.password)
				}

				// Verify that the generated hash is valid using bcrypt's own verification
				// This ensures our HashPassword wrapper is producing valid bcrypt hashes
				err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tt.password))
				if err != nil {
					t.Errorf("HashPassword() generated an invalid hash for '%s': %v", tt.password, err)
				}
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	// Generate a known good hash for testing once.
	validPassword := "testpassword123"
	// Use the HashPassword function from the package being tested to generate a valid hash.
	// This ensures the test hash is created with the same logic/cost as production.
	validHash, err := passwordHashing.HashPassword(validPassword)
	if err != nil {
		t.Fatalf("Failed to generate test hash for VerifyPassword tests: %v", err)
	}

	tests := []struct {
		name         string
		password     string
		hash         string
		wantVerified bool
	}{
		{
			name:         "Correct Password",
			password:     validPassword,
			hash:         validHash,
			wantVerified: true,
		},
		{
			name:         "Incorrect Password",
			password:     "wrongpassword",
			hash:         validHash,
			wantVerified: false,
		},
		{
			name:     "Empty Password (correct against empty string hash)",
			password: "",
			// Generate a hash for an empty string directly for this specific case
			// since VerifyPassword doesn't hash, it just compares.
			hash:         func() string { h, _ := bcrypt.GenerateFromPassword([]byte(""), 14); return string(h) }(),
			wantVerified: true,
		},
		{
			name:         "Empty Password (incorrect against non-empty hash)",
			password:     "",
			hash:         validHash,
			wantVerified: false,
		},
		{
			name:         "Empty Hash",
			password:     validPassword,
			hash:         "",
			wantVerified: false, // bcrypt.CompareHashAndPassword returns an error for invalid hash
		},
		{
			name:         "Malformed Hash",
			password:     validPassword,
			hash:         "notavalidhash",
			wantVerified: false, // bcrypt.CompareHashAndPassword returns an error for invalid hash
		},
		{
			name:         "Hash from different password",
			password:     "anotherpassword",
			hash:         validHash, // Using validHash for a different password
			wantVerified: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function from the imported package
			gotVerified := passwordHashing.VerifyPassword(tt.password, tt.hash)
			if gotVerified != tt.wantVerified {
				t.Errorf("VerifyPassword() gotVerified = %v, want %v for password '%s' and hash '%s'", gotVerified, tt.wantVerified, tt.password, tt.hash)
			}
		})
	}
}

func TestHashPassword_Concurrency(t *testing.T) {
	// This test checks if HashPassword behaves correctly under concurrent access.
	numGoroutines := 100
	password := "concurrentpassword"
	hashes := make(chan string, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			// Call HashPassword from the imported package
			hash, err := passwordHashing.HashPassword(password)
			if err != nil {
				errors <- err
				return
			}
			hashes <- hash
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-errors:
			t.Errorf("Concurrency test failed with error: %v", err)
			return
		case hash := <-hashes:
			// Call VerifyPassword from the imported package
			if !passwordHashing.VerifyPassword(password, hash) {
				t.Errorf("Concurrency test: Failed to verify a concurrently generated hash.")
				return
			}
		}
	}
	close(hashes)
	close(errors)
}
