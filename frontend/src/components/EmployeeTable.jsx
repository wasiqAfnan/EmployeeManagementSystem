export default function EmployeeTable({ employees, onEdit, onDelete, loading }) {
  if (loading) {
    return (
      <div className="flex justify-center items-center py-20">
        <svg className="w-6 h-6 animate-spin text-blue-500" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
        </svg>
        <span className="ml-2 text-gray-400 text-sm">Loading employees...</span>
      </div>
    );
  }

  if (employees.length === 0) {
    return (
      <div className="text-center py-20">
        <svg className="w-12 h-12 text-gray-200 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M17 20h5v-2a4 4 0 00-5.197-3.797M9 20H4v-2a4 4 0 015.197-3.797m0 0A4 4 0 1112 8a4 4 0 012.197 6.203m0 0A4 4 0 1117 20" />
        </svg>
        <p className="text-gray-400 text-sm">No employees found</p>
      </div>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="bg-gray-50 border-b border-gray-100">
            {['Emp ID', 'Full Name', 'Job Title', 'Department', 'Salary', 'Actions'].map((h) => (
              <th
                key={h}
                className="px-4 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wide"
              >
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-50">
          {employees.map((emp) => (
            <tr key={emp.empId} className="hover:bg-gray-50 transition-colors">
              <td className="px-4 py-3">
                <span className="inline-block bg-blue-50 text-blue-700 text-xs font-medium px-2 py-0.5 rounded">
                  {emp.empId}
                </span>
              </td>
              <td className="px-4 py-3 font-medium text-gray-800">{emp.fullName}</td>
              <td className="px-4 py-3 text-gray-600">{emp.jobTitle}</td>
              <td className="px-4 py-3">
                <span className="inline-block bg-gray-100 text-gray-600 text-xs px-2 py-0.5 rounded">
                  {emp.department}
                </span>
              </td>
              <td className="px-4 py-3 text-gray-700 font-medium">
                INR {Number(emp.salary).toLocaleString()}
              </td>
              <td className="px-4 py-3">
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => onEdit(emp)}
                    className="text-xs px-3 py-1.5 rounded-md border border-gray-200 text-gray-600 hover:border-blue-300 hover:text-blue-600 transition-colors"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => onDelete(emp)}
                    className="text-xs px-3 py-1.5 rounded-md border border-gray-200 text-gray-600 hover:border-red-300 hover:text-red-600 transition-colors"
                  >
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
