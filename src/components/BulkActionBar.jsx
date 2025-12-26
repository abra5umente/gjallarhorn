import React from 'react'
import { useServices } from '../context/ServiceContext'

const BulkActionBar = ({ onEdit, onDelete }) => {
  const { selectedIds, clearSelection } = useServices()

  if (selectedIds.size === 0) return null

  return (
    <div className="fixed bottom-6 left-1/2 transform -translate-x-1/2 bg-white rounded-lg shadow-lg border border-gray-200 p-4 flex items-center space-x-4 z-50">
      <span className="text-gray-700 font-medium">
        {selectedIds.size} service{selectedIds.size > 1 ? 's' : ''} selected
      </span>
      <button onClick={onEdit} className="btn-secondary">
        Edit Selected
      </button>
      <button onClick={onDelete} className="btn-danger">
        Delete Selected
      </button>
      <button
        onClick={clearSelection}
        className="text-gray-500 hover:text-gray-700 px-3 py-2"
      >
        Cancel
      </button>
    </div>
  )
}

export default BulkActionBar
