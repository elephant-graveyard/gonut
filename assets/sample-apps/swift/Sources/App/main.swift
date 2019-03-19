import Kitura
import Foundation

// Disable Buffering to write directly to stdout
setbuf(stdout, nil)

let endpoint = Router()
let logmessage = "Hello, Homeport!"
endpoint.get("/"){
    request, response, next in
    print(logmessage)
    response.send("\(logmessage)")
    next()
}

Kitura.addHTTPServer(onPort: 8080, with: endpoint)
Kitura.run()
