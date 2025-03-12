package mod_test

import (
	"testing"

	"github.com/donovanmods/projectdaedalus-db-tool/lib/mod"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	m := (&mod.Mod{}).New()
	assert.NotNil(t, m)
	assert.Equal(t, mod.StateNew, m.State())
	assert.Empty(t, m.Meta.Status.Errors)
	assert.Empty(t, m.Meta.Status.Warnings)
}

func TestString(t *testing.T) {
	m := &mod.Mod{Name: "Test Mod", Version: "1.0", Author: "Test Author"}
	assert.Equal(t, "Test Mod v1.0 by Test Author", m.String())
}

func TestState(t *testing.T) {
	m := &mod.Mod{}
	assert.Equal(t, mod.StateUnmodified, m.State())
}

func TestSetState(t *testing.T) {
	m := &mod.Mod{}
	m.SetState(mod.StateUpdated)
	assert.Equal(t, mod.StateUpdated, m.State())
}

func TestStateString(t *testing.T) {
	m := &mod.Mod{}
	assert.Equal(t, "Unmodified", m.StateString())
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
