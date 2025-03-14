package mod_test

import (
	"testing"

	"github.com/donovanmods/projectdaedalus-db-tool/lib/mod"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	m := mod.New()
	assert.NotNil(t, m)
	assert.Equal(t, mod.StateNew, m.State())
	assert.Empty(t, m.Meta.Status.Errors)
	assert.Empty(t, m.Meta.Status.Warnings)
	assert.Equal(t, "0.1.0", m.Version)
}

func TestReset(t *testing.T) {
	m := &mod.Mod{Name: "Test Mod", Version: "1.1", Author: "Test Author"}
	assert.Equal(t, "Test Mod", m.Name)
	assert.Equal(t, "Test Author", m.Author)
	m.Reset()
	assert.Empty(t, m.Name)
	assert.Empty(t, m.Version)
	assert.Empty(t, m.Author)
	assert.Equal(t, mod.StateFresh, m.State())
	assert.Empty(t, m.Meta.Status.Errors)
	assert.Empty(t, m.Meta.Status.Warnings)
}

func TestClean(t *testing.T) {
	m := &mod.Mod{Name: "Test Mod", Version: "1.1", Author: "Test Author"}

	m.ReadmeURL = "http://example.com/readme.md?foo=bar"
	m.Files.Pak = "http://example.com/pak.zip/"
	m.Files.Exmodz = "http://example.com/exmodz.zip"

	m.Clean()

	assert.Equal(t, "Test Mod", m.Name)
	assert.Equal(t, "", m.ImageURL)
	assert.Equal(t, "http://example.com/readme.md", m.ReadmeURL)
	assert.Equal(t, "http://example.com/pak.zip", m.Files.Pak)
	assert.Equal(t, "http://example.com/exmodz.zip", m.Files.Exmodz)
	assert.Equal(t, mod.StateFresh, m.State())
	assert.Empty(t, m.Meta.Status.Errors)
	assert.Empty(t, m.Meta.Status.Warnings)
}

func TestString(t *testing.T) {
	m := &mod.Mod{Name: "Test Mod", Version: "1.0", Author: "Test Author"}
	assert.Equal(t, "Test Mod v1.0 by Test Author", m.String())
}

func TestState(t *testing.T) {
	m := &mod.Mod{}
	assert.Equal(t, mod.StateFresh, m.State())
}

func TestSetState(t *testing.T) {
	m := &mod.Mod{}
	m.SetState(mod.StateUpdated)
	assert.Equal(t, mod.StateUpdated, m.State())
}

func TestStateString(t *testing.T) {
	m := &mod.Mod{}
	assert.Equal(t, "Fresh", m.StateString())
	m.SetState(mod.StateNew)
	assert.Equal(t, "New", m.StateString())
}

func TestDirty(t *testing.T) {
	m := &mod.Mod{}
	assert.False(t, m.Dirty())
	m.SetState(mod.StateUpdated)
	assert.True(t, m.Dirty())
}

func TestValid(t *testing.T) {
	m := &mod.Mod{}
	assert.True(t, m.Valid())
	m.AddWarning("Test Warning")
	assert.True(t, m.Valid())
	m.AddError("Test Error")
	assert.False(t, m.Valid())
}
