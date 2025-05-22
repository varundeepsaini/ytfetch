import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

const getTheme = (mode) => createTheme({
  palette: {
    mode,
    ...(mode === 'dark' ? {
      primary: {
        main: '#90caf9',
      },
      secondary: {
        main: '#f48fb1',
      },
      background: {
        default: '#121212',
        paper: '#1e1e1e',
      },
    } : {
      primary: {
        main: '#1976d2',
      },
      secondary: {
        main: '#dc004e',
      },
    }),
  },
});

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <ThemeProvider theme={getTheme('light')}>
      <CssBaseline />
      <App />
    </ThemeProvider>
  </React.StrictMode>
); 