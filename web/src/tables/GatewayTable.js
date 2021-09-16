import React from 'react'

const GatewayTable = (props) => (
  <table>
    <thead>
      <tr>
        <th>ID</th>
        <th>Host Name</th>
        <th>Actions</th>
      </tr>
    </thead>
    <tbody>
      {props.gateways.length > 0 ? (
        props.gateways.map((gateway) => (
          <tr key={gateway.ID}>
            <td>{gateway.ID}</td>
            <td>{gateway.Hostname}</td>
            <td>
              <button
                onClick={() => props.deleteGateway(gateway.ID)}
                className="button muted-button"
              >
                Delete
              </button>
            </td>
          </tr>
        ))
      ) : (
        <tr>
          <td colSpan={4}>No gateways</td>
        </tr>
      )}
    </tbody>
  </table>
)

export default GatewayTable
