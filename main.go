package main

//fuck you world. You're just a peice of shit
import (
	"fmt"
	"log"
	"math"
	"runtime"
	"syscall"
	"time"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

const (
	SM_CXSCREEN = uintptr(0) // X Size of screen
	SM_CYSCREEN = uintptr(1) // Y Size of screen
)

type Matrix [4][4]float32

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func (v Vector3) Dist(other Vector3) float32 {
	return float32(math.Abs(float64(v.X-other.X)) + math.Abs(float64(v.Y-other.Y)) + math.Abs(float64(v.Z-other.Z)))
}

type Vector2 struct {
	X float32
	Y float32
}

type Rectangle struct {
	Top    float32
	Left   float32
	Right  float32
	Bottom float32
}

type Entity struct {
	Health   int32
	Team     int32
	Name     string
	Position Vector2
	Bones    map[string]Vector2
	HeadPos  Vector3
	Distance float32
	Rect     Rectangle
}

type Offset struct {
	DwViewMatrix           uintptr `json:"dwViewMatrix"`
	DwLocalPlayerPawn      uintptr `json:"dwLocalPlayerPawn"`
	DwEntityList           uintptr `json:"dwEntityList"`
	M_hPlayerPawn          uintptr `json:"m_hPlayerPawn"`
	M_iHealth              uintptr `json:"m_iHealth"`
	M_lifeState            uintptr `json:"m_lifeState"`
	M_iTeamNum             uintptr `json:"m_iTeamNum"`
	M_vOldOrigin           uintptr `json:"m_vOldOrigin"`
	M_pGameSceneNode       uintptr `json:"m_pGameSceneNode"`
	M_modelState           uintptr `json:"m_modelState"`
	M_boneArray            uintptr `json:"m_boneArray"`
	M_nodeToWorld          uintptr `json:"m_nodeToWorld"`
	M_sSanitizedPlayerName uintptr `json:"m_sSanitizedPlayerName"`
}

var (
	user32                     = windows.NewLazySystemDLL("user32.dll")
	gdi32                      = windows.NewLazySystemDLL("gdi32.dll")
	getSystemMetrics           = user32.NewProc("GetSystemMetrics")
	setLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	showCursor                 = user32.NewProc("ShowCursor")
	setTextAlign               = gdi32.NewProc("SetTextAlign")
	createFont                 = gdi32.NewProc("CreateFontW")
	createCompatibleDC         = gdi32.NewProc("CreateCompatibleDC")
	createSolidBrush           = gdi32.NewProc("CreateSolidBrush")
	createPen                  = gdi32.NewProc("CreatePen")
	procCreatePen              = modGdi32.NewProc("CreatePen")
	modGdi32                   = syscall.NewLazyDLL("gdi32.dll")
)

// changed up
var (
	teamCheck           bool   = true
	headCircle          bool   = true
	skeletonRendering   bool   = true
	boxRendering        uint8  = 3 //change bool to uint8
	nameRendering       bool   = false
	healthBarRendering  bool   = true
	healthTextRendering bool   = false
	line                bool   = true
	frameDelay          uint32 = 0 //change 15 to 0
)

func init() {
	// Ensure main() runs on the main thread.
	runtime.LockOSThread()
}

func logAndSleep(message string, err error) {
	log.Printf("%s: %v\n", message, err)
	time.Sleep(5 * time.Second)
}

//   ________                            _ \n
//  |   _____|                          | | \n
//  |  |_____ _____  _____   ____  ____ | | \n
//  |   _____/ ___/ |  __ \ / ___ / __ \| | \n
//  |  |_____\____ \|  ___/| (___| ____|| |__ \n
//  |________|_____/| |     \____ \_____|___/ \n
//                  |_| \n

//   _________                            _ \n |   ______|                          | | \n |  |______ _____  _____   ____  ____ | | \n |   ______/ ___/ |  __ \ / ___ / __ \| | \n |  |______\____ \|  ___/| (___| ____|| |__ \n |_________|_____/| |     \____ \_____|___/ \n                  |_| \n

// fuck yyou fuck me fuck evry
func main() {
	go cliMenu()

	screenWidth, _, _ := getSystemMetrics.Call(0)
	screenHeight, _, _ := getSystemMetrics.Call(1)

	hwnd := initWindow(screenWidth, screenHeight)
	if hwnd == 0 {
		logAndSleep("Error creating window", fmt.Errorf("%v", win.GetLastError()))
		return
	}
	defer win.DestroyWindow(hwnd)

	// win.SetCursor()

	pid, err := findProcessId("cs2.exe")
	if err != nil {
		logAndSleep("Error finding process ID", err)
		return
	}

	clientDll, err := getModuleBaseAddress(pid, "client.dll")
	if err != nil {
		logAndSleep("Error getting client.dll base address", err)
		return
	}

	procHandle, err := getProcessHandle(pid)
	if err != nil {
		logAndSleep("Error getting process handle", err)
		return
	}

	hdc := win.GetDC(hwnd)
	if hdc == 0 {
		logAndSleep("Error getting device context", fmt.Errorf("%v", win.GetLastError()))
		return
	}

	bgBrush, _, _ := createSolidBrush.Call(uintptr(0x000000))
	if bgBrush == 0 {
		logAndSleep("Error creating brush", fmt.Errorf("%v", win.GetLastError()))
		return
	}
	defer win.DeleteObject(win.HGDIOBJ(bgBrush))
	redPen, _, _ := createPen.Call(win.PS_SOLID, 1, 0xFF0000)
	if redPen == 0 {
		logAndSleep("Error creating pen", fmt.Errorf("%v", win.GetLastError()))
		return
	}
	//make red more deep
	defer win.DeleteObject(win.HGDIOBJ(redPen))
	greenPen, _, _ := createPen.Call(win.PS_SOLID, 1, 0x00FF00)
	if greenPen == 0 {
		logAndSleep("Error creating pen", fmt.Errorf("%v", win.GetLastError()))
		return
	}
	defer win.DeleteObject(win.HGDIOBJ(greenPen))
	//changed bluepen color more deep
	bluePen, _, _ := createPen.Call(win.PS_SOLID, 1, 0x0000ff)
	if bluePen == 0 {
		logAndSleep("Error creating pen", fmt.Errorf("%v", win.GetLastError()))
		return
	}
	defer win.DeleteObject(win.HGDIOBJ(bluePen))
	//changed bone color to cyan
	bonePen, _, _ := createPen.Call(win.PS_SOLID, 1, 0x00FFFF)
	if bonePen == 0 {
		logAndSleep("Error creating pen", fmt.Errorf("%v", win.GetLastError()))
		return
	}
	defer win.DeleteObject(win.HGDIOBJ(bonePen))
	outlinePen, _, _ := createPen.Call(win.PS_SOLID, 1, 0x000001)
	if outlinePen == 0 {
		logAndSleep("Error creating pen", fmt.Errorf("%v", win.GetLastError()))
		return
	}
	defer win.DeleteObject(win.HGDIOBJ(outlinePen))

	font, _, _ := createFont.Call(12, 0, 0, 0, win.FW_HEAVY, 0, 0, 0, win.DEFAULT_CHARSET, win.OUT_DEFAULT_PRECIS, win.CLIP_DEFAULT_PRECIS, win.DEFAULT_QUALITY, win.DEFAULT_PITCH|win.FF_DONTCARE, 0)

	offsets := getOffsets()

	var msg win.MSG

	for win.GetMessage(&msg, 0, 0, 0) > 0 {
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)

		win.SetTimer(hwnd, 1, frameDelay, 0)

		memhdc, _, _ := createCompatibleDC.Call(uintptr(hdc))
		memBitmap := win.CreateCompatibleBitmap(hdc, int32(screenWidth), int32(screenHeight))
		win.SelectObject(win.HDC(memhdc), win.HGDIOBJ(memBitmap))
		win.SelectObject(win.HDC(memhdc), win.HGDIOBJ(bgBrush))
		win.SetBkMode(win.HDC(memhdc), win.TRANSPARENT)
		win.SelectObject(win.HDC(memhdc), win.HGDIOBJ(font))

		entities := getEntitiesInfo(procHandle, clientDll, screenWidth, screenHeight, offsets)
		//bone (skeleton) 颜色随着队伍变化

		w, _, _ := syscall.NewLazyDLL(`User32.dll`).NewProc(`GetSystemMetrics`).Call(SM_CXSCREEN)

		for _, entity := range entities {
			if entity.Distance < 35 {
				continue
			}
			if skeletonRendering {
				if entity.Team == 2 {
					drawSkeleton(win.HDC(memhdc), bluePen, entity.Bones)
				} else {
					drawSkeleton(win.HDC(memhdc), redPen, entity.Bones)
				}
			}
			if entity.Team == 2 {
				renderEntityInfo(win.HDC(memhdc), bluePen, greenPen, outlinePen, bonePen, entity.Rect, entity.Health, entity.Name, entity.HeadPos, int32(w))
			} else {
				renderEntityInfo(win.HDC(memhdc), redPen, greenPen, outlinePen, bonePen, entity.Rect, entity.Health, entity.Name, entity.HeadPos, int32(w))
			}
		}
		win.BitBlt(hdc, 0, 0, int32(screenWidth), int32(screenHeight), win.HDC(memhdc), 0, 0, win.SRCCOPY)

		// Delete the memory bitmap and device context
		win.DeleteObject(win.HGDIOBJ(memBitmap))
		win.DeleteDC(win.HDC(memhdc))
	}
}
