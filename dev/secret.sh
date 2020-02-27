#!/bin/bash

kubectl create secret generic gateway-access \
-n funceasy \
--from-file=gateway.public.key \
--from-file=gateway.token