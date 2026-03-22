TARGET      := bin/main
SRC_DIR     := src

# Compiler and Flags
GO          := go
GOFLAGS     := -v
LDFLAGS		 	:= -w -s

# Define phony targets
.PHONY: all build clean 

all: clean build

build: 
	$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(TARGET) $(SRC_DIR)/main.go

clean:
	rm -f $(TARGET)
