
 ::

  cpu0 - cpu65c816 or m68k
  cpu1 
  gpu0 - graphics unit, usually Vicky
  gabe - universal IO unit, maybe renamed in future
  bus0 - two instances of bus, responsible for
  bus1   translation addresses and routing to
         responsible devices 

         cpu* itself can use local ram withouth bus,
         as well as via bus, it is a compromise between 
         "clean" design and performance (especially
         when Go and C code is mixed, like for m68k)

         TODO: reconsider part marked as '?'

   gui
     |
   platform
          |
          |     
          + cpu0 +- bus0 --------------------+  (65c816)
          |            |                     |
          |         ram0                     |
          |                                  |
          + cpu1 +- bus1 --------------------+  (m68k)
          |      |     |?                    |
          |      +- ram1                     |
          |                                  |
          + gpu0 ----------------------------+
          |                                  |
          + gabe ----------------------------+
          |                                  |
          + ram2 [shared, optional] ---------+



