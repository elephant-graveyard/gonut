import java.net.*;
import java.io.*;
import java.util.*;
import java.util.regex.*;

public class App{

    public static void main(String args[]) {
        try{
            ServerSocket server = new ServerSocket(8080);
            while (true) {
                try{
                    Socket client = server.accept();

                    InputStream in = client.getInputStream();
                    Scanner s = new Scanner(in, "UTF-8");
                    OutputStream out = new BufferedOutputStream(client.getOutputStream());

                    if(!s.hasNext()){
                        continue;
                    }

                    String data = s.useDelimiter("\\r\\n\\r\\n").next();
                    Matcher get = Pattern.compile("^GET").matcher(data);

                    if(get.find()){
                        String message = "Hello, Homeport!";

                        byte[] response = ("HTTP/1.0 200 OK\r\n" +
                            "Content-Type: text/plain\r\n" +
                            "Date: " + new Date() + "\r\n" +
                            "Content-length: " + message.length() + "\r\n\r\n"+
                            message).getBytes("UTF-8");
                        out.write(response, 0, response.length);

                    }else{
                        out.write(("HTTP/1.0 400 Bad Request\r\n\r\n").getBytes("UTF-8"));
                    }

                    out.close();
                    s.close();

                }catch (Throwable tri) {
                    System.err.println("Error handling request: " + tri);
                }
            }
        }catch (Throwable tr) {
            System.err.println("Could not start server: " + tr);
        }
    }
}