package image

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	// Setup
	filename := "./testdata/full-valid-example.yaml"
	configData, err := os.ReadFile(filename)
	require.NoError(t, err)

	// Test
	definition, err := ParseDefinition(configData)

	// Verify
	require.NoError(t, err)

	// - Definition
	assert.Equal(t, "1.0", definition.APIVersion)
	assert.EqualValues(t, "x86_64", definition.Image.Arch)
	assert.Equal(t, "iso", definition.Image.ImageType)

	// - Image
	assert.Equal(t, "slemicro5.5.iso", definition.Image.BaseImage)
	assert.Equal(t, "eibimage.iso", definition.Image.OutputImageName)

	// - Operating System -> Kernel Arguments
	expectedKernelArgs := []string{
		"alpha=foo",
		"beta=bar",
		"baz",
	}
	assert.Equal(t, expectedKernelArgs, definition.OperatingSystem.KernelArgs)

	// Operating System -> Users
	userConfigs := definition.OperatingSystem.Users
	require.Len(t, userConfigs, 3)
	assert.Equal(t, "alpha", userConfigs[0].Username)
	assert.Equal(t, "$6$bZfTI3Wj05fdxQcB$W1HJQTKw/MaGTCwK75ic9putEquJvYO7vMnDBVAfuAMFW58/79abky4mx9.8znK0UZwSKng9dVosnYQR1toH71", userConfigs[0].EncryptedPassword)
	assert.Contains(t, userConfigs[0].SSHKey, "ssh-rsa AAAAB3")
	assert.Equal(t, "beta", userConfigs[1].Username)
	assert.Equal(t, "$6$GHjiVHm2AT.Qxznz$1CwDuEBM1546E/sVE1Gn1y4JoGzW58wrckyx3jj2QnphFmceS6b/qFtkjw1cp7LSJNW1OcLe/EeIxDDHqZU6o1", userConfigs[1].EncryptedPassword)
	assert.Equal(t, "", userConfigs[1].SSHKey)
	assert.Equal(t, "gamma", userConfigs[2].Username)
	assert.Equal(t, "", userConfigs[2].EncryptedPassword)
	assert.Contains(t, userConfigs[2].SSHKey, "ssh-rsa BBBBB3")

	// Operating System -> Systemd
	systemd := definition.OperatingSystem.Systemd
	require.Len(t, systemd.Enable, 2)
	assert.Equal(t, "enable0", systemd.Enable[0])
	assert.Equal(t, "enable1", systemd.Enable[1])
	require.Len(t, systemd.Disable, 1)
	assert.Equal(t, "disable0", systemd.Disable[0])

	// Operating System -> Suma
	suma := definition.OperatingSystem.Suma
	assert.Equal(t, "suma.edge.suse.com", suma.Host)
	assert.Equal(t, "slemicro55", suma.ActivationKey)
	assert.Equal(t, false, suma.GetSSL)

	// EmbeddedArtifactRegistry
	embeddedArtifactRegistry := definition.EmbeddedArtifactRegistry
	assert.Equal(t, "hello-world:latest", embeddedArtifactRegistry.ContainerImages[0].Name)
	assert.Equal(t, "rgcrprod.azurecr.us/longhornio/longhorn-ui:v1.5.1", embeddedArtifactRegistry.ContainerImages[1].Name)
	assert.Equal(t, "carbide-key.pub", embeddedArtifactRegistry.ContainerImages[1].SupplyChainKey)
	assert.Equal(t, "rancher", embeddedArtifactRegistry.HelmCharts[0].Name)
	assert.Equal(t, "https://releases.rancher.com/server-charts/stable", embeddedArtifactRegistry.HelmCharts[0].RepoURL)
	assert.Equal(t, "2.8.0", embeddedArtifactRegistry.HelmCharts[0].Version)

	// Kubernetes
	kubernetes := definition.Kubernetes
	assert.Equal(t, "v1.29.0+rke2r1", kubernetes.Version)
	assert.Equal(t, "server", kubernetes.NodeType)
	assert.Equal(t, "cilium", kubernetes.CNI)
	assert.Equal(t, true, kubernetes.MultusEnabled)
	assert.Equal(t, false, kubernetes.VSphereEnabled)
}

func TestParseBadConfig(t *testing.T) {
	// Setup
	badData := []byte("Not actually YAML")

	// Test
	_, err := ParseDefinition(badData)

	// Verify
	require.Error(t, err)
	assert.ErrorContains(t, err, "could not parse the image definition")
}

func TestArch_Short(t *testing.T) {
	assert.Equal(t, "amd64", ArchTypeX86.Short())
	assert.Equal(t, "arm64", ArchTypeARM.Short())
	assert.PanicsWithValue(t, "unknown arch: abc", func() {
		arch := Arch("abc")
		arch.Short()
	})
}
