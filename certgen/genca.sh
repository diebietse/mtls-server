#!/bin/bash
cfssl genkey -initca ./ca.json | cfssljson -bare root
