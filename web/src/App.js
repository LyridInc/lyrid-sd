import React, { useState, useEffect } from 'react'
import ExporterTable from './tables/ExporterTable'
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
  const [lyridConnection, setLyridConnection] = useState({"status":"Checking Lyrid account ..."})
  const [time, setTime] = useState(Date.now())
  const exportersData = []
  const [exporters, setExporters] = useState(exportersData)
  
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
        setLyridConnection({"status":"Checking Lyrid account ..."})
        checkLyridConnection()
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
  
  const checkLyridConnection = () => {
    fetch(ROOT_URL+"/status")
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        setLyridConnection({status: "Connected to Lyrid under accout name " + result[0].name + "."})
    })
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
  
  useEffect(() => {
    const interval = setInterval(() => setTime(Date.now()), 60000)
    
    fetch(ROOT_URL+"/exporters")
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then(
      (result) => {
        //console.log(result)
        const keys = Object.keys(result)
        let eps = [];
        for (const key of keys) {
          eps.push(result[key])
        }
        setExporters(eps)
      },
      (error) => {
        console.log(error)
      }
    )
    
    checkLyridConnection()
    
    return () => {
        clearInterval(interval);
      }
  }, [time])
  return (
    <div className="container">
      <h1>Lyrid Service Discovery Configuration</h1>
      <form
        onSubmit={(event) => {
          event.preventDefault()
          updateConfiguration()
        }}
      >
      <label className="switch">
        <input type="checkbox" checked={configuration.Is_Local} onChange={toggleLocal} />
        <div className="slider"></div>
      </label>
      { !configuration.Is_Local ? (
      <div>
      <small>{lyridConnection.status}</small>
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
      <button>Save</button>
      </form>
      
      <div className="flex-large">
        <h2>List exporters 
          <button
            onClick={() => setTime(Date.now())}
            className="button muted-button"
          >
            Refresh
          </button>
        </h2>
        <ExporterTable exporters={exporters}/>
      </div>
      
    </div>
  )
}

export default App
