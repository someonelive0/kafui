package backend

import "testing"

func TestSetting(t *testing.T) {
	myconfig, err := LoadConfig("../kafui.toml")
	if err != nil {
		t.Fatalf("LoadSetting failed: %s", err)
	}

	t.Logf("myconfig: %#v", *myconfig)
	t.Logf("myconfig: %s", myconfig.Dump())

	if err = SaveConfig(myconfig); err != nil {
		t.Fatalf("SaveConfig failed: %s", err)
	}
}
