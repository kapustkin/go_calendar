protoc --proto_path=api/proto/v1 --proto_path=scripts --go_out=plugins=grpc:pkg/api/v1 calendar.proto


cd D:\Documents\kapustkin\protoc\bin\
.\protoc -I=d:\Documents\kapustkin\go_calendar\api\proto\calendarpb\ --go_out=plugins=grpc:d:\Documents\kapustkin\go_calendar\api\proto\calendarpb\ calendar.proto