## CommonMark

Golang wrapper for CommonMark using cgo.

This project requires libcmark, which can be built thusly:
```sh
git clone https://github.com/commonmark/cmark
cd cmark
mkdir build
cd build
cmake -GNinja ..
ninja
ninja install
```