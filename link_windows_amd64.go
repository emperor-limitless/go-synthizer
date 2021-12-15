package synthizer

/*
#cgo windows LDFLAGS: -Llibs/windows ${SRCDIR}/libs/windows/libsynthizer.a -static-libgcc -static-libstdc++  -Wl,-Bstatic -lstdc++ -lpthread -Wl,-Bdynamic -lkernel32 -luser32 -lole32 -lwinmm -ldsound -ldxguid
#cgo windows CPPFLAGS: -O1 -I${SRCDIR}/libs/include
#cgo windows CFLAGS: -O1 -I${SRCDIR}/libs/include
*/
import "C"
