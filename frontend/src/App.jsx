import { useCallback, useEffect, useRef, useState } from 'react';
import {
  createEmployee,
  deleteEmployee,
  getAllEmployees,
  searchEmployees,
  updateEmployee,
} from './api/employeeApi';
import ConfirmDialog from './components/ConfirmDialog';
import EmployeeForm from './components/EmployeeForm';
import EmployeeTable from './components/EmployeeTable';
import Modal from './components/Modal';
import Navbar from './components/Navbar';
import Toast from './components/Toast';

export default function App() {
  const [employees, setEmployees] = useState([]);
  const [tableLoading, setTableLoading] = useState(true);
  const [formLoading, setFormLoading] = useState(false);

  // Modal state: null | 'create' | 'edit' | 'delete'
  const [modal, setModal] = useState(null);
  const [selectedEmployee, setSelectedEmployee] = useState(null);

  // Search
  const [searchQuery, setSearchQuery] = useState('');
  const [isSearching, setIsSearching] = useState(false);
  const searchTimeout = useRef(null);

  // Toast
  const [toast, setToast] = useState(null); // { message, type }

  function showToast(message, type = 'success') {
    setToast({ message, type });
  }

  // Fetch all employees
  const fetchAll = useCallback(async () => {
    setTableLoading(true);
    try {
      const res = await getAllEmployees();
      setEmployees(res.data || []);
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setTableLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchAll();
  }, [fetchAll]);

  // Live search with 400ms debounce
  useEffect(() => {
    clearTimeout(searchTimeout.current);

    if (!searchQuery.trim()) {
      fetchAll();
      return;
    }

    searchTimeout.current = setTimeout(async () => {
      setIsSearching(true);
      setTableLoading(true);
      try {
        const res = await searchEmployees(searchQuery.trim());
        setEmployees(res.data || []);
      } catch (err) {
        showToast(err.message, 'error');
      } finally {
        setTableLoading(false);
        setIsSearching(false);
      }
    }, 400);

    return () => clearTimeout(searchTimeout.current);
  }, [searchQuery, fetchAll]);

  // Create
  async function handleCreate(data) {
    setFormLoading(true);
    try {
      await createEmployee(data);
      showToast('Employee created successfully');
      setModal(null);
      fetchAll();
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setFormLoading(false);
    }
  }

  // Edit
  function openEdit(emp) {
    setSelectedEmployee(emp);
    setModal('edit');
  }

  async function handleUpdate(data) {
    setFormLoading(true);
    try {
      await updateEmployee(selectedEmployee.empId, data);
      showToast('Employee updated successfully');
      setModal(null);
      setSelectedEmployee(null);
      fetchAll();
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setFormLoading(false);
    }
  }

  // Delete
  function openDelete(emp) {
    setSelectedEmployee(emp);
    setModal('delete');
  }

  async function handleDelete() {
    setFormLoading(true);
    try {
      await deleteEmployee(selectedEmployee.empId);
      showToast('Employee deleted successfully');
      setModal(null);
      setSelectedEmployee(null);
      fetchAll();
    } catch (err) {
      showToast(err.message, 'error');
    } finally {
      setFormLoading(false);
    }
  }

  function closeModal() {
    if (formLoading) return;
    setModal(null);
    setSelectedEmployee(null);
  }

  // Stats
  const totalSalary = employees.reduce((sum, e) => sum + Number(e.salary || 0), 0);
  const departments = new Set(employees.map((e) => e.department)).size;

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />

      <main className="max-w-6xl mx-auto px-4 py-8">

        {/* ── Stats row ── */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-8">
          <StatCard label="Total Employees" value={employees.length} icon="👥" />
          <StatCard label="Departments" value={departments} icon="🏢" />
          <StatCard
            label="Total Salary"
            value={`INR ${totalSalary.toLocaleString()}`}
            icon="💰"
          />
        </div>

        {/* ── Toolbar ── */}
        <div className="bg-white rounded-xl border border-gray-100 shadow-sm">
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 px-5 py-4 border-b border-gray-100">
            <div className="relative flex-1 max-w-sm">
              <svg
                className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                strokeWidth={2}
              >
                <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-4.35-4.35M17 11A6 6 0 111 11a6 6 0 0116 0z" />
              </svg>
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder="Search by name, ID, title, department…"
                className="w-full pl-9 pr-4 py-2 text-sm rounded-lg border border-gray-200 bg-gray-50 focus:outline-none focus:border-blue-400 focus:bg-white transition-colors"
              />
              {isSearching && (
                <svg className="absolute right-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 animate-spin text-blue-400" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z" />
                </svg>
              )}
            </div>

            <button
              onClick={() => setModal('create')}
              className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors font-medium shrink-0"
            >
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2.5}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M12 4v16m8-8H4" />
              </svg>
              Add Employee
            </button>
          </div>

          {/* ── Table ── */}
          <EmployeeTable
            employees={employees}
            loading={tableLoading}
            onEdit={openEdit}
            onDelete={openDelete}
          />

          {/* ── Footer count ── */}
          {!tableLoading && employees.length > 0 && (
            <div className="px-5 py-3 border-t border-gray-50 text-xs text-gray-400">
              {searchQuery
                ? `${employees.length} result${employees.length !== 1 ? 's' : ''} for "${searchQuery}"`
                : `${employees.length} employee${employees.length !== 1 ? 's' : ''} total`}
            </div>
          )}
        </div>
      </main>

      {/* ── Create Modal ── */}
      {modal === 'create' && (
        <Modal title="Add New Employee" onClose={closeModal}>
          <EmployeeForm
            onSubmit={handleCreate}
            onCancel={closeModal}
            loading={formLoading}
          />
        </Modal>
      )}

      {/* ── Edit Modal ── */}
      {modal === 'edit' && selectedEmployee && (
        <Modal title="Edit Employee" onClose={closeModal}>
          <EmployeeForm
            initial={selectedEmployee}
            onSubmit={handleUpdate}
            onCancel={closeModal}
            isEdit
            loading={formLoading}
          />
        </Modal>
      )}

      {/* ── Delete Confirm ── */}
      {modal === 'delete' && selectedEmployee && (
        <ConfirmDialog
          employee={selectedEmployee}
          onConfirm={handleDelete}
          onCancel={closeModal}
          loading={formLoading}
        />
      )}

      {/* ── Toast ── */}
      {toast && (
        <Toast
          message={toast.message}
          type={toast.type}
          onClose={() => setToast(null)}
        />
      )}
    </div>
  );
}

function StatCard({ label, value, icon }) {
  return (
    <div className="bg-white rounded-xl border border-gray-100 shadow-sm px-5 py-4 flex items-center gap-4">
      <span className="text-2xl">{icon}</span>
      <div>
        <p className="text-2xl font-bold text-gray-800">{value}</p>
        <p className="text-xs text-gray-400 mt-0.5">{label}</p>
      </div>
    </div>
  );
}
