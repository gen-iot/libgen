//+build parent_lib

package libgen

func (this *sys) sysWorkMode() workMode {
	return modeParent
}

func (this *sys) SysGrabChild() error {
	return nil
}
