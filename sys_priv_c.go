//+build child_lib

package libgen

func (this *sys) sysWorkMode() workMode {
	return modeChild
}

func (this *sys) sysGetLinkId() (string, error) {
	return "", nil
}
