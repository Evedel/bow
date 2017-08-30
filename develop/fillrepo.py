#!/usr/bin/env python3

import subprocess as sp
import random as rd
import os as os

addr = 'localhost:5001'
base = 'evedel/bow'
dftn = 'Dockerfile-test'

BranchCoef = 8
lParent = base + ":latest"

def shell(c):
    sto, ste = sp.Popen((c).split(), stdout=sp.PIPE).communicate()
    if ste == None:
        print("Done:  " + c)
    else:
        print("Error: " + c + '\n' + sto + '\n' + ste)



def genDockerfile(parent, adddir):
    DF='FROM      ' + parent + '\n' + \
       'ADD       ./golang ' +  adddir + ' \n'
    f = open(dftn, 'w')
    f.write(DF)

def fillRepoInherited():
    numPushed = 0
    l1name = addr + '/' + lParent
    shell('docker tag ' + lParent + ' ' + l1name)
    shell('docker push ' + l1name)
    numPushed += 1

    for i in range(BranchCoef):
        l2name = addr + '/' + str(i) + '-' + base + ':latest'
        genDockerfile(l1name, '/' + str(i))
        shell('docker build . -f ' + dftn + ' -t ' + l2name)
        shell('docker push ' + l2name)
        numPushed += 1
        for j in range(BranchCoef):
            l3name = addr + '/' + str(i) + '-' + base + '-' + str(j) + ':latest'
            genDockerfile(l2name, '/' + str(i) + '/' + str(j))
            shell('docker build . -f ' + dftn + ' -t ' + l3name)
            shell('docker push ' + l3name)
            numPushed += 1
            shell('docker rmi ' + l3name)
            for k in range(BranchCoef):
                l3name = addr + '/' + str(i) + '-' + base + '-' + str(j) + ':' + str(k)
                genDockerfile(l2name, '/' + str(i) + '/' + str(j))
                shell('docker build . -f ' + dftn + ' -t ' + l3name)
                shell('docker push ' + l3name)
                numPushed += 1
                shell('docker rmi ' + l3name)
        shell('docker rmi ' + l2name)
    os.remove(dftn)
    shell('docker rmi ' + l1name)
    print('Done: pushed ' + str(numPushed) + ' images ')

def fillRepoInline():
    numPushed = 0
    l1name = addr + '/' + lParent
    shell('docker tag ' + lParent + ' ' + l1name)
    shell('docker push ' + l1name)
    numPushed += 1

    for i in range(BranchCoef**3):
        l2name = addr + '/' + str(i).zfill(4) + '-' + base + ':latest'
        genDockerfile(l1name, '/' + str(i))
        shell('docker build . -f ' + dftn + ' -t ' + l2name)
        shell('docker push ' + l2name)
        shell('docker rmi ' + l2name)
        numPushed += 1

    os.remove(dftn)
    shell('docker rmi ' + l1name)
    print('Done: pushed ' + str(numPushed) + ' images ')

sto, ste = sp.Popen(("docker login " + addr + " -p test -u test").split(), stdout=sp.PIPE).communicate()

if (ste == None):
    print(sto.decode('utf-8')[:-1] + " [" + addr + "]")
    # fillRepoInherited()
    fillRepoInline()
