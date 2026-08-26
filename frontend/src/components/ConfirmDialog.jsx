import Modal from './Modal';

export default function ConfirmDialog({ employee, onConfirm, onCancel, loading }) {
  return (
    <Modal title="Delete Employee" onClose={onCancel}>
      <p className="text-sm text-gray-600 mb-1">
        Are you sure you want to delete{' '}
        <span className="font-semibold text-gray-800">{employee.fullName}</span>?
      </p>
      <p className="text-xs text-gray-400 mb-6">
        Employee ID: {employee.empId} · This action cannot be undone.
      </p>
      <div className="flex justify-end gap-2">
        <button
          onClick={onCancel}
          disabled={loading}
          className="px-4 py-2 text-sm rounded-lg border border-gray-200 text-gray-600 hover:bg-gray-50 transition-colors disabled:opacity-50"
        >
          Cancel
        </button>
        <button
          onClick={onConfirm}
          disabled={loading}
          className="px-4 py-2 text-sm rounded-lg bg-red-500 text-white hover:bg-red-600 transition-colors disabled:opacity-60 flex items-center gap-2"
        >
          {loading && (
            <svg className="w-3.5 h-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
            </svg>
          )}
          Delete
        </button>
      </div>
    </Modal>
  );
}
