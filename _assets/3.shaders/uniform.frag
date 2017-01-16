
#version 330 core
out vec4 color;

uniform vec4 ourColor; // We set this variable in the OpenGL code.

void main()
{
    color = ourColor;
}