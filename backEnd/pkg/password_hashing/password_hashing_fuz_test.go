package passwordHashing

import (
	"strings"
	"testing"
)

// FuzzGetArgon2IdHashParts fuzzer for the getArgon2IdHashParts function
func FuzzGetArgon2IdHashParts(f *testing.F) {
	// Add a comprehensive set of "seed" inputs to guide the fuzzer.
	// These inputs cover both valid and various invalid formats.
	f.Add("$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNraWxs$dG9naXRodWI")                 // Valid hash
	f.Add("$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNraWlsbA$dG9naXRodWI")               // Valid with different salt length
	f.Add("not a hash string")                                                       // Completely invalid format
	f.Add("$argon2id$v=19$m=65536,t=1,p=4")                                          // Missing salt and hash
	f.Add("$argon2id$v=invalid$m=65536,t=1,p=4$c2FsdHNraWxs$dG9naXRodWI")            // Invalid version number
	f.Add("$argon2id$v=19$m=65536,t=1$c2FsdHNraWxs$dG9naXRodWI")                     // Missing 'p' cost factor
	f.Add("$argon2id$v=19$m=65536,p=4$c2FsdHNraWxs$dG9naXRodWI")                     // Missing 't' cost factor
	f.Add("$argon2id$v=19$t=1,p=4$c2FsdHNraWxs$dG9naXRodWI")                         // Missing 'm' cost factor
	f.Add("$argon2id$v=19$m=65536,t=1,p=4,extra=1$c2FsdHNraWxs$dG9naXRodWI")         // Extra cost factor
	f.Add("$argon2id$v=19$m=65536,t=1,p=$c2FsdHNraWxs$dG9naXRodWI")                  // Missing cost factor value
	f.Add("$argon2id$v=19$m=65536,t=1,p=4$invalid-base64-salt$dG9naXRodWI")          // Invalid base64 for salt
	f.Add("$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNraWxs$invalid-base64-hash")         // Invalid base64 for hash
	f.Add("$argon2id$v=19$m=65536,t=1,p=4$$dG9naXRodWI")                             // Empty salt
	f.Add("$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNraWxs$")                            // Empty hash
	f.Add("too$many$parts$in$this$hash$string")                                      // Too many parts
	f.Add("no-dollar-prefix$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNraWxs$dG9naXRodWI") // Hash starts with a character instead of a dollar sign
	f.Add("$argon2d$v=19$m=65536,t=1,p=4$c2FsdHNraWxs$dG9naXRodWI")                  // Incorrect algorithm name ("argon2d" instead of "argon2id")
	f.Add("$argon2id$v=abc$m=65536,t=1,p=4$c2FsdHNraWxs$dG9naXRodWI")                // Non-numeric version
	f.Add("$argon2id$v=19$m=65536,t=1$c2FsdHNraWxs$dG9naXRodWI")                     // Missing cost factor 'p'
	f.Add("$argon2id$v=19$t=1,p=4,m=65536$c2FsdHNraWxs$dG9naXRodWI")                 // Cost factors are out of order
	f.Add("$argon2id$v=19$m=065536,t=1,p=4$c2FsdHNraWxs$dG9naXRodWI")                // Leading zero on a cost factor value
	f.Add("$argon2id$v=19$m=65536,t=1a,p=4$c2FsdHNraWxs$dG9naXRodWI")                // Non-numeric cost factor value
	f.Add("$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNraWxs==$dG9naXRodWI")               // Base64 salt with padding, which is not allowed
	f.Add("$argon2id$v=19$m=65536,t=1,p=4$c2FsdHNraWxs$dG9naXRodWI==")               // Base64 hash with padding, which is not allowed

	// This is the fuzzing target. It receives the test inputs.
	f.Fuzz(func(t *testing.T, fullhash string) {
		// Call the function being fuzzed.
		result, err := getArgon2IdHashParts(fullhash)

		// This is a crucial check for a robust fuzz test.
		// If no error occurred, we assume the input was a valid hash.
		// We then perform a "round-trip" test by calling the String() method
		// on the parsed result. The reconstructed string should be identical
		// to the original input. This validates not just crash resistance,
		// but also the correctness and completeness of the parsing logic.
		if err == nil {
			if reconstructedHash := result.String(); reconstructedHash != fullhash {
				t.Errorf("Round-trip failed: original hash '%s', reconstructed hash '%s'", fullhash, reconstructedHash)
			}
		} else {
			// If an error occurred, we check for a specific type of failure.
			// This helps to confirm that the function handles invalid inputs
			// gracefully and returns the expected error type.
			if result != nil {
				t.Errorf("Expected nil result on error, but got a non-nil result")
			}
			// Note: You can add more specific error type checks here if needed,
			// for example, to verify which type of custom error was returned.
		}
	})
}

// FuzzSplitParamOnEquals is a fuzzer for the splitParamOnEquals function.
// It generates test inputs and calls the function being fuzzed.
func FuzzSplitParamOnEquals(f *testing.F) {
	f.Add("key=123", "key")
	f.Add("r=8", "r")
	f.Add("t=2", "t")
	f.Add("m=1024", "m")
	f.Add("p=1", "p")
	f.Add("k=123456789", "k")
	f.Add("bad-format", "")
	f.Add("wrong=prefix", "right")

	f.Fuzz(func(t *testing.T, str string, paramPrefix string) {
		parts, err := splitParamOnEquals(str, paramPrefix)
		if err == nil {
			if len(parts) != 2 {
				t.Fatalf("Expected 2 parts, got %d for input '%s' and prefix '%s'", len(parts), str, paramPrefix)
			}
		}
	})
}

// FuzzGetCostFactors is a fuzzer for the getCostFactors function.
func FuzzGetCostFactors(f *testing.F) {
	f.Add("m=8192,t=2,p=1")         // Valid input
	f.Add("m=1024,t=1,p=4")         // Another valid input
	f.Add("m=123,t=456,p=789")      // Another valid input
	f.Add("m=123,t=456")            // Too few parts
	f.Add("m=123,t=456,p=789,o=10") // Too many parts
	f.Add("x=123,y=456,z=789")      // Wrong prefixes
	f.Add("m=abc,t=2,p=1")          // Non-numeric value
	f.Add("m=01,t=2,p=1")           // Leading zero
	f.Add("m=1,t=,p=1")             // Missing value
	f.Add("m=1,t=2p=1")             // Missing equals sign

	f.Fuzz(func(t *testing.T, costFactors string) {
		costFactorsMap, err := getCostFactors(costFactors)
		if err == nil {
			if len(costFactorsMap) != 3 {
				t.Fatalf("Expected 3 cost factors, got %d for input '%s'", len(costFactorsMap), costFactors)
			}
		}
	})
}

// FuzzGenerateNew is a fuzzer for the GenerateNew function.
func FuzzGenerateNew(f *testing.F) {
	pepperSecret := []byte("a-secure-pepper-for-testing")
	kripto, err := NewKriptoArgon2ID(
		32,
		4,
		64,
		1,
		32,
		pepperSecret,
	)
	if err != nil {
		f.Fatalf("Failed to create KriptoArgon2ID instance: %v", err)
	}
	f.Add("")                                   // Empty string
	f.Add("password")                           // A simple, common password
	f.Add("S3cUr3Pa$$w0rd!123")                 // A password with special characters
	f.Add(strings.Repeat("a", 1000))            // A very long password
	f.Add("测试密码")                               // Unicode characters
	f.Add("\x00\x01\x02\x03")                   // Binary data
	f.Add("!@#$%^&*()_+-=")                     // All special characters
	f.Add("1234567890")                         // All numeric characters
	f.Add(strings.Repeat("B", 1024))            // A password that crosses a memory boundary
	f.Add("`~!@#$%^&*()_-+=|\\}[{]:;'<>,./?\"") // Punctuation characters
	f.Add("\n\r\t")                             // Whitespace characters

	f.Fuzz(func(t *testing.T, password string) {
		fullHash, err := kripto.GenerateNew(password)

		if err != nil {
			t.Fatalf("GenerateNew failed for password '%s': %v", password, err)
		}
		verifyErr := kripto.Verify(password, fullHash)

		if verifyErr != nil {
			t.Fatalf("Verification failed for password '%s' with hash '%s': %v", password, fullHash, verifyErr)
		}
	})
}

// FuzzGenerateSalt is the fuzz test function.
func FuzzGenerateSalt(f *testing.F) {
	f.Add(uint8(0))   // A zero length
	f.Add(uint8(16))  // A common salt length
	f.Add(uint8(32))  // Another common length
	f.Add(uint8(100)) // A larger length
	f.Add(uint8(255)) // The maximum possible length for uint8

	f.Fuzz(func(t *testing.T, saltLength uint8) {
		argon := &KriptoArgon2ID{
			saltLength: saltLength,
		}
		salt, err := argon.generateSalt()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}
		if len(salt) != int(saltLength) {
			t.Fatalf("Expected salt of length %d, but got %d", saltLength, len(salt))
		}
	})
}

// FuzzNewKriptoArgon2ID is a fuzzer for the NewKriptoArgon2ID function.
func FuzzNewKriptoArgon2ID(f *testing.F) {
	f.Add(uint8(16), uint32(1), uint32(64), uint8(1), uint32(32), []byte("a very secure pepper!"))                              // A valid, standard case
	f.Add(uint8(15), uint32(1), uint32(64), uint8(1), uint32(32), []byte("a very secure pepper!"))                              // Invalid saltLength
	f.Add(uint8(16), uint32(1), uint32(63), uint8(1), uint32(32), []byte("a very secure pepper!"))                              // Invalid memoryCost
	f.Add(uint8(16), uint32(1), uint32(64), uint8(1), uint32(31), []byte("a very secure pepper!"))                              // Invalid keyLength
	f.Add(uint8(16), uint32(1), uint32(64), uint8(1), uint32(32), []byte("short"))                                              // Invalid pepperSecret length
	f.Add(uint8(255), uint32(4294967295), uint32(4294967295), uint8(255), uint32(4294967295), []byte("pepper with max length")) // Max values

	f.Fuzz(func(t *testing.T, saltLength uint8, iterations uint32, memoryCostMegaBytes uint32, threads uint8, keyLengthB uint32, pepperSecret []byte) {
		argon2Id, err := NewKriptoArgon2ID(saltLength, iterations, memoryCostMegaBytes, threads, keyLengthB, pepperSecret)
		switch {
		case saltLength >= 16 && memoryCostMegaBytes >= 64 && keyLengthB >= 32 && len(pepperSecret) >= 16:
			if err != nil {
				t.Errorf("Expected no error for valid parameters, but got %v", err)
			}
			if argon2Id == nil {
				t.Error("Expected a non-nil KriptoArgon2ID instance for valid parameters")
			}
			if argon2Id != nil && (argon2Id.saltLength != saltLength || argon2Id.timeCost != iterations || argon2Id.memoryCostKiloBytes != memoryCostMegaBytes*1024 || argon2Id.threads != threads || argon2Id.outputkeyLengthBytes != keyLengthB) {
				t.Errorf("Returned struct has incorrect values")
			}
		case saltLength < 16 || memoryCostMegaBytes < 64 || keyLengthB < 32 || len(pepperSecret) < 16:
			if err == nil {
				t.Errorf("Expected an error for invalid parameters, but got none")
			}
			if argon2Id != nil {
				t.Errorf("Expected a nil KriptoArgon2ID instance for invalid parameters, but got %v", argon2Id)
			}
		}
	})
}
