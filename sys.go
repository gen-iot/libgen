package libgen

func SysSendMsg(msg *IOMsg) error {
	if msg == nil {
		return ErrIllegalParam
	}
	return nil
}

func SysRecvMsg() (*IOMsg, error) {
	return nil, nil
}
