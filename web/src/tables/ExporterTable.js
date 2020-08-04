import React from 'react'

const ExporterTable = (props) => (
  <table>
    <thead>
      <tr>
        <th>ID</th>
        <th>Port</th>
        <th>Metric Endpoint</th>
        <th>Original Endpoint</th>
      </tr>
    </thead>
    <tbody>
      {props.exporters.length > 0 ? (
        props.exporters.map((exporter) => (
          <tr key={exporter.id}>
            <td>{exporter.ID}</td>
            <td>{exporter.Port}</td>
            <td><a target="_blank" href={exporter.MetricEndpoint}>{exporter.MetricEndpoint}</a></td>
            <td>{exporter.URL}</td>
          </tr>
        ))
      ) : (
        <tr>
          <td colSpan={4}>No exporters</td>
        </tr>
      )}
    </tbody>
  </table>
)

export default ExporterTable
