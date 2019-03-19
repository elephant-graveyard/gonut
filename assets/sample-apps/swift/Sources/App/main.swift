import Kitura
import Foundation

// Disable Buffering to write directly to stdout
setbuf(stdout, nil)
// simplest log writer ever
print("Gentleman start the engine ....")
print("This is a Swift based app!")

let endpoint = Router()
let logmessage = "Hello. I'm Johnny 5."
endpoint.get("/"){
    request, response, next in
    print(logmessage)
    response.send("\(logmessage)")
    next()
}

Kitura.addHTTPServer(onPort: 8080, with: endpoint)
Kitura.run()
