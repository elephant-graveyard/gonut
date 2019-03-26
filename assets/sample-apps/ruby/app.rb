Bundler.require :web

get '/' do
  [200, {'Content-Type' => 'text/plain'}, ["Hello, Homeport!\n"]]
end
