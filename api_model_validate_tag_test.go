package libgen

import (
	"github.com/gen-iot/std"
	"testing"
)

func Test_RegisterDevicesRequest_Validate(t *testing.T) {
	req := &RegisterDevicesRequest{
		Devices: []*Device{
			{
				ModelInfo: &ModelInfo{
					Package: "com.pujie88.com",
					Name:    "test",
				},
				Id:         "1",
				Name:       "2",
				Room:       "3",
				Properties: nil,
				MetaData:   nil,
			},
		},
	}
	err := std.ValidateStruct(req)
	std.Assert(err == nil, "1 validate should be success")
	req = &RegisterDevicesRequest{
		Devices: nil,
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "2 validate should be failed")
	req = &RegisterDevicesRequest{
		Devices: []*Device{},
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "3 validate should be failed")
	req = &RegisterDevicesRequest{
		Devices: []*Device{
			{
				ModelInfo: &ModelInfo{
					Package: "com.pujie88.iot",
					Name:    "test",
				},
				Id:         "",
				Name:       "",
				Room:       "",
				Properties: nil,
				MetaData:   nil,
			},
		},
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "4 validate should be failed")
}

func Test_ListDevicesByRoomRequest_Validate(t *testing.T) {
	req := &ListDevicesByRoomRequest{
		Rooms:    []string{"101"},
		Includes: nil,
	}
	err := std.ValidateStruct(req)
	std.Assert(err == nil, "1 validate should be success")
	req = &ListDevicesByRoomRequest{
		Rooms:    []string{},
		Includes: nil,
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "2 validate should be failed")

	req = &ListDevicesByRoomRequest{
		Rooms:    []string{"101"},
		Includes: []*ModelInfo{},
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "3 validate should be failed")

	req = &ListDevicesByRoomRequest{
		Rooms: []string{"101"},
		Includes: []*ModelInfo{
			{
				Package: "",
				Name:    "",
			},
		},
	}
	err = std.ValidateStruct(req)
	std.Assert(err != nil, "4 validate should be failed")
	req = &ListDevicesByRoomRequest{
		Rooms: []string{"101"},
		Includes: []*ModelInfo{
			{
				Package: "com.pujie88.iot",
				Name:    "test",
			},
		},
	}
	err = std.ValidateStruct(req)
	std.Assert(err == nil, "5 validate should be success")
}
