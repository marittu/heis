package elevatorDriver

import (
	"errors"
	"fmt"
)

const MOTOR_SPEED = 2800

var lampMatrix = [N_FLOORS][N_BUTTONS]int{
	[3]int{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	[3]int{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	[3]int{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	[3]int{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}

var buttonMatrix = [N_FLOORS][N_BUTTONS]int{
	[3]int{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	[3]int{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	[3]int{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	[3]int{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

func ElevInit() {
	if !ioInit() {
		errors.New("Elevator not initialized")
		return
	}
	fmt.Println("Initialized")
	for floor := 0; floor < N_FLOORS; floor++ {
			for button := BUTTON_CALL_UP; button < N_BUTTONS; button++{
				ElevSetButtonLamp(floor, button, 0) //no orders before initialization
				ElevSetDoorOpenLamp(0)
			}
	}
			
	floor := ElevGetFloorSensorSignal()
	if floor == -1 {
		ElevDrive(-1)
		
	}else{
		ElevDrive(0)
	}
	
	
}

func ElevDrive(dir ElevMotorDirection) {
	if dir == 0 {
		ioWriteAnalog(MOTOR, 0)

	} else if dir == 1 {
		ioClearBit(MOTORDIR)
		ioWriteAnalog(MOTOR, MOTOR_SPEED)

	} else if dir == -1 {
		ioSetBit(MOTORDIR)
		ioWriteAnalog(MOTOR, MOTOR_SPEED)

	}
}

func ElevSetButtonLamp(floor int, button ElevButtonType, value int) {
	if floor >= 0 && floor < N_FLOORS && button >= 0 && button < N_BUTTONS {
		if value == 1 {
			ioSetBit(lampMatrix[floor][button])
		} else {
			ioClearBit(lampMatrix[floor][button])
		}
	}
}

func ElevSetFloorIndicator(floor int) {
	if !(floor >= 0 && floor < N_FLOORS) {
		errors.New("Floor not valid")
		return
	}

	// Binary encoding. One light must always be on.
	if floor&0x02 != 0 {
		ioSetBit(LIGHT_FLOOR_IND1)
	} else {
		ioClearBit(LIGHT_FLOOR_IND1)
	}

	if floor&0x01 != 0 {
		ioSetBit(LIGHT_FLOOR_IND2)
	} else {
		ioClearBit(LIGHT_FLOOR_IND2)
	}
}

func ElevGetButtonSignal(button ElevButtonType, floor int) int {
	if floor >= 0 && floor < N_FLOORS && button >= 0 && button < N_BUTTONS {
		if ioReadBit(buttonMatrix[floor][button]) {
			return 1
		} else {
			return 0
		}

		return 0
	}
	return 0
}

func ElevGetFloorSensorSignal() int {
	if ioReadBit(SENSOR_FLOOR1) {
		return 0
	} else if ioReadBit(SENSOR_FLOOR2) {
		return 1
	} else if ioReadBit(SENSOR_FLOOR3) {
		return 2
	} else if ioReadBit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func ElevSetDoorOpenLamp(value int) {
	if value != 0 {
		ioSetBit(LIGHT_DOOR_OPEN)
	} else {
		ioClearBit(LIGHT_DOOR_OPEN)
	}
}