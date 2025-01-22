package signup

import (
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/codeready-toolchain/registration-service/pkg/log"
	"github.com/gin-gonic/gin"
)

// Signup represents Signup resource which is a wrapper of K8s UserSignup
// and the corresponding MasterUserRecord resources.
type Signup struct {
	// The UserSignup resource name
	Name string `json:"name"`
	// The Web Console URL of the cluster which the user was provisioned to
	ConsoleURL string `json:"consoleURL,omitempty"`
	// The Che Dashboard URL of the cluster which the user was provisioned to
	CheDashboardURL string `json:"cheDashboardURL,omitempty"`
	// The proxy URL of the cluster
	ProxyURL string `json:"proxyURL,omitempty"`
	// The RHODS URL for the user's cluster
	RHODSMemberURL string `json:"rhodsMemberURL,omitempty"`
	// The server api URL of the cluster which the user was provisioned to
	APIEndpoint string `json:"apiEndpoint,omitempty"`
	// The name of the cluster which the user was provisioned to
	ClusterName string `json:"clusterName,omitempty"`
	// The user's default namespace
	DefaultUserNamespace string `json:"defaultUserNamespace,omitempty"`
	// The complaint username.  This may differ from the corresponding Identity Provider username, because of the
	// limited character set available for naming (see RFC1123) in K8s. If the username contains characters which are
	// disqualified from the resource name, the username is transformed into an acceptable resource name instead.
	// For example, johnsmith@redhat.com -> johnsmith
	CompliantUsername string `json:"compliantUsername"`
	// Original username from the Identity Provider
	Username string `json:"username"`
	// GivenName from the Identity Provider
	GivenName string `json:"givenName"`
	// FamilyName from the Identity Provider
	FamilyName string `json:"familyName"`
	// Company from the Identity Provider
	Company string `json:"company"`
	Status  Status `json:"status,omitempty"`
	// StartDate is the date that the user's current subscription started, in RFC3339 format
	StartDate string `json:"startDate,omitempty"`
	// End Date is the date that the user's current subscription will end, in RFC3339 format
	EndDate string `json:"endDate,omitempty"`
}

// Status represents UserSignup resource status
type Status struct {
	// If true then the corresponding user's account is ready to be used
	Ready bool `json:"ready"`
	// Brief reason for the status last transition.
	Reason string `json:"reason"`
	// Human readable message indicating details about last transition.
	Message string `json:"message,omitempty"`
	// VerificationRequired is used to determine if a user requires phone verification.
	// The user should not be provisioned if VerificationRequired is set to true.
	// VerificationRequired is set to false when the user is ether exempt from phone verification or has already successfully passed the verification.
	// Default value is false.
	VerificationRequired bool `json:"verificationRequired"`
}

// PollUpdateSignup will attempt to execute the provided updater function, and if it fails
// will reattempt the update for a limited number of retries
func PollUpdateSignup(ctx *gin.Context, updater func() error) error {
	// Attempt to execute an update function, retrying a number of times if the update fails
	attempts := 0
	for {
		attempts++

		// Attempt the update
		updateErr := updater()

		// If there was an error, then only log it for now
		if updateErr != nil {
			log.Error(ctx, updateErr, fmt.Sprintf("error while executing updating, attempt #%d", attempts))
		} else {
			// Otherwise if there was no error executing the update, then break here
			break
		}

		// If we've exceeded the number of attempts, then return a useful error to the user.  We won't return the actual
		// error to the user here, as we've already logged it
		if attempts > 4 {
			return updateErr
		}
	}

	return nil
}

const DNS1123NameMaximumLength = 63

// EncodeUserIdentifier transforms a subject value (the user's username) to make it DNS-1123 compliant,
// by removing invalid characters, trimming the length and prefixing with a CRC32 checksum if required.
// ### WARNING ### changing this function will cause breakage, as it is used to lookup existing UserSignup
// resources.  If a change is absolutely required, then all existing UserSignup instances must be migrated
// to the new value
func EncodeUserIdentifier(subject string) string {
	// Sanitize subject to be compliant with DNS labels format (RFC-1123)
	encoded := sanitizeDNS1123(subject)

	// Add a checksum prefix if the encoded value is different to the original subject value
	if encoded != subject {
		encoded = fmt.Sprintf("%x-%s", crc32.Checksum([]byte(subject), crc32.IEEETable), encoded)
	}

	// Trim if the length exceeds the maximum
	if len(encoded) > DNS1123NameMaximumLength {
		encoded = encoded[0:DNS1123NameMaximumLength]
	}

	return encoded
}

func sanitizeDNS1123(str string) string {
	// convert to lowercase
	lstr := strings.ToLower(str)

	// remove unwanted characters
	b := strings.Builder{}
	for _, r := range lstr {
		switch {
		case r >= '0' && r <= '9':
			fallthrough
		case r >= 'a' && r <= 'z':
			fallthrough
		case r == '-':
			b.WriteRune(r)
		}
	}

	// remove leading and trailing '-'
	return strings.Trim(b.String(), "-")
}
