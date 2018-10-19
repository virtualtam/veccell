veccell - cellular automata - power gliders
===========================================

*veccell* is a pet project to toy with `Cellular Automata`_ and the `Go`_
programming language.

Building
--------

::

   # get the sources
   $ git clone https://github.com/virtualtam/veccell.git

   # go >= 1.11: enable module support
   $ export GO111MODULE=on

   # build binaries
   $ make build


Automatons
----------

- ``elementary``: an `Elementary Cellular Automaton`_
- ``game-of-life``: a simple implementation of `Conway's Game of Life`_
- ``game-of-life-dx``: a simple implementation of `Conway's Game of Life`_, with
  additional bells and whistles!

Controls
--------

====== =====================
Key    Action
====== =====================
ctrl+c quit
q      quit
r      randomize state
up     increase render delay
down   decrease render delay
====== =====================


.. _Cellular Automata: https://en.wikipedia.org/wiki/Cellular_automaton
.. _Conway's Game of Life: https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life
.. _Elementary Cellular Automaton: https://en.wikipedia.org/wiki/Elementary_cellular_automaton
.. _Go: https://golang.org/
