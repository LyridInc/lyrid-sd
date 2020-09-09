import React, { useState, useEffect } from 'react'
import ExporterTable from './tables/ExporterTable'
import GatewayTable from './tables/GatewayTable'
import Select from 'react-select'
const App = () => {

  const ROOT_URL = '';
  const initialState = {
    Discovery_Port_Start: 	9001,
    Max_Discovery: 			1024,
    Discovery_Poll_Interval: 	"10s",
    Scrape_Valid_Timeout:       "5m",
    Lyrid_Key:                  "",
    Lyrid_Secret:               "",
    Local_Serverless_Url:       "http://localhost:8080",
    Is_Local:                   true,
    Noc_App_Name:               ""
  }
  const [configuration, setConfiguration] = useState(initialState)
  const [lyridConnection, setLyridConnection] = useState({"status":"Checking Lyrid account ..."})
  const [time, setTime] = useState(Date.now())
  const exportersData = []
  const [exporters, setExporters] = useState(exportersData)
  const [gateways, setGateways] = useState([])
  
  const [apps, setApps] = useState([])
  const options = [
  { value: 'chocolate', label: 'Chocolate' },
  { value: 'strawberry', label: 'Strawberry' },
  { value: 'vanilla', label: 'Vanilla' }
]
 const [selectedApp, setSelectedApp] = useState({value: configuration.Noc_App_Name, label: configuration.Noc_App_Name})
  
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
  
  const handleAppSelect = (app) => {
    setSelectedApp(app)
    setConfiguration({ ...configuration, ["Noc_App_Name"]: app.value })
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
        listApps()
    })
  }
  
  const listApps = () => {
    fetch(ROOT_URL+"/apps")
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then((result) => {
        let apps = [];
        result.forEach((element) => {
          apps.push({ value: element, label: element })
        })
        setApps(apps)
        console.log(apps)
    })
  }
  
  const deleteExporter = (id) => {
    const requestOptions = {
      method: 'DELETE'
    };
    fetch(ROOT_URL+'/exporter/delete/'+id, requestOptions)
    setExporters(exporters.filter((exporter) => exporter.ID !== id))
  }
  
  const deleteGateway = (id) => {
    const requestOptions = {
      method: 'DELETE'
    };
    fetch(ROOT_URL+'/gateway/delete/'+id, requestOptions)
    setGateways(gateways.filter((gateway) => gateway.ID !== id))
  }
  
  useEffect(() => {
    fetch(ROOT_URL+"/config")
    .then(res => res.json())
    .then(
      (result) => {
        console.log(result)
        setConfiguration(result)
        setSelectedApp({ value: result.Noc_App_Name, label: result.Noc_App_Name })
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
    
    fetch(ROOT_URL+"/gateways")
    .then((res) => {
        if(!res.ok) { throw new Error(res.status) } else {return res.json()}
    })
    .then(
      (result) => {
        setGateways(result)
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
      <label>NOC App</label>
      <Select options={apps} value={selectedApp} onChange={handleAppSelect}/>
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
      <label>Scrape Result Valid Timeout</label>
      <input
        type="text"
        name="Scrape_Valid_Timeout"
        value={configuration.Scrape_Valid_Timeout}
        onChange={handleInputChange}
      />
      <button>Save</button>
      </form>
      <div className="flex-large">
        <h2>List gateways 
          <button
            onClick={() => setTime(Date.now())}
            className="button muted-button"
          >
            Refresh
          </button>
        </h2>
        <GatewayTable gateways={gateways} deleteGateway={deleteGateway}/>
      </div>
      
      <div className="flex-large">
        <h2>List exporters 
          <button
            onClick={() => setTime(Date.now())}
            className="button muted-button"
          >
            Refresh
          </button>
        </h2>
        <ExporterTable exporters={exporters} deleteExporter={deleteExporter}/>
      </div>
      
    </div>
  )
}

export default App
