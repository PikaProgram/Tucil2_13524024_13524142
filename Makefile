BIN_DIR         := bin
TARGET_VOX      := $(BIN_DIR)/voxelizer
TARGET_VIEW     := $(BIN_DIR)/viewer
TARGET_VOX_WIN  := $(BIN_DIR)/voxelizer.exe
TARGET_VIEW_WIN := $(BIN_DIR)/viewer.exe

SRC_VOX         := src/voxelizerdriver/main.go
SRC_VIEW        := src/viewerdriver/main.go

# Compiler and Flags
GO              := go
GOFLAGS         := -v
LDFLAGS         := -w -s

# Define phony targets
.PHONY: all build build-voxelizer build-viewer build-windows build-voxelizer-windows build-viewer-windows clean

all: clean build

build: build-voxelizer build-viewer

build-windows: build-voxelizer-windows build-viewer-windows

build-voxelizer:
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(TARGET_VOX) $(SRC_VOX)

build-viewer:
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(TARGET_VIEW) $(SRC_VIEW)

build-voxelizer-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(TARGET_VOX_WIN) $(SRC_VOX)

build-viewer-windows:
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(TARGET_VIEW_WIN) $(SRC_VIEW)

clean:
	rm -f $(TARGET_VOX) $(TARGET_VIEW) $(TARGET_VOX_WIN) $(TARGET_VIEW_WIN)