const BASE_URL = '/employee';
const BASE_URL_PLURAL = '/employees';

async function handleResponse(res) {
  const json = await res.json();
  if (!res.ok) {
    throw new Error(json.message || 'Something went wrong');
  }
  return json;
}

export async function getAllEmployees() {
  const res = await fetch(BASE_URL_PLURAL);
  return handleResponse(res);
}

export async function getEmployee(empId) {
  const res = await fetch(`${BASE_URL}/${empId}`);
  return handleResponse(res);
}

export async function createEmployee(data) {
  const res = await fetch(BASE_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  return handleResponse(res);
}

export async function updateEmployee(empId, data) {
  const res = await fetch(`${BASE_URL}/${empId}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  return handleResponse(res);
}

export async function deleteEmployee(empId) {
  const res = await fetch(`${BASE_URL}/${empId}`, {
    method: 'DELETE',
  });
  return handleResponse(res);
}

export async function searchEmployees(query) {
  const res = await fetch(`${BASE_URL_PLURAL}/search?q=${encodeURIComponent(query)}`);
  return handleResponse(res);
}
