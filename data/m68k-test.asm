
	org    $00100000

init:
	moveq  #$43,d0
	move.b d0,$00AFA000
	bra    init
