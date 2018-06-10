# To Do List

* Refactor color/buf.go: this is not how sampled colors should be added together. See https://github.com/fogleman/pt/blob/master/pt/buffer.go
* Stop using "image/color" in favour of own type. Converting between the two is a pain.
* Refactor color.Color to use 32bit colors instead of 64bit. No visual difference between the two.
