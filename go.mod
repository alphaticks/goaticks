module github.com/alphaticks/goaticks

go 1.16

require (
	gitlab.com/tachikoma.ai/tickobjects v0.0.0-20210513105529-119529ed7173 // indirect
	gitlab.com/tachikoma.ai/tickstore-go-client v0.0.0-20210513105541-31ab5ec7270b
	google.golang.org/grpc v1.35.0
)

replace gitlab.com/tachikoma.ai/tickstore-go-client => ../../tachikoma.ai/tickstore-go-client
