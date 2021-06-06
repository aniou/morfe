```
git clone https://github.com/kstenerud/Musashi/
cd  Musashi
make
cp *.o     into go65c816/cmd/gui
cp m68k.h  into go65c816/cmd/gui
go build -o gui *go
./gui m68k.ini
<press F10 to start m68k>
```

