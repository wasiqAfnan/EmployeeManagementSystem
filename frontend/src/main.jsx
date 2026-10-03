import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { Amplify } from 'aws-amplify';
import './index.css'
import App from './App.jsx'

// Extract user pool ID from authority
const authority = (import.meta.env.VITE_COGNITO_AUTHORITY || '').trim();
const userPoolId = authority.split('/').pop();

// Extract domain without https://
const domainUrl = (import.meta.env.VITE_COGNITO_DOMAIN || '').trim();
const domain = domainUrl.replace('https://', '').replace('http://', '');

Amplify.configure({
  Auth: {
    Cognito: {
      userPoolId: userPoolId,
      userPoolClientId: (import.meta.env.VITE_COGNITO_CLIENT_ID || '').trim(),
      loginWith: {
        oauth: {
          domain: domain,
          scopes: ['email', 'openid', 'profile'],
          redirectSignIn: [window.location.origin],
          redirectSignOut: [window.location.origin],
          responseType: 'code'
        }
      }
    }
  }
});

createRoot(document.getElementById('root')).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
