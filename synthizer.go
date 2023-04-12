//MIT License
//
//!
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

// Synthizer main object interface.
type ObjectBase interface {
	GetHandleChecked() (*C.syz_Handle, error)
	GetHandle() *C.syz_Handle
	Destroy() error
	ConfigDeleteBehavior(int, float32) error
}
// The synthizer main object struct.
type objectBase struct {
	handle *C.syz_Handle
}
func _NewObjectBase(hdl *C.syz_Handle) objectBase {
	return objectBase { handle: hdl }
}

func (self *objectBase) GetHandleChecked() (*C.syz_Handle, error) {
	// Any attempt to reference a non-existing object should raise an error
	handle := self.GetHandle()
	if *handle == 0 {
		return nil, errors.New("SynthizerError: Object no longer exist")
	}
	return handle, nil
}

func (self *objectBase) Destroy() error {
	err := CHECKED(C.syz_handleDecRef(*self.handle))
	//*self.handle = 0
	return err
}

func (self *objectBase) GetHandle() *C.syz_Handle {
	return self.handle
}
func (self *objectBase) ConfigDeleteBehavior(linger int, timeout float32) error {
	cfg := C.create_DeleteBehaviorConfig(C.int(linger), C.double(timeout))
	C.syz_initDeleteBehaviorConfig(&cfg)
	return CHECKED(C.syz_configDeleteBehavior(*self.handle, &cfg))
}

type PropertyBase struct {
	instance *objectBase
	property C.int
}
func (self *PropertyBase) GetHandleChecked() (*C.syz_Handle, error) {
	// Any attempt to reference a non-existing object should raise an error
	obj := reflect.ValueOf(*self.instance)
	if obj.IsZero() {
		return nil, errors.New("SynthizerError: Object no longer exist")
	}
	handle := (*self.instance).GetHandle()
	if *handle == 0 {
		return nil, errors.New("SynthizerError: Object no longer exist")
	}
	return handle, nil
}
func (self *PropertyBase) GetProperty() C.int {
	return self.property
}

type IntProperty struct {
	PropertyBase
}

func NewIntProperty(instance *objectBase, property C.int) IntProperty {
	return IntProperty { PropertyBase { instance, property } }
}

func (self *IntProperty) Get() (int, error) {
	var val C.int = 0
	handle, err := self.GetHandleChecked()
	if err != nil {
		return 0, err
	}
	err = CHECKED(C.syz_getI(&val, *handle, self.property))
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func (self *IntProperty) Set(value int) error {
	handle, err := self.GetHandleChecked()
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

func NewBoolProperty(instance *objectBase, property C.int) BoolProperty {
	return BoolProperty { PropertyBase { instance, property } }
}

func (self *BoolProperty) Get() (bool, error) {
	var val C.int = 0
	handle, err := self.GetHandleChecked()
	if err != nil {
		return false, err
	}
	err = CHECKED(C.syz_getI(&val, *handle, self.property))
	if err != nil {
		return false, err
	}
	return bool(C.bld(val)), nil
}

func (self *BoolProperty) Set(val bool) error {
	var value int
	if val == true {
		value = 1
	} else {
		value = 0
	}
	handle, err := self.GetHandleChecked()
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
func NewDoubleProperty(instance *objectBase, property C.int) DoubleProperty {
	return DoubleProperty { PropertyBase { instance, property } }
}

func (self *DoubleProperty) Get() (float32, error) {
	var val C.double = 0.0
	handle, err := self.GetHandleChecked()
	if err != nil {
		return 0.0, err
	}
	err = CHECKED(C.syz_getD(&val, *handle, self.property))
	if err != nil {
		return 0.0, err
	}
	return float32(val), nil
}

func (self *DoubleProperty) Set(value float32) error {
	handle, err := self.GetHandleChecked()
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
func NewDouble3Property(instance *objectBase, property C.int) Double3Property {
	return Double3Property { PropertyBase { instance, property } }
}

func (self *Double3Property) Get() (float32, float32, float32, error) {
	var x, y, z C.double = 0.0, 0.0, 0.0
	handle, err := self.GetHandleChecked()
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}
	err = CHECKED(C.syz_getD3(&x, &y, &z, *handle, self.property))
	if err != nil {
		return 0.0, 0.0, 0.0, err
	}
	return float32(x), float32(y), float32(z), err
}

func (self *Double3Property) Set(x float32, y float32, z float32) error {
	handle, err := self.GetHandleChecked()
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
func NewDouble6Property(instance *objectBase, property C.int) Double6Property {
	return Double6Property { PropertyBase { instance, property } }
}

func (self *Double6Property) Get() (float32, float32, float32, float32, float32, float32, error) {
	var x1, y1, z1, x2, y2, z2 C.double = 0.0, 0.0, 0.0, 0.0, 0.0, 0.0
	handle, err := self.GetHandleChecked()
	if err != nil {
		return 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, err
	}
	err = CHECKED(C.syz_getD6(&x1, &y1, &z1, &x2, &y2, &z2, *handle, self.property))
	if err != nil {
		return 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, err
	}
	return float32(x1), float32(y1), float32(z1), float32(x2), float32(y2), float32(z2), err
}

func (self *Double6Property) Set(x1 float32, y1 float32, z1 float32, x2 float32, y2 float32, z2 float32) error {
	handle, err := self.GetHandleChecked()
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
func NewObjectProperty(instance *objectBase, property C.int) ObjectProperty {
	return ObjectProperty { PropertyBase { instance, property } }
}

func (self *ObjectProperty) Set(value ObjectBase) error {
	handle, err := self.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_setO(*handle, self.property, *value.GetHandle()))
	if err != nil {
		return err
	}
	return nil
}

type BiquadProperty struct {
	PropertyBase
}
func NewBiquadProperty(instance *objectBase, property C.int) BiquadProperty {
	return BiquadProperty { PropertyBase { instance, property } }
}

func (self *BiquadProperty) Set(value BiquadConfig) error {
	handle, err := self.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_setBiquad(*handle, self.property, value.config))
	if err != nil {
		return err
	}
	return nil
}

type Pausable struct {
	objectBase
	current_time, suggested_automation_time DoubleProperty
}

func NewPausable(handle *C.syz_Handle) Pausable {
	pb := Pausable {}
	pb.objectBase = objectBase { handle } 
	pb.current_time = NewDoubleProperty(&pb.objectBase, P_CURRENT_TIME)
	pb.suggested_automation_time = NewDoubleProperty(&pb.objectBase, P_SUGGESTED_AUTOMATION_TIME)
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
func NewContext() (*Context, error) {
	handle := C.create_handle()
	err := CHECKED(C.syz_createContext(&handle, nil, nil))
	if err != nil {
		return nil, err
	}
	self := Context{}
	self.Pausable = NewPausable(&handle)
	self.Gain = NewDoubleProperty(&self.objectBase, P_GAIN)
	self.Position = NewDouble3Property(&self.objectBase, P_POSITION)
	self.Orientation = NewDouble6Property(&self.objectBase, P_ORIENTATION)
	self.Default_distance_model = NewIntProperty(&self.objectBase, P_DEFAULT_DISTANCE_MODEL)
	self.Default_distance_ref = NewDoubleProperty(&self.objectBase, P_DEFAULT_DISTANCE_REF)
	self.Default_distance_max = NewDoubleProperty(&self.objectBase, P_DEFAULT_DISTANCE_MAX)
	self.Default_rolloff = NewDoubleProperty(&self.objectBase, P_DEFAULT_ROLLOFF)
	self.Default_closeness_boost = NewDoubleProperty(&self.objectBase, P_DEFAULT_CLOSENESS_BOOST)
	self.Default_closeness_boost_distance = NewDoubleProperty(&self.objectBase, P_DEFAULT_CLOSENESS_BOOST_DISTANCE)
	self.Default_panner_strategy = NewIntProperty(&self.objectBase, P_DEFAULT_PANNER_STRATEGY)
	return &self, nil
}

func (self *Context) ConfigRoute(output ObjectBase, input ObjectBase, optionalV ...float64) error {
	config := C.struct_syz_RouteConfig{}
	err := CHECKED(C.syz_initRouteConfig(&config))
	if err != nil {
		return err
	}
	if len(optionalV) > 0 {
		config.gain = C.double(optionalV[0])
		if len(optionalV) > 1 {
			config.fade_time = C.double(optionalV[1])
		}
	}
	err = CHECKED(C.syz_routingConfigRoute(*self.handle, *output.GetHandle(), *input.GetHandle(), &config))
	if err != nil {
		return err
	}
	return nil
}
func (self *Context) RemoveRoute(output objectBase, input objectBase, optionalV ...float64) error {
	fade_time := 0.01
	if len(optionalV) > 0 {
		fade_time = optionalV[0]
	}
	err := CHECKED(C.syz_routingRemoveRoute(*self.handle, *output.handle, *input.handle, C.double(fade_time)))
	if err != nil {
		return err
	}
	return nil
}

type StreamHandle struct {
	objectBase
}
func newStreamHandle(handle *C.syz_Handle) *StreamHandle {
	return &StreamHandle { objectBase { handle } }
}

func StreamHandleFromFile(path string) (*StreamHandle, error) {
	handle := C.create_handle()
	ph := C.CString(path)
	defer C.free(unsafe.Pointer(ph))
	err := CHECKED(C.syz_createStreamHandleFromFile(&handle, ph, nil, nil))
	if err != nil {
		return nil, err
	}
	return newStreamHandle(&handle), err
}

// Probably needs improvement, Now it only accepts a string because I can't figure a way to safely pass everything from a byte array as const *char, But if you do, Then please submit a pull request.
func StreamHandleFromMemory(data string) (*StreamHandle, error) {
	dt := C.CString(data)
	handle := C.create_handle()
	err := CHECKED(C.syz_createStreamHandleFromMemory(&handle, C.ulonglong(len(data)), dt, nil, nil))
	if err != nil {
		return nil, err
	}
	sh := newStreamHandle(&handle)
	// We're doing this instead of synthizer FreeCallbacks because you can't pass go functions to C.
	// IF you have a better way for doing this, Please submit a pull request.
	runtime.SetFinalizer(sh, func(Sh *StreamHandle) {
		C.free(unsafe.Pointer(dt))
	})
	return sh, nil
}
// We're missing custom stream handles and protocols because I'm not properly sure how to implement them, If you do, Then please send a pull request.

type Generator interface {
	// We use useless functions just to know what implements generator.
	synthizerGeneratorImpl() bool
	ObjectBase
}
type generator struct {
	Pausable
	Gain, PitchBend DoubleProperty
}
func newGenerator(handle *C.syz_Handle) *generator {
	self := generator{}
	self.Pausable = NewPausable(handle)
	self.Gain = NewDoubleProperty(&self.objectBase, P_GAIN)
	self.PitchBend = NewDoubleProperty(&self.objectBase, P_PITCH_BEND)
	return &self
}

func (gen *generator) synthizerGeneratorImpl() bool {
	return true
}

type StreamingGenerator struct {
	generator
	Looping BoolProperty
	PlaybackPosition DoubleProperty
}
func newStreamingGenerator(handle *C.syz_Handle) *StreamingGenerator {
	self := StreamingGenerator{}
	self.generator = *newGenerator(handle)
	self.Looping = NewBoolProperty(&self.objectBase, P_LOOPING)
	self.PlaybackPosition = NewDoubleProperty(&self.objectBase, P_PLAYBACK_POSITION)
	return &self
}

func StreamingGeneratorFromFile(ctx *Context, path string) (*StreamingGenerator, error) {
	ph := C.CString(path)
	defer C.free(unsafe.Pointer(ph))
	out := C.create_handle()
	err := CHECKED(C.syz_createStreamingGeneratorFromFile(&out, *ctx.GetHandle(), ph, nil, nil, nil))
	if err != nil {
		return nil, err
	}
	return newStreamingGenerator(&out), err
}

func StreamingGeneratorFromHandle(ctx *Context, stream *StreamHandle) (*StreamingGenerator, error) {
	handle := C.create_handle()
	err := CHECKED(C.syz_createStreamingGeneratorFromStreamHandle(&handle, *ctx.handle, *stream.handle, nil, nil, nil))
	if err != nil {
		return nil, err
	}
	return newStreamingGenerator(&handle), err
}
type Source interface {
	AddGenerator(Generator) error
	RemoveGenerator(Generator) error
	ObjectBase
}

type source struct {
	Pausable
	Gain DoubleProperty
}

func newSource(handle *C.syz_Handle) *source {
	self := source{}
	self.Pausable = NewPausable(handle)
	self.Gain = NewDoubleProperty(&self.objectBase, P_GAIN)
	return &self
}

func (self *source) AddGenerator(gen Generator) error {
	h, err := gen.GetHandleChecked()
	if err != nil {
		return err
	}
	err = CHECKED(C.syz_sourceAddGenerator(*self.handle, *h))
	if err != nil {
		return err
	}
	return nil
}
func (self *source) RemoveGenerator(gen Generator) error {
	h, err := gen.GetHandleChecked()
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
	source
}
func NewDirectSource(ctx *Context) (*DirectSource, error) {
	out := C.create_handle()
	err := CHECKED(C.syz_createDirectSource(&out, *ctx.handle, nil, nil, nil))
	if err != nil {
		return nil, err
	}
	return &DirectSource { *newSource(&out) }, nil
}

type AngularPannedSource struct {
	source
	Azimuth, Elevation DoubleProperty
}

func NewAngularPannedSource(ctx *Context, panner_strategy C.int, OV ...float32) (*AngularPannedSource, error) {
	var azimuth float32 = 0.0
	var elevation float32 = 0.0
	if len(OV) > 0 {
		azimuth = OV[0]
		if len(OV) > 1 {
			elevation = OV[1]
		}
	}
	out := C.create_handle()
	err := CHECKED(C.syz_createAngularPannedSource(&out, *ctx.handle, panner_strategy, C.double(azimuth), C.double(elevation), nil, nil, nil))
	if err != nil {
		return nil, err
	}
	self := AngularPannedSource{}
	self.source = *newSource(&out)
	self.Azimuth = NewDoubleProperty(&self.objectBase, P_AZIMUTH)
	self.Elevation = NewDoubleProperty(&self.objectBase, P_ELEVATION)
	return &self, nil
}

type ScalarPannedSource struct {
	source
	PanningScalar DoubleProperty
}

func NewScalarPannedSource(ctx *Context, panner_strategy C.int, panning_scalar float32) (*ScalarPannedSource, error) {
	out := C.create_handle()
	err := CHECKED(C.syz_createScalarPannedSource(&out, *ctx.handle, panner_strategy, C.double(panning_scalar), nil, nil, nil))
	if err != nil {
		return nil, err
	}
	self := ScalarPannedSource{}
	self.source = *newSource(&out)
	self.PanningScalar = NewDoubleProperty(&self.objectBase, P_PANNING_SCALAR)
	return &self, err
}

type Source3D struct {
	source
	DistanceRef, DistanceMax, Rolloff, ClosenessBoost, ClosenessBoostDistance DoubleProperty
	Position Double3Property
	Orientation Double6Property
	DistanceModel IntProperty
}

func NewSource3D(ctx *Context) (*Source3D, error) {
	out := C.create_handle()
	err := CHECKED(C.syz_createSource3D(&out, *ctx.handle, C.int(PANNER_STRATEGY_DELEGATE), C.double(0.0), C.double(0.0), C.double(0.0), nil, nil, nil))
	if err != nil {
		return nil, err
	}
	self := Source3D{}
	self.source = *newSource(&out)
	self.DistanceModel = NewIntProperty(&self.objectBase, P_DISTANCE_MODEL)
	self.DistanceRef = NewDoubleProperty(&self.objectBase, P_DISTANCE_REF)
	self.DistanceMax = NewDoubleProperty(&self.objectBase, P_DISTANCE_MAX)
	self.Rolloff = NewDoubleProperty(&self.objectBase, P_ROLLOFF)
	self.ClosenessBoost = NewDoubleProperty(&self.objectBase, P_CLOSENESS_BOOST)
	self.ClosenessBoostDistance = NewDoubleProperty(&self.objectBase, P_CLOSENESS_BOOST_DISTANCE)
	self.Position = NewDouble3Property(&self.objectBase, P_POSITION)
	self.Orientation = NewDouble6Property(&self.objectBase, P_ORIENTATION)
	return &self, nil
}



type Buffer struct {
	objectBase
}

func newBuffer(handle *C.syz_Handle) *Buffer {
	return &Buffer { objectBase { handle } }
}
func BufferFromFile(path string) (*Buffer, error) {
	handle := NewHandle()
	ph := C.CString(path)
	defer C.free(unsafe.Pointer(ph))
		err := CHECKED(C.syz_createBufferFromFile(&handle, ph, nil, nil))
	if err != nil {
		return nil, err
	}
	return newBuffer(&handle), err
}

func BufferFromEncodedData(data string) (*Buffer, error) {
	handle := NewHandle()
	dt := C.CString(data)
	defer C.free(unsafe.Pointer(dt))
	length := len(data)
	if length == 0 {
		return nil, errors.New("Cannot safely pass empty arrays to synthizer.")
	}
	err := CHECKED(C.syz_createBufferFromEncodedData(&handle, C.ulonglong(length), dt, nil, nil))
	if err != nil {
		return nil, err
	}
	return newBuffer(&handle), err
}
// Todo: Add BufferFromFloatArray function.

func BufferFromStreamHandle(stream StreamHandle) (*Buffer, error) {
	handle := NewHandle()
	err := CHECKED(C.syz_createBufferFromStreamHandle(&handle, *stream.handle, nil, nil))
	if err != nil {
		return nil, err
	}
	return newBuffer(&handle), err
}
func (self *Buffer) GetChannels() (int, error) {
	var ret C.uint
	err := CHECKED(C.syz_bufferGetChannels(&ret, *self.handle))
	if err != nil {
		return 0, err
	}
	return int(ret), nil
}

func (self *Buffer) GetLengthInSamples() (int, error) {
	var ret C.uint
	err := CHECKED(C.syz_bufferGetLengthInSamples(&ret, *self.handle))
	if err != nil {
		return 0, err
	}
	return int(ret), nil
}

func (self *Buffer) GetLengthInSeconds() (float64, error) {
	var ret C.double
	err := CHECKED(C.syz_bufferGetLengthInSeconds(&ret, *self.handle))
	if err != nil {
		return 0.0, err
	}
	return float64(ret), err
}

// Todo: Add Buffer.GetLengthInBytes
type BufferGenerator struct {
	generator
	Looping BoolProperty
	Buffer ObjectProperty
	PlaybackPosition DoubleProperty
}

func NewBufferGenerator(ctx *Context) (*BufferGenerator, error) {
	handle := NewHandle()
	ctx_h, err := ctx.GetHandleChecked()
	if err != nil {
		return nil, err
	}
	err = CHECKED(C.syz_createBufferGenerator(&handle, *ctx_h, nil, nil, nil))
	if err != nil {
		return nil, err
	}
	self := BufferGenerator{}
	self.generator = *newGenerator(&handle)
	self.Buffer = NewObjectProperty(&self.objectBase, P_BUFFER)
	self.PlaybackPosition = NewDoubleProperty(&self.objectBase, P_PLAYBACK_POSITION)
	self.Looping = NewBoolProperty(&self.objectBase, P_LOOPING)
	return &self, err
}


type NoiseGenerator struct {
	generator
	NoiseType IntProperty
}

func NewNoiseGenerator(ctx *Context, channels int) (*NoiseGenerator, error) {
	handle := NewHandle()
	ctx_h, err := ctx.GetHandleChecked()
	if err != nil {
		return nil, err
	}
	err = CHECKED(C.syz_createNoiseGenerator(&handle, *ctx_h, C.uint(channels), nil, nil, nil))
	if err != nil {
		return nil, err
	}
	self := NoiseGenerator{}
	self.generator = *newGenerator(&handle)
	self.NoiseType = NewIntProperty(&self.objectBase, P_NOISE_TYPE)
	return &self, nil
}


type BiquadConfig struct {
	config *C.struct_syz_BiquadConfig
}

func newBiquadConfig() (*BiquadConfig, error) {
	self := BiquadConfig { &C.struct_syz_BiquadConfig{} }
	err := CHECKED(C.syz_biquadDesignIdentity(self.config))
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func BiquadConfigDesignIdentity() (*BiquadConfig, error) {
	out, err := newBiquadConfig()
	if err != nil {
		return nil, err
	}
	err = CHECKED(C.syz_biquadDesignIdentity(out.config))
	if err != nil {
		return nil, err
	}
	return out, err
}
func BiquadConfigDesignLowpass(frequency float64, OV ...float64) (*BiquadConfig, error) {
	out, err := newBiquadConfig()
	if err != nil {
		return nil, err
	}
	q := 0.7071135624381276
	if len(OV) > 0 {
		q = OV[0]
	}
	err = CHECKED(C.syz_biquadDesignLowpass(out.config, C.double(frequency), C.double(q)))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func BiquadConfigDesignHighpass(frequency float64, OV ...float64) (*BiquadConfig, error) {
	out, err := newBiquadConfig()
	if err != nil {
		return nil, err
	}
	q := 0.7071135624381276
	if len(OV) > 0 {
		q = OV[0]
	}
	err = CHECKED(C.syz_biquadDesignHighpass(out.config, C.double(frequency), C.double(q)))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func BiquadDesignBandpass(frequency float64, bandwidth float64) (*BiquadConfig, error) {
	out, err := newBiquadConfig()
	if err != nil {
		return nil, err
	}
	err = CHECKED(C.syz_biquadDesignBandpass(out.config, C.double(frequency), C.double(bandwidth)))
	if err != nil {
		return nil, err
	}
	return out, nil
}

type GlobalEffect struct {
	objectBase
	FilterInput BiquadProperty
	Gain DoubleProperty
}

func newGlobalEffect(handle *C.syz_Handle) *GlobalEffect {
	self := GlobalEffect{}
	self.objectBase = objectBase { handle }
	self.Gain = NewDoubleProperty(&self.objectBase, P_GAIN)
	self.FilterInput = NewBiquadProperty(&self.objectBase, P_FILTER_INPUT)
	return &self
}
func (self *GlobalEffect) Reset() error {
	return CHECKED(C.syz_effectReset(*self.handle))
}

type EchoTapConfig struct {
	Delay float64
	GainL float64
	GainR float64
}
func NewEchoTapConfig(delay float64, gain_l float64, gain_r float64) *EchoTapConfig {
	return &EchoTapConfig { delay, gain_l, gain_r }
}

// Todo, Add GlobalEcho.
type GlobalFdnReverb struct {
	GlobalEffect
	MeanFreePath, T60, LateReflectionsLfRolloff, LateReflectionsLfReference, LateReflectionsHfRolloff, LateReflectionsHfReference, LateReflectionsDiffusion, LateReflectionsModulationDepth, LateReflectionsModulationFrequency, LateReflectionsDelay DoubleProperty	
}

func NewGlobalFdnReverb(ctx *Context) (*GlobalFdnReverb, error) {
	handle := NewHandle()
	h, err := ctx.GetHandleChecked()
	if err != nil {
		return nil, err
	}
	err = CHECKED(C.syz_createGlobalFdnReverb(&handle, *h, nil, nil, nil))
	if err != nil {
		return nil, err
	}
	self := GlobalFdnReverb {}
	self.GlobalEffect = *newGlobalEffect(&handle)
	self.MeanFreePath = NewDoubleProperty(&self.objectBase, P_MEAN_FREE_PATH)
	self.T60 = NewDoubleProperty(&self.objectBase, P_T60)
	self.LateReflectionsLfRolloff = NewDoubleProperty(&self.objectBase, P_LATE_REFLECTIONS_LF_ROLLOFF)
	self.LateReflectionsLfReference = NewDoubleProperty(&self.objectBase, P_LATE_REFLECTIONS_LF_REFERENCE)
	self.LateReflectionsHfRolloff = NewDoubleProperty(&self.objectBase, P_LATE_REFLECTIONS_HF_ROLLOFF)
	self.LateReflectionsHfReference = NewDoubleProperty(&self.objectBase, P_LATE_REFLECTIONS_HF_REFERENCE)
	self.LateReflectionsDiffusion = NewDoubleProperty(&self.objectBase, P_LATE_REFLECTIONS_DIFFUSION)
	self.LateReflectionsModulationDepth = NewDoubleProperty(&self.objectBase, P_LATE_REFLECTIONS_MODULATION_DEPTH)
	self.LateReflectionsModulationFrequency = NewDoubleProperty(&self.objectBase, P_LATE_REFLECTIONS_MODULATION_FREQUENCY)
	self.LateReflectionsDelay = NewDoubleProperty(&self.objectBase, P_LATE_REFLECTIONS_DELAY)
	return &self, nil
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