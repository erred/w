--- title
service ui.md
--- description
different ways to present a service
--- main


It's been 4 months
and I've already refactored a service 4 times

# _server side_ templating

pros:

- one thing to manage
- internally consistent
- can be pretty (if you can bear writing css)

cons:

- tightly bound
- write html and css

# client side rendering _json api_

pros:

- separating front and back end
- can be pretty (if you can bear writing css)

cons:

- **javascript**, so heavy
- json
- consistency is hard
- html, css, and javascript

# client side rendering _grpc-web api_

pros:

- separating front and back end
- versioned, backward compatible apis
- can be pretty (if you can bear writing css)

cons:

- **javascript**, so heavy
- grpc-web is still weird wrt load balancing / routing / caching
- html, css and moer javascript

# chatbot _Dialogflow_

pros:

- backend only, no need to write frontend
- auth by default
- magic processing
- multiplatform

cons:

- magic processing
- least common denominator of platforms
- no push
- ugly

# chatbot _Telegram_

pros:

- backend only, no need to write frontend
- auth by default
- push messages

cons:

- ugly
