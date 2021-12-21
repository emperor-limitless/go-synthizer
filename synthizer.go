//MIT License
//
//Copyright (c) 2021 mohamedSulaimanAlmarzooqi
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.
//
package synthizer

/*
#include <synthizer.h>
#include <synthizer_constants.h>
#include <stdlib.h>
#include <stdbool.h>
bool bld(int value) {
	return value == 0;
}
struct syz_DeleteBehaviorConfig create_DeleteBehaviorConfig(int linger, double linger_timeout) {
	struct syz_DeleteBehaviorConfig cfg;
	cfg.linger = linger;
	cfg.linger_timeout = linger_timeout;
	return cfg;
}
struct syz_LibraryConfig create_LibraryConfig(unsigned int log_level, unsigned int logging_backend, const char *libsndfile_path) {
	struct syz_LibraryConfig config;
	config.log_level = log_level;
	config.logging_backend = logging_backend;
	config.libsndfile_path = libsndfile_path;
	return config;
}
syz_Handle create_handle() {
	syz_Handle handle = 0;
	return handle;
}
*/
import "C"
import (
	"unsafe"
	"errors"
	"strconv"
	"reflect"
	"runtime"
)

// The synthizer main object struct.
type ObjectBase struct {
	handle *C.syz_Handle
}
func _NewObjectBase(hdl *C.syz_Handle) ObjectBase {
	return ObjectBase { handle: hdl }
}

func (self *ObjectBase) GetHandleChecked() (error, *C.syz_Handle) {
	// Any attempt to reference a non-existing object should raise an error
	handle := self.GetHandle()
	if *handle == 0 {
		return errors.New("SynthizerError: Object no longer exist"), nil
	}
	return nil, handle
}

func (self *ObjectBase) Destroy() error {
	*self.handle = 0
	return CHECKED(C.syz_handleDecRef(*self.handle))
}

func (self *ObjectBase) GetHandle() *C.syz_Handle {
	return self.handle
}
func (self *ObjectBase) ConfigDeleteBehavior(linger int, timeout float32) error {
	cfg := C.create_DeleteBehaviorConfig(C.int(linger), C.double(timeout))
	C.syz_initDeleteBehaviorConfig(&cfg)
	return CHECKED(C.syz_configDeleteBehavior(*self.handle, &cfg))
}

type PropertyBase struct {
	instance *ObjectBase
	property C.int
}
func (self *PropertyBase) GetHandleChecked() (error, *C.syz_Handle) {
	// Any attempt to reference a non-existing object should raise an error
	obj := reflect.ValueOf(*self.instance)
	if obj.IsZero() {
		return errors.New("SynthizerError: Object no longer exist"), nil
	}
	handle := (*self.instance).GetHandle()
	if *handle == 0 {
		return errors.New("SynthizerError: Object no longer exist"), nil
	}
	return nil, handle
}
func (self *PropertyBase) GetProperty() C.int {
	return self.property
}

type IntProperty struct {
	PropertyBase
}

func NewIntProperty(instance *ObjectBase, property C.int) IntProperty {
	return IntProperty { PropertyBase { instance, property } }
}

func (self *IntProperty) Get() (error, int) {
	var val C.int = 0
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err, 0
	}
	err = CHECKED(C.syz_getI(&val, *handle, self.property))
	if err != nil {
		return err, 0
	}
	return nil, int(val)
}

func (self *IntProperty) Set(value int) error {
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_setI(*handle, self.property, C.int(value)))
	if err != nil {
		return err
	}
	return nil
}


// A boolean property, For stuff like looping so we don't need to let the user do something like object.looping.Set(1) OR something like that.
type BoolProperty struct {
	PropertyBase
}

func NewBoolProperty(instance *ObjectBase, property C.int) BoolProperty {
	return BoolProperty { PropertyBase { instance, property } }
}

func (self *BoolProperty) Get() (error, bool) {
	var val C.int = 0
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err, false
	}
	err = CHECKED(C.syz_getI(&val, *handle, self.property))
	if err != nil {
		return err, false
	}
	return nil, bool(C.bld(val))
}

func (self *BoolProperty) Set(val bool) error {
	var value int
	if val == true {
		value = 1
	} else {
		value = 0
	}
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_setI(*handle, self.property, C.int(value)))
	if err != nil {
		return err
	}
	return nil
}

type DoubleProperty struct {
	PropertyBase
}
func NewDoubleProperty(instance *ObjectBase, property C.int) DoubleProperty {
	return DoubleProperty { PropertyBase { instance, property } }
}

func (self *DoubleProperty) Get() (error, float32) {
	var val C.double = 0.0
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err, 0.0
	}
	err = CHECKED(C.syz_getD(&val, *handle, self.property))
	if err != nil {
		return err, 0.0
	}
	return nil, float32(val)
}

func (self *DoubleProperty) Set(value float32) error {
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_setD(*handle, self.property, C.double(value)))
	if err != nil {
		return err
	}
	return nil
}

type Double3Property struct {
	PropertyBase
}
func NewDouble3Property(instance *ObjectBase, property C.int) Double3Property {
	return Double3Property { PropertyBase { instance, property } }
}

func (self *Double3Property) Get() (error, float32, float32, float32) {
	var x, y, z C.double = 0.0, 0.0, 0.0
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err, 0.0, 0.0, 0.0
	}
	err = CHECKED(C.syz_getD3(&x, &y, &z, *handle, self.property))
	if err != nil {
		return err, 0.0, 0.0, 0.0
	}
	return nil, float32(x), float32(y), float32(z)
}

func (self *Double3Property) Set(x float32, y float32, z float32) error {
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_setD3(*handle, self.property, C.double(x), C.double(y), C.double(z)))
	if err != nil {
		return err
	}
	return nil
}

type Double6Property struct {
	PropertyBase
}
func NewDouble6Property(instance *ObjectBase, property C.int) Double6Property {
	return Double6Property { PropertyBase { instance, property } }
}

func (self *Double6Property) Get() (error, float32, float32, float32, float32, float32, float32) {
	var x1, y1, z1, x2, y2, z2 C.double = 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
	}
	err = CHECKED(C.syz_getD6(&x1, &y1, &z1, &x2, &y2, &z2, *handle, self.property))
	if err != nil {
		return err, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
	}
	return nil, float32(x1), float32(y1), float32(z1), float32(x2), float32(y2), float32(z2)
}

func (self *Double6Property) Set(x1 float32, y1 float32, z1 float32, x2 float32, y2 float32, z2 float32) error {
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_setD6(*handle, self.property, C.double(x1), C.double(y1), C.double(z1), C.double(x2), C.double(y2), C.double(z2)))
	if err != nil {
		return err
	}
	return nil
}

type ObjectProperty struct {
	PropertyBase
}
func NewObjectProperty(instance *ObjectBase, property C.int) ObjectProperty {
	return ObjectProperty { PropertyBase { instance, property } }
}

func (self *ObjectProperty) Set(value ObjectBase) error {
	err, handle := self.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_setO(*handle, self.property, *value.handle))
	if err != nil {
		return err
	}
	return nil
}

type Pausable struct {
	ObjectBase
	current_time, suggested_automation_time DoubleProperty
}

func NewPausable(handle *C.syz_Handle) Pausable {
	pb := Pausable {}
	pb.ObjectBase = ObjectBase { handle } 
	pb.current_time = NewDoubleProperty(&pb.ObjectBase, P_CURRENT_TIME)
	pb.suggested_automation_time = NewDoubleProperty(&pb.ObjectBase, P_SUGGESTED_AUTOMATION_TIME)
	return pb
}
func (self *Pausable) Play() error {
	return CHECKED(C.syz_play(*self.handle))
}

func (self *Pausable) Pause() error {
	return CHECKED(C.syz_pause(*self.handle))
}


type Context struct {
	Pausable
	Gain, Default_distance_ref, Default_distance_max, Default_rolloff, Default_closeness_boost, Default_closeness_boost_distance DoubleProperty
	Position Double3Property
	Orientation Double6Property
	Default_distance_model, Default_panner_strategy IntProperty
}
func NewContext() (error, *Context) {
	handle := C.create_handle()
	err := CHECKED(C.syz_createContext(&handle, nil, nil))
	if err != nil {
		return err, nil
	}
	self := Context{}
	self.Pausable = NewPausable(&handle)
	self.Gain = NewDoubleProperty(&self.ObjectBase, P_GAIN)
	self.Position = NewDouble3Property(&self.ObjectBase, P_POSITION)
	self.Orientation = NewDouble6Property(&self.ObjectBase, P_ORIENTATION)
	self.Default_distance_model = NewIntProperty(&self.ObjectBase, P_DEFAULT_DISTANCE_MODEL)
	self.Default_distance_ref = NewDoubleProperty(&self.ObjectBase, P_DEFAULT_DISTANCE_REF)
	self.Default_distance_max = NewDoubleProperty(&self.ObjectBase, P_DEFAULT_DISTANCE_MAX)
	self.Default_rolloff = NewDoubleProperty(&self.ObjectBase, P_DEFAULT_ROLLOFF)
	self.Default_closeness_boost = NewDoubleProperty(&self.ObjectBase, P_DEFAULT_CLOSENESS_BOOST)
	self.Default_closeness_boost_distance = NewDoubleProperty(&self.ObjectBase, P_DEFAULT_CLOSENESS_BOOST_DISTANCE)
	self.Default_panner_strategy = NewIntProperty(&self.ObjectBase, P_DEFAULT_PANNER_STRATEGY)
	return nil, &self
}

type StreamHandle struct {
	ObjectBase
}
func newStreamHandle(handle *C.syz_Handle) *StreamHandle {
	return &StreamHandle { ObjectBase { handle } }
}

func StreamHandleFromFile(path string) (error, *StreamHandle) {
	handle := C.create_handle()
	ph := C.CString(path)
	defer C.free(unsafe.Pointer(ph))
	err := CHECKED(C.syz_createStreamHandleFromFile(&handle, ph, nil, nil))
	if err != nil {
		return err, nil
	}
	return nil, newStreamHandle(&handle)
}

// Probably needs improvement, Now it only accepts a string because I can't figure a way to safely pass everything from a byte array as const *char, But if you do, Then please submit a pull request.
func StreamHandleFromMemory(data string) (error, *StreamHandle) {
	dt := C.CString(data)
	handle := C.create_handle()
	err := CHECKED(C.syz_createStreamHandleFromMemory(&handle, C.ulonglong(len(data)), dt, nil, nil))
	if err != nil {
		return err, nil
	}
	sh := newStreamHandle(&handle)
	// We're doing this instead of synthizer FreeCallbacks because you can't pass go functions to C.
	// IF you have a better way for doing this, Please submit a pull request.
	runtime.SetFinalizer(sh, func(Sh *StreamHandle) {
		C.free(unsafe.Pointer(dt))
	})
	return nil, sh
}
// We're missing custom stream handles and protocols because I'm not properly sure how to implement them, If you do, Then please send a pull request.

type Generator struct {
	Pausable
	Gain, PitchBend DoubleProperty
}
func newGenerator(handle *C.syz_Handle) *Generator {
	self := Generator{}
	self.Pausable = NewPausable(handle)
	self.Gain = NewDoubleProperty(&self.ObjectBase, P_GAIN)
	self.PitchBend = NewDoubleProperty(&self.ObjectBase, P_PITCH_BEND)
	return &self
}


type StreamingGenerator struct {
	Generator
	Looping BoolProperty
	PlaybackPosition DoubleProperty
}
func newStreamingGenerator(handle *C.syz_Handle) *StreamingGenerator {
	self := StreamingGenerator{}
	self.Generator = *newGenerator(handle)
	self.Looping = NewBoolProperty(&self.ObjectBase, P_LOOPING)
	self.PlaybackPosition = NewDoubleProperty(&self.ObjectBase, P_PLAYBACK_POSITION)
	return &self
}

func StreamingGeneratorFromFile(ctx *Context, path string) (error, *StreamingGenerator) {
	ph := C.CString(path)
	defer C.free(unsafe.Pointer(ph))
	out := C.create_handle()
	err := CHECKED(C.syz_createStreamingGeneratorFromFile(&out, *ctx.GetHandle(), ph, nil, nil, nil))
	if err != nil {
		return err, nil
	}
	return nil, newStreamingGenerator(&out)
}

func StreamingGeneratorFromHandle(ctx *Context, stream *StreamHandle) (error, *StreamingGenerator) {
	handle := C.create_handle()
	err := CHECKED(C.syz_createStreamingGeneratorFromStreamHandle(&handle, *ctx.handle, *stream.handle, nil, nil, nil))
	if err != nil {
		return err, nil
	}
	return nil, newStreamingGenerator(&handle)
}

type Source struct {
	Pausable
	Gain DoubleProperty
}

func newSource(handle *C.syz_Handle) *Source {
	self := Source{}
	self.Pausable = NewPausable(handle)
	self.Gain = NewDoubleProperty(&self.ObjectBase, P_GAIN)
	return &self
}

func (self *Source) AddGenerator(gen Generator) error {
	err, h := gen.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_sourceAddGenerator(*self.handle, *h))
	if err != nil {
		return err
	}
	return nil
}
func (self *Source) RemoveGenerator(gen Generator) error {
	err, h := gen.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_sourceRemoveGenerator(*self.handle, *h))
	if err != nil {
		return err
	}
	return nil
}


type DirectSource struct {
	Source
}
func NewDirectSource(ctx *Context) (error, *DirectSource) {
	out := C.create_handle()
	err := CHECKED(C.syz_createDirectSource(&out, *ctx.handle, nil, nil, nil))
	if err != nil {
		return err, nil
	}
	return nil, &DirectSource { *newSource(&out) }
}

type AngularPannedSource struct {
	Source
	Azimuth, Elevation DoubleProperty
}

func NewAngularPannedSource(ctx *Context, panner_strategy C.int, azimuth float32, elevation float32) (error, *AngularPannedSource) {
	out := C.create_handle()
	err := CHECKED(C.syz_createAngularPannedSource(&out, *ctx.handle, panner_strategy, C.double(azimuth), C.double(elevation), nil, nil, nil))
	if err != nil {
		return err, nil
	}
	self := AngularPannedSource{}
	self.Source = *newSource(&out)
	self.Azimuth = NewDoubleProperty(&self.ObjectBase, P_AZIMUTH)
	self.Elevation = NewDoubleProperty(&self.ObjectBase, P_ELEVATION)
	return nil, &self
}

type ScalarPannedSource struct {
	Source
	PanningScalar DoubleProperty
}

func NewScalarPannedSource(ctx *Context, panner_strategy C.int, panning_scalar float32) (error, *ScalarPannedSource) {
	out := C.create_handle()
	err := CHECKED(C.syz_createScalarPannedSource(&out, *ctx.handle, panner_strategy, C.double(panning_scalar), nil, nil, nil))
	if err != nil {
		return err, nil
	}
	self := ScalarPannedSource{}
	self.Source = *newSource(&out)
	self.PanningScalar = NewDoubleProperty(&self.ObjectBase, P_PANNING_SCALAR)
	return nil, &self
}

type Source3D struct {
	Source
	DistanceRef, DistanceMax, Rolloff, ClosenessBoost, ClosenessBoostDistance DoubleProperty
	Position Double3Property
	Orientation Double6Property
	DistanceModel IntProperty
}

func NewSource3D(ctx *Context) (error, *Source3D) {
	out := C.create_handle()
	err := CHECKED(C.syz_createSource3D(&out, *ctx.handle, C.int(PANNER_STRATEGY_DELEGATE), C.double(0.0), C.double(0.0), C.double(0.0), nil, nil, nil))
	if err != nil {
		return err, nil
	}
	self := Source3D{}
	self.Source = *newSource(&out)
	self.DistanceModel = NewIntProperty(&self.ObjectBase, P_DISTANCE_MODEL)
	self.DistanceRef = NewDoubleProperty(&self.ObjectBase, P_DISTANCE_REF)
	self.DistanceMax = NewDoubleProperty(&self.ObjectBase, P_DISTANCE_MAX)
	self.Rolloff = NewDoubleProperty(&self.ObjectBase, P_ROLLOFF)
	self.ClosenessBoost = NewDoubleProperty(&self.ObjectBase, P_CLOSENESS_BOOST)
	self.ClosenessBoostDistance = NewDoubleProperty(&self.ObjectBase, P_CLOSENESS_BOOST_DISTANCE)
	self.Position = NewDouble3Property(&self.ObjectBase, P_POSITION)
	self.Orientation = NewDouble6Property(&self.ObjectBase, P_ORIENTATION)
	return nil, &self
}



type Buffer struct {
	ObjectBase
}

func newBuffer(handle *C.syz_Handle) *Buffer {
	return &Buffer { ObjectBase { handle } }
}
func BufferFromFile(path string) (error, *Buffer) {
	handle := NewHandle()
	ph := C.CString(path)
	defer C.free(unsafe.Pointer(ph))
		err := CHECKED(C.syz_createBufferFromFile(&handle, ph, nil, nil))
	if err != nil {
		return err, nil
	}
	return nil, newBuffer(&handle)
}

func BufferFromEncodedData(data string) (error, *Buffer) {
	handle := NewHandle()
	dt := C.CString(data)
	defer C.free(unsafe.Pointer(dt))
	length := len(data)
	if length == 0 {
		return errors.New("Cannot safely pass empty arrays to synthizer."), nil
	}
	err := CHECKED(C.syz_createBufferFromEncodedData(&handle, C.ulonglong(length), dt, nil, nil))
	if err != nil {
		return err, nil
	}
	return nil, newBuffer(&handle)
}
// Todo: Add BufferFromFloatArray function.

func BufferFromStreamHandle(stream StreamHandle) (error, *Buffer) {
	handle := NewHandle()
	err := CHECKED(C.syz_createBufferFromStreamHandle(&handle, *stream.handle, nil, nil))
	if err != nil {
		return err, nil
	}
	return nil, newBuffer(&handle)
}
func (self *Buffer) GetChannels() (error, int) {
	var ret C.uint
	err := CHECKED(C.syz_bufferGetChannels(&ret, *self.handle))
	if err != nil {
		return err, 0
	}
	return nil, int(ret)
}

func (self *Buffer) GetLengthInSamples() (error, int) {
	var ret C.uint
	err := CHECKED(C.syz_bufferGetLengthInSamples(&ret, *self.handle))
	if err != nil {
		return err, 0
	}
	return nil, int(ret)
}

func (self *Buffer) GetLengthInSeconds() (error, float64) {
	var ret C.double
	err := CHECKED(C.syz_bufferGetLengthInSeconds(&ret, *self.handle))
	if err != nil {
		return err, 0.0
	}
	return nil, float64(ret)
}

// Todo: Add Buffer.GetLengthInBytes
type BufferGenerator struct {
	Generator
	Looping BoolProperty
	Buffer ObjectProperty
	PlaybackPosition DoubleProperty
}

func NewBufferGenerator(ctx *Context) (error, *BufferGenerator) {
	handle := NewHandle()
	err, ctx_h := ctx.GetHandleChecked()
	if err != nil {
		return err, nil
	}
	err = CHECKED(C.syz_createBufferGenerator(&handle, *ctx_h, nil, nil, nil))
	if err != nil {
		return err, nil
	}
	self := BufferGenerator{}
	self.Generator = *newGenerator(&handle)
	self.Buffer = NewObjectProperty(&self.ObjectBase, P_BUFFER)
	self.PlaybackPosition = NewDoubleProperty(&self.ObjectBase, P_PLAYBACK_POSITION)
	self.Looping = NewBoolProperty(&self.ObjectBase, P_LOOPING)
	return nil, &self
}


type LibraryConfig struct {
	log_level LogLevel
	logging_backend LoggingBackend
	libsndfile_path string
}

func NewLibraryConfig(log_level LogLevel, logging_backend LoggingBackend) LibraryConfig {
	libc := LibraryConfig {
		log_level: log_level,
		logging_backend: logging_backend,
		libsndfile_path: "",
	}
	return libc
}

func (self *LibraryConfig) SetLibsndfilePath(path string) {
	self.libsndfile_path = path
}
func InitializeWithConfig(config *LibraryConfig) {
	if config.libsndfile_path == "" {
		conf := C.create_LibraryConfig(C.uint(config.log_level), C.uint(config.logging_backend), nil)
		CHECKED(C.syz_initializeWithConfig(&conf))
	} else {
		str := C.CString(config.libsndfile_path)
		conf := C.create_LibraryConfig(C.uint(config.log_level), C.uint(config.logging_backend), str)
		CHECKED(C.syz_initializeWithConfig(&conf))
		C.free(unsafe.Pointer(str))
	}
}
func GOCHECK(err error) error {
	if err != nil {
		panic(err)
	}
	return nil
}

// Convert A C char* to a GoString, Freeing the C char and returning the go String.
func CCharToGoString(ch *C.char) string {
	str := C.GoString(ch)
	C.free(unsafe.Pointer(ch))
	return str
}

func CHECKED(x C.syz_ErrorCode) error {
	var ret int = int(x)
	if ret > 0 {
		return errors.New("SynthizerError: " + CCharToGoString(C.syz_getLastErrorMessage()) + " [" + strconv.Itoa(ret) + "]")
	}
	return nil
}

func NewHandle() C.syz_Handle {
	return C.create_handle()
}

func Initialize() error {
	err := CHECKED(C.syz_initialize())
	return err
}

func Shutdown() error {
	err := CHECKED(C.syz_shutdown())
	return err
}