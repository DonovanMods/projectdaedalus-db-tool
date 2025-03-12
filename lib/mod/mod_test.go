package mod_test

import (
	"testing"

	"github.com/donovanmods/projectdaedalus-db-tool/lib/mod"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	m := (&mod.Mod{}).New()
	assert.NotNil(t, m)
	assert.Equal(t, mod.New, m.Meta.State)
	assert.Empty(t, m.Meta.Status["Errors"])
	assert.Empty(t, m.Meta.Status["Warnings"])
}

func TestString(t *testing.T) {
	m := &mod.Mod{Name: "Test Mod", Version: "1.0", Author: "Test Author"}
	assert.Equal(t, "Test Mod v1.0 by Test Author", m.String())
}

func TestToJSON(t *testing.T) {
	m := &mod.Mod{Name: "Test Mod"}
	data, err := m.ToJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"Test Mod"`)
}

func TestFromJSON(t *testing.T) {
	data := []byte(`{"name":"Test Mod"}`)
	m := &mod.Mod{}
	err := m.FromJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, "Test Mod", m.Name)
}

func TestState(t *testing.T) {
	m := &mod.Mod{}
	assert.Equal(t, mod.Unmodified, m.Meta.State)
}

func TestDirty(t *testing.T) {
	m := &mod.Mod{}
	assert.False(t, m.Dirty())
	m.Meta.State = mod.Updated
	assert.True(t, m.Dirty())
}
