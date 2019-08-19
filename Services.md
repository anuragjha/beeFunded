# To implement applications on blockchain as service

1. Service interface
``` 
Interface Service {

  Service
  
  toServe()
  served()
}
```
- Service is a state
- toServe() method will take in a Service 
- served() will execute the request and return the result

#
2. Service structure
```
Service {
  serviceCode - helps to identify type of service
  serviceObj - contains the service requested
  
}
```
- serviceCode and serviceObj are state

#
For example -  Funding service 
