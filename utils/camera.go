package utils

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
	"github.com/go-gl/mathgl/mgl64"
)

const (
	FORWARD = iota
	BACKWARD
	LEFT
	RIGHT
)

const (
	YAW        = -90.0
	PITCH      = 0.0
	SPEED      = 3.0
	SENSITIVTY = 0.25
	ZOOM       = 45.0
)

// Camera is the camera object maintaing the stae
type Camera struct {
	Position mgl32.Vec3
	Front    mgl32.Vec3
	Up       mgl32.Vec3
	Right    mgl32.Vec3
	WorldUp  mgl32.Vec3

	// Eular Angles
	Yaw   float64
	Pitch float64

	// Camera options
	MovementSpeed    float64
	MouseSensitivity float64
	Zoom             float64
}

func NewCamera(position, up mgl32.Vec3, yaw, pitch float64) Camera {
	c := Camera{
		Position:         position,
		WorldUp:          up,
		Yaw:              yaw,
		Pitch:            pitch,
		Front:            mgl32.Vec3{0.0, 0.0, -1.0},
		MovementSpeed:    SPEED,
		MouseSensitivity: SENSITIVTY,
		Zoom:             ZOOM,
	}
	c.updateCameraVectors()
	return c
}

// Constructor with scalar values
func NewCameraWithScalars(posX, posY, posZ, upX, upY, upZ float32,  yaw, pitch float64) Camera {
	c := Camera{
		Position:         mgl32.Vec3{posX, posY, posZ},
		WorldUp:          mgl32.Vec3{upX, upY, upZ},
		Yaw:              yaw,
		Pitch:            pitch,
		Front:            mgl32.Vec3{0.0, 0.0, -1.0},
		MovementSpeed:    SPEED,
		MouseSensitivity: SENSITIVTY,
		Zoom:             ZOOM,
	}
	c.updateCameraVectors()
	return c
}

// GetViewMatrix returns the view natrix
func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	eye := c.Position
	center := c.Position.Add(c.Front)
	up := c.Up
	return mgl32.LookAt(
		eye.X(), eye.Y(), eye.Z(),
		center.X(), center.Y(), center.Z(),
		up.X(), up.Y(), up.Z())
}

// Processes input received from any keyboard-like input system. Accepts input parameter in the form of camera defined ENUM (to abstract it from windowing systems)
func (c *Camera) ProcessKeyboard(direction int, deltaTime float64) {
	velocity := float32(c.MovementSpeed * deltaTime)
	if direction == FORWARD {
		c.Position = c.Position.Add(c.Front.Mul(velocity))
	}
	if direction == BACKWARD {
		c.Position = c.Position.Sub(c.Front.Mul(velocity))
	}
	if direction == LEFT {
		c.Position = c.Position.Sub(c.Right.Mul(velocity))
	}
	if direction == RIGHT {
		c.Position = c.Position.Add(c.Right.Mul(velocity))
	}
}

// Processes input received from a mouse input system. Expects the offset value in both the x and y direction.
func (c *Camera) ProcessMouseMovement(xoffset, yoffset float64, constrainPitch bool) {
	xoffset *= c.MouseSensitivity
	yoffset *= c.MouseSensitivity

	c.Yaw += xoffset
	c.Pitch += yoffset

	// Make sure that when pitch is out of bounds, screen doesn't get flipped
	if constrainPitch {
		if c.Pitch > 89.0 {
			c.Pitch = 89.0
		}
		if c.Pitch < -89.0 {
			c.Pitch = -89.0
		}
	}
	// Update Front, Right and Up Vectors using the updated Eular angles
	c.updateCameraVectors()
}

// Processes input received from a mouse scroll-wheel event. Only requires input on the vertical wheel-axis
func (c *Camera) ProcessMouseScroll(yoffset float64) {
	if c.Zoom >= 1.0 && c.Zoom <= 45.0 {
		c.Zoom -= yoffset
	}
	if c.Zoom <= 1.0 {
		c.Zoom = 1.0
	}
	if c.Zoom >= 45.0 {
		c.Zoom = 45.0
	}
}

func (c *Camera) updateCameraVectors() {

	x := float32(math.Cos(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch)))
	y := float32(math.Sin(mgl64.DegToRad(c.Pitch)))
	z := float32(math.Sin(mgl64.DegToRad(c.Yaw)) * math.Cos(mgl64.DegToRad(c.Pitch)))
	front := mgl32.Vec3{x, y, z}
	front = front.Normalize()
	// Also re-calculate the Right and Up vector
	// Normalize the vectors, because their length gets closer to 0 the more you look up or down which results in slower movement.
	c.Right = front.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Front).Normalize()
}

