package utils

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

const (
	FORWARD = iota
	BACKWARD
	LEFT
	RIGHT
)
const YAW = -90.0
const PITCH = 0.0
const SPEED = 3.0
const SENSITIVTY = 0.25
const ZOOM = 45.0

const (
    RadToDeg = 180/math.Pi
    DegToRad = math.Pi/180
)

type Camera struct {
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Up       mgl32.Vec3
	Right    mgl32.Vec3
	WorldUp  mgl32.Vec3
	// Eular Angles
	Yaw   float32
	Pitch float32
	// Camera options
	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

func NewCamera(position, up mgl32.Vec3, yaw, pitch float32) Camera{
	cam := Camera{
		Position: position,
		WorldUp: up,
		Yaw: yaw,
		Pitch: pitch,
		Front: mgl32.Vec3{0.0, 0.0, -1.0},
		MovementSpeed: SPEED,
		MouseSensitivity: SENSITIVTY,
		Zoom: ZOOM,
	}
	cam.updateCameraVectors()
	return cam
}

// Constructor with scalar values
func NewCameraWithScalars(posX, posY, posZ, upX, upY, upZ, yaw, pitch float32) Camera {
	cam := Camera{
		Position: mgl32.Vec3{posX, posY, posZ},
		WorldUp: mgl32.Vec3{upX, upY, upZ},
		Yaw: yaw,
		Pitch: pitch,
		Front: mgl32.Vec3{0.0, 0.0, -1.0},
		MovementSpeed: SPEED,
		MouseSensitivity: SENSITIVTY,
		Zoom: ZOOM,
	}
	cam.updateCameraVectors()
	return cam
}

func (cam *Camera)  GetViewMatrix() mgl32.Mat4 {
	eye := cam.Position
	center := cam.Position.Add(cam.Front)
	up := cam.Up
	return mgl32.LookAt(
		eye.X(), eye.Y(), eye.Z(),
		center.X(), center.Y(), center.Z(),
		up.X(), up.Y(), up.Z())
}

// Processes input received from any keyboard-like input system. Accepts input parameter in the form of camera defined ENUM (to abstract it from windowing systems)
func (cam *Camera) ProcessKeyboard(direction int, deltaTime float32) {
	velocity := float32(cam.MovementSpeed) * deltaTime
	if (direction == FORWARD) {
		cam.Position = cam.Position.Add(cam.Front.Mul(velocity))
	}
	if (direction == BACKWARD) {
		cam.Position = cam.Position.Sub(cam.Front.Mul(velocity))
	}
	if (direction == LEFT) {
		cam.Position = cam.Position.Sub(cam.Right.Mul(velocity))
	}
	if (direction == RIGHT) {
		cam.Position = cam.Position.Add(cam.Right.Mul(velocity))
	}
}

// Processes input received from a mouse input system. Expects the offset value in both the x and y direction.
func (cam *Camera) ProcessMouseMovement(xoffset, yoffset float32, constrainPitch bool) {
	xoffset *= cam.MouseSensitivity
	yoffset *= cam.MouseSensitivity

	cam.Yaw   += xoffset
	cam.Pitch += yoffset

	// Make sure that when pitch is out of bounds, screen doesn't get flipped
	if (constrainPitch){
		if (cam.Pitch > 89.0) {
			cam.Pitch = 89.0
		}
		if (cam.Pitch < -89.0) {
			cam.Pitch = -89.0
		}
	}
	// Update Front, Right and Up Vectors using the updated Eular angles
	cam.updateCameraVectors()
}

// Processes input received from a mouse scroll-wheel event. Only requires input on the vertical wheel-axis
func (cam *Camera) ProcessMouseScroll(yoffset float32) {
	if (cam.Zoom >= 1.0 && cam.Zoom <= 45.0) {
		cam.Zoom -= yoffset
	}
	if (cam.Zoom <= 1.0) {
		cam.Zoom = 1.0
	}
	if (cam.Zoom >= 45.0) {
		cam.Zoom = 45.0
	}
}

func (cam *Camera) updateCameraVectors() {
	
	x := cos(cam.Yaw * DegToRad) * cos(cam.Pitch * DegToRad)
	y := sin(cam.Pitch * DegToRad)
	z := sin(cam.Yaw * DegToRad) * cos(cam.Pitch * DegToRad)
	front := mgl32.Vec3{x, y, z}
	front = front.Normalize()
	// Also re-calculate the Right and Up vector
	// Normalize the vectors, because their length gets closer to 0 the more you look up or down which results in slower movement.
	cam.Right = front.Cross(cam.WorldUp).Normalize()
	cam.Up = cam.Right.Cross(cam.Front).Normalize()
}

func cos(f float32) float32 {
	return float32(math.Cos(float64(f)))
}
func sin(f float32) float32 {
	return float32(math.Sin(float64(f)))
}