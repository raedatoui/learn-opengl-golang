#version 330 core
// in vec3 ourColor;
in vec3 ourPosition;

out vec4 color;

void main()
{
    color = vec4(ourPosition, 1.0f);    // Note how the position value is linearly interpolated to get all the different colors
}