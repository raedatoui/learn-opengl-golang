package utils

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
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
func NewCameraWithScalars(posX, posY, posZ, upX, upY, upZ, yaw, pitch float32) {
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

func (cam *Camera) updateCameraVectors() {
	
	x := math.Cos(cam.Yaw * DegToRad) * math.Cos(cam.Pitch * DegToRad)
	y := math.Sin(cam.Pitch * DegToRad)
	z := math.Sin(cam.Yaw * DegToRad) * math.Cos(cam.Pitch * DegToRad)
	front := mgl32.Vec3{x, y, z}
	front = front.Normalize()
	// Also re-calculate the Right and Up vector
	// Normalize the vectors, because their length gets closer to 0 the more you look up or down which results in slower movement.
	cam.Right = front.Cross(cam.WorldUp).Normalize()
	cam.Up = cam.Right.Cross(cam.Front).Normalize()
}
