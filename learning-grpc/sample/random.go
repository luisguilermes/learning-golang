package sample

import (
	"math/rand"

	"github.com/google/uuid"
	"github.com/luisguilermes/learning-golang/learning-grpc/pb"
)

// randomKeyboardLayout returns a random keyboard layout
func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

// randomCPUBrand returns a random CPU brand
func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD")
}

// randomCPUName returns a random CPU name
func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFromSet("Xeon E", "Core i9", "Core i7", "Core i5", "Core i3")
	}
	return randomStringFromSet("Ryzen 9", "Ryzen 7", "Ryzen 5", "Ryzen 3")
}

// randomGPUBrand returns a random GPU brand
func randomGPUBrand() string {
	return randomStringFromSet("Nvidia", "AMD")
}

// randomGPUName returns a random GPU name
func randomGPUName(brand string) string {
	if brand == "Nvidia" {
		return randomStringFromSet("RTX 3090", "RTX 3080", "RTX 3070", "RTX 3060")
	}
	return randomStringFromSet("RX 6900 XT", "RX 6800 XT", "RX 6700 XT")
}

func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4320)
	width := height * 16 / 9

	return &pb.Screen_Resolution{
		Width:  uint32(height),
		Height: uint32(width),
	}
}

// randomScreenPanel returns a random screen panel
func randomScreenPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Lenovo")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("MacBook Air", "MacBook Pro")
	case "Dell":
		return randomStringFromSet("XPS", "Latitude", "Inspiron")
	default:
		return randomStringFromSet("ThinkPad", "IdeaPad", "Yoga")
	}
}

// randomStringFromSet returns a random string from a set of strings
func randomStringFromSet(strings ...string) string {
	n := len(strings)
	if n == 0 {
		return ""
	}
	return strings[rand.Intn(n)]
}

// randomBool returns a random boolean
func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat64(min float64, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func randomFloat32(min float32, max float32) float32 {
	return min + rand.Float32()*(max-min)
}

func randomID() string {
	return uuid.New().String()
}
