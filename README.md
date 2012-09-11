#gogame
I am working on a "toy" game project and thought it would be interesting to use golang for it.
This project is just an experimental side-project of mine and is not intended to be generic.

#concurrent scene
Currently it only consists of a minimalist entity-component system (scene/actor/property), some tests and a demo of how to use them.

The model makes heavy use of Go's CSP-based concurrency features using goroutines and channels. For very small scenes, this may be a bit overkill - but for large scenes it means the game logic can scale across as many CPU's that are made available to it on the host machine.

The scene itself is very simple and will probably see performance issues in large scenes as described above. Once I get to that point I'll consider new data structures for the scene/actor. For now it should do.

#next steps
I want to grow the demo into a simple scene that handles user inputs, renders graphics and plays audio. That scene then can be used to stress test the engine (in particular the scene) and see how the concurrency model works out.

#license
This software is licensed under the "MIT license". A copy can be found in the "LICENSE" file.
