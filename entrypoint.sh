#!/bin/sh -l
echo $1 $2 $3
eval "/codon-optimizer -i $1 -o $2 -t $3"
