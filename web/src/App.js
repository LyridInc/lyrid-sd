import React, { useState, useEffect } from 'react'

const App = () => {

  const ROOT_URL = '';
  const initialState = {
    Discovery_Port_Start: 	9001,
    Max_Discovery: 			1024,
    Discovery_Poll_Interval: 	"10s",
    Lyrid_Key:                  "",
    Lyrid_Secret:               "",
    Local_Serverless_Url:       "http://localhost:8080",
    Is_Local:                   true
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
    if (name == "Discovery_Port_Start" || name == "Max_Discovery") {
        setConfiguration({ ...configuration, [name]: parseInt(value) })
    } else {
        setConfiguration({ ...configuration, [name]: value })
    }
  }
  
  const toggleLocal = () => {
    setConfiguration({ ...configuration, ["Is_Local"]: !configuration.Is_Local })
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
      <label>Is Local</label>
      <label className="switch">
        <input type="checkbox" checked={configuration.Is_Local} onChange={toggleLocal} />
        <div className="slider"></div>
      </label>
      { !configuration.Is_Local ? (
      <div>
      <label>Lyrid Key</label>
      <input
        type="text"
        name="Lyrid_Key"
        value={configuration.Lyrid_Key}
        onChange={handleInputChange}
      />
      <label>Lyrid Secret</label>
      <input
        type="text"
        name="Lyrid_Secret"
        value={configuration.Lyrid_Secret}
        onChange={handleInputChange}
      />
      </div>
      ) : (
      <div>
      <label>Local_Serverless_Url</label>
      <input
        type="text"
        name="Local_Serverless_Url"
        value={configuration.Local_Serverless_Url}
        onChange={handleInputChange}
      />
      </div>
      )}
      <button>Save</button>
      </form>
    </div>
  )
}

export default App
