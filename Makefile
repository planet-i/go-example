PROJECT="example"  # APP=blog 好像加不加引号都没关系
APP=blog

default:
	echo ${PROJECT}

app:
	@echo ${APP}