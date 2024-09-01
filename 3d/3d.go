package three_d

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

func Run_3D_Terrain() {
	// 앱 초기 설정
	a := app.App()
	a.IWindow.(*window.GlfwWindow).SetTitle("3D Terrain")
	a.IWindow.(*window.GlfwWindow).SetSize(500, 500)
	scene := core.NewNode()

	// 카메라 설정
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)
	camera.NewOrbitControl(cam)

	// 오브젝트
	geom := geometry.NewTorus(1, .4, 12, 32, math32.Pi*2)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)
	scene.Add(mesh)

	// 빛 설정
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	//scene.Add(helper.NewAxes(0.5)) //Axis

	// 배경 색 설정
	a.Gls().ClearColor(0, 0, 0, 1.0)

	// 실행
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		renderer.Render(scene, cam)
	})
}
