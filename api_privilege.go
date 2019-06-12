package libgen

/**
 * Created by xuchao on 2019-06-12 .
 */

type PrivilegeCode uint32

const (
	//privilege for manage device model ,
	//DeclareDeviceModel | RemoveDeviceModels | UpdateDeviceModel
	ManageDeviceModel PrivilegeCode = iota

	//privilege for manage device ,
	//

	ManageDevice
	//privilege for retrieve devices from other domain,
	RetrieveDevice
)
