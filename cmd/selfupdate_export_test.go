package cmd

// DetectInstallMethodForTest exposes the unexported detectInstallMethod
// function so it can be exercised from tests without network or exec calls.
func DetectInstallMethodForTest(exePath string) string {
	return detectInstallMethod(exePath)
}
