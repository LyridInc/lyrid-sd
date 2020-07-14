import React, { useState, useEffect } from 'react'

const App = () => {

  const ROOT_URL = 'http://localhost:8000';
  const initialState = {
    Bind_Address:  	  		"127.0.0.1",
    Mngt_Port: 				9000,
    Discovery_Port_Start: 	9001,
    Max_Discovery: 			1024,
    Discovery_Poll_Interval: 	"10s",
    Discovery_Interface: 		"127.0.0.1"
  }
  const [configuration, setConfiguration] = useState(initialState)
  
  const updateConfiguration = () => {
    const requestOptions = {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(configuration)
    };
    fetch(ROOT_URL+'/config', requestOptions)
    .then(res => res.json())
    .then(
      (result) => {
        console.log(result)
        setConfiguration(result)
      },
      (error) => {
        console.log(error)
      }
    )
  }
  
  const handleInputChange = (event) => {
    const { name, value } = event.target
    setConfiguration({ ...configuration, [name]: value })
  }
  
  useEffect(() => {
    fetch(ROOT_URL+"/config")
    .then(res => res.json())
    .then(
      (result) => {
        console.log(result)
        setConfiguration(result)
      },
      (error) => {
        console.log(error)
      }
    )
  }, [])
  
  return (
    <div className="container">
      <h1>Lyrid Service Discovery Configuration</h1>
      <form
      onSubmit={(event) => {
        event.preventDefault()
        updateConfiguration()
      }}
    >
      <label>Bind Address</label>
      <input
        type="text"
        name="Bind_Address"
        value={configuration.Bind_Address}
        onChange={handleInputChange}
      />
      <label>Mngt Port</label>
      <input
        type="text"
        name="Mngt_Port"
        value={configuration.Mngt_Port}
        onChange={handleInputChange}
      />
      <label>Discovery Port Start</label>
      <input
        type="text"
        name="Discovery_Port_Start"
        value={configuration.Discovery_Port_Start}
        onChange={handleInputChange}
      />
      <label>Max Discovery</label>
      <input
        type="text"
        name="Max_Discovery"
        value={configuration.Max_Discovery}
        onChange={handleInputChange}
      />
      <label>Discovery Poll_Interval</label>
      <input
        type="text"
        name="Discovery_Poll_Interval"
        value={configuration.Discovery_Poll_Interval}
        onChange={handleInputChange}
      />
      <label>Discovery Interface</label>
      <input
        type="text"
        name="Discovery_Interface"
        value={configuration.Discovery_Interface}
        onChange={handleInputChange}
      />
      <button>Save</button>
      </form>
    </div>
  )
}

export default App
