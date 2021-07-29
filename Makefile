
musashi_dir      = external/Musashi
musashi_objects  = $(musashi_dir)/m68kcpu.o $(musashi_dir)/m68kdasm.o $(musashi_dir)/m68kops.o
shim_objects     = lib/musashi-c-wrapper/shim.o

all: help

help:
	@echo ""
	@echo "make morfe      - build 65c816-only version"
	@echo "make morfe-m68k - build emulator with additional m68k support"
	@echo "make clean      - clean-up binaries"
	@echo "make clean-all  - clean-up binaries and object files"

$(musashi_dir):
	@echo "attempting to clone Musashi into $(musashi_dir)"
	git clone https://github.com/kstenerud/Musashi/ $(musashi_dir)

$(musashi_objects): lib/musashi-c-wrapper/m68kconf.h
	cp -vf $? $(musashi_dir)/
	$(MAKE) -C $(musashi_dir) clean
	$(MAKE) -C $(musashi_dir)

$(shim_objects): lib/musashi-c-wrapper/shim.c
	$(CC) -std=c99 -I $(musashi_dir) -Wall -c $< -o $@

morfe-m68k: $(musashi_dir) $(musashi_objects) $(shim_objects)
	CGO_LDFLAGS_ALLOW=".+/(Musashi|musashi-c-wrapper)/.+\.o" \
	go build --tags m68k -o $@ cmd/gui/*.go
	@echo "type ./$(@) conf/m68k.ini to run"

morfe:
	go build             -o $@ cmd/gui/*.go
	@echo "type ./$(@) conf/c256.ini to run"

clean:
	rm -fv morfe morfe-m68k

clean-all: $(musashi_objects) $(shim_objects)
	rm -fv $^


.PHONY: all help

