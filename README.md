#gogame
I am working on a "toy" game project and thought it would be interesting to use golang for it.
This project is just an experimental side-project of mine and is not intended to be generic.

#concurrent scene
Currently it only consists of a minimalist entity-component system (scene/actor/property), some tests and a demo of how to use them.

The model makes heavy use of Go's CSP-based concurrency features using goroutines and channels. For very small scenes, this may be a bit overkill - but for large scenes it means the game logic can scale across as many CPU's that are made available to it on the host machine.

The scene itself is very simple and will probably see performance issues in large scenes as described above. Once I get to that point I'll consider new data structures for the scene/actor. For now it should do.

#run it
To run the (not very exciting) demo:

```
$ cd demo
$ go run demo.go
```

Then you'll see a bunch of this stuff:

```
movement: 16.139ms
render: 32.154ms
movement: 15.999ms
before push: &{Do:0xf840001aa0 X:1 Y:2 Z:3} &{Do:0xf840001af0 X:0 Y:0 Z:0}
after push: &{Do:0xf840001aa0 X:1 Y:2 Z:3} &{Do:0xf840001af0 X:0 Y:0 Z:0}
render one frame
graphics.Pull
movement: 16.089ms
movement: 15.879ms
render: 31.968ms
before push: &{Do:0xf840001aa0 X:4 Y:8 Z:12} &{Do:0xf840001af0 X:2 Y:4 Z:6}
after push: &{Do:0xf840001aa0 X:4 Y:8 Z:12} &{Do:0xf840001af0 X:2 Y:4 Z:6}
```

It's really just debug prints for me, but it may be interesting to you, too.
"movement" and "render" represent an update by the MovementSystem and RenderSystem, which are set to tick about every 16ms and 32ms respectively. The graphcis system needs to transfer the positional data from the Physical component into the Graphical component before it renders, which is the "push" part you can see there. It prints the physical and graphcial properties "before" and "after" the push to compare values at that point. Note that every single system and component runs in it's own goroutine and communicate via channels, making the entire thing work without a single mutex.

This is really just the beginning, when the code grows I'll add better documentation and demos.

#next steps
I want to grow the demo into a simple scene that handles user inputs, renders graphics and plays audio. That scene then can be used to stress test the engine (in particular the scene) and see how the concurrency model works out.

#license
This software is licensed under the "MIT license". A copy can be found in the "LICENSE" file.
