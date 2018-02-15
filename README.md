### The original tutorial
https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl



Hi, this is my attempt to create render using Go. 
I've never written any program with Go, so, let's get it!

### How to run

You need graphic adapter **supporting OpenGL ~4'th** version, **otherwise**, you have to shange lib's and code on **OpenGL 2** by yourself.

Make sure you have Go on your machine. 

Lets get required libs 
```
go get github.com/go-gl/gl/v4.1-core/gl
go get github.com/go-gl/glfw/v3.2/glfw

```

Now, go into project folder and run the following:

```
go run simple_triangle.go
```
Or just run _run.bat_

You did it!


