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
)

type IObject interface {
	Destroy() error
	GetHandle() *C.syz_Handle
}
// The synthizer main object struct.
type ObjectBase struct {
	handle *C.syz_Handle
}
func _NewObjectBase(hdl *C.syz_Handle) ObjectBase {
	return ObjectBase { handle: hdl }
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
	instance *IObject
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

func NewIntProperty(instance *IObject, property C.int) IntProperty {
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


type DoubleProperty struct {
	PropertyBase
}
func NewDoubleProperty(instance *IObject, property C.int) DoubleProperty {
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
func NewDouble3Property(instance *IObject, property C.int) Double3Property {
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
func NewDouble6Property(instance *IObject, property C.int) Double6Property {
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


type Context struct {
	handle *C.syz_Handle
	gain, default_distance_ref, default_distance_max, default_rolloff, default_closeness_boost, default_closeness_boost_distance DoubleProperty
	position Double3Property
	orientation Double6Property
	default_distance_model DISTANCE_MODEL
	default_panner_strategy PANNER_STRATEGY
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