#! /bin/bash

(cd ImageProject && make)

(cd ImageProject/scr/bmpReader && go build readBMPSaveCMP.go)
(cd ImageProject/scr/bmpReader && go build saveBMPreadCMP.go)








