# go_os - демонстрация работы с операционной системой

## 1. Библиотека "flag" (флаги без дополнительных опций):
___
### 1.1 Без флагов
```bash
go run ./1_flags/main.go
```
### 1.2 С флагами
```bash
go run ./1_flags/main.go -name=John -age=25
```
___

## 2. Библиотека "pflag" (расширенные опции для флагов):
___
### 2.1 Без флагов
```bash
cd ./2_pflag/ && go run main.go
```
### 2.2 С флагами
```bash
cd ./2_pflag/ && go run main.go --name John -a 25
```
___