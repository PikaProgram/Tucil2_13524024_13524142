# Tucil2_13524024_13524142

An app where you can provide an .obj file containing the
data of vertices and faces, and turns it into a voxelized
version. This project also provides a rudimentary 3D viewer written using Ebitengine.

# HOW TO BUILD

To compile and run this project, you will need:
- Go (1.18 or higher recommended)
- Ebitengine (github.com/hajimehoshi/ebiten/v2)

If you are compiling native Linux binaries, Ebitengine requires standard X11 and OpenGL Cgo dependencies.
For Ubuntu/Debian (or WSL Ubuntu):
```bash
sudo apt update
sudo apt install libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
```
This project uses makefile to compile. 
To compile for Linux:
```bash
make build
```
To compile for Windows:
```bash
make build-windows
```
To clean the build directory:
```bash
make clean
```

# HOW TO RUN

The built binary will be stored in the bin folder. To run the program you will have to run the command
```bash
bin/voxelizer <input.obj> <maxDepth> <output.obj>
```
where you name the path of the output yourself.
Example
```bash
bin/voxelizer convert test/input/cow.obj 5 test/output/cow.obj
```
Remember that the path is relative to the terminal you are running the binary on. Since the example runs the binary from the root project folder (not from the bin folder), the path is relative to the root project folder.

To run the 3D viewer, you will need to run the command:
```bash
bin/viewer <path_to_file.obj>
```
For example:
```bash
bin/viewer test/output/cow_voxel.obj
```
Viewer Controls:
- Left Click + Drag: Rotate the model (X and Y axis).
- Mouse Wheel: Zoom in / Zoom out.
- 'C' Key: Toggle dynamic Backface Culling on/off.




