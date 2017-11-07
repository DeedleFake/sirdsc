package sdl

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

// #cgo pkg-config: sdl2
//
// #include <SDL.h>
import "C"

const (
	WINDOWPOS_UNDEFINED = C.SDL_WINDOWPOS_UNDEFINED
	WINDOWPOS_CENTERED  = C.SDL_WINDOWPOS_CENTERED

	K_UP    = C.SDLK_UP
	K_DOWN  = C.SDLK_DOWN
	K_LEFT  = C.SDLK_LEFT
	K_RIGHT = C.SDLK_RIGHT
	K_w     = C.SDLK_w
	K_s     = C.SDLK_s

	PIXELFORMAT_ARGB8888 = C.SDL_PIXELFORMAT_ARGB8888
	PIXELFORMAT_RGBA32   = C.SDL_PIXELFORMAT_RGBA32
	PIXELFORMAT_ABGR8888 = C.SDL_PIXELFORMAT_ABGR8888

	TEXTUREACCESS_STREAMING = C.SDL_TEXTUREACCESS_STREAMING
)

func Init() error {
	if C.SDL_Init(C.SDL_INIT_EVERYTHING) < 0 {
		return errors.New(C.GoString(C.SDL_GetError()))
	}
	return nil
}

func Quit() {
	C.SDL_Quit()
}

type Window struct {
	c *C.SDL_Window
}

func CreateWindow(title string, x, y, w, h int, flags uint32) (*Window, error) {
	c := C.SDL_CreateWindow(
		C.CString(title),
		C.int(x),
		C.int(y),
		C.int(w),
		C.int(h),
		C.Uint32(flags),
	)
	if c == nil {
		return nil, errors.New(C.GoString(C.SDL_GetError()))
	}

	win := &Window{
		c: c,
	}
	runtime.SetFinalizer(win, (*Window).Destroy)
	return win, nil
}

func CreateWindowAndRenderer(w, h int, flags uint32) (*Window, *Renderer, error) {
	var win *C.SDL_Window
	var ren *C.SDL_Renderer
	if C.SDL_CreateWindowAndRenderer(C.int(w), C.int(h), C.Uint32(flags), &win, &ren) < 0 {
		return nil, nil, errors.New(C.GoString(C.SDL_GetError()))
	}

	return &Window{c: win}, &Renderer{
		c: ren,

		pix: make([]color.Color, w*h),
		w:   w,
		h:   h,
	}, nil
}

func (win *Window) Destroy() {
	C.SDL_DestroyWindow(win.c)
}

func (win *Window) UpdateSurface() {
	C.SDL_UpdateWindowSurface(win.c)
}

type Renderer struct {
	c *C.SDL_Renderer

	m    sync.RWMutex
	pix  []color.Color
	w, h int
}

func (ren *Renderer) Destroy() {
	C.SDL_DestroyRenderer(ren.c)
}

func (ren *Renderer) Copy(tex *Texture, src image.Rectangle, dst image.Rectangle) {
	C.SDL_RenderCopy(ren.c, tex.c, sdlRect(src), sdlRect(dst))
}

func (ren *Renderer) ColorModel() color.Model {
	return color.RGBAModel
}

func (ren *Renderer) Bounds() image.Rectangle {
	return image.Rect(0, 0, ren.w, ren.h)
}

func (ren *Renderer) At(x, y int) color.Color {
	ren.m.RLock()
	defer ren.m.RUnlock()

	return ren.pix[(y*ren.w)+x]
}

func (ren *Renderer) Set(x, y int, c color.Color) {
	ren.m.Lock()
	defer ren.m.Unlock()

	r, g, b, a := c.RGBA()
	C.SDL_SetRenderDrawColor(ren.c, C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a))
	C.SDL_RenderDrawPoint(ren.c, C.int(x), C.int(y))
	ren.pix[(y*ren.w)+x] = c
}

func (ren *Renderer) Present() {
	C.SDL_RenderPresent(ren.c)
}

type Texture struct {
	c *C.SDL_Texture

	format *C.SDL_PixelFormat
	w, h   int
}

func (ren *Renderer) CreateTexture(format uint32, access, w, h int) (*Texture, error) {
	t := C.SDL_CreateTexture(ren.c, C.Uint32(format), C.int(access), C.int(w), C.int(h))
	if t == nil {
		return nil, errors.New(C.GoString(C.SDL_GetError()))
	}

	return &Texture{
		c: t,
		w: w,
		h: h,
	}, nil
}

func (t *Texture) Destroy() {
	C.SDL_DestroyTexture(t.c)
}

type TextureImage struct {
	t *Texture

	w, h int
	pix  []C.Uint32
}

func (t *Texture) Image() *TextureImage {
	var pixels unsafe.Pointer
	var pitch C.int
	C.SDL_LockTexture(t.c, nil, &pixels, &pitch)

	return &TextureImage{
		t: t,
		w: t.w,
		h: t.h,
		pix: *(*[]C.Uint32)(unsafe.Pointer(&reflect.SliceHeader{
			Data: uintptr(pixels),
			Len:  t.w * t.h,
			Cap:  t.w * t.h,
		})),
	}
}

func (img *TextureImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (img *TextureImage) Bounds() image.Rectangle {
	return image.Rect(0, 0, img.w, img.h)
}

func (img *TextureImage) At(x, y int) color.Color {
	c := img.pix[(y*img.w)+x]
	a := (c & 0xFF000000) >> 24
	b := (c & 0xFF0000) >> 16
	g := (c & 0xFF00) >> 8
	r := c & 0xFF

	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

func (img *TextureImage) Set(x, y int, c color.Color) {
	r, g, b, a := c.RGBA()
	cc := uint32((a*255/0xFFFF)<<24) | uint32((b*255/0xFFFF)<<16) | uint32((g*255/0xFFFF)<<8) | uint32(r*255/0xFFFF)

	img.pix[(y*img.w)+x] = C.Uint32(cc)
}

func (img *TextureImage) Close() {
	C.SDL_UnlockTexture(img.t.c)
}

type Surface struct {
	c *C.SDL_Surface
}

func (win *Window) GetSurface() (*Surface, error) {
	screen := C.SDL_GetWindowSurface(win.c)
	if screen == nil {
		return nil, errors.New(C.GoString(C.SDL_GetError()))
	}

	return &Surface{
		c: screen,
	}, nil
}

func (s *Surface) ColorModel() color.Model {
	return color.RGBAModel
}

func (s *Surface) Bounds() image.Rectangle {
	return image.Rect(0, 0, int(s.c.w), int(s.c.h))
}

func (s *Surface) At(x, y int) color.Color {
	C.SDL_LockSurface(s.c)
	defer C.SDL_UnlockSurface(s.c)

	c := s.pix()[(y*int(s.c.w))+x]

	var r, g, b C.Uint8
	C.SDL_GetRGB(C.Uint32(c), s.c.format, &r, &g, &b)
	return color.RGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: 255,
	}
}

func (s *Surface) Set(x, y int, c color.Color) {
	r, g, b, _ := c.RGBA()
	cc := C.SDL_MapRGB(s.c.format, C.Uint8(r), C.Uint8(g), C.Uint8(b))

	C.SDL_LockSurface(s.c)
	defer C.SDL_UnlockSurface(s.c)

	s.pix()[(y*int(s.c.w))+x] = uint32(cc)
}

func (s *Surface) pix() []uint32 {
	return *(*[]uint32)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(s.c.pixels)),
		Len:  int(s.c.w * s.c.h),
		Cap:  int(s.c.w * s.c.h),
	}))
}

func sdlRect(r image.Rectangle) (cr *C.SDL_Rect) {
	if r == image.ZR {
		return nil
	}

	return &C.SDL_Rect{
		x: C.int(r.Min.X),
		y: C.int(r.Min.Y),
		w: C.int(r.Dx()),
		h: C.int(r.Dy()),
	}
}

func PollEvent() interface{} {
	var ev C.SDL_Event
	ok := C.SDL_PollEvent(&ev)
	if ok == 0 {
		return nil
	}

	switch t := *(*uint32)(unsafe.Pointer(&ev)); t {
	case C.SDL_QUIT:
		return QuitEvent{c: *(*C.SDL_QuitEvent)(unsafe.Pointer(&ev))}

	case C.SDL_KEYUP:
		return KeyUpEvent{c: *(*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev))}
	case C.SDL_KEYDOWN:
		return KeyDownEvent{c: *(*C.SDL_KeyboardEvent)(unsafe.Pointer(&ev))}

	default:
		return fmt.Errorf("Unsupported event type: %v", t)
	}
}

type QuitEvent struct {
	c C.SDL_QuitEvent
}

type KeyUpEvent struct {
	c C.SDL_KeyboardEvent
}

func (ev KeyUpEvent) Keysym() Keysym {
	return Keysym{c: ev.c.keysym}
}

type KeyDownEvent struct {
	c C.SDL_KeyboardEvent
}

func (ev KeyDownEvent) Keysym() Keysym {
	return Keysym{c: ev.c.keysym}
}

type Keysym struct {
	c C.SDL_Keysym
}

func (ks Keysym) Sym() uint32 {
	return uint32(ks.c.sym)
}
