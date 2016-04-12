package elevatorDriver

/*
#cgo CFLAGS: -std=gnu11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "elev.h"
*/
import "C"

func ElevInit(elevatorType int) {
	C.elev_init(C.elev_type(elevatorType))
}

func ElevDrive(dir ElevMotorDirection) {
	C.elev_set_motor_direction(C.elev_motor_direction_t(dir))
}

func ElevSetButtonLamp(floor int, button ElevButtonType, value int) {
	C.elev_set_button_lamp(C.elev_button_type_t(button), C.int(floor), C.int(value))
}

func ElevSetFloorIndicator(floor int) {
	C.elev_set_floor_indicator(C.int(floor))
}

func ElevGetButtonSignal(button ElevButtonType, floor int) int {
	return int(C.elev_get_button_signal(C.elev_button_type_t(button), C.int(floor)))
}

func ElevGetFloorSensorSignal() int {
	return int(C.elev_get_floor_sensor_signal())
}

func ElevSetDoorOpenLamp(value int) {
	C.elev_set_door_open_lamp(C.int(value))
}
