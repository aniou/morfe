```
git clone https://github.com/aniou/go65c816
cd go65c816
git checkout m68k
cd ..
git clone https://github.com/kstenerud/Musashi/
cd  Musashi
make
cp *.o     ../go65c816/cmd/gui
cp m68k.h  ../go65c816/cmd/gui
cd ../go65c816/cmd/gui
go build -o gui *go
./gui m68k.ini
<press F10 to start m68k and execute few steps>
```

