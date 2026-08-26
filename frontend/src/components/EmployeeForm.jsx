import { useState } from 'react';

const EMPTY_FORM = {
  empId: '',
  fullName: '',
  jobTitle: '',
  department: '',
  salary: '',
};

export default function EmployeeForm({ initial = {}, onSubmit, onCancel, isEdit = false, loading = false }) {
  const [form, setForm] = useState({ ...EMPTY_FORM, ...initial });
  const [errors, setErrors] = useState({});

  function validate() {
    const e = {};
    if (!form.empId.trim()) e.empId = 'Employee ID is required';
    if (!form.fullName.trim()) e.fullName = 'Full name is required';
    if (!form.jobTitle.trim()) e.jobTitle = 'Job title is required';
    if (!form.department.trim()) e.department = 'Department is required';
    if (!form.salary || Number(form.salary) <= 0) e.salary = 'Salary must be greater than 0';
    return e;
  }

  function handleChange(e) {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
    if (errors[name]) setErrors((prev) => ({ ...prev, [name]: '' }));
  }

  function handleSubmit(e) {
    e.preventDefault();
    const errs = validate();
    if (Object.keys(errs).length > 0) { setErrors(errs); return; }
    onSubmit({ ...form, salary: Number(form.salary) });
  }

  const fields = [
    { name: 'empId', label: 'Employee ID', placeholder: 'e.g. EMP001', disabled: isEdit },
    { name: 'fullName', label: 'Full Name', placeholder: 'e.g. John Doe' },
    { name: 'jobTitle', label: 'Job Title', placeholder: 'e.g. Software Engineer' },
    { name: 'department', label: 'Department', placeholder: 'e.g. Engineering' },
    { name: 'salary', label: 'Salary (INR)', placeholder: 'e.g. 85000', type: 'number' },
  ];

  return (
    <form onSubmit={handleSubmit} noValidate>
      <div className="flex flex-col gap-4">
        {fields.map(({ name, label, placeholder, disabled, type = 'text' }) => (
          <div key={name}>
            <label className="block text-sm font-medium text-gray-700 mb-1">{label}</label>
            <input
              type={type}
              name={name}
              value={form[name]}
              onChange={handleChange}
              placeholder={placeholder}
              disabled={disabled || loading}
              className={`w-full px-3 py-2 text-sm rounded-lg border outline-none transition-colors
                ${errors[name]
                  ? 'border-red-400 bg-red-50 focus:border-red-500'
                  : 'border-gray-200 bg-gray-50 focus:border-blue-400 focus:bg-white'}
                ${disabled ? 'bg-gray-100 text-gray-400 cursor-not-allowed' : ''}
              `}
            />
            {errors[name] && (
              <p className="mt-1 text-xs text-red-500">{errors[name]}</p>
            )}
          </div>
        ))}
      </div>

      <div className="flex justify-end gap-2 mt-6">
        <button
          type="button"
          onClick={onCancel}
          disabled={loading}
          className="px-4 py-2 text-sm rounded-lg border border-gray-200 text-gray-600 hover:bg-gray-50 transition-colors disabled:opacity-50"
        >
          Cancel
        </button>
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 text-sm rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors disabled:opacity-60 flex items-center gap-2"
        >
          {loading && (
            <svg className="w-3.5 h-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
            </svg>
          )}
          {isEdit ? 'Update Employee' : 'Add Employee'}
        </button>
      </div>
    </form>
  );
}
