import { useState, useEffect } from "react";
import { signInWithRedirect, signOut, getCurrentUser, fetchAuthSession } from 'aws-amplify/auth';
import { Hub } from 'aws-amplify/utils';
import Dashboard from "./Dashboard";
import api from './utils/api';

function App() {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    checkUser();

    const unsubscribe = Hub.listen('auth', ({ payload }) => {
      switch (payload.event) {
        case 'signInWithRedirect':
          checkUser();
          break;
        case 'signInWithRedirect_failure':
          setError('An error has occurred during the OAuth flow.');
          break;
        case 'signedOut':
          setUser(null);
          break;
      }
    });

    return unsubscribe;
  }, []);

  async function checkUser() {
    try {
      setIsLoading(true);
      await getCurrentUser(); // Ensure user is signed in locally

      const session = await fetchAuthSession();
      const idToken = session.tokens?.idToken?.toString();

      // Call backend /me endpoint and pass ID token in X-Id-Token header
      // The Axios interceptor automatically attaches the access token
      const response = await api.get('/me', {
        headers: idToken ? { 'X-Id-Token': idToken } : {}
      });

      const userData = response.data?.data;
      if (!userData) {
        throw new Error('No user data returned from /me');
      }

      setUser(userData);
      setError(null);
    } catch (e) {
      console.error('User check failed', e);
      setUser(null);
    } finally {
      setIsLoading(false);
    }
  }

  const handleSignIn = async () => {
    try {
      await signInWithRedirect();
    } catch (err) {
      console.error(err);
    }
  };

  const handleSignOut = async () => {
    try {
      await signOut();
    } catch (err) {
      console.error(err);
    }
  };

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <p className="text-gray-500">Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <p className="text-red-500">Encountering error... {error}</p>
      </div>
    );
  }

  if (user) {
    const mockAuth = {
      user: {
        profile: {
          email: user.email || user.username || "",
          role: user.role || "user",
        }
      }
    };
    
    return (
      <Dashboard auth={mockAuth} onSignOut={handleSignOut} />
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="w-full max-w-md p-8 space-y-8 bg-white shadow-lg rounded-xl text-center border border-gray-100">
        <div>
          <div className="bg-blue-600 text-white rounded-xl w-16 h-16 flex items-center justify-center font-bold text-2xl mx-auto shadow-sm">
            EMS
          </div>
          <h2 className="mt-6 text-3xl font-extrabold text-gray-900">Welcome to EMS</h2>
          <p className="mt-2 text-sm text-gray-500">Employee Management System</p>
        </div>
        <button 
          onClick={handleSignIn}
          className="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none cursor-pointer transition-colors"
        >
          Log in or Sign up with Amazon Cognito
        </button>
      </div>
    </div>
  );
}

export default App;
