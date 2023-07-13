package marshal

import "github.com/hiroara/carbo/messaging/message"

type Marshaller[S any] func(v S) message.Message
