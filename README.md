# Netron1-Go
See [Netron1](https://docs.google.com/document/d/16yAbgGwr8msl_BopI4x0BRbj4XC2RE0nVQDWUQzYnAw/edit?usp=sharing) google doc for most up to date information.

This iteration is an evaluation phase of Percolation theory and Criticality of virus spreading. It is based on the article [Going critical](http://35.161.88.15/interactive/going-critical/).

# Running simulation
Just run: "```$ go run .```" inside the *Netron1-Go* directory

# Notes
The app is built in two parts: gui and simulation coroutine.

The simulation generates images and saves them to PNG files. The gui shows them.

The simulation messages the gui when a *frame* is done, the gui then renders it.

# Current tasks
- Draw an image with pixels randomly on/off


The app shows individual frames of a simulation. The left/right keys move backwards and forwards respectively. The "End" key always attempts to move to the most current frame.



# Dependencies
