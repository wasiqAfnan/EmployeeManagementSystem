export default function Navbar() {
  return (
    <nav className="bg-white border-b border-gray-200 px-6 py-4 flex items-center gap-3">
      <div className="bg-blue-600 text-white rounded-lg w-8 h-8 flex items-center justify-center font-bold text-sm">
        EMS
      </div>
      <div>
        <h1 className="text-gray-900 font-semibold text-lg leading-tight">
          Employee Management
        </h1>
        <p className="text-gray-400 text-xs">Manage your team records</p>
      </div>
    </nav>
  );
}
