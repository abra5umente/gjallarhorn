import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { useServices } from '../context/ServiceContext'
import StatusBadge from './StatusBadge'
import LoadingSpinner from './LoadingSpinner'
import BulkActionBar from './BulkActionBar'
import BulkEditModal from './BulkEditModal'

const ServiceList = () => {
  const {
    services,
    loading,
    error,
    deleteService,
    selectedIds,
    toggleSelection,
    selectAll,
    clearSelection,
    bulkDeleteServices,
  } = useServices()

  const [showBulkEdit, setShowBulkEdit] = useState(false)

  const handleDelete = async (id, name) => {
    if (window.confirm(`Are you sure you want to delete "${name}"?`)) {
      try {
        await deleteService(id)
      } catch (err) {
        alert(`Failed to delete service: ${err.message}`)
      }
    }
  }

  const handleBulkDelete = async () => {
    if (window.confirm(`Are you sure you want to delete ${selectedIds.size} services?`)) {
      try {
        await bulkDeleteServices(Array.from(selectedIds))
      } catch (err) {
        alert(`Failed to delete services: ${err.message}`)
      }
    }
  }

  const formatDate = (dateString) => {
    if (!dateString) return 'Never'
    return new Date(dateString).toLocaleString()
  }

  if (loading) {
    return <LoadingSpinner />
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <div className="text-red-600 mb-4">
          <svg className="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 19.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">Error loading services</h3>
        <p className="text-gray-600">{error}</p>
      </div>
    )
  }

  if (services.length === 0) {
    return (
      <div className="text-center py-12">
        <div className="text-gray-400 mb-4">
          <svg className="mx-auto h-12 w-12" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h3 className="text-lg font-medium text-gray-900 mb-2">No services monitored</h3>
        <p className="text-gray-600 mb-6">Get started by adding your first service to monitor.</p>
        <Link to="/add" className="btn-primary">
          Add Service
        </Link>
      </div>
    )
  }

  const allSelected = services.length > 0 && selectedIds.size === services.length

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div className="flex items-center space-x-4">
          <h1 className="text-2xl font-bold text-gray-900">Monitored Services</h1>
          {services.length > 0 && (
            <button
              onClick={allSelected ? clearSelection : selectAll}
              className="text-sm text-primary-600 hover:text-primary-700"
            >
              {allSelected ? 'Deselect All' : 'Select All'}
            </button>
          )}
        </div>
        <Link to="/add" className="btn-primary">
          Add Service
        </Link>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {services.map((service) => (
          <div
            key={service.id}
            className={`card relative ${selectedIds.has(service.id) ? 'ring-2 ring-primary-500' : ''}`}
          >
            <div className="absolute top-4 left-4">
              <input
                type="checkbox"
                checked={selectedIds.has(service.id)}
                onChange={() => toggleSelection(service.id)}
                className="w-4 h-4 rounded border-gray-300 text-primary-600 focus:ring-primary-500"
              />
            </div>
            <div className="flex justify-between items-start mb-4 ml-6">
              <h3 className="text-lg font-semibold text-gray-900">{service.name}</h3>
              <StatusBadge status={service.status} />
            </div>
            
            <div className="space-y-2 mb-4">
              <div className="flex items-center text-sm text-gray-600">
                <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                </svg>
                <a 
                  href={service.url} 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="text-primary-600 hover:text-primary-700 truncate"
                >
                  {service.url}
                </a>
              </div>
              
              <div className="flex items-center text-sm text-gray-600">
                <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Check every {service.interval}s
              </div>
              
              <div className="flex items-center text-sm text-gray-600">
                <svg className="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Last checked: {formatDate(service.lastChecked)}
              </div>
            </div>
            
            <div className="flex space-x-2">
              <Link
                to={`/edit/${service.id}`}
                className="flex-1 btn-secondary text-center"
              >
                Edit
              </Link>
              <button
                onClick={() => handleDelete(service.id, service.name)}
                className="flex-1 btn-danger"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>

      <BulkActionBar
        onEdit={() => setShowBulkEdit(true)}
        onDelete={handleBulkDelete}
      />

      <BulkEditModal
        isOpen={showBulkEdit}
        onClose={() => setShowBulkEdit(false)}
      />
    </div>
  )
}

export default ServiceList
