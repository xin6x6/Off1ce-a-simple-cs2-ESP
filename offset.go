package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

func getOffsets() Offset {
	var offsets Offset

	// Open the file
	offsetsJson, err := os.Open("offsets.json")
	if err != nil {
		fmt.Println("Error opening offsets.json", err)
		return offsets
	}
	defer offsetsJson.Close()

	// Decode the JSON
	err = json.NewDecoder(offsetsJson).Decode(&offsets)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return offsets
	}
	return offsets
}

func getEntitiesInfo(procHandle windows.Handle, clientDll uintptr, screenWidth uintptr, screenHeight uintptr, offsets Offset) []Entity {
	var entityList uintptr
	var entities []Entity
	err := read(procHandle, clientDll+offsets.DwEntityList, &entityList)
	if err != nil {
		return entities
	}
	var (
		localPlayerP           uintptr
		localPlayerGameScene   uintptr
		localPlayerSceneOrigin Vector3
		localTeam              int32
		listEntry              uintptr
		gameScene              uintptr
		entityController       uintptr
		entityControllerPawn   uintptr
		entityPawn             uintptr
		entityNameAddress      uintptr
		entityBoneArray        uintptr
		entityTeam             int32
		entityHealth           int32
		entityLifeState        int32
		entityName             string
		sanitizedNameStr       string
		entityOrigin           Vector3
		viewMatrix             Matrix
	)
	bones := map[string]int{
		"head":        6,
		"neck_0":      5,
		"spine_1":     4,
		"spine_2":     2,
		"pelvis":      0,
		"arm_upper_L": 8,
		"arm_lower_L": 9,
		"hand_L":      10,
		"arm_upper_R": 13,
		"arm_lower_R": 14,
		"hand_R":      15,
		"leg_upper_L": 22,
		"leg_lower_L": 23,
		"ankle_L":     24,
		"leg_upper_R": 25,
		"leg_lower_R": 26,
		"ankle_R":     27,
	}
	var (
		currentBone      Vector3
		entityHead       Vector3
		entityHeadTop    Vector3
		entityHeadBottom Vector3
	)
	// localPlayerP
	err = read(procHandle, clientDll+offsets.DwLocalPlayerPawn, &localPlayerP)
	if err != nil {
		return entities
	}
	// localPlayerGameScene
	err = read(procHandle, localPlayerP+offsets.M_pGameSceneNode, &localPlayerGameScene)
	if err != nil {
		return entities
	}
	// localPlayerSceneOrigin
	err = read(procHandle, localPlayerGameScene+offsets.M_nodeToWorld, &localPlayerSceneOrigin)
	if err != nil {
		return entities
	}
	// viewMatrix
	err = read(procHandle, clientDll+offsets.DwViewMatrix, &viewMatrix)
	if err != nil {
		return entities
	}
	for i := 0; i < 64; i++ {
		var tempEntity Entity
		var entityBones map[string]Vector2 = make(map[string]Vector2)
		var sanitizedName strings.Builder
		// listEntry
		err = read(procHandle, entityList+uintptr((8*(i&0x7FFF)>>9)+16), &listEntry)
		if err != nil {
			return entities
		}
		if listEntry == 0 {
			continue
		}
		// entityController
		err = read(procHandle, listEntry+uintptr(120)*uintptr(i&0x1FF), &entityController)
		if err != nil {
			return entities
		}
		if entityController == 0 {
			continue
		}
		// entityControllerPawn
		err = read(procHandle, entityController+offsets.M_hPlayerPawn, &entityControllerPawn)
		if err != nil {
			return entities
		}
		if entityControllerPawn == 0 {
			continue
		}
		// listEntry
		err = read(procHandle, entityList+uintptr(0x8*((entityControllerPawn&0x7FFF)>>9)+16), &listEntry)
		if err != nil {
			return entities
		}
		if listEntry == 0 {
			continue
		}
		// entityPawn
		err = read(procHandle, listEntry+uintptr(120)*uintptr(entityControllerPawn&0x1FF), &entityPawn)
		if err != nil {
			return entities
		}
		if entityPawn == 0 {
			continue
		}
		if entityPawn == localPlayerP {
			continue
		}
		// entityLifeState
		err = read(procHandle, entityPawn+offsets.M_lifeState, &entityLifeState)
		if err != nil {
			return entities
		}
		if entityLifeState != 256 {
			continue
		}
		// entityTeam
		err = read(procHandle, entityPawn+offsets.M_iTeamNum, &entityTeam)
		if err != nil {
			return entities
		}
		if entityTeam == 0 {
			continue
		}
		if teamCheck {
			// localTeam
			err = read(procHandle, localPlayerP+offsets.M_iTeamNum, &localTeam)
			if err != nil {
				return entities
			}
			if localTeam == entityTeam {
				continue
			}
		}
		// entityHealth
		err = read(procHandle, entityPawn+offsets.M_iHealth, &entityHealth)
		if err != nil {
			return entities
		}
		if entityHealth < 1 || entityHealth > 100 {
			continue
		}
		// entityNameAddress
		err = read(procHandle, entityController+offsets.M_sSanitizedPlayerName, &entityNameAddress)
		if err != nil {
			return entities
		}
		// entityName
		err = read(procHandle, entityNameAddress, &entityName)
		if err != nil {
			return entities
		}
		if entityName == "" {
			continue
		}
		for _, c := range entityName {
			if unicode.IsLetter(c) || unicode.IsDigit(c) || unicode.IsPunct(c) || unicode.IsSpace(c) {
				sanitizedName.WriteRune(c)
			}
		}
		sanitizedNameStr = sanitizedName.String()
		// gameScene
		err = read(procHandle, entityPawn+offsets.M_pGameSceneNode, &gameScene)
		if err != nil {
			return entities
		}
		if gameScene == 0 {
			continue
		}
		// entityBoneArray
		err = read(procHandle, gameScene+offsets.M_modelState+offsets.M_boneArray, &entityBoneArray)
		if err != nil {
			return entities
		}
		if entityBoneArray == 0 {
			continue
		}
		// entityOrigin
		err = read(procHandle, entityPawn+offsets.M_vOldOrigin, &entityOrigin)
		if err != nil {
			return entities
		}
		// boneArray
		for boneName, boneIndex := range bones {
			err = read(procHandle, entityBoneArray+uintptr(boneIndex)*32, &currentBone)
			if err != nil {
				return entities
			}
			if boneName == "head" {
				entityHead = currentBone
				if !skeletonRendering {
					break
				}
			}
			boneX, boneY := worldToScreen(viewMatrix, currentBone)
			entityBones[boneName] = Vector2{boneX, boneY}
		}
		entityHeadTop = Vector3{entityHead.X, entityHead.Y, entityHead.Z + 7}
		entityHeadBottom = Vector3{entityHead.X, entityHead.Y, entityHead.Z - 5}
		screenPosHeadX, screenPosHeadTopY := worldToScreen(viewMatrix, entityHeadTop)
		_, screenPosHeadBottomY := worldToScreen(viewMatrix, entityHeadBottom)
		screenPosFeetX, screenPosFeetY := worldToScreen(viewMatrix, entityOrigin)
		entityBoxTop := Vector3{entityOrigin.X, entityOrigin.Y, entityOrigin.Z + 70}
		_, screenPosBoxTop := worldToScreen(viewMatrix, entityBoxTop)
		if screenPosHeadX <= -1 || screenPosFeetY <= -1 || screenPosHeadX >= float32(screenWidth) || screenPosHeadTopY >= float32(screenHeight) {
			continue
		}
		boxHeight := screenPosFeetY - screenPosBoxTop

		tempEntity.Health = entityHealth
		tempEntity.Team = entityTeam
		tempEntity.Name = sanitizedNameStr
		tempEntity.Distance = entityOrigin.Dist(localPlayerSceneOrigin)
		tempEntity.Position = Vector2{screenPosFeetX, screenPosFeetY}
		tempEntity.Bones = entityBones
		tempEntity.HeadPos = Vector3{screenPosHeadX, screenPosHeadTopY, screenPosHeadBottomY}
		tempEntity.Rect = Rectangle{screenPosBoxTop, screenPosFeetX - boxHeight/4, screenPosFeetX + boxHeight/4, screenPosFeetY}

		entities = append(entities, tempEntity)
	}
	return entities
}

func worldToScreen(viewMatrix Matrix, position Vector3) (float32, float32) {
	var screenX float32
	var screenY float32
	screenX = viewMatrix[0][0]*position.X + viewMatrix[0][1]*position.Y + viewMatrix[0][2]*position.Z + viewMatrix[0][3]
	screenY = viewMatrix[1][0]*position.X + viewMatrix[1][1]*position.Y + viewMatrix[1][2]*position.Z + viewMatrix[1][3]
	w := viewMatrix[3][0]*position.X + viewMatrix[3][1]*position.Y + viewMatrix[3][2]*position.Z + viewMatrix[3][3]
	if w < 0.01 {
		return -1, -1
	}
	invw := 1.0 / w
	screenX *= invw
	screenY *= invw
	width, _, _ := getSystemMetrics.Call(0)
	height, _, _ := getSystemMetrics.Call(1)
	widthFloat := float32(width)
	heightFloat := float32(height)
	x := widthFloat / 2
	y := heightFloat / 2
	x += 0.5*screenX*widthFloat + 0.5
	y -= 0.5*screenY*heightFloat + 0.5
	return x, y
}

// here's change

func CreatePen(style, width int32, color uint32) uintptr {
	ret, _, _ := procCreatePen.Call(uintptr(style), uintptr(width), uintptr(color))
	return ret
}

func RGB(r, g, b byte) uint32 {
	return uint32(r) | (uint32(g) << 8) | (uint32(b) << 16)
}

func windowProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_TIMER:
		return 0
	case win.WM_DESTROY:
		win.PostQuitMessage(0)
		return 0
	default:
		return win.DefWindowProc(hwnd, msg, wParam, lParam)
	}
}
